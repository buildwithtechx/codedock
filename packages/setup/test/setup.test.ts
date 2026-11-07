import { describe, expect, it, vi } from 'vitest';
import {
  buildCreateServerPayload,
  buildSetupPayload,
  buildSshSettings,
  HETZNER_WIZARD_STEPS,
  isFirstStep,
  isLastStep,
  isPrivateIp,
  nextStep,
  prevStep,
  pushServerSetup,
  readSshPrivateKey,
  validateProviderConfig,
  validateServerAddress,
  validateSshConnection,
  validateSshPayload,
  waitForApi,
} from '../src/index';

describe('@codedock/setup', () => {
  describe('validation', () => {
    it('detects private and loopback IP addresses correctly', () => {
      expect(isPrivateIp('127.0.0.1')).toBe(true);
      expect(isPrivateIp('localhost')).toBe(true);
      expect(isPrivateIp('10.0.1.5')).toBe(true);
      expect(isPrivateIp('172.20.10.2')).toBe(true);
      expect(isPrivateIp('192.168.1.100')).toBe(true);
      expect(isPrivateIp('169.254.1.1')).toBe(true);
      expect(isPrivateIp('::1')).toBe(true);
      expect(isPrivateIp('fc00::1')).toBe(true);
      expect(isPrivateIp('8.8.8.8')).toBe(false);
      expect(isPrivateIp('1.1.1.1')).toBe(false);
      expect(isPrivateIp('codedock.example')).toBe(false);
    });

    it('validates server addresses and hostnames', () => {
      expect(validateServerAddress('192.168.1.1')).toBeNull();
      expect(validateServerAddress('server.example.com')).toBeNull();
      expect(validateServerAddress('localhost')).toBeNull();
      expect(validateServerAddress('')).toBe('Host address is required');
      expect(validateServerAddress('   ')).toBe('Host address is required');
      expect(validateServerAddress('invalid$$host@@')).toBe(
        'Invalid host address or domain format'
      );
    });

    it('validates SSH connection parameters', () => {
      expect(
        validateSshConnection('key', '-----BEGIN OPENSSH PRIVATE KEY-----', '', '22')
      ).toBeNull();
      expect(validateSshConnection('password', '', 'secret-pass', 22)).toBeNull();
      expect(validateSshConnection('key', '', '', '22')).toBe(
        'Provide the selected SSH credential'
      );
      expect(validateSshConnection('password', '', '', '22')).toBe(
        'Provide the selected SSH credential'
      );
      expect(validateSshConnection('key', 'valid-key', '', '0')).toBe(
        'SSH port must be between 1 and 65535'
      );
      expect(validateSshConnection('key', 'valid-key', '', '70000')).toBe(
        'SSH port must be between 1 and 65535'
      );
    });

    it('validates full SSH payloads', () => {
      expect(
        validateSshPayload({
          host: '10.0.0.1',
          authMethod: 'password',
          password: 'pass',
          port: 22,
        })
      ).toBeNull();

      expect(
        validateSshPayload({
          host: '',
          authMethod: 'password',
          password: 'pass',
        })
      ).toBe('Host address is required');
    });

    it('reads and validates SSH private key file text', async () => {
      const validKey = '-----BEGIN RSA PRIVATE KEY-----\nMIIE...\n-----END RSA PRIVATE KEY-----';
      await expect(readSshPrivateKey(validKey)).resolves.toBe(validKey);

      await expect(readSshPrivateKey('not a private key')).rejects.toThrow(
        'Select an SSH private key file'
      );

      const fakeFile = {
        size: 70 * 1024,
        text: async () => 'test',
      };
      await expect(readSshPrivateKey(fakeFile)).rejects.toThrow(
        'Key file exceeds maximum size of 64KB'
      );
    });

    it('validates provider configuration', () => {
      expect(validateProviderConfig({ provider: 'hetzner', credentialId: 'cred-1' })).toBeNull();
      expect(validateProviderConfig({ provider: 'hetzner' })).toBe(
        'Hetzner credential or API token is required'
      );
    });
  });

  describe('ssh normalization', () => {
    it('normalizes SSH settings with sensible defaults', () => {
      const normalized = buildSshSettings({
        host: '  10.0.0.5  ',
        authMethod: 'key',
        privateKey: '  key-content  ',
      });
      expect(normalized.sshHost).toBe('10.0.0.5');
      expect(normalized.sshPort).toBe(22);
      expect(normalized.sshUser).toBe('root');
      expect(normalized.sshKey).toBe('key-content');
      expect(normalized.sshTransport).toBe('direct');
    });

    it('builds server creation payload', () => {
      const payload = buildCreateServerPayload('Production Node', {
        host: '10.0.0.5',
        authMethod: 'password',
        password: 'secure',
        port: 2222,
      });
      expect(payload.name).toBe('Production Node');
      expect(payload.ipAddress).toBe('10.0.0.5');
      expect(payload.sshPort).toBe(2222);
      expect(payload.sshPassword).toBe('secure');
    });
  });

  describe('wizard flow state machine', () => {
    const steps = HETZNER_WIZARD_STEPS.map((s) => s.id);

    it('transitions forward and backward through steps', () => {
      expect(isFirstStep(steps, 'credential')).toBe(true);
      expect(nextStep(steps, 'credential')).toBe('plan');
      expect(nextStep(steps, 'plan')).toBe('review');
      expect(nextStep(steps, 'review')).toBe('progress');
      expect(nextStep(steps, 'progress')).toBeNull();
      expect(isLastStep(steps, 'progress')).toBe(true);

      expect(prevStep(steps, 'progress')).toBe('review');
      expect(prevStep(steps, 'credential')).toBeNull();
    });
  });

  describe('api client', () => {
    it('builds setup payload', () => {
      const payload = buildSetupPayload({
        name: 'Worker 1',
        host: '192.168.1.50',
        sshSettings: { sshPort: 2200, sshUser: 'ubuntu', sshAuthMethod: 'key' },
      });
      expect(payload.name).toBe('Worker 1');
      expect(payload.ipAddress).toBe('192.168.1.50');
      expect(payload.sshPort).toBe(2200);
      expect(payload.sshUser).toBe('ubuntu');
    });

    it('pushes server setup over HTTP', async () => {
      const mockFetch = vi.fn().mockResolvedValue({
        ok: true,
        status: 200,
        json: async () => ({ data: { id: 'srv-100', name: 'Worker 1' } }),
      });

      const result = await pushServerSetup(
        'http://localhost:8080/api',
        {
          name: 'Worker 1',
          ipAddress: '10.0.0.1',
          sshPort: 22,
          sshUser: 'root',
          sshAuthMethod: 'key',
        },
        'token-123',
        mockFetch
      );

      expect(result.id).toBe('srv-100');
      expect(mockFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/servers',
        expect.objectContaining({
          method: 'POST',
          headers: expect.objectContaining({
            Authorization: 'Bearer token-123',
          }),
        })
      );
    });

    it('waits for API availability with timeout', async () => {
      let attempts = 0;
      const mockFetch = vi.fn().mockImplementation(async () => {
        attempts++;
        if (attempts < 2) throw new Error('connection refused');
        return { ok: true, status: 200 };
      });

      const isReady = await waitForApi('http://localhost:8080/health', 5000, 10, mockFetch);
      expect(isReady).toBe(true);
      expect(attempts).toBe(2);
    });
  });
});
