package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Post struct {
	Title, PublicDate, ViewId, Body string
}

const (
	time_layout = "2006-01-02"
	time_format = "January 02, 2006"
)

func main() {
	router := httprouter.New()

	router.ServeFiles("/static/*filepath", http.Dir("static/"))
	router.GET("/", Index)
	router.GET("/signin/", SignIn)
	router.POST("/signin/", SignInPost)
	router.GET("/entry/:entry/", Entry)

	log.Fatal(http.ListenAndServe(":8888", router))
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func MarkDown(x string) interface{} {
	return template.HTML(x)
}
