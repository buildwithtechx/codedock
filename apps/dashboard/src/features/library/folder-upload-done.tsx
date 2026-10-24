import { Link } from '@tanstack/react-router';
import { Package } from 'lucide-react';
import { Button } from '@codedock/ui/button';

interface FolderUploadDoneProps {
  name: string;
  projectId: string;
  onReset: () => void;
}

export function FolderUploadDone({ name, projectId, onReset }: FolderUploadDoneProps) {
  return (
    <div className="rounded-2xl border border-border/50 bg-card p-8 text-center">
      <div className="mx-auto mb-4 flex h-14 w-14 items-center justify-center rounded-2xl bg-emerald-500/10">
        <Package className="size-7 text-emerald-500" />
      </div>
      <h3 className="font-semibold text-base text-foreground">Archive deployed</h3>
      <p className="mx-auto mt-1 max-w-sm text-muted-foreground text-sm">
        {name} is deploying in your project.
      </p>
      <div className="mt-6 flex flex-wrap items-center justify-center gap-3">
        <Button asChild>
          <Link to="/projects/$projectId" params={{ projectId }}>
            Open project
          </Link>
        </Button>
        <Button variant="outline" onClick={onReset}>
          Upload another
        </Button>
      </div>
    </div>
  );
}
