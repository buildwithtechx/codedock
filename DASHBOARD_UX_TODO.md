# Dashboard UX Refactor TODO

## Purpose

Rebuild the Codedock dashboard information architecture using OpenShip as a structural reference, not as a visual or copywriting template. Preserve every Codedock capability and make its location discoverable before removing or moving any navigation item.

## Reference Findings

OpenShip uses a small persistent navigation surface and distributes deeper functionality through page ownership.

- The primary sidebar contains Home, Projects, Apps, Deployments, Backups, Settings, and optional product-specific items.
- Theme and sidebar collapse controls live in the sidebar header.
- The create-project action and organization/account control live at the bottom of the sidebar.
- Servers, domains, sources, activity, and operational tools are not all primary navigation items.
- List pages use a main content column and an optional 340px contextual right rail.
- Every desktop contextual rail is sticky at the page top while its related main content scrolls.
- Detail pages use contextual navigation in a right rail on desktop and horizontal navigation on smaller screens.
- Settings use one content heading and a contextual settings rail. They do not repeat the page title, active tab title, and section title at the same hierarchy level.
- The Home page has a dominant projects surface, a narrow operational context column, and separate shortcut cards.
- Cards are independent surfaces. Page content is not wrapped inside a generic outer card.
- Full creation and configuration workflows are routes. Dialogs are reserved for bounded choices, confirmations, and lightweight edits.

## Codedock Constraints

- Keep Codedock branding, terminology, violet accent, cloud mark, and Manrope typography.
- Do not copy OpenShip copy, logos, illustrations, or gradients.
- Do not remove a route until its replacement route and discovery path are implemented.
- Do not make a feature reachable only through an unimplemented command-palette item.
- Keep desktop right rails keyboard accessible and provide an equivalent mobile navigation pattern.
- Keep one page title per screen. Sections may have headings, but they must not duplicate the page or active-tab title.

## Route Discovery Inventory

| Codedock capability | Current route | Intended discovery path | Refactor outcome |
| --- | --- | --- | --- |
| Home | `/` | Primary sidebar | Keep primary navigation |
| Projects | `/projects` | Primary sidebar | Keep primary navigation |
| Project overview, canvas, compose, resources, schedules, settings | `/projects/$projectId/*` | Project list, Home, project contextual rail | Keep nested and add complete right/mobile contextual navigation |
| Apps | `/apps` | Primary sidebar | Keep primary navigation |
| App and database service settings | `/services/$serviceId/*` | Project canvas/list, Apps, service contextual rail | Keep nested and add complete right/mobile contextual navigation |
| Deployments | `/deployments` | Primary sidebar | Keep primary navigation |
| Backups | `/backups` | Primary sidebar | Standalone page; link S3 destinations from its own page |
| S3 destinations | `/s3-destinations` | Backups page and backup settings | Remove from primary sidebar, never leave orphaned |
| DNS audit and provider settings | `/dns` | Project/service domains and a Home shortcut | Keep a global audit page; use Domains and DNS terminology consistently |
| Source providers | `/sources` | Home shortcut, project create/import flow, project source settings | Keep global page but not primary sidebar |
| Servers | `/servers`, `/servers/$serverId` | Home shortcut, project deploy target actions, server-required empty states | Keep global pages but not primary sidebar |
| Organizations | `/organizations`, `/organizations/$organizationId` | Organization switcher and settings/team links | Keep out of primary sidebar |
| Members | `/users` | Organization detail and settings/team section | Keep out of primary sidebar |
| API access | `/api-access` | Primary sidebar | Keep primary navigation because it is an explicit operator workflow |
| Profile and security | `/profile` | Account identity area | Keep direct route; do not create a separate user-settings dropdown |
| Audit logs | `/audit-logs` | Settings security section | Keep route and link from Settings |
| Maintenance | `/maintenance` | Settings instance section | Keep route and link from Settings |
| Updates | `/updates` | Settings instance section | Keep route and link from Settings |
| Migration | `/migrations` | Settings instance section | Keep route and link from Settings |
| Scheduled tasks | `/scheduled-tasks` | Project/service schedules and Settings operations | Confirm global purpose before retaining standalone route |
| AI configuration | `/ai` | Settings integrations or service-level AI tools | Confirm ownership and add a real entry point |
| New project | Current `CreateProjectModal` | Sidebar create action, Home, Projects empty state | Replace with a dedicated route that supports source, template, runtime, and deployment choices |
| New server | Current create-server dialog | Server-required calls to action and operations settings | Replace with a dedicated setup route when the flow includes credential, install, or verification steps |
| New backup destination | Current destination dialog | Backups page | Keep as a dialog only if the choice and credentials fit a short bounded flow |

## Primary Sidebar Target

### Header

- Codedock logo and product name.
- Theme toggle and sidebar collapse control aligned at the right.
- One visible divider beneath the header.

### Primary navigation

1. Home
2. Projects
3. Apps
4. Deployments

### Platform navigation

1. Backups
2. API Access
3. Settings

### Footer

- New project action.
- One workspace switcher trigger.
- Organization dropdown contains switch organization, manage organizations, signed-in identity, and sign out.
- Do not render a second permanent user row or a generic user-settings popup.

### Explicitly excluded from the primary sidebar

- Servers
- Domains and DNS
- Sources
- AI
- S3 destinations
- Organizations
- Members
- Maintenance
- Updates
- Migration
- Audit logs
- Documentation

Each excluded capability must have at least one implemented route-level discovery path listed in the inventory above.

## Page Architecture Target

### Shared page frame

- Introduce one dashboard page-frame primitive with max content width, responsive padding, standard header spacing, and a 340px desktop rail option.
- Do not use this primitive to create a card around page content.
- Require `lg:sticky lg:top-6 lg:self-start` for every right rail unless the rail contains an independently scrolling table or editor.
- Standardize loading, empty, unavailable, and permission-denied states.

### Reference page patterns to adapt

| Reference pattern | Structural rule for Codedock |
| --- | --- |
| Settings | One `Settings` page header and subtitle. Main content uses domain-specific section titles. The sticky rail has a static settings/account identity card plus navigation. Do not render the selected tab again as a second main heading. |
| Backups | A standalone page with header, supported-storage empty state, and a focused destination-picker dialog. The page owns backup browsing; the dialog only chooses a destination. |
| Deployments | Header, filters, dominant list or empty state, and a sticky operational summary rail with a real next action. |
| Projects | Header and create action. Empty state is a full page experience that launches the project-creation route rather than a modal. |
| New project/library | A dedicated route with source-mode navigation, a dominant setup panel, and a sticky connection/summary rail. |
| Detail screens | Main working area plus a sticky contextual rail. On mobile, the same destinations move into a horizontal nav or sheet above the main content. |

### Home

- Keep one compact greeting and subtitle.
- Use a two-column desktop grid: dominant projects content plus a 340px operational rail.
- Keep project creation in the dominant empty state and sidebar action.
- Surface sources, servers, domains, and platform configuration as independent shortcut cards or contextual calls to action.
- Make every Home shortcut point to a real, supported route.

### Projects

- List page: heading, create action, search/filter controls, project cards/list, optional deploy-target rail.
- Detail page: project identity/context card plus section navigation on the right; mobile horizontal section navigation.
- Keep canvas, compose, resource creation, schedules, and settings discoverable from the project rail.
- Do not duplicate project title in the parent layout and child screen.
- Make project creation a route, not a modal. The route must support source import, template selection, runtime choice, and an explicit return path.

### Apps

- Treat Apps as deployed workloads across the active organization.
- Keep the app list focused on status, deployment target, source, and direct service navigation.
- Provide a contextual rail for templates, suggested deployments, or operational summary only when it has real data or a next action.
- Provide a dedicated app creation route for catalog/template selection and preserve progress in the URL or draft state.

### Services

- Preserve all existing service routes: configuration, build, deployments, metrics, webhooks, schedules, storage, domains, route rules, variables, terminal, serverless editor, previews, log drains, and danger zone.
- Group sections in the service rail by lifecycle rather than exposing a flat list of fifteen links.
- Keep the right rail sticky on desktop and the same grouped controls horizontally scrollable or sheet-based on mobile.
- Ensure project context and a back path are always visible.

### Deployments

- Keep filters and deployments list in the main column.
- Use the right rail for real deployment health, active environment, deployment counts, and recovery actions.
- Do not wrap filters, list, and rail in one outer card.

### Backups

- Keep `/backups` independent from Settings.
- Link S3 destinations, backup retention, restore history, and storage status from the backups page.
- Do not make Settings and Backups active simultaneously.
- Use a centered destination-picker dialog only for a short destination selection flow; move long credential or verification flows into a dedicated route.

### DNS and domains

- Keep `/dns` as the organization-wide audit and provider configuration surface.
- Keep service-level domain management on `/services/$serviceId/domains`.
- Link the audit from Home shortcuts, service domains, project rail, and any DNS warning state.
- Decide whether the public name is `Domains and DNS` or `DNS`; use one label consistently in navigation and headings.

### Sources and servers

- Sources: discover through Home, project create/import, and project source configuration.
- Servers: discover through Home, deploy target selection, and no-server project/app empty states.
- Neither belongs in the primary sidebar unless product usage data proves otherwise.
- Treat server onboarding as a page when it includes install commands, agent registration, SSH credentials, validation, or progressive status. Keep destructive disconnect/remove actions as dialogs.

### Settings and administration

- Keep `/settings` for general, notifications, OAuth, integrations, security, and instance configuration.
- Give Settings one `Settings` page heading and subtitle. Content cards own domain-specific headings; the selected navigation label must not be repeated as another page heading.
- Use a sticky desktop settings rail with a static Codedock settings identity card and a separate settings-navigation card.
- Move maintenance, updates, migration, audit logs, members, and AI to named Settings subsections or explicitly linked settings cards before removing their current routes.
- Keep API Access separate because it is an operator workflow, not profile preference.

## Component Work Packages

### 1. Navigation contract

- Create a typed navigation registry that records title, route, owner, primary/sidebar location, contextual location, mobile fallback, and permission requirement.
- Drive sidebar, command palette, Home shortcuts, and settings links from the registry where appropriate.
- Add a test that every internal navigation target exists in TanStack Router.
- Replace obsolete command palette entries: `/teams`, `/templates`, `/notifications`, and `/terminal` currently do not correspond to Codedock dashboard routes.

### 2. Page shell and rails

- Create `PageFrame`, `PageHeader`, `ContextRail`, and `MobileContextNav` components.
- Use semantic `main`, `aside`, `nav`, headings, and labelled navigation regions.
- Keep rails optional and data-driven; do not render empty decorative cards.
- Make every populated desktop rail sticky using the shared component rather than per-page ad hoc classes.
- Use the same 340px rail breakpoint, 24px gap, and top offset across list and detail pages.
- Use a horizontal mobile contextual navigation pattern before the main content. Do not simply hide destinations on mobile.

### 3. Route versus dialog policy

| Interaction | Use a route when | Use a dialog when | Codedock decision |
| --- | --- | --- | --- |
| Create project | The user selects source, template, build/runtime, environment, or deployment intent | Never for the full workflow | Add `/projects/new` and retire `CreateProjectModal` after parity |
| Create app | The user selects a catalog app, project, storage, credentials, or configuration | Never for catalog deployment | Add `/apps/new` and `/apps/new/$templateId` |
| Add server | The user configures credentials, agent install, verification, or multiple steps | A short name-only action with no follow-up | Add `/servers/new`; preserve existing command handoff and status polling |
| Add backup destination | Credentials or verification span several steps | Picker plus compact credential form | Keep dialog only after validating it remains bounded; otherwise add `/backups/new` |
| Edit service/project settings | The edit needs direct URLs, deep links, review, or multiple sections | A small single-field edit | Keep as detail-page sections |
| Delete, revoke, detach, reset | Never | Confirmation with clear consequence and error recovery | Keep dialogs |
| Quick create from empty state | It launches a primary workflow | It creates a trivial local object with no configuration | Link to the corresponding creation route |

### 4. Visual token pass

- Define dark-surface levels for canvas, sidebar, card, elevated popover, muted control, and input.
- Increase contrast through token values, not scattered opacity overrides.
- Remove `bg-card/40`, `bg-background/50`, and low-opacity border use from primary information surfaces unless translucency is intentional.
- Standardize card border, radius, hover, focus, and active states.
- Audit every existing card component for solid surface use and visible borders in dark mode.

### 5. Header hierarchy pass

- Define a single header ownership rule for every route.
- Parent layouts provide context only; leaf pages own the visible title unless the context itself is the title.
- Remove duplicate labels such as `Settings` plus `General` plus a repeated `General` rail heading.
- Apply the same rule to project and service layouts.

### 6. Responsive interaction pass

- Sidebar opens as an overlay under desktop breakpoint and closes on route change.
- Desktop rails become horizontal navigation or a sheet on mobile.
- Verify 320px, 375px, 768px, 1024px, 1440px, and 1920px layouts.
- Verify keyboard flow, focus restore, roving navigation where relevant, and no clipped content.

### 7. Permission and state pass

- Centralize visibility rules for admin, owner, member, and viewer navigation.
- Provide unavailable and permission-denied states instead of silently hiding the destination after a user follows an existing link.
- Confirm organization switching invalidates and refetches organization-scoped data.
- Keep command palette results permission-aware.

## Delivery Order

1. Freeze sidebar removals and create the typed navigation registry.
2. Complete the route-discovery matrix with a product owner decision for each ambiguous route.
3. Build shared page frame, header, sticky right rail, and mobile contextual navigation primitives.
4. Add the project, app, and server creation routes before retiring their dialogs.
5. Refactor Home, Projects, Apps, Deployments, and Backups using those primitives.
6. Refactor project and service detail layouts without losing nested routes.
7. Refactor Settings and administration routes into explicit settings subsections.
8. Apply the visual token pass across all dashboard surfaces.
9. Test every route, breakpoint, role, empty state, sticky rail, and command-palette target.

## Acceptance Checklist

- Every Codedock route has a documented owner and at least one working discovery path.
- Primary sidebar has no more than seven persistent product entries plus the project action and workspace footer.
- Backups and Settings are separate routes with mutually exclusive active states.
- Servers, domains, sources, organizations, and instance tools are reachable without being forced into primary navigation.
- No screen has duplicated page, tab, and section headings.
- No generic outer card wraps unrelated page cards.
- Dark mode surfaces have visible, intentional hierarchy at normal display brightness.
- Every desktop contextual rail remains visible while its main-column content scrolls, and each has a mobile equivalent.
- Long-running or multi-step creation workflows have URL-addressable routes rather than dialogs.
- Command palette contains only real routes and implemented actions.
- Desktop and mobile contextual navigation expose the same destinations.
- Route generation, formatting, typecheck, navigation tests, and visual regression checks pass.
