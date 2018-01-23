package main

import (

	//"database/sql"

	"html/template"
	"net/http" //_ "github.com/go-sql-driver/mysql"
	"path/filepath"
)

type Sc struct {
	scid          int
	ip            string
	location      string
	version       string
	nlc           int
	scnumber      int
	transportmode string
	environment   string
	priconc       string
	secconc       string
	devicesactive int
	//list      []MortyList
}

type Example struct {
	Ip           string
	LocationName string
}

type ScList struct {
	Ip           []string
	LocationName []string
}

func handleIndex(w http.ResponseWriter, req *http.Request) {
	// below code is from http://www.alexedwards.net/blog/serving-static-sites-with-go
	lp := filepath.Join("static/templates", "layout.html")
	fp := filepath.Join("static/templates", filepath.Clean(req.URL.Path))

	t, _ := template.ParseFiles(lp, fp)
	t.ExecuteTemplate(w, "layout", nil)
}

func handleListOfAllSCs(w http.ResponseWriter, req *http.Request) {
	t, _ := template.ParseFiles("static/templates/table.html")
	p := ScList{Ip: []string{"192.168.181.14", "192.168.181.15"},
		LocationName: []string{"Campsie", "Domestic"}}
	t.Execute(w, p)
}

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/itf/", handleListOfAllSCs)
	http.ListenAndServe(":8080", nil)
}
