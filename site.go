package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
)

var templates = template.Must(template.ParseFiles("base.html", "home.html"))

func main() {
	handle := func(err error) {
		if err != nil {
			panic(err)
		}
	}
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := struct{}{}
		err := templates.Execute(w, data)
		handle(err)
	})
	port := os.Getenv("PORT")
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	handle(err)
}
