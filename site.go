package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

type worksheet struct {
	WorksheetName, Problems, Solutions string
}

type unit struct {
	UnitName   string
	Worksheets []worksheet
}

const (
	BaseTemplatePath string = "templates/base.html"
)

var (
	homeTemplate      = getTemplate("home.html")
	unitTemplate      = getTemplate("unit.html")
	worksheetTemplate = getTemplate("worksheet.html")
)

func getTemplate(path string) *template.Template {
	return template.Must(template.ParseFiles(BaseTemplatePath, fmt.Sprintf("templates/%s", path)))
}

func main() {
	handle := func(err error) {
		if err != nil {
			panic(err)
		}
	}
	var units []unit
	bytes, err := ioutil.ReadFile("homework_data.json")
	handle(err)
	err = json.Unmarshal(bytes, &units)
	handle(err)
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := struct {
			Units []unit
		}{
			units,
		}
		err := homeTemplate.Execute(w, data)
		handle(err)
	})
	r.HandleFunc("/units/{unit}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		index, err := strconv.Atoi(vars["unit"])
		handle(err)
		data := struct {
			Unit      unit
			UnitIndex int
			Units     []unit
		}{
			units[index],
			index,
			units,
		}
		err = unitTemplate.Execute(w, data)
		handle(err)
	})
	r.HandleFunc("/units/{unit}/worksheets/{worksheet}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		unitIdx, err := strconv.Atoi(vars["unit"])
		handle(err)
		worksheetIdx, err := strconv.Atoi(vars["worksheet"])
		handle(err)
		data := struct {
			Unit           unit
			UnitIndex      int
			Worksheet      worksheet
			WorksheetIndex int
			Units          []unit
		}{
			units[unitIdx],
			unitIdx,
			units[unitIdx].Worksheets[worksheetIdx],
			worksheetIdx,
			units,
		}
		err = worksheetTemplate.Execute(w, data)
		handle(err)
	})
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.Handle("/", r)
	port := os.Getenv("PORT")
	err = http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	handle(err)
}
