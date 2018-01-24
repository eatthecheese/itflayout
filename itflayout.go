package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	_ "github.com/go-sql-driver/mysql"
)

const NUMSCFIELDS = 10

type Sc struct {
	Scid          int
	Ip            string
	Location      string
	Version       string
	Nlc           int
	Scnumber      int
	Transportmode string
	Environment   string
	Priconc       string
	Secconc       string
	Devicesactive int
	//list      []MortyList
}

type ScList []Sc

type Example struct {
	Ip           string
	LocationName string
}

type ExampleList []Example

func handleIndex(w http.ResponseWriter, req *http.Request) {
	// below code is from http://www.alexedwards.net/blog/serving-static-sites-with-go
	lp := filepath.Join("static/templates", "layout.html")
	fp := filepath.Join("static/templates", filepath.Clean(req.URL.Path))

	t, _ := template.ParseFiles(lp, fp)
	t.ExecuteTemplate(w, "layout", nil)
}

func handleListOfAllSCs(w http.ResponseWriter, req *http.Request) {

	db, err := testlabConnectDb()
	if req.Method == "POST" {
		if err := req.ParseForm(); err != nil {
			log.Println(err)
		}
		//updateSc := Sc{}
		for key, values := range req.PostForm {
			fmt.Println(key, values)
		}
		//UpdateIntoScs(db, )
	}

	r, err := db.Query("SELECT * FROM list_of_scs m;")
	checkErr(err)
	defer r.Close()
	var listOfScs ScList = ScList{}
	listOfScs, _ = getListOfScs(r, listOfScs)

	t, _ := template.ParseFiles("static/templates/table.html")
	t.Execute(w, listOfScs)
}

// Connect to the testlab database
func testlabConnectDb() (*sql.DB, error) {
	userDb := "morty"
	pwDb := "True-cube1"
	connDb := "192.168.181.121:3306"
	schemaDb := "testlab"

	db, err := mySQLConnect(userDb, pwDb, connDb, schemaDb)
	checkErr(err)
	return db, err
}

// Cleanly connect to a mySQL database
func mySQLConnect(userDb string, pwDb string, connDb string, schemaDb string) (*sql.DB, error) {
	db, err := sql.Open("mysql", userDb+":"+pwDb+"@tcp"+"("+connDb+")/"+schemaDb)
	fmt.Println("Connecting to db...")
	checkErr(err)

	fmt.Println("Connected!")
	return db, err
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func insertIntoScs(db *sql.DB, Ip string, Location string, Version string, Nlc int, Scnumber int, Transportmode string, Environment string, Priconc string, Secconc string, Devicesactive int) {
	// prepare to insert some entries into list_of_scs
	stmt, err := db.Prepare("insert into list_of_scs (Ip, Location, Version, Nlc, Scnumber, Transportmode, Environment, Priconc, Secconc, Devicesactive) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	checkErr(err)
	result, err := stmt.Exec(Ip, Location, Version, Nlc, Scnumber, Transportmode, Environment, Priconc, Secconc, Devicesactive)
	checkErr(err)
	resRows, _ := result.RowsAffected()
	fmt.Println(resRows, "rows affected")
}

// Update records in the db - WIP, require data checking conditions
func UpdateIntoScs(db *sql.DB, Scid int, Ip string, Location string, Version string, Nlc int, Scnumber int, Transportmode string, Environment string, Priconc string, Secconc string, Devicesactive int) {
	// prepare to update some entries into list_of_scs
	stmt, err := db.Prepare("update list_of_scs set Ip=?, Location=?, Version=?, Nlc=?, Scnumber=?, Transportmode=?, Environment=?, Priconc=?, Secconc=?, Devicesactive=? where Scid=?")
	checkErr(err)
	result, err := stmt.Exec(Ip, Location, Version, Nlc, Scnumber, Transportmode, Environment, Priconc, Secconc, Devicesactive, Scid)
	checkErr(err)
	resRows, _ := result.RowsAffected()
	fmt.Println(resRows, "rows affected")
}

func getListOfScs(r *sql.Rows, listOfScs []Sc) (listOfScsOut []Sc, err error) {
	for i := 0; r.Next(); i++ {
		newSc := Sc{}
		err = r.Scan(&newSc.Scid, &newSc.Ip, &newSc.Location, &newSc.Version, &newSc.Nlc, &newSc.Scnumber, &newSc.Transportmode, &newSc.Environment, &newSc.Priconc, &newSc.Secconc, &newSc.Devicesactive)
		checkErr(err)
		listOfScs = append(listOfScs, newSc)
		//fmt.Println(listOfScs[i])
	}
	return listOfScs, err
}

func main() {
	// -------------------------------------
	// Handle HTML
	//---------------------------------------
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/itf/", handleListOfAllSCs)
	http.ListenAndServe(":8080", nil)
}
