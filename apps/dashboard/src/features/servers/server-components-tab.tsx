import { Container, ShieldCheck } from 'lucide-react';
import {
  Empty,
  EmptyDescription,
  EmptyHeader,
  EmptyMedia,
  EmptyTitle,
} from '#/components/ui/empty';

export function ServerComponentsTab() {
  return (
    <div className="space-y-6">
      <div className="rounded-2xl bg-card">
        <div className="flex flex-wrap items-center gap-3 border-border/50 border-b px-5 py-4">
          <div className="flex h-9 w-9 items-center justify-center rounded-xl bg-success/10">
            <ShieldCheck className="size-[18px] text-success" />
          </div>
          <div className="min-w-0 flex-1">
            <h2 className="font-semibold text-[15px] text-foreground">System health</h2>
            <p className="text-muted-foreground text-xs">Required components on this node</p>
          </div>
        </div>
        <div className="p-5">
          <Empty>
            <EmptyHeader>
              <EmptyMedia variant="icon">
                <Container />
              </EmptyMedia>
              <EmptyTitle>Component checks are not available yet</EmptyTitle>
              <EmptyDescription>
                Health checks for Docker, Git, and the reverse proxy are planned. Telemetry on the
                Overview tab already reports live CPU, memory, and disk for this node.
              </EmptyDescription>
            </EmptyHeader>
          </Empty>
        </div>
      </div>
    </div>
  );
}
