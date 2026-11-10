// Package http provides the pure HTTP transport and routing infrastructure for Codedock.
//
// # Architectural Standards:
//
// This package is strictly concerned with HTTP transport mechanics and contains
// zero domain models, database queries, or long-running worker loops.
//
// # Responsibilities:
//   - router.go: Configures the Echo instance, global middleware pipeline, and group routes.
//   - echo_middleware.go: Request logging, panic recovery, gzip compression, and client IP resolution.
//   - csrf.go: CSRF cookie generation and header validation for state-changing endpoints.
//   - spa.go: Single-Page Application (SPA) HTML routing fallback for the web dashboard.
//   - docs.go: Embedded documentation serving.
//   - middleware/: JWT bearer authentication, session guards, token authentication, and rate limiting.
//   - apischema/: OpenAPI 3.0 specification generation and interactive documentation routes.
package http
