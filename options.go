package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

// Option contains details
type Option struct {
	ID          string `json:"Id"`
	ProductID   string `json:"ProductId"`
	Name        string `json:"Name"`
	Description string `json:"Description"`
}

// GET /products/{id}/options - finds all options for a specified product.
func returnOptionsByProductID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnOptionsForProductID")

	// get {id} from URL
	vars := mux.Vars(r)
	key := vars["id"]

	// Create read-only transaction
	txn := OptDB.Txn(false)
	defer txn.Abort()

	it, err := txn.Get("option", "productId", key)
	if err != nil {
		returnError(w, http.StatusBadRequest, err.Error())
	}

	var options []Option
	// iterate through the product DB and add ALL objects to Options[]
	for obj := it.Next(); obj != nil; obj = it.Next() {
		o := obj.(*Option)
		options = append(options, *o)
	}

	// encode as json and return
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(options)
}

// GET /products/{id}/options/{optionId} - finds the specified product option for the specified product.
func returnOptionForProduct(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnOptionForProduct")

	// get {id} from URL
	vars := mux.Vars(r)
	pkey := vars["id"]
	okey := vars["optionId"]

	// Create read-only transaction
	txn := OptDB.Txn(false)
	defer txn.Abort()

	// get all of the options for product {id}
	it, err := txn.Get("option", "productId", pkey)
	if err != nil {
		returnError(w, http.StatusBadRequest, err.Error())
	}

	var options []Option
	// iterate through the product DB and add options with {optionId}
	for obj := it.Next(); obj != nil; obj = it.Next() {
		o := obj.(*Option)
		if o.ID == okey {
			options = append(options, *o)
		}
	}

	// encode as json and return
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(options)
}

// POST /products/{id}/options - adds a new product option to the specified product.
func createNewOption(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: createNewOption")

	// get the new option from the request
	reqBody, _ := ioutil.ReadAll(r.Body)
	var opt Option
	json.Unmarshal(reqBody, &opt)

	option := []*Option{
		{opt.ID, opt.ProductID, opt.Name, opt.Description},
	}

	// Create a write Transaction
	txn := OptDB.Txn(true)

	// insert new Option in the database
	for _, o := range option {
		if err := txn.Insert("option", o); err != nil {
			returnError(w, http.StatusNotFound, err.Error())
		}
	}

	// Commit the transaction
	txn.Commit()
}

// PUT /products/{id}/options/{optionId} - updates the specified product option.
func updateOption(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: updateOption")

	// get {id} from URL
	vars := mux.Vars(r)
	pkey := vars["id"]
	okey := vars["optionId"]

	// Create read-only transaction
	txn := OptDB.Txn(false)
	defer txn.Abort()

	// get all options for product {id}
	it, err := txn.Get("option", "productId", pkey)
	if err != nil {
		returnError(w, http.StatusBadRequest, err.Error())
	}

	var options []Option
	// iterate through the product DB and add objects with {optionId}
	for obj := it.Next(); obj != nil; obj = it.Next() {
		o := obj.(*Option)
		if o.ID == okey {
			options = append(options, *o)
		}
	}

	// get new option from request
	reqBody, _ := ioutil.ReadAll(r.Body)

	var opt Option
	// Unmarshal to create option
	json.Unmarshal(reqBody, &opt)

	option := []*Option{
		{opt.ID, opt.ProductID, opt.Name, opt.Description},
	}

	// check that the optionID is the same before update
	if options[0].ID == option[0].ID {
		// Create a write Transaction
		txn = OptDB.Txn(true)

		// update option in the database
		for _, o := range option {
			if err := txn.Insert("option", o); err != nil {
				returnError(w, http.StatusNotFound, err.Error())
			}
		}

		// Commit the transaction
		txn.Commit()
	} else {
		returnError(w, http.StatusNotFound, "Update Fail: Product Option Not Found")
	}
}

// DELETE /products/{id}/options/{optionId} - deletes the specified product option.
func deleteOption(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: deleteOption")

	// get {id} from URL
	vars := mux.Vars(r)
	pkey := vars["id"]
	okey := vars["optionId"]

	// Create read-only transaction
	txn := OptDB.Txn(false)
	defer txn.Abort()

	// get all options for product {id}
	it, err := txn.Get("option", "productId", pkey)
	if err != nil {
		returnError(w, http.StatusBadRequest, err.Error())
	}

	var options []Option
	// iterate through the product DB and add objects with {optionId}
	for obj := it.Next(); obj != nil; obj = it.Next() {
		o := obj.(*Option)
		if o.ID == okey {
			options = append(options, *o)
		}
	}

	// delete option
	// Create a write transaction
	txn = OptDB.Txn(true)

	// delete product in the database
	err = txn.Delete("option", options[0])
	if err != nil {
		// return 404 error product not found
		returnError(w, http.StatusNotFound, err.Error())
	}

	// Commit the transaction
	txn.Commit()
}
