import { Loader2 } from 'lucide-react';
import type { ReactNode } from 'react';
import { useEffect, useRef, useState } from 'react';
import { getApiBaseUrl } from '#/lib/api-client';
import { useSystemStore } from '#/stores/system-store';

const POLL_MS = 2000;
const HINT_AFTER_ATTEMPTS = 15;

async function probeBackend(): Promise<boolean> {
  try {
    const response = await fetch(`${getApiBaseUrl()}/system/public`, {
      credentials: 'include',
    });
    return response.ok;
  } catch {
    return false;
  }
}

export function BackendGate({ children }: { children: ReactNode }) {
  const [ready, setReady] = useState(false);
  const [attempts, setAttempts] = useState(0);
  const siteName = useSystemStore((state) => state.siteName);
  const timer = useRef<number | undefined>(undefined);

  useEffect(() => {
    let cancelled = false;
    const ping = async () => {
      if (await probeBackend()) {
        if (!cancelled) setReady(true);
        return;
      }
      if (!cancelled) {
        setAttempts((count) => count + 1);
        timer.current = window.setTimeout(ping, POLL_MS);
      }
    };
    ping();
    return () => {
      cancelled = true;
      window.clearTimeout(timer.current);
    };
  }, []);

  if (ready) return <>{children}</>;

  return (
    <main className="grid min-h-dvh place-items-center bg-background px-6 text-center">
      <div className="w-full max-w-sm">
        <div className="mb-5 flex items-center justify-center gap-2.5">
          <img src="/apple-touch-icon.png" alt="" className="size-8 rounded-lg" />
          <span className="font-bold text-foreground text-lg tracking-[-0.04em]">{siteName}</span>
        </div>
        <div className="flex items-center justify-center gap-3 text-muted-foreground text-sm">
          <Loader2 className="size-4 animate-spin text-primary" />
          Starting {siteName}…
        </div>
        <p className="mt-3 text-muted-foreground text-xs leading-5">
          {attempts >= HINT_AFTER_ATTEMPTS
            ? 'Still starting — first boot pulls the database image. Check the daemon terminal if this persists.'
            : 'Waiting for the daemon to respond.'}
        </p>
      </div>
    </main>
  );
}
