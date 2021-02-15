# Stories

## At a glance
| Story | Name | Status |
| :--- | :--- | ---: |
| 0 | Set up handler                | done |
| 1 | GET /products                 | done |
| 2 | GET /products?name={name}     | done |
| 3 | GET /products/{productId}     | done |
| 4 | POST /products                | done |
| 5 | PUT /products/{productId}     | done |
| 6 | DELETE /products/{productId}  | done |
| 7 | GET /products/{productId}/options                | done |
| 8 | GET /products/{productId}/options/{optionId}     | done |
| 9 | POST /products/{productId}/options               | done |
| 10 | PUT /products/{productId}/options/{optionId}    | done |
| 11 | DELETE /products/{productId}/options/{optionId} | done |
| 12 | REFACTOR                     | in progress |
| 13 | Design & Connect Database    |  Done |
| 14 | Write Tests                  |   |

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

### GET /products/{productId}
    - [x] gets the product that matches the specified ID 
    - [x] ID is a GUID.
    - [x] Use Postman to Test
    - [x] Update documentation

### POST /products
    - [x] creates a new product.
    - [x] Use Postman to Test
    - [x] Update documentation

### PUT /products/{productId}
    - [x] updates a product.
    - [x] Use Postman to Test
    - [x] Update documentation

### DELETE /products/{productId}
    - [x] deletes a product.
    - [x] Use Postman to Test
    - [x] Update documentation

### GET /products/{productId}/options
    - [x] finds all options for a specified product.
    - [x] Use Postman to Test
    - [x] Update documentation
    Questions:
    - Can a product exist with no Options?

### GET /products/{productId}/options/{optionId}
    - [x] finds the specified product option for the specified product.
    - [x] Use Postman to Test
    - [x] Update documentation

### POST /products/{productId}/options
    - [x] adds a new product option to the specified product.
    - [x] Use Postman to Test
    - [x] Update documentation

### PUT /products/{productId}/options/{optionId}
    - [x] updates the specified product option.
    - [x] Use Postman to Test
    - [x] Update documentation

### DELETE /products/{productId}/options/{optionId}
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

### Design & Connect Database
    - [x] design and set up database
    - [x] connect to API
    - [x] Use Postman to Test

### Write Tests
    - [-] write unit tests
    - [ ] Save Postman events
    - [ ] run tests
    - [ ] update documentation