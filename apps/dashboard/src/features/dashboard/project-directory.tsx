import { Link } from '@tanstack/react-router';
import { Loader2, Plus } from 'lucide-react';
import { PageHeader } from '#/components/layout/page-header';
import { Button } from '#/components/ui/button';
import { QueryErrorState } from '#/components/ui/query-error-state';
import { ProjectCard } from '#/features/projects/project-card';
import { useListCanvasSummaries } from '#/hooks/use-canvas';
import { ProjectEmptyState } from './project-empty-state';

export function ProjectDirectory() {
  const { data, isLoading, isError, refetch } = useListCanvasSummaries();
  const projects = data?.data || [];

  return (
    <div className="space-y-6">
      <PageHeader
        title="Projects"
        description={
          isLoading
            ? 'Loading projects...'
            : `${projects.length} project${projects.length === 1 ? '' : 's'}`
        }
        action={
          <Button asChild className="gap-2">
            <Link to="/projects/new">
              <Plus className="h-4 w-4" />
              New project
            </Link>
          </Button>
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
        <ProjectEmptyState />
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
