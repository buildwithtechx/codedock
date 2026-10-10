import { Cloud, Pencil, Server } from 'lucide-react';
import { Button } from '#/components/ui/button';

interface DeployDestinationCardProps {
  runtimeMode: 'web' | 'worker' | 'static';
  onEdit: () => void;
}

export function DeployDestinationCard({ runtimeMode, onEdit }: DeployDestinationCardProps) {
  const isStatic = runtimeMode === 'static';
  const targetLabel = isStatic
    ? 'Cloud Pages (default)'
    : runtimeMode === 'web'
      ? 'Production Server'
      : 'Background Worker';

  const description = isStatic
    ? 'Files published to the edge · Zero running containers'
    : 'Dedicated container process · Allocated server resources';

  return (
    <section className="rounded-2xl border border-border/60 bg-card p-5">
      <div className="flex items-center justify-between gap-4">
        <div className="flex items-center gap-3.5">
          <div className="grid size-10 place-items-center rounded-xl bg-primary/10 text-primary">
            {isStatic ? <Cloud className="size-5" /> : <Server className="size-5" />}
          </div>
          <div>
            <h3 className="font-semibold text-foreground text-sm">Where do you want to deploy?</h3>
            <p className="text-muted-foreground text-xs">Choose where your application will run</p>
          </div>
        </div>
        <Button
          type="button"
          variant="outline"
          size="sm"
          onClick={onEdit}
          className="gap-1.5 text-xs"
        >
          <Pencil className="size-3.5" />
          Edit
        </Button>
      </div>

      <div className="mt-4 flex flex-col gap-2 rounded-xl border border-border/40 bg-muted/20 px-3.5 py-2.5 sm:flex-row sm:items-center sm:justify-between">
        <div className="flex items-center gap-2.5 text-xs">
          <span className="font-semibold text-foreground">{targetLabel}</span>
          <span className="text-muted-foreground/60">·</span>
          <span className="text-muted-foreground">{description}</span>
        </div>
        <span className="inline-flex w-fit items-center rounded-full bg-emerald-500/10 px-2 py-0.5 font-medium text-[10px] text-emerald-600 dark:text-emerald-400">
          Active target
        </span>
      </div>
    </section>
  );
}
