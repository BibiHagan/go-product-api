package main

import (
	"encoding/json"
	"net/http"
)

// Error details
type Error struct {
	Code        int
	Description error
}

func returnError(w http.ResponseWriter, code int, err error) {
	var error Error
	error.Code = code
	error.Description = err

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(error)
}
