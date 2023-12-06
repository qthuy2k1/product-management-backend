# API

## **User APIs**

1. **GetUser** (Method: GET)

    - **Success**
        * URL: localhost:3000/users/1
        * Status code: 200 OK
        * Result:
            {
                "id": 1,
                "name": "Thuy Nguyen",
                "email": "qthuy@gmail.com",
                "password": "$2a$14$1VRhclX4P/JkiqJ7nKXYvuC/4NdvKXreBPay9sDJv1CpWs4eFhXAC",
                "role": "user",
                "status": "activated",
                "created_at": "2023-05-11T09:01:53.102071Z",
                "updated_at": "2023-05-11T09:01:53.102071Z"
            }

    - **Errors**
        1. User not found:
            * URL: localhost:3000/users/100000
            * Status code: 404 Not Found
            * Result:
                {
                    "message": "user not found"
                }

        2. Invalid user ID
            * URL: localhost:3000/users/-1
            * Status code: 400 Bad Request
            * Result:
                {
                    "message": "invalid ID"
                }


2. **CreateUser** (Method: POST)

    - **Success**
        * URL: localhost:3000/users/
        * Status code: 201 Created
        * Input:
            {
                "name": "Quang Thuy",
                "email": "qthuy@example.com",
                "password": "123123"
            }
        * Result:
            {
                "success": true
            }

    - **Errors**
        1. Invalid json:
            * URL: localhost:3000/users/
            * Status code: 400 Bad Request
            * Input:
                {
                    "name": "Quang Thuy",
                    "email": "qthuy@example.com",
                    "password": "123123 // missing end of "
                }
            * Result:
                {
                    "message": "invalid json"
                }

        2. Missing name field
            * URL: localhost:3000/users/
            * Status code: 400 Bad Request
            * Input:
                {
                    "email": "qthuy@example.com",
                    "password": "123123"
                }
            * Result:
                {
                    "message": "name cannot be blank"
                }

        3. Missing email field
            * URL: localhost:3000/users/
            * Status code: 400 Bad Request
            * Input:
                {
                    "name": "Quang Thuy",
                    "password": "123123"
                }
            * Result:
                {
                    "message": "invalid email"
                }

        4. Missing password field
            * URL: localhost:3000/users/
            * Status code: 400 Bad Request
            * Input:
                {
                    "email": "qthuy@example.com",
                    "name": "Quang Thuy",
                }
            * Result:
                {
                    "message": "password must between 6 and 72 characters"
                }

        5. Invalid password
            * URL: localhost:3000/users/
            * Status code: 400 Bad Request
            * Input:
                {
                    "email": "qthuy@example.com",
                    "name": "Quang Thuy",
                    "password": "12345"
                }
            * Result:
                {
                    "message": "password must be between 6 and 72 characters"
                }

        6. Wrong type input
            * URL: localhost:3000/users/
            * Status code: 400 Bad Request
            * Input:
                {
                    "email": "qthuy@example.com",
                    "name": "Quang Thuy",
                    "password": 123456 // this should be a string
                }
            * Result:
                {
                    "message": "invalid json"
                }
        7. Duplicate email
            * URL: localhost:3000/users/
            * Status code: 400 Bad Request
            * Input:
                {
                    "email": "qthuy@example.com",
                    "name": "Quang Thuy",
                    "password": 123456"
                }
            * Result:
                {
                    "message": "email already exists"
                }

        8. Internal server error
            * URL: localhost:3000/users/
            * Status code: 500 Internal Server Error
            * Input:
                {
                    "email": "qthuy@example.com",
                    "name": "Quang Thuy",
                    "password": "123456"
                }
            * Result:
                {
                    "message": "internal server error"
                }


 

## **Product Category APIs**
1. **CreateProductCategory** (Method: POST)

    - **Success**
        * URL: localhost:3000/product-categories/
        * Status code: 201 Created
        * Input:
            {
                "name": "Cellphone",
                "description" "Cellphone"
            }
        * Result:
            {
                "success": true
            }

    - **Errors**
        1. Invalid json:
            * URL: localhost:3000/product-categories/
            * Status code: 400 Bad Request
            * Input:
                {
                    "name": "Cellphone",
                    "description": "Cellphone // missing end of "
                }
            * Result:
                {
                    "message": "invalid json"
                }

        2. Missing name field
            * URL: localhost:3000/product-categories/
            * Status code: 400 Bad Request
            * Input:
                {
                    "description" "Cellphone"
                }
            * Result:
                {
                    "message": "name cannot be blank"
                }

        3. Missing description field
            * URL: localhost:3000/product-categories/
            * Status code: 400 Bad Request
            * Input:
                {
                    "name" "Cellphone"
                }
            * Result:
                {
                    "message": "description cannot be blank"
                }

        4. Wrong type input
            * URL: localhost:3000/product-categories/
            * Status code: 400 Bad Request
            * Input:
                {
                    "name": "Cellphone",
                    "description:: 123456 // this should be a string
                }
            * Result:
                {
                    "message": "invalid json"
                }

        5. Internal server error
            * URL: localhost:3000/product-categories/
            * Status code: 500 Internal Server Error
            * Input:
                {
                    "name": "Cellphone",
                    "description": "Cellphone"
                }
            * Result:
                {
                    "message": "internal server error"
                }


## **Product APIs**
1. **CreateProduct** (Method: Post)

    - **Success**
        * URL: localhost:3000/products/
        * Status code: 201 Created
        * Input:
            {
                "name": "iPhone 14",
                "description": "An Apple cellphone with A15 Bionic chip, 6GB RAM and 128GB storage",
                "price": 1500,
                "quantity": 10,
                "category_id": 1,
                "author_id":1
            }
        * Result: 
            {
                "success": true
            }

    - **Errors**
        1. Invalid json:
            * URL: localhost:3000/product/
            * Status code: 400 Bad Request
            * Input:
                {
                    "name": "iPhone 14",
                    "description": "An Apple cellphone, // missing end of "
                    "price": 1500,
                    "quantity": 15,
                    "category_id": 1,
                    "author_id":1
                }
            * Result:
                {
                    "message": "invalid json"
                }

        2. Missing name and description field
            * URL: localhost:3000/product/
            * Status code: 400 Bad Request
            * Input:
                {
                    "price": 1500,
                    "quantity": 10,
                    "category_id": 1,
                    "author_id":1
                }
            * Result:
                {
                    "message": "name cannot be blank"
                }

        3. Value in quantity field is less than 0
            * URL: localhost:3000/product/
            * Status code: 400 Bad Request
            * Input:
                {
                    "name": "iPhone 14",
                    "description": "An Apple cellphone with A15 Bionic chip, 6GB RAM and 128GB storage",
                    "price":1500,
                    "quantity": -1,
                    "category_id": 1,
                    "author_id":1
                }
            * Result:
                {
                    "message": "quantity must be greater than 0"
                }

        4. Wrong type input
            * URL: localhost:3000/product/
            * Status code: 400 Bad Request
            * Input:
                {
                    "name": 123123, // this should be string
                    "description": "An Apple cellphone with A15 Bionic chip, 6GB RAM and 128GB storage",
                    "price": 1500,
                    "quantity": 10,
                    "category_id": 1,
                    "author_id":1
                }
            * Result:
                {
                    "message": "invalid json"
                }

        5. Internal server error
            * URL: localhost:3000/product/
            * Status code: 500 Internal Server Error
            * Input:
                {
                    "name": "iPhone 14",
                    "description": "An Apple cellphone with A15 Bionic chip, 6GB RAM and 128GB storage",
                    "price": 1500,
                    "quantity": 10,
                    "category_id": 1,
                    "author_id":1
                }
            * Result:
                {
                    "message": "internal server error"
                }


2. **GetProduct** (Method: Get)

    - **Success**
        * URL: localhost:3000/products/1
        * Status code: 200 OK
        * Result:
            {
                "name": "iPhone 14",
                "description": "An Apple cellphone with A15 Bionic chip, 6GB RAM and 128GB storage",
                "price": 1500,
                "quantity": 10,
                "category_id": 1,
                "author_id":1
            }

    - **Errors**
        1. Product not found:
            * URL: localhost:3000/products/100000
            * Status code: 404 Not Found
            * Result:
                {
                    "message": "product not found"
                }

        2. Invalid product ID
            * URL: localhost:3000/products/-1
            * Status code: 400 Bad Request
            * Result:
                {
                    "message": "invalid ID"
                }

3. **DeleteProduct** (Method: Delete)

    - **Success**
        * URL: localhost:3000/products/1
        * Status code: 200 OK
        * Result:
            {
                "success": true
            }
        
    - **Errors**
        1. Product not found:
            * URL: localhost:3000/products/100000
            * Status code: 404 Not Found
            * Result:
                {
                    "message": "product not found"
                }

        2. Invalid product ID
            * URL: localhost:3000/products/-1
            * Status code: 400 Bad Request
            * Result:
                {
                    "message": "invalid ID"
                }
        3. Internal server error
            * URL: localhost:3000/product/1
            * Status code: 500 Internal Server Error
            * Result:
                {
                    "message": "internal server error"
                }

4. **UpdateProduct** (Method: Put)

    - **Success**
        * URL: localhost:3000/products/1
        * Status code: 200 OK
        * Input:
            {
                "name": "iPhone 14",
                "description": "An Apple cellphone with A15 Bionic chip, 6GB RAM and 128GB storage",
                "price": 1500,
                "quantity": 15,
                "category_id": 1,
                "author_id":1
            }
        * Result: 
            {
                "success": true
            }

    - **Errors**
        1. Invalid json:
            * URL: localhost:3000/product/1
            * Status code: 400 Bad Request
            * Input:
                {
                    "name": "iPhone 14",
                    "description": "An Apple cellphone, // missing end of "
                    "price": 1500,
                    "quantity": 15,
                    "category_id": 1,
                    "author_id":1
                }
            * Result:
                {
                    "message": "invalid json"
                }

        2. Missing name
            * URL: localhost:3000/product/1
            * Status code: 400 Bad Request
            * Input:
                {
                    "price": 1500,
                    "quantity": 10,
                    "category_id": 1,
                    "author_id":1
                }
            * Result:
                {
                    "message": "name cannot be blank"
                }

        3. Value in quantity field is less than 0
            * URL: localhost:3000/product/1
            * Status code: 400 Bad Request
            * Input:
                {
                    "name": "iPhone 14",
                    "description": "An Apple cellphone with A15 Bionic chip, 6GB RAM and 128GB storage",
                    "price": 1500,
                    "quantity": -1,
                    "category_id": 1,
                    "author_id":1
                }
            * Result:
                {
                    "message": "quantity must be greater than 0"
                }

        4. Wrong type input
            * URL: localhost:3000/product/1
            * Status code: 400 Bad Request
            * Input:
                {
                    "name": 123123, // this should be string
                    "description": "An Apple cellphone with A15 Bionic chip, 6GB RAM and 128GB storage",
                    "price": 1500,
                    "quantity": 10,
                    "category_id": 1,
                    "author_id":1
                }
            * Result:
                {
                    "message": "invalid json"
                }

        5. Internal server error
            * URL: localhost:3000/product/1
            * Status code: 500 Internal Server Error
            * Input:
                {
                    "name": "iPhone 14",
                    "description": "An Apple cellphone with A15 Bionic chip, 6GB RAM and 128GB storage",
                    "price": 1500,
                    "quantity": 10,
                    "category_id": 1,
                    "author_id":1
                }
            * Result:
                {
                    "message": "internal server error"
                }
