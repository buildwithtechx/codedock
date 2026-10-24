# App catalogue

One catalogue backs the dashboard one-click list (`GET /api/one-click`), the database
provisioning options, and the marketing template index. Templates live in
`internal/engine/compose/*.yaml`, one file per app, validated at load by
`TestEmbeddedCatalogueValidates`. There is no registry upload step: adding a file to
that directory ships the template.

## Template shape

```yaml
x-codedock:
  is_one_click: true
  verified: true
  name: Example
  description: What the app does
  category: analytics
  icon: example
  secrets:
    - var: APP_SECRET
      label: App secret
      length: 32
  inputs:
    - var: BASE_URL
      label: Public URL
      default: http://localhost:8000
services:
  web:
    image: example/web:1.0
    ports:
      - "8080"
    environment:
      - APP_SECRET=${secret:APP_SECRET}
      - BASE_URL=${input:BASE_URL}
    healthcheck:
      test: ["CMD-SHELL", "wget -qO- http://localhost:8080/ || exit 1"]
      interval: 30s
      timeout: 5s
      retries: 3
    volumes:
      - web-data:/data
volumes:
  web-data:
```

`secrets` are generated at install unless the user provides a value. `inputs` are
human settings collected by the deploy dialog: each needs a `default` or
`required: true`, never both. `${secret:VAR}` and `${input:VAR}` placeholders must
reference declared entries, and every declared entry must be referenced. Ports use
`host:container` or bare container ports; volumes must be declared named volumes,
never bind mounts. `healthcheck` follows Compose syntax and is rendered verbatim
into the install plan. Database templates set `is_database: true` with a
`connection_string` instead and keep the database provisioning path.

## Verified bar

`verified: true` means a maintainer installed the template on a clean instance from
this catalogue and confirmed the app boots and serves traffic. Flip the flag only
after that install, and re-verify after image bumps that change major versions.
Unverified templates install normally; the dashboard just omits the badge.

## Exclusion criteria

A service is a bad one-click fit when a correct install cannot be expressed in one
self-contained Compose file. Concretely, Codedock does not ship templates that need:

- bind-mounted config files or directories from the host,
- Compose `extends` chains or multi-file `include` layouts,
- more than a handful of cooperating services (roughly: app plus datastores),
- host-level setup steps outside the Compose document.

PostHog is excluded on all four counts: its self-hosted distribution is a base file
plus an environment overlay that `extends` shared service definitions, seven-plus
bind-mounted ClickHouse configuration files, and around twenty-five services
including Kafka, Zookeeper, Temporal, SeaweedFS, and supporting workers. Shrinking
that to one file would produce an install that looks complete but silently drops
pipelines. Use the documented Docker deployment instead.

Full Supabase is excluded for the same reason: the upstream self-hosted stack spans
Postgres plus roughly a dozen services (Auth, Storage, Realtime, Kong, Studio,
PostgREST, pg-meta, Edge Runtime) with shared secrets and storage volumes that must
be provisioned in order. The catalogue ships `supabase-studio` only, which pairs
Studio and Kong with a single Postgres for UI-first exploration.

Redis-protocol alternatives (Dragonfly, KeyDB) stay database-only templates: they
have no web UI to brand or one-click into, so they install through database
provisioning and intentionally carry no logo tile.
