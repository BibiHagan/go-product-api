package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/mux"
)

// Option contains details
type Option struct {
	ID          string `json:"Id"`
	ProductID   string `json:"ProductId"`
	Name        string `json:"Name"`
	Description string `json:"Description"`
}

func returnOptionsByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnOptionsByID")
	vars := mux.Vars(r)
	key := vars["id"]

	var optionsList []Option

	for _, option := range Options {
		if option.ProductID == key {
			optionsList = append(optionsList, option)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(optionsList)
}

func returnOptionForProduct(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnOptionForProduct")
	vars := mux.Vars(r)
	pkey := vars["id"]
	okey := vars["optionId"]

	for _, option := range Options {
		if option.ProductID == pkey {
			if option.ID == okey {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(option)
			}
		}
	}

}

func createNewOption(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: createNewOption")

	vars := mux.Vars(r)
	key := vars["id"]
	var prodExists = false
	var optionsList []Option

	for _, product := range Products {
		if product.ID == key {
			prodExists = true
		}
	}

	if prodExists {
		reqBody, _ := ioutil.ReadAll(r.Body)
		var option Option
		json.Unmarshal(reqBody, &option)
		Options = append(Options, option)

		for _, option := range Options {
			if option.ProductID == key {
				optionsList = append(optionsList, option)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(optionsList)
	} else {
		json.NewEncoder(w).Encode("Product does not exist")

	}
}

// Options is a list of Option
var Options []Option
