package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mysteryboy73/Superman-Detector/models"
)

func TestHandleLoginRequest(t *testing.T) {
	login := models.LoginRequest{
		Username:      "Bob",
		UnixTimestamp: 1586223780,
		EventUUID:     "85ad929a-db03-4bf4-9541-8f728fa12e42",
		IPAddress:     "1.32.196.0/24",
	}
	loginAsJSON, _ := json.Marshal(login)
	request := httptest.NewRequest("POST", "/loginRequest", bytes.NewBuffer(loginAsJSON))
	response := httptest.ResponseRecorder{}

	HandleLoginRequest(&response, request)

	expectedContentType := "application/json; charset=UTF-8"
	if contentType := response.Header().Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("handler returned wrong content-type: get %v want %v", contentType, expectedContentType)
	}

	if status := response.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

}
