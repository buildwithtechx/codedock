import { Gauge, ScanSearch } from 'lucide-react';
import {
  Empty,
  EmptyDescription,
  EmptyHeader,
  EmptyMedia,
  EmptyTitle,
} from '#/components/ui/empty';

export function ServerSecurityTab() {
  return (
    <div className="space-y-6">
      <div className="rounded-2xl bg-card">
        <div className="flex flex-wrap items-center justify-between gap-3 border-border/50 border-b px-5 py-4">
          <div className="flex min-w-0 items-center gap-3">
            <div className="flex h-9 w-9 items-center justify-center rounded-xl bg-primary/10">
              <ScanSearch className="size-[18px] text-primary" />
            </div>
            <div>
              <h2 className="font-semibold text-[15px] text-foreground">Exposed ports</h2>
              <p className="text-muted-foreground text-xs">
                Listeners reachable from outside this node
              </p>
            </div>
          </div>
        </div>
        <div className="p-5">
          <Empty>
            <EmptyHeader>
              <EmptyMedia variant="icon">
                <ScanSearch />
              </EmptyMedia>
              <EmptyTitle>Port scans are not available yet</EmptyTitle>
              <EmptyDescription>
                On-demand exposure scans of this node are planned. Until then, audit listening ports
                directly over SSH.
              </EmptyDescription>
            </EmptyHeader>
          </Empty>
        </div>
      </div>

      <div className="rounded-2xl bg-card">
        <div className="flex flex-wrap items-center justify-between gap-3 border-border/50 border-b px-5 py-4">
          <div className="flex min-w-0 items-center gap-3">
            <div className="flex h-9 w-9 items-center justify-center rounded-xl bg-warning/10">
              <Gauge className="size-[18px] text-warning" />
            </div>
            <div>
              <h2 className="font-semibold text-[15px] text-foreground">Rate limiting</h2>
              <p className="text-muted-foreground text-xs">Edge request policy for this node</p>
            </div>
          </div>
        </div>
        <div className="p-5">
          <Empty>
            <EmptyHeader>
              <EmptyMedia variant="icon">
                <Gauge />
              </EmptyMedia>
              <EmptyTitle>Rate limiting is not available yet</EmptyTitle>
              <EmptyDescription>
                Per-node request policies are planned. Edge defaults currently apply to every route
                on this server.
              </EmptyDescription>
            </EmptyHeader>
          </Empty>
        </div>
      </div>
    </div>
  );
}
