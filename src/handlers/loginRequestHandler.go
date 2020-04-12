package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/mysteryboy73/Superman-Detector/builders"
	"github.com/mysteryboy73/Superman-Detector/models"
)

// LoginRequestHandler takes a builder and returns a response
type LoginRequestHandler struct {
	responseBuilder builders.ResponseBuilder
}

// NewLoginRequestHandler creates loginRequestHandler with default responseBuilder
func NewLoginRequestHandler() *LoginRequestHandler {
	loginRequestHandler := LoginRequestHandler{responseBuilder: builders.NewLoginResponseBuilder()}
	return &loginRequestHandler
}

// HandleLoginRequest retrieves data for response
func (lrh *LoginRequestHandler) HandleLoginRequest(w http.ResponseWriter, r *http.Request) {
	var request models.LoginRequest

	// Decode the json response body into a LoginRequest model
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		return
	}

	// Call the builder to build the response
	response, err := lrh.responseBuilder.Build(request)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		return
	}

	// Send back the response in json format
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}
