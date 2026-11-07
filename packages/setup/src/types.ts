export type SshAuthMethod = 'key' | 'password';
export type SshTransport = 'direct' | 'cloudflare';

export interface SshPayload {
  host: string;
  port?: number | string;
  user?: string;
  authMethod: SshAuthMethod;
  privateKey?: string;
  password?: string;
  jumpHost?: string;
  transport?: SshTransport;
}

export interface ServerSshSettings {
  sshHost: string;
  sshPort: number;
  sshUser: string;
  sshAuthMethod: SshAuthMethod;
  sshKey?: string;
  sshPassword?: string;
  sshJumpHost?: string;
  sshTransport: SshTransport;
}

export type SetupTarget = 'managed' | 'connected' | 'local';

export interface ProviderConfig {
  provider: 'hetzner' | 'custom';
  credentialId?: string;
  location?: string;
  serverType?: string;
  image?: string;
  apiKey?: string;
}

export type SetupStep =
  | 'target'
  | 'credential'
  | 'plan'
  | 'ssh'
  | 'review'
  | 'progress'
  | 'complete';

export interface SetupState {
  step: SetupStep;
  target: SetupTarget;
  serverName?: string;
  ssh?: SshPayload;
  provider?: ProviderConfig;
  error?: string | null;
}

export interface SetupServerPayload {
  name: string;
  ipAddress: string;
  sshPort: number;
  sshUser: string;
  sshAuthMethod: SshAuthMethod;
  sshKey?: string;
  sshPassword?: string;
  sshJumpHost?: string;
}
