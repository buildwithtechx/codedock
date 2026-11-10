// Package deploy contains the workload execution engine, container manager, and zero-downtime rollout runtime.
//
// # Architectural Purpose:
//
// This package is an execution engine consumed by domain modules (such as deployments and apps)
// to coordinate image builds, container lifecycle transitions, and blue-green or rolling updates.
//
// # Key Responsibilities:
//   - deployer.go: Core pipeline coordinating build outputs, volume mounting, and container creation.
//   - rollout.go & rollout_recovery.go: Zero-downtime deployment execution and automatic rollback on failure.
//   - container_manager.go: Direct container lifecycle operations (start, stop, restart, destroy).
//   - cutover_decision.go: Traffic cutover validation checking health checks and port responsiveness.
//   - database.go: Provisioning, lifecycle management, and migration execution for database containers.
//   - interfaces.go: Dependency storage contracts required by the deployment runtime.
package deploy
