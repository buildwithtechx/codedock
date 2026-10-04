import { Check, Copy, Cpu, RefreshCw, ServerIcon, Shield } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import { Card } from '#/components/ui/card';
import { useTestSSH } from '#/hooks/use-servers';
import type { Server } from '#/interfaces/server';

interface ServerConnectionCardProps {
  server: Server;
  onRefresh?: () => void;
}

export function ServerConnectionCard({ server, onRefresh }: ServerConnectionCardProps) {
  const [copied, setCopied] = useState(false);
  const { mutateAsync: testSSH, isPending: isTesting } = useTestSSH();

  const copyToken = async () => {
    if (!server.workerToken) return;
    await navigator.clipboard.writeText(server.workerToken);
    setCopied(true);
    toast.success('Worker token copied to clipboard');
    setTimeout(() => setCopied(false), 2000);
  };

  const handleTestConnection = async () => {
    if (server.isLocal) {
      toast.success('Local Docker socket is active');
      return;
    }
    try {
      const res = await testSSH({
        sshHost: server.sshHost || server.ipAddress,
        sshPort: server.sshPort || 22,
        sshUser: server.sshUser || 'root',
      });
      if (res.success) {
        toast.success('SSH connection verified successfully');
      } else {
        toast.error(`SSH check failed: ${res.message}`);
      }
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'SSH connection failed');
    }
  };

  return (
    <Card className="p-5">
      <div className="flex flex-wrap items-center justify-between gap-3">
        <div className="flex items-center gap-3">
          <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-primary/10 text-primary">
            {server.isLocal ? <Cpu className="h-5 w-5" /> : <ServerIcon className="h-5 w-5" />}
          </div>
          <div>
            <h3 className="font-semibold text-sm">{server.name}</h3>
            <p className="font-mono text-muted-foreground text-xs">{server.ipAddress}</p>
          </div>
        </div>

        <div className="flex items-center gap-2">
          {!server.isLocal && (
            <Button
              variant="outline"
              size="sm"
              onClick={handleTestConnection}
              disabled={isTesting}
              className="h-8 gap-1.5 text-xs"
            >
              <RefreshCw className={`h-3.5 w-3.5 ${isTesting ? 'animate-spin' : ''}`} />
              Test Connection
            </Button>
          )}
          {onRefresh && (
            <Button variant="ghost" size="sm" onClick={onRefresh} className="h-8 text-xs">
              Refresh
            </Button>
          )}
        </div>
      </div>

      <div className="mt-5 grid grid-cols-2 gap-4 border-border/60 border-t pt-4 sm:grid-cols-4">
        <div>
          <span className="text-muted-foreground text-xs">Node Type</span>
          <p className="mt-0.5 font-medium text-xs">
            {server.isLocal ? 'Local Daemon' : 'Remote SSH'}
          </p>
        </div>
        <div>
          <span className="text-muted-foreground text-xs">SSH Endpoint</span>
          <p className="mt-0.5 font-mono text-xs">
            {server.isLocal
              ? 'Local socket'
              : `${server.sshUser || 'root'}@${server.sshHost || server.ipAddress}:${server.sshPort || 22}`}
          </p>
        </div>
        <div>
          <span className="text-muted-foreground text-xs">Transport</span>
          <p className="mt-0.5 font-medium text-xs capitalize">{server.sshTransport || 'direct'}</p>
        </div>
        <div>
          <span className="text-muted-foreground text-xs">Auth Method</span>
          <p className="mt-0.5 font-medium text-xs capitalize">{server.sshAuthMethod || 'key'}</p>
        </div>
      </div>

      {server.workerToken && (
        <div className="mt-4 flex items-center justify-between rounded-xl border border-border/70 bg-muted/20 px-3.5 py-2.5">
          <div className="flex items-center gap-2">
            <Shield className="h-3.5 w-3.5 text-muted-foreground" />
            <span className="text-muted-foreground text-xs">Worker Token:</span>
            <span className="font-mono text-xs">••••••••••••••••</span>
          </div>
          <Button variant="ghost" size="icon" className="h-7 w-7" onClick={copyToken}>
            {copied ? (
              <Check className="h-3.5 w-3.5 text-emerald-500" />
            ) : (
              <Copy className="h-3.5 w-3.5" />
            )}
          </Button>
        </div>
      )}
    </Card>
  );
}
