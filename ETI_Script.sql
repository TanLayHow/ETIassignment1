/* Passengers */
DROP DATABASE IF EXISTS ETI_Passengers;
CREATE DATABASE ETI_Passengers;
USE ETI_Passengers;
DROP TABLE IF EXISTS Passengers;
CREATE TABLE Passengers (
    PassengerID VARCHAR(255) NOT NULL,
    FirstName VARCHAR(255) NOT NULL,
    LastName VARCHAR(255),
    MobileNumber VARCHAR(255) UNIQUE,
    Email VARCHAR(255) UNIQUE,
    PRIMARY KEY (PassengerID)
);
INSERT INTO Passengers (PassengerID, FirstName, LastName, MobileNumber, Email) VALUES("PID12345678", "Tom", "William", "88096390", "TomWilliam@gmail.com");
INSERT INTO Passengers (PassengerID, FirstName, LastName, MobileNumber, Email) VALUES("PID55566688", "Daniel", "Dan", "88223344", "DanielDan@gmail.com");
INSERT INTO Passengers (PassengerID, FirstName, LastName, MobileNumber, Email) VALUES("PID98712388", "John", "Doe", "91192628", "JohnDoe@gmail.com");
SELECT * FROM Passengers;
/* Passenger Rides */
DROP TABLE IF EXISTS PassengerRides;
CREATE TABLE PassengerRides (
    RideID VARCHAR(255) NOT NULL,
    RideDate DATETIME,
    PassengerID VARCHAR(255) NOT NULL,
	DriverID VARCHAR(255) NOT NULL,
    PickUpPostal VARCHAR(255),
    DropOffPostal VARCHAR(255),
    CarLicense VARCHAR(255),
    RideStatus VARCHAR(255),
    PRIMARY KEY (RideID),
	FOREIGN KEY (PassengerID) REFERENCES Passengers (PassengerID)
);	

/* USE ETI_Passengers;
SELECT * FROM Passengers
SELECT * FROM PassengerRides; */

/* Drivers */
DROP DATABASE IF EXISTS ETI_Drivers;
CREATE DATABASE ETI_Drivers;
USE ETI_Drivers;
DROP TABLE IF EXISTS Drivers;
CREATE TABLE Drivers (
    DriverID VARCHAR(255) NOT NULL,
    FirstName VARCHAR(255) NOT NULL,
    LastName VARCHAR(255),
    MobileNumber VARCHAR(255) UNIQUE,
    Email VARCHAR(255) UNIQUE,
    IdentificationNumber VARCHAR(255) NOT NULL UNIQUE,
    CarLicense VARCHAR(255),
    DriverStatus VARCHAR(255),
    PRIMARY KEY (DriverID)
);
INSERT INTO Drivers (DriverID, FirstName, LastName, MobileNumber, Email, IdentificationNumber, CarLicense, DriverStatus) VALUES("DID44448888", "Rheanna", "Yong", "93825712", "RheannaYong@gmail.com", "T0318985B", "L34567", 'Available');
INSERT INTO Drivers (DriverID, FirstName, LastName, MobileNumber, Email, IdentificationNumber, CarLicense, DriverStatus) VALUES("DID33337777", "Mongus", "Luke", "32148591", "MongusLuke@gmail.com", "T0318484A", "L98765", 'Available');
INSERT INTO Drivers (DriverID, FirstName, LastName, MobileNumber, Email, IdentificationNumber, CarLicense, DriverStatus) VALUES("DID12345679", "Pom", "Pom", "98375822", "PomPom@gmail.com", "S349143T", "L12345", 'Available');

/* USE ETI_Driver;
SELECT * FROM Drivers; */
