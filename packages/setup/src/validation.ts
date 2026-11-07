import type { ProviderConfig, SshAuthMethod, SshPayload } from './types';

export function isPrivateIp(ip: string): boolean {
  const trimmed = ip.trim().toLowerCase();
  if (trimmed === 'localhost' || trimmed === '::1' || trimmed === '127.0.0.1') {
    return true;
  }

  // IPv4 checks
  const ipv4Match = trimmed.match(/^(\d{1,3})\.(\d{1,3})\.(\d{1,3})\.(\d{1,3})$/);
  if (ipv4Match) {
    const [, a, b] = ipv4Match.map(Number);
    if (a === 10) return true;
    if (a === 127) return true;
    if (a === 169 && b === 254) return true;
    if (a === 172 && b >= 16 && b <= 31) return true;
    if (a === 192 && b === 168) return true;
    return false;
  }

  // IPv6 unique local (fc00::/7) or link-local (fe80::/10)
  if (trimmed.startsWith('fc') || trimmed.startsWith('fd') || trimmed.startsWith('fe80:')) {
    return true;
  }

  return false;
}

export function validateServerAddress(address: string): string | null {
  const trimmed = address.trim();
  if (!trimmed) {
    return 'Host address is required';
  }

  // Allow standard domain names, hostnames, IPv4, and IPv6
  const domainRegex =
    /^(?:(?:[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}|localhost)$/;
  const ipv4Regex =
    /^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/;
  const ipv6Regex = /^[0-9a-fA-F:]+$/;

  if (
    !ipv4Regex.test(trimmed) &&
    !domainRegex.test(trimmed) &&
    !(ipv6Regex.test(trimmed) && trimmed.includes(':'))
  ) {
    return 'Invalid host address or domain format';
  }

  return null;
}

export function validateSshConnection(
  method: SshAuthMethod,
  key: string,
  password: string,
  port: string | number
): string | null {
  if ((method === 'key' && !key.trim()) || (method === 'password' && !password)) {
    return 'Provide the selected SSH credential';
  }
  const value = Number(port);
  if (!Number.isInteger(value) || value < 1 || value > 65535) {
    return 'SSH port must be between 1 and 65535';
  }
  return null;
}

export function validateSshPayload(payload: SshPayload): string | null {
  const addressError = validateServerAddress(payload.host);
  if (addressError) return addressError;

  return validateSshConnection(
    payload.authMethod,
    payload.privateKey || '',
    payload.password || '',
    payload.port ?? 22
  );
}

export async function readSshPrivateKey(
  fileOrText: { size: number; text: () => Promise<string> } | string
): Promise<string> {
  let content = '';
  if (typeof fileOrText === 'string') {
    content = fileOrText;
  } else {
    if (fileOrText.size > 64 * 1024) {
      throw new Error('Key file exceeds maximum size of 64KB');
    }
    content = await fileOrText.text();
  }

  if (!content.includes('PRIVATE KEY-----')) {
    throw new Error('Select an SSH private key file');
  }
  return content;
}

export function validateProviderConfig(config: ProviderConfig): string | null {
  if (!config.provider) {
    return 'Provider is required';
  }
  if (config.provider === 'hetzner') {
    if (!config.credentialId && !config.apiKey) {
      return 'Hetzner credential or API token is required';
    }
  }
  return null;
}
