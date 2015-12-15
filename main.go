package main

import (
	"log"
	"net/http"

	"github.com/AntonMezhanskiy/changelog/routes"
	"github.com/julienschmidt/httprouter"
)

func main() {
	router := httprouter.New()

	router.ServeFiles("/static/*filepath", http.Dir("static/"))

	router.GET("/", routes.Index)
	router.GET("/signin/", routes.SignIn)
	router.GET("/entry/:entry/", routes.Entry)

	router.POST("/signin/", routes.SignInPost)

	log.Fatal(http.ListenAndServe(":8888", router))
}
