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
	myRouter.HandleFunc("/products", returnAllProducts)
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func main() {
	fmt.Println("Rest API v2.0 - Mux Routers")
	Products = []Product{
		// {Id: "01234567-89ab-cdef-0123-456789abcdef", Name: "Product name", Description: "Product description", Price: 123.45, DeliveryPrice: 12.34},
		{
			Id:            "01234567-89ef-ghij-9876-543210ghijkl",
			Name:          "Product name",
			Description:   "Product description",
			Price:         678.89,
			DeliveryPrice: 12.34,
		},
	}
	handleRequests()
}
