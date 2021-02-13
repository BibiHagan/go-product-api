package main

import (
	"encoding/json"
	"net/http"
)

// Error details
type Error struct {
	Code        int
	Description string
}

func returnError(w http.ResponseWriter, code int, description string) {
	var error Error
	error.Code = code
	error.Description = description

	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(error)
}
