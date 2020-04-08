package main

import "github.com/mysteryboy73/Superman-Detector/models"

// ResponseBuilder interface for building response
type ResponseBuilder interface {
	build(loginRequest models.LoginRequest) models.TravelResponse
}

// LoginResponseBuilder retrieves necessary data for response
type LoginResponseBuilder struct {
}

func (lrb *LoginResponseBuilder) build(request models.LoginRequest) models.TravelResponse {
	return models.TravelResponse{CurrentLocation: GetCurrentLocation(request)}
}
