package main

import (
	"log"
	"net/http"

	"github.com/AntonMezhanskiy/changelog/models"
	"github.com/julienschmidt/httprouter"
)

func main() {
	router := httprouter.New()

	router.ServeFiles("/static/*filepath", http.Dir("static/"))
	router.GET("/", models.Index)
	router.GET("/signin/", models.SignIn)
	router.POST("/signin/", models.SignInPost)
	router.GET("/entry/:entry/", models.Entry)

	log.Fatal(http.ListenAndServe(":8888", router))
}
