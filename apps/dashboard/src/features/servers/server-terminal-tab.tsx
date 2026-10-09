import { TerminalSquare } from 'lucide-react';
import { BlurIp } from '#/components/ui/blur-ip';
import {
  Empty,
  EmptyDescription,
  EmptyHeader,
  EmptyMedia,
  EmptyTitle,
} from '#/components/ui/empty';
import type { Server } from '#/interfaces/server';

export function ServerTerminalTab({ server }: { server: Server }) {
  const host = server.sshHost ?? server.ipAddress;

  return (
    <div className="overflow-hidden rounded-2xl bg-card">
      <div className="flex items-center gap-3 border-border/50 border-b px-5 py-4">
        <div className="flex size-9 shrink-0 items-center justify-center rounded-xl bg-muted">
          <TerminalSquare className="size-[18px] text-muted-foreground" />
        </div>
        <div className="min-w-0 flex-1">
          <h2 className="truncate font-semibold text-[15px] text-foreground">
            Terminal <span className="text-muted-foreground">· {server.name}</span>
          </h2>
          <p className="text-muted-foreground text-xs">Interactive shell on this node</p>
        </div>
      </div>
      <div className="p-5">
        <Empty>
          <EmptyHeader>
            <EmptyMedia variant="icon">
              <TerminalSquare />
            </EmptyMedia>
            <EmptyTitle>Server shell unavailable</EmptyTitle>
            <EmptyDescription>
              An interactive shell for{' '}
              {server.isLocal ? (
                'local shell'
              ) : (
                <>
                  {server.sshUser ?? 'root'}@<BlurIp>{host}</BlurIp>
                </>
              )}{' '}
              is not connected yet. Use your SSH client to access this node.
            </EmptyDescription>
          </EmptyHeader>
        </Empty>
      </div>
    </div>
  );
}
