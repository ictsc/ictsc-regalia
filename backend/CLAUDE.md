# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Building
- **Run scoreserver**: `./scripts/local-exec go run ./cmd/scoreserver/ -dev` (starts with local DB/Redis)
- **Build binary**: `go build ./cmd/scoreserver/`

### Testing
- **Run all tests**: `go test ./...`
- **Run single test**: `go test -run TestName ./path/to/package`
- Tests use Testcontainers for isolated environment

### Linting
- **Lint check**: `golangci-lint run`
- Configuration: `.golangci.yaml`

### Database
- **Access local DB**: `./scripts/local-exec psql`
- **Run migrations**: `./scripts/local-exec ./scripts/migrate`
- **Access K8s DB**: `./scripts/kube-exec psql`

### Dependencies
- **Start local services**: `docker compose up -d` (PostgreSQL, Redis, Jaeger)

### Code Generation
- **Generate code**: `go generate ./...` (runs buf for protobuf generation)

## Architecture Overview

This is a Go backend for a CTF scoring server using Connect RPC (gRPC-Web compatible).

### Core Structure
- **cmd/**: Entry points
  - `scoreserver/`: Main scoring server
  - `batch/`: Batch processing tools
  - `create-team/`: Team creation utility
  - `retrieve-token/`: Token retrieval utility

- **scoreserver/**: Main application logic
  - `domain/`: Core business entities (User, Team, Problem, Answer, Score, etc.)
  - `admin/`: Admin API handlers (Connect RPC)
  - `contestant/`: Contestant API handlers (Connect RPC)
  - `infra/`: Infrastructure implementations
    - `pg/`: PostgreSQL repositories
    - `discord/`: Discord OAuth integration
    - `sstate/`: Redis session state
  - `config/`: Configuration management
  - `batch/`: Batch job implementations

- **pkg/proto/**: Generated protobuf/Connect code
  - `admin/`: Admin API definitions
  - `contestant/`: Contestant API definitions

### Key Design Patterns
- **Repository Pattern**: Domain models have corresponding repositories in `infra/pg/`
- **Connect RPC**: API uses Connect protocol (gRPC-Web compatible) with protobuf
- **Session Management**: Redis-based sessions via `sstate` package
- **Transactions**: Uses `domain.Tx` interface for database transactions
- **Testing**: Tests use Testcontainers for isolated PostgreSQL instances

### Authentication & Authorization
- Discord OAuth2 for user authentication
- Session-based auth using Redis
- Admin/Contestant role separation via different endpoints

### Key Technologies
- **API**: Connect RPC (protobuf-based)
- **Database**: PostgreSQL with pgx driver
- **Session Store**: Redis
- **Observability**: OpenTelemetry (traces via Jaeger)
- **Testing**: Testcontainers for integration tests