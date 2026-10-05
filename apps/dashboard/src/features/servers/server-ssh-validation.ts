export function validateSshConnection(
  method: 'key' | 'password',
  key: string,
  password: string,
  port: string
): string | null {
  if ((method === 'key' && !key.trim()) || (method === 'password' && !password))
    return 'Provide the selected SSH credential';
  const value = Number(port);
  if (!Number.isInteger(value) || value < 1 || value > 65535)
    return 'SSH port must be between 1 and 65535';
  return null;
}

export async function readSshPrivateKey(file: File): Promise<string> {
  if (file.size > 64 * 1024) throw new Error('Key file exceeds maximum size of 64KB');
  const text = await file.text();
  if (!text.includes('PRIVATE KEY-----')) throw new Error('Select an SSH private key file');
  return text;
}
