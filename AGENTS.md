# Agent Instructions — Codedock Website & Documentation

## Code Style

- **Max 350 lines per file.** If a file exceeds this, extract components, helpers, or types into separate files.
- **One component per file.** Never cram thousands of lines into a single file. Break components down into individual files.
- **Name files in `kebab-case`** (e.g. `feature-card.astro`, `stack-carousel.astro`).
- Use named exports over default exports where applicable.
- **No comments or JSDoc allowed.** Code should be self-explanatory. If logic is non-obvious, refactor to make it clear rather than adding a comment.

## Workflow

- **Do not run build or test commands after every change.** Just make the code change. If something breaks, the user will say so.
- Run `npm run format:fix` (`npx @biomejs/biome check --write .`) before committing or finishing a session. NEVER run `prettier` (`npx prettier`) — Biome is our strict formatter.
- Prefer `read`/`grep`/`glob` tools over `bash` for file exploration.
- When making edits, read the file first, then use `edit` for targeted changes.

## Stack

| Layer           | Tech                                       |
| --------------- | ------------------------------------------ |
| Marketing (web) | Astro 7, Tailwind CSS v4                   |
| Docs            | Astro 7, Starlight                         |
| Styling         | Tailwind CSS v4 `@theme` directives        |
| Monorepo        | npm workspaces (`apps/web/`, `apps/docs/`) |
| Formatter       | Biome (`@biomejs/biome`)                   |

## Conventions

- **Marketing pages** live in `apps/web/src/pages/`, components in `apps/web/src/components/`.
- **Pure Astro + Tailwind CSS v4** — NO React/Vue/Svelte islands on marketing pages.
- **Docs pages** live in `apps/docs/src/content/docs/` as `.md` or `.mdx` files following Starlight file-based routing.
- **Sidebar config** lives in `apps/docs/astro.config.mjs`. All sidebar groups MUST have `collapsed: false` so categories stay permanently open.
- Use Tailwind CSS v4 `@theme` directives for design tokens; avoid custom CSS where Tailwind utilities suffice.
- **Format strictly with Biome** (`npm run format:fix` / `biome check --write .`). NEVER use Prettier (`npx prettier`).

## Commands

```sh
npm run dev:web        # starts marketing site at http://localhost:4321
npm run dev:docs       # starts documentation portal at http://localhost:4322
npm run build:all      # builds both apps/web and apps/docs
npm run format:fix     # runs Biome formatter on all workspace files
```
