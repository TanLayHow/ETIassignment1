DROP DATABASE ETI_Passengers;
CREATE DATABASE ETI_Passengers;
USE ETI_Passengers;

DROP TABLE Passengers;
CREATE TABLE Passengers (
    PassengerID varchar(255) NOT NULL,
    FirstName varchar(255) NOT NULL,
    LastName varchar(255),
    MobileNumber varchar(255) UNIQUE,
    Email varchar(255) UNIQUE,
    PRIMARY KEY (PassengerID)
);

DROP TABLE PassengerRides;
CREATE TABLE PassengerRides (
    RideID varchar(255) NOT NULL,
    RideDate DateTime,
    PassengerID varchar(255) NOT NULL,
    PickupCode varchar(255),
    DropoffCode varchar(255),
    DriverID varchar(255) NOT NULL,
    CarLicense varchar(255),
    RideStatus varchar(255),
    primary key (RideID),
 foreign key (PassengerID) references Passengers (PassengerID)
); 

DROP DATABASE ETI_Drivers;
CREATE DATABASE ETI_Drivers;
USE ETI_Drivers;

DROP TABLE Drivers;
CREATE TABLE Drivers (
    DriverID varchar(255) NOT NULL,
    FirstName varchar(255) NOT NULL,
    LastName varchar(255),
    MobileNumber varchar(255) UNIQUE,
    Email varchar(255),
    IdentificationNumber varchar(255) NOT NULL UNIQUE,
    CarLicense varchar(255),
    DriverStatus varchar(255),
    PRIMARY KEY (DriverID)
);