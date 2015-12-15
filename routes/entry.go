package routes

import (
	"database/sql"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	_ "github.com/mattn/go-sqlite3"
	"github.com/russross/blackfriday"
)

func Entry(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	entry := ps.ByName("entry")

	htmlFuncMap := make(template.FuncMap)
	htmlFuncMap["markDown"] = MarkDown
	t, err := template.New("").Funcs(htmlFuncMap).ParseFiles("templates/entry.html", "templates/header.html", "templates/footer.html")
	checkErr(err)

	db, err := sql.Open("sqlite3", "db/main.db")
	checkErr(err)

	rows, err := db.Query(`SELECT
		title, 
		publicDate,
		description
		FROM posts WHERE viewid=?;`, entry)
	checkErr(err)

	// Create replacer with pairs as arguments.
	replacer := strings.NewReplacer(
		"January", "Январь",
		"February", "Февраль",
		"March", "Март",
		"April", "Апрель",
		"May", "Май",
		"June", "Июнь",
		"July", "Июль",
		"August", "Август",
		"September", "Сентябрь",
		"October", "Октябрь",
		"November", "Ноябрь",
		"December", "Декабрь")

	var Title, PublicDate, Body string

	for rows.Next() {
		err := rows.Scan(
			&Title,
			&PublicDate,
			&Body)
		checkErr(err)

		parsed_time, _ := time.Parse(time_layout, PublicDate)
		PublicDate = parsed_time.Format(time_format)

		// Replace all pairs.
		PublicDate = replacer.Replace(PublicDate)

		Body = string(blackfriday.MarkdownBasic([]byte(Body)))
	}

	rows.Close()
	db.Close()

	if Title == "" {
		http.Redirect(w, r, "/", 301)
	}

	data := struct {
		Title, PublicDate, Body string
	}{
		Title:      Title,
		PublicDate: PublicDate,
		Body:       Body,
	}

	t.ExecuteTemplate(w, "entry", data)
}
