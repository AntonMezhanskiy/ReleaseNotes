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

func Edit(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	t, err := template.ParseFiles("templates/edit.html", "templates/header.html", "templates/footer.html")
	checkErr(err)

	db, err := bolt.Open("my.db", 0600, nil)
	checkErr(err)
	defer db.Close()

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("overview"))
		if err != nil {
			return err
		}
		return nil
	})

	var Body string

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("overview"))
		value := b.Get([]byte("overview"))

		buf := bytes.NewBuffer(value)
		dec := gob.NewDecoder(buf)

		err = dec.Decode(&Body)
		if err != nil {
			// log.Fatal("decode error 1:", err)
		}

		return nil
	})

	data := struct {
		Title string
		Body  string
	}{
		Title: "Admin | Release Notes",
		Body:  Body,
	}

	t.ExecuteTemplate(w, "edit", data)
}

func EditPost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	body := r.FormValue("body")

	var bb bytes.Buffer
	enc := gob.NewEncoder(&bb)
	err := enc.Encode(body)
	if err != nil {
		log.Fatal("encode:", err)
	}

	db, err := bolt.Open("my.db", 0600, nil)
	checkErr(err)
	defer db.Close()

	db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("overview"))
		if err != nil {
			return err
		}
		return b.Put([]byte("overview"), bb.Bytes())
	})

	http.Redirect(w, r, "/", 301)
}
