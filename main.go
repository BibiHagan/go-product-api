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
	myRouter.HandleFunc("/products/{productId}", returnProductByID).Methods("GET")
	myRouter.HandleFunc("/products/{productId}", updateProductByID).Methods("PUT")
	myRouter.HandleFunc("/products/{productId}", deleteProductByID).Methods("DELETE")

	myRouter.HandleFunc("/products/{productId}/options", returnOptionsByProductID).Methods("GET")
	myRouter.HandleFunc("/products/{productId}/options", createNewOption).Methods("POST")
	myRouter.HandleFunc("/products/{productId}/options/{optionId}", returnOptionForProduct).Methods("GET")
	myRouter.HandleFunc("/products/{productId}/options/{optionId}", updateOption).Methods("PUT")
	myRouter.HandleFunc("/products/{productId}/options/{optionId}", deleteOption).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":10000", myRouter))
}
