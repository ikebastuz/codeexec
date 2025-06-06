# CodeExec

A code execution service built with Go. Run code snippets in 7 different languages through a web interface or REST API.

## Features

- **Multi-Language Support**: Go, JavaScript, Python, Ruby, Rust, PHP, Java
- **Web Interface**: Clean UI with syntax highlighting
- **Docker Isolation**: Each execution runs in separate containers
- **Caching**: PostgreSQL caching for faster repeat runs
- **Metrics**: Prometheus monitoring
- **Rate Limiting**: Prevents abuse

## Tech Stack

- **Go 1.23** - Main application
- **PostgreSQL** - Caching database
- **Docker** - Code execution containers
- **Prometheus** - Metrics
- **SQLC** - Type-safe SQL generation
- **HTML + Tailwind CSS** - Frontend

## How to Run

### Prerequisites

- Go 1.23+
- Docker
- PostgreSQL
- SQLC: `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`

### Setup

```bash
# Clone and install
git clone <repo-url>
cd codeexec
go mod tidy

# Generate database code
sqlc generate

# Create database
createdb codeexec
psql -d codeexec -f schema.sql

# Set environment variables
export PORT=1450
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=your_username
export DB_PASS=your_password
export DB_NAME=codeexec

# Run
go run cmd/server/*.go
```

Visit `http://localhost:1450` to use the web interface.

## API Example

```bash
curl -X POST http://localhost:1450/run \
  -F "language=go" \
  -F "code=package main
import \"fmt\"
func main() {
    fmt.Println(\"Hello World!\")
}"
```

## Go Concepts Demonstrated

- **Dependency Injection** - Clean architecture with interface abstractions
- **Concurrent Programming** - Goroutines, channels, and synchronization
- **Docker Integration** - Container management and execution
- **Database Patterns** - SQLC code generation and caching strategies
- **HTTP Server Patterns** - Middleware, routing, and API design
- **Testing Strategies** - Unit tests with mocking and integration tests
- **Monitoring & Observability** - Prometheus metrics and structured logging
