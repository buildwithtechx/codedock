---
title: API Endpoint Catalogue
description: Registered API methods, paths and middleware requirements.
---

This catalogue follows the current route registrations. Access labels describe route middleware only; handlers also check resource ownership, account state, credentials or signatures. WebSockets authenticate in their own handlers. `Handler / middleware` does not mean anonymous access is allowed.

Read [API conventions](/api/) first. Paths use `:parameter` placeholders. Linked definitions identify each handler; request DTOs and response behavior live in the handler and model files. The exhaustive machine-readable schema is `docs/api/openapi.json` in the repository, also served at runtime as `GET /api/docs/openapi.json`; the Go suite fails when a route lacks a registry entry.

## Auth

| Method | Path | Route access | Handler wiring |
| --- | --- | --- | --- |
| GET | `/api/auth/csrf` | CSRF bootstrap; no login required | [`bootstrapCSRF`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| POST | `/api/auth/signup` | Handler / middleware | [`authHandler.Register`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| POST | `/api/auth/signin` | Handler / middleware | [`authHandler.Login`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| POST | `/api/auth/refresh` | Handler / middleware | [`authHandler.Refresh`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| POST | `/api/auth/forgot-password` | Handler / middleware | [`authHandler.ForgotPassword`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| POST | `/api/auth/reset-password` | Handler / middleware | [`authHandler.ResetPassword`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| POST | `/api/auth/email/resend` | Handler / middleware | [`authHandler.ResendVerificationEmail`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| POST | `/api/auth/email/verify` | Handler / middleware | [`authHandler.VerifyEmail`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| POST | `/api/auth/logout` | Handler / middleware | [`authHandler.Logout`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| GET | `/api/auth/me` | Authenticated | [`userHandler.GetProfile`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| GET | `/api/auth/oauth/providers/enabled` | Handler / middleware | [`oauthHandler.ListEnabledProviders`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| GET | `/api/auth/oauth/:provider` | Handler / middleware | [`oauthHandler.OAuthRedirect`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| GET | `/api/auth/oauth/:provider/callback` | Handler / middleware | [`oauthHandler.OAuthCallback`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| POST | `/api/auth/2fa/setup` | Authenticated | [`oauthHandler.Setup2FA`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| POST | `/api/auth/2fa/verify` | Authenticated | [`oauthHandler.Verify2FA`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| POST | `/api/auth/2fa/disable` | Authenticated | [`oauthHandler.Disable2FA`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |

## System

| Method | Path | Route access | Handler wiring |
| --- | --- | --- | --- |
| GET | `/api/system/public` | Handler / middleware | [`settingsHandler.GetPublicSettings`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| GET | `/api/system/setup-status` | Handler / middleware | [`onboardingHandler.SetupStatus`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| POST | `/api/system/setup/import` | Handler / middleware | [`migrationHandler.ImportDuringSetup`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| GET | `/api/system/stats` | Authenticated | [`systemHandler.GetStats`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| POST | `/api/system/restart` | Handler / middleware; instance admin | [`systemHandler.Restart`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| POST | `/api/system/maintenance/cleanup` | Handler / middleware; instance admin | [`systemHandler.Cleanup`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| POST | `/api/system/export` | Handler / middleware; instance admin | [`migrationHandler.Export`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| POST | `/api/system/import` | Handler / middleware; instance admin | [`migrationHandler.Import`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| POST | `/api/system/takeover/scan` | Handler / middleware; instance admin | [`takeoverHandler.Scan`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| POST | `/api/system/takeover/adopt` | Handler / middleware; instance admin | [`takeoverHandler.Adopt`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| GET | `/api/system/takeover/runs` | Authenticated; instance admin | [`takeoverHandler.ListRuns`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| GET | `/api/system/takeover/runs/:id` | Authenticated; instance admin | [`takeoverHandler.GetRun`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |

## User

| Method | Path | Route access | Handler wiring |
| --- | --- | --- | --- |
| GET | `/api/users` | Authenticated; instance admin | [`userHandler.ListUsers`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| POST | `/api/users/invite` | Authenticated; instance admin | [`authHandler.AdminInviteUser`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| DELETE | `/api/users/:id` | Authenticated; instance admin | [`userHandler.DeleteUser`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| GET | `/api/profile` | Authenticated | [`userHandler.GetProfile`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| PUT | `/api/profile` | Authenticated | [`userHandler.UpdateProfile`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| POST | `/api/profile/email/request` | Authenticated | [`userHandler.RequestEmailChange`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| POST | `/api/profile/email/verify` | Authenticated | [`userHandler.VerifyEmailChange`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| PUT | `/api/profile/password` | Authenticated | [`userHandler.ChangePassword`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| GET | `/api/profile/tokens` | Authenticated | [`userHandler.ListPATs`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| POST | `/api/profile/tokens` | Authenticated | [`userHandler.CreatePAT`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| DELETE | `/api/profile/tokens/:id` | Authenticated | [`userHandler.DeletePAT`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |

## Project

| Method | Path | Route access | Handler wiring |
| --- | --- | --- | --- |
| GET | `/api/projects` | Authenticated | [`projectHandler.ListProjects`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| POST | `/api/projects` | Authenticated | [`projectHandler.CreateProject`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| GET | `/api/projects/:id` | Authenticated; project access | [`projectHandler.GetProject`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| DELETE | `/api/projects/:id` | Authenticated; project owner | [`projectHandler.DeleteProject`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| GET | `/api/domains` | Authenticated | [`domainHandler.ListAll`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| GET | `/api/services/:id/domains` | Authenticated | [`domainHandler.ListByService`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| POST | `/api/services/:id/domains` | Authenticated | [`domainHandler.Create`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| DELETE | `/api/domains/:id` | Authenticated | [`domainHandler.Delete`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| POST | `/api/domains/:id/verify` | Authenticated | [`domainHandler.Verify`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| GET | `/api/domains/:id/verify` | Authenticated | [`domainHandler.Verify`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| GET | `/api/projects/:id/env` | Authenticated; project access; env:read | [`projectEnvHandler.GetVars`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| PUT | `/api/projects/:id/env` | Authenticated; project administrator; env:write | [`projectEnvHandler.SetVars`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| POST | `/api/projects/:id/env` | Authenticated; project administrator; env:write | [`projectEnvHandler.SetVars`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| POST | `/api/projects/:id/environments` | Authenticated; project administrator | [`environmentHandler.Create`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| GET | `/api/projects/:id/environments` | Authenticated; project access | [`environmentHandler.ListByProject`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| GET | `/api/projects/:id/apps` | Authenticated; project access | [`appServiceHandler.ListByProject`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| GET | `/api/projects/:id/services` | Authenticated; project access | [`appServiceHandler.ListByProject`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| GET | `/api/projects/:id/deployments` | Authenticated; project access | [`deploymentHandler.ListProjectDeployments`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| POST | `/api/projects/:id/deploy` | Authenticated; project administrator | [`deploymentHandler.TriggerProject`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| GET | `/api/project-apps` | Authenticated | [`projectAppHandler.List`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| POST | `/api/project-apps` | Authenticated | [`projectAppHandler.Create`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| GET | `/api/project-apps/:id` | Authenticated | [`projectAppHandler.Get`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| PUT | `/api/project-apps/:id` | Authenticated | [`projectAppHandler.Update`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| DELETE | `/api/project-apps/:id` | Authenticated | [`projectAppHandler.Delete`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| GET | `/api/projects/:projectId/tokens` | Authenticated; project administrator; env:read | [`projectSettingsHandler.ListTokens`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| POST | `/api/projects/:projectId/tokens` | Authenticated; project administrator; env:write | [`projectSettingsHandler.CreateToken`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| DELETE | `/api/projects/:projectId/tokens/:id` | Authenticated; project administrator; env:write | [`projectSettingsHandler.DeleteToken`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| GET | `/api/projects/:projectId/registries` | Authenticated; project administrator | [`registryHandler.List`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| POST | `/api/projects/:projectId/registries` | Authenticated; project administrator | [`registryHandler.Create`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| DELETE | `/api/projects/:projectId/registries/:id` | Authenticated; project administrator | [`registryHandler.Delete`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| GET | `/api/services/:serviceId/route-rules` | Authenticated; project/service access | [`routeRuleHandler.List`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| POST | `/api/services/:serviceId/route-rules` | Authenticated; project/service administrator | [`routeRuleHandler.Create`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| PATCH | `/api/services/:serviceId/route-rules/:ruleId` | Authenticated; project/service administrator | [`routeRuleHandler.Update`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| DELETE | `/api/services/:serviceId/route-rules/:ruleId` | Authenticated; project/service administrator | [`routeRuleHandler.Delete`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |

## Server

| Method | Path | Route access | Handler wiring |
| --- | --- | --- | --- |
| GET | `/api/servers` | Authenticated; server:read | [`serverHandler.List`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| GET | `/api/servers/:id` | Authenticated; server:read | [`serverHandler.Get`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| POST | `/api/servers` | Authenticated; server:write | [`serverHandler.Create`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| PATCH | `/api/servers/:id` | Authenticated; server:write | [`serverHandler.Update`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| POST | `/api/servers/test-ssh` | Authenticated; server:write | [`serverHandler.TestSSH`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| DELETE | `/api/servers/:id` | Authenticated; server:write | [`serverHandler.Delete`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| GET | `/api/ws/servers/:serverId/metrics` | Handler / middleware | [`serverMetricsWSHandler.Handle`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |

## Organization

| Method | Path | Route access | Handler wiring |
| --- | --- | --- | --- |
| GET | `/api/organizations` | Authenticated | [`orgHandler.List`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| POST | `/api/organizations` | Authenticated | [`orgHandler.Create`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| GET | `/api/organizations/:id` | Authenticated; organization member | [`orgHandler.Get`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| DELETE | `/api/organizations/:id` | Authenticated; organization owner | [`orgHandler.Delete`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| GET | `/api/organizations/:id/members` | Authenticated; organization member | [`orgHandler.ListMembers`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| POST | `/api/organizations/:id/members` | Authenticated; organization admin | [`orgHandler.InviteMember`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| PUT | `/api/organizations/:id/members/:userId` | Authenticated; organization admin | [`orgHandler.UpdateMember`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| DELETE | `/api/organizations/:id/members/:memberId` | Authenticated; organization admin | [`orgHandler.RemoveMember`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |

## Database

| Method | Path | Route access | Handler wiring |
| --- | --- | --- | --- |
| GET | `/api/databases` | Authenticated; database:manage | [`dbHandler.ListDatabases`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| POST | `/api/databases` | Authenticated; database:manage | [`dbHandler.CreateDatabase`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| GET | `/api/databases/:id` | Authenticated; database:manage | [`dbHandler.GetDatabase`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| PUT | `/api/databases/:id` | Authenticated; database:manage | [`dbHandler.UpdateDatabase`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| DELETE | `/api/databases/:id` | Authenticated; database:manage | [`dbHandler.DeleteDatabase`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| POST | `/api/databases/:id/start` | Authenticated; database:manage | [`dbHandler.StartDatabase`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| POST | `/api/databases/:id/stop` | Authenticated; database:manage | [`dbHandler.StopDatabase`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| POST | `/api/databases/:id/restart` | Authenticated; database:manage | [`dbHandler.RestartDatabase`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| POST | `/api/databases/:id/credentials/reveal` | Authenticated; database:manage | [`dbHandler.RevealCredentials`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| POST | `/api/databases/:id/query` | Authenticated; database:manage | [`dbHandler.QueryDatabase`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| POST | `/api/databases/:id/import` | Authenticated; database:manage | [`dbHandler.ImportData`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| GET | `/api/databases/:id/schemas` | Authenticated; database:manage | [`dbHandler.GetSchemas`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| GET | `/api/databases/:id/data/:table` | Authenticated; database:manage | [`dbHandler.GetTableData`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| POST | `/api/databases/:id/data/:table` | Authenticated; database:manage | [`dbHandler.InsertTableRow`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| PUT | `/api/databases/:id/data/:table` | Authenticated; database:manage | [`dbHandler.UpdateTableRow`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| DELETE | `/api/databases/:id/data/:table` | Authenticated; database:manage | [`dbHandler.DeleteTableRow`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| GET | `/api/databases/:id/backups` | Authenticated; database:manage | [`backupHandler.ListRecordsByDatabase`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |
| POST | `/api/databases/:id/backups` | Authenticated; database:manage | [`backupHandler.TriggerDatabaseBackup`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes.go) |

## App

| Method | Path | Route access | Handler wiring |
| --- | --- | --- | --- |
| GET | `/api/environments/:id/apps` | Authenticated | [`appServiceHandler.ListByEnvironment`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| POST | `/api/environments/:id/apps` | Authenticated | [`appServiceHandler.Create`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| DELETE | `/api/environments/:id` | Authenticated | [`environmentHandler.Delete`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| GET | `/api/apps` | Authenticated | [`appServiceHandler.ListByOrganization`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| GET | `/api/apps/:id` | Authenticated; project/service access | [`appServiceHandler.Get`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| PUT | `/api/apps/:id` | Authenticated; project/service administrator | [`appServiceHandler.Update`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| DELETE | `/api/apps/:id` | Authenticated; project/service owner | [`appServiceHandler.Delete`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| POST | `/api/apps/:id/stop` | Authenticated; project/service administrator | [`appServiceHandler.StopService`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| POST | `/api/apps/:id/redeploy` | Authenticated; project/service administrator | [`appServiceHandler.RedeployService`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| POST | `/api/apps/:id/restart` | Authenticated; project/service administrator | [`appServiceHandler.RestartService`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| GET | `/api/apps/:id/webhooks` | Authenticated; project/service access | [`appServiceHandler.ListWebhooks`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| POST | `/api/apps/:id/webhooks` | Authenticated; project/service administrator | [`appServiceHandler.CreateWebhook`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| DELETE | `/api/apps/:id/webhooks/:webhookId` | Authenticated; project/service administrator | [`appServiceHandler.DeleteWebhook`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| GET | `/api/apps/:id/volumes` | Authenticated; project/service access | [`appServiceHandler.ListVolumes`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| POST | `/api/apps/:id/volumes` | Authenticated; project/service administrator | [`appServiceHandler.CreateVolume`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| DELETE | `/api/apps/:id/volumes/:volumeId` | Authenticated; project/service administrator | [`appServiceHandler.DeleteVolume`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| GET | `/api/apps/:id/log-drains` | Authenticated; project/service access | [`appServiceHandler.ListLogDrains`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| POST | `/api/apps/:id/log-drains` | Authenticated; project/service administrator | [`appServiceHandler.CreateLogDrain`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| DELETE | `/api/apps/:id/log-drains/:drainId` | Authenticated; project/service administrator | [`appServiceHandler.DeleteLogDrain`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| GET | `/api/services/:serviceId/variables` | Authenticated; project/service access | [`serviceVarHandler.List`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| GET | `/api/services/:serviceId/env-suggestions` | Authenticated; project/service access | [`serviceVarHandler.Suggest`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| POST | `/api/services/:serviceId/variables` | Authenticated; project/service administrator | [`serviceVarHandler.Create`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| PUT | `/api/services/:serviceId/variables/:id` | Authenticated; project/service administrator | [`serviceVarHandler.Update`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| DELETE | `/api/services/:serviceId/variables/:id` | Authenticated; project/service administrator | [`serviceVarHandler.Delete`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| GET | `/api/services/:serviceId/serverless/code` | Authenticated; project/service access | [`serverlessHandler.GetCode`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| POST | `/api/services/:serviceId/serverless/code` | Authenticated; project/service administrator | [`serverlessHandler.SaveCode`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| GET | `/api/services/:serviceId/logs` | Handler / middleware | [`serviceLogsWSHandler.Handle`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |

## Deployment

| Method | Path | Route access | Handler wiring |
| --- | --- | --- | --- |
| GET | `/api/deployments` | Authenticated | [`deploymentHandler.ListOrganizationDeployments`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| GET | `/api/services/:serviceId/deployments` | Authenticated; project/service access | [`deploymentHandler.ListServiceDeployments`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| GET | `/api/services/:serviceId/previews` | Authenticated; project/service access | [`deploymentHandler.ListPRPreviews`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| POST | `/api/services/:serviceId/deploy` | Authenticated; project/service administrator | [`deploymentHandler.Trigger`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| POST | `/api/deployments/:id/rollback` | Authenticated | [`deploymentHandler.Rollback`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| GET | `/api/deployments/:id/logs` | Authenticated; logs:read | [`deploymentHandler.GetLogs`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| GET | `/api/deployments/:id/explain` | Authenticated | [`deploymentHandler.ExplainFailure`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| GET | `/api/services/:serviceId/metrics` | Authenticated; project/service access | [`deploymentHandler.GetMetrics`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| GET | `/api/services/:serviceId/metrics/historical` | Authenticated; project/service access | [`metricsHandler.GetHistoricalMetrics`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| GET | `/api/services/:serviceId/logs/historical` | Authenticated; project/service access | [`logHandler.GetHistoricalLogs`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |

## Backup

| Method | Path | Route access | Handler wiring |
| --- | --- | --- | --- |
| GET | `/api/backups` | Authenticated | [`backupHandler.List`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| POST | `/api/backups` | Authenticated | [`backupHandler.Create`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| GET | `/api/backups/:id` | Authenticated | [`backupHandler.Get`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| PUT | `/api/backups/:id` | Authenticated; backup:write | [`backupHandler.Update`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| DELETE | `/api/backups/:id` | Authenticated | [`backupHandler.Delete`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| POST | `/api/backups/:id/trigger` | Authenticated | [`backupHandler.Trigger`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| POST | `/api/backups/:id/restore` | Authenticated | [`backupHandler.Restore`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| GET | `/api/backups/:id/records` | Authenticated | [`backupHandler.ListRecords`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| GET | `/api/backups/:id/records/:recordId/download` | Authenticated | [`backupHandler.DownloadRecord`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| DELETE | `/api/backups/:id/records/:recordId` | Authenticated | [`backupHandler.DeleteRecord`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| GET | `/api/backup-records` | Authenticated | [`backupHandler.ListAllRecords`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| GET | `/api/s3-destinations` | Authenticated; instance admin | [`backupHandler.ListS3Destinations`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| POST | `/api/s3-destinations` | Authenticated; instance admin | [`backupHandler.CreateS3Destination`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| POST | `/api/s3-destinations/verify` | Authenticated; instance admin | [`backupHandler.VerifyS3Draft`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| PUT | `/api/s3-destinations/:id` | Authenticated; instance admin | [`backupHandler.UpdateS3Destination`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| POST | `/api/s3-destinations/:id/verify` | Authenticated; instance admin | [`backupHandler.VerifyS3Destination`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| POST | `/api/s3-destinations/:id/default` | Authenticated; instance admin | [`backupHandler.SetDefaultS3Destination`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| DELETE | `/api/s3-destinations/:id` | Authenticated; instance admin | [`backupHandler.DeleteS3Destination`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |

## Settings

| Method | Path | Route access | Handler wiring |
| --- | --- | --- | --- |
| GET | `/api/settings` | Authenticated | [`settingsHandler.GetSettings`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| PUT | `/api/settings` | Handler / middleware; instance admin | [`settingsHandler.UpdateSettings`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| GET | `/api/ai` | Authenticated | [`aiSettingsHandler.GetAISettings`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| POST | `/api/ai/diagnose` | Authenticated | [`aiSettingsHandler.DiagnoseLogs`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| PUT | `/api/ai` | Handler / middleware; instance admin | [`aiSettingsHandler.UpdateAISettings`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| GET | `/api/notifications` | Authenticated | [`notifSettingsHandler.GetNotificationSettings`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| PUT | `/api/notifications` | Handler / middleware; instance admin | [`notifSettingsHandler.UpdateNotificationSettings`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| GET | `/api/settings/updates/status` | Authenticated | [`updaterHandler.GetUpdateStatus`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| POST | `/api/settings/updates/check` | Handler / middleware; instance admin | [`updaterHandler.CheckUpdate`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| POST | `/api/settings/updates/deploy` | Handler / middleware; instance admin | [`updaterHandler.DeployUpdate`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| GET | `/api/settings/oauth/providers` | Handler / middleware; instance admin | [`oauthHandler.ListProviders`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| PUT | `/api/settings/oauth/providers` | Handler / middleware; instance admin | [`oauthHandler.SaveProvider`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| POST | `/api/settings/git_apps/github/manifest-callback` | Handler / middleware; instance admin | [`gitAppsHandler.ExchangeGithubManifestCode`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| GET | `/api/settings/git_apps/github` | Authenticated | [`gitAppsHandler.ListGithubApps`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| GET | `/api/settings/git_apps/github/:id` | Authenticated | [`gitAppsHandler.GetGithubApp`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| PUT | `/api/settings/git_apps/github` | Handler / middleware; instance admin | [`gitAppsHandler.SaveGithubApp`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| DELETE | `/api/settings/git_apps/github/:id` | Handler / middleware; instance admin | [`gitAppsHandler.DeleteGithubApp`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| POST | `/api/settings/notifications/test` | Handler / middleware; instance admin | [`notificationHandler.TestNotification`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |

## Misc

| Method | Path | Route access | Handler wiring |
| --- | --- | --- | --- |
| POST | `/api/compose/deploy` | Authenticated | [`composeHandler.Deploy`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| POST | `/api/compose/analyze` | Authenticated | [`composeHandler.Analyze`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| POST | `/api/deploy/archive` | Authenticated | [`archiveHandler.DeployArchive`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| GET | `/api/examples` | Authenticated | [`exampleHandler.List`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| GET | `/api/one-click` | Authenticated | [`oneClickHandler.List`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| POST | `/api/one-click/deploy` | Authenticated | [`oneClickHandler.Deploy`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| POST | `/api/dns` | Authenticated; instance admin | [`dnsHandler.Create`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| GET | `/api/dns` | Authenticated; instance admin | [`dnsHandler.List`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| PUT | `/api/dns/:id` | Authenticated; instance admin | [`dnsHandler.Update`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| DELETE | `/api/dns/:id` | Authenticated; instance admin | [`dnsHandler.Delete`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| GET | `/api/scheduled-tasks` | Authenticated | [`scheduledTaskHandler.ListProjectScheduledTasks`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| POST | `/api/scheduled-tasks` | Authenticated | [`scheduledTaskHandler.Create`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| GET | `/api/scheduled-tasks/:id` | Authenticated | [`scheduledTaskHandler.Get`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| DELETE | `/api/scheduled-tasks/:id` | Authenticated | [`scheduledTaskHandler.Delete`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| POST | `/api/scheduled-tasks/:id/trigger` | Authenticated | [`scheduledTaskHandler.Run`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| POST | `/api/git/connect` | Authenticated | [`gitHandler.Connect`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| GET | `/api/git/status` | Authenticated | [`gitHandler.Status`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| DELETE | `/api/git/connect/:provider` | Authenticated | [`gitHandler.Disconnect`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| GET | `/api/git/repos` | Authenticated | [`gitHandler.ListRepos`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| POST | `/api/webhooks/git/services/:serviceId` | Handler / middleware | [`webhookHandler.HandleServiceWebhook`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| POST | `/api/webhooks/github/services/:serviceId` | Handler / middleware | [`webhookHandler.HandleGitHubWebhook`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| GET | `/api/canvas/projects` | Authenticated | [`canvasHandler.ListCanvasSummaries`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| GET | `/api/projects/:id/summary` | Authenticated | [`canvasHandler.GetCanvasSummary`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| GET | `/api/environments/:id/canvas` | Authenticated | [`canvasHandler.GetEnvironmentCanvas`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| GET | `/api/audit-logs` | Authenticated; instance admin | [`auditLogHandler.List`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| GET | `/api/mcp/sse` | Authenticated | [`HandleMCPSSE`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| POST | `/api/mcp/messages` | Authenticated | [`HandleMCPMessage`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| GET | `/api/ws/terminal/:id` | Handler / middleware | [`terminalHandler.HandleWebSocket`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| GET | `/api/ws/services/:id/terminal` | Handler / middleware | [`terminalHandler.HandleWebSocket`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |

## Billing

| Method | Path | Route access | Handler wiring |
| --- | --- | --- | --- |
| POST | `/api/billing/webhook` | Handler / middleware; cloud only | [`billingHandler.Webhook`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| GET | `/api/billing/config` | Authenticated; cloud only | [`billingHandler.GetConfig`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |
| POST | `/api/billing/checkout` | Authenticated; cloud only | [`billingHandler.CreateCheckoutSession`](https://github.com/buildwithtechx/codedock/blob/main/internal/http/routes_app.go) |

## Outside the API prefix

`GET /healthz` returns health status. Static assets and the dashboard fallback are registered outside this catalogue. The API-prefixed OAuth callbacks are included in the Auth section. Billing routes return `404` when cloud mode is disabled.
