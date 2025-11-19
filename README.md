## Devices API

The Devices API is a REST service for registering, querying, updating, and deleting devices. The project follows a simple, predictable, and easy-to-maintain architecture based on the separation between route, business rule, and persistence layers.

It includes automatic documentation via Swagger, uses PostgreSQL as a database, is containerized with Docker, and applies domain validations to ensure data consistency.

## Objective

The objective of this service is to provide a central point for managing "Device" resources, containing the essential operations:

- Create a new device
- Partially or fully update an existing device
- Retrieve a device by ID
- List devices with optional filters by brand and state
- Automatic pagination for unfiltered listings
- Delete devices respecting domain rules

## Architecture

The API is organized into three main layers:

Handlers (HTTP Layer)
Receive requests, validate the format, and bridge to the service. They are also responsible for serialization and status codes.

Service (Business Rule)
Centralizes all domain rules:

- devices with `state` "in-use" cannot have their `name` or `brand` changed
- "in-use" devices cannot be deleted
- `state` values must be valid
- `created_at` is not altered

Repository (Persistence)
Executes the SQL operations using PostgreSQL. Includes queries by `brand`, `state`, `ID`, and paginated listing.

This separation facilitates testing, project evolution, and overall clarity.

## Routes

The API exposes the following endpoints:

- `POST /devices`
- `GET /devices`
- `GET /devices/{id}`
- `PATCH /devices/{id}`
- `DELETE /devices/{id}`

Detailed documentation is available via Swagger.

## Documentation (Swagger)

After starting the application, the documentation will be available at:
`http://localhost:8080/swagger/index.html`

If it is necessary to regenerate the files:
`swag init -g cmd/server/main.go -o internal/docs`

## Domain Rules

The rules applied in the service are:
- `created_at` cannot be modified under any circumstances.
- Devices in `state` "in-use" cannot have their `name` or `brand` altered.
- Devices in `state` "in-use" cannot be deleted.
- The `state` field only accepts the values: available, in-use, inactive.
- `name` and `brand` are mandatory upon creation.

## How to Run Without Docker

It is necessary to have Go 1.23+ and PostgreSQL installed.

1.  `export DATABASE_URL="postgres://postgres:postgres@localhost:5432/devices?sslmode=disable"`
2.  `go run cmd/server/main.go`

The API will be available at:
`http://localhost:8080`

## How to Run With Docker

The project already includes a docker-compose configured with PostgreSQL and a healthcheck. To bring everything up:
`docker-compose up --build`

The API will be available at:
`http://localhost:8080`

And PostgreSQL at:
`localhost:5432`

## Tests

The service layer is testable through mocks.
To run the tests:

`go test ./...`

## Future Improvements

Some possible improvements:

- Integration tests using an isolated database
- API versioning (e.g., `/api/v1`)
- More structured error handling
- Inclusion of `PUT` for complete updates
- Expansion of filters in the listing