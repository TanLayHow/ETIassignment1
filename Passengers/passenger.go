package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// For passenger, ride variables
var (
	passengerValue  string
	passenger_list  = map[string]Passenger{}
	rides_list      = map[string]Ride{}
	driverRidesList = map[string]Ride{}
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
	// To retrieve passengers
	router.HandleFunc("/api/v1/passengers", getPassengers)
	// To get, create and update passengers
	router.HandleFunc("/api/v1/passengers/{passengerid}", newPassenger).Methods("GET", "POST", "PUT")
	// To create a ride booking for the passenger
	router.HandleFunc("/api/v1/passengers/{passengerid}/ride/{pickuppostal}/{dropoffpostal}", booking)
	// To look at passenger's ride history
	router.HandleFunc("/api/v1/passengers/{passengerid}/history", passengerHistory)
	// To obtain list of driver rides for drivers to start ride
	router.HandleFunc("/api/v1/ride/{driverid}/driver/{ridestatus}", driverRideStart)
	// To obtain list of driver rides for drivers to end ride
	router.HandleFunc("/api/v1/ride/{rideid}/{driverid}/{status}", driverRidesEnd)

	// Open ports
	fmt.Println("Listening at port 5000")
	log.Fatal(http.ListenAndServe(":5000", router))
}

// Get list of passengers
func getPassengers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/ETI_Passengers")
	if err != nil {
		db.Close()
		panic(err.Error())
	}
	defer db.Close()

	// Get passenger data from database
	getPassengerData(db)
	data, err := json.Marshal(map[string]map[string]Passenger{"Passengers": passenger_list})
	if err != nil {
		log.Fatal(err)
	}
	// To write data to header
	if len(passenger_list) != 0 {
		fmt.Fprintf(w, "%s\n", data)
	}
}

// Get passenger data from database
func getPassengerData(db *sql.DB) {
	results, err := db.Query("SELECT * FROM Passengers")
	if err != nil {
		db.Close()
		panic(err.Error())
	}
	// Add every passenger into a list
	for results.Next() {
		var passenger Passenger
		var passengerID string
		err = results.Scan(&passengerID, &passenger.FirstName, &passenger.LastName, &passenger.MobileNumber, &passenger.Email)
		if err != nil {
			db.Close()
			panic(err.Error())
		}
		// Add passenger to passenger list
		passenger_list[passengerID] = passenger
	}
}

// Create/Update passenger
func newPassenger(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	params := mux.Vars(r)
	passengerValue = params["passengerid"]

	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/ETI_Passengers")
	if err != nil {
		db.Close()
		panic(err.Error())
	}
	defer db.Close()

	// Get passenger data to create/update passenger
	getPassengerData(db)

	passengerDetails, exists := passenger_list[passengerValue]

	if !exists {
		// To create passenger
		if r.Method == "POST" {
			newPassenger := Passenger{}
			reqBody, _ := ioutil.ReadAll(r.Body)
			json.Unmarshal(reqBody, &newPassenger)

			_, err := db.Exec("INSERT INTO Passengers (PassengerID, FirstName, LastName, MobileNumber, Email) values(?, ?, ?, ?, ?)", passengerValue, newPassenger.FirstName, newPassenger.LastName, newPassenger.MobileNumber, newPassenger.Email)
			if err != nil {
				panic(err.Error())
			}
			passenger_list[passengerValue] = newPassenger
			w.WriteHeader(http.StatusAccepted)
		}
	} else if exists {
		// To get passenger
		if r.Method == "GET" {
			data, _ := json.Marshal(passengerDetails)
			fmt.Fprintf(w, "%s\n", data)

		} else if r.Method == "PUT" {
			// To update passenger
			updatePassenger := Passenger{}
			reqBody, _ := ioutil.ReadAll(r.Body)
			json.Unmarshal(reqBody, &updatePassenger)
			_, err := db.Exec("UPDATE Passengers SET FirstName=?, LastName=?, MobileNumber=?, Email=? WHERE PassengerID=?", updatePassenger.FirstName, updatePassenger.LastName, updatePassenger.MobileNumber, updatePassenger.Email, passengerValue)
			if err != nil {
				panic(err.Error())
			}
			passenger_list[passengerValue] = updatePassenger
			w.WriteHeader(http.StatusAccepted)
		}
	}
}

// To create a ride booking for the passenger
func booking(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	params := mux.Vars(r)
	// To store parameters from the header
	passengerID := params["passengerid"]
	pickUpPostal := params["pickuppostal"]
	dropOffPostal := params["dropoffpostal"]

	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/ETI_Passengers")
	if err != nil {
		db.Close()
		panic(err.Error())
	}
	defer db.Close()

	// To get passenger details for booking
	getPassengerData(db)

	_, exists := passenger_list[passengerID]
	if exists {
		driverList := map[string]Driver{}
		// To book a driver, by changing status of driver
		resp, _ := http.Get("http://localhost:3000/api/v1/ride")

		body, _ := ioutil.ReadAll(resp.Body)

		json.Unmarshal([]byte(body), &driverList)

		var bookedDriverID string
		var bookedDriver Driver

		for key, driver := range driverList {
			bookedDriverID = key
			bookedDriver = driver
		}

		generatedRideID := "RID" + time.Now().Format("15052007250513")
		// Create PassengerRide with PassengerID and DriverID
		_, _ = db.Exec("INSERT INTO PassengerRides (RideID, RideDate, PassengerID, DriverID, PickUpPostal, DropOffPostal, CarLicense, RideStatus) values(?, ?, ?, ?, ?, ?, ?, ?)", generatedRideID, time.Now(), passengerID, bookedDriverID, pickUpPostal, dropOffPostal, bookedDriver.CarLicense, "Awaiting")

		// Write data to header
		data, _ := json.Marshal(map[string]Driver{bookedDriverID: bookedDriver})
		fmt.Fprintf(w, "%s\n", data)

	}
}

// Get list of passenger past rides
func passengerHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	params := mux.Vars(r)
	// Retrieves the passengerid which will be used to retrieve the passenger rides
	passengerID := params["passengerid"]

	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/ETI_Passengers")
	defer db.Close()

	rides_list = map[string]Ride{}

	rides, err := db.Query("SELECT * FROM PassengerRides WHERE PassengerID = ? ORDER BY RideDate DESC", passengerID)

	// Make list of rides with passengerID
	for rides.Next() {
		var ride Ride
		var rideID string
		err = rides.Scan(&rideID, &ride.RideDate, &ride.PassengerID, &ride.DriverID, &ride.PickUpPostal, &ride.DropOffPostal, &ride.CarLicense, &ride.RideStatus)
		rides_list[rideID] = ride
	}

	// Write data to header
	data, err := json.Marshal(map[string]map[string]Ride{"Rides": rides_list})
	if err != nil {
		log.Fatal(err)
	}
	if len(rides_list) != 0 {
		fmt.Fprintf(w, "%s\n", data)
	}

}

// To create passengerRide with driverID
func driverRideStart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	params := mux.Vars(r)
	driverID := params["driverid"]
	ridestatus := params["ridestatus"]

	db, _ := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/ETI_Passengers")

	driverRidesList = map[string]Ride{}

	results, err := db.Query("SELECT * FROM PassengerRides WHERE DriverID = ? AND RideStatus = ? ORDER BY RideDate DESC", driverID, ridestatus)

	for results.Next() {
		var ride Ride
		var rideID string
		err = results.Scan(&rideID, &ride.RideDate, &ride.PassengerID, &ride.DriverID, &ride.PickUpPostal, &ride.DropOffPostal, &ride.CarLicense, &ride.RideStatus)
		driverRidesList[rideID] = ride
	}

	// Write data to header
	data, err := json.Marshal(map[string]map[string]Ride{"Rides": driverRidesList})
	if err != nil {
		log.Fatal(err)
	}
	if len(driverRidesList) != 0 {
		fmt.Fprintf(w, "%s\n", data)
	}
}

// To end a ride by changing ridestatus
func driverRidesEnd(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	params := mux.Vars(r)

	rideID := params["rideid"]
	driverID := params["driverid"]
	status := params["status"]

	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/ETI_Passengers")
	if err != nil {
		db.Close()
		panic(err.Error())
	}
	defer db.Close()

	// Update PassengerRide status to Done
	_, err = db.Exec("UPDATE PassengerRides SET RideStatus=? WHERE RideID=?", status, rideID)
	driverRidesList = map[string]Ride{}
	results, err := db.Query("SELECT * FROM PassengerRides WHERE DriverID = ? AND RideStatus = ? ORDER BY RideDate DESC", driverID, status)

	for results.Next() {
		var ride Ride
		var rideID string
		err = results.Scan(&rideID, &ride.RideDate, &ride.PassengerID, &ride.DriverID, &ride.PickUpPostal, &ride.DropOffPostal, &ride.CarLicense, &ride.RideStatus)

		driverRidesList[rideID] = ride
	}

	// Update Driver status
	if strings.TrimSpace(status) == "Done" {
		_, err := http.Get("http://localhost:3000/api/v1/ride/" + driverID)
		if err != nil {
			log.Fatalln(err)
		}
	}

	// Write data to header
	data, _ := json.Marshal(map[string]map[string]Ride{"Rides": driverRidesList})

	if len(driverRidesList) != 0 {
		fmt.Fprintf(w, "%s\n", data)
	}
}
