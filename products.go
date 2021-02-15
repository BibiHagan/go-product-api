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

		products = getAllProducts(w, "id", "all")
	}

	// encode as json and return
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// GET /products/{id} - gets the product that matches the specified ID - ID is a GUID.
func returnProductByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnProductByID")

	// get {id} from URL
	vars := mux.Vars(r)
	key := vars["id"]

	product := getSingleProduct(w, "id", key)

	// encode as json and return
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
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

	writeProductToDB(w, product)
}

// PUT /products/{id} - updates a product.
func updateProductByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: updateProduct")

	// get {id} from URL
	vars := mux.Vars(r)
	key := vars["id"]

	// check Product exists
	prodExist := getSingleProduct(w, "id", key)
	if prodExist != nil {
		// get product from request
		reqBody, _ := ioutil.ReadAll(r.Body)

		var prod Product
		// Unmarshal to create product
		json.Unmarshal(reqBody, &prod)
		product := []*Product{
			{prod.ID, prod.Name, prod.Description, prod.Price, prod.DeliveryPrice},
		}

		writeProductToDB(w, product)

	} else {
		returnError(w, http.StatusNotFound, "Update Fail: Product not found")
	}
}

// DELETE /products/{id} - deletes a product and its options.
func deleteProductByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: deleteProductByID")

	// find Product
	// get {id} from URL
	vars := mux.Vars(r)
	key := vars["id"]

	product := getSingleProduct(w, "id", key)

	// delete product
	deleteProduct(w, product)

	/* 	// delete all options for productId
	   	// Create read-only transaction
	   	txn := OptDB.Txn(false)
	   	defer txn.Abort()

	   	// get all of the options for product {id}
	   	it, err := txn.Get("option", "productId", key)
	   	if err != nil {
	   		returnError(w, http.StatusBadRequest, err.Error())
	   	}

	   	// delete Options
	   	// Create a write transaction
	   	txn = OptDB.Txn(true)

	   	// iterate through the product DB and delete options
	   	for obj := it.Next(); obj != nil; obj = it.Next() {
	   		o := obj.(*Option)
	   		if o.ProductID == key {
	   			err = txn.Delete("option", o)
	   			if err != nil {
	   				// return 404 error product not found
	   				returnError(w, http.StatusNotFound, err.Error())
	   			}
	   		}
	   	}

	   	// Commit the transaction
	   	txn.Commit() */
}

func getSingleProduct(w http.ResponseWriter, index, key string) interface{} {
	// Create read-only transaction
	txn := ProdDB.Txn(false)
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

	var products []Product
	var it memdb.ResultIterator
	var err error

	// Create read-only transaction
	txn := ProdDB.Txn(false)

	if key == "all" {
		it, err = txn.Get("product", index)
	} else {
		it, err = txn.Get("product", index, key)
	}

	if err != nil {
		returnError(w, http.StatusBadRequest, err.Error())
	}

	// iterate through the product DB and add ALL objects to Products[]
	for obj := it.Next(); obj != nil; obj = it.Next() {
		p := obj.(*Product)
		products = append(products, *p)
	}

	return products
}

func writeProductToDB(w http.ResponseWriter, product []*Product) {
	// Create a write transaction
	txn := ProdDB.Txn(true)

	// insert new product in the database
	for _, p := range product {
		if err := txn.Insert("product", p); err != nil {
			returnError(w, http.StatusNotFound, err.Error())
		}
	}

	// Commit the transaction
	txn.Commit()
}

func deleteProduct(w http.ResponseWriter, product interface{}) {
	if product != nil {
		// Create a write transaction
		txn := ProdDB.Txn(true)

		// delete product in the database
		err := txn.Delete("product", product)
		if err != nil {
			// return 404 error product not found
			returnError(w, http.StatusNotFound, err.Error())
		}

		// Commit the transaction
		txn.Commit()
	} else {
		returnError(w, http.StatusNotFound, "Delete Fail: Product not found")
	}
}
