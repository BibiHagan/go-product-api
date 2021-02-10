package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Product struct {
	Id            string  `json:"Id"`
	Name          string  `json:"Name"`
	Description   string  `json:"Description"`
	Price         float32 `json:"Price"`
	DeliveryPrice float32 `json:"DeliveryPrice"`
}

func returnAllProducts(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllProducts")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Products)
}

var Products []Product
