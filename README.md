#  Marvel online shop project

## Team members
Moldash Dariya 22B030560,
Bexeit Alua Armankyzy 22B030283,
Muratova Alnura 22B030399,
Izteleuova Arailym Altaikyzy 22B030369

## Project description 
Our project is an online shop specializing in the sale of comics, figures, clothing, and other products related to the world of Marvel. Users will be able to view the range of goods, add them to the basket, place orders, and track delivery status.  

## Shop REST API
```
POST /categories
GET /categories/:id
PUT /categories/:id
DELETE /categories/:id

POST /products 
GET /products/:id
PUT /products/:id
DELETE /products/:id

POST /orders
GET /orders/:id
PUT /orders/:id
DELETE /orders/:id

POST /users
GET /users/:id
PUT /users/:id
DELETE /users/:id

```

## DB Structure

```
Product: 
id (primary key) 
title 
description 
price 
category_id (foreign key to the Category table) 

Category: 
id (primary key) 
name 

User: 
id (primary key) 
username 
email 
password_hash 

Order: 
id (primary key) 
user_id (foreign key to table User) 
total_price 
status (in processing, fulfilled, etc.) 

```

## Relationships 
Each item can have only one category, but one category can have many items. 
Each order can have only one user, but one user can have many orders.

