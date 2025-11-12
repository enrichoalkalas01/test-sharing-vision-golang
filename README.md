# Test Sharing Vision - Golang REST API

A robust RESTful API built with Go (Golang) and Fiber framework for managing articles. This project implements Clean Architecture principles and supports multiple database options.

## Features

- **Article Management**: Full CRUD operations for articles
- **Clean Architecture**: Separation of concerns with domain, repository, usecase, and handler layers
- **Multiple Database Support**: MySQL, PostgreSQL, SQLite, SQL Server, MongoDB, Cassandra, Redis, Elasticsearch
- **Structured Logging**: Using Uber Zap logger
- **Configuration Management**: Viper for flexible configuration
- **Middleware Support**: Request ID, Recovery, and custom middleware
- **Input Validation**: Built-in validation for article entities
- **Filtering**: Support for filtering, sorting, and pagination

## Tech Stack

- **Framework**: [Fiber v2](https://gofiber.io/) - Express-inspired web framework
- **ORM**: [GORM](https://gorm.io/) - Go Object Relational Mapping
- **Logger**: [Zap](https://github.com/uber-go/zap) - Structured logging
- **Configuration**: [Viper](https://github.com/spf13/viper) - Configuration management
- **Database Drivers**:
  - MySQL/MariaDB
  - PostgreSQL
  - SQLite
  - SQL Server
  - MongoDB
  - Cassandra
  - Redis
  - Elasticsearch

## Project Structure

```
.
├── cmd/
│   └── api/
│       └── server.go           # Application entry point
├── configs/
│   └── viper.go               # Configuration loader
├── internal/
│   ├── domain/
│   │   └── article.go         # Domain entities and business rules
│   ├── dto/
│   │   └── article/           # Data Transfer Objects
│   │       ├── article_request.go
│   │       ├── article_response.go
│   │       ├── article_filter.go
│   │       └── errors.go
│   ├── handler/
│   │   └── article_handler.go # HTTP handlers
│   ├── repository/
│   │   └── article/           # Data access layer
│   │       ├── article_interface.go
│   │       └── article_repository.go
│   ├── routes/
│   │   ├── route.go           # Route configuration
│   │   └── article_routes.go  # Article routes
│   └── usecase/
│       └── article/           # Business logic layer
│           ├── article_interface.go
│           └── article_usecase.go
├── pkg/
│   ├── common/
│   │   ├── filter/            # Filter utilities
│   │   └── response/          # Response utilities
│   ├── database/              # Database connections
│   │   ├── cassandra.go
│   │   ├── elasticsearch.go
│   │   ├── mongodb.go
│   │   ├── mssql.go
│   │   ├── mysql.go
│   │   ├── postgres.go
│   │   ├── redis.go
│   │   └── sqlite.go
│   ├── logger/
│   │   └── logger.go          # Logger configuration
│   ├── middlewares/           # HTTP middlewares
│   │   ├── recovery.go
│   │   └── request_id.go
│   └── servers/               # Server setup
│       ├── fiber.go
│       ├── middleware.go
│       └── routes.go
├── .env.example               # Environment variables example
├── .gitignore                 # Git ignore rules
├── Dockerfile                 # Docker image definition
├── docker-compose.yml         # Docker compose configuration
├── go.mod                     # Go module dependencies
├── go.sum                     # Go module checksums
└── makefile                   # Build and development commands
```

## Getting Started

### Prerequisites

#### For Local Development:
- Go 1.24.0 or higher
- MySQL (or other supported database)
- Make (optional, for using Makefile commands)

#### For Docker Deployment:
- Docker Engine 20.10+
- Docker Compose 2.0+

### Installation

1. Clone the repository:
```bash
git clone https://github.com/enrichoalkalas01/test-sharing-vision-golang.git
cd test-sharing-vision-golang
```

2. Install dependencies:
```bash
make deps
# or
go mod download && go mod tidy
```

3. Copy the environment variables file:
```bash
cp .env.example .env
```

4. Configure your environment variables in `.env`:
```env
APP_NAME=sharing-vision-golang
APP_ENV=development
APP_VERSION=1.0.0

# Database Configuration
MYSQL_HOST=localhost
MYSQL_PORT=3306
MYSQL_USER=your_username
MYSQL_PASSWORD=your_password
MYSQL_DATABASE=your_database
MYSQL_MAX_IDLE_CONNS=10
MYSQL_MAX_OPEN_CONNS=100
MYSQL_SSL_MODE=false
```

### Running the Application

#### Development Mode (with hot-reload):
```bash
make dev
```

#### Production Mode:
```bash
make build
./bin/server.exe
```

#### Run without building:
```bash
make start
# or
go run cmd/api/server.go
```

## Docker Deployment

This project includes Docker support for easy deployment and development.

### Prerequisites for Docker

- Docker Engine 20.10+
- Docker Compose 2.0+

### Using Docker Compose (Recommended)

The easiest way to run the application with all dependencies:

1. Make sure you have `.env` file configured (or use default values in docker-compose.yml)

2. Start all services:
```bash
docker-compose up -d
```

3. View logs:
```bash
docker-compose logs -f app
```

4. Stop all services:
```bash
docker-compose down
```

5. Stop and remove volumes (caution: deletes database data):
```bash
docker-compose down -v
```

### Services Included

When you run `docker-compose up`, the following services will be started:

- **app**: The Go API application (port 3000)
- **mysql**: MySQL 8.0 database (port 3306)
- **phpmyadmin**: Web interface for MySQL management (port 8080)

Access the services:
- API: http://localhost:3000
- phpMyAdmin: http://localhost:8080

### Building Docker Image Only

If you want to build the Docker image without using docker-compose:

```bash
# Build the image
docker build -t sharing-vision-api .

# Run the container (you'll need to provide database connection separately)
docker run -p 3000:3000 \
  -e MYSQL_HOST=your_mysql_host \
  -e MYSQL_USER=your_user \
  -e MYSQL_PASSWORD=your_password \
  -e MYSQL_DATABASE=your_database \
  sharing-vision-api
```

### Docker Environment Variables

You can customize the deployment by setting environment variables in your `.env` file:

```env
# Application
APP_NAME=sharing-vision-golang
APP_ENV=production
APP_VERSION=1.0.0

# Server
SERVER_PORT=3000
SERVER_HOST=0.0.0.0

# MySQL
MYSQL_USER=appuser
MYSQL_PASSWORD=your_secure_password
MYSQL_DATABASE=sharing_vision

# Logging
LOG_LEVEL=info
```

### Production Deployment Tips

1. **Change default passwords**: Update `MYSQL_PASSWORD` in `.env` file
2. **Use secrets management**: For production, use Docker secrets or external secret managers
3. **Configure volumes**: Persist data using named volumes (already configured in docker-compose.yml)
4. **Health checks**: The Dockerfile includes health checks for monitoring
5. **Resource limits**: Add resource limits in docker-compose.yml for production:
```yaml
deploy:
  resources:
    limits:
      cpus: '0.5'
      memory: 512M
```

### Troubleshooting Docker

If the app container fails to start:
```bash
# Check logs
docker-compose logs app

# Check MySQL is ready
docker-compose logs mysql

# Restart services
docker-compose restart
```

If you need to rebuild after code changes:
```bash
docker-compose up -d --build
```

## API Endpoints

### Articles

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/article` | Get list of articles (with filtering, sorting, pagination) |
| POST | `/article` | Create a new article |
| GET | `/article/:article_id` | Get article by ID |
| PUT | `/article/:article_id` | Update article by ID |
| DELETE | `/article/:article_id` | Delete article by ID |

### Article Schema

```json
{
  "id": 1,
  "title": "Article Title",
  "content": "Article content goes here...",
  "category": "Technology",
  "status": "Publish",
  "created_date": "2025-01-01T00:00:00Z",
  "updated_date": "2025-01-01T00:00:00Z"
}
```

#### Status Values:
- `Publish` - Published article
- `Draft` - Draft article
- `Thrash` - Trashed article

### Example Requests

#### Create Article
```bash
curl -X POST http://localhost:3000/article \
  -H "Content-Type: application/json" \
  -d '{
    "title": "My First Article",
    "content": "This is the content of my article",
    "category": "Technology",
    "status": "Publish"
  }'
```

#### Get All Articles
```bash
curl http://localhost:3000/article
```

#### Get Article by ID
```bash
curl http://localhost:3000/article/1
```

#### Update Article
```bash
curl -X PUT http://localhost:3000/article/1 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Updated Title",
    "content": "Updated content",
    "category": "Technology",
    "status": "Publish"
  }'
```

#### Delete Article
```bash
curl -X DELETE http://localhost:3000/article/1
```

## Validation Rules

Articles must meet the following criteria:

- **Title**: Required, 3-200 characters
- **Content**: Required, minimum 10 characters
- **Category**: Required
- **Status**: Must be one of: `Publish`, `Draft`, `Thrash`

## Development

### Available Make Commands

```bash
make dev       # Run with hot-reload (requires Air)
make build     # Build production binary
make start     # Run without hot-reload
make test      # Run tests
make clean     # Clean build artifacts
make deps      # Install dependencies
make fmt       # Format code
make lint      # Run linter (requires golangci-lint)
```

### Code Formatting
```bash
make fmt
# or
go fmt ./...
```

### Running Tests
```bash
make test
# or
go test -v ./...
```

## Architecture

This project follows **Clean Architecture** principles:

1. **Domain Layer**: Contains business entities and rules
2. **Repository Layer**: Handles data access and persistence
3. **Use Case Layer**: Implements business logic
4. **Handler Layer**: Handles HTTP requests and responses

### Dependency Flow:
```
Handler → Use Case → Repository → Domain
```

## Logging

The application uses Uber Zap for structured logging. Logs include:
- Request IDs for tracing
- Structured fields for better searchability
- Different log levels (debug, info, warn, error)

## Middleware

- **Request ID**: Adds unique ID to each request for tracing
- **Recovery**: Recovers from panics and returns proper error responses
- Custom middleware support available

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License.

## Contact

Enrico Alkalas - [@enrichoalkalas01](https://github.com/enrichoalkalas01)

Project Link: [https://github.com/enrichoalkalas01/test-sharing-vision-golang](https://github.com/enrichoalkalas01/test-sharing-vision-golang)
