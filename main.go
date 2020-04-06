package main

import (
	"database/sql"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type currentGeo struct {
	Lat    float64 `json:"lat"`
	Lon    float64 `json:"lon"`
	Radius int     `json:"radius"`
}

type travelResponse struct {
	Current currentGeo `json:"currentGeo"`
}

type loginRequest struct {
	Username      string `json:"username"`
	UnixTimestamp int    `json:"unix_timestamp"`
	EventUUID     string `json:"event_uuid"`
	IPAddress     string `json:"ip_address"`
}

func getCurrentLocation(login loginRequest) currentGeo {
	var currentLocation currentGeo

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
	var request loginRequest
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

	response := getCurrentLocation(request)

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
