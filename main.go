package main

import (
	"log"
	"net/http"

	"github.com/mysteryboy73/Superman-Detector/handlers"
)

func main() {
	loginRequestHandler := handlers.NewLoginRequestHandler()

	http.HandleFunc("/loginRequest", loginRequestHandler.HandleLoginRequest)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Could not listen on port 8080 %v", err)
	}
}
