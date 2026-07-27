import { createEnv } from "@t3-oss/env-core";
import { z } from "zod";

export const env = createEnv({
  server: {
    SERVER_URL: z.string().url().optional(),
  },

  clientPrefix: "VITE_",

  client: {
    VITE_SITE_NAME: z.string().min(1).default("Codedock"),
    VITE_SITE_URL: z.string().url().default("https://codedock.run"),
    VITE_PUBLIC_POSTHOG_KEY: z.string().optional(),
    VITE_PUBLIC_POSTHOG_HOST: z.string().default("https://us.i.posthog.com"),
    VITE_PUBLIC_POSTHOG_DEFAULTS: z.string().default("2026-01-30"),
  },

  runtimeEnv: import.meta.env,
  emptyStringAsUndefined: true,
});
