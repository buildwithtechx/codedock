import { Link, useNavigate } from '@tanstack/react-router';
import { ArrowLeft, FolderKanban, Server } from 'lucide-react';
import { useEffect, useState } from 'react';
import { toast } from 'sonner';
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

export function ProjectCreationPage() {
  const navigate = useNavigate();
  const activeOrganizationId = useOrganizationStore((state) => state.activeOrganizationId);
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [serverId, setServerId] = useState('local');
  const [organizationId, setOrganizationId] = useState('');
  const { data: servers = [] } = useListServers();
  const { data: organizations = [] } = useListOrganizations();
  const { mutateAsync: createProject, isPending } = useCreateProject();

  useEffect(() => {
    if (activeOrganizationId) {
      setOrganizationId(activeOrganizationId);
    }
  }, [activeOrganizationId]);

  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault();

    try {
      const response = await createProject({
        payload: {
          name,
          description,
          ...(serverId !== 'local' ? { serverId } : {}),
          ...(organizationId ? { organizationId } : {}),
        },
      });
      toast.success('Project created');
      await navigate({ to: '/projects/$projectId', params: { projectId: response.data.id } });
    } catch {
      toast.error('Failed to create project');
    }
  };

  return (
    <div className="grid min-w-0 gap-6 xl:grid-cols-[minmax(0,1fr)_21.25rem]">
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
            Start with a project, then choose its services, sources, and deployment setup.
          </p>
        </header>

        <form onSubmit={handleSubmit} className="mt-8 max-w-3xl space-y-6">
          <section className="rounded-2xl border border-border bg-card p-6 shadow-sm">
            <div className="flex items-center gap-3">
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
                      .filter((server) => !server.isControlPlane)
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

          <div className="flex items-center justify-between gap-3">
            <Link to="/projects">
              <Button type="button" variant="ghost">
                Cancel
              </Button>
            </Link>
            <Button type="submit" disabled={isPending}>
              {isPending ? 'Creating project...' : 'Create project'}
            </Button>
          </div>
        </form>
      </main>

      <aside className="hidden xl:sticky xl:top-6 xl:block xl:self-start">
        <section className="rounded-2xl border border-border bg-card p-5 shadow-sm">
          <div className="flex items-center gap-2">
            <Server className="h-4 w-4 text-primary" />
            <h2 className="font-semibold text-sm">What happens next</h2>
          </div>
          <ol className="mt-5 space-y-4 text-sm">
            <li className="flex gap-3">
              <span className="font-semibold text-primary">01</span>
              <span className="text-muted-foreground">
                Choose an application, database, image, or compose service.
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
      </aside>
    </div>
  );
}
