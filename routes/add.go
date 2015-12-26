package routes

import (
	"bytes"
	"encoding/gob"
	"html/template"
	"log"
	"net/http"

	"github.com/boltdb/bolt"
	"github.com/julienschmidt/httprouter"
)

func Add(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	t, err := template.ParseFiles("templates/add.html", "templates/header.html", "templates/footer.html")
	checkErr(err)

	data := struct {
		Title string
	}{
		Title: "Admin | Release Notes",
	}

	t.ExecuteTemplate(w, "add", data)
}

func AddPost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	title := r.FormValue("title")
	publicDate := r.FormValue("publicDate")
	body := r.FormValue("body")
	releaseNumber := r.FormValue("releaseNumber")

	note := Note{title, title, publicDate, body, releaseNumber}
	var bb bytes.Buffer
	enc := gob.NewEncoder(&bb)
	err := enc.Encode(note)
	if err != nil {
		log.Fatal("encode:", err)
	}

	db, err := bolt.Open("my.db", 0600, nil)
	checkErr(err)
	defer db.Close()

	db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("notes"))
		if err != nil {
			return err
		}
		return b.Put([]byte(title), bb.Bytes())
	})

	http.Redirect(w, r, "/release-notes/", 301)
}
