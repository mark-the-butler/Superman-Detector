package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mysteryboy73/Superman-Detector/models"
)

type testbuilder struct {
	testResponse models.TravelResponse
}

func (tb *testbuilder) build(request models.LoginRequest) models.TravelResponse {
	return tb.testResponse
}

var loginRequestHandler = LoginRequestHandler{responseBuilder: &testbuilder{}}

func TestHandleLoginRequest(t *testing.T) {
	loginRequestHandler.responseBuilder = &testbuilder{testResponse: models.TravelResponse{
		CurrentLocation: models.CurrentGeo{
			Lat:    1.23,
			Lon:    -4.56,
			Radius: 20,
		},
	}}
	t.Run("Sends back 200 okay with json body", func(t *testing.T) {
		login := models.LoginRequest{
			Username:      "Bob",
			UnixTimestamp: 1586223780,
			EventUUID:     "85ad929a-db03-4bf4-9541-8f728fa12e42",
			IPAddress:     "1.32.196.0/24",
		}
		loginAsJSON, _ := json.Marshal(login)
		request := httptest.NewRequest("POST", "/loginRequest", bytes.NewBuffer(loginAsJSON))
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

}
