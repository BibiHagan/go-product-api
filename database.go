package main

import (
	"fmt"

	"github.com/hashicorp/go-memdb"
)

func createDatabase() *memdb.MemDB {
	DBSchema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"product": {
				Name: "product",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "ProductID"},
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
			"option": {
				Name: "option",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "OptionID"},
					},
					"productId": {
						Name:    "productId",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "ProductID"},
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
				},
			},
		},
	}

	// Create a new product data base
	Database, err := memdb.NewMemDB(DBSchema)
	if err != nil {
		panic(err)
	}

	// Create a write transaction
	txn := Database.Txn(true)

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

	// Insert some Options
	Options := []*Option{
		{"1", "1", "Chocolate Croissant", "Chocolate filled croissant"},
		{"2", "1", "Savoury Croissant", "Ham & cheese filled Croissane"},
		{"3", "2", "Chocolate Crunt", "Chocolate flavoured Crunt"},
		{"4", "2", "Apple Crunt", "Apple flavoured Crunt"},
	}
	for _, o := range Options {
		if err := txn.Insert("option", o); err != nil {
			panic(err)
		}
	}

	// Commit the transaction
	txn.Commit()
	fmt.Println("DB created - go-memdb")

	return Database
}
