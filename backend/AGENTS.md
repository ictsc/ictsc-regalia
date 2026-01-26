# Development Commands

## Running the Server
- **Local dev server**: `./scripts/local-exec go run ./cmd/scoreserver/ -dev -fake.schedule`
- **Build**: `go build ./cmd/scoreserver/`

## Testing
- **All tests**: `go test ./...`
- **Single test**: `go test -run TestName ./path/to/package`
- IMPORTANT: Tests use Testcontainers, so Docker must be running

## Code Quality
- **Lint**: `golangci-lint run`
- **Format**: `golangci-lint fmt`
- **Generate protobuf**: `go generate ./...`

## Database Access
- **Local DB**: `./scripts/local-exec psql`
- **Run migrations**: `./scripts/local-exec ./scripts/migrate`
- **K8s DB**: `./scripts/kube-exec psql`

## Local Services
- **Start dependencies**: `docker compose up -d` (PostgreSQL, Redis, Jaeger)

# Architecture Patterns

## Repository Pattern
- Domain models in `scoreserver/domain/` have matching repositories in `scoreserver/infra/pg/`
- Always use repositories for data access, never direct SQL in handlers

## Transactions
- Use `domain.Tx` interface for database transactions
- Pass transaction context through repository methods

## API Handlers
- Admin endpoints in `scoreserver/admin/`, contestant endpoints in `scoreserver/contestant/`
- Both use Connect RPC (gRPC-Web compatible) with protobuf definitions in `pkg/proto/`
- IMPORTANT: Admin and contestant have separate authentication/authorization flows

## Testing
- Integration tests MUST use Testcontainers for isolated PostgreSQL instances
- Never rely on shared test database state