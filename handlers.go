package main

import (
	"encoding/json"
	"io/ioutil"
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

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(body, &request); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		return
	}

	response := lrh.responseBuilder.build(request)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}
