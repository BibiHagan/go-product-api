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
	myRouter.HandleFunc("/products/{id}", returnProductByID).Methods("GET")
	myRouter.HandleFunc("/products/{id}", updateProductByID).Methods("PUT")
	myRouter.HandleFunc("/products/{id}", deleteProductByID).Methods("DELETE")

	myRouter.HandleFunc("/products/{id}/options", returnOptionsByID).Methods("GET")
	myRouter.HandleFunc("/products/{id}/options/{optionId}", returnOptionForProduct).Methods("GET")

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

	Options = []Option{
		{
			ID:          "1",
			ProductID:   "1",
			Name:        "Chocolate Croissant",
			Description: "Chocolate",
		},
		{
			ID:          "2",
			ProductID:   "1",
			Name:        "Croissant",
			Description: "Plain",
		},
		{
			ID:          "3",
			ProductID:   "2",
			Name:        "Chocolate Crunt",
			Description: "Chocolate flavoured Crunt",
		},
		{
			ID:          "4",
			ProductID:   "2",
			Name:        "Apple Crunt",
			Description: "Apple flavoured Crunt",
		},
	}

	handleRequests()
}
