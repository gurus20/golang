package main

import (
	"net/http"
	// "controllers/controller"
)

func Routes() {
	// http.HandleFunc("/", HomeHandler)
	http.ListenAndServe(":8000", nil)
}