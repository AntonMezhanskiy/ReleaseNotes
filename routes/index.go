package routes

import (
	"database/sql"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	_ "github.com/mattn/go-sqlite3"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	t, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")
	checkErr(err)

	db, err := sql.Open("sqlite3", "db/main.db")
	checkErr(err)

	// query
	rows, err := db.Query(`SELECT
		title,
		publicDate,
		viewid
		FROM posts WHERE visibility = '1' ORDER BY id DESC;`)
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

	posts := []Post{}

	for rows.Next() {
		post := Post{}
		err := rows.Scan(
			&post.Title,
			&post.PublicDate,
			&post.ViewId)
		checkErr(err)

		parsed_time, _ := time.Parse(time_layout, post.PublicDate)
		post.PublicDate = parsed_time.Format(time_format)

		// Replace all pairs.
		post.PublicDate = replacer.Replace(post.PublicDate)

		posts = append(posts, post)
	}
	rows.Close()
	db.Close()

	data := struct {
		Title string
		Posts []Post
	}{
		Title: "Title",
		Posts: posts,
	}

	t.ExecuteTemplate(w, "index", data)
}
