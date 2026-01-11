# Go Opportunities

## Overview
Go Opportunities is a project designed to manage job openings using Go and SQLite. It provides a RESTful API for creating, updating, deleting, and listing job openings.

## Project Structure
- **config/**: Contains configuration files and database connection logic.
- **db/**: Database-related files, including models and migrations.
- **docs/**: Documentation files, including Swagger specifications.
- **handler/**: Contains the HTTP handlers for managing job openings.
- **router/**: Defines the routes for the API.
- **schemas/**: Contains the data schemas for job openings.

## Getting Started
### Prerequisites
- Go 1.16 or later
- SQLite

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

### Running the Application
To run the application, use the following command:
```bash
make
```

### API Endpoints
- `POST /openings`: Create a new job opening.
- `GET /openings`: List all job openings.
- `GET /openings/{id}`: Get a specific job opening by ID.
- `PUT /openings/{id}`: Update a job opening by ID.
- `DELETE /openings/{id}`: Delete a job opening by ID.

## Contributing
Contributions are welcome! Please open an issue or submit a pull request.

## License
This project is licensed under the MIT License.