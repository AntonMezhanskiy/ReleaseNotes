package models

import (
	"html/template"
	"net/http"

	"github.com/julienschmidt/httprouter"
	_ "github.com/mattn/go-sqlite3"
)

func SignIn(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	t, err := template.ParseFiles("templates/signin.html", "templates/header.html", "templates/footer.html")
	checkErr(err)

	data := struct {
		Title string
	}{
		Title: "Авторизация",
	}

	t.ExecuteTemplate(w, "signin", data)
}

func SignInPost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// pass := r.FormValue("pass")
	// email := r.FormValue("email")

	cookie := http.Cookie{Name: "user_session", Value: "test", Path: "/", HttpOnly: true}
	http.SetCookie(w, &cookie)

	http.Redirect(w, r, "/", 301)
}
