# Go Test Server for PR Review Agent

A simple Go HTTP server designed to test the PR review agent's support for Go language features in vector embeddings (pgvector) and knowledge graph (neo4j).

## Project Structure

```
test-repo-golang-support/
├── go.mod                    # Module definition with external dependency
├── main.go                   # Entry point, package main
├── models/
│   └── models.go             # Structs with embedding, type aliases
├── interfaces/
│   └── interfaces.go         # Interfaces with embedded interfaces
├── services/
│   └── service.go            # Methods with value/pointer receivers
├── handlers/
│   └── handlers.go           # HTTP handlers, external imports
└── README.md                 # This file
```

## Go Language Features Demonstrated

### 1. Package Declarations
Each file demonstrates proper Go package organization:
- `package main` - Application entry point
- `package models` - Data models
- `package interfaces` - Interface definitions
- `package services` - Business logic
- `package handlers` - HTTP handlers

### 2. Functions and Methods

**Value Receivers** (work on a copy):
```go
func (u User) FullName() string
func (u User) IsActive() bool
func (s UserService) Count() int
```

**Pointer Receivers** (can modify the original):
```go
func (u *User) UpdateEmail(email string)
func (u *User) SetRole(role string)
func (s *UserService) Read(ctx context.Context, id string) (*User, error)
```

### 3. Structs with Embedded Types (Composition)
```go
type BaseEntity struct {
    ID        string
    CreatedAt time.Time
}

type User struct {
    BaseEntity  // Embedded struct
    Timestamps  // Another embedded struct
    FirstName   string
    Email       string
}
```

### 4. Interfaces with Embedded Interfaces
```go
type Reader interface {
    Read(ctx context.Context, id string) (*User, error)
}

type Writer interface {
    Write(ctx context.Context, user *User) error
}

type ReadWriter interface {
    Reader  // Embedded interface
    Writer  // Embedded interface
}
```

### 5. Type Aliases
```go
type UserID = string           // Type alias (=)
type UserList = []User         // Type alias for slice
type ResponseCode int          // Type definition (not alias)
```

### 6. Import Relationships

**Standard Library:**
- `fmt` - Formatting
- `net/http` - HTTP server
- `encoding/json` - JSON encoding/decoding
- `time` - Time handling
- `context` - Context management
- `sync` - Synchronization primitives
- `log` - Logging
- `os` - OS interactions

**External Packages:**
- `github.com/gorilla/mux` - HTTP router

**Internal Packages:**
- Cross-package imports between models, services, handlers

## Running the Server

```bash
# Download dependencies
go mod tidy

# Run the server
go run main.go

# Or build and run
go build -o server
./server
```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check |
| GET | `/api/v1/users` | List all users |
| GET | `/api/v1/users/{id}` | Get user by ID |
| POST | `/api/v1/users` | Create a new user |
| PUT | `/api/v1/users/{id}` | Update a user |
| DELETE | `/api/v1/users/{id}` | Delete a user |

## Example Requests

**Create a user:**
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"first_name": "Alice", "last_name": "Wonder", "email": "alice@example.com", "role": "admin"}'
```

**Get all users:**
```bash
curl http://localhost:8080/api/v1/users
```

**Get a specific user:**
```bash
curl http://localhost:8080/api/v1/users/user_1
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | Server port |

## Testing the PR Review Agent

When you create a PR with this code, the PR review agent should:
1. Parse and extract all Go language features
2. Store embeddings for functions, structs, interfaces, etc.
3. Build knowledge graph relationships for imports, type dependencies, method receivers
4. Enable semantic search across the codebase
