package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hashicorp/go-memdb"
)

// Option contains details
type Option struct {
	OptionID    string `json:"OptionId"`
	ProductID   string `json:"ProductId"`
	Name        string `json:"Name"`
	Description string `json:"Description"`
}

// GET /products/{productId}/options - finds all options for a specified product.
func returnOptionsByProductID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnOptionsForProductID")

	// get {productId} from URL
	pkey := mux.Vars(r)["productId"]

	// get all the options for product {productId}
	options := getAllOptions(w, "productId", pkey)

	if options != nil {
		Encode(w, options)
	} else {
		Encode(w, []Option{})
	}
}

// GET /products/{productId}/options/{optionId} - finds the specified product option for the specified product.
func returnOptionForProduct(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnOptionForProduct")

	// get {productId} from URL
	vars := mux.Vars(r)
	pkey := vars["productId"]
	okey := vars["optionId"]

	option := getSingleOption(w, pkey, okey)

	if option != nil {
		Encode(w, option)
	} else {
		Encode(w, []Option{})
	}
}

// POST /products/{productId}/options - adds a new product option to the specified product.
func createNewOption(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: createNewOption")

	writeOptionToDB(w, r)
}

// PUT /products/{productId}/options/{optionId} - updates the specified product option.
func updateOption(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: updateOption")

	// get {productId} from URL
	vars := mux.Vars(r)
	pkey := vars["productId"]
	okey := vars["optionId"]

	// check Option exists
	option := getSingleOption(w, pkey, okey)

	if option == nil {
		returnError(w, http.StatusNotFound, "Update Fail: Product Option Not Found")
	} else {
		writeOptionToDB(w, r)
	}
}

// DELETE /products/{productId}/options/{optionId} - deletes the specified product option.
func deleteOption(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: deleteOption")

	// get {productId} from URL
	vars := mux.Vars(r)
	pkey := vars["productId"]
	okey := vars["optionId"]

	// check Option exists
	option := getSingleOption(w, pkey, okey)

	if option != nil {
		// Create a write transaction
		txn := Database.Txn(true)

		// delete Option
		if !deleteOptionFromDB(w, txn, option) {
			// Commit the transaction
			txn.Commit()
		}

	} else {
		returnError(w, http.StatusNotFound, "Delete Fail: Option not found")
	}

}

func getSingleOption(w http.ResponseWriter, pkey, okey string) interface{} {
	// Create read-only transaction
	txn := Database.Txn(false)
	defer txn.Abort()

	// search for the options for product {productId}
	it, err := txn.Get("option", "productId", pkey)
	if err != nil {
		// return DB error
		returnError(w, http.StatusInternalServerError, err.Error())
	}

	// iterate through the Options returned and retrn the option with {optionId}
	for obj := it.Next(); obj != nil; obj = it.Next() {
		o := obj.(*Option)
		if o.OptionID == okey {
			return *o
		}
	}

	return nil
}

func getAllOptions(w http.ResponseWriter, index, key string) []Option {
	// Create read-only transaction
	txn := Database.Txn(false)
	defer txn.Abort()

	var it memdb.ResultIterator
	var err error
	it, err = txn.Get("option", index, key)
	if err != nil {
		returnError(w, http.StatusBadRequest, err.Error())
	}

	// iterate through the option DB and add ALL options
	var options []Option
	for obj := it.Next(); obj != nil; obj = it.Next() {
		o := obj.(*Option)
		options = append(options, *o)
	}

	return options
}

func writeOptionToDB(w http.ResponseWriter, r *http.Request) {
	// get the new option from the request
	reqBody, _ := ioutil.ReadAll(r.Body)
	var opt Option
	json.Unmarshal(reqBody, &opt)

	option := []*Option{
		{opt.OptionID, opt.ProductID, opt.Name, opt.Description},
	}

	// Create a write transaction
	txn := Database.Txn(true)

	// insert new product in the database
	for _, p := range option {
		if err := txn.Insert("option", p); err != nil {
			returnError(w, http.StatusInternalServerError, err.Error())
		}
	}

	// Commit the transaction
	txn.Commit()
}

func deleteOptionFromDB(w http.ResponseWriter, txn *memdb.Txn, option interface{}) bool {
	errs := false
	// delete option in the database
	err := txn.Delete("option", option)
	if err != nil {
		// return 404 error product not found
		returnError(w, http.StatusNotFound, err.Error())
		errs = true
	}

	return errs
}
