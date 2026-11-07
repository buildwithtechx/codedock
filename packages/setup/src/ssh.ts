import type { ServerSshSettings, SetupServerPayload, SshPayload } from './types';

export function buildSshSettings(payload: SshPayload): ServerSshSettings {
  const port = Number(payload.port ?? 22);
  const normalizedPort = Number.isInteger(port) && port > 0 && port <= 65535 ? port : 22;

  return {
    sshHost: payload.host.trim(),
    sshPort: normalizedPort,
    sshUser: payload.user?.trim() || 'root',
    sshAuthMethod: payload.authMethod,
    sshKey: payload.authMethod === 'key' ? payload.privateKey?.trim() : undefined,
    sshPassword: payload.authMethod === 'password' ? payload.password : undefined,
    sshJumpHost: payload.jumpHost?.trim() || undefined,
    sshTransport: payload.transport || 'direct',
  };
}

export function buildCreateServerPayload(name: string, ssh: SshPayload): SetupServerPayload {
  const settings = buildSshSettings(ssh);
  return {
    name: name.trim() || settings.sshHost,
    ipAddress: settings.sshHost,
    sshPort: settings.sshPort,
    sshUser: settings.sshUser,
    sshAuthMethod: settings.sshAuthMethod,
    sshKey: settings.sshKey,
    sshPassword: settings.sshPassword,
    sshJumpHost: settings.sshJumpHost,
  };
}
