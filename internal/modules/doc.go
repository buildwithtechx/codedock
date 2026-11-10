// Package modules contains all domain-specific business capabilities for Codedock.
//
// # Architecture Standards:
//
// Each subpackage under modules represents a bounded domain context (e.g. apps, auth,
// deployments, databases, projects, servers, monitoring).
//
// Standard File Structure per Module:
//   - models.go: Domain data structures, database entity models, and core status types.
//   - interface.go: Repository and Service contracts defining persistence and business capabilities.
//   - errors.go: Sentinel error definitions specific to the domain.
//   - request.go: HTTP incoming JSON request DTOs and parameter validations.
//   - response.go: HTTP outgoing JSON response DTOs.
//   - handlers.go: HTTP Echo controller handlers mapping incoming requests to service calls.
//   - service.go: Pure business logic coordinating persistence, platform events, and validations.
//   - repository.go: SQL database queries and data store operations.
//   - routes.go: Echo router endpoint registrations via RegisterRoutes.
//
// # Isolation Rules:
//   - Modules must remain self-contained and loosely coupled.
//   - Modules NEVER import internal/adapters or internal/app.
//   - Cross-module dependencies are fulfilled through interfaces injected at startup by internal/app.
package modules
