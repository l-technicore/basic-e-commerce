# basic-e-commerce
Basic ecommerce micro services

## How to run services?
1. Users Service:
   - **go run user/user_api_launch.go**
2. Orders Service:
   - **go run order/order_api_launch.go**

## Packages:
**db** Database Connector for Oracle
	- Used the following connection url by default in both services: **admin/admin@localhost:1521/mydb**
	- i.e. user: admin, pass: admin, database:mydb on local machine

**controllers/user** handler functions and data objects for users service

**controllers/order** handler functions and data objects for orders service

## Database Schema:
1. Users Table: Users(id number, name varchar, mobile number, address varchar, PRIMARY KEY (id))
2. Orders Table: Orders(order_id number, fk_user_id number, details varchar, value number, PRIMARY KEY (order_id), FOREIGN KEY (fk_user_id) REFERENCES Orders(id))


## Services and hoe to use them:
1. Users Service (Public):
	- GET|POST /users
	- POST /users/create
		-- Post Payload sample: {"name":"Hello", "mobile":9000000001, "address":"Bell Homes"}
	- GET|POST /users/{id}
	- POST /users/{id}/update
		-- Post Payload sample: {"name":"Hello", "mobile":9000000001, "address":"Bell Homes"}
	- GET|POST /users/{id}/delete
	- GET|POST /users/{id}/orders
	- GET|POST /users/{id}/orders/{order_id}
	- POST /users/{id}/orders/{order_id}/create
		-- Post Payload sample: {"order_details":"5 pcs, Galaxy A21", "value":29999}
	- POST /users/{id}/orders/{order_id}/update
		-- Post Payload sample: {"order_details":"5 pcs, Galaxy A21", "value":29999}
	- GET|POST /users/{id}/orders/{order_id}/delete
	- GET|POST /users/{id}/delete/orders

2. Orders Service (Private):
	- GET|POST /users/{id}/orders
	- GET|POST /users/{id}/orders/{order_id}
	- POST /users/{id}/orders/{order_id}/create
		-- Post Payload sample: {"order_details":"5 pcs, Galaxy A21", "value":29999}
	- POST /users/{id}/orders/{order_id}/update
		-- Post Payload sample: {"order_details":"5 pcs, Galaxy A21", "value":29999}
	- GET|POST /users/{id}/orders/{order_id}/delete
	- GET|POST /users/{id}/delete/orders