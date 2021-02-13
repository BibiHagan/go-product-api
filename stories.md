# Stories

## At a glance
| Story | Name | Status |
| :--- | :--- | ---: |
| 0 | Set up handler            | done |
| 1 | GET /products             | done |
| 2 | GET /products?name={name} | done |
| 3 | GET /products/{id}        | done |
| 4 | POST /products            | done |
| 5 | PUT /products/{id}        | done |
| 6 | DELETE /products/{id}     | done |
| 7 | GET /products/{id}/options                | done |
| 8 | GET /products/{id}/options/{optionId}     | done |
| 9 | POST /products/{id}/options               | done |
| 10 | PUT /products/{id}/options/{optionId}    | done |
| 11 | DELETE /products/{id}/options/{optionId} | done |
| 12 | REFACTOR                 | in progress |
| 13 | Vet Data                 |  |
|  | Design & Connect Database     |   |

## In detail
### Set up handler
    - [x] Set up Listen and Serve.
    - [x] Use Postman to Test
    - [x] Start documentation

### GET /products
    - [x] gets all products.
    - [x] Use Postman to Test
    - [x] Update documentation

### GET /products?name={name}
    - [x] finds all products matching the specified name.
    - [x] Use Postman to Test
    - [x] Update documentation

### GET /products/{id}
    - [x] gets the product that matches the specified ID 
    - [x] ID is a GUID.
    - [x] Use Postman to Test
    - [x] Update documentation

### POST /products
    - [x] creates a new product.
    - [x] Use Postman to Test
    - [x] Update documentation

### PUT /products/{id}
    - [x] updates a product.
    - [x] Use Postman to Test
    - [x] Update documentation

### DELETE /products/{id}
    - [x] deletes a product.
    - [x] Use Postman to Test
    - [x] Update documentation

### GET /products/{id}/options
    - [x] finds all options for a specified product.
    - [x] Use Postman to Test
    - [x] Update documentation
    Questions:
    - Can a product exist with no Options?

### GET /products/{id}/options/{optionId}
    - [x] finds the specified product option for the specified product.
    - [x] Use Postman to Test
    - [x] Update documentation

### POST /products/{id}/options
    - [x] adds a new product option to the specified product.
    - [x] Use Postman to Test
    - [x] Update documentation

### PUT /products/{id}/options/{optionId}
    - [x] updates the specified product option.
    - [x] Use Postman to Test
    - [x] Update documentation

### DELETE /products/{id}/options/{optionId}
    - [x] delete the specified product option.
    - [x] Use Postman to Test
    - [x] delete all options if product is deleted
    - [x] Use Postman to Test
    - [x] Update documentation

### Error handling
    - [x] create error type.
    - [x] Use Postman to Test
    - [x] Update documentation

### REFACTOR
    - [ ] convert all slices to maps
    - [ ] Use Postman to Test
    - [ ] Update documentation

### Vet Data
    - [ ] add checks to make sure all data coming in is valid
    - [ ] Use Postman to Test
    - [ ] Update documentation

### Design & Connect Database
    - [ ] design and set up database
    - [ ] connect to API
    - [ ] Test

