# Codedock Website & Documentation Monorepo 🌐

Official repository for the **Codedock Marketing Website** (`codedock.run`) and **Documentation Portal** (`docs.codedock.run`).

Built with **Astro 7**, **Starlight**, and **Tailwind CSS v4**.

---

## 📁 Repository Structure

```text
codedock-website/
├── apps/
│   ├── web/          # Marketing landing site, comparison suite, solutions, and feature deep-dives
│   └── docs/         # Full Starlight documentation portal with uncollapsed navigation hierarchy
├── package.json      # npm workspace configuration
└── biome.json        # Biome code formatting & linting configuration
```

---

## 🚀 Development Quickstart

### Prerequisites

- **Node.js**: `^20.0.0` or `^22.0.0`
- **npm**: `^10.0.0`

### 1. Install Dependencies

```bash
npm install
```

### 2. Run Local Dev Servers

#### Marketing Website (`apps/web`)

Runs at `http://localhost:4321`:

```bash
npm run dev:web
```

#### Documentation Portal (`apps/docs`)

Runs at `http://localhost:4322`:

```bash
npm run dev:docs
```

---

## 🛠️ Build & Deployment Commands

```bash
# Build both marketing website and docs
npm run build:all

# Build marketing site only
npm run build:web

# Build documentation portal only
npm run build:docs

# Format and check code with Biome
npm run format:fix
```

---

## 📖 Related Repositories

- **[codedock](https://github.com/buildwithtechx/codedock)**: Main engine repository containing the Go daemon (`codedockd`) and React control panel dashboard (`dashboard/`).

---

## 📄 License

This repository is licensed under the [MIT License](LICENSE).
