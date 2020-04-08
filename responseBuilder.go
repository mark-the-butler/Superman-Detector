package main

import "github.com/mysteryboy73/Superman-Detector/models"

// ResponseBuilder interface for building response
type ResponseBuilder interface {
	build(loginRequest models.LoginRequest) (response models.TravelResponse, err error)
}

// LoginResponseBuilder retrieves necessary data for response
type LoginResponseBuilder struct {
}

func (lrb *LoginResponseBuilder) build(request models.LoginRequest) (response models.TravelResponse, err error) {
	repsonse := models.TravelResponse{CurrentLocation: GetCurrentLocation(request)}
	return repsonse, nil
}
