package main

import (
	"html/template"
	"net/http"
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
	err := http.ListenAndServe("localhost:8080", nil)
	handle(err)
}
