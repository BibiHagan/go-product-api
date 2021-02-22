package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func Equal(a, b []Product) bool {
	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}

	return true
}

func TestProductsEndPoint(t *testing.T) {
	expect := []Product{
		{"1", "Croissant", "Buttery french pastry", 123.45, 12.34},
		{"2", "Crunt", "Croissant meets bundt cake", 678.89, 12.34},
	}

	ts := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Database = createDatabase()
		returnAllProducts(w, r)
	}))

	ts.EnableHTTP2 = true
	ts.StartTLS()
	defer ts.Close()

	res, err := ts.Client().Get(ts.URL)
	if err != nil {
		log.Fatal(err)
	}

	productList, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	expectString, _ := fmt.Printf("%v\n", expect)
	prodString, _ := fmt.Printf("%s\n", productList)
	if expectString == prodString {
		t.Errorf("Expected: %+v, Got: %+v", expect, productList)
	}
}

func TestReturnAllProducts(t *testing.T) {
	fmt.Println("Test GET all products")

	expect := []Product{
		{"1", "Croissant", "Buttery french pastry", 123.45, 12.34},
		{"2", "Crunt", "Croissant meets bundt cake", 678.89, 12.34},
	}

	var productList []Product
	handler := func(w http.ResponseWriter, r *http.Request) {
		Database = createDatabase()
		productList = getAllProducts(w, "id", "")
	}

	req := httptest.NewRequest("GET", "http://localhost:10000/products", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	if !Equal(expect, productList) {
		t.Errorf("Expected: %+v, Got: %+v", expect, productList)
	}
}

func TestGetSingleProduct(t *testing.T) {
	fmt.Println("Test GET single product")

	expect := Product{
		ProductID:     "2",
		Name:          "Crunt",
		Description:   "Croissant meets bundt cake",
		Price:         678.89,
		DeliveryPrice: 12.34,
	}

	var product interface{}
	handler := func(w http.ResponseWriter, r *http.Request) {
		Database = createDatabase()
		product = getSingleProduct(w, "id", "2")
	}

	req := httptest.NewRequest("GET", "http://localhost:10000/products", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	if reflect.DeepEqual(expect, product) {
		t.Errorf("Expected: %+v, Got: %+v", expect, product)
	}
}
