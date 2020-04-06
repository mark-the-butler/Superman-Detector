package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
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

func handleLoginRequest(w http.ResponseWriter, r *http.Request) {
	var request loginRequest
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048567)) // Limiting the size of the request body
	if err != nil {
		panic(err)
	}
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

	response := travelResponse{currentGeo{Lat: 39.1702, Lon: -76.8538, Radius: 200}}

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
