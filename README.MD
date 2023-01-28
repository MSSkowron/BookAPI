# BookRESTAPI

BookRESTAPI is a simple REST API written in Go language to manage books. We can perform basic CRUD operations on the database via the API.

## Technologies

- Go
- PostgreSQL
- Docker

## Requirements

We need to have Docker installed in order to run the application.

## Installation

`git clone https://github.com/MSSkowron/GoBookAPI`

## How to run

```
cd GoBookAPI
docker-compose up
```

## How to use the API

GoBookAPI listens on the port 8080. It provides the following endpoints and methods:

- **POST /register**

  JSON Request Body:

  ```json
  {
    "email": "string",
    "password": "string",
    "firstName": "string",
    "lastName": "string",
    "age": "int64"
  }
  ```

  JSON Response Body:

  ```json
  {
    "id": "int64",
    "email": "string",
    "password": "string",
    "firstName": "string",
    "lastName": "string",
    "age": "int64"
  }
  ```

- **POST /login**
  JSON Request Body:

  ```json
  {
    "email": "string",
    "password": "string"
  }
  ```

  Response: authentication token of type string
  <br/>

- In methods below we need to include the authentication token which we've received while logging in. The token should be set with the key "Token" in a request header.

  - **GET /books**
  - **POST /books**
    JSON Request Body:

  ```json
  {
    "author": "string",
    "title": "string"
  }
  ```

  - **GET /books/{id}**
  - **PUT /books/{id}**
    JSON Request Body:

  ```json
  {
    "author": "string",
    "title": "string"
  }
  ```

  - **DELETE /books/{id}**
