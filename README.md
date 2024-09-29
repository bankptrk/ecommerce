#Starter Project

```
docker-compose up -d
go run .
```

- REGISTER (POST METHOD)

  http://localhost:8000/users/register

```
{
    "first_name":"bank",
    "last_name":"eiei",
    "email":"bankptrk@gmail.com",
    "phone":"095",
    "password":"1234"
}

```
Respone should be "User registered successfully!"

If user already exists respone should be "Email already exists"

- LOGIN (POST METHOD)
  
  http://localhost:8000/users/login

```
{
    "email":"bankptrk@gmail.com",
    "password":"1234"
}
```

- CREATE PRODUCT (POST METHOD)
http://localhost:8000/admin/addproduct

```
{
  "product_name":"example product",
  "price": 1234.5,
  "rating":5,
  "inage":"example.png"
}
```

- GET ALL PRODUCT (GET METHOD)
http://localhost:8000/users/products
respone should be :

```
[
    {
        "Product_ID": "66f79d2e074e13a1a7019687",
        "product_name": "Example Product",
        "price": 100,
        "rating": 4,
        "image": "https://example.com/image.png"
    },
    {
        "Product_ID": "66f79d32074e13a1a7019688",
        "product_name": "Example Product",
        "price": 100,
        "rating": 4,
        "image": ""
    },
    {
        "Product_ID": "66f8e116adce4c7ebed0ac55",
        "product_name": "example product",
        "price": 1234.5,
        "rating": 5,
        "image": ""
    }
]
```

- GET PRODUCT BY NAME (GET METHOD)
ex : http://localhost:8000/users/search?name=Example --> (partial search)
respone should be :

```
[
    {
        "Product_ID": "66f79d2e074e13a1a7019687",
        "product_name": "Example Product",
        "price": 100,
        "rating": 4,
        "image": "https://example.com/image.png"
    },
    {
        "Product_ID": "66f79d32074e13a1a7019688",
        "product_name": "Example Product",
        "price": 100,
        "rating": 4,
        "image": ""
    }
]

```

- ADD ADDRESS (POST METHOD)
ex : http://localhost:8000/addaddress?id=xxxx

```
{
    "house_name": "4/40",
    "street_name": "123 Main St",
    "city_name": "Bangkok",
    "zip_code": "11130"
}

```
  
- UPDATE ADDRESS FOR BILLING (PUT METHOD)
ex : http://localhost:8000/users/addresses/billing?id=xxxxxx

```
{
    "house_name": "4/41",
    "street_name": "123st",
    "city_name": "Nonthaburi",
    "zip_code": "11130"
}
```

- UPDATE ADDRESS FOR SHIPPING (PUT METHOD)
  ex : http://localhost:8000/users/addresses/shipping?id=xxxxxx
  
```
{
    "house_name": "4/40",
    "street_name": "123 Main St",
    "city_name": "Bangkok",
    "zip_code": "11130"
}

```

- DELETE ADDRESS (DELETE METHOD)
ex : http://localhost:8000/users/addresses?id=xxxxxx


- ADD PRODUCT TO CART (POST METHOD)
ex : http://localhost:8000/users/cart?pid=xxxxxx&uid=xxxxxx


- DELETE ITEM FROM CART (DELETE METHOD)
ex : http://localhost:8000/users/cart/item?pid=xxxxxx&uid=xxxxxx


- GET ITEM FROM CART (GET METHOD)
ex : http://localhost:8000/users/cart?id=xxxxxx
respone should be :

```
{
    "total": 2469,
    "userCart": [
        {
            "Product_ID": "66f8e116adce4c7ebed0ac55",
            "product_name": "example product",
            "price": 1234.5,
            "rating": 5,
            "image": ""
        },
        {
            "Product_ID": "66f8e116adce4c7ebed0ac55",
            "product_name": "example product",
            "price": 1234.5,
            "rating": 5,
            "image": ""
        }
    ]
}

```

- BUY PRODUCT FROM CART (POST METHOD)
ex : http://localhost:8000/users/cart/purchase?id=xxxxxx


- BUY PRODUCT INSTANTBUY (POST METHOD)
ex : http://localhost:8000/users/cart/instant-buy?uid=xxxxxx&pid=xxxxxx
