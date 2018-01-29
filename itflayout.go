package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

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
type Device struct {
	ParentSc      *Sc
	Deviceid      int
	Ip            string
	VersionDevice string
	VersionRtd    string
	DeviceType    string
	DoppIp        string
	DoppPort      string
	Plinth        int
	VersionEprom  string
	IsFLR         bool
	IsE2Gate      bool
	IsRLG         bool
}

type ScList []Sc
type DeviceList []Device

func handleIndex(w http.ResponseWriter, req *http.Request) {
	// below code is from http://www.alexedwards.net/blog/serving-static-sites-with-go
	lp := filepath.Join("static/templates", "layout.html")
	fp := filepath.Join("static/templates", filepath.Clean(req.URL.Path))

	t, _ := template.ParseFiles(lp, fp)
	t.ExecuteTemplate(w, "layout", nil)
}

func handleListOfAllDevices(w http.ResponseWriter, req *http.Request) {
	db, err := testlabConnectDb()
	r, err := db.Query(`SELECT d.deviceid, d.ip, d.version, d.version_rtd, d.devicetype, d.doppip, d.doppport, d.plinth, d.scip, d.version_eprom, 
						s.location, s.nlc, s.environment, s.transportmode from list_of_devices d, list_of_scs s 
						where d.scip = s.ip;`)
	checkErr(err)
	defer r.Close()

	var listOfDevices DeviceList = DeviceList{}
	listOfDevices, _ = getListOfDevices(r, listOfDevices)

	t, _ := template.ParseFiles("static/templates/listofdevices.html")
	t.Execute(w, listOfDevices)
}

func handleListOfAllSCs(w http.ResponseWriter, req *http.Request) {
	db, err := testlabConnectDb()
	if req.Method == "POST" {
		if err := req.ParseForm(); err != nil {
			log.Println(err)
		}
		updateSc := Sc{}
		newInputs := req.PostForm
		updateScid := newInputs.Get("id") // Get the SC ID to be updated
		updateScidInt, _ := strconv.Atoi(updateScid)
		// Push the updated SC values
		updateSc.Ip = newInputs.Get("new_ip" + updateScid)
		updateSc.Location = newInputs.Get("new_location" + updateScid)
		updateSc.Version = newInputs.Get("new_version" + updateScid)
		updateSc.Nlc, _ = strconv.Atoi(newInputs.Get("new_nlc" + updateScid))
		updateSc.Scnumber, _ = strconv.Atoi(newInputs.Get("new_scnumber" + updateScid))
		updateSc.Transportmode = newInputs.Get("new_transportmode" + updateScid)
		updateSc.Environment = newInputs.Get("new_environment" + updateScid)
		updateSc.Priconc = newInputs.Get("new_priconc" + updateScid)
		updateSc.Secconc = newInputs.Get("new_secconc" + updateScid)
		updateSc.Devicesactive, _ = strconv.Atoi(newInputs.Get("new_devicesactive" + updateScid))

		if updateScidInt != 0 {
			updateIntoScs(db, updateScidInt, updateSc.Ip, updateSc.Location, updateSc.Version, updateSc.Nlc, updateSc.Scnumber, updateSc.Transportmode, updateSc.Environment, updateSc.Priconc, updateSc.Secconc, updateSc.Devicesactive)
		} else {
			insertIntoScs(db, updateSc.Ip, updateSc.Location, updateSc.Version, updateSc.Nlc, updateSc.Scnumber, updateSc.Transportmode, updateSc.Environment, updateSc.Priconc, updateSc.Secconc, updateSc.Devicesactive)
		}
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
func updateIntoScs(db *sql.DB, Scid int, Ip string, Location string, Version string, Nlc int, Scnumber int, Transportmode string, Environment string, Priconc string, Secconc string, Devicesactive int) {
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

func getListOfDevices(r *sql.Rows, listOfDevices []Device) (listOfDevicesOut []Device, err error) {
	for i := 0; r.Next(); i++ {
		newDevice := Device{}
		newSc := Sc{}
		newDevice.ParentSc = &newSc
		err = r.Scan(&newDevice.Deviceid, &newDevice.Ip, &newDevice.VersionDevice, &newDevice.VersionRtd, &newDevice.DeviceType,
			&newDevice.DoppIp, &newDevice.DoppPort, &newDevice.Plinth, &newDevice.ParentSc.Ip, &newDevice.VersionEprom,
			&newDevice.ParentSc.Location, &newDevice.ParentSc.Nlc, &newDevice.ParentSc.Environment, &newDevice.ParentSc.Transportmode)
		checkErr(err)
		if newDevice.DeviceType == "RLG MSTRP2" || newDevice.DeviceType == "RLG CTP" {
			newDevice.IsRLG = true
		} else if newDevice.DeviceType == "FLR" || newDevice.DeviceType == "FLR CTP" {
			newDevice.IsFLR = true
		} else if newDevice.DeviceType == "E2 Gate" {
			newDevice.IsE2Gate = true
		}

		listOfDevices = append(listOfDevices, newDevice)
	}
	return listOfDevices, err
}

func main() {
	// -------------------------------------
	// Handle HTML
	//---------------------------------------
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/itf/", handleListOfAllSCs)
	http.HandleFunc("/itf/devices", handleListOfAllDevices)
	http.ListenAndServe(":8080", nil)
}
