package main

import (
	"encoding/json"
	"net/http"
)

// Encode sets the header and encode to JSON
func Encode(w http.ResponseWriter, product interface{}) {
	// encode as json and return
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}
