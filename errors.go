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

func returnError(w http.ResponseWriter, code int, err string) {
	var error Error
	error.Code = code
	error.Description = err

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(error)
}
