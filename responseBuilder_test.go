package main

import (
	"errors"
	"reflect"
	"testing"

	"github.com/mysteryboy73/Superman-Detector/models"
)

type testRepo struct {
	response   models.GeoLocation
	err        error
	saved      bool
	saveCalled *bool
}

func (tr *testRepo) saveLogin(request models.LoginRequest) bool {
	*tr.saveCalled = true
	return tr.saved
}

func (tr *testRepo) getLocation(request models.LoginRequest) (models.GeoLocation, error) {
	return tr.response, tr.err
}

var responseBuilder = LoginResponseBuilder{geoRepository: &testRepo{}}

func TestResponseBuilder(t *testing.T) {
	loginRequest := models.LoginRequest{
		Username:      "Bob",
		UnixTimestamp: 3847561023,
		EventUUID:     "763 - 283",
		IPAddress:     "128.92.3.1/24",
	}

	defaultLocation := models.GeoLocation{
		Lat:    34.2,
		Lon:    -83.2,
		Radius: 20,
	}

	var defaultSave bool

	responseBuilder.geoRepository = &testRepo{response: defaultLocation, saved: true, saveCalled: &defaultSave}

	t.Run("Returns TravelResponse with current location info", func(t *testing.T) {
		expectedLocation := models.GeoLocation{
			Lat:    34.2,
			Lon:    -83.2,
			Radius: 20,
		}

		responseBuilder.geoRepository = &testRepo{response: expectedLocation, saved: true, saveCalled: &defaultSave}

		actual, _ := responseBuilder.build(loginRequest)

		if reflect.DeepEqual(actual.CurrentLocation, expectedLocation) != true {
			t.Errorf("builder returned incorrect current location: got %+v\n want %+v\n", actual.CurrentLocation, expectedLocation)
		}

		if actual.TravelFromCurrentGeoSuspicious != false {
			t.Errorf("builder returned incorrect travel from current geo suspicous: got %v want %v", actual.TravelFromCurrentGeoSuspicious, false)
		}

		if actual.TravelToCurrentGeoSuspicious != false {
			t.Errorf("builder returned incorrect travel to current geo suspicous: got %v want %v", actual.TravelToCurrentGeoSuspicious, false)
		}
	})

	t.Run("Returns error when repository gives back error", func(t *testing.T) {
		expectedError := "no records found"
		expectedResponse := models.TravelResponse{}

		responseBuilder.geoRepository = &testRepo{response: models.GeoLocation{}, err: errors.New(expectedError), saveCalled: &defaultSave, saved: true}

		actual, err := responseBuilder.build(loginRequest)

		if actual != expectedResponse {
			t.Errorf("builder returned a response but wanted it to be nil")
		}

		if err.Error() != expectedError {
			t.Errorf("builder returned incorrect error: got %v want %v", err.Error(), expectedError)
		}
	})

	t.Run("Saves login attempt to the database", func(t *testing.T) {
		saveCalled := false

		responseBuilder.geoRepository = &testRepo{saveCalled: &saveCalled}

		actual, _ := responseBuilder.build(loginRequest)

		if saveCalled != true {
			t.Errorf("builder did not make any save calls")
		}

		if reflect.DeepEqual(actual.CurrentLocation, defaultLocation) {
			t.Errorf("current location is different that expected")
		}
	})

	t.Run("Returns error when login cannot be saved to database", func(t *testing.T) {
		expectedError := "login could not be saved"
		expectedResponse := models.TravelResponse{}

		responseBuilder.geoRepository = &testRepo{saved: false, saveCalled: &defaultSave}

		actual, err := responseBuilder.build(loginRequest)

		if err.Error() != expectedError {
			t.Errorf("builder returned incorrect error: got %v want %v", err.Error(), expectedError)
		}

		if actual != expectedResponse {
			t.Errorf("builder returned a response but wanted it to be nil")
		}
	})
}
