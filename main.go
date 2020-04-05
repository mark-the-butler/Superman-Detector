package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type travel struct {
	Lat    float64 `json:"lat"`
	Lon    float64 `json:"lon"`
	Radius int     `json:"radius"`
}

func handleLoginRequest(w http.ResponseWriter, r *http.Request) {
	response := travel{Lat: 39.1702, Lon: -76.8538, Radius: 200}

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
