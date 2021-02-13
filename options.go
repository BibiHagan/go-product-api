package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/mux"
)

// Options is a list of Option
var Options []Option

// Option contains details
type Option struct {
	ID          string `json:"Id"`
	ProductID   string `json:"ProductId"`
	Name        string `json:"Name"`
	Description string `json:"Description"`
}

// GET /products/{id}/options - finds all options for a specified product.
func returnOptionsByProductID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnOptionsByProductID")

	// get {id} from URL
	vars := mux.Vars(r)
	key := vars["id"]

	var optionsList []Option

	// find all the options for a product and add to list
	for _, option := range Options {
		if option.ProductID == key {
			optionsList = append(optionsList, option)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(optionsList)
}

// GET /products/{id}/options/{optionId} - finds the specified product option for the specified product.
func returnOptionForProduct(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnOptionForProduct")

	// get {id} and {optionId} from URL
	vars := mux.Vars(r)
	pkey := vars["id"]
	okey := vars["optionId"]
	optionExists := false

	// search through the Options[] the option with the right combination
	for _, option := range Options {
		if option.ProductID == pkey {
			if option.ID == okey {
				// return the option if found
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(option)
				optionExists = true
			}
		}
	}

	if !optionExists {
		// return 404 error product not found
		returnError(w, http.StatusNotFound, "Option Not Found")
	}
}

// POST /products/{id}/options - adds a new product option to the specified product.
func createNewOption(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: createNewOption")

	// get {id} from URL
	vars := mux.Vars(r)
	key := vars["id"]
	var prodExists = false
	var optionsList []Option

	// search Product[] to ensure Product exists
	for _, product := range Products {
		if product.ID == key {
			prodExists = true
		}
	}

	// If product exists create new Option
	if prodExists {
		reqBody, _ := ioutil.ReadAll(r.Body)
		var option Option
		json.Unmarshal(reqBody, &option)

		// Todo: if invalid or exists, return 400 bad request

		Options = append(Options, option)

		// create new Option
		for _, option := range Options {
			if option.ProductID == key {
				optionsList = append(optionsList, option)
			}
		}
	} else {
		// return 404 error product not found
		returnError(w, http.StatusNotFound, "Product does not exist Option not created")
	}
}

// PUT /products/{id}/options/{optionId} - updates the specified product option.
func updateOption(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: updateOption")

	// get {id} and {optionId} from URL
	vars := mux.Vars(r)
	pkey := vars["id"]
	okey := vars["optionId"]
	optionExists := false

	// search Options[] for the option with the matching combination
	for index, option := range Options {
		if option.ProductID == pkey {
			if option.ID == okey {
				// get new option from request
				reqBody, _ := ioutil.ReadAll(r.Body)
				var option Option
				json.Unmarshal(reqBody, &option)

				// Todo: if invalid or exists, return 400 bad request

				Options[index] = option
				optionExists = true
			}
		}
	}

	if !optionExists {
		// return 404 error product not found
		returnError(w, http.StatusNotFound, "Update error: Option Not Found")
	}
}

// DELETE /products/{id}/options/{optionId} - deletes the specified product option.
func deleteOption(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: deleteOption")

	// get {id} and {optionId} from URL
	vars := mux.Vars(r)
	pkey := vars["id"]
	okey := vars["optionId"]
	optionExists := false

	index := 0
	for _, option := range Options {
		if option.ProductID == pkey {
			if option.ID == okey {
				Options = append(Options[:index], Options[index+1:]...)
				index--
				optionExists = true
			}
		}
		index++
	}

	if !optionExists {
		// return 404 error product not found
		returnError(w, http.StatusNotFound, "Delete error: Option Not Found")
	}
}
