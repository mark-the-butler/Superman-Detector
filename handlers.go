package main

import (
	"encoding/json"
	"net/http"

	"github.com/mysteryboy73/Superman-Detector/models"
)

// LoginRequestHandler takes a builder and returns a response
type LoginRequestHandler struct {
	responseBuilder ResponseBuilder
}

// NewLoginRequestHandler creates loginRequestHandler with default responseBuilder
func NewLoginRequestHandler() *LoginRequestHandler {
	loginRequestHandler := LoginRequestHandler{responseBuilder: &LoginResponseBuilder{}}
	return &loginRequestHandler
}

// HandleLoginRequest retrieves data for response
func (lrh *LoginRequestHandler) HandleLoginRequest(w http.ResponseWriter, r *http.Request) {
	var request models.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		return
	}

	response, err := lrh.responseBuilder.build(request)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}
