# Development Commands

- `export PATH="$HOME/.bun/bin:$PATH"; bun install` - Install dependencies
- `export PATH="$HOME/.bun/bin:$PATH"; bun run dev` - Start Nuxt dev server
- `export PATH="$HOME/.bun/bin:$PATH"; bun run build` - Build the Nuxt app
- `export PATH="$HOME/.bun/bin:$PATH"; bun run typecheck` - Run Nuxt type checks

# Architecture

This frontend is a Bun-managed Nuxt app.

- `app.vue` - Root application view
- `plugins/vue-query.ts` - TanStack Query client registration
- `nuxt.config.ts` - Nuxt configuration

# Conventions

- Use TanStack Query through `@tanstack/vue-query` composables.
- Keep API base URL configuration in Nuxt runtime config.
