package main

import (
	"fmt"

	"github.com/hashicorp/go-memdb"
)

func createProductsDatabase() *memdb.MemDB {
	//create Product Schema
	productSchema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"product": {
				Name: "product",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "ID"},
					},
					"name": {
						Name:    "name",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "Name"},
					},
					"description": {
						Name:    "description",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "Description"},
					},
					"price": {
						Name:    "price",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "Price"},
					},
					"deliveryPrice": {
						Name:    "deliveryPrice",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "DeliveryPrice"},
					},
				},
			},
		},
	}

	// Create a new product data base
	ProdDB, err := memdb.NewMemDB(productSchema)
	if err != nil {
		panic(err)
	}

	// Create a write transaction
	txn := ProdDB.Txn(true)

	// Insert some products
	Products := []*Product{
		{"1", "Croissant", "Buttery french pastry", 123.45, 12.34},
		{"2", "Crunt", "Croissant meets bundt cake", 678.89, 12.34},
	}
	for _, p := range Products {
		if err := txn.Insert("product", p); err != nil {
			panic(err)
		}
	}

	// Commit the transaction
	txn.Commit()
	fmt.Println("Product DB created - go-memdb")

	return ProdDB
}
