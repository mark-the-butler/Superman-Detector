package builders

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

func (tr *testRepo) SaveLogin(request models.LoginRequest) bool {
	*tr.saveCalled = true
	return tr.saved
}

func (tr *testRepo) GetLocation(ipAddress string) (models.GeoLocation, error) {
	switch ipAddress {
	case "1.2.3.4/5":
		return models.GeoLocation{Lat: 16.4792, Lon: 104.6583, Radius: 500}, nil
	case "5.4.3.2/1":
		return models.GeoLocation{Lat: 34.7725, Lon: 113.7266, Radius: 50}, nil
	default:
		return tr.response, tr.err
	}
}

func (tr *testRepo) GetPreviousAndFutureIPAdress(username string, currentIP string, currentTimeStamp int) (previousLogin, futureLogin models.LoginRequest) {
	previousLocation := models.LoginRequest{IPAddress: "1.2.3.4/5", UnixTimestamp: 1611689190, EventUUID: "775CB1A3-FF0F-14E6-15D9-3FDB98B5C530", Username: "Bob"}
	futureLocation := models.LoginRequest{IPAddress: "5.4.3.2/1", UnixTimestamp: 1611696390, EventUUID: "5CC1A437-F1AE-409A-22CB-7E3769858939", Username: "Bob"}
	return previousLocation, futureLocation
}

var responseBuilder = LoginResponseBuilder{geoRepository: &testRepo{}}

func TestResponseBuilderBuild(t *testing.T) {
	loginRequest := models.LoginRequest{
		Username:      "Bob",
		UnixTimestamp: 1611692790,
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

		actual, _ := responseBuilder.Build(loginRequest)

		if reflect.DeepEqual(actual.CurrentLocation, expectedLocation) != true {
			t.Errorf("builder returned incorrect current location: got %+v\n want %+v\n", actual.CurrentLocation, expectedLocation)
		}
	})

	t.Run("Returns error when repository gives back error", func(t *testing.T) {
		expectedError := "could not retrieve location"
		expectedResponse := models.TravelResponse{}

		responseBuilder.geoRepository = &testRepo{response: models.GeoLocation{}, err: errors.New(expectedError), saveCalled: &defaultSave, saved: true}

		actual, err := responseBuilder.Build(loginRequest)

		if actual != expectedResponse {
			t.Errorf("builder returned a response but wanted it to be nil")
		}

		if err.Error() != expectedError {
			t.Errorf("builder returned incorrect error: got %v want %v", err.Error(), expectedError)
		}
	})

	t.Run("Saves login attempt to the database", func(t *testing.T) {
		saveCalled := false

		responseBuilder.geoRepository = &testRepo{saveCalled: &saveCalled, saved: true}

		_, err := responseBuilder.Build(loginRequest)

		if saveCalled != true {
			t.Errorf("builder did not make any save calls")
		}

		if err != nil {
			t.Errorf("an error was returned but expected there not to be: %v", err.Error())
		}
	})

	t.Run("Returns error when login cannot be saved to database", func(t *testing.T) {
		expectedError := "login could not be saved"
		expectedResponse := models.TravelResponse{}

		responseBuilder.geoRepository = &testRepo{saved: false, saveCalled: &defaultSave}

		actual, err := responseBuilder.Build(loginRequest)

		if err.Error() != expectedError {
			t.Errorf("builder returned incorrect error: got %v want %v", err.Error(), expectedError)
		}

		if actual != expectedResponse {
			t.Errorf("builder returned a response but wanted it to be nil")
		}
	})

	t.Run("Returns entire response", func(t *testing.T) {
		expectedResponse := models.TravelResponse{
			CurrentLocation: models.GeoLocation{
				Lat:    34.2,
				Lon:    -83.2,
				Radius: 20,
			},
			TravelToCurrentGeoSuspicious:   true,
			TravelFromCurrentGeoSuspicious: true,
			PreviousLocation: &models.GeoLocation{
				IP:        "1.2.3.4/5",
				Speed:     8895,
				Lat:       16.4792,
				Lon:       104.6583,
				Radius:    500,
				TimeStamp: 1611689190,
			},
			FutureLocation: &models.GeoLocation{
				IP:        "5.4.3.2/1",
				Speed:     7545,
				Lat:       34.7725,
				Lon:       113.7266,
				Radius:    50,
				TimeStamp: 1611696390,
			},
		}

		responseBuilder.geoRepository = &testRepo{response: defaultLocation, saved: true, saveCalled: &defaultSave}

		actual, _ := responseBuilder.Build(loginRequest)

		if reflect.DeepEqual(actual.CurrentLocation, expectedResponse.CurrentLocation) != true {
			t.Errorf("builder returned incorrect current location: got %+v\n want %+v\n", actual.CurrentLocation, expectedResponse.CurrentLocation)
		}

		if reflect.DeepEqual(actual.PreviousLocation, expectedResponse.PreviousLocation) != true {
			t.Errorf("builder returned incorrect previous location: got %+v\n want %+v\n", actual.PreviousLocation, expectedResponse.PreviousLocation)
		}

		if reflect.DeepEqual(actual.FutureLocation, expectedResponse.FutureLocation) != true {
			t.Errorf("builder returned incorrect future location: got %+v\n want %+v\n", actual.FutureLocation, expectedResponse.FutureLocation)
		}

		if actual.TravelFromCurrentGeoSuspicious != expectedResponse.TravelFromCurrentGeoSuspicious {
			t.Errorf("builder returned incorrect travel from current geo suspicous: got %v want %v", actual.TravelFromCurrentGeoSuspicious, false)
		}

		if actual.TravelToCurrentGeoSuspicious != expectedResponse.TravelToCurrentGeoSuspicious {
			t.Errorf("builder returned incorrect travel to current geo suspicous: got %v want %v", actual.TravelToCurrentGeoSuspicious, false)
		}
	})
}

func TestResponseBuilderCalculateGeoSuspiciousSpeed(t *testing.T) {
	fromLogin := models.LoginRequest{
		Username:      "Bob",
		UnixTimestamp: 1611692790,
		EventUUID:     "763 - 283",
		IPAddress:     "128.92.3.1/24",
	}

	toLogin := models.LoginRequest{
		IPAddress:     "5.4.3.2/1",
		UnixTimestamp: 1611696390,
		EventUUID:     "5CC1A437-F1AE-409A-22CB-7E3769858939",
		Username:      "Bob",
	}

	fromLocation := models.GeoLocation{
		Lat:    34.2,
		Lon:    -83.2,
		Radius: 20,
	}

	toLocation := models.GeoLocation{
		Lat:    34.7725,
		Lon:    113.7266,
		Radius: 50,
	}

	speed := calculateGeoSuspiciousSpeed(fromLogin, toLogin, fromLocation, toLocation)

	if speed != 7545 {
		t.Errorf("incorrect speed: got %v want %v", speed, 7545)
	}
}
