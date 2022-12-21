package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// For Passenger and Driver variables
var (
	passenger_list map[string]map[string]Passenger
	driver_list    map[string]map[string]Driver
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
outer:
	for {
		// Show menu of program
		fmt.Println("----------------")
		fmt.Println("|  Drive N Go  |")
		fmt.Println("|     Menu     |")
		fmt.Println("----------------")
		fmt.Println("[1] Log In")
		fmt.Println("[2] Account Creation")
		fmt.Println("[0] Quit")

		/* var choice int
		fmt.Scanf("%d", &choice) */

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("\nEnter an option: ")
		input, _ := reader.ReadString('\n')
		option, _ := strconv.Atoi(strings.TrimSpace(input))

		// Options of program
		switch option {
		case 0:
			fmt.Println("Program Exited")
			break outer
		case 1:
			Login()
		case 2:
			AccountCreation()
		default:
			fmt.Println("Please Select A Listed Option")
		}
	}
}

func Login() {
outer:
	for {
		// Login options, with roles
		fmt.Println("----------------")
		fmt.Println("Select Login Type")
		fmt.Println("----------------")
		fmt.Println("[1] Passenger")
		fmt.Println("[2] Driver")
		fmt.Println("[0] Exit")

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("\nEnter an option: ")
		userInput, _ := reader.ReadString('\n')
		option, _ := strconv.Atoi(strings.TrimSpace(userInput))

		switch option {
		case 0:
			break outer
		case 1:
			// Passengers
			resp, err := http.Get("http://localhost:5000/api/v1/passengers")
			if err != nil {
				fmt.Println("Passenger list was unable to be obtained.")
			}
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatalln(err)
			}
			// Read JSON into usable data
			json.Unmarshal([]byte(body), &passenger_list)

			// Check for passenger accounts
			if len(passenger_list) != 0 {
				// List passenger list
				fmt.Println("----------------")
				fmt.Print("\nPassenger List")

				for _, item := range passenger_list["Passengers"] {
					fmt.Println(item.FirstName + " " + item.LastName + ", " + item.MobileNumber + ", " + item.Email)
				}

				// Ask user for Email input
				reader := bufio.NewReader(os.Stdin)
				fmt.Print("\nEnter Email: ")
				userInput, _ := reader.ReadString('\n')
				var userEmail = strings.TrimSpace(userInput)
				var foundUser bool = false
				var passenger Passenger
				var passengerID string

				// To find if email matches any passenger email
				for key, item := range passenger_list["Passengers"] {
					if userEmail != item.Email {
					} else {
						foundUser = true
						passenger = item
						passengerID = key
					}
				}

				// If user is found, log in
				if foundUser == true {
					fmt.Println("Successfully logged in.")
					PassengerMenu(passengerID, passenger)
				} else {
					fmt.Println("User not found.")
					continue
				}
			}
		case 2:
			// Drivers
			resp, err := http.Get("http://localhost:3000/api/v1/drivers")
			if err != nil {
				fmt.Println("Driver list was unable to be obtained.")
			}
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatalln(err)
			}
			// Read JSON into usable data
			json.Unmarshal([]byte(body), &driver_list)

			// Check for driver accounts
			if len(driver_list) != 0 {
				// List driver list
				fmt.Println("----------------")
				fmt.Print("\nDriver List")

				for _, item := range driver_list["Drivers"] {
					fmt.Println(item.FirstName + " " + item.LastName + ", " + item.MobileNumber + ", " + item.Email)
				}

				// Ask user for email input
				reader := bufio.NewReader(os.Stdin)
				fmt.Print("\nEnter Email: ")
				userInput, _ := reader.ReadString('\n')
				var userEmail = strings.TrimSpace(userInput)
				var foundUser bool = false
				var driver Driver
				var driverID string

				// To match if email is in driver email list
				for key, item := range driver_list["Drivers"] {
					if userEmail != item.Email {
					} else {
						foundUser = true
						driver = item
						driverID = key
					}
				}

				// If email is found, log in
				if foundUser == true {
					fmt.Println("Successfully logged in.")
					DriverMenu(driverID, driver)
				} else {
					fmt.Println("User not found.")
					continue
				}
			}
		default:
			fmt.Println("Please Select A Listed Option")
		}
	}
}

func AccountCreation() {
outer:
	for {
		// Create account menu to choose roles
		fmt.Println("----------------")
		fmt.Println("Choose account role")
		fmt.Println("[1] Passenger")
		fmt.Println("[2] Driver")
		fmt.Println("[0] Exit")

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("\nEnter an option: ")
		userInput, _ := reader.ReadString('\n')
		option, _ := strconv.Atoi(strings.TrimSpace(userInput))

		switch option {
		case 0:
			break outer
		case 1:
			// Create passenger
			newPassenger := Passenger{}
			var newString string

			// Ask user to input passenger details
			reader = bufio.NewReader(os.Stdin)
			fmt.Println("Enter First Name: ")
			inputString, _ := reader.ReadString('\n')
			newString = strings.TrimSpace(inputString)
			newPassenger.FirstName = newString

			reader = bufio.NewReader(os.Stdin)
			fmt.Println("Enter Last Name: ")
			inputString, _ = reader.ReadString('\n')
			newString = strings.TrimSpace(inputString)
			newPassenger.LastName = newString

			reader = bufio.NewReader(os.Stdin)
			fmt.Println("Enter Mobile Number: ")
			inputString, _ = reader.ReadString('\n')
			newString = strings.TrimSpace(inputString)
			newPassenger.MobileNumber = newString

			reader = bufio.NewReader(os.Stdin)
			fmt.Println("Enter Email: ")
			inputString, _ = reader.ReadString('\n')
			newString = strings.TrimSpace(inputString)
			newPassenger.Email = newString

			// Store the data in JSON format to send to server
			jsonBody, _ := json.Marshal(newPassenger)

			// Create unique ID
			passengerID := "PID" + time.Now().Format("15052007250513")

			// Create new passenger and submit request to server
			client := &http.Client{}
			if req, err := http.NewRequest("POST", "http://localhost:5000/api/v1/passengers/"+passengerID, bytes.NewBuffer(jsonBody)); err == nil {
				if _, err := client.Do(req); err == nil {
					fmt.Println("\nPassenger " + newPassenger.FirstName + " " + newPassenger.LastName + " Created.")
					break outer
				} else {
					fmt.Println("\nPassenger account failed to be created.")
					continue
				}
			}

		case 2:
			// Create driver
			newDriver := Driver{}
			var newString string

			// Ask user to input driver details
			reader = bufio.NewReader(os.Stdin)
			fmt.Println("Enter First Name: ")
			inputString, _ := reader.ReadString('\n')
			newString = strings.TrimSpace(inputString)
			newDriver.FirstName = newString

			reader = bufio.NewReader(os.Stdin)
			fmt.Println("Enter Last Name: ")
			inputString, _ = reader.ReadString('\n')
			newString = strings.TrimSpace(inputString)
			newDriver.LastName = newString

			reader = bufio.NewReader(os.Stdin)
			fmt.Println("Enter Mobile Number: ")
			inputString, _ = reader.ReadString('\n')
			newString = strings.TrimSpace(inputString)
			newDriver.MobileNumber = newString

			reader = bufio.NewReader(os.Stdin)
			fmt.Println("Enter Email: ")
			inputString, _ = reader.ReadString('\n')
			newString = strings.TrimSpace(inputString)
			newDriver.Email = newString

			reader = bufio.NewReader(os.Stdin)
			fmt.Println("Enter Car License: ")
			inputString, _ = reader.ReadString('\n')
			newString = strings.TrimSpace(inputString)
			newDriver.CarLicense = newString

			// Store the data in JSON format to send to server
			jsonBody, _ := json.Marshal(newDriver)

			// Create unique ID
			driverID := "DID" + time.Now().Format("15052007250513")

			// Create new driver and submit request to server
			client := &http.Client{}
			if req, err := http.NewRequest("POST", "http://localhost:3000/api/v1/drivers/"+driverID, bytes.NewBuffer(jsonBody)); err == nil {
				if _, err := client.Do(req); err == nil {
					fmt.Println("\nDriver " + newDriver.FirstName + " " + newDriver.LastName + " Created.")
					break outer
				} else {
					fmt.Println("\nDriver account failed to be created.")
					continue
				}
			}
		default:
			fmt.Println("Please Select A Listed Option")
		}
	}
}

// List passenger menu
func PassengerMenu(passengerID string, p Passenger) {
outer:
	for {
		// Passenger menu to show what passenger can do
		fmt.Println("----------------")
		fmt.Println("Passenger Menu")
		fmt.Println("[1] Update Account Details")
		fmt.Println("[2] Book Ride")
		fmt.Println("[3] View Rides")
		fmt.Println("[0] Exit")

		// Read input for menu
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("\nEnter an option: ")
		userInput, _ := reader.ReadString('\n')
		option, _ := strconv.Atoi(strings.TrimSpace(userInput))

		// menu
		switch option {
		case 0:
			break outer
		case 1:
			// Ask user to input passenger details
			updatedPassenger := Passenger{}
			var newString string

			reader = bufio.NewReader(os.Stdin)
			fmt.Println("Enter First Name: ")
			inputString, _ := reader.ReadString('\n')
			newString = strings.TrimSpace(inputString)
			updatedPassenger.FirstName = newString

			reader = bufio.NewReader(os.Stdin)
			fmt.Println("Enter Last Name: ")
			inputString, _ = reader.ReadString('\n')
			newString = strings.TrimSpace(inputString)
			updatedPassenger.LastName = newString

			reader = bufio.NewReader(os.Stdin)
			fmt.Println("Enter Mobile Number: ")
			inputString, _ = reader.ReadString('\n')
			newString = strings.TrimSpace(inputString)
			updatedPassenger.MobileNumber = newString

			reader = bufio.NewReader(os.Stdin)
			fmt.Println("Enter Email: ")
			inputString, _ = reader.ReadString('\n')
			newString = strings.TrimSpace(inputString)
			updatedPassenger.Email = newString

			// Store the data in JSON format to send to server
			jsonBody, _ := json.Marshal(updatedPassenger)

			// Submit request to server
			client := &http.Client{}
			if req, err := http.NewRequest("PUT", "http://localhost:5000/api/v1/passengers/"+passengerID, bytes.NewBuffer(jsonBody)); err == nil {
				if _, err := client.Do(req); err == nil {
					fmt.Print("\nPassenger " + updatedPassenger.FirstName + " was updated.")
				}
			}

		case 2:
			// Book a ride
			// Ask user to input ride details
			var stringPickUpPostal string
			reader := bufio.NewReader(os.Stdin)
			fmt.Println("Please input the pickup postal code: ")
			pickupPostal, _ := reader.ReadString('\n')
			stringPickUpPostal = strings.TrimSpace(pickupPostal)

			var stringDropOffPostal string
			reader = bufio.NewReader(os.Stdin)
			fmt.Println("Input the Drop-Off Postal Code: ")
			dropoffPostal, _ := reader.ReadString('\n')
			stringDropOffPostal = strings.TrimSpace(dropoffPostal)

			// Get a driver randomly
			var assignedDriver map[string]Driver
			resp, err := http.Get("http://localhost:5000/api/v1/passengers/" + passengerID + "/ride/" + stringPickUpPostal + "/" + stringDropOffPostal)

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatalln(err)
			}
			// Read JSON into usable data
			json.Unmarshal([]byte(body), &assignedDriver)

			fmt.Println("Searching For a Driver...\n ")

			// Check if there is an available driver to assign to ride
			if len(assignedDriver) != 0 {
				fmt.Println("Found a driver for you!")
				for _, driver := range assignedDriver {
					fmt.Println("Driver Name: " + driver.FirstName + " " + driver.LastName)
					fmt.Println("Mobile Number: " + driver.MobileNumber)
					fmt.Println("Car License: " + driver.CarLicense)
					fmt.Println("Pick-Up Postal Code: " + stringPickUpPostal)
					fmt.Println("Drop-Off Postal Code: " + stringDropOffPostal)
				}
			} else {
				fmt.Println("Can't find an available driver.")
				continue
			}

		case 3:
			// View ride past history
			resp, _ := http.Get("http://localhost:5000/api/v1/passengers/" + passengerID + "/history")

			body, _ := ioutil.ReadAll(resp.Body)

			var bookingHistory map[string]map[string]Ride
			// Read JSON into usable data
			json.Unmarshal([]byte(body), &bookingHistory)

			// If user has booked before
			if len(bookingHistory) != 0 {
				fmt.Println("Ride History: " + p.FirstName + " " + p.LastName)

				// Print ride details
				for _, item := range bookingHistory["Rides"] {
					fmt.Println("Ride Date: " + item.RideDate)
					fmt.Println("Pick-Up Postal Code: " + item.PickUpPostal)
					fmt.Println("Drop-Off Postal Code: " + item.DropOffPostal)
					fmt.Println("Car License: " + item.CarLicense)
				}

			} else {
				// If there is no ride made yet
				fmt.Println("\nYou Have No Rides In Your Ride History!")
			}

		default:
			fmt.Println("Please Select A Listed Option")
		}

	}
}

// List driver menu
func DriverMenu(driverID string, d Driver) {
outer:
	for {
		// Choose driver menu
		fmt.Println("----------------")
		fmt.Println("Driver Menu")
		fmt.Println("[1] Update Account Details")
		fmt.Println("[2] Start trip")
		fmt.Println("[3] End trip")
		fmt.Println("[0] Exit")

		// Read driver choice
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("\nEnter an option: ")
		userInput, _ := reader.ReadString('\n')
		option, _ := strconv.Atoi(strings.TrimSpace(userInput))

		switch option {
		case 0:
			break outer
		case 1:
			// Creating driver
			// Ask user to input driver details
			updatedDriver := Driver{}
			var newString string

			reader = bufio.NewReader(os.Stdin)
			fmt.Println("Enter First Name: ")
			inputString, _ := reader.ReadString('\n')
			newString = strings.TrimSpace(inputString)
			updatedDriver.FirstName = newString

			reader = bufio.NewReader(os.Stdin)
			fmt.Println("Enter Last Name: ")
			inputString, _ = reader.ReadString('\n')
			newString = strings.TrimSpace(inputString)
			updatedDriver.LastName = newString

			reader = bufio.NewReader(os.Stdin)
			fmt.Println("Enter Mobile Number: ")
			inputString, _ = reader.ReadString('\n')
			newString = strings.TrimSpace(inputString)
			updatedDriver.MobileNumber = newString

			reader = bufio.NewReader(os.Stdin)
			fmt.Println("Enter Email: ")
			inputString, _ = reader.ReadString('\n')
			newString = strings.TrimSpace(inputString)
			updatedDriver.Email = newString

			reader = bufio.NewReader(os.Stdin)
			fmt.Println("Enter Car License: ")
			inputString, _ = reader.ReadString('\n')
			newString = strings.TrimSpace(inputString)
			updatedDriver.CarLicense = newString

			// Store the data in JSON format to send to server
			jsonBody, _ := json.Marshal(updatedDriver)

			// Submit request to server to update driver
			client := &http.Client{}
			if req, err := http.NewRequest("PUT", "http://localhost:3000/api/v1/drivers/"+driverID, bytes.NewBuffer(jsonBody)); err == nil {
				if _, err := client.Do(req); err == nil {
					fmt.Print("\nDriver " + updatedDriver.FirstName + " was updated.")
				}
			}

		case 2:
			// Starting a trip
			// Assigning driver to trip
			resp, _ := http.Get("http://localhost:5000/api/v1/ride/" + driverID + "/driver/Awaiting")

			body, _ := ioutil.ReadAll(resp.Body)

			rideList := map[string]map[string]Ride{}
			// Read JSON into usable data
			json.Unmarshal([]byte(body), &rideList)

			// If there are rides for the driver
			if len(rideList["Rides"]) != 0 {
				// Print ride information for driver
				var rideID string
				for key, item := range rideList["Rides"] {
					rideID = key
					fmt.Println("Ride Date: " + item.RideDate)
					fmt.Println("Pick-Up Postal Code: " + item.PickUpPostal)
					fmt.Println("Drop-Off Postal Code: " + item.DropOffPostal)
					fmt.Println("Car License: " + item.CarLicense)
				}

				// Option for driver to start trip
				var stringInput string
				reader = bufio.NewReader(os.Stdin)
				fmt.Println("Ready for trip?")
				fmt.Println("Yes/No")
				input, _ := reader.ReadString('\n')
				stringInput = strings.TrimSpace(input)

				// Driver starts trip
				if stringInput == "Yes" {
					resp, _ = http.Get("http://localhost:5000/api/v1/ride/" + rideID + "/" + driverID + "/Active")

					body, _ = ioutil.ReadAll(resp.Body)

					rideList = map[string]map[string]Ride{}
					// Read JSON into usable data
					json.Unmarshal([]byte(body), &rideList)

					fmt.Println("Ride started.")
				} else if stringInput == "No" {
					// Driver does not start trip
					fmt.Println("Ride not started...")
				} else {
					// User inputs something else than Yes/No
					fmt.Println("Please enter a correct input.")
				}
			} else {
				// There are no rides for the driver
				fmt.Println("\nThere are no trips available.")
				continue
			}
		case 3:
			// Ending a trip
			resp, _ := http.Get("http://localhost:5000/api/v1/ride/" + driverID + "/driver/Active")

			body, _ := ioutil.ReadAll(resp.Body)

			rideList := map[string]map[string]Ride{}
			// Read JSON into usable data
			json.Unmarshal([]byte(body), &rideList)

			// If there is a ride, display it
			if len(rideList["Rides"]) != 0 {
				var rideID string
				for key, item := range rideList["Rides"] {
					rideID = key
					fmt.Println("Ride Date: " + item.RideDate)
					fmt.Println("Pick-Up Postal Code: " + item.PickUpPostal)
					fmt.Println("Drop-Off Postal Code: " + item.DropOffPostal)
					fmt.Println("Car License: " + item.CarLicense)
				}

				// Option to end ride
				var stringInput string
				reader = bufio.NewReader(os.Stdin)
				fmt.Println("End trip?")
				fmt.Println("Yes/No")
				input, _ := reader.ReadString('\n')
				stringInput = strings.TrimSpace(input)

				// Driver ends ride
				if stringInput == "Yes" {
					resp, _ = http.Get("http://localhost:5000/api/v1/ride/" + rideID + "/" + driverID + "/Done")

					body, _ = ioutil.ReadAll(resp.Body)

					rideList = map[string]map[string]Ride{}

					// Read JSON into usable data
					json.Unmarshal([]byte(body), &rideList)

					fmt.Println("Ride ended.")
				} else if stringInput == "No" {
					// Driver does not end ride
					fmt.Println("Ride not ended...")
				} else {
					// Other inputs
					fmt.Println("Please enter a correct input.")
				}

			} else {
				fmt.Println("There are no trips to end.")
				continue
			}
		default:
			fmt.Println("Please Select A Listed Option")
		}

	}
}
