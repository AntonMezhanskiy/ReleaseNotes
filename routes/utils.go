package routes

import (
	"html/template"
)

type Post struct {
	Title, PublicDate, ViewId, Body string
}

const (
	time_layout = "2006-01-02"
	time_format = "January 02, 2006"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func MarkDown(x string) interface{} {
	return template.HTML(x)
}
