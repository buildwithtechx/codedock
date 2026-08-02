import { Link } from '@tanstack/react-router';
import { FolderKanban, Loader2, Plus } from 'lucide-react';
import { PageHeader } from '#/components/layout/page-header';
import { Button } from '#/components/ui/button';
import { QueryErrorState } from '#/components/ui/query-error-state';
import { WorkspaceEmptyState } from '#/components/ui/workspace-empty-state';
import { ProjectCard } from '#/features/projects/project-card';
import { useListCanvasSummaries } from '#/hooks/use-canvas';

export function ProjectDirectory() {
  const { data, isLoading, isError, refetch } = useListCanvasSummaries();
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
        <div className="flex min-h-[25rem] items-center justify-center">
          <Loader2 className="h-5 w-5 animate-spin text-muted-foreground" />
        </div>
      ) : isError ? (
        <QueryErrorState
          title="Projects are unavailable"
          description="Codedock could not load projects for the active workspace."
          onRetry={() => void refetch()}
        />
      ) : projects.length === 0 ? (
        <WorkspaceEmptyState
          icon={FolderKanban}
          title="Start with a project"
          description="Projects keep the applications, environments, and releases your workspace runs together."
          action={
            <Link to="/projects/new">
              <Button className="gap-2">
                <Plus className="h-4 w-4" />
                Create project
              </Button>
            </Link>
          }
        />
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
