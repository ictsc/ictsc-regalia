# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Monorepo Management
- **Install dependencies**: `pnpm install`
- **Build all packages**: `pnpm build`
- **Run all tests**: `pnpm ci:test`
- **Run all lints**: `pnpm ci:lint`
- **Generate protobuf code**: `pnpm generate`

### Contestant Dashboard (main app)
Run from `packages/contestant/`:
- **Development server**: `pnpm dev` (runs on localhost:3000, proxies API calls to localhost:8080)
- **Build**: `pnpm build`
- **Test**: `pnpm test` (uses Vitest)
- **Lint check**: `pnpm lint` (runs TypeScript, ESLint, and Prettier checks in parallel)
- **Fix linting**: `pnpm lint-fix`
- **Storybook**: `pnpm story` (runs on localhost:6006)
- **Build Storybook**: `pnpm story:build`

### Admin Dashboard
Run from `packages/admin/`:
- **Development server**: `pnpm dev`
- **Build**: `pnpm build`
- **Lint check**: `pnpm lint` (TypeScript, ESLint, Prettier)
- **Fix linting**: `pnpm lint-fix`

### Proto Package
Run from `packages/proto/`:
- **Generate protobuf code**: `pnpm generate` (runs buf generate and creates exports)

## Architecture Overview

This is a React frontend monorepo for the ICTSC platform using pnpm workspaces.

### Package Structure
- **packages/contestant/**: Competition participant dashboard
  - `app/components/`: Reusable UI components without business logic
  - `app/features/`: Business logic, API calls, and state management
  - `app/routes/`: Page definitions using TanStack Router (file-based routing)
  - `assets/`: Static assets
  - Uses Tailwind CSS for styling

- **packages/admin/**: Administration dashboard
  - Uses Mantine UI components
  - Similar architecture to contestant app

- **packages/proto/**: Generated protobuf/Connect RPC code
  - Exports: `@ictsc/proto/admin/v1` and `@ictsc/proto/contestant/v1`
  - Generated from proto files in backend repository

- **packages/config/**: Shared configuration
  - ESLint config: `@ictsc/config/eslint`
  - Prettier config: `@ictsc/config/prettier`

### Key Technologies
- **Framework**: React 19 with TypeScript
- **Routing**: TanStack Router with file-based routing (prefix `~` for route files)
- **API Communication**: Connect RPC (gRPC-Web compatible)
- **Styling**:
  - Contestant: Tailwind CSS
  - Admin: Mantine UI
- **Build Tool**: Vite
- **Testing**: Vitest
- **Component Development**: Storybook (contestant app only)
- **Code Quality**: ESLint, Prettier, TypeScript strict mode

### Development Patterns
- **File-based routing**: Routes are auto-generated from files in `app/routes/` prefixed with `~`
- **API proxy**: Development server proxies `/api` requests to backend at `localhost:8080`
- **Shared proto types**: Both apps use `@ictsc/proto` for type-safe API communication
- **Monorepo commands**: Root commands run recursively on all packages with `--recursive --parallel`

### Backend Integration
- Backend runs on `localhost:8080` (see `../backend/CLAUDE.md` for backend details)
- Uses Connect RPC protocol for API communication
- Authentication via Discord OAuth2 with session management