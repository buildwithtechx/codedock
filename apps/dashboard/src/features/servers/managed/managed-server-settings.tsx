import { Link } from '@tanstack/react-router';
import { FileText, Settings } from 'lucide-react';
import { useEffect, useState } from 'react';
import { Button } from '#/components/ui/button';
import type { ManagedServerInfo } from '#/interfaces/server';
import { CapacitySummary } from './capacity-summary';
import { ConnectionValue } from './connection-value';
import { getManagedBootLogs, getManagedInfo } from './managed-api';
import { ManagedControlPanel } from './managed-control-panel';

export function ManagedServerSettings({ serverId }: { serverId: string }) {
  const [info, setInfo] = useState<ManagedServerInfo | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const [bootLogs, setBootLogs] = useState<string | null>(null);
  const [logsLoading, setLogsLoading] = useState(false);

  const fetchInfo = async () => {
    setLoading(true);
    setError(null);
    try {
      const data = await getManagedInfo(serverId);
      setInfo(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load server settings');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    void fetchInfo();
  }, [serverId]);

  const handleFetchBootLogs = async () => {
    setLogsLoading(true);
    try {
      const res = await getManagedBootLogs(serverId, 200);
      setBootLogs(res.logs);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch boot logs');
    } finally {
      setLogsLoading(false);
    }
  };

  const resources = info?.resources ?? {
    cpuCores: 4,
    memoryMb: 8192,
    diskMb: 80000,
  };

  return (
    <ManagedControlPanel
      title="Server Settings"
      description="Specifications, hardware virtualization limits, and host boot telemetry."
      icon={Settings}
      busy={loading}
      error={error}
      refresh={() => void fetchInfo()}
    >
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
        <ConnectionValue label="Base Runner Image" value={info?.image ?? 'codedock-runner:v1.4'} />
        <ConnectionValue
          label="Instance Identifier"
          value={info?.workspaceId ?? `ws-${serverId}`}
        />
        <div className="rounded-xl bg-muted/40 p-3">
          <p className="font-medium text-muted-foreground text-xs">Operating System</p>
          <p className="mt-1 font-medium text-foreground text-sm">
            {info?.operatingSystem ?? 'Debian GNU/Linux 12 (bookworm)'}
          </p>
        </div>
        <div className="rounded-xl bg-muted/40 p-3">
          <p className="font-medium text-muted-foreground text-xs">Storage Lifecycle</p>
          <p className="mt-1 font-medium text-foreground text-sm">
            {info?.mode === 'permanent' ? 'Persistent SSD' : 'Persistent Storage'}
          </p>
        </div>
      </div>

      <div className="space-y-3">
        <h3 className="font-medium text-foreground text-sm">Allocated Capacity</h3>
        <CapacitySummary resources={resources} />
        <p className="text-muted-foreground text-xs">
          Hardware cgroup resource limits provisioned for container runtimes on this node.
        </p>
        <Button variant="secondary" size="sm" asChild>
          <Link to="/billing">Resize capacity</Link>
        </Button>
      </div>

      <div className="space-y-3 pt-2">
        <div className="flex items-center justify-between">
          <h3 className="font-medium text-foreground text-sm">System Diagnostics</h3>
          <Button
            variant="outline"
            size="sm"
            disabled={logsLoading}
            onClick={() => void handleFetchBootLogs()}
          >
            <FileText className="mr-2 size-3.5" />
            {logsLoading ? 'Reading boot logs...' : 'View boot logs'}
          </Button>
        </div>

        {bootLogs && (
          <div className="space-y-2 rounded-xl bg-muted/20 p-4">
            <div className="flex items-center justify-between">
              <span className="font-medium text-foreground text-xs">
                journalctl -b (recent boot)
              </span>
              <Button
                variant="ghost"
                size="sm"
                className="h-7 text-xs"
                onClick={() => setBootLogs(null)}
              >
                Close
              </Button>
            </div>
            <pre className="max-h-80 overflow-auto whitespace-pre-wrap rounded-lg bg-background p-3 font-mono text-muted-foreground text-xs leading-relaxed">
              {bootLogs}
            </pre>
          </div>
        )}
      </div>
    </ManagedControlPanel>
  );
}
