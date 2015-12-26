package routes

import (
	"bytes"
	"encoding/gob"
	"html/template"
	"net/http"

	"github.com/boltdb/bolt"
	"github.com/julienschmidt/httprouter"
	"github.com/russross/blackfriday"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	htmlFuncMap := make(template.FuncMap)
	htmlFuncMap["markDown"] = MarkDown
	t, err := template.New("").Funcs(htmlFuncMap).ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")
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
		Title, Body string
	}{
		Title: "Admin | Release Notes",
		Body:  string(blackfriday.MarkdownBasic([]byte(Body))),
	}

	t.ExecuteTemplate(w, "index", data)
}
