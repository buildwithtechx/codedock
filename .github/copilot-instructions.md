# Codedock — Agent Instructions

## What This Is

This is the Codedock codebase — an open-source, self-hosted PaaS that turns any bare-metal VPS into a private Vercel & Railway. It's written in Go (backend daemon) and TypeScript (React 19 control panel dashboard in `dashboard/`).

- **Language (backend)**: Go — `cmd/`, `internal/`, `pkg/`
- **Language (frontend)**: TypeScript/TSX — React 19 in `dashboard/`
- **Dashboard Stack**: React 19 + TanStack Router + Vite + Tailwind CSS v4
- **Database**: embedded SQLite (`modernc.org/sqlite`, CGO-free)
- **Container runtime**: Docker SDK (`github.com/docker/docker/client`)

## Architecture

### Backend (`cmd/`, `internal/`)

1. **Entrypoint** — `cmd/codedockd/main.go`: HTTP server daemon, wires dependencies
2. **Models** — `internal/models/`: Domain structs, DTOs, database entities (no circular imports)
3. **Repositories** — `internal/repositories/`: Database persistence, SQL interfaces, SQLite implementations
4. **Services** — `internal/services/`: Business logic, external integrations
5. **Handlers** — `internal/handlers/`: HTTP controllers, Echo route handlers
6. **HTTP Setup** — `internal/http/`: Server setup, routes, CORS, auth middleware
7. **Engine** — `internal/engine/`: Container engine, Docker deployer, runtime management, cron, backup workers

### Dashboard (`dashboard/`)

React 19 + TanStack Router + TanStack Query + Zustand + Radix UI + Tailwind CSS v4. The self-hosted control panel where users deploy apps, manage databases, view logs, and configure settings.

## Code Style & Conventions

### Go & Architecture

- **Layered Monolith Architecture (`internal/`)**: All Go code inside `internal/` must be organized by clean functional layers: `models/`, `repositories/`, `services/`, `handlers/`, `http/`, and `engine/`. Avoid domain-driven vertical slices.
- **Consumer-Defined Interfaces**: Define narrow interfaces where they are _consumed_ (`Accept interfaces, return structs`).
- **Files**: `snake_case.go` — `container_health.go`
- **Packages**: short, lowercase, single word — `cron`, `auth`, `apikeys`
- **No inline comments** — only GoDoc when logic is non-obvious
- **No GoDoc** on self-explanatory types (`// User represents a user`) or HTTP handlers
- **No global state** — pass deps via struct fields, wire in `main.go`
- **Always check errors** — wrap with `fmt.Errorf("context: %w", err)`
- **JSON tags** on every exported struct field
- **No `init()`** — use explicit constructors

### TypeScript (Dashboard)

- **Files**: `kebab-case.tsx` — `project-card.tsx`
- **Named exports** over default exports
- **One component per file**
- Routes in `dashboard/src/routes/` (TanStack Router file conventions)
- Do **not** edit `routeTree.gen.ts` by hand
- **State Management:** Use standard Zustand (`create`) for global UI state. No wrappers, shortcuts, or legacy APIs.
- **Data Tables:** Use `@tanstack/react-table` for data grid components.
- **Telemetry:** Use `posthog-js` and `@posthog/react`. Integrations go in `dashboard/src/integrations/`.
- Use `tailwind-merge` + `clsx` + `class-variance-authority` for class composition

### Commands (`Makefile`)

| Command       | Action                                                                               |
| ------------- | ------------------------------------------------------------------------------------ |
| `make dev`    | Start Go daemon + Vite dashboard dev server concurrently (port 3000)                 |
| `make build`  | Build dashboard GUI to `dashboard/dist` and Go binary `bin/codedockd`                |
| `make format` | Format Go code (`go fmt ./...`) and dashboard (`cd dashboard && npm run format:fix`) |
| `make check`  | Run Go fmt + vet and Biome format check                                              |
| `make test`   | Run Go unit tests and dashboard tests                                                |

## Navigation Tips

- **Find a handler**: `internal/handlers/<domain>/<handler>.go`
- **Find a repository**: `internal/repositories/<entity>.go`
- **Find a service**: `internal/services/<domain>/<service>.go`
- **Find a dashboard route**: `dashboard/src/routes/` (TanStack Router file-based)
- **Find a dashboard component**: `dashboard/src/components/<domain>/`

## Constraints

- Do NOT run build or test commands after every change — the user will say if something breaks
- Do NOT commit unless explicitly asked
- Do NOT edit `routeTree.gen.ts` by hand
- Do NOT add `init()` in Go — use explicit constructors
- Do NOT use `mattn/go-sqlite3` — use `modernc.org/sqlite`
- Do NOT add unnecessary dependencies
- Run `make format` before finishing a session
