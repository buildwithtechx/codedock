# Hetzner live validation drill

Disposable end-to-end drill for Phase 2: provision a real Hetzner Cloud server through
Codedock, verify bootstrap, resize it, then tear it down without leftovers. Run from a
checkout with `curl` and `jq` installed.

## Environment contract

Export these before running. Use a disposable Hetzner Cloud project; the drill creates
and deletes a real server billed to that project (a few cents for a cx23).

- `CODEDOCK_URL` — daemon base URL, e.g. `http://localhost:8080/api`.
- `CODEDOCK_TOKEN` — login JWT for an organization admin (`server:read` + `server:write`).
- `CODEDOCK_ORG_ID` — organization that owns the drill credential and quota.
- `HETZNER_TOKEN` — Hetzner Cloud API token for the disposable project.
- `DRILL_REGION` — optional, default `fsn1`.
- `DRILL_TYPE` — optional, default `cx23`.
- `DRILL_RESIZE_TYPE` — optional, default `cx33` (must be an upgrade from `DRILL_TYPE`).
- `DRILL_IMAGE` — optional, default `ubuntu-24.04`.

The drill manages its own credential (label `live-drill-<timestamp>`) and deletes it at
the end. Ensure the organization quota allows one extra server before starting.

## Automated drill

```sh
./scripts/hetzner-live-drill.sh
```

The script walks the same API path as the dashboard wizard:

1. Creates a credential (fails fast when Hetzner rejects the token).
2. Fetches live offers and asserts the drill region and type are purchasable.
3. Reviews and applies a provision, then polls the operation to `COMPLETED`.
4. Refreshes the server and asserts it is online with a public IP.
5. Reviews and applies a resize, asserting the external ID is unchanged (same disk).
6. Reviews and applies a delete, then asserts the server and its SSH key are gone
   from Hetzner Cloud.
7. Deletes the drill credential.

Timings are printed per phase. Expected on a healthy path: provision completes in
under 10 minutes with bootstrap (first boot to Docker ready) well under that; resize
and delete each complete in under 5 minutes.

## Manual wizard path

Prefer the dashboard when validating UX instead of timings:
`Servers → Add server → Hetzner Cloud`, then follow credential → plan → review →
live terminal. The terminal streams the same operation the script polls.

## Sign-off checklist (roadmap P2_D / P2_E)

- [ ] Provision reaches `COMPLETED`; server is online with a public IP.
- [ ] Operation logs show `PROVISIONING → BOOTSTRAP → READY` with no SSH/Docker timeout.
- [ ] Resize to a larger type completes and keeps the same Hetzner server ID.
- [ ] Delete completes; no `codedock-server` labelled server or `codedock-m-*` SSH key
      remains in the Hetzner project.
- [ ] Drill credential removed from the organization.
