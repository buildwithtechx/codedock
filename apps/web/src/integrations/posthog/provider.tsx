import { PostHogProvider as BasePostHogProvider } from '@posthog/react';
import posthog from 'posthog-js';
import type { ReactNode } from 'react';
import { env } from '../../env';

if (typeof window !== 'undefined' && env.VITE_PUBLIC_POSTHOG_KEY) {
  posthog.init(env.VITE_PUBLIC_POSTHOG_KEY, {
    api_host: env.VITE_PUBLIC_POSTHOG_HOST,
    person_profiles: 'identified_only',
    capture_pageview: false,
    defaults: '2025-11-30',
  });
}

interface PostHogProviderProps {
  children: ReactNode;
}

export function PostHogProvider({ children }: PostHogProviderProps) {
  return <BasePostHogProvider client={posthog}>{children}</BasePostHogProvider>;
}
