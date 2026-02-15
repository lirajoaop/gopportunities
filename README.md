# Go Opportunities

## Overview

Go Opportunities is a RESTful API for managing job openings, built with Go, Gin, GORM and PostgreSQL. The project follows a layered architecture (Handler → Service → Repository) with Swagger documentation.

## Architecture

```
handler/      → HTTP layer (receives requests, validates, sends responses)
service/      → Business logic (rules, orchestration)
repository/   → Data access (GORM queries to PostgreSQL)
schemas/      → Data models (Opening, OpeningResponse)
config/       → Configuration (database, logger, environment variables)
router/       → Route definitions and middleware (CORS)
docs/         → Swagger documentation (auto-generated)
```

### Request flow

```
Client → Router → Handler → Service → Repository → PostgreSQL
```

## Tech Stack

- **Language:** Go 1.25
- **Framework:** Gin
- **ORM:** GORM
- **Database:** PostgreSQL (Neon.tech)
- **Documentation:** Swagger (swag)
- **Containerization:** Docker (multi-stage build)

## Getting Started

### Prerequisites

- Go 1.25 or later
- PostgreSQL database (local or cloud)

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/lirajoaop/gopportunities.git
   cd gopportunities
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Create a `.env` file based on the example:
   ```bash
   cp .env.example .env
   ```

4. Configure your database URL in `.env`:
   ```
   DATABASE_URL=postgresql://user:password@host/database?sslmode=require
   PORT=8080
   CORS_ORIGINS=http://localhost:3000
   ```

### Running the Application

```bash
make              # runs with Swagger docs generation
make run          # runs without regenerating docs
make build        # builds the binary
make docs         # generates Swagger docs only
```

### Docker

```bash
docker build -t gopportunities .
docker run -p 8080:8080 --env-file .env gopportunities
```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/v1/openings` | List all job openings |
| `GET` | `/api/v1/openings/:id` | Get a job opening by ID |
| `POST` | `/api/v1/openings` | Create a new job opening |
| `PUT` | `/api/v1/openings/:id` | Update a job opening by ID |
| `DELETE` | `/api/v1/openings/:id` | Delete a job opening by ID |

### Swagger Documentation

After running the application, access the Swagger UI at:

```
http://localhost:8080/swagger/index.html
```

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.
