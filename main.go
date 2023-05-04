package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

type ChildInfo struct {
	Child_id      int    `json:"child_id"`
	Name          string `json:"name"`
	Phone         string `json:"phone"`
	Year_of_birth string `json:"year_of_birth"`
}
type ClubInfo struct {
	Club_id   int    `json:"club_id"`
	Club_name string `json:"club_name"`
	Price     string `json:"price"`
}

type server struct {
	db *sql.DB
}

func dbConnect() server {
	db, err := sql.Open("sqlite3", "database.db")
	fmt.Println("Opening database")
	if err != nil {
		log.Fatal(err)
	}

	s := server{db: db}

	return s
}
func (s *server) selectUsers() []ChildInfo {
	rows, err := s.db.Query("select child_id, name, phone, year_of_birth from child_users;")
	if err != nil {
		log.Fatal(err)
	}

	var users []ChildInfo
	for rows.Next() {
		var user ChildInfo
		err := rows.Scan(&user.Child_id, &user.Name, &user.Phone, &user.Year_of_birth)
		if err != nil {
			log.Fatal("selectUsers", err)
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		log.Fatal("selectUsers2", err)
	}

	return users
}

func (s *server) selectClubs() []ClubInfo {
	rows, err := s.db.Query("select club_id, club_name, price from clubs")
	if err != nil {
		log.Fatal(err)
	}

	var clubs []ClubInfo
	for rows.Next() {
		var club ClubInfo
		err := rows.Scan(&club.Club_id, &club.Club_name, &club.Price)
		if err != nil {
			log.Fatal("selectClubs", err)
		}
		clubs = append(clubs, club)
	}
	if err = rows.Err(); err != nil {
		log.Fatal("selectClubs2", err)
	}

	return clubs
}

func (s *server) selectUser(id int) ChildInfo {
	rows := s.db.QueryRow("select child_id, name, phone,year_of_birth from child_users where child_id=?;", id)

	var user ChildInfo
	err := rows.Scan(&user.Child_id, &user.Name, &user.Phone, &user.Year_of_birth)
	if err != nil {
		log.Fatal("selectUsers", err)
	}

	return user
}

func (s *server) selectClub(index int) ClubInfo {
	rows := s.db.QueryRow("select club_id, club_name, price from clubs where club_id=?;", index)

	var club ClubInfo
	err := rows.Scan(&club.Club_id, &club.Club_name, &club.Price)
	if err != nil {
		log.Fatal("selectClubs", err)
	}

	return club
}

func (s *server) allUsersHandle(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./static/users.html")
	if err != nil {
		log.Fatal("allUsersHandle", err)
	}

	allUsers := s.selectUsers()
	errExecute := t.Execute(w, allUsers)
	fmt.Println(allUsers[0].Name)
	if errExecute != nil {
		log.Fatal("allUsersHandle2", err)
	}
}
func (s *server) allClubsHandle(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./static/clubs.html")
	if err != nil {
		log.Fatal("allClubsHandle", err)
	}

	allClubs := s.selectClubs()
	errExecute := t.Execute(w, allClubs)
	fmt.Println(allClubs[0].Club_name)
	if errExecute != nil {
		log.Fatal("allClubsHandle2", err)
	}
}

func (s *server) updateUserByID(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	id := r.FormValue("id")
	idInt, err := strconv.Atoi(id)
	Name := r.FormValue("name")
	Phone := r.FormValue("phone")
	Year_of_birth := r.FormValue("year_of_birth")

	updateUser(Name, Phone, Year_of_birth, idInt, s)
	http.Redirect(w, r, "/users", http.StatusSeeOther)
}

func (s *server) updateClubByID(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	index := r.FormValue("index")
	idInteger, err := strconv.Atoi(index)
	Club_name := r.FormValue("club_name")
	Price := r.FormValue("price")

	updateClub(Club_name, Price, idInteger, s)
	http.Redirect(w, r, "/clubs", http.StatusSeeOther)
}

func (s *server) updateUserForm(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./static/updUser.html")
	if err != nil {
		log.Fatal("allUsersHandle", err)
	}

	err = r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	id := r.FormValue("id")
	idInt, err := strconv.Atoi(id)
	user := s.selectUser(idInt)

	t.Execute(w, user)
}

func (s *server) updateClubForm(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./static/updUser.html")
	if err != nil {
		log.Fatal("allClubsHandle", err)
	}

	err = r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	id := r.FormValue("index")
	idInteger, err := strconv.Atoi(id)
	club := s.selectClub(idInteger)

	t.Execute(w, club)
}

func (s *server) allUserChangeHandle(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./static/updUsers.html")
	if err != nil {
		log.Fatal("allUsersHandle", err)
	}

	allUsers := s.selectUsers()
	errExecute := t.Execute(w, allUsers)
	// fmt.Println(allUsers[0].FullName)
	if errExecute != nil {
		log.Fatal("allUsersHandle2", err)
	}
}
func (s *server) allClubChangeHandle(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./static/updUsers.html")
	if err != nil {
		log.Fatal("allClubsHandle", err)
	}

	allClubs := s.selectClubs()
	errExecute := t.Execute(w, allClubs)
	// fmt.Println(allClubs[0].Name)
	if errExecute != nil {
		log.Fatal("allClubsHandle2", err)
	}
}
func (s *server) deleteUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	id := r.FormValue("id")
	idInt, err := strconv.Atoi(id)
	deleteUser(idInt, s)
	http.Redirect(w, r, "/index.html", http.StatusSeeOther)
}

func (s *server) deleteClub(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}
	index := r.FormValue("index")
	idInteger, err := strconv.Atoi(index)
	deleteClub(idInteger, s)
	http.Redirect(w, r, "/index.html", http.StatusSeeOther)
}

func (s *server) formHandle(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	Name := r.FormValue("name")
	Phone := r.FormValue("phone")
	Year_of_birth := r.FormValue("year_of_birth")
	userId := createUser(Name, Phone, Year_of_birth, s)

	person := ChildInfo{
		Child_id:      userId,
		Name:          Name,
		Phone:         Phone,
		Year_of_birth: Year_of_birth,
	}
	outputHTML(w, "./static/formComplete.html", person)
}
func (s *server) formClub(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	Club_name := r.FormValue("club_name")
	Price := r.FormValue("price")
	clubId := createClub(Club_name, Price, s)

	club := ClubInfo{
		Club_id:   clubId,
		Club_name: Club_name,
		Price:     Price,
	}

	outputHTML2(w, "./static/formComplete.html", club)
}
func createUser(name string, phone string, year_of_birth string, s *server) int {
	res, err := s.db.Exec("INSERT INTO child_users (Name, Phone, Year_of_birth) VALUES ($1, $2 , $3)", name, phone, year_of_birth)
	if err != nil {
		log.Fatal(err)
	}

	child_id, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}

	return int(child_id)
}

func updateUser(name string, phone string, year_of_birth string, id int, s *server) int {
	res, err := s.db.Exec("update child_users set Name=?, phone=? , year_of_birth=? where child_id=?", name, phone, year_of_birth, id)
	if err != nil {
		log.Fatal(err)
	}
	child_id, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	return int(child_id)
}

func deleteUser(id int, s *server) {
	_, err := s.db.Exec("delete from child_users where child_id=?", id)
	if err != nil {
		log.Fatal(err)
	}
}

func outputHTML(w http.ResponseWriter, filename string, person ChildInfo) {
	t, err := template.ParseFiles(filename)
	if err != nil {
		log.Fatal(err)
	}

	errExecute := t.Execute(w, person)

	if errExecute != nil {
		log.Fatal(errExecute)
	}
}
func createClub(club_name string, price string, s *server) int {
	res, err := s.db.Exec("INSERT INTO clubs (Club_ame, Price) VALUES ($1, $2 )", club_name, price)
	if err != nil {
		log.Fatal(err)
	}

	club_id, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}

	return int(club_id)
}

func updateClub(club_name string, price string, index int, s *server) int {
	res, err := s.db.Exec("update clubs set club_name=?, price=?  where club_id=?", club_name, price, index)
	if err != nil {
		log.Fatal(err)
	}
	club_id, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	return int(club_id)
}

func deleteClub(index int, s *server) {
	_, err := s.db.Exec("delete from clubs where club_id=?", index)
	if err != nil {
		log.Fatal(err)
	}
}

func outputHTML2(w http.ResponseWriter, filename string, club ClubInfo) {
	t, err := template.ParseFiles(filename)
	if err != nil {
		log.Fatal(err)
	}

	errExecute := t.Execute(w, club)

	if errExecute != nil {
		log.Fatal(errExecute)
	}
}

func main() {
	// Connecting database
	s := dbConnect()
	defer s.db.Close()
	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/", fileServer)
	http.HandleFunc("/form", s.formHandle)
	http.HandleFunc("/users", s.allUsersHandle)
	http.HandleFunc("/change", s.allUserChangeHandle)
	http.HandleFunc("/update", s.updateUserForm)
	http.HandleFunc("/delete", s.deleteUser)
	http.HandleFunc("/updateUserByID", s.updateUserByID)
	fmt.Println("Server running...")
	http.ListenAndServe(":8081", nil)

}
