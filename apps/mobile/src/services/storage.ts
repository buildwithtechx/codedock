import * as SecureStore from 'expo-secure-store';
import type { CodedockConnection } from '@/interfaces';

const KEY_URL = 'codedock_instance_url';
const KEY_TOKEN = 'codedock_pat_token';
const KEY_USER = 'codedock_user_email';
const KEY_MODE = 'codedock_connection_mode';

export const MobileStorage = {
  async getConnection(): Promise<CodedockConnection | null> {
    try {
      const [url, token, userEmail, mode] = await Promise.all([
        SecureStore.getItemAsync(KEY_URL),
        SecureStore.getItemAsync(KEY_TOKEN),
        SecureStore.getItemAsync(KEY_USER),
        SecureStore.getItemAsync(KEY_MODE),
      ]);

      if (!url || !token) {
        return null;
      }

      return {
        url,
        token,
        mode:
          (mode as 'cloud' | 'self-hosted') || (url.includes('codedock') ? 'cloud' : 'self-hosted'),
        userEmail: userEmail || undefined,
      };
    } catch {
      return null;
    }
  },

  async saveConnection(
    url: string,
    token: string,
    userEmail?: string,
    mode: 'cloud' | 'self-hosted' = 'self-hosted'
  ): Promise<void> {
    await Promise.all([
      SecureStore.setItemAsync(KEY_URL, url),
      SecureStore.setItemAsync(KEY_TOKEN, token),
      userEmail ? SecureStore.setItemAsync(KEY_USER, userEmail) : Promise.resolve(),
      SecureStore.setItemAsync(KEY_MODE, mode),
    ]);
  },

  async clearConnection(): Promise<void> {
    await Promise.all([
      SecureStore.deleteItemAsync(KEY_URL),
      SecureStore.deleteItemAsync(KEY_TOKEN),
      SecureStore.deleteItemAsync(KEY_USER),
      SecureStore.deleteItemAsync(KEY_MODE),
    ]);
  },
};
