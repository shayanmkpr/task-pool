# Task Pool

A simple task queue system written in Go with a REST API.

## Running the Server

Default configuration:

```bash
go run cmd/main.go
```

Configurable:

```bash
go run cmd/main.go \
  -pool-size=10 \
  -workers=5 \
  -port=9090 \
  -stdout-log
```

## API Endpoints

- `POST /tasks` - Create a new task
- `GET /tasks/{id}` - Get task by ID
- `GET /tasks` - Get all tasks

## Example API usage

Submit a task:

```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "title": "My Task",
    "description": "Task description"
  }'
```

Get a task by ID:

```bash
curl -X GET http://localhost:8080/tasks/{id}
```

Get all tasks:

```bash
curl -X GET http://localhost:8080/tasks
```

## Assumptions
- Uses in-memory storage (everything will be lost on restart)
- Default pool size is 10 tasks
- Default worker count is 5
- Tasks have random processing time (1-5 seconds)

## Building

```bash
go build cmd/main.go
```

## Running Unit Tests

```bash
go test ./internal/...
```

## Docker

```bash
# Build the Docker image
docker build -t task-pool .

# Run the container with default settings
docker run -p 8080:8080 task-pool

# Run the container with custom pool size, workers, port, and logging to stdout
docker run -p 8080:8080 task-pool \
  -pool-size=4 \
  -workers=3 \
  -port=8080 \
  -stdout-log
