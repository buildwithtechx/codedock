import { Link } from '@tanstack/react-router';
import { Code2, Search } from 'lucide-react';
import { useEffect, useMemo, useState } from 'react';
import { Button } from '#/components/ui/button';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import { QueryErrorState } from '#/components/ui/query-error-state';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '#/components/ui/select';
import { Skeleton } from '#/components/ui/skeleton';
import { useListAllProjects } from '#/features/projects';
import { CreateGitAppModal } from '#/features/sources/create-git-app-modal';
import { useListExampleApps } from '#/hooks/use-templates';
import type { ExampleApp } from '#/interfaces/templates';

const EXAMPLES_REPO_URL = 'https://github.com/buildwithtechx/codedock-examples.git';
const EXAMPLES_BRANCH = 'main';
const GENERIC_LOGO = '/app-logos/_generic.svg';

function ExampleLogo({ example }: { example: ExampleApp }) {
  const [stage, setStage] = useState(0);
  if (!example.logo || stage > 1) {
    return (
      <span
        aria-hidden="true"
        className="flex size-10 shrink-0 items-center justify-center rounded-xl bg-primary/10 text-primary"
      >
        <Code2 className="size-5" />
      </span>
    );
  }
  return (
    <span className="flex size-10 shrink-0 items-center justify-center overflow-hidden rounded-xl bg-muted/60">
      <img
        src={stage === 0 ? example.logo : GENERIC_LOGO}
        alt=""
        aria-hidden="true"
        className="size-3/5 object-contain"
        onError={() => setStage((current) => current + 1)}
      />
    </span>
  );
}

export function TemplatesGallery() {
  const { data: examplesResponse, isLoading, isError, refetch } = useListExampleApps();
  const projectsQuery = useListAllProjects();
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
    const needle = query.trim().toLowerCase();
    if (!needle) return examples;
    return examples.filter((example) =>
      `${example.name} ${example.description}`.toLowerCase().includes(needle)
    );
  }, [examples, query]);

  if (projectsQuery.isLoading) {
    return <Skeleton className="h-10 w-full max-w-md rounded-xl" />;
  }

  if (projects.length === 0 && !projectsQuery.isLoading) {
    return (
      <div className="rounded-2xl border border-border/50 bg-card px-5 py-12 text-center text-muted-foreground text-sm">
        <p>You need a project before deploying a template.</p>
        <Button asChild size="sm" variant="secondary" className="mt-4">
          <Link to="/projects/new">Create a project</Link>
        </Button>
      </div>
    );
  }

  return (
    <div>
      <div className="space-y-4">
        <div className="w-full max-w-md space-y-2">
          <Label htmlFor="templates-project">Project</Label>
          <Select value={projectId} onValueChange={setProjectId}>
            <SelectTrigger id="templates-project" className="w-full">
              <SelectValue placeholder="Select a project" />
            </SelectTrigger>
            <SelectContent>
              {projects.map((project) => (
                <SelectItem key={project.id} value={project.id}>
                  {project.name}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
          {projectsQuery.isError && (
            <p role="alert" className="text-destructive text-xs">
              Projects could not be loaded.
            </p>
          )}
        </div>
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
            <div
              key={example.id}
              className="flex w-full items-start gap-3 rounded-2xl border border-border/50 bg-card p-5 transition-all hover:border-primary/40 hover:shadow-md"
            >
              <ExampleLogo example={example} />
              <div className="min-w-0 flex-1">
                <p className="truncate font-medium text-foreground">{example.name}</p>
                <p className="mt-1 line-clamp-2 text-muted-foreground text-sm">
                  {example.description}
                </p>
                <div className="mt-3 flex gap-2">
                  <Button
                    size="sm"
                    className="flex-1"
                    disabled={!projectId}
                    onClick={() => setDeployTarget(example)}
                  >
                    Deploy
                  </Button>
                  <Button
                    size="sm"
                    variant="outline"
                    onClick={() => window.open(example.repo, '_blank')}
                  >
                    GitHub
                  </Button>
                </div>
              </div>
            </div>
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
