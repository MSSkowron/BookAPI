# BookRESTAPI

BookRESTAPI is a simple REST API written in Go language that allows you to manage books. It provides basic CRUD operations on a PostgreSQL database via the HTTP API.

## Technologies

- Go
- PostgreSQL
- Docker

## Requirements

To run the application, make sure you have Docker installed on your system.

## Installation

Clone the repository:

```
git clone https://github.com/MSSkowron/BookRESTAPI
```

## How to run

Navigate to the project directory:

```
cd BookRESTAPI
```

Start the application using Docker Compose:

```
docker-compose up
```

The application will be accessible at http://localhost:8080.

## How to use the API

The BookRESTAPI exposes the following endpoints and methods:

### User Registration
  Endpoint: `\register` Method: `POST`

  Registers a new user.

  Request Body:
  ```json
  {
    "email": "string",
    "password": "string",
    "firstName": "string",
    "lastName": "string",
    "age": "int64"
  }
  ```

  Response Body:
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

### User Login
  Endpoint: `\login` Method: `POST`

  Authenticates a user and returns an authentication token.

  Request Body:
  ```json
  {
    "email": "string",
    "password": "string"
  }
  ```

  Response Body:
  ```json
  {
    "token": "string"
  }
  ```
  **Note**: The authentication token received in the response must be included in the request headers for the following endpoints.


### Book Management
  - Endpoint: `\books` Method: `GET`

    Retrieves a list of all books.

  - Endpoint: `\books` Method: `POST`

    Creates a new book.

    Request Body:
    ```json
      {
        "author": "string",
        "title": "string"
      }
    ```

  - Endpoint: `\books\{id}` Method: `POST`

    Retrieves details of a specific book by ID.

  - Endpoint: `\books\{id}` Method: `PUT`

    Updates the details of a specific book by ID.
    
    Request Body:
    ```json
    {
      "author": "string",
      "title": "string"
    }
    ```

  - Endpoint: `\books\{id}` Method: `DELETE`
  
    Deletes a specific book by ID.

Remember to include the authentication token received during login by setting the token with the key "Token" in the request header for the above endpoints.
