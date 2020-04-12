package handlers

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mysteryboy73/Superman-Detector/models"
)

type testbuilder struct {
	testResponse models.TravelResponse
	err          error
}

func (tb *testbuilder) Build(request models.LoginRequest) (models.TravelResponse, error) {
	return tb.testResponse, tb.err
}

var loginRequestHandler = LoginRequestHandler{responseBuilder: &testbuilder{}}
var validJSON = []byte(`{"username":"Bob","unix_timestamp":1514764800,"event_uuid":"85ad929a-db03-4bf4-9541-8f728fa12e42","ip_address":"1.32.196.0/24"}`)

func TestHandleLoginRequest(t *testing.T) {
	loginRequestHandler.responseBuilder = &testbuilder{testResponse: models.TravelResponse{
		CurrentLocation: models.GeoLocation{
			Lat:    1.23,
			Lon:    -4.56,
			Radius: 20,
		},
	}}
	t.Run("Sends back 200 okay with json body", func(t *testing.T) {
		request := httptest.NewRequest("POST", "/loginRequest", bytes.NewBuffer(validJSON))
		response := httptest.NewRecorder()

		loginRequestHandler.HandleLoginRequest(response, request)

		expectedContentType := "application/json; charset=UTF-8"
		if contentType := response.Header().Get("Content-Type"); contentType != expectedContentType {
			t.Errorf("handler returned wrong content-type: get %v want %v", contentType, expectedContentType)
		}

		if status := response.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		expectedBody := `{"currentGeo":{"lat":1.23,"lon":-4.56,"radius":20},"travelToCurrentGeoSuspicious":false,"traveFromCurrentGeoSuspicious":false}`
		if strings.TrimRight(response.Body.String(), "\n") != expectedBody {
			t.Errorf("handler returned unexpected body: got %v want %v", response.Body.String(), expectedBody)
		}
	})

	t.Run("Sends back full response body ", func(t *testing.T) {
		loginRequestHandler.responseBuilder = &testbuilder{testResponse: models.TravelResponse{
			CurrentLocation: models.GeoLocation{
				Lat:    1.23,
				Lon:    -4.56,
				Radius: 20,
			},
			TravelToCurrentGeoSuspicious:   true,
			TravelFromCurrentGeoSuspicious: false,
			PreviousLocation: &models.GeoLocation{
				IP:        "123.342.5.6/34",
				Speed:     783,
				Lat:       5,
				Lon:       1,
				Radius:    5,
				TimeStamp: 1586370808,
			},
			FutureLocation: &models.GeoLocation{
				IP:        "123.274.5.6/36",
				Speed:     200,
				Lat:       -4,
				Lon:       3,
				Radius:    10,
				TimeStamp: 1586370927,
			},
		}}
		request := httptest.NewRequest("POST", "/loginRequest", bytes.NewBuffer(validJSON))
		response := httptest.NewRecorder()

		loginRequestHandler.HandleLoginRequest(response, request)

		expectedBody := `{"currentGeo":{"lat":1.23,"lon":-4.56,"radius":20},"travelToCurrentGeoSuspicious":true,"traveFromCurrentGeoSuspicious":false,"precedingIpAccess":{"ip":"123.342.5.6/34","speed":783,"lat":5,"lon":1,"radius":5,"timestamp":1586370808},"subsequentIpAccess":{"ip":"123.274.5.6/36","speed":200,"lat":-4,"lon":3,"radius":10,"timestamp":1586370927}}`
		if strings.TrimRight(response.Body.String(), "\n") != expectedBody {
			t.Errorf("handler returned unexpected body: got %v want %v", response.Body.String(), expectedBody)
		}

	})

	t.Run("Sends back 400 when request is not formatted properly", func(t *testing.T) {
		invalidJSON := []byte(`{"username":"Bob",incorrectparameter:123}`)
		request := httptest.NewRequest("POST", "/loginRequest", bytes.NewBuffer(invalidJSON))
		response := httptest.NewRecorder()

		loginRequestHandler.HandleLoginRequest(response, request)

		if status := response.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusBadRequest)
		}
	})

	t.Run("Sends back 500 when responseBuilder returns error", func(t *testing.T) {
		loginRequestHandler.responseBuilder = &testbuilder{testResponse: models.TravelResponse{}, err: errors.New("Something went wrong")}
		request := httptest.NewRequest("POST", "/loginRequest", bytes.NewBuffer(validJSON))
		response := httptest.NewRecorder()

		loginRequestHandler.HandleLoginRequest(response, request)

		if status := response.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusInternalServerError)
		}
	})
}
