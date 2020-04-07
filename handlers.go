package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/mysteryboy73/Superman-Detector/models"
)

// HandleLoginRequest retrieves data for response
func HandleLoginRequest(w http.ResponseWriter, r *http.Request) {
	var request models.LoginRequest

	body, err := ioutil.ReadAll(r.Body)
	CheckErr(err)

	if err := json.Unmarshal(body, &request); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(400)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	response := models.TravelResponse{CurrentLocation: GetCurrentLocation(request)}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}
