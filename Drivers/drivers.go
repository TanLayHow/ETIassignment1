package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// For driver variable
var (
	driverValue string
	driver_list = map[string]Driver{}
)

// Passenger details
type Passenger struct {
	FirstName    string `json:"FirstName"`
	LastName     string `json:"LastName"`
	MobileNumber string `json:"MobileNumber"`
	Email        string `json:"Email"`
}

// Ride details
type Ride struct {
	RideDate      string `json:"RideDate"`
	PassengerID   string `json:"PassengerID"`
	DriverID      string `json:"DriverID"`
	PickUpPostal  string `json:"PickUpPostal"`
	DropOffPostal string `json:"DropOffPostal"`
	CarLicense    string `json:"CarLicense"`
	RideStatus    string `json:"RideStatus"`
}

// Driver details
type Driver struct {
	FirstName            string `json:"FirstName"`
	LastName             string `json:"LastName"`
	MobileNumber         string `json:"MobileNumber"`
	Email                string `json:"Email"`
	IdentificationNumber string `json:"IdentificationNumber"`
	CarLicense           string `json:"CarLicense"`
	DriverStatus         string `json:"DriverStatus"`
}

func main() {
	router := mux.NewRouter()
	// To retrieve drivers
	router.HandleFunc("/api/v1/drivers", getDrivers)
	// To get, create and update drivers
	router.HandleFunc("/api/v1/drivers/{driverid}", driver).Methods("GET", "POST", "PUT")
	// To update status of driver once the trip has been completed
	router.HandleFunc("/api/v1/ride/{driverid}", rideCompleted)
	// To book a driver
	router.HandleFunc("/api/v1/ride", bookDriver)

	// Open port
	fmt.Println("Listening at port 3000")
	log.Fatal(http.ListenAndServe(":3000", router))
}

// To get the list of drivers
func getDrivers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	db, _ := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/ETI_Drivers")
	defer db.Close()
	// Get driver data from database
	getDriverData(db)

	data, err := json.Marshal(map[string]map[string]Driver{"Drivers": driver_list})
	if err != nil {
		log.Fatal(err)
	}

	// To write data to header
	if len(driver_list) != 0 {
		fmt.Fprintf(w, "%s\n", data)
	}
}

// Get driver data from database
func getDriverData(db *sql.DB) {
	results, err := db.Query("SELECT * FROM Drivers")
	if err != nil {
		db.Close()
		panic(err.Error())
	}
	// For each driver in list
	for results.Next() {
		var d Driver
		var driverID string
		err = results.Scan(&driverID, &d.FirstName, &d.LastName, &d.MobileNumber, &d.Email, &d.IdentificationNumber, &d.CarLicense, &d.DriverStatus)
		if err != nil {
			db.Close()
			panic(err.Error())
		}
		// Add driver to driver list
		driver_list[driverID] = d
	}
}

// To get, create and update driver
func driver(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	params := mux.Vars(r)
	driverValue = params["driverid"]

	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/ETI_Drivers")
	if err != nil {
		db.Close()
		panic(err.Error())
	}
	defer db.Close()

	// To get driver data from database
	getDriverData(db)

	driverDetails, exists := driver_list[driverValue]

	if !exists {
		if r.Method == "POST" {
			// To create driver
			newDriver := Driver{}
			reqBody, _ := ioutil.ReadAll(r.Body)
			json.Unmarshal(reqBody, &newDriver)
			_, err := db.Exec("INSERT INTO Drivers (DriverID, FirstName, LastName, MobileNumber, Email, IdentificationNumber, CarLicense, DriverStatus) values(?, ?, ?, ?, ?, ?, ?, ?)", driverValue, newDriver.FirstName, newDriver.LastName, newDriver.MobileNumber, newDriver.Email, newDriver.IdentificationNumber, newDriver.CarLicense, "Available")
			if err != nil {
				panic(err.Error())
			}
			driver_list[driverValue] = newDriver
			w.WriteHeader(http.StatusAccepted)
		}
	} else if exists {
		if r.Method == "GET" {
			// To get driver
			data, _ := json.Marshal(driverDetails)
			fmt.Fprintf(w, "%s\n", data)

		} else if r.Method == "PUT" {
			// To update driver
			updateDriver := Driver{}
			reqBody, _ := ioutil.ReadAll(r.Body)
			json.Unmarshal(reqBody, &updateDriver)
			// Update database
			_, err := db.Exec("UPDATE Drivers SET FirstName=?, LastName=?, MobileNumber=?, Email=?, CarLicense=? WHERE DriverID=?", updateDriver.FirstName, updateDriver.LastName, updateDriver.MobileNumber, updateDriver.Email, updateDriver.CarLicense, driverValue)
			if err != nil {
				panic(err.Error())
			}
			driver_list[driverValue] = updateDriver
			w.WriteHeader(http.StatusAccepted)
		}

	}
}

// To book a driver
func bookDriver(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/ETI_Drivers")
	if err != nil {
		db.Close()
		panic(err.Error())
	}
	defer db.Close()

	// Getting 1 max driver that is available for the passenger
	results, err := db.Query("SELECT * FROM Drivers WHERE DriverStatus = 'Available' ORDER BY RAND() LIMIT 1")
	if err != nil {
		db.Close()
		panic(err.Error())
	}
	for results.Next() {
		var aDriver Driver
		var driverID string
		err = results.Scan(&driverID, &aDriver.FirstName, &aDriver.LastName, &aDriver.MobileNumber, &aDriver.Email, &aDriver.IdentificationNumber, &aDriver.CarLicense, &aDriver.DriverStatus)
		if err != nil {
			db.Close()
			panic(err.Error())
		}

		// Set Driver Status to Busy
		_, err = db.Exec("UPDATE Drivers SET DriverStatus='Active' WHERE DriverID=?", driverID)
		if err != nil {
			panic(err.Error())
		}
		aDriver.DriverStatus = "Active"

		// Write data to header
		data, err := json.Marshal(map[string]Driver{driverID: aDriver})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, "%s\n", data)
	}
}

// To update status of driver once the trip has been completed
func rideCompleted(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	params := mux.Vars(r)
	driverValue = params["driverid"]

	// Update database
	db, _ := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/ETI_Drivers")
	_, _ = db.Exec("UPDATE Drivers SET DriverStatus='Available' WHERE DriverID=?", driverValue)
}
