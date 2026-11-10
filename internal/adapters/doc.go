// Package adapters provides cross-module bridge implementations and structural glue.
//
// # Architectural Purpose:
//
// To prevent cyclic package imports and tight coupling between domain modules,
// modules declare the interfaces they require (e.g. apps.AccessChecker, canvas.ServerLookup).
//
// This package implements those interfaces by translating and delegating calls
// between the consumer module and the provider module (or execution engine).
//
// # Rules:
//   - Domain modules (internal/modules/*) NEVER import internal/adapters.
//   - Execution engines (internal/deploy) NEVER import internal/adapters.
//   - Adapters are instantiated only within the application composition root (internal/app)
//     during dependency injection at server startup.
package adapters
