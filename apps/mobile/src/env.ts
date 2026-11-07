import { z } from 'zod';

const envSchema = z.object({
  EXPO_PUBLIC_CODEDOCK_CLOUD_URL: z.string().default('https://app.codedock.run'),
  EXPO_PUBLIC_APP_SCHEME: z.string().default('codedock'),
});

export const env = envSchema.parse({
  EXPO_PUBLIC_CODEDOCK_CLOUD_URL: process.env.EXPO_PUBLIC_CODEDOCK_CLOUD_URL,
  EXPO_PUBLIC_APP_SCHEME: process.env.EXPO_PUBLIC_APP_SCHEME,
});
