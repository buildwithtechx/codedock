# Codedock — Agent Instructions

## What This Is

This is the Codedock codebase — an open-source, self-hosted PaaS that turns any bare-metal VPS into a private Vercel & Railway. It's written in Go (backend daemon) and TypeScript (React 19 control panel dashboard in `apps/dashboard/`, Tauri 2.0 desktop app in `apps/desktop/`).

- **Language (backend)**: Go — `cmd/`, `internal/`, `pkg/`
- **Language (frontend)**: TypeScript/TSX — React 19 in `apps/dashboard/`
- **Desktop Application**: Tauri 2.0 (Rust native shell wrapping `apps/dashboard/dist`) in `apps/desktop/`
- **Dashboard Stack**: React 19 + TanStack Router + Vite + Tailwind CSS v4
- **Database**: embedded SQLite (`modernc.org/sqlite`, CGO-free)
- **Container runtime**: Docker SDK (`github.com/docker/docker/client`)

## Architecture

### Backend (`cmd/`, `internal/`)

1. **Entrypoint** — `cmd/codedockd/main.go`: HTTP server daemon, wires dependencies
2. **Domain Modules** — `internal/modules/`: Self-contained vertical slices (`apps`, `auth`, `clusters`, `databases`, `deployments`, `git`, `monitoring`, `projects`, `servers`, `system`, `users`, etc.) containing models, repositories, services, handlers, and routes.
3. **Platform Infrastructure** — `internal/platform/`: Infrastructure drivers and runtime execution utilities (`bare`, `build`, `database`, `docker`, `networking`, `notifications`, `observability`, `scheduler`, `ssh`, `storage`, `vault`).
4. **Deploy Execution Engine** — `internal/deploy/`: Container deployment runtime, rollout state machine, cutover decider, autoscaler.
5. **Adapters** — `internal/adapters/`: Glue adapters between modules and runtime engines (modules never import `adapters`).
6. **HTTP Transport** — `internal/http/`: Server setup, route registration, Echo middleware, and presentation.

### Apps Monorepo (`apps/`)

- **Dashboard (`apps/dashboard/`)**: React 19 + TanStack Router + TanStack Query + Zustand + Radix UI + Tailwind CSS v4. The self-hosted control panel where users deploy apps, manage databases, view logs, and configure settings.
- **Desktop (`apps/desktop/`)**: Tauri 2.0 cross-platform application shell (macOS, Windows, Linux).

## Code Style & Conventions

### Go & Architecture

- **Domain-Oriented Modular Monolith (`internal/`)**: All Go code inside `internal/` is organized by self-contained domain modules (`internal/modules/`), platform drivers (`internal/platform/`), deploy runtime (`internal/deploy/`), adapters (`internal/adapters/`), and thin HTTP router (`internal/http/`).
- **Strict File Size Cap**: Every `.go` file **MUST be <= 300 LOC**.
- **Zero Legacy Imports**: Zero imports of legacy directories (`internal/models`, `internal/repositories`, `internal/services`, `internal/handlers`, `internal/engine`).
- **Consumer-Defined Interfaces**: Define narrow interfaces where they are _consumed_ (`Accept interfaces, return structs`).
- **Files**: `snake_case.go` — `container_health.go`
- **Packages**: short, lowercase, single word — `cron`, `auth`, `apikeys`
- **No inline comments** — only GoDoc when logic is non-obvious
- **No GoDoc** on self-explanatory types (`// User represents a user`) or HTTP handlers
- **No global state** — pass deps via struct fields, wire in `main.go`
- **Always check errors** — wrap with `fmt.Errorf("context: %w", err)`
- **JSON tags** on every exported struct field
- **No `init()`** — use explicit constructors

### TypeScript (Dashboard & Desktop)

- **Files**: `kebab-case.tsx` — `project-card.tsx`
- **Named exports** over default exports
- **One component per file**
- Routes in `apps/dashboard/src/routes/` (TanStack Router file conventions)
- Do **not** edit `routeTree.gen.ts` by hand
- **State Management:** Use standard Zustand (`create`) for global UI state. No wrappers, shortcuts, or legacy APIs.
- **Data Tables:** Use `@tanstack/react-table` for data grid components.
- **Telemetry:** Use `posthog-js` and `@posthog/react`. Integrations go in `apps/dashboard/src/integrations/`.
- Use `tailwind-merge` + `clsx` + `class-variance-authority` for class composition

### Commands (`npm run ...`)

| Command       | Action                                                                               |
| ------------- | ------------------------------------------------------------------------------------ |
| `npm run dev`    | Start Go daemon + Vite dashboard dev server concurrently (port 3000)                 |
| `npm run build`  | Build dashboard GUI to `apps/dashboard/dist` and Go binary `bin/codedockd`           |
| `npm run fmt`    | Format Go code (`go fmt ./...`) and dashboard (`npm run format:fix`)                 |
| `npm run check`  | Run Go fmt + vet and Biome format check                                              |
| `npm run test`   | Run Go unit tests and dashboard tests                                                |

## Navigation Tips

- **Find a domain module**: `internal/modules/<domain>/` (handlers, models, repos, services, routes)
- **Find a platform driver**: `internal/platform/<capability>/`
- **Find deployment runtime**: `internal/deploy/`
- **Find a dashboard route**: `apps/dashboard/src/routes/` (TanStack Router file-based)
- **Find a dashboard component**: `apps/dashboard/src/components/<domain>/`

## Constraints

- Do NOT run build or test commands after every change — the user will say if something breaks
- Do NOT commit unless explicitly asked
- Do NOT edit `routeTree.gen.ts` by hand
- Do NOT add `init()` in Go — use explicit constructors
- Do NOT use `mattn/go-sqlite3` — use `modernc.org/sqlite`
- Do NOT add unnecessary dependencies
- Run `npm run fmt` before finishing a session
