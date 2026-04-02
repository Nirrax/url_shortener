# URL Shortener

A simple, efficient, and scalable URL shortener service written in Go. This service provides API endpoints to create, retrieve, and delete shortened URLs, which redirect to their original long-form counterparts. It is built with a clean, layered architecture and is fully containerized with Docker for easy setup and deployment.

## Features

*   **Create Short URLs**: Generate a unique, short identifier for any long URL.
*   **Redirect**: Automatically redirects users from the short URL to the original long URL.
*   **Retrieve URL Details**: Fetch information about a shortened URL, including the original link.
*   **Delete URLs**: Remove a shortened URL from the service.
*   **Base62 Encoding**: Uses a custom Base62 encoding for compact and URL-safe short links.
*   **Database Migrations**: Manages database schema changes using Goose.
*   **Containerized**: Easy to run and deploy using Docker and Docker Compose.

## Tech Stack

*   **Language**: Go (Golang)
*   **Database**: PostgreSQL
*   **Containerization**: Docker, Docker Compose
*   **API**: Go Standard Library (`net/http`)
*   **Database Driver**: `sqlx`, `pq`
*   **Configuration**: `viper`
*   **Migrations**: `goose`
*   **Testing**: `testify`

## Project Structure

The project follows a standard layered architecture to separate concerns and improve maintainability.

```
.
├── cmd/                # Application entry points
│   └── api/            # HTTP server, handlers, and routes
│   └── main.go         # Main application setup and startup
├── internal/           # Internal application logic
│   ├── config/         # Configuration loading (Viper)
│   ├── database/       # Database connection, migrations, and interfaces
│   ├── entity/         # Core data structures (structs)
│   ├── repository/     # Data Access Layer (interacts with the database)
│   ├── service/        # Business logic layer
│   └── utils/          # Utility functions (Base62, HTTP helpers, etc.)
├── docker-compose.yaml # Docker Compose configuration for running the service and database
├── dockerfile          # Multi-stage Dockerfile for building the application
└── .env.example        # Example environment variables file
```

## Getting Started

### Prerequisites

*   Docker
*   Docker Compose

### Installation & Running

1.  **Clone the repository:**
    ```sh
    git clone https://github.com/nirrax/url_shortener.git
    cd url_shortener
    ```

2.  **Create an environment file:**
    Copy the example environment file to create your own configuration.
    ```sh
    cp .env.example .env
    ```

3.  **Configure environment variables:**
    Open the `.env` file and set the required variables. The database variables will be used by both the Go application and the `postgres` service in `docker-compose.yaml`.

    ```env
    # Application settings
    SERVER_PORT=8080
    RUN_UP_MIGRATIONS=true

    # PostgreSQL Database settings
    DB_HOST=postgres
    DB_PORT=5432
    DB_USER=myuser
    DB_PASSWORD=mypassword
    DB_NAME=url_shortener_db
    ```
    *Note: `DB_HOST` should be the name of the PostgreSQL service defined in `docker-compose.yaml`, which is `postgres`.*

4.  **Build and run the application:**
    Use Docker Compose to build the application image and start the services.
    ```sh
    docker-compose up --build
    ```
    The API server will be running on the port specified by `SERVER_PORT` (e.g., `http://localhost:8080`).

## API Endpoints

### Create a Short URL

Creates a new short URL from a given long URL.

*   **Endpoint**: `POST /`
*   **Request Body**:
    ```json
    {
      "url": "https://example.com/a-very-long-url-that-needs-to-be-shortened"
    }
    ```
*   **Response**: `201 Created`
    ```json
    {
      "id": 1,
      "shortUrl": "1",
      "longUrl": "https://example.com/a-very-long-url-that-needs-to-be-shortened",
      "createdAt": "2024-05-21T10:00:00Z"
    }
    ```

### Redirect to Long URL

Redirects the user to the original long URL associated with the short URL.

*   **Endpoint**: `GET /{shortUrl}`
*   **Example**: `GET /1`
*   **Response**: `301 Moved Permanently` with a `Location` header pointing to the long URL.

### Get URL Details

Retrieves the details for a given short URL without redirecting.

*   **Endpoint**: `GET /url/{shortUrl}`
*   **Example**: `GET /url/1`
*   **Response**: `200 OK`
    ```json
    {
      "id": 1,
      "shortUrl": "1",
      "longUrl": "https://example.com/a-very-long-url-that-needs-to-be-shortened",
      "createdAt": "2024-05-21T10:00:00Z"
    }
    ```

### Delete a URL

Deletes a short URL from the database.

*   **Endpoint**: `DELETE /{shortUrl}`
*   **Example**: `DELETE /1`
*   **Response**: `200 OK`
    ```json
    {}
    ```

## Running Tests

To run the unit and integration tests for the project, execute the following command from the root directory:

```sh
go test ./...
