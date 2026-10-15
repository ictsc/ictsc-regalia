# ICTSC Score Server

`backend/openapi.json` is the canonical API contract. The backend and both
frontend applications must be regenerated from that file after contract
changes.

See `@frontend/AGENTS.md` for frontend-specific commands.

# Workflow

## OpenAPI Changes
- Run `task generate` from the repository root (not from subdirectories).
- Commit the generated Go transport and TypeScript API types with the contract.
- Do not add Connect RPC or Protocol Buffers endpoints; REST paths under
  `/api/v1` are the public API.
