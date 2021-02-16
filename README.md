# go-product-api

These are the endpoints:

1. `GET /products` - gets all products.
2. `GET /products?name={name}` - finds all products matching the specified name.
3. `GET /products/{productId}` - gets the product that matches the specified ID - ID is a GUID.
4. `POST /products` - creates a new product.
5. `PUT /products/{productId}` - updates a product.
6. `DELETE /products/{productId}` - deletes a product and its options.

7. `GET /products/{productId}/options` - finds all options for a specified product.
8. `GET /products/{productId}/options/{optionId}` - finds the specified product option for the specified product.
9. `POST /products/{productId}/options` - adds a new product option to the specified product.
10. `PUT /products/{productId}/options/{optionId}` - updates the specified product option.
11. `DELETE /products/{productId}/options/{optionId}` - deletes the specified product option.

## All `Models` should conform to :

**Product:**
```json
{
  "Id": "01234567-89ab-cdef-0123-456789abcdef",
  "Name": "Product name",
  "Description": "Product description",
  "Price": "123.45",
  "DeliveryPrice": "12.34"
}
```

**Products:**
```json
[
  {
    // product
  },
  {
    // product
  }
]
```

**Product Option:**
```json
{
  "Id": "01234567-89ab-cdef-0123-456789abcdef",
  "ProductID": "3d820bfc-ce08-45c7-87d8-5f6d902da960",
  "Name": "Product name",
  "Description": "Product description"
}
```

**Product Options:**
```json
[
  {
    // product option
  },
  {
    // product option
  }
]
```

# Error Handling:

**Error:**
```json
{
  "Code": "400",
  "Description": "Bad Request"
}
```

# How to Install

    > go install 
- usage: go install [-i] [build flags] [packages]

compiles and installs the packages name by the import paths.
The -i flag installs the dependencies of the named packages as well
   
    > go build -i 
- usage: go build [-o output] [-i] [build flags] [packages]

Build compiles the packages named by the import paths,
along with their dependencies, but it does not install the results 
The -i flag installs the packages that are dependencies
    

# How to Run app

To start the webserver:

    > go run .

This can also be activated through the vscode launch package
Results can be viewd in either a browser or in Postman

# How to Run tests

    > go test -v

Not many tests. I used postman to create my packets and run debugging 
I chose not to mock http protocols in for testing as I wanted to 
spend more time on the program itself and the documentation