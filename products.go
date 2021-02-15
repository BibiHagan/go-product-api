package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/hashicorp/go-memdb"
)

// Product contains details
type Product struct {
	ID            string  `json:"Id"`
	Name          string  `json:"Name"`
	Description   string  `json:"Description"`
	Price         float32 `json:"Price"`
	DeliveryPrice float32 `json:"DeliveryPrice"`
}

// GET /products - gets all products.
// GET /products?name={name} - finds all products matching the specified name.
func returnAllProducts(w http.ResponseWriter, r *http.Request) {
	// Check URL for params
	u, err := url.Parse(r.RequestURI)
	if err != nil {
		returnError(w, http.StatusBadRequest, err.Error())
	}
	params, _ := url.ParseQuery(u.RawQuery)

	var products []Product

	// if there are params
	if len(params) > 0 {
		// if param is a {name}
		name := params["name"][0]
		if name != "" {
			fmt.Println("Endpoint Hit: returnProductByName")
			products = getAllProducts(w, "name", name)
		} else {
			// return error message and 400 bad request
			returnError(w, http.StatusBadRequest, "Unknown Params")
		}
	} else {
		// else return ALL products in json format
		fmt.Println("Endpoint Hit: returnAllProducts")

		products = getAllProducts(w, "id", "")
	}

	// encode as json and return
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// GET /products/{id} - gets the product that matches the specified ID - ID is a GUID.
func returnProductByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnProductByID")

	// get {id} from URL
	key := mux.Vars(r)["id"]

	product := getSingleProduct(w, "id", key)

	// encode as json and return
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

// POST /products - creates a new product.
func createNewProduct(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: createNewProduct")

	writeProductToDB(w, r)
}

// PUT /products/{id} - updates a product.
func updateProductByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: updateProduct")

	// get {id} from URL
	vars := mux.Vars(r)
	key := vars["id"]

	// check Product exists
	prodExist := getSingleProduct(w, "id", key)
	if prodExist == nil {
		returnError(w, http.StatusNotFound, "Update Fail: Product not found")
	} else {
		writeProductToDB(w, r)
	}
}

// DELETE /products/{id} - deletes a product and its options.
func deleteProductByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: deleteProductByID")

	// get {id} from URL
	vars := mux.Vars(r)
	key := vars["id"]

	// find the product to delete
	product := getSingleProduct(w, "id", key)

	// delete product
	deleteProduct(w, product, key)
}

func getSingleProduct(w http.ResponseWriter, index, key string) interface{} {
	// Create read-only transaction
	txn := Database.Txn(false)
	defer txn.Abort()

	// search for {id}
	product, err := txn.First("product", index, key)
	if err != nil {
		// return 404 error product not found
		returnError(w, http.StatusNotFound, err.Error())
	}

	return product
}

func getAllProducts(w http.ResponseWriter, index, key string) []Product {
	// Create read-only transaction
	txn := Database.Txn(false)
	defer txn.Abort()

	var it memdb.ResultIterator
	var err error
	if key == "" {
		it, err = txn.Get("product", index)
	} else {
		it, err = txn.Get("product", index, key)
	}

	if err != nil {
		returnError(w, http.StatusInternalServerError, err.Error())
	}

	var products []Product
	// iterate through the product DB and add ALL objects
	for obj := it.Next(); obj != nil; obj = it.Next() {
		p := obj.(*Product)
		products = append(products, *p)
	}

	return products
}

func writeProductToDB(w http.ResponseWriter, r *http.Request) {
	// get product from request
	reqBody, _ := ioutil.ReadAll(r.Body)

	var prod Product
	// Unmarshal to create product
	json.Unmarshal(reqBody, &prod)
	product := []*Product{
		{prod.ID, prod.Name, prod.Description, prod.Price, prod.DeliveryPrice},
	}

	// Create a write transaction
	txn := Database.Txn(true)

	// insert new product in the database
	for _, p := range product {
		if err := txn.Insert("product", p); err != nil {
			returnError(w, http.StatusInternalServerError, err.Error())
		}
	}

	// Commit the transaction
	txn.Commit()
}

func deleteProduct(w http.ResponseWriter, product interface{}, key string) {
	if product != nil {
		// Create a write transaction
		txn := Database.Txn(true)

		// delete product in the database
		err := txn.Delete("product", product)
		if err != nil {
			// return 404 error product not found
			returnError(w, http.StatusNotFound, err.Error())
		}

		// Commit the transaction
		txn.Commit()

		// Get all options for productId
		options := getAllOptions(w, "productId", key)

		// Delete all options[]
		for _, opt := range options {
			deleteOptionFromDB(w, opt)
		}

	} else {
		returnError(w, http.StatusNotFound, "Delete Fail: Product not found")
	}
}
