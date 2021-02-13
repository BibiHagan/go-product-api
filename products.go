package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/mux"
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
		panic(err)
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
			returnError(w, http.StatusBadRequest, "Unknown Params")
		}
	} else {
		// else return ALL products in json format
		fmt.Println("Endpoint Hit: returnAllProducts")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Products)
	}
}

// GET /products/{id} - gets the product that matches the specified ID - ID is a GUID.
func returnProductByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnProductByID")

	// get {id} from URL
	vars := mux.Vars(r)
	key := vars["id"]
	prodExists := false

	// search through Products[] for {id}
	for _, product := range Products {
		if product.ID == key {
			// If found return product in json format
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(product)
			prodExists = true
		}
	}

	if !prodExists {
		// return 404 error product not found
		returnError(w, http.StatusNotFound, "Product Not Found")
	}
}

// GET /products?name={name} - finds all products matching the specified name.
func returnProductByName(w http.ResponseWriter, name string) {
	fmt.Println("Endpoint Hit: returnProductByName")

	// create new list with product or requested name
	var prodsList []Product

	// go through All the products and add Product to new list
	for _, product := range Products {
		if product.Name == name {
			prodsList = append(prodsList, product)
		}
	}

	// If there are products in the return list
	if len(prodsList) > 0 {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(prodsList)
	} else {
		// no products found return error message and 404 not found
		returnError(w, http.StatusNotFound, "Product Not Found")
	}
}

// POST /products - creates a new product.
func createNewProduct(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: createNewProduct")

	// get the new product from the request
	reqBody, _ := ioutil.ReadAll(r.Body)
	var product Product
	// Unmarshal to create a Product and add to Products[]
	json.Unmarshal(reqBody, &product)

	// Todo: if invalid or exists, return 400 bad request

	Products = append(Products, product)
}

// PUT /products/{id} - updates a product.
func updateProductByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: updateProduct")

	// get product from request
	reqBody, _ := ioutil.ReadAll(r.Body)
	// TODO: if invalid, return 400 bad request

	var product Product
	// Unmarshal to create product
	json.Unmarshal(reqBody, &product)

	// get {id} from URL
	vars := mux.Vars(r)
	id := vars["id"]
	prodExists := false

	// Search through Products[] and overwrite item if found
	for index, prod := range Products {
		if prod.ID == id {
			Products[index] = product
			prodExists = true
		}
	}

	if !prodExists {
		// no products found return error message and 404 not found
		returnError(w, http.StatusNotFound, "Update Fail: Product Not Found")
	}
}

// DELETE /products/{id} - deletes a product and its options.
func deleteProductByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: deleteProductByID")

	// get {id} from URL
	vars := mux.Vars(r)
	id := vars["id"]
	prodExists := false

	// Search Products[] for {id}
	prodIndex := 0
	for _, product := range Products {
		// if found delete it
		if product.ID == id {
			Products = append(Products[:prodIndex], Products[prodIndex+1:]...)
			prodIndex--
			prodExists = true

			// search through Options[] and delete ALL options for that product
			optIndex := 0
			for _, option := range Options {
				if option.ProductID == id {
					Options = append(Options[:optIndex], Options[optIndex+1:]...)
					optIndex--
				}
				optIndex++
			}
		}
		prodIndex++
	}

	if !prodExists {
		// no products found return error message and 404 not found
		returnError(w, http.StatusNotFound, "Delete Fail: Product Not Found")
	}
}
