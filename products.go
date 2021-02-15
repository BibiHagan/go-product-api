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
	ProductID     string  `json:"ProductId"`
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
			if products == nil {
				// return error message and 404
				returnError(w, http.StatusNotFound, "Product Not found")
			} else {
				Encode(w, products)
			}
		} else {
			// return error message and 400 bad request
			returnError(w, http.StatusBadRequest, "Unknown Params")
		}
	} else {
		// else return ALL products in json format
		fmt.Println("Endpoint Hit: returnAllProducts")

		products = getAllProducts(w, "id", "")
		if products != nil {
			Encode(w, products)
		}
	}
}

// GET /products/{productId} - gets the product that matches the specified ID - ID is a GUID.
func returnProductByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnProductByID")

	// get {productId} from URL
	pkey := mux.Vars(r)["productId"]

	product := getSingleProduct(w, "id", pkey)

	if product != nil {
		Encode(w, product)
	} else {
		Encode(w, []Product{})
	}
}

// POST /products - creates a new product.
func createNewProduct(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: createNewProduct")

	writeProductToDB(w, r)
}

// PUT /products/{productId} - updates a product.
func updateProductByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: updateProduct")

	// get {productId} from URL
	pkey := mux.Vars(r)["productId"]

	// check Product exists
	prodExist := getSingleProduct(w, "id", pkey)
	if prodExist == nil {
		returnError(w, http.StatusNotFound, "Update Fail: Product not found")
	} else {
		writeProductToDB(w, r)
	}
}

// DELETE /products/{productId} - deletes a product and its options.
func deleteProductByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: deleteProductByID")

	// get {productId} from URL
	pkey := mux.Vars(r)["productId"]

	// find the product to delete
	product := getSingleProduct(w, "id", pkey)

	// delete product
	deleteProduct(w, product, pkey)
}

func getSingleProduct(w http.ResponseWriter, index, pkey string) interface{} {
	// Create read-only transaction
	txn := Database.Txn(false)
	defer txn.Abort()

	// search for {productId}
	product, err := txn.First("product", index, pkey)
	if err != nil {
		// return DB error
		returnError(w, http.StatusInternalServerError, err.Error())
	}

	return product
}

func getAllProducts(w http.ResponseWriter, index, pkey string) []Product {
	// Create read-only transaction
	txn := Database.Txn(false)
	defer txn.Abort()

	var it memdb.ResultIterator
	var err error
	if pkey == "" {
		it, err = txn.Get("product", index)
	} else {
		it, err = txn.Get("product", index, pkey)
	}

	var products []Product
	if err == nil {
		// iterate through the product DB and add ALL objects
		for obj := it.Next(); obj != nil; obj = it.Next() {
			p := obj.(*Product)
			products = append(products, *p)
		}
	} else {
		// return DB error
		returnError(w, http.StatusInternalServerError, err.Error())
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
		{prod.ProductID, prod.Name, prod.Description, prod.Price, prod.DeliveryPrice},
	}

	// Create a write transaction
	txn := Database.Txn(true)

	// insert new product in the database
	for _, p := range product {
		if err := txn.Insert("product", p); err != nil {
			// return DB error
			returnError(w, http.StatusInternalServerError, err.Error())
		}
	}

	// Commit the transaction
	txn.Commit()
}

func deleteProduct(w http.ResponseWriter, product interface{}, pkey string) {
	if product != nil {
		// Create a write transaction
		txn := Database.Txn(true)

		// Get all options for productId
		options := getAllOptions(w, "productId", pkey)

		// Delete all options[]
		for _, opt := range options {
			deleteOptionFromDB(w, txn, opt)
		}

		// delete product in the database
		err := txn.Delete("product", product)
		if err != nil {
			// return DB error
			returnError(w, http.StatusInternalServerError, err.Error())
		}

		// Commit the transaction
		txn.Commit()

	} else {
		returnError(w, http.StatusNotFound, "Delete Fail: Product not found")
	}
}

// Encode sets the header and encode to JSON
func Encode(w http.ResponseWriter, product interface{}) {
	// encode as json and return
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}
