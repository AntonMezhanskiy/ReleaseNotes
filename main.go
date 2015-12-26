package main

import (
	"log"
	"net/http"

	"github.com/AntonMezhanskiy/ReleaseNotes/routes"
	"github.com/julienschmidt/httprouter"
)

func main() {
	router := httprouter.New()

	// static files
	router.ServeFiles("/static/*filepath", http.Dir("static/"))

	// GET
	router.GET("/", routes.Index)
	router.GET("/release-notes/", routes.ReleaseNotes)
	router.GET("/release-notes/:entry/", routes.Entry)
	router.GET("/add/", routes.Add)
	router.GET("/edit/", routes.Edit)

	// POST
	router.POST("/add/", routes.AddPost)
	router.POST("/edit/", routes.EditPost)

	log.Fatal(http.ListenAndServe(":8888", router))
}
