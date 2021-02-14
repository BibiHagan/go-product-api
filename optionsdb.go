package main

import (
	"fmt"

	"github.com/hashicorp/go-memdb"
)

func createOptionsDatabase() *memdb.MemDB {
	//create Option Schema
	optionSchema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"option": {
				Name: "option",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "ID"},
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
	OptDB, err := memdb.NewMemDB(optionSchema)
	if err != nil {
		panic(err)
	}

	// Create a write transaction
	txn := OptDB.Txn(true)

	// Insert some products
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
	fmt.Println("Option DB created - go-memdb")

	return OptDB
}
