import { Link } from '@tanstack/react-router';
import { FolderKanban, Loader2, Plus } from 'lucide-react';
import { PageHeader } from '#/components/layout/page-header';
import { Button } from '#/components/ui/button';
import { ProjectCard } from '#/features/projects/project-card';
import { useListCanvasSummaries } from '#/hooks/use-canvas';

export function ProjectDirectory() {
  const { data, isLoading } = useListCanvasSummaries();
  const projects = data?.data || [];

  return (
    <div className="space-y-6">
      <PageHeader
        title="Projects"
        description="Projects group the services and environments your team runs."
        action={
          <Link to="/projects/new">
            <Button className="gap-2">
              <Plus className="h-4 w-4" />
              New project
            </Button>
          </Link>
        }
      />

      {isLoading ? (
        <div className="flex min-h-80 items-center justify-center rounded-2xl border border-border/80 bg-card">
          <Loader2 className="h-5 w-5 animate-spin text-muted-foreground" />
        </div>
      ) : projects.length === 0 ? (
        <div className="flex min-h-80 flex-col items-center justify-center rounded-2xl border border-border border-dashed bg-card px-6 text-center">
          <div className="flex h-12 w-12 items-center justify-center rounded-2xl bg-primary/10 text-primary">
            <FolderKanban className="h-5 w-5" />
          </div>
          <h2 className="mt-4 font-semibold text-lg">No projects in this organization</h2>
          <p className="mt-1 max-w-sm text-muted-foreground text-sm">
            Create a project to organize the services you want Codedock to run.
          </p>
          <Link to="/projects/new" className="mt-5">
            <Button className="gap-2">
              <Plus className="h-4 w-4" />
              Create project
            </Button>
          </Link>
        </div>
      ) : (
        <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
          {projects.map((project) => (
            <ProjectCard key={project.id} project={project} />
          ))}
        </div>
      )}
    </div>
  );
}
