package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hashicorp/go-memdb"
)

// Database points to the Product and Option Tables
var Database *memdb.MemDB

func main() {
	fmt.Println("Rest API v2.0 - Mux Routers")

	// create DB
	Database = createDatabase()

	handleRequests()
}

func handleRequests() {
	myRouter := mux.NewRouter()
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
