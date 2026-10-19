import { describe, expect, it, vi } from 'vitest';
import { passkeyService } from './passkey';
import { apiClient } from '#/lib/api-client';

vi.mock('#/lib/api-client', () => ({
  apiClient: {
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    delete: vi.fn(),
  },
}));

describe('passkeyService', () => {
  it('unwraps publicKey envelope in getRegisterOptions if wrapped by server', async () => {
    const rawOptions = { challenge: 'test-challenge', rp: { name: 'Codedock' } };
    vi.mocked(apiClient.post).mockResolvedValueOnce({
      publicKey: rawOptions,
    });

    const result = await passkeyService.getRegisterOptions();
    expect(result).toEqual(rawOptions);
    expect(apiClient.post).toHaveBeenCalledWith('/auth/passkeys/register/options');
  });

  it('returns bare options in getRegisterOptions if server returns unwrapped object', async () => {
    const rawOptions = { challenge: 'test-challenge-2' };
    vi.mocked(apiClient.post).mockResolvedValueOnce(rawOptions);

    const result = await passkeyService.getRegisterOptions();
    expect(result).toEqual(rawOptions);
  });

  it('unwraps publicKey envelope in getAuthenticateOptions if wrapped', async () => {
    const authOptions = { challenge: 'auth-challenge' };
    vi.mocked(apiClient.post).mockResolvedValueOnce({
      publicKey: authOptions,
    });

    const result = await passkeyService.getAuthenticateOptions();
    expect(result).toEqual(authOptions);
    expect(apiClient.post).toHaveBeenCalledWith('/auth/passkeys/authenticate/options');
  });

  it('uses PUT method for renamePasskey matching the server endpoint', async () => {
    vi.mocked(apiClient.put).mockResolvedValueOnce({
      success: true,
      message: 'Passkey renamed successfully',
      data: null,
    });

    await passkeyService.renamePasskey('key-123', 'Office YubiKey');
    expect(apiClient.put).toHaveBeenCalledWith('/auth/passkeys/key-123', {
      name: 'Office YubiKey',
    });
  });
});
