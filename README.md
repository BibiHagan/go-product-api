# go-product-api

These are the endpoints:

1. `GET /products` - gets all products.
2. `GET /products?name={name}` - finds all products matching the specified name.
3. `GET /products/{id}` - gets the product that matches the specified ID - ID is a GUID.
4. `POST /products` - creates a new product.
5. `PUT /products/{id}` - updates a product.
6. `DELETE /products/{id}` - deletes a product and its options.

7. `GET /products/{id}/options` - finds all options for a specified product.
8. `GET /products/{id}/options/{optionId}` - finds the specified product option for the specified product.
9. `POST /products/{id}/options` - adds a new product option to the specified product.
10. `PUT /products/{id}/options/{optionId}` - updates the specified product option.
11. `DELETE /products/{id}/options/{optionId}` - deletes the specified product option.

## All `Models` should conform to :

**Product:**
```json
{
  "Id": "01234567-89ab-cdef-0123-456789abcdef",
  "Name": "Product name",
  "Description": "Product description",
  "Price": 123.45,
  "DeliveryPrice": 12.34
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

## Error Handling:

**Error:**
```json
{
  "Code": "400",
  "Description": "Bad Request"
}
```

**Error List**
- 400, "Unknown Params"
- 404, "Product Not Found"
- 404, "Update Fail: Product Not Found"
- 404, "Delete Fail: Product Not Found"

- 404, "Option Not Found"
- 404, "Product does not exist Option not created"
- 404, "Update error: Option Not Found"
- 404, "Delete error: Option Not Found"
