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

func returnAllProducts(w http.ResponseWriter, r *http.Request) {

	u, err := url.Parse(r.RequestURI)
	if err != nil {
		panic(err)
	}
	params, _ := url.ParseQuery(u.RawQuery)

	if len(params) > 0 {
		returnProductByName(w, params["name"][0])
	} else {
		fmt.Println("Endpoint Hit: returnAllProducts")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Products)
	}
}

func returnProductByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnProductByID")
	vars := mux.Vars(r)
	key := vars["id"]

	for _, product := range Products {
		if product.ID == key {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(product)
		}
	}
}

func returnProductByName(w http.ResponseWriter, name string) {
	fmt.Println("Endpoint Hit: returnProductByName")

	var prodsList []Product

	for _, product := range Products {
		if product.Name == name {
			prodsList = append(prodsList, product)
		}
	}

	if len(prodsList) > 0 {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(prodsList)
	} else {
		json.NewEncoder(w).Encode("No Items found with that Name")
	}
}

func createNewProduct(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: createNewProduct")

	reqBody, _ := ioutil.ReadAll(r.Body)
	var product Product
	json.Unmarshal(reqBody, &product)
	Products = append(Products, product)

	json.NewEncoder(w).Encode("Product Create Successful")
}

func updateProductByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: updateProduct")

	reqBody, _ := ioutil.ReadAll(r.Body)
	var product Product
	json.Unmarshal(reqBody, &product)

	vars := mux.Vars(r)
	id := vars["id"]

	for index, prod := range Products {
		if prod.ID == id {
			Products[index] = product
		}
	}

	json.NewEncoder(w).Encode("Product Update Successful")
}

func deleteProductByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: deleteProductByID")

	vars := mux.Vars(r)
	id := vars["id"]

	prodIndex := 0
	for _, product := range Products {
		if product.ID == id {
			Products = append(Products[:prodIndex], Products[prodIndex+1:]...)
			json.NewEncoder(w).Encode("Product Delete Successful")
			prodIndex--

			optIndex := 0
			for _, option := range Options {
				if option.ProductID == id {
					Options = append(Options[:optIndex], Options[optIndex+1:]...)
					optIndex--
				}
				optIndex++
			}
			json.NewEncoder(w).Encode("Delete of Options Successful")
		}
		prodIndex++
	}
}
