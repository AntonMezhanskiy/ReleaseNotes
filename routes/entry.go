package routes

import (
	"bytes"
	"encoding/gob"
	// "fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
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

	db, err := bolt.Open("my.db", 0600, nil)
	checkErr(err)
	defer db.Close()

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("posts"))
		c := b.Cursor()
		min := []byte(entry)

		for k, v := c.Seek(min); k != nil && bytes.Compare(k, min) <= 0; k, v = c.Next() {
			buf := bytes.NewBuffer(v)
			dec := gob.NewDecoder(buf)

			var q Post
			err = dec.Decode(&q)
			if err != nil {
				log.Fatal("decode error 1:", err)
			}
			Title = q.Title
			PublicDate = q.PublicDate
			Body = q.Body

			parsed_time, _ := time.Parse(time_layout, PublicDate)
			PublicDate = parsed_time.Format(time_format)

			// Replace all pairs.
			PublicDate = replacer.Replace(PublicDate)

			Body = string(blackfriday.MarkdownBasic([]byte(Body)))
			// fmt.Printf("%s: %s\n", new(big.Int).SetBytes(k), q.Title)
		}

		return nil
	})

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
