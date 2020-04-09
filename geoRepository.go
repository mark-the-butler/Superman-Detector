package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mysteryboy73/Superman-Detector/models"
)

type DataRepo interface {
	getCurrentLocation(login models.LoginRequest) (models.GeoLocation, error)
}

type GeoRepository struct {
}

// GetCurrentLocation retrieves users current location from db
func (gr *GeoRepository) getCurrentLocation(login models.LoginRequest) (models.GeoLocation, error) {
	var currentLocation models.GeoLocation

	database, err := sql.Open("sqlite3", "./db/geolite2.db")
	checkErr(err)

	statement, err := database.Prepare("SELECT latitude, longitude, accuracy_radius FROM blocks WHERE network =?")
	checkErr(err)

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
