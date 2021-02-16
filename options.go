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

	isNew := true
	writeOptionToDB(w, r, isNew)
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
		return
	}

	isNew := false
	writeOptionToDB(w, r, isNew)
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

	if option == nil {
		returnError(w, http.StatusNotFound, "Delete Fail: Option not found")
		return
	}

	// Create a write transaction
	txn := Database.Txn(true)

	// delete Option
	if deleteOptionFromDB(w, txn, option) {
		// Commit the transaction
		txn.Commit()
	}
}

// gets a single option with {productId} and {optionId}
func getSingleOption(w http.ResponseWriter, pkey, okey string) interface{} {
	// Create read-only transaction
	txn := Database.Txn(false)
	defer txn.Abort()

	// search for the options for product {productId}
	it, err := txn.Get("option", "productId", pkey)
	if err != nil {
		// return DB error
		returnError(w, http.StatusInternalServerError, err.Error())
		return nil
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

// gets all the option for {productId}
// if none found returns an empty []Option
func getAllOptions(w http.ResponseWriter, index, key string) []Option {
	// Create read-only transaction
	txn := Database.Txn(false)
	defer txn.Abort()

	var it memdb.ResultIterator
	var err error
	it, err = txn.Get("option", index, key)
	if err != nil {
		// return DB error
		returnError(w, http.StatusInternalServerError, err.Error())
		return nil
	}

	// iterate through the option DB and add ALL options
	var options []Option
	for obj := it.Next(); obj != nil; obj = it.Next() {
		o := obj.(*Option)
		options = append(options, *o)
	}

	return options
}

// Creates and updates a record in the Options Table
// if create record exists it returns an error
// can't create a new option for a product that does not exist
// can't update if product does not exist
func writeOptionToDB(w http.ResponseWriter, r *http.Request, isNew bool) {
	// get the new option from the request
	reqBody, _ := ioutil.ReadAll(r.Body)
	var opt Option
	json.Unmarshal(reqBody, &opt)

	newOption := []*Option{
		{opt.OptionID, opt.ProductID, opt.Name, opt.Description},
	}

	option := getSingleOption(w, opt.ProductID, opt.OptionID)

	if isNew {
		// do not create Option if Product does not exist
		product := getSingleProduct(w, "id", opt.ProductID)

		if product == nil {
			returnError(w, http.StatusNotFound, "New Option fail: product {productId} does not exist")
			return
		}
		// do not create new option if already exists
		if option != nil {
			returnError(w, http.StatusBadRequest, "New Option fail: option exists for this {productId}{optionId}")
			return
		}
	} else {
		// do not update option if it does not exist
		if option == nil {
			returnError(w, http.StatusNotFound, "Update Option fail: option not found for this {productId}{optionId}")
			return
		}
	}

	// Create a write transaction
	txn := Database.Txn(true)

	// insert new product in the database
	for _, p := range newOption {
		if err := txn.Insert("option", p); err != nil {
			returnError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	// Commit the transaction
	txn.Commit()
}

// Deletes option from Option Table and returns if delete was successful or not
func deleteOptionFromDB(w http.ResponseWriter, txn *memdb.Txn, option interface{}) bool {
	errs := false
	// delete option in the database
	err := txn.Delete("option", option)
	if err != nil {
		// return DB error
		returnError(w, http.StatusInternalServerError, err.Error())
		errs = true
	}

	// returns success
	return !errs
}
