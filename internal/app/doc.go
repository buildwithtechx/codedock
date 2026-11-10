// Package app serves as the central composition root and application lifecycle orchestrator for Codedock.
//
// # Architectural Purpose:
//
// Decoupled from transport mechanics, this package wires together domain modules,
// platform drivers, cross-module adapters, and background workers into a coherent application.
//
// # Responsibilities:
//   - app.go: Bootstraps the application, opens the database, and drives lifecycle startup/shutdown.
//   - services.go: Instantiates all domain module repositories and business services.
//   - bindings.go: Wires cross-module adapters and runtime engines into domain services.
//   - listeners.go: Registers event bus dispatchers, notification workers, and metrics collectors.
//   - recovery.go: Restores in-flight containers, SSL certificates, and scheduled tasks on daemon start.
//   - routes.go & routes_*.go: Mounts domain routes onto the HTTP router without polluting transport.
package app
