import { Link } from '@tanstack/react-router';
import { RotateCw, Settings, WifiOff } from 'lucide-react';
import { Button } from '#/components/ui/button';
import type { Server } from '#/interfaces/server';
import { reachabilityOf } from './server-list-row';

interface ServerConnectionBannerProps {
  server: Server;
  retrying: boolean;
  onRetry: () => void;
}

export function ServerConnectionBanner({ server, retrying, onRetry }: ServerConnectionBannerProps) {
  const reach = reachabilityOf(server);
  if (reach === 'online') return null;

  const host = server.sshHost || server.ipAddress;
  const offline = reach === 'offline';

  return (
    <div
      role="alert"
      className={
        offline
          ? 'mb-6 rounded-2xl border border-destructive/30 bg-destructive/5 p-4'
          : 'mb-6 rounded-2xl border border-warning/30 bg-warning/5 p-4'
      }
    >
      <div className="flex items-start gap-3">
        <div
          className={
            offline
              ? 'flex size-9 shrink-0 items-center justify-center rounded-lg bg-destructive/10 text-destructive'
              : 'flex size-9 shrink-0 items-center justify-center rounded-lg bg-warning/10 text-warning'
          }
        >
          <WifiOff className="size-4" />
        </div>
        <div className="min-w-0 flex-1">
          <p className="font-semibold text-foreground text-sm">
            {offline
              ? `Cannot reach ${server.name}`
              : `Connection status for ${server.name} is unknown`}
          </p>
          <p className="mt-1 text-[13px] text-muted-foreground leading-relaxed">
            {offline
              ? `The last heartbeat from ${host} failed. Telemetry and deploys to this node are paused until it comes back.`
              : 'This server has not reported a heartbeat yet. Telemetry appears once the node checks in.'}
          </p>
          {offline && !server.isLocal && (
            <ul className="mt-2 list-disc space-y-0.5 ps-5 text-muted-foreground/80 text-xs">
              <li>Check that the host is powered on</li>
              <li>
                <code className="font-mono">ping {host}</code> to verify network reachability
              </li>
              <li>
                <code className="font-mono">
                  ssh {server.sshUser ?? 'root'}@{host} -p {server.sshPort ?? 22}
                </code>{' '}
                to verify SSH access
              </li>
            </ul>
          )}
          <div className="mt-3 flex flex-wrap items-center gap-2">
            <Button size="sm" variant="secondary" onClick={onRetry} disabled={retrying}>
              <RotateCw className={`size-3 ${retrying ? 'animate-spin' : ''}`} />
              {retrying ? 'Checking…' : 'Retry'}
            </Button>
            <Button size="sm" variant="secondary" asChild>
              <Link
                to="/servers/$serverId"
                params={{ serverId: server.id }}
                search={{ tab: 'settings' }}
              >
                <Settings className="size-3" />
                Edit server
              </Link>
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
}
