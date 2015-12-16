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
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	t, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")
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

	db, err := bolt.Open("my.db", 0600, nil)
	checkErr(err)
	defer db.Close()

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("posts"))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {

			buf := bytes.NewBuffer(v)
			dec := gob.NewDecoder(buf)

			var q Post
			err = dec.Decode(&q)
			if err != nil {
				log.Fatal("decode error 1:", err)
			}

			// fmt.Println(new(big.Int).SetBytes(k), q.Title, q.Content, q.Created)
			post := Post{q.Title, q.PublicDate, q.ViewId, q.Body}
			parsed_time, _ := time.Parse(time_layout, post.PublicDate)
			post.PublicDate = parsed_time.Format(time_format)

			// Replace all pairs.
			post.PublicDate = replacer.Replace(post.PublicDate)

			posts = append(posts, post)
		}

		return nil
	})

	data := struct {
		Title string
		Posts []Post
	}{
		Title: "Title",
		Posts: posts,
	}

	t.ExecuteTemplate(w, "index", data)
}
