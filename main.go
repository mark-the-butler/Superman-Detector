package main

import (
	"log"
	"net/http"
)

func main() {
	loginRequestHandler := NewLoginRequestHandler()

	http.HandleFunc("/loginRequest", loginRequestHandler.HandleLoginRequest)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Could not listen on port 8080 %v", err)
	}
}
