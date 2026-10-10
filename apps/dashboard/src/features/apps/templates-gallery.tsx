import { ArrowRight, Search } from 'lucide-react';
import { useEffect, useMemo, useState } from 'react';
import { toast } from 'sonner';
import { Input } from '#/components/ui/input';
import { QueryErrorState } from '#/components/ui/query-error-state';
import { Skeleton } from '#/components/ui/skeleton';
import { useCreateProject, useListAllProjects } from '#/features/projects';
import { CreateGitAppModal } from '#/features/sources/create-git-app-modal';
import { useListExampleApps } from '#/hooks/use-templates';
import type { ExampleApp } from '#/interfaces/templates';
import { useOrganizationStore } from '#/stores/organization-store';
import { ExampleLogo } from './example-logo';

const EXAMPLES_REPO_URL = 'https://github.com/buildwithtechx/codedock-examples.git';
const EXAMPLES_BRANCH = 'main';

export function TemplatesGallery() {
  const { data: examplesResponse, isLoading, isError, refetch } = useListExampleApps();
  const projectsQuery = useListAllProjects();
  const { mutateAsync: createProject, isPending: isCreatingProject } = useCreateProject();
  const activeOrgId = useOrganizationStore((s) => s.activeOrganizationId);

  const [query, setQuery] = useState('');
  const [projectId, setProjectId] = useState('');
  const [deployTarget, setDeployTarget] = useState<ExampleApp | null>(null);

  const examples = Array.isArray(examplesResponse) ? examplesResponse : [];
  const projects = projectsQuery.data ?? [];

  useEffect(() => {
    if (!projectId && projects.length > 0) {
      setProjectId(projects[0].id);
    }
  }, [projectId, projects]);

  const filtered = useMemo(() => {
    const needle = (query ?? '').trim().toLowerCase();
    if (!needle) return examples;
    return examples.filter((example) =>
      `${example.name} ${example.description}`.toLowerCase().includes(needle)
    );
  }, [examples, query]);

  const handleDeploy = async (example: ExampleApp) => {
    let targetProjectId = projectId;
    if (!targetProjectId) {
      if (projects.length > 0) {
        targetProjectId = projects[0].id;
        setProjectId(targetProjectId);
      } else {
        try {
          const res = await createProject({
            payload: {
              name: example.name,
              organizationId: activeOrgId || '',
            },
          });
          targetProjectId = res.data.id;
          setProjectId(targetProjectId);
          toast.success(`Project "${example.name}" created`);
        } catch {
          toast.error('Failed to create project for template');
          return;
        }
      }
    }
    setDeployTarget(example);
  };

  return (
    <div>
      <div className="relative w-full max-w-md">
        <Search className="pointer-events-none absolute top-1/2 left-3.5 size-4 -translate-y-1/2 text-muted-foreground" />
        <Input
          type="search"
          value={query}
          onChange={(event) => setQuery(event.target.value)}
          placeholder="Search templates…"
          aria-label="Search templates"
          className="pr-4 pl-10"
        />
      </div>

      {isError ? (
        <QueryErrorState
          title="Examples are unavailable"
          description="Codedock could not load the example project catalogue."
          onRetry={() => void refetch()}
        />
      ) : isLoading ? (
        <div className="mt-6 grid grid-cols-1 gap-4 md:grid-cols-2 xl:grid-cols-3">
          {Array.from({ length: 6 }, (_, index) => (
            <Skeleton key={index} className="h-32 rounded-2xl" />
          ))}
        </div>
      ) : filtered.length === 0 ? (
        <div className="mt-6 rounded-2xl border border-border/50 bg-card px-5 py-12 text-center text-muted-foreground text-sm">
          <p>
            {examples.length === 0
              ? 'No example projects are available on this instance yet.'
              : 'No templates match your search.'}
          </p>
        </div>
      ) : (
        <div className="mt-6 grid grid-cols-1 gap-4 md:grid-cols-2 xl:grid-cols-3">
          {filtered.map((example) => (
            <button
              key={example.id}
              type="button"
              disabled={isCreatingProject}
              onClick={() => void handleDeploy(example)}
              className="group flex w-full cursor-pointer items-start gap-3 rounded-2xl border border-border/50 bg-card p-5 text-left transition-all hover:border-primary/40 hover:shadow-md disabled:cursor-not-allowed disabled:opacity-60"
            >
              <ExampleLogo example={example} />
              <div className="min-w-0 flex-1">
                <div className="flex items-center justify-between gap-2">
                  <p className="truncate font-medium text-foreground">{example.name}</p>
                  <ArrowRight className="size-4 shrink-0 text-muted-foreground/40 transition-colors group-hover:text-foreground" />
                </div>
                <p className="mt-1 line-clamp-2 text-muted-foreground text-sm">
                  {example.description}
                </p>
              </div>
            </button>
          ))}
        </div>
      )}

      <CreateGitAppModal
        isOpen={deployTarget !== null}
        onOpenChange={(open) => {
          if (!open) {
            setDeployTarget(null);
          }
        }}
        projectId={projectId}
        initialName={deployTarget?.id ?? ''}
        initialRepositoryUrl={EXAMPLES_REPO_URL}
        initialBranch={EXAMPLES_BRANCH}
        initialRootDirectory={deployTarget?.id ?? ''}
      />
    </div>
  );
}
