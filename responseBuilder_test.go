package main

import (
	"errors"
	"testing"

	"github.com/mysteryboy73/Superman-Detector/models"
)

type testRepo struct {
	response models.GeoLocation
	err      error
}

func (tr *testRepo) getLocation(request models.LoginRequest) (models.GeoLocation, error) {
	return tr.response, tr.err
}

var responseBuilder = LoginResponseBuilder{geoRepository: &testRepo{}}

func TestResponseBuilder(t *testing.T) {
	var loginRequest = models.LoginRequest{
		Username:      "Bob",
		UnixTimestamp: 3847561023,
		EventUUID:     "763 - 283",
		IPAddress:     "128.92.3.1/24",
	}

	t.Run("Returns TravelResponse with current location info", func(t *testing.T) {
		expectedLat := 34.2
		expectedLon := -83.2
		expectedRadius := 20

		responseBuilder.geoRepository = &testRepo{response: models.GeoLocation{
			Lat:    expectedLat,
			Lon:    expectedLon,
			Radius: expectedRadius,
		}}

		actual, _ := responseBuilder.build(loginRequest)

		if actual.CurrentLocation.Lat != expectedLat {
			t.Errorf("builder returned incorrect Lat for CurrentLocation: got %v want %v", actual.CurrentLocation.Lat, expectedLat)
		}

		if actual.CurrentLocation.Lon != expectedLon {
			t.Errorf("builder returned incorrect Lon for CurrentLocation: got %v want %v", actual.CurrentLocation.Lon, expectedLon)
		}

		if actual.CurrentLocation.Radius != expectedRadius {
			t.Errorf("builder returned incorrect Radius for CurrentLocation: got %v want %v", actual.CurrentLocation.Radius, expectedRadius)
		}
	})

	t.Run("Returns error when repository gives back error", func(t *testing.T) {
		expectedError := "no records found"
		expectedResponse := models.TravelResponse{}

		responseBuilder.geoRepository = &testRepo{response: models.GeoLocation{}, err: errors.New(expectedError)}

		actual, err := responseBuilder.build(loginRequest)

		if actual != expectedResponse {
			t.Errorf("builder returned a response but wanted it to be nil")
		}

		if err.Error() != expectedError {
			t.Errorf("builder returned incorrect error: got %v want %v", err.Error(), expectedError)
		}

	})
}
