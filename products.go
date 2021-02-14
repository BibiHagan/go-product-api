package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
)

// Products is a list of product
var Products []Product

// Product contains details
type Product struct {
	ID            string  `json:"Id"`
	Name          string  `json:"Name"`
	Description   string  `json:"Description"`
	Price         float32 `json:"Price"`
	DeliveryPrice float32 `json:"DeliveryPrice"`
}

// GET /products - gets all products.
func returnAllProducts(w http.ResponseWriter, r *http.Request) {
	// Check URL for params
	u, err := url.Parse(r.RequestURI)
	if err != nil {
		returnError(w, http.StatusBadRequest, err)
	}
	params, _ := url.ParseQuery(u.RawQuery)

	// if there are params
	if len(params) > 0 {
		// if param is a {name}
		name := params["name"][0]
		if name != "" {
			returnProductByName(w, name)
		} else {
			// return error message and 400 bad request
			//returnError(w, http.StatusBadRequest, "Unknown Params")
		}
	} else {
		// else return ALL products in json format
		fmt.Println("Endpoint Hit: returnAllProducts")
		// Create read-only transaction
		txn := ProdDB.Txn(false)
		it, err := txn.Get("product", "id")
		if err != nil {
			returnError(w, http.StatusBadRequest, err)
		}

		var products []Product
		// iterate through the product DB and add ALL objects to Products[]
		for obj := it.Next(); obj != nil; obj = it.Next() {
			p := obj.(*Product)
			products = append(products, *p)
		}

		// encode as json and return
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(products)
	}
}

// GET /products/{id} - gets the product that matches the specified ID - ID is a GUID.
func returnProductByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnProductByID")

	// get {id} from URL
	vars := mux.Vars(r)
	key := vars["id"]

	// Create read-only transaction
	txn := ProdDB.Txn(false)
	defer txn.Abort()

	// search for {id}
	product, err := txn.First("product", "id", key)
	if err != nil {
		// return 404 error product not found
		returnError(w, http.StatusNotFound, err)
	}

	// encode as json and return
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

// GET /products?name={name} - finds all products matching the specified name.
func returnProductByName(w http.ResponseWriter, name string) {
	fmt.Println("Endpoint Hit: returnProductByName")

	// Create read-only transaction
	txn := ProdDB.Txn(false)
	defer txn.Abort()

	// get all Products with {name}
	prodsList, err := txn.Get("product", "name", name)
	if err != nil {
		// no products found return error message and 404 not found
		returnError(w, http.StatusNotFound, err)
	}

	var products []Product
	// iterate through the prodList created by query
	for obj := prodsList.Next(); obj != nil; obj = prodsList.Next() {
		p := obj.(*Product)
		products = append(Products, *p)
	}

	// encode as json and return
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// POST /products - creates a new product.
func createNewProduct(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: createNewProduct")

	// get the new product from the request
	reqBody, _ := ioutil.ReadAll(r.Body)
	var prod Product
	json.Unmarshal(reqBody, &prod)

	product := []*Product{
		{prod.ID, prod.Name, prod.Description, prod.Price, prod.DeliveryPrice},
	}

	// Create a write transaction
	txn := ProdDB.Txn(true)

	// insert new product in the database
	for _, p := range product {
		if err := txn.Insert("product", p); err != nil {
			returnError(w, http.StatusNotFound, err)
		}
	}

	// Commit the transaction
	txn.Commit()
}

// PUT /products/{id} - updates a product.
func updateProductByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: updateProduct")

	// get product from request
	reqBody, _ := ioutil.ReadAll(r.Body)

	var prod Product
	// Unmarshal to create product
	json.Unmarshal(reqBody, &prod)

	product := []*Product{
		{prod.ID, prod.Name, prod.Description, prod.Price, prod.DeliveryPrice},
	}

	// Create a write transaction
	txn := ProdDB.Txn(true)

	// update product in the database
	for _, p := range product {
		if err := txn.Insert("product", p); err != nil {
			returnError(w, http.StatusNotFound, err)
		}
	}

	// Commit the transaction
	txn.Commit()
}

// DELETE /products/{id} - deletes a product and its options.
func deleteProductByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: deleteProductByID")

	// find Product
	// get {id} from URL
	vars := mux.Vars(r)
	key := vars["id"]

	// Create read-only transaction
	txn := ProdDB.Txn(false)
	defer txn.Abort()

	// search for {id}
	product, err := txn.First("product", "id", key)
	if err != nil {
		// return 404 error product not found
		returnError(w, http.StatusNotFound, err)
	}

	// delete product
	// Create a write transaction
	txn = ProdDB.Txn(true)

	// delete product in the database
	err = txn.Delete("product", product)
	if err != nil {
		// return 404 error product not found
		returnError(w, http.StatusNotFound, err)
	}

	// Commit the transaction
	txn.Commit()
}
