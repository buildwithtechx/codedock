import { Code2, KeyRound, RefreshCw } from 'lucide-react';
import { useEffect, useState } from 'react';
import { Button } from '#/components/ui/button';
import type { ManagedRuntimeCredential, ManagedRuntimeStatus } from '#/interfaces/server';
import { ConnectionValue } from './connection-value';
import {
  enableManagedRuntime,
  getManagedRuntimeCredential,
  getManagedRuntimeStatus,
  rotateManagedRuntimeCredential,
} from './managed-api';
import { ManagedControlPanel } from './managed-control-panel';

export function ManagedServerRuntime({ serverId }: { serverId: string }) {
  const [status, setStatus] = useState<ManagedRuntimeStatus | null>(null);
  const [credential, setCredential] = useState<ManagedRuntimeCredential | null>(null);
  const [loading, setLoading] = useState(true);
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [confirmRotate, setConfirmRotate] = useState(false);

  const fetchStatus = async () => {
    setLoading(true);
    setError(null);
    try {
      const data = await getManagedRuntimeStatus(serverId);
      setStatus(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load runtime status');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    void fetchStatus();
  }, [serverId]);

  const handleEnable = async () => {
    setBusy(true);
    setError(null);
    try {
      const updated = await enableManagedRuntime(serverId);
      setStatus(updated);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to enable runtime API');
    } finally {
      setBusy(false);
    }
  };

  const handleReveal = async () => {
    setBusy(true);
    setError(null);
    try {
      const cred = await getManagedRuntimeCredential(serverId);
      setCredential(cred);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to generate runtime credentials');
    } finally {
      setBusy(false);
    }
  };

  const handleRotate = async () => {
    setBusy(true);
    setError(null);
    try {
      const cred = await rotateManagedRuntimeCredential(serverId);
      setCredential(cred);
      setConfirmRotate(false);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to rotate runtime token');
    } finally {
      setBusy(false);
    }
  };

  return (
    <ManagedControlPanel
      title="Runtime API"
      description="Direct runner API endpoints for Codedock CLI tooling, container inspection, and CI pipelines."
      icon={Code2}
      busy={loading || busy}
      error={error}
      refresh={() => void fetchStatus()}
    >
      <div className="flex flex-wrap items-center justify-between gap-3 rounded-xl bg-muted/30 p-4">
        <div>
          <span className="font-medium text-foreground text-sm">
            {status?.enabled ? 'Runtime API is Active' : 'Runtime API is Disabled'}
          </span>
          <p className="mt-0.5 text-muted-foreground text-xs">
            {status?.enabled
              ? 'Runner listens for authorized control-plane build dispatch and status probes.'
              : 'Runner API endpoints are disabled.'}
          </p>
        </div>
        {!status?.enabled && (
          <Button size="sm" disabled={busy} onClick={() => void handleEnable()}>
            Enable Runtime API
          </Button>
        )}
      </div>

      <p className="rounded-xl bg-muted/20 p-4 text-muted-foreground text-xs leading-relaxed">
        The Runtime API allows orchestrating daemon tasks directly on this server using
        cryptographic tokens.
      </p>

      {status?.enabled && !credential && (
        <div className="space-y-3">
          <Button variant="secondary" size="sm" disabled={busy} onClick={() => void handleReveal()}>
            <KeyRound className="mr-2 size-3.5" />
            Reveal runtime credentials
          </Button>
        </div>
      )}

      {credential && (
        <div className="space-y-4 rounded-xl bg-muted/20 p-4">
          <div className="flex items-center justify-between">
            <h3 className="font-medium text-foreground text-sm">Authentication Credentials</h3>
            <Button
              variant="ghost"
              size="sm"
              className="h-7 text-xs"
              onClick={() => setCredential(null)}
            >
              Hide credentials
            </Button>
          </div>

          <ConnectionValue label="Runtime Endpoint URL" value={credential.endpoint} />
          <ConnectionValue label="Bearer Runner Token" value={credential.token} secret />

          {credential.expiresAt && (
            <p className="text-muted-foreground text-xs">
              Token expires on {new Date(credential.expiresAt).toLocaleString()}
            </p>
          )}

          <div className="pt-2">
            {!confirmRotate ? (
              <Button
                variant="secondary"
                size="sm"
                disabled={busy}
                onClick={() => setConfirmRotate(true)}
              >
                <RefreshCw className="mr-2 size-3.5" />
                Rotate token
              </Button>
            ) : (
              <div className="space-y-3 rounded-lg bg-warning/10 p-3">
                <p className="text-foreground text-xs">
                  Rotating this token immediately invalidates any active CLI or agent sessions.
                </p>
                <div className="flex gap-2">
                  <Button
                    variant="destructive"
                    size="sm"
                    disabled={busy}
                    onClick={() => void handleRotate()}
                  >
                    Confirm Rotate
                  </Button>
                  <Button
                    variant="ghost"
                    size="sm"
                    disabled={busy}
                    onClick={() => setConfirmRotate(false)}
                  >
                    Cancel
                  </Button>
                </div>
              </div>
            )}
          </div>
        </div>
      )}
    </ManagedControlPanel>
  );
}
