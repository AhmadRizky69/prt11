package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kataras/go-sessions"
)

var err error
var db *sql.DB

type MyMux struct{}

type Employee struct {
	Id       int
	Username string
	Password string
	Age      int
	Email    string
}

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := ""
	dbName := "company"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		log.Fatal(err)
	}
	// if err != nil {
	// 	panic(err.Error())
	// }
	return db
}

var tmpl = template.Must(template.ParseGlob("form/*"))

func checkErr(w http.ResponseWriter, r *http.Request, err error) bool {
	if err != nil {

		fmt.Println(r.Host + r.URL.Path)

		http.Redirect(w, r, r.Host+r.URL.Path, 301)
		return false
	}

	return true
}

func login(w http.ResponseWriter, r *http.Request) {
	session := sessions.Start(w, r)
	if len(session.GetString("username")) != 0 && checkErr(w, r, err) {
		http.Redirect(w, r, "/", 302)
	}
	if r.Method != "POST" {
		http.ServeFile(w, r, "form/login.html")
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")

	users := QueryUser(username)

	//deskripsi dan compare password
	// var password_tes = bcrypt.CompareHashAndPassword([]byte(users.Password), []byte(password))
	// fmt.Println(password_tes)
	if password == users.Password {
		// login success
		//  session := sessions.Start(w, r)

		//  session.Set("username", users.Username)
		//  session.Set("age", users.Age)
		http.Redirect(w, r, "/index", 302)
	} else {
		//login failed
		http.Redirect(w, r, "/login", 302)
	}

}
func QueryUser(username string) Employee {
	db := dbConn()
	var users = Employee{}
	err = db.QueryRow(`
		SELECT id, 
		username, 
		password, 
		age
		FROM Employee WHERE email=?
		`, username).
		Scan(
			&users.Id,
			&users.Username,
			&users.Password,
			&users.Age,
		)
	return users
}

func logout(w http.ResponseWriter, r *http.Request) {
	session := sessions.Start(w, r)
	session.Clear()
	sessions.Destroy(w, r)
	http.Redirect(w, r, "/", 302)
}

func Index(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	selDB, err := db.Query("SELECT * FROM Employee ORDER BY id DESC")
	if err != nil {
		log.Fatal(err)
	}
	// if err != nil {
	// 	panic(err.Error())
	// }
	emp := Employee{}
	res := []Employee{}
	for selDB.Next() {
		var id, age int
		var username, password, email string
		err = selDB.Scan(&id, &username, &password, &age, &email)
		if err != nil {
			log.Fatal(err)
		}
		// if err != nil {
		// 	panic(err.Error())
		// }
		emp.Id = id
		emp.Username = username
		emp.Password = password
		emp.Age = age
		emp.Email = email
		res = append(res, emp)
	}
	tmpl.ExecuteTemplate(w, "Index", res)
	defer db.Close()
}

func Show(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM Employee WHERE id=?", nId)
	if err != nil {
		log.Fatal(err)
	}
	// if err != nil {
	// 	panic(err.Error())
	// }
	emp := Employee{}
	for selDB.Next() {
		var id, age int
		var username, password, email string
		err = selDB.Scan(&id, &username, &password, &age, &email)
		if err != nil {
			log.Fatal(err)
		}
		// if err != nil {
		// 	panic(err.Error())
		// }
		emp.Id = id
		emp.Username = username
		emp.Password = password
		emp.Age = age
		emp.Email = email
	}
	tmpl.ExecuteTemplate(w, "Show", emp)
	defer db.Close()
}

func New(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "New", nil)
}

func Edit(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM Employee WHERE id=?", nId)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	if err != nil {
		panic(err.Error())
	}
	emp := Employee{}
	for selDB.Next() {
		var id, age int
		var username, password, email string
		err = selDB.Scan(&id, &username, &password, &age, &email)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		if err != nil {
			panic(err.Error())
		}
		emp.Id = id
		emp.Username = username
		emp.Password = password
		emp.Age = age
		emp.Email = email
	}
	tmpl.ExecuteTemplate(w, "Edit", emp)
	defer db.Close()
}

func Insert(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")
		age := r.FormValue("age")
		email := r.FormValue("email")
		insForm, err := db.Prepare("INSERT INTO Employee(username,password,age, email) VALUES(?,?,?,?)")
		if err != nil {
			log.Fatal(err)
		}
		// if err != nil {
		// 	panic(err.Error())
		// }
		insForm.Exec(username, password, age, email)
		log.Println("Name: " + username + "Password :" + password + "age :" + age + "email" + email)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func register(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.ServeFile(w, r, "form/register.html")
		username := r.FormValue("username")
		password := r.FormValue("password")
		age := r.FormValue("age")
		email := r.FormValue("email")
		insForm, err := db.Prepare("INSERT INTO Employee(username,password,age, email) VALUES(?,?,?,?)")
		if err != nil {
			log.Fatal(err)
		}
		insForm.Exec(username, password, age, email)
		log.Println("Name: " + username + "Password :" + password + "age :" + age + "email" + email)
		//return
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)

	// username := r.FormValue("username")
	// password := r.FormValue("password")
	// age := r.FormValue("age")
	// email := r.FormValue("email")

	// users := QueryUser(email)

	// if (Employee{}) == users {
	// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	// 	if len(hashedPassword) != 0 && checkErr(w, r, err) {
	// 		stmt, err := db.Prepare("INSERT INTO employee SET username=?, password=?, age=?, email=?")
	// 		if err == nil {
	// 			_, err := stmt.Exec(&username, &hashedPassword, &age, &email)
	// 			if err != nil {
	// 				http.Error(w, err.Error(), http.StatusInternalServerError)
	// 				return
	// 			}

	// 			http.Redirect(w, r, "/login", http.StatusSeeOther)
	// 			return
	// 		}
	// 	}
	// } else {
	// 	http.Redirect(w, r, "/register", 302)
	// }
}

func Update(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")
		age := r.FormValue("age")
		email := r.FormValue("email")
		id := r.FormValue("uid")
		insForm, err := db.Prepare("UPDATE Employee SET username=?, password=?, age=?, email=? WHERE id=?")
		if err != nil {
			log.Fatal(err)
		}
		// if err != nil {
		// 	panic(err.Error())
		// }
		insForm.Exec(username, password, age, email, id)
		log.Println("Update Name: " + username + "Password :" + password + "age :" + age + "email" + email)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	emp := r.URL.Query().Get("id")
	delForm, err := db.Prepare("DELETE FROM Employee WHERE id=?")
	if err != nil {
		log.Fatal(err)
	}
	// if err != nil {
	// 	panic(err.Error())
	// }
	delForm.Exec(emp)
	log.Println("DELETE")
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}
func (p *MyMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		login(w, r)
	case "/show":
		Show(w, r)
	case "/new":
		New(w, r)
	case "/edit":
		Edit(w, r)
	case "/index":
		Index(w, r)
	case "/register":
		register(w, r)
	case "/insert":
		Insert(w, r)
	case "/update":
		Update(w, r)
	case "/delete":
		Delete(w, r)
	// case "/login":
	// 	login(w, r)
	default:
		http.NotFound(w, r)
	}
	return
}
func validateForm(val url.Values) bool {
	username := val.Get("username")
	password := val.Get("password")
	age := val.Get("age")
	email := val.Get("email")

	var result bool = true

	if len(username) < 5 {
		fmt.Println("Minimal karakter username adalah 5: ", username)
		result = false
	}
	if len(password) < 8 {
		fmt.Println("Minimal karakter password adalah 8: ", password)
		result = false
	}
	m, _ := regexp.MatchString(`/^(([^<>()\[\]\.,;:\s@\"]+(\.[^<>()\[\]\.,;:\s@\"]+)*)|(\".+\"))@(([^<>()[\]\.,;:\s@\"]+\.)+[^<>()[\]\.,;:\s@\"]{2,})$/i`, email)

	if !m {
		fmt.Println("Format email salah: ", email)
		result = false
	}

	valInt, err := strconv.Atoi(age)
	if err != nil || valInt < 0 {
		fmt.Println("Age bukan bilangan bulat positif: ", age)
		result = false
	}

	return result
}

func main() {
	mux := &MyMux{}

	err := http.ListenAndServe(":9090", mux)

	if err != nil {
		log.Fatal("Something error")
	}

}
