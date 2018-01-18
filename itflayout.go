package main

import (
	//"database/sql"

	"html/template"
	"net/http" //_ "github.com/go-sql-driver/mysql"
	"path/filepath"
)

func handleIndex(w http.ResponseWriter, req *http.Request) {
	// below code is from http://www.alexedwards.net/blog/serving-static-sites-with-go
	lp := filepath.Join("static/templates", "layout.html")
	fp := filepath.Join("static/templates", filepath.Clean(req.URL.Path))

	t, _ := template.ParseFiles(lp, fp)
	t.ExecuteTemplate(w, "layout", nil)
}

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", handleIndex)
	http.ListenAndServe(":8080", nil)
}
