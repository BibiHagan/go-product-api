package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hashicorp/go-memdb"
)

// ProdDB points to the Product Database
var ProdDB *memdb.MemDB

// OptDB points to the Options Database
var OptDB *memdb.MemDB

func main() {
	fmt.Println("Rest API v2.0 - Mux Routers")

	// create DB
	ProdDB = createProductsDatabase()
	OptDB = createOptionsDatabase()

	handleRequests()
}

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
