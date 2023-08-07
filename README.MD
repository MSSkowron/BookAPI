# BookRESTAPI

BookRESTAPI is a straightforward REST API implemented in Go language, designed to facilitate book management. It offers fundamental CRUD operations on a PostgreSQL database through its HTTP API.

## Technologies

- Go 1.19
- PostgreSQL 15
- Docker

## Architecture Overview

The BookRESTAPI adopts a well-structured layered architecture to ensure clean separation of concerns and promote modularity. The system consists of the following key components:

![Architecture](./docs/architecture.png)

1. **Transport Layer**: Responsible for handling incoming HTTP requests and routing them appropriately through the Gorilla Mux router. It actively manages request authentication, authorization, and validation before onward transmission to the next layer. Implementation can be found in the `./internal/api` directory.

2. **Business Logic Layer**: The core of the application's business logic resides here. It interfaces with the data access layer to execute CRUD operations on both book and user data. Additionally, this layer governs user registration, login, and token generation functionalities. Implementation can be found in the `./internal/services` directory.

3. **Data Access Layer**: This layer shoulders the task of communicating with the PostgreSQL database. Utilizing the pgx library, it establishes connections and executes queries, abstracting the underlying database interactions from higher-level layers. Implementation can be found in the `./internal/database` directory.

Following this architectural model, the BookRESTAPI application attains crucial attributes such as separation of concerns, maintainability, and scalability. Each layer can be developed, tested, and modified independently, ensuring a well-organized and cohesive codebase.

## Key Dependencies

- `bcrypt` (golang.org/x/crypto): Used for bcrypt hashing.
- `jwt` (github.com/golang-jwt/jwt): Used for simplifying JWT token handling.
- `mux` (github.com/gorilla/mux): Used for creating the API server.
- `viper` (github.com/spf13/viper): Used for managing configurations.
- `testify` (github.com/stretchr/testify): Used for enhanced testing.
- `pqx` (github.com/jackc/pgx/v5): Used for working with PostgreSQL databases.

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

`\register` Method: `POST`

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
  "created_at": "time.Time",
  "email": "string",
  "password": "string",
  "firstName": "string",
  "lastName": "string",
  "age": "int64"
}
```

### User Login

`\login` Method: `POST`

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

### Book Management

Authentication bearer token received during login must be included to perform book management requests. Set the bearer token with the key "Authorization" in the request header for the above endpoints. The bearer authentication header should be set as follows:

```
"Authorization": "Bearer <token>"
```

Available endpoints:

- `\books` Method: `GET`

  Retrieves a list of all books.

- `\books` Method: `POST`

  Creates a new book.

  Request Body:

  ```json
  {
    "author": "string",
    "title": "string"
  }
  ```

- `\books\{id}` Method: `GET`

  Retrieves details of a specific book by ID.

- `\books\{id}` Method: `PUT`

  Updates the details of a specific book by ID.

  Request Body:

  ```json
  {
    "author": "string",
    "title": "string"
  }
  ```

- `\books\{id}` Method: `DELETE`

  Deletes a specific book by ID.

### Errors

In case of errors, the server returns an appropriate status code and JSON in the following format:

```json
{
  "error": "string"
}
```

## License

The project is available as open source under the terms of the MIT License.
