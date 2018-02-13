package main

/*
Todo ----------------------------
List of bus rigs page
List of DOPPs page, should be similar to List Of SCs page
Clean up unused .html and related go code
Show MM6 BVs (Newcastle)
Show CLD
Clean up top menu buttons into a ribbon of options
**Fix bug for when Device has unknown SC, disappears from List of Devices view**
Add warnings for if invalid details are entered in List of SCs page
*Add filters for Cubic Only, ETS Only, All Devices, etc. for List of Devices/List of SCs page*
Aesthetic improvements to the table views
Fix clipping issue on visual layout page
Aesthetic improvements to the visual layout page
Add clickable options for the visual layout page
Add Collaborative Test Team-facing view for List of devices
Separate javascript from html files
Add integration with change requests
Add integration with Device Versions in the ITF
**Add Cancel button for when Edit button is pressed accidentally in List of SCs/Devices pages**

*/

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
const NewRLG = 100
const NewE2G = 101
const NewFLR = 102

type Sc struct {
	Scname        string
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
	Changeby      string
	Changetime    string
	Comment       string
	//list      []MortyList
}
type Device struct {
	ParentSc      *Sc
	Devicename    string
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
	Changeby      string
	Changetime    string
	Comment       string
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

	// Work in Progress below
	if req.Method == "POST" {
		if err := req.ParseForm(); err != nil {
			log.Println(err)
		}
		updateDevice := Device{}
		updateDeviceSc := Sc{}
		updateDevice.ParentSc = &updateDeviceSc
		newInputs := req.PostForm
		updateDevice.Devicename = newInputs.Get("id") // Get the Device ID to be updated
		//updateDeviceidInt, _ := strconv.Atoi(updateDeviceid)
		// Push the updated Device values
		updateDevice.Ip = newInputs.Get("new_ip")
		updateDevice.ParentSc.Ip = newInputs.Get("new_scip")
		updateDevice.VersionDevice = newInputs.Get("new_versiondevice")
		updateDevice.VersionRtd = newInputs.Get("new_versionrtd")
		updateDevice.Plinth, _ = strconv.Atoi(newInputs.Get("new_plinth"))
		updateDevice.Changeby = newInputs.Get("new_changeby")
		updateDevice.Comment = newInputs.Get("new_comment")
		updateDevice.VersionEprom = newInputs.Get("new_versioneprom")
		updateDevice.DoppPort = newInputs.Get("new_doppport")
		updateDevice.DoppIp = newInputs.Get("new_doppip")

		if newInputs.Get("devicetype") == "rlg" {
			updateDevice.IsRLG = true
			updateDevice.DeviceType = "RLG MSTRP2"
			/*updateDevice.VersionEprom = newInputs.Get("new_versioneprom")
			updateDevice.DoppPort = newInputs.Get("new_doppport")
			updateDevice.DoppIp = newInputs.Get("new_doppip")*/
		} else if newInputs.Get("devicetype") == "flr" {
			updateDevice.IsFLR = true
			updateDevice.DeviceType = "FLR"
			/*updateDevice.DoppPort = newInputs.Get("new_doppport")
			updateDevice.DoppIp = newInputs.Get("new_doppip")*/
		} else if newInputs.Get("devicetype") == "e2gate" { // is E2 gate
			updateDevice.IsE2Gate = true
			updateDevice.DeviceType = "E2 Gate"
		}

		updateIntoDevices(db, updateDevice.Devicename, updateDevice.Ip, updateDevice.VersionDevice, updateDevice.VersionRtd, updateDevice.DeviceType, updateDevice.DoppIp,
			updateDevice.DoppPort, updateDevice.Plinth, updateDevice.ParentSc.Ip, updateDevice.VersionEprom, updateDevice.Changeby)
		updateIntoDevicesComments(db, updateDevice.Comment)
	}

	r, err := db.Query(`SELECT d.devicename, d.deviceid, d.ip, d.version, d.version_rtd, d.devicetype, d.doppip, d.doppport, d.plinth, d.scip, d.version_eprom,
		s.location, s.nlc, s.environment, s.transportmode, d.changeby, d.changetime, ifnull(c.comment, "no comment") as comment
		from
		(select *
		from list_of_devices i
		where changetime=(select max(j.changetime) from list_of_devices j where i.devicename=j.devicename)
		order by ip asc) d
		left join list_of_devices_comments c on c.id = d.deviceid,
		(select ip, location, nlc, environment, transportmode
		from list_of_scs a
		where changetime=(select max(b.changetime) from list_of_scs b where a.scname=b.scname)
		order by ip asc) s
		where d.scip = s.ip
		order by d.ip asc;`)
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
		updateSc.Scname = newInputs.Get("id") // Get the SC Name to be updated
		//updateScidInt, _ := strconv.Atoi(updateScid)
		// Push the updated (and old) SC values
		updateSc.Ip = newInputs.Get("new_ip")
		updateSc.Location = newInputs.Get("new_location")
		updateSc.Version = newInputs.Get("new_version")
		updateSc.Nlc, _ = strconv.Atoi(newInputs.Get("new_nlc"))
		updateSc.Scnumber, _ = strconv.Atoi(newInputs.Get("new_scnumber"))
		updateSc.Transportmode = newInputs.Get("new_transportmode")
		updateSc.Environment = newInputs.Get("new_environment")
		updateSc.Priconc = newInputs.Get("new_priconc")
		updateSc.Secconc = newInputs.Get("new_secconc")
		updateSc.Changeby = newInputs.Get("new_changeby")
		updateSc.Comment = newInputs.Get("new_comment")
		//updateSc.Devicesactive, _ = strconv.Atoi(newInputs.Get("new_devicesactive" + updateScid))
		updateIntoScs(db, updateSc.Scname, updateSc.Ip, updateSc.Location, updateSc.Version, updateSc.Nlc, updateSc.Scnumber, updateSc.Transportmode, updateSc.Environment, updateSc.Priconc, updateSc.Secconc, updateSc.Changeby)
		updateIntoScsComments(db, updateSc.Comment)
		/*
			if updateScidInt != 0 {
				updateIntoScs(db, updateScidInt, updateSc.Ip, updateSc.Location, updateSc.Version, updateSc.Nlc, updateSc.Scnumber, updateSc.Transportmode, updateSc.Environment, updateSc.Priconc, updateSc.Secconc)
			} else {
				insertIntoScs(db, updateSc.Ip, updateSc.Location, updateSc.Version, updateSc.Nlc, updateSc.Scnumber, updateSc.Transportmode, updateSc.Environment, updateSc.Priconc, updateSc.Secconc)
			}
		*/
	}
	// FIX THIS QUERY TO ACCOUNT FOR NEW COLUMN 'SCNAME'
	r, err := db.Query(`select s.scname, s.scid, s.ip, s.location, s.version, s.nlc, s.scnumber, s.transportmode, s.environment, s.priconc, s.secconc, s.changeby, s.changetime, d.devices, ifnull(c.comment, "no comment") as comment
		from (select *
		from list_of_scs a
		where changetime=(select max(b.changetime) from list_of_scs b where a.scname=b.scname)
		order by ip asc) s
		left join
		(select list_of_scs.ip, count(distinct list_of_devices.scip) as devices
		from list_of_scs
		left join list_of_devices on list_of_scs.ip = list_of_devices.scip
		group by list_of_scs.ip) d on s.ip = d.ip
		left join list_of_scs_comments c on c.id = s.scid
		order by s.ip asc;`)
	checkErr(err)
	defer r.Close()
	var listOfScs ScList = ScList{}
	listOfScs, _ = getListOfScs(r, listOfScs)

	t, _ := template.ParseFiles("static/templates/table.html")
	t.Execute(w, listOfScs)
}

func handleListOfAllBusRigs(w http.ResponseWriter, req *http.Request) {
	t, _ := template.ParseFiles("static/templates/listofbusrigs.html")
	t.Execute(w, t)

}

func handleLayout(w http.ResponseWriter, req *http.Request) {
	db, err := testlabConnectDb()

	r, err := db.Query(`SELECT d.devicename, d.deviceid, d.ip, d.version, d.version_rtd, d.devicetype, d.doppip, d.doppport, d.plinth, d.scip, d.version_eprom,
		s.location, s.nlc, s.environment, s.transportmode, d.changeby, d.changetime, ifnull(c.comment, "no comment") as comment
		from
		(select *
		from list_of_devices i
		where changetime=(select max(j.changetime) from list_of_devices j where i.devicename=j.devicename)
		order by ip asc) d
		left join list_of_devices_comments c on c.id = d.deviceid,
		(select ip, location, nlc, environment, transportmode
		from list_of_scs a
		where changetime=(select max(b.changetime) from list_of_scs b where a.scname=b.scname)
		order by ip asc) s
		where d.scip = s.ip
		order by d.ip asc;`)
	checkErr(err)
	defer r.Close()

	var listOfDevices DeviceList = DeviceList{}
	listOfDevices, _ = getListOfDevices(r, listOfDevices)

	t, _ := template.ParseFiles("static/templates/layoutdiagram.html")
	t.Execute(w, listOfDevices)
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

// FUNCTION DEPRECATED, DO NOT MAKE CHANGES
func insertIntoScs(db *sql.DB, Ip string, Location string, Version string, Nlc int, Scnumber int, Transportmode string, Environment string, Priconc string, Secconc string) {
	// prepare to insert some entries into list_of_scs
	stmt, err := db.Prepare("insert into list_of_scs (Ip, Location, Version, Nlc, Scnumber, Transportmode, Environment, Priconc, Secconc) values (?, ?, ?, ?, ?, ?, ?, ?, ?)")
	checkErr(err)
	result, err := stmt.Exec(Ip, Location, Version, Nlc, Scnumber, Transportmode, Environment, Priconc, Secconc)
	checkErr(err)
	resRows, _ := result.RowsAffected()
	fmt.Println(resRows, "rows affected")
}

// Update records in the db - WIP, require data checking conditions
func updateIntoScs(db *sql.DB, Scname string, Ip string, Location string, Version string, Nlc int, Scnumber int, Transportmode string, Environment string, Priconc string, Secconc string, Changeby string) {
	// prepare to update some entries into list_of_scs
	stmt, err := db.Prepare(`insert into list_of_scs (scname, ip, location, version, nlc, scnumber, transportmode, environment, priconc, secconc, changeby, changetime)
							 values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, now())`)
	checkErr(err)
	result, err := stmt.Exec(Scname, Ip, Location, Version, Nlc, Scnumber, Transportmode, Environment, Priconc, Secconc, Changeby)
	checkErr(err)
	resRows, _ := result.RowsAffected()
	fmt.Println(resRows, "rows affected")
}

// Update List of SCs comments table
func updateIntoScsComments(db *sql.DB, Comment string) {
	// prepare to update some entries into list_of_scs
	stmt, err := db.Prepare(`insert into list_of_scs_comments (comment)
							 values (?)`)
	checkErr(err)
	result, err := stmt.Exec(Comment)
	checkErr(err)
	resRows, _ := result.RowsAffected()
	fmt.Println(resRows, "rows affected")
}

// DEPERECATED FUNCTIONS BELOW, DON'T MESS WITH IT
func insertIntoDevicesRLG(db *sql.DB, Ip string, Scip string, VersionDevice string, VersionRtd string, Plinth int, VersionEprom string, DoppPort string, DoppIp string, Devicetype string) {
	// prepare to insert some device into list_of_devices
	stmt, err := db.Prepare("insert into list_of_devices (Ip, Scip, Version, Version_rtd, Plinth, Version_eprom, Doppport, Doppip, Devicetype) values (?, ?, ?, ?, ?, ?, ?, ?, ?)")
	checkErr(err)
	result, err := stmt.Exec(Ip, Scip, VersionDevice, VersionRtd, Plinth, VersionEprom, DoppPort, DoppIp, Devicetype)
	checkErr(err)
	resRows, _ := result.RowsAffected()
	fmt.Println(resRows, "rows affected")
}

func insertIntoDevicesE2G(db *sql.DB, Ip string, Scip string, VersionDevice string, VersionRtd string, Plinth int, Devicetype string) {
	// prepare to insert some device into list_of_devices
	stmt, err := db.Prepare("insert into list_of_devices (Ip, Scip, Version, Version_rtd, Plinth, Devicetype) values (?, ?, ?, ?, ?, ?)")
	checkErr(err)
	result, err := stmt.Exec(Ip, Scip, VersionDevice, VersionRtd, Plinth, Devicetype)
	checkErr(err)
	resRows, _ := result.RowsAffected()
	fmt.Println(resRows, "rows affected")
}

func insertIntoDevicesFLR(db *sql.DB, Ip string, Scip string, VersionDevice string, VersionRtd string, Plinth int, DoppPort string, DoppIp string, Devicetype string) {
	// prepare to insert some device into list_of_devices
	stmt, err := db.Prepare("insert into list_of_devices (Ip, Scip, Version, Version_rtd, Plinth, Doppport, Doppip, Devicetype) values (?, ?, ?, ?, ?, ?, ?, ?)")
	checkErr(err)
	result, err := stmt.Exec(Ip, Scip, VersionDevice, VersionRtd, Plinth, DoppPort, DoppIp, Devicetype)
	checkErr(err)
	resRows, _ := result.RowsAffected()
	fmt.Println(resRows, "rows affected")
}

// Update Device DB - WIP, require data checking conditions
func updateIntoDevices(db *sql.DB, Devicename string, Ip string, VersionDevice string, VersionRtd string, Devicetype string, DoppIp string, DoppPort string, Plinth int, Scip string, VersionEprom string, Changeby string) {
	// prepare to update some entries into list_of_scs
	stmt, err := db.Prepare(`insert into list_of_devices (devicename, ip, version, version_rtd, devicetype, doppip, doppport, plinth, scip, version_eprom, changeby, changetime)
	values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, now())`)
	checkErr(err)
	result, err := stmt.Exec(Devicename, Ip, VersionDevice, VersionRtd, Devicetype, DoppIp, DoppPort, Plinth, Scip, VersionEprom, Changeby)
	checkErr(err)
	resRows, _ := result.RowsAffected()
	fmt.Println(resRows, "rows affected")
}

func updateIntoDevicesComments(db *sql.DB, Comment string) {
	// prepare to update some entries into list_of_scs
	stmt, err := db.Prepare(`insert into list_of_devices_comments (comment)
	values (?)`)
	checkErr(err)
	result, err := stmt.Exec(Comment)
	checkErr(err)
	resRows, _ := result.RowsAffected()
	fmt.Println(resRows, "rows affected")
}

// E2 Gate function
func updateIntoDevicesE2G(db *sql.DB, Deviceid int, Ip string, Scip string, VersionDevice string, VersionRtd string, Plinth int) {
	// prepare to update some entries into list_of_scs
	stmt, err := db.Prepare("update list_of_devices set Ip=?, Version=?, Version_rtd=?, Plinth=?, Scip=? where Deviceid=?")
	checkErr(err)
	result, err := stmt.Exec(Ip, VersionDevice, VersionRtd, Plinth, Scip, Deviceid)
	checkErr(err)
	resRows, _ := result.RowsAffected()
	fmt.Println(resRows, "rows affected")
}

func updateIntoDevicesRLG(db *sql.DB, Deviceid int, Ip string, Scip string, VersionDevice string, VersionRtd string, Plinth int, VersionEprom string, DoppPort string, DoppIp string) {
	// prepare to update some entries into list_of_scs
	stmt, err := db.Prepare("update list_of_devices set Ip=?, Version=?, Version_rtd=?, Plinth=?, Scip=?, Version_eprom=?, Doppip=?, Doppport=?  where Deviceid=?")
	checkErr(err)
	result, err := stmt.Exec(Ip, VersionDevice, VersionRtd, Plinth, Scip, VersionEprom, DoppIp, DoppPort, Deviceid)
	checkErr(err)
	resRows, _ := result.RowsAffected()
	fmt.Println(resRows, "rows affected")
}

func updateIntoDevicesFLR(db *sql.DB, Deviceid int, Ip string, Scip string, VersionDevice string, VersionRtd string, Plinth int, DoppPort string, DoppIp string) {
	// prepare to update some entries into list_of_scs
	stmt, err := db.Prepare("update list_of_devices set Ip=?, Version=?, Version_rtd=?, Plinth=?, Scip=?, Doppip=?, Doppport=? where Deviceid=?")
	checkErr(err)
	result, err := stmt.Exec(Ip, VersionDevice, VersionRtd, Plinth, Scip, DoppIp, DoppPort, Deviceid)
	checkErr(err)
	resRows, _ := result.RowsAffected()
	fmt.Println(resRows, "rows affected")
}

func getListOfScs(r *sql.Rows, listOfScs []Sc) (listOfScsOut []Sc, err error) {
	for i := 0; r.Next(); i++ {
		newSc := Sc{}
		err = r.Scan(&newSc.Scname, &newSc.Scid, &newSc.Ip, &newSc.Location, &newSc.Version, &newSc.Nlc, &newSc.Scnumber, &newSc.Transportmode, &newSc.Environment, &newSc.Priconc, &newSc.Secconc, &newSc.Changeby, &newSc.Changetime, &newSc.Devicesactive, &newSc.Comment)
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
		err = r.Scan(&newDevice.Devicename, &newDevice.Deviceid, &newDevice.Ip, &newDevice.VersionDevice, &newDevice.VersionRtd, &newDevice.DeviceType,
			&newDevice.DoppIp, &newDevice.DoppPort, &newDevice.Plinth, &newDevice.ParentSc.Ip, &newDevice.VersionEprom,
			&newDevice.ParentSc.Location, &newDevice.ParentSc.Nlc, &newDevice.ParentSc.Environment, &newDevice.ParentSc.Transportmode,
			&newDevice.Changeby, &newDevice.Changetime, &newDevice.Comment)
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

	//http.Handle("/static/assets/", http.StripPrefix("static/assets/", http.FileServer(http.Dir("assets"))))

	//http.HandleFunc("/", handleIndex)
	http.HandleFunc("/itf/", handleListOfAllSCs)
	http.HandleFunc("/itf/devices", handleListOfAllDevices)
	http.HandleFunc("/itf/layout", handleLayout)
	http.HandleFunc("/itf/busrigs", handleListOfAllBusRigs)
	http.ListenAndServe(":8000", nil)
}
