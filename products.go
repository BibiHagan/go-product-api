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
		return
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
				return
			}

			Encode(w, products)
		} else {
			// return error message and 400 bad request
			returnError(w, http.StatusBadRequest, "Unknown Params")
			return
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

	isNew := true
	writeProductToDB(w, r, isNew)
}

// PUT /products/{productId} - updates a product.
func updateProductByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: updateProduct")

	isNew := false
	writeProductToDB(w, r, isNew)
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

// getSingleProduct returns a single record with {productId}
func getSingleProduct(w http.ResponseWriter, index, pkey string) interface{} {
	// Create read-only transaction
	txn := Database.Txn(false)
	defer txn.Abort()

	// search for {productId}
	product, err := txn.First("product", index, pkey)
	if err != nil {
		// return DB error
		returnError(w, http.StatusInternalServerError, err.Error())
		return nil
	}

	return product
}

// getAllProducts returns a []product of all products in the DB
// Pagination limit is 10
func getAllProducts(w http.ResponseWriter, index, pkey string, offset int) []Product {
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

	if err != nil {
		// return DB error
		returnError(w, http.StatusInternalServerError, err.Error())
		return []Product{}
	}

	var products []Product
	count := offset
	i := 0
	// iterate through the product DB and add ALL objects
	for obj := it.Next(); obj != nil; obj = it.Next() {
		if i > count {
			p := obj.(*Product)
			products = append(products, *p)
		}
		i++
		if i > count+10 {
			break
		}
	}

	return products
}

// Handles both create and update product
// new product returns error if it exists
// update product returns error if it does not exist
func writeProductToDB(w http.ResponseWriter, r *http.Request, isNew bool) {
	// get product from request
	reqBody, _ := ioutil.ReadAll(r.Body)

	var prod Product
	// Unmarshal to create product
	json.Unmarshal(reqBody, &prod)
	newProduct := []*Product{
		{prod.ProductID, prod.Name, prod.Description, prod.Price, prod.DeliveryPrice},
	}

	// search for {productId}
	product := getSingleProduct(w, "id", prod.ProductID)

	if isNew {
		// do not create new product if already exists
		if product != nil {
			returnError(w, http.StatusBadRequest, "New Product fail: product exists for this {productId}")
			return
		}
	} else {
		// do not update product if it does not exist
		if product == nil {
			returnError(w, http.StatusNotFound, "Update fail: product not found for this {productId}")
			return
		}
	}

	// Create a write transaction
	txn := Database.Txn(true)

	// insert new product in the database
	for _, p := range newProduct {
		if err := txn.Insert("product", p); err != nil {
			// return DB error
			returnError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	// Commit the transaction
	txn.Commit()
}

// Delete returns an error if no product found to delete
// gets all options for {productId} and deletes them
// This is done in one transaction so if product has an error
// transaction does not get commited
func deleteProduct(w http.ResponseWriter, product interface{}, pkey string) {
	if product == nil {
		returnError(w, http.StatusNotFound, "Delete Fail: Product not found")
		return
	}

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
		return
	}

	// Commit the transaction
	txn.Commit()
}

// Encode sets the header and encode to JSON
// used by both Products and Options
func Encode(w http.ResponseWriter, product interface{}) {
	// encode as json and return
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}
