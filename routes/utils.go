package routes

import (
	"html/template"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func MarkDown(x string) interface{} {
	return template.HTML(x)
}
