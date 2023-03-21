package main

import (
	"database/sql"
	"net/http"
	"text/template"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

// Session
var (
	key   = []byte("APnyYC2IKcDXb3IS")
	store = sessions.NewCookieStore(key)
)

// Controller
func HomeHandler(response http.ResponseWriter, request *http.Request) {
	session, _ := store.Get(request, "authentication")

	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(response, request, "/login", http.StatusSeeOther)
	}

	tmpl := template.Must(template.ParseFiles("templates/home.html"))
	context := map[string]interface{}{
		"user": session.Values["user"],
	}

	tmpl.Execute(response, context)
}

type user struct {
	username string
	password string
}

func LoginHandler(response http.ResponseWriter, request *http.Request) {
	session, _ := store.Get(request, "authentication")

	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		if request.Method == http.MethodPost {
			db, dberr := sql.Open("mysql", "gurdayal-s:xmenguru77@(127.0.0.1:3306)/gowebdb")

			if dberr == nil {
				var form_user user
				form_user.username = request.FormValue("username")
				form_user.password = request.FormValue("password")

				var db_user user
				query := "SELECT username, password FROM users WHERE (username, password) = (?, ?)"

				if err := db.QueryRow(query, form_user.username, form_user.password).Scan(&db_user.username, &db_user.password); err != nil {
					http.Redirect(response, request, "/login", http.StatusSeeOther)
				}

				session.Values["authenticated"] = true
				session.Values["user"] = db_user.username
				session.Save(request, response)
				http.Redirect(response, request, "/", http.StatusSeeOther)
			}
		}
	} else {
		http.Redirect(response, request, "/", http.StatusSeeOther)
	}
	tmpl := template.Must(template.ParseFiles("templates/login.html"))
	tmpl.Execute(response, nil)
}

func SignupHandler(response http.ResponseWriter, request *http.Request) {

	if request.Method == http.MethodPost {
		db, dberr := sql.Open("mysql", "gurdayal-s:xmenguru77@(127.0.0.1:3306)/gowebdb")

		if dberr == nil {
			var form_user user

			form_user.username = request.FormValue("username")
			form_user.password = request.FormValue("password")
			created_at := time.Now()

			_, err := db.Exec(`INSERT INTO users (username, password, created_at) VALUES (?, ?, ?)`, form_user.username, form_user.password, created_at)

			if err == nil {
				http.Redirect(response, request, "/login", http.StatusSeeOther)
			}
		}
	}

	tmpl := template.Must(template.ParseFiles("templates/signup.html"))
	tmpl.Execute(response, nil)
}

func LogoutHandler(response http.ResponseWriter, request *http.Request) {
	session, _ := store.Get(request, "authentication")

	// Revoke users authentication
	session.Values["authenticated"] = false
	session.Values["user"] = "anonymous"
	session.Save(request, response)
	http.Redirect(response, request, "/login", http.StatusSeeOther)
}


func MarkAttendence(response http.ResponseWriter, request *http.Request) {
	session, _ := store.Get(request, "authentication")

	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(response, request, "/login", http.StatusSeeOther)
	} else {
		_, dberr := sql.Open("mysql", "gurdayal-s:xmenguru77@(127.0.0.1:3306)/gowebdb")
		
		if dberr == nil {
			username := session.Values["user"]
			// created_at := time.Now()
			month, year, _ := time.Now().Date()

			println(username)
			// println(created_at)
			println(month)
			println(year)

			// _, err := db.Exec(`INSERT INTO users (username, created_at, att_month, att_year) VALUES (?, ?, ?, ?)`, username, created_at, month, year)
			// if err == nil {
			// 	http.Redirect(response, request, "/login", http.StatusSeeOther)
			// }
		}

	}


}

// Url Routings
func Routes() {
	routes := mux.NewRouter()

	routes.HandleFunc("/", HomeHandler)
	routes.HandleFunc("/login", LoginHandler)
	routes.HandleFunc("/signup", SignupHandler)
	routes.HandleFunc("/logout", LogoutHandler)
	routes.HandleFunc("/mark-attendence", MarkAttendence)

	fs := http.FileServer(http.Dir("static/"))
	routes.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	http.ListenAndServe(":8000", routes)
}

// Entrypoint
func main() {
	println("Server started at :8000")
	Routes()
}
