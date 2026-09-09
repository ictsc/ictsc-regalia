# Frontend development

The frontend is a pnpm workspace containing two Nuxt 4 / Vue / TypeScript SPAs.

## Commands (from frontend/)

- `pnpm install --frozen-lockfile`: install and prepare Nuxt types.
- `pnpm build`: check API types and generate both static SPAs.
- `pnpm ci:test`: run Vitest tests across API, UI, and both applications.
- `pnpm ci:lint`: typecheck, ESLint, and Prettier checks.
- `pnpm e2e`: contestant and admin Chromium tests with mocked APIs.
- `pnpm --filter @ictsc/ui story` / `story:build`: Vue Storybook.
- `pnpm --filter @ictsc/competition dev`: localhost:3000.
- `pnpm --filter @ictsc/admin dev`: localhost:3001/admin/.

## Architecture

- `packages/api/`: generated OpenAPI types, typed REST client, ApiError, mappers, SSE.
- `packages/ui/`: copied design assets and shared Vue components; no dependency on reference directories.
- `packages/contestant/nuxt/` and `packages/admin/nuxt/`: Nuxt source directories.
- Use Nuxt `pages/` routing, Vue composables and middleware, not React or TanStack routes.
- Both apps use `ssr: false`; deployment artifacts are `.output/public/`.

## Contracts and conventions

- `backend/openapi.json` is canonical. Run `task generate` at repository root after contract changes and include generated Go and TypeScript outputs.
- Public API is REST under `/api/v1`; no Connect RPC or Protocol Buffers.
- All REST requests use `@ictsc/api` with credentials enabled.
- Proxy `/api` unchanged to localhost:8080; never prepend the admin base path or strip `/api`.
- Deployment updates use SSE and unsubscribe on unmount. Do not add polling or refetch intervals.
- Nuxt AsyncData is shallow by default: replace the root object or use `deep: true` when updating nested SSE state.
- Convert snake_case only at mapper / feature boundaries.
- Draft storage must be scoped to contestant, team, and problem; do not migrate unowned legacy drafts.
- Team colors must be selected from the contract palette and edited only through admin APIs.
