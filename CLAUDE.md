# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Cubby is a full-stack activity tracking application (trackers, gym workouts, market prices, timers) with family sharing. Monorepo with a Go API backend and SvelteKit frontend, with Android support via Capacitor.

## Commands

### API (`/api`)

```bash
make dev          # Run API with hot reload (uses wgo)
make build        # Compile binary with stripped symbols
make lint         # Run golangci-lint
make migrate      # Run database migrations
make reset        # Reset database
```

### Web (`/web`)

```bash
pnpm dev          # Vite dev server (port 5173)
pnpm build        # Production build (static adapter)
pnpm check        # Svelte type checking
pnpm lint         # Prettier + ESLint
pnpm format       # Auto-format with Prettier
pnpm mob          # Build and run on Android (Capacitor)
pnpm mob:prod     # Production Android build
```

### Docker

```bash
docker compose up # Run full stack (API + Web + Postgres)
```

## Architecture

### Backend (`/api`)

- **Go 1.26**, PostgreSQL via `pgx/v5` + `sqlx`
- **Auth**: Stytch (magic link + OTP), JWT in cookies
- **Notifications**: Firebase Cloud Messaging
- **Entry point**: `cmd/api/main.go`, routes in `cmd/api/routes.go`
- **Internal packages**:
  - `server/` — HTTP handlers, middleware (CORS, auth, request ID), service layer
  - `tracker/`, `gym/`, `market/`, `timer/`, `user/` — domain packages
  - `database/` — Postgres connection setup
  - `migration/` — SQL migrations
  - `notifier/` — FCM push notifications
  - `response/` — HTTP response helpers and validation
- Handler pattern: domain handlers in `server/*_handlers.go`, auth via `RequireAuthentication` middleware wrapper

### Frontend (`/web`)

- **SvelteKit** (static adapter with fallback), **Svelte 5**, **TypeScript 6**
- **Styling**: Tailwind CSS 4 + DaisyUI
- **State**: TanStack SvelteQuery for server state
- **HTTP**: `ky` client configured in `src/lib/api.ts` with Capacitor cookie handling
- **Query definitions**: `src/lib/queries.ts` — all TanStack Query factories
- **Routes** (`src/routes/`):
  - `(auth)/` — login/logout
  - `app/trackers/`, `app/gym/`, `app/market/`, `app/profile/`, `app/count/` — protected app pages
- **Mobile**: Capacitor for Android, service worker for PWA/offline

### Cross-cutting

- Environment config in root `.env` (DB, API port, Stytch, Firebase credentials)
- CORS configured for dev (localhost:5173), production web, and Android emulator (10.0.2.2)
- Node >=24 required
