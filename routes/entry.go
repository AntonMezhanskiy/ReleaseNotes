package routes

import (
	"bytes"
	"encoding/gob"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/boltdb/bolt"
	"github.com/julienschmidt/httprouter"
	"github.com/russross/blackfriday"
)

func Entry(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	entry := ps.ByName("entry")

	htmlFuncMap := make(template.FuncMap)
	htmlFuncMap["markDown"] = MarkDown
	t, err := template.New("").Funcs(htmlFuncMap).ParseFiles("templates/entry.html", "templates/header.html", "templates/footer.html")
	checkErr(err)

	var Title, PublicDate, Body, ReleaseNumber string

	db, err := bolt.Open("my.db", 0600, nil)
	checkErr(err)
	defer db.Close()

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("notes"))
		c := b.Cursor()
		min := []byte(entry)

		for k, v := c.Seek(min); k != nil && bytes.Compare(k, min) <= 0; k, v = c.Next() {
			buf := bytes.NewBuffer(v)
			dec := gob.NewDecoder(buf)

			var q Note
			err = dec.Decode(&q)
			if err != nil {
				log.Fatal("decode error 1:", err)
			}
			Title = q.Title
			PublicDate = q.PublicDate
			Body = q.Body
			ReleaseNumber = q.ReleaseNumber

			parsed_time, _ := time.Parse(time_layout, PublicDate)
			PublicDate = parsed_time.Format(time_format)

			Body = string(blackfriday.MarkdownBasic([]byte(Body)))
		}

		return nil
	})

	if Title == "" {
		http.Redirect(w, r, "/", 301)
	}

	data := struct {
		Title, PublicDate, Body, ReleaseNumber string
	}{
		Title:         Title,
		PublicDate:    PublicDate,
		Body:          Body,
		ReleaseNumber: ReleaseNumber,
	}

	t.ExecuteTemplate(w, "entry", data)
}
