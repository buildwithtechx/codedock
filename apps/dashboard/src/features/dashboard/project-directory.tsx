import { Link } from '@tanstack/react-router';
import { Plus, Search, Server } from 'lucide-react';
import { useEffect, useMemo, useState } from 'react';
import { HelpMenu } from '#/components/layout/help-menu';
import { PageHeader } from '#/components/layout/page-header';
import { Button } from '#/components/ui/button';
import { Input } from '#/components/ui/input';
import { QueryErrorState } from '#/components/ui/query-error-state';
import { ProjectCard } from '#/features/projects/project-card';
import {
  buildHostingFilterOptions,
  type ProjectHostingFilter,
  ProjectHostingFilters,
  projectMatchesHostingFilter,
} from '#/features/projects/project-hosting-filters';
import { ProjectListRow } from '#/features/projects/project-list-row';
import { type ProjectView, ProjectViewToggle } from '#/features/projects/project-view-toggle';
import { useListCanvasSummaries } from '#/hooks/use-canvas';
import { EmptyIllustration } from './empty-illustration';
import { ProjectEmptyState } from './project-empty-state';

const viewStorageKey = 'codedock-projects-view';

export function ProjectDirectory() {
  const { data, isLoading, isError, refetch } = useListCanvasSummaries();
  const projects = data?.data || [];
  const [search, setSearch] = useState('');
  const [filter, setFilter] = useState<ProjectHostingFilter>({ kind: 'all' });
  const [view, setView] = useState<ProjectView>('grid');

  useEffect(() => {
    try {
      const saved = localStorage.getItem(viewStorageKey);
      if (saved === 'grid' || saved === 'list') setView(saved);
    } catch {
      /* Preferences are optional. */
    }
  }, []);

  const changeView = (next: ProjectView) => {
    setView(next);
    try {
      localStorage.setItem(viewStorageKey, next);
    } catch {
      /* Preferences are optional. */
    }
  };

  const filters = useMemo(() => buildHostingFilterOptions(projects), [projects]);
  const showFilters = filters.length > 1;
  const showServerCta = !projects.some((project) => project.deployTarget === 'server');
  const showSidebar = showFilters || showServerCta;
  const filtered = projects.filter(
    (project) =>
      projectMatchesHostingFilter(project, filter) &&
      [project.name, project.description].some((value) =>
        value?.toLowerCase().includes(search.toLowerCase())
      )
  );

  return (
    <div className="space-y-6 pb-20">
      <PageHeader
        title="Projects"
        description={
          isLoading
            ? 'Loading projects...'
            : `${projects.length} project${projects.length === 1 ? '' : 's'}`
        }
        action={
          <div className="flex items-center gap-2">
            <Button asChild className="gap-2">
              <Link to="/projects/new">
                <Plus className="h-4 w-4" />
                New project
              </Link>
            </Button>
            <HelpMenu />
          </div>
        }
      />

      {isLoading ? (
        <div className="divide-y divide-border/50 rounded-2xl bg-card" aria-busy="true">
          {Array.from({ length: 5 }, (_, index) => (
            <div key={index} className="flex animate-pulse items-center gap-4 px-5 py-4">
              <div className="size-10 rounded-xl bg-muted" />
              <div className="flex-1 space-y-2">
                <div className="h-4 w-32 rounded-lg bg-muted" />
                <div className="h-3 w-48 max-w-full rounded-lg bg-muted/60" />
              </div>
              <div className="h-6 w-16 rounded-full bg-muted/60" />
            </div>
          ))}
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
        <div
          className={`grid grid-cols-1 gap-x-6 gap-y-4 ${showSidebar ? 'min-[60rem]:grid-cols-[minmax(0,1fr)_340px]' : ''}`}
        >
          <div className="flex min-w-0 items-center gap-3 min-[60rem]:col-start-1 min-[60rem]:row-start-1">
            <div className="relative min-w-0 flex-1">
              <Search className="pointer-events-none absolute top-1/2 left-3.5 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
              <Input
                type="search"
                placeholder="Search projects"
                aria-label="Search projects"
                value={search}
                onChange={(event) => setSearch(event.target.value)}
                className="h-10 bg-muted/60 ps-10 pe-4"
              />
            </div>
            <ProjectViewToggle value={view} onChange={changeView} />
          </div>

          <div className="min-w-0 min-[60rem]:col-start-1 min-[60rem]:row-start-2">
            {filtered.length > 0 ? (
              view === 'grid' ? (
                <div className="grid grid-cols-[repeat(auto-fill,minmax(min(100%,20rem),1fr))] gap-3">
                  {filtered.map((project) => (
                    <ProjectCard key={project.id} project={project} />
                  ))}
                </div>
              ) : (
                <div className="divide-y divide-border/50 rounded-2xl bg-card">
                  {filtered.map((project) => (
                    <ProjectListRow key={project.id} project={project} />
                  ))}
                </div>
              )
            ) : (
              <div className="flex min-h-80 flex-col items-center justify-center px-6 py-12 text-center">
                <EmptyIllustration className="relative mx-auto mb-6 h-40 w-56" />
                {search ? (
                  <p className="max-w-sm text-muted-foreground text-sm">
                    No projects found for &ldquo;{search}&rdquo;.
                  </p>
                ) : (
                  <>
                    <h2 className="mb-2 font-medium text-foreground text-xl">
                      No projects on this target
                    </h2>
                    <p className="max-w-sm text-muted-foreground text-sm leading-relaxed">
                      Nothing is deployed here yet. Import a repository or move a project to this
                      target to get started.
                    </p>
                  </>
                )}
              </div>
            )}
          </div>

          {showSidebar && (
            <aside className="space-y-4 min-[60rem]:sticky min-[60rem]:top-6 min-[60rem]:col-start-2 min-[60rem]:row-span-2 min-[60rem]:row-start-1 min-[60rem]:self-start">
              {showFilters && (
                <ProjectHostingFilters options={filters} active={filter} onChange={setFilter} />
              )}
              {showServerCta && (
                <div className="rounded-2xl bg-card p-5">
                  <Server className="mb-3 h-5 w-5 text-muted-foreground" />
                  <h2 className="font-medium text-foreground text-sm">Deploy to a server</h2>
                  <p className="mt-1 text-muted-foreground text-sm leading-relaxed">
                    Connect a server over SSH to run production workloads outside the control plane.
                  </p>
                  <Button asChild variant="secondary" size="sm" className="mt-3">
                    <Link to="/servers/new">
                      <Plus className="h-4 w-4" />
                      Add server
                    </Link>
                  </Button>
                </div>
              )}
            </aside>
          )}
        </div>
      )}
    </div>
  );
}
