import { ArrowRight, Plus } from 'lucide-react';
import { Button } from '#/components/ui/button';

export function HomeNextStep({
  onCreateProject,
  projectCount,
}: {
  onCreateProject: () => void;
  projectCount: number;
}) {
  if (projectCount > 0) {
    return null;
  }

  return (
    <section className="rounded-2xl border border-border bg-card p-5 shadow-sm">
      <div className="flex items-center gap-2">
        <ArrowRight className="h-4 w-4 text-primary" />
        <h2 className="font-semibold text-sm">Next step</h2>
      </div>
      <p className="mt-3 text-muted-foreground text-sm leading-6">
        Start with a project, then add the services and domains that belong to it.
      </p>
      <Button
        variant="ghost"
        className="mt-3 h-auto gap-2 px-0 text-primary hover:bg-transparent"
        onClick={onCreateProject}
      >
        <Plus className="h-4 w-4" />
        Create a project
      </Button>
    </section>
  );
}
