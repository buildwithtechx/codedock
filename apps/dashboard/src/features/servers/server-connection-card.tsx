import { Globe, Hash, KeyRound, Network, Server as ServerIcon, User } from 'lucide-react';
import type { ReactNode } from 'react';
import { BlurIp } from '#/components/ui/blur-ip';
import type { Server } from '#/interfaces/server';

function ConnectionRow({
  icon,
  label,
  children,
}: {
  icon: ReactNode;
  label: string;
  children: ReactNode;
}) {
  return (
    <div className="flex items-center justify-between">
      <div className="flex items-center gap-2">
        <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-muted/60">
          {icon}
        </div>
        <span className="text-muted-foreground text-sm">{label}</span>
      </div>
      {children}
    </div>
  );
}

export function ServerConnectionCard({ server }: { server: Server }) {
  const host = server.sshHost || server.ipAddress;

  return (
    <div className="rounded-2xl border border-border/50 bg-card p-5">
      <div className="mb-4 flex items-center gap-2">
        <ServerIcon className="size-4 text-muted-foreground" />
        <h3 className="font-semibold text-foreground text-sm">Connection</h3>
      </div>
      <div className="space-y-3">
        <ConnectionRow icon={<Globe className="size-4 text-muted-foreground" />} label="Host">
          <span className="max-w-[150px] truncate text-right font-medium font-mono text-foreground text-sm">
            {server.isLocal ? 'localhost' : <BlurIp>{host}</BlurIp>}
          </span>
        </ConnectionRow>
        {!server.isLocal && (
          <ConnectionRow icon={<Hash className="size-4 text-muted-foreground" />} label="Port">
            <span className="font-medium font-mono text-foreground text-sm">
              {server.sshPort ?? 22}
            </span>
          </ConnectionRow>
        )}
        {!server.isLocal && (
          <ConnectionRow icon={<User className="size-4 text-muted-foreground" />} label="User">
            <span className="font-medium font-mono text-foreground text-sm">
              {server.sshUser ?? 'root'}
            </span>
          </ConnectionRow>
        )}

        <div className="my-2 h-px bg-border/60" />

        <ConnectionRow icon={<KeyRound className="size-4 text-muted-foreground" />} label="Auth">
          <span className="font-medium text-foreground text-sm">
            {server.isLocal
              ? 'Local socket'
              : server.sshAuthMethod === 'password'
                ? 'Password'
                : server.sshAuthMethod === 'agent'
                  ? 'SSH agent'
                  : 'SSH key'}
          </span>
        </ConnectionRow>
        {!server.isLocal && (
          <ConnectionRow
            icon={<Network className="size-4 text-muted-foreground" />}
            label="Transport"
          >
            <span className="font-medium text-foreground text-sm capitalize">
              {server.sshTransport || 'direct'}
            </span>
          </ConnectionRow>
        )}
      </div>
    </div>
  );
}
