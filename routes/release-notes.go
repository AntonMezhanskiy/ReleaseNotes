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
)

func ReleaseNotes(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	t, err := template.ParseFiles("templates/release-notes.html", "templates/header.html", "templates/footer.html")
	checkErr(err)

	notes := []Note{}

	db, err := bolt.Open("my.db", 0600, nil)
	checkErr(err)
	defer db.Close()

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("notes"))
		if err != nil {
			return err
		}
		return nil
	})

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("notes"))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {

			buf := bytes.NewBuffer(v)
			dec := gob.NewDecoder(buf)

			var q Note
			err = dec.Decode(&q)
			if err != nil {
				log.Fatal("decode error 1:", err)
			}

			note := Note{q.ViewId, q.Title, q.PublicDate, q.Body, q.ReleaseNumber}
			parsed_time, _ := time.Parse(time_layout, note.PublicDate)
			note.PublicDate = parsed_time.Format(time_format)

			notes = append(notes, note)
		}

		return nil
	})

	data := struct {
		Title string
		Notes []Note
	}{
		Title: "Admin | Release Notes",
		Notes: notes,
	}

	t.ExecuteTemplate(w, "release-notes", data)
}
