package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mysteryboy73/Superman-Detector/models"
)

// DataRepo is an interface for retrieving data from a db
type DataRepo interface {
	getLocation(login models.LoginRequest) (models.GeoLocation, error)
	saveLogin(login models.LoginRequest) bool
}

// GeoRepository implements the DataRepo interface for retrieving data from a db
type GeoRepository struct {
}

func (gr *GeoRepository) saveLogin(login models.LoginRequest) bool {
	return true
}

// GetCurrentLocation retrieves users current location from db
func (gr *GeoRepository) getLocation(login models.LoginRequest) (models.GeoLocation, error) {
	var currentLocation models.GeoLocation

	database, err := sql.Open("sqlite3", "./db/geolite2.db")
	checkErr(err)

	statement, err := database.Prepare("SELECT latitude, longitude, accuracy_radius FROM blocks WHERE network =?")
	checkErr(err)

	statement.Exec()

	rows, err := statement.Query(login.IPAddress)
	checkErr(err)

	if rows.Next() {
		err = rows.Scan(&currentLocation.Lat, &currentLocation.Lon, &currentLocation.Radius)
		checkErr(err)
	}

	rows.Close()
	return currentLocation, nil
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
