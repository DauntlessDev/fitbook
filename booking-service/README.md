# Booking Service

This is the booking service for the Fitbook application, responsible for managing gym bookings and related events.

## Architecture

The service follows Clean Architecture principles and implements Domain-Driven Design (DDD) patterns. It uses CQRS (Command Query Responsibility Segregation) for separating read and write operations.

### Key Components

- **Domain Layer**: Core business logic and entities (bookings, events)
- **Application Layer**: Use cases and business rules
- **Interface Layer**: API handlers and repository implementations
- **Infrastructure Layer**: Technical concerns (database, cache, event publishing)

## Getting Started

### Prerequisites

- Go 1.22.1 or later
- Docker and Docker Compose
- PostgreSQL 15
- Redis 7

### Local Development

1. Clone the repository
2. Start the services:
   ```bash
   docker-compose up
   ```
3. The API will be available at `http://localhost:8080`

### API Endpoints

- `POST /api/v1/bookings`: Create a new booking
- `GET /api/v1/bookings`: List bookings
- `GET /api/v1/bookings/{gym_id}`: List bookings by gym ID
- `DELETE /api/v1/bookings/{id}`: Cancel a booking
- `GET /health`: Health check endpoint

## Project Structure

```
booking-service/
├── cmd/
│   └── api/              # Application entry point
├── internal/
│   ├── domain/          # Domain layer
│   ├── application/     # Application layer
│   ├── interfaces/      # Interface adapters
│   └── infrastructure/  # Infrastructure concerns
├── pkg/                 # Shared packages
├── api/                 # API definitions
├── configs/            # Configuration files
├── deployments/        # Deployment files
└── migrations/         # Database migrations
```

## Development

### Running Tests

```bash
go test ./...
```

### Building

```bash
go build -o booking-service ./cmd/api
```

## Deployment

The service can be deployed using the provided Docker and Kubernetes configurations in the `deployments` directory.

## License

This project is licensed under the MIT License - see the LICENSE file for details. 

