# go-product-api

This project requires golang to be installed (see http://golang.org).

Once that's available, extract the zip file or clone the repo (http://github.com/BibiHagan/go-product-api), 
and run the following from the root of the project folder.

## How to Run

To download dependencies (see go.mod) and start the webserver:

    > go run .

This can also be activated through the Code launch package (see .vscode/launch.json).
Logs are visible at the command line (or in Code). Responses can be viewed in either a browser or in Postman.

A list of other commands is available via:

    > go help

### Endpoints

I used Postman to create my requests, test and run debugging (see products-api.postman_collection.json).

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

### Request/response Model

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

### Error Model

**Error:**
```json
{
  "Code": "400",
  "Description": "Bad Request"
}
```
