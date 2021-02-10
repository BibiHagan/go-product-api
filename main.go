package main

import (
	"fmt"
	"log"
	"net/http"

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
	myRouter.HandleFunc("/products/{id}", returnProductByID)
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func main() {
	fmt.Println("Rest API v2.0 - Mux Routers")
	Products = []Product{
		{
			ID:            "1",
			Name:          "Croissant",
			Description:   "Buttery french pastry",
			Price:         123.45,
			DeliveryPrice: 12.34,
		},
		{
			ID:            "2",
			Name:          "Crunt",
			Description:   "Croissant vs bundt cake",
			Price:         678.89,
			DeliveryPrice: 12.34,
		},
	}
	handleRequests()
}
