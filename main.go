package main

import (
	"database/sql"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mysteryboy73/Superman-Detector/models"
)

func getCurrentLocation(login models.LoginRequest) models.CurrentGeo {
	var currentLocation models.CurrentGeo

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
	return currentLocation
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func handleLoginRequest(w http.ResponseWriter, r *http.Request) {
	var request models.LoginRequest
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048567)) // Limiting the size of the request body
	checkErr(err)
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &request); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(400)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	response := models.TravelResponse{CurrentLocation: getCurrentLocation(request)}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}

func main() {
	http.HandleFunc("/loginRequest", handleLoginRequest)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
