# TODO

## Outpipe tunnel integration

Optional integration with the standalone Outpipe tunnel platform. The platform
contract (API endpoints, webhook events and signatures, metadata convention)
is implemented and documented in the outpipe repo:
`docs/codedock-integration.md` in ~/Dev/TechX/codedock-tunnel.

- [ ] Define an optional Codedock adapter that consumes the standalone tunnel API.
- [ ] Allow Codedock to create scoped tunnel credentials through the public integration API.
- [ ] Add optional tunnel metadata links to Codedock projects and services.
- [ ] Add tunnel status synchronization through polling or signed webhooks.
- [ ] Add optional tunnel creation after local-directory deployment.
- [ ] Add optional MCP tools for listing, creating, inspecting, and revoking tunnels.
- [ ] Keep Codedock usable when the tunnel service is unavailable.
- [ ] Ensure uninstalling or disabling Codedock integration never deletes standalone tunnel data.
- [ ] Add integration tests proving a standalone tunnel server works without Codedock.
