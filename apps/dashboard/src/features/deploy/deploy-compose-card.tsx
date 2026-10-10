import { ChevronDown, ChevronUp, Layers, Scan } from 'lucide-react';
import { useState } from 'react';
import { Button } from '#/components/ui/button';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';

interface DeployComposeCardProps {
  composePath?: string;
  onComposePathChange?: (p: string) => void;
}

export function DeployComposeCard({
  composePath = '',
  onComposePathChange,
}: DeployComposeCardProps) {
  const [open, setOpen] = useState(false);
  const [path, setPath] = useState(composePath);

  const handleChange = (val: string) => {
    setPath(val);
    onComposePathChange?.(val);
  };

  return (
    <section className="rounded-2xl border border-border/60 bg-card p-6">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3.5">
          <div className="flex size-10 shrink-0 items-center justify-center rounded-xl bg-muted text-muted-foreground">
            <Layers className="size-5" />
          </div>
          <div>
            <h2 className="font-semibold text-foreground text-sm">Compose file</h2>
            <p className="text-muted-foreground text-xs">
              Deploying with Docker Compose? Point to the file.
            </p>
          </div>
        </div>

        <button
          type="button"
          onClick={() => setOpen((prev) => !prev)}
          className="p-1 text-muted-foreground transition-colors hover:text-foreground"
          aria-label={open ? 'Collapse compose file section' : 'Expand compose file section'}
        >
          {open ? <ChevronUp className="size-4" /> : <ChevronDown className="size-4" />}
        </button>
      </div>

      {open && (
        <div className="mt-5 space-y-3 border-border/50 border-t pt-5">
          <Label htmlFor="compose-path" className="text-muted-foreground text-xs">
            Path to your compose file <span className="text-muted-foreground/60">(Optional)</span>
          </Label>
          <div className="flex items-center gap-2">
            <Input
              id="compose-path"
              value={path}
              onChange={(e) => handleChange(e.target.value)}
              placeholder="deploy/docker-compose/docker-compose.yml"
              className="h-10 rounded-xl border-border/60 bg-background/50 font-mono text-xs"
            />
            <Button
              type="button"
              variant="outline"
              size="sm"
              className="h-10 gap-1.5 rounded-xl px-4 text-xs"
              onClick={() => {}}
            >
              <Scan className="size-3.5" />
              Scan
            </Button>
          </div>
          <p className="text-[11px] text-muted-foreground leading-relaxed">
            For repos that keep compose outside the root — give the file (deploy/stack.yml) or the
            folder holding it (deploy/docker-compose). Re-scans the repo and deploys the services it
            declares.
          </p>
        </div>
      )}
    </section>
  );
}
