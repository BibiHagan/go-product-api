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
	myRouter := mux.NewRouter() //.strictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/products", returnAllProducts).Methods("GET")
	myRouter.HandleFunc("/products", createNewProduct).Methods("POST")
	myRouter.HandleFunc("/products/{id}", returnProductByID).Methods("GET")
	myRouter.HandleFunc("/products/{id}", updateProductByID).Methods("PUT")
	myRouter.HandleFunc("/products/{id}", deleteProductByID).Methods("DELETE")

	myRouter.HandleFunc("/products/{id}/options", returnOptionsByID).Methods("GET")
	myRouter.HandleFunc("/products/{id}/options", createNewOption).Methods("POST")
	myRouter.HandleFunc("/products/{id}/options/{optionId}", returnOptionForProduct).Methods("GET")
	myRouter.HandleFunc("/products/{id}/options/{optionId}", updateOption).Methods("PUT")
	myRouter.HandleFunc("/products/{id}/options/{optionId}", deleteOption).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func main() {
	fmt.Println("Rest API v2.0 - Mux Routers")

	// Open products data file
	productJsonFile, err := os.Open("data/products.json")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened products.json")
	defer productJsonFile.Close()

	productByteValue, _ := ioutil.ReadAll(productJsonFile)
	json.Unmarshal(productByteValue, &Products)

	// Open options data file
	optionsJsonFile, err := os.Open("data/options.json")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened options.json")
	defer optionsJsonFile.Close()

	optionsByteValue, _ := ioutil.ReadAll(optionsJsonFile)
	json.Unmarshal(optionsByteValue, &Options)

	handleRequests()
}
