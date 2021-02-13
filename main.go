package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/mux"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func handleRequests() {
	myRouter := mux.NewRouter()
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/products", returnAllProducts).Methods("GET")
	myRouter.HandleFunc("/products", createNewProduct).Methods("POST")
	myRouter.HandleFunc("/products/{id}", returnProductByID).Methods("GET")
	myRouter.HandleFunc("/products/{id}", updateProductByID).Methods("PUT")
	myRouter.HandleFunc("/products/{id}", deleteProductByID).Methods("DELETE")

	myRouter.HandleFunc("/products/{id}/options", returnOptionsByProductID).Methods("GET")
	myRouter.HandleFunc("/products/{id}/options", createNewOption).Methods("POST")
	myRouter.HandleFunc("/products/{id}/options/{optionId}", returnOptionForProduct).Methods("GET")
	myRouter.HandleFunc("/products/{id}/options/{optionId}", updateOption).Methods("PUT")
	myRouter.HandleFunc("/products/{id}/options/{optionId}", deleteOption).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func main() {
	fmt.Println("Rest API v2.0 - Mux Routers")

	// Open products data file
	productJSONFILE, err := os.Open("data/products.json")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened products.json")
	// defer file close so that file can be parsed
	defer productJSONFILE.Close()

	// read json file into byte struct
	productByteValue, _ := ioutil.ReadAll(productJSONFILE)
	// Unmarshal into to Products struct
	json.Unmarshal(productByteValue, &Products)

	// Open options data file
	optionsJSONFILE, err := os.Open("data/options.json")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened options.json")
	// defer file close so that file can be parsed
	defer optionsJSONFILE.Close()

	// read json file into byte struct
	optionsByteValue, _ := ioutil.ReadAll(optionsJSONFILE)
	// Unmarshal into to Products struct
	json.Unmarshal(optionsByteValue, &Options)

	handleRequests()
}
