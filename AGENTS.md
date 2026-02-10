# ICTSC Score Server

See `@backend/AGENTS.md` and `@frontend/AGENTS.md` for directory-specific commands.

# Workflow

## Protobuf Changes
- Run `task generate` from repository root (NOT from subdirectories)
- Regenerates both Go code in `backend/pkg/proto/` AND TypeScript in `frontend/packages/proto/`
