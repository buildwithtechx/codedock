# Codedock Security & Architectural Roadmap

This roadmap tracks all active security remediations, architectural refactoring, and code-style compliance tasks across the Codedock codebase (`internal/`, `cmd/`, `pkg/`, `bootstrap/`, `scripts/`, and `apps/dashboard/src/`). All historical and completed tasks have been pruned.

---

## 1. Critical Security Findings

- [x] **1. Service log WebSocket authorization bypass:** Add project/service membership verification in `internal/handlers/service_logs_ws.go:42` so arbitrary JWTs cannot connect to logs of unrelated services.
- [x] **2. Unauthorised cross-project deployments:** Enforce project membership checks in Compose (`internal/handlers/compose.go:67`), Archive (`internal/handlers/archive.go:25`), and One-click (`internal/handlers/oneclick.go:30`) deployments.
- [x] **3. OAuth client secrets exposed:** Make `GET /settings/oauth/providers` admin-only and always redact `clientSecret` (`internal/handlers/oauth.go:30`, `internal/models/auth.go:21`).
- [x] **4. OAuth login CSRF/account-linking vulnerability:** Validate returned OAuth `state` values against session/cookie state in `internal/handlers/oauth.go:60, 81`.
- [x] **5. Canvas and audit-log tenant leakage:** Restrict global canvas summaries (`internal/handlers/canvas.go:21`) and audit logs (`internal/handlers/audit_log.go:21`) by project, organization, or admin role.

---

## 2. High-Risk Findings

- [x] **6. App services can be reassigned across projects:** Validate that `EnvironmentID` belongs to `ProjectID` on app service create/update (`internal/handlers/app_services.go:59, 153`).
- [x] **7. Rollback authorization is weaker than deployment authorization:** Enforce consistent admin-level service access for rollback (`internal/handlers/deployment.go:130`) matching deployment triggering (`internal/http/routes.go:232`).
- [x] **8. Database passwords returned in API models:** Prevent returning decrypted password fields in `json:"password"` (`internal/models/database.go:35`, `internal/repositories/database.go:56`).
- [x] **9. 2FA is not enforced during login:** Implement interactive TOTP second-factor challenge during login when `TOTPEnabled` is true (`internal/services/auth_service.go:126`).
- [x] **10. JWT sessions survive user disablement/deletion:** Verify account status, password-version revocation, or disablement during JWT token validation/authorization.
- [x] **11. Refresh tokens are replayable:** Implement server-side storage and revocation tracking for refresh tokens (`internal/services/auth_service.go:156`).
- [x] **12. Unbounded archive/upload reads:** Enforce maximum size limits on Compose and Archive upload readers (`internal/handlers/compose.go:98`, `internal/handlers/archive.go:60`).

---

## 3. Bootstrap & Operational Issues

- [ ] **13. Installer checksum verification is optional:** Require checksum validation for Compose, control binaries, Docker installer, and worker binaries (`bootstrap/install.sh:12`, `bootstrap/worker.sh:117`).
- [ ] **14. Generated .env permissions are not explicitly restricted:** Enforce `0600` permissions when writing `.env` secrets (`bootstrap/install.sh:152`).
- [ ] **15. Docker socket grants control-plane-level host access:** Tighten container permissions or document security implications of Docker socket access in `docker-compose.prod.yml:18, 64`.
- [ ] **16. Production healthcheck appears invalid:** Ensure backend route `/healthz` exists and matches `docker-compose.prod.yml:37` healthcheck definition.
- [ ] **17. Restore script should validate archive entry types:** Reject symlinks/hardlinks in archive extraction in addition to path traversal (`scripts/restore.sh:22`).

---

## 4. Medium-Risk & Dashboard Security Findings

- [ ] **AI Diagnosis Rate Limiting:** Add per-user rate or cost limits for AI deployment failure diagnosis.
- [ ] **Database Query Endpoints:** Require strict authorization and add audit logging/safeguards for destructive SQL/Redis query endpoints.
- [ ] **WebSocket Origin Protection:** Reject empty `Origin` headers where appropriate to strengthen CSRF protection.
- [ ] **Strict JWT HMAC Algorithm:** Explicitly restrict token validation to `HS256` only.
- [ ] **Fail-Safe Handler Lookups:** Ensure handlers fail closed rather than open when database errors occur during lookup before destructive operations.
- [ ] **Backup Record Permissions:** Enforce strict admin authorization and secure filesystem permissions for backup records containing sensitive data.
- [ ] **WebSocket Authentication:** Send explicit authentication token/header for live logs WebSocket rather than relying solely on cookies.
- [ ] **Backend Authorization Source of Truth:** Ensure dashboard UX guards are backed by rigorous backend authorization checks.
- [ ] **Frontend Credential Redaction:** Ensure sensitive credential fields in frontend models remain redacted by default.

---

## 5. Architectural & Structural Cleanup

### Go Backend

- [ ] **Broad Handlers:** `internal/handlers/` is becoming too broad. Split by domain (e.g., `auth/`, `users/`, `projects/`, `deployments/`, `backups/`).
- [ ] **Bloated Services:** `internal/services/` has too many unrelated responsibilities. Group related services similarly to handlers.
- [ ] **Mixed Engine Package:** `internal/engine/` is a large mixed package. Separate infrastructure concerns (e.g., `docker/`, `cron/`, `networking/`, `observability/`).
- [ ] **CLI Mix:** `cmd/codedockd/` mixes CLI parsing, setup, deployment, DB ops, and password management. Move command implementations into `cmd/codedockd/commands/`.
- [ ] **Script Overlap:** `bootstrap/` and `scripts/` overlap in responsibilities (upgrade, restore, install, backup). Define one canonical lifecycle location.
- [ ] **Public SDK Leakage:** `pkg/http/` imports internal models, making it less reusable as a public SDK.

### Dashboard Structure

- [ ] **Colocate Domain Concepts:** `services/`, `hooks/`, and `interfaces/` duplicate domain concepts. Colocate them within `features/<domain>/` (e.g., `features/projects/api.ts`, `hooks.ts`, `types.ts`).
- [ ] **Feature Granularity:** `features/instance/` is too broad. Split into `features/settings/`, `backups/`, `dns/`, `notifications/`, `users/`, `registries/`.

---

## 6. Code Style: 350-Line Limit Compliance

The following files exceed the project's 350-line limit and require refactoring:

- [ ] `apps/dashboard/src/features/instance/notifications-settings.tsx` — 436 lines
- [ ] `internal/handlers/backup.go` — 403 lines
- [ ] `apps/dashboard/src/features/instance/s3-destinations-list.tsx` — 398 lines
- [ ] `apps/dashboard/src/features/instance/maintenance-settings.tsx` — 387 lines
- [ ] `internal/engine/backup_manager.go` — 383 lines
- [ ] `apps/dashboard/src/features/instance/backups-list.tsx` — 352 lines
