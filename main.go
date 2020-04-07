package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/loginRequest", HandleLoginRequest)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
