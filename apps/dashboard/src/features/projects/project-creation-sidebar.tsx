import { Server } from 'lucide-react';
import type { DeployPath } from './project-creation-paths';

export function ProjectCreationSidebar({ selectedPath }: { selectedPath?: DeployPath | null }) {
  return (
    <section className="rounded-2xl border border-border bg-card p-5 shadow-sm">
      <div className="flex items-center gap-2">
        <Server className="h-4 w-4 text-primary" />
        <h2 className="font-semibold text-sm">What happens next</h2>
      </div>
      <ol className="mt-5 space-y-4 text-sm">
        <li className="flex gap-3">
          <span className="font-semibold text-primary">01</span>
          <span className="text-muted-foreground">
            {selectedPath
              ? `Codedock opens the ${selectedPath.title.toLowerCase()} setup for this project.`
              : 'Codedock opens the empty project overview.'}
          </span>
        </li>
        <li className="flex gap-3">
          <span className="font-semibold text-primary">02</span>
          <span className="text-muted-foreground">
            Connect a source and configure build settings.
          </span>
        </li>
        <li className="flex gap-3">
          <span className="font-semibold text-primary">03</span>
          <span className="text-muted-foreground">
            Attach domains and monitor deployments from the project.
          </span>
        </li>
      </ol>
    </section>
  );
}
