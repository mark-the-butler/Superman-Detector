package main

import (
	"errors"

	"github.com/mysteryboy73/Superman-Detector/models"
)

// ResponseBuilder interface for building response
type ResponseBuilder interface {
	build(loginRequest models.LoginRequest) (response models.TravelResponse, err error)
}

// LoginResponseBuilder retrieves necessary data for response
type LoginResponseBuilder struct {
	geoRepository DataRepo
}

// NewLoginResponseBuilder returns a LoginResponseBuilder with repository dependency
func NewLoginResponseBuilder() *LoginResponseBuilder {
	loginResponseBuilder := LoginResponseBuilder{geoRepository: &GeoRepository{}}
	return &loginResponseBuilder
}

// Build constructs a TravelResponse
func (lrb *LoginResponseBuilder) build(request models.LoginRequest) (models.TravelResponse, error) {
	var response models.TravelResponse

	saved := lrb.geoRepository.saveLogin(request)

	if saved != true {
		return models.TravelResponse{}, errors.New("login could not be saved")
	}

	currentLocation, err := lrb.geoRepository.getLocation(request)

	if err != nil {
		return models.TravelResponse{}, err
	}

	response.CurrentLocation = currentLocation

	return response, nil
}
