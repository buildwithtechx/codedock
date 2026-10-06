---
title: Autoscaling
description: Opt in to local Docker replica changes with CPU thresholds and recorded decisions.
---

Autoscaling is disabled unless you explicitly enable it in an application's overview. It currently supports local Docker targets. Remote SSH, cluster and managed cloud targets require a separate runtime capability.

Set minimum and maximum replicas between 1 and 10, with the maximum at least the minimum. Choose CPU thresholds satisfying `0 <= down < up <= 100` and a cooldown between 120 seconds and 24 hours. The primary replica supplies the CPU measurement; this is not an average across all replicas. The worker evaluates policies every two minutes.

A threshold crossing changes the desired replica count by one within the configured bounds. Counts above zero that fall outside the bounds reconcile after cooldown. Legacy counts below one mean one effective replica; the stored value changes only when scaling requests a deployment. Scaling creates a deployment, so build and rollout behavior follows the application's deployment configuration. An active deployment, changed policy or concurrently changed replica count prevents a scaling reservation.

The policy panel shows the last decision, CPU sample, evaluation time and scale request time. A scale request is not a completed deployment: follow its deployment status and logs. A deployment creation failure restores the prior desired count and records the failure. Turning autoscaling off prevents future reservations; it does not cancel an already accepted deployment.

Use `GET /api/services/:serviceId/autoscaling` to inspect the policy and `PUT` to save it. Saving requires project administrator access. Existing services have disabled default policies without changing their replica counts.
