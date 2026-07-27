import { useEffect } from "react";

declare global {
  interface Window {
    __codedock_posthog_initialized?: boolean;
    posthog?: {
      init: (key: string, config: Record<string, unknown>) => void;
    };
  }
}

export function PosthogTracking() {
  const posthogKey = import.meta.env.VITE_PUBLIC_POSTHOG_KEY;
  const posthogHost = import.meta.env.VITE_PUBLIC_POSTHOG_HOST || "https://us.i.posthog.com";
  const posthogDefaults = import.meta.env.VITE_PUBLIC_POSTHOG_DEFAULTS || "2026-01-30";

  useEffect(() => {
    if (!posthogKey || window.__codedock_posthog_initialized) return;

    window.__codedock_posthog_initialized = true;

    const script = document.createElement("script");
    script.type = "text/javascript";
    script.async = true;
    script.crossOrigin = "anonymous";
    script.src = `${posthogHost.replace(".i.posthog.com", "-assets.i.posthog.com")}/static/array.js`;

    script.onload = () => {
      if (window.posthog) {
        window.posthog.init(posthogKey, {
          api_host: posthogHost,
          defaults: posthogDefaults,
          autocapture: true,
          capture_pageview: true,
          capture_pageleave: true,
        });
      }
    };

    document.head.appendChild(script);
  }, []);

  return null;
}
