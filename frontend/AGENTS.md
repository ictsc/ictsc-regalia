# Development Commands

## Monorepo (run from root)
- `pnpm install` - Install dependencies
- `pnpm build` - Build all packages
- `pnpm ci:test` - Run all tests
- `pnpm ci:lint` - Run all lints
- `pnpm generate` - Generate protobuf code from backend definitions

## Contestant Dashboard (packages/contestant/)
- `pnpm dev` - Dev server on localhost:3000, proxies `/api` to localhost:8080
- `pnpm lint` - Check TypeScript, ESLint, and Prettier
- `pnpm lint-fix` - Auto-fix linting issues
- `pnpm story` - Run Storybook on localhost:6006

## Admin Dashboard (packages/admin/)
- `pnpm dev` - Dev server
- `pnpm lint` / `pnpm lint-fix` - Linting

# Architecture

ICTSC competition platform frontend monorepo (pnpm workspaces).

## Package Structure
- `packages/contestant/` - Participant dashboard (Tailwind CSS)
  - `app/components/` - UI components
  - `app/features/` - Business logic and API calls
  - `app/routes/` - TanStack Router pages
- `packages/admin/` - Admin dashboard (Mantine UI)
- `packages/proto/` - Generated Connect RPC code
  - Exports: `@ictsc/proto/admin/v1` and `@ictsc/proto/contestant/v1`
  - Generated from backend proto files
- `packages/config/` - Shared ESLint/Prettier configs

## Important Conventions
- **File-based routing**: Route files in `app/routes/` MUST be prefixed with `~` (e.g., `~index.tsx`)
- **API proxy**: Dev server proxies `/api` requests to `localhost:8080` backend
- **Proto imports**: Use `@ictsc/proto/contestant/v1` or `@ictsc/proto/admin/v1` for type-safe API calls

## Backend Integration
- Backend runs on `localhost:8080` (see `@../backend/CLAUDE.md`)
- Uses Connect RPC protocol (gRPC-Web compatible)
- Auth: Discord OAuth2 with session management

# Workflow
- Always run linting after making code changes
- Regenerate proto code after backend proto changes with `pnpm generate`