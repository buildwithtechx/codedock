import { Link, useNavigate, useSearch } from '@tanstack/react-router';
import { ArrowLeft, ArrowRight, FolderKanban, LibraryBig, Server } from 'lucide-react';
import { useEffect, useState } from 'react';
import { toast } from 'sonner';
import { ContextRail } from '#/components/layout/context-rail';
import { Button } from '#/components/ui/button';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '#/components/ui/select';
import { useListOrganizations } from '#/features/organizations';
import { useCreateProject } from '#/features/projects';
import { useListServers } from '#/hooks/use-servers';
import { useOrganizationStore } from '#/stores/organization-store';
import { DEPLOY_PATHS, type DeployPathId, deployPathById } from './project-creation-paths';

type CreationTarget = DeployPathId | 'empty';

const TEMPLATE_PATHS: Record<string, CreationTarget> = {
  'one-click': 'one-click',
  examples: 'examples',
};

export function ProjectCreationPage() {
  const navigate = useNavigate();
  const search = useSearch({ from: '/_dashboard/projects/new' });
  const activeOrganizationId = useOrganizationStore((state) => state.activeOrganizationId);
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [serverId, setServerId] = useState('local');
  const [organizationId, setOrganizationId] = useState('');
  const [target, setTarget] = useState<CreationTarget>('git');
  const { data: servers = [] } = useListServers();
  const { data: organizations = [] } = useListOrganizations();
  const { mutateAsync: createProject, isPending } = useCreateProject();

  useEffect(() => {
    if (activeOrganizationId) {
      setOrganizationId(activeOrganizationId);
    } else if (organizations.length > 0 && !organizationId) {
      setOrganizationId(organizations[0].id);
    }
  }, [activeOrganizationId, organizations, organizationId]);

  useEffect(() => {
    if (search.template && TEMPLATE_PATHS[search.template]) {
      setTarget(TEMPLATE_PATHS[search.template]);
    }
  }, [search.template]);

  const selectedPath = target === 'empty' ? null : deployPathById(target);

  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault();

    if (!organizationId) {
      toast.error('An active organization is required to create a project');
      return;
    }

    try {
      const response = await createProject({
        payload: {
          name,
          description,
          organizationId,
          ...(serverId !== 'local' ? { serverId } : {}),
        },
      });
      const projectId = response.data.id;
      toast.success('Project created');

      if (target === 'empty') {
        await navigate({ to: '/projects/$projectId', params: { projectId } });
        return;
      }
      if (target === 'compose') {
        await navigate({ to: '/projects/$projectId/compose', params: { projectId } });
        return;
      }
      if (target === 'one-click' || target === 'examples') {
        await navigate({
          to: '/projects/$projectId/new',
          params: { projectId },
          search: { tab: target },
        });
        return;
      }
      await navigate({
        to: '/projects/$projectId/new',
        params: { projectId },
        search: { tab: 'resources', resource: target },
      });
    } catch {
      toast.error('Failed to create project');
    }
  };

  return (
    <div className="grid min-w-0 gap-6 lg:grid-cols-[minmax(0,1fr)_21.25rem]">
      <main className="min-w-0">
        <Link
          to="/projects"
          className="inline-flex items-center gap-2 text-muted-foreground text-sm transition-colors hover:text-foreground"
        >
          <ArrowLeft className="h-4 w-4" />
          Projects
        </Link>
        <header className="mt-6">
          <p className="font-medium text-muted-foreground text-sm">New workspace</p>
          <h1 className="mt-1 font-semibold text-2xl tracking-tight">Create a project</h1>
          <p className="mt-1 max-w-2xl text-muted-foreground text-sm">
            Name the workspace first, then pick how its first service gets deployed. You can mix
            paths later inside the project.
          </p>
        </header>

        <form onSubmit={handleSubmit} className="mt-8 max-w-3xl space-y-6">
          <section className="rounded-2xl border border-border bg-card p-6 shadow-sm">
            <div className="flex items-center gap-3">
              <span className="font-semibold text-primary text-sm">01</span>
              <div className="flex h-9 w-9 items-center justify-center rounded-xl bg-primary/12 text-primary">
                <FolderKanban className="h-4 w-4" />
              </div>
              <div>
                <h2 className="font-semibold text-sm">Project details</h2>
                <p className="text-muted-foreground text-xs">
                  Name the workspace your services belong to.
                </p>
              </div>
            </div>

            <div className="mt-6 grid gap-5 sm:grid-cols-2">
              <div className="space-y-2 sm:col-span-2">
                <Label htmlFor="project-name">Project name</Label>
                <Input
                  id="project-name"
                  value={name}
                  onChange={(event) => setName(event.target.value)}
                  placeholder="Acme platform"
                  required
                />
              </div>
              <div className="space-y-2 sm:col-span-2">
                <Label htmlFor="project-description">Description</Label>
                <Input
                  id="project-description"
                  value={description}
                  onChange={(event) => setDescription(event.target.value)}
                  placeholder="Services, environments, and releases for the platform"
                />
              </div>
              <div className="space-y-2">
                <Label>Organization</Label>
                <Select value={organizationId} onValueChange={setOrganizationId}>
                  <SelectTrigger>
                    <SelectValue placeholder="Select an organization" />
                  </SelectTrigger>
                  <SelectContent>
                    {organizations.map((organization) => (
                      <SelectItem key={organization.id} value={organization.id}>
                        {organization.name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
              <div className="space-y-2">
                <Label>Deployment target</Label>
                <Select value={serverId} onValueChange={setServerId}>
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="local">Local control plane</SelectItem>
                    {servers
                      .filter((server) => !server.isControlPlane && server.id !== 'local')
                      .map((server) => (
                        <SelectItem key={server.id} value={server.id}>
                          {server.name} ({server.ipAddress})
                        </SelectItem>
                      ))}
                  </SelectContent>
                </Select>
              </div>
            </div>
          </section>

          <section className="rounded-2xl border border-border bg-card p-6 shadow-sm">
            <div className="flex items-center gap-3">
              <span className="font-semibold text-primary text-sm">02</span>
              <div>
                <h2 className="font-semibold text-sm">Choose a deploy path</h2>
                <p className="text-muted-foreground text-xs">
                  Each path opens the right setup flow right after the project is created.
                </p>
              </div>
            </div>

            <div className="mt-6 grid gap-3 sm:grid-cols-2">
              {DEPLOY_PATHS.map((path) => {
                const active = target === path.id;
                return (
                  <button
                    key={path.id}
                    type="button"
                    aria-pressed={active}
                    onClick={() => setTarget(path.id)}
                    className={`flex flex-col gap-2 rounded-xl border p-4 text-left transition-colors ${
                      active
                        ? 'border-primary/60 bg-primary/5 shadow-sm'
                        : 'border-border/70 bg-background hover:border-primary/35'
                    }`}
                  >
                    <span className="flex items-center gap-2">
                      <path.icon
                        className={`h-4 w-4 ${active ? 'text-primary' : 'text-muted-foreground'}`}
                      />
                      <span className="font-semibold text-sm">{path.title}</span>
                    </span>
                    <span className="text-muted-foreground text-xs leading-5">
                      {path.description}
                    </span>
                    <span className="text-[11px] text-muted-foreground/80">{path.hint}</span>
                  </button>
                );
              })}
              <button
                type="button"
                aria-pressed={target === 'empty'}
                onClick={() => setTarget('empty')}
                className={`flex flex-col gap-2 rounded-xl border border-dashed p-4 text-left transition-colors sm:col-span-2 ${
                  target === 'empty'
                    ? 'border-primary/60 bg-primary/5'
                    : 'border-border/70 bg-background hover:border-primary/35'
                }`}
              >
                <span className="font-semibold text-sm">Decide later</span>
                <span className="text-muted-foreground text-xs leading-5">
                  Create an empty project and add services from the project overview whenever ready.
                </span>
              </button>
            </div>

            <Link
              to="/library"
              className="mt-4 inline-flex items-center gap-2 text-primary text-sm hover:underline"
            >
              <LibraryBig className="h-4 w-4" />
              Prefer to start from code? Browse the library first.
            </Link>
          </section>

          <div className="flex items-center justify-between gap-3">
            <Button asChild type="button" variant="ghost">
              <Link to="/projects">Cancel</Link>
            </Button>
            <Button type="submit" disabled={isPending || !name.trim() || !organizationId}>
              {isPending
                ? 'Creating project...'
                : selectedPath
                  ? `Create and add ${selectedPath.title.toLowerCase()}`
                  : 'Create project'}
              {!isPending && <ArrowRight className="h-4 w-4" />}
            </Button>
          </div>
        </form>
      </main>

      <ContextRail>
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
      </ContextRail>
    </div>
  );
}
