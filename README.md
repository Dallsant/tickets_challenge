
## Tickets challenge

Tickets challenge 

## Requirements: 
- Golang 12.9.1+
- Docker
- Docker compose

## Run DB container

> docker-compose up -d

## Install Golang Dependencies
> go get -u github.com/jinzhu/gorm
> go get -u github.com/gorilla/mux
> go get -u github.com/jinzhu/gorm/dialects/postgres
> go get -u github.com/dgrijalva/jwt-go
> go get -u golang.org/x/crypto/bcrypt
> go get -u github.com/gorilla/handlers

## Run Go test server on localhost
> go run . 
.. (server will run on port 8000)

## Generate build
> go build

## Run build
> ./tickets_challenge

## Endpoints
> GET /tickets will return all tickets saved inside the DB

> GET /ticket will return a specified ticket by its id

example: GET ticket?id=1 

> POST /ticket to create a new ticket by the current logged in user (requires "Authorization" header with token)

example: Authorization eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VyX2lkIjoxLCJleHAiOjIzOTUxMDYwMjR9.yBe84hw6QkWn4kJ0VMetVSKX6Y4qpDUnLwjwGPAZpYk


> DELETE /ticket delete a specified ticket inside the DB

example: DELETE ticket?id=1 

> PUT /tickets Update a ticket's status

example: PUT ticket?id=1 
{
    status:"closed"
}
> POST /register create a new user

example:
{
    username:"eduardo",
    email:"test@gmail.com",
    password:"password"
}
> POST /login will return all tickets saved inside the DB

example:
{
    username:"eduardo",
    password:"password"
}





