# Dashboard UX Audit

Reference implementation inspected: `/home/eminisolomon/Dev/openship/apps/dashboard`.

The reference informs layout hierarchy, page composition, navigation density, state handling, and responsive behavior. Codedock keeps its own identity, content, and visual assets.

## Global Rules

- Keep the global sidebar to nine destinations: Home, Projects, Apps, Deployments, Servers, Domains & DNS, Backups, API Access, and Settings.
- Place workspace switching and account actions at the bottom of the sidebar.
- Keep the primary sidebar focused on destinations, not operational subpages.
- Use a contextual side rail only for a detail page or when it provides current-page decisions, metrics, or navigation.
- Make contextual rails sticky on desktop and replace them with a horizontal or menu-based control on mobile.
- Use bordered surfaces for lists, tables, forms, and summaries. Do not wrap first-use, loading, empty, or recoverable error states in a generic card.
- Keep one page title and one page description. Tab contents start with their task content unless the tab needs a distinct object-level heading.
- Use page routes for creation and multi-step workflows. Use dialogs only for focused, reversible actions.

## Route Map

### Home `/`

Current: `home-overview`, `home-project-list`, `home-runtime-summary`, and `home-shortcuts`.

Target:

- Keep the workspace greeting and sticky runtime rail.
- Use a full-width, unframed first-project state when no projects exist.
- Show the project list as the primary data surface only after projects exist.
- Reduce shortcut cards to a secondary action strip. Do not let them compete with the first-project action.

### Projects `/projects`

Current: `project-directory` with a grid of project cards.

Target:

- Keep a concise heading, count, and one primary creation action.
- Keep an unframed first-project state.
- Add search, sort, and a target filter only when project volume requires them.
- Use project cards for actual projects, not as a wrapper around the empty state.

### Project Detail `/projects/:projectId/*`

Current: a project context rail with overview, canvas, compose, settings, members, tokens, and scheduled-task routes.

Target:

- Retain the contextual rail and mobile context navigation.
- Normalize the overview into a compact project identity header, environment selector, resource list, and empty-resource state.
- Keep project-level resources, environments, tasks, and members in this context instead of global navigation.
- Split the current oversized overview component before visual work; it mixes header, empty state, resource cards, and status rendering.

### Apps `/apps`

Current: `app-directory` lists deployed application services.

Target:

- Keep Apps as the organization-wide application inventory.
- Keep an unframed first-app state.
- Add status and project filtering only when there are applications to filter.
- Render app cards as concise inventory rows or tiles with status, project, endpoint, and latest deployment information.

### New App `/apps/new` and `/projects/:projectId/new`

Current: project picker followed by project-scoped resource creation.

Target:

- Treat this as Codedock's deployment library, not a modal.
- Provide source choices first: Git, Docker image, Compose, database, and one-click application.
- Keep a sticky contextual rail for source connection and a concise deployment summary when it improves the decision.
- Reuse the project picker only when the entry point does not already establish a project.

### Service Detail `/services/:serviceId/*`

Current: a contextual rail with fifteen service routes.

Target:

- Retain the contextual detail layout.
- Group tabs into Setup, Observe, Operate, and Danger without showing every route at once on narrow screens.
- Keep global Deployments separate from service-scoped deployment history.
- Use service pages for configuration and observability, not for organization navigation.

### Deployments `/deployments`

Current: a project selector, app selector, then a service-scoped deployment table.

Target:

- Replace selector-first browsing with an organization-wide deployment feed.
- Add search, project filter, and status filter above the feed.
- Show a sticky overview rail with total, successful, active, and failed counts.
- Use an unframed first-deployment state and a distinct filtered-empty state.

Backend requirement:

- Add an authorization-scoped, paginated deployment list by organization, with project, service, status, and search filters.
- Return deployment rows with service name, project name, status, trigger, branch, commit, and timestamps so the dashboard does not fan out per row.

### Servers `/servers`, `/servers/new`, `/servers/:serverId`

Current: server list uses legacy card, empty, and error styling.

Target:

- Replace the legacy dashed empty card and error card with flat states.
- Keep server cards only for actual fleet entries.
- Add a sticky fleet summary on the list when servers exist.
- Preserve the server detail page as a contextual layout for connection, metrics, components, and terminal operations.

### Domains & DNS `/dns`

Current: a three-tab audit, provider-credentials, and global-domain-settings page.

Target:

- Keep this Codedock-specific resource page and its three tabs.
- Preserve drafts when changing tabs.
- Replace generic error panels with the shared flat recovery state.
- Keep statistics as compact metrics, not as a dashboard-wide card grid.

### Backups `/backups` and `/s3-destinations`

Current: the backup page is a global configuration form while destinations are a separate route.

Target:

- Keep Backups independent in the sidebar.
- Merge storage destinations into the Backups page as a tab or primary list section.
- Reserve the backup configuration form for a focused settings surface or database context.
- Keep restore history and destination details as data surfaces; use a page or large dialog for destination setup.
- Redirect the legacy `/s3-destinations` route after the merged page exists.

### API Access `/api-access`

Current: standalone personal access token page.

Target:

- Keep it as an independent sidebar destination.
- Do not duplicate personal access tokens in profile or general Settings.

### Settings `/settings?tab=*`

Current: ten right-rail tabs: General, Notifications, OAuth, Team, Audit, AI, Sources, Maintenance, Updates, and Migration.

Target:

- Keep one Settings page heading and a sticky right-side tab rail.
- Keep General, Notifications, OAuth, Team, Audit, AI, Sources, Maintenance, Updates, and Migration as URL tabs.
- Remove repeated title and subtitle blocks from tab content unless a tab is an independent workflow.
- Give each settings section its own save or action control rather than a global page save action.
- Keep audit logs, team management, source integrations, migration, maintenance, and updates out of the global sidebar.

### Legacy and Context Routes

- `/audit-logs` redirects to Settings Audit.
- `/scheduled-tasks` remains a project or service operation; do not place it in global navigation.
- `/profile` remains account-only and must not duplicate API Access.
- `/organizations` and `/users` remain owner/admin management routes until their capabilities are folded into Settings Team. Do not expose them as global sidebar items.

## Shared State System

- Loading states use local skeletons for data surfaces and an unframed spinner for full-page states.
- Empty states use a centered Codedock system-map visual, title, explanatory copy, and at most two actions.
- Recoverable errors use the same centered composition with an error tone and retry action.
- Filtered-empty states remain inside the data surface because filters are still actionable context.
- 403, 404, and failed API calls require distinct copy and available recovery actions.

## Implementation Order

1. Add organization-scoped deployment listing and replace the selector-first global deployment page.
2. Normalize Home, Projects, Apps, Deployments, Servers, and Backups page-state patterns.
3. Merge backup destinations into Backups and retire the standalone destination route.
4. Split project overview and service-detail shells into focused components, then align their contextual rails.
5. Remove duplicated Settings tab headings and move save actions to their owning sections.
6. Validate desktop, tablet, and mobile states for every route before additional visual changes.
