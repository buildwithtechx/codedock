import { Link } from '@tanstack/react-router';
import { Check, GitBranch, GitCommitHorizontal, Loader2, Plus, Settings2 } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Badge } from '#/components/ui/badge';
import { Button } from '#/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '#/components/ui/dialog';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import { QueryErrorState } from '#/components/ui/query-error-state';
import { useListByProject as useListAppsByProject } from '#/hooks/use-apps';
import { useListByOrganization as useListDeployments } from '#/hooks/use-deployments';
import { useConnect, useDisconnect, useGetStatus } from '#/hooks/use-git';
import { ProjectGitServiceSource } from './project-git-service-source';

const PROVIDERS = [
  { id: 'github', name: 'GitHub' },
  { id: 'gitlab', name: 'GitLab' },
  { id: 'bitbucket', name: 'Bitbucket' },
  { id: 'gitea', name: 'Gitea' },
];

interface ProjectGitTabProps {
  projectId: string;
  environmentId?: string;
}

export function ProjectGitTab({ projectId, environmentId }: ProjectGitTabProps) {
  const { data: statusRes, isLoading: statusLoading } = useGetStatus();
  const connectMutation = useConnect();
  const disconnectMutation = useDisconnect();
  const [connecting, setConnecting] = useState<string | null>(null);
  const [accountName, setAccountName] = useState('');
  const [accessToken, setAccessToken] = useState('');

  const {
    data: appsRes,
    isLoading: appsLoading,
    isError: appsError,
    refetch: refetchApps,
  } = useListAppsByProject(projectId, environmentId, Boolean(environmentId));
  const { data: deploymentsRes } = useListDeployments({ projectId, limit: 8 });

  const statuses = statusRes?.data ?? [];
  const services = appsRes?.data ?? [];
  const commits = (deploymentsRes?.data?.records ?? []).filter(
    (deployment) => deployment.commitHash
  );

  const handleConnect = (event: React.FormEvent) => {
    event.preventDefault();
    if (!connecting) return;
    connectMutation.mutate(
      { provider: connecting, accessToken, accountName: accountName || 'Personal' },
      {
        onSuccess: () => {
          setConnecting(null);
          setAccessToken('');
          setAccountName('');
          toast.success('Git account connected');
        },
        onError: (error: Error) => {
          toast.error(error.message || 'Failed to connect git account');
        },
      }
    );
  };

  const handleDisconnect = (provider: string) => {
    disconnectMutation.mutate(provider, {
      onSuccess: () => toast.success('Git account disconnected'),
      onError: (error: Error) => toast.error(error.message || 'Failed to disconnect'),
    });
  };

  return (
    <div className="space-y-5">
      <section className="overflow-hidden rounded-2xl border border-border/60 bg-card">
        <div className="flex flex-wrap items-start justify-between gap-3 border-border/50 border-b px-5 py-4">
          <div className="flex items-start gap-3">
            <div className="flex h-9 w-9 shrink-0 items-center justify-center rounded-xl bg-primary/10 text-primary">
              <GitBranch className="h-4 w-4" />
            </div>
            <div>
              <h3 className="font-semibold text-sm">Connected git accounts</h3>
              <p className="mt-0.5 text-muted-foreground text-sm">
                Connect providers once, then attach their repositories to any service.
              </p>
            </div>
          </div>
          <Button variant="outline" size="sm" asChild>
            <Link to="/settings">
              <Settings2 className="h-3.5 w-3.5" />
              GitHub Apps
            </Link>
          </Button>
        </div>
        <div className="grid gap-3 px-5 py-4 sm:grid-cols-2">
          {statusLoading ? (
            <div className="flex items-center gap-2 text-muted-foreground text-sm sm:col-span-2">
              <Loader2 className="h-4 w-4 animate-spin" />
              Loading git accounts
            </div>
          ) : (
            PROVIDERS.map((provider) => {
              const status = statuses.find(
                (entry) => entry.provider === provider.id && entry.connected
              );
              return (
                <div
                  key={provider.id}
                  className="flex items-center justify-between gap-3 rounded-xl border border-border/60 bg-muted/20 px-4 py-3"
                >
                  <div className="flex min-w-0 items-center gap-3">
                    <img
                      src={`/git-providers/${provider.id}-icon.svg`}
                      alt=""
                      className="h-6 w-6 shrink-0 object-contain"
                    />
                    <div className="min-w-0">
                      <p className="font-medium text-sm">{provider.name}</p>
                      <p className="truncate text-muted-foreground text-xs">
                        {status ? 'Connected' : 'Not connected'}
                      </p>
                    </div>
                  </div>
                  {status ? (
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => handleDisconnect(provider.id)}
                      disabled={disconnectMutation.isPending}
                      className="text-destructive hover:text-destructive"
                    >
                      Disconnect
                    </Button>
                  ) : (
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => {
                        setConnecting(provider.id);
                        setAccessToken('');
                        setAccountName('');
                      }}
                    >
                      Connect
                    </Button>
                  )}
                </div>
              );
            })
          )}
        </div>
      </section>

      <section className="overflow-hidden rounded-2xl border border-border/60 bg-card">
        <div className="flex items-start gap-3 border-border/50 border-b px-5 py-4">
          <div className="flex h-9 w-9 shrink-0 items-center justify-center rounded-xl bg-primary/10 text-primary">
            <GitCommitHorizontal className="h-4 w-4" />
          </div>
          <div>
            <h3 className="font-semibold text-sm">Service sources</h3>
            <p className="mt-0.5 text-muted-foreground text-sm">
              Each service tracks a repository and branch. Push webhooks redeploy on every push.
            </p>
          </div>
        </div>
        <div className="space-y-3 px-5 py-4">
          {appsLoading ? (
            <div className="flex items-center gap-2 text-muted-foreground text-sm">
              <Loader2 className="h-4 w-4 animate-spin" />
              Loading services
            </div>
          ) : appsError ? (
            <QueryErrorState
              title="Services are unavailable"
              description="Codedock could not load this environment."
              onRetry={() => void refetchApps()}
            />
          ) : services.length === 0 ? (
            <div className="flex flex-col items-center rounded-xl border border-border/70 border-dashed bg-muted/20 px-6 py-8 text-center">
              <p className="font-medium text-sm">No services in this environment yet</p>
              <p className="mt-1 max-w-sm text-muted-foreground text-xs leading-5">
                Add a git-backed service to manage its repository, branch, and push webhooks here.
              </p>
              <Button size="sm" className="mt-4" asChild>
                <Link to="/projects/$projectId/new" params={{ projectId }}>
                  <Plus className="h-4 w-4" />
                  Add service
                </Link>
              </Button>
            </div>
          ) : (
            services.map((service) => (
              <ProjectGitServiceSource key={service.id} service={service} />
            ))
          )}
        </div>
      </section>

      <section className="overflow-hidden rounded-2xl border border-border/60 bg-card">
        <div className="border-border/50 border-b px-5 py-4">
          <h3 className="font-semibold text-sm">Recent commits</h3>
          <p className="mt-0.5 text-muted-foreground text-sm">
            Commits that produced the latest deployments in this project.
          </p>
        </div>
        <div className="px-5 py-4">
          {commits.length === 0 ? (
            <div className="rounded-xl border border-border/60 border-dashed bg-muted/20 px-4 py-5 text-center">
              <p className="font-medium text-sm">No commits deployed yet</p>
              <p className="mt-1 text-muted-foreground text-sm">
                Deploy a git-backed service to see its commit history here.
              </p>
            </div>
          ) : (
            <div className="divide-y divide-border/50 overflow-hidden rounded-xl border border-border/50">
              {commits.map((deployment) => (
                <div key={deployment.id} className="flex items-center gap-3 px-4 py-2.5">
                  <div className="min-w-0 flex-1">
                    <div className="flex flex-wrap items-center gap-x-1.5 gap-y-1 text-xs">
                      <span className="max-w-full truncate font-medium text-muted-foreground">
                        {deployment.serviceName}
                      </span>
                      <span className="text-muted-foreground/40">·</span>
                      <span className="text-muted-foreground">
                        {new Date(deployment.createdAt).toLocaleString(undefined, {
                          month: 'short',
                          day: 'numeric',
                          hour: '2-digit',
                          minute: '2-digit',
                        })}
                      </span>
                      {deployment.commitHash && (
                        <code className="rounded-full bg-muted/60 px-1.5 py-px font-medium text-muted-foreground text-xs">
                          {deployment.commitHash.slice(0, 7)}
                        </code>
                      )}
                      {deployment.branch && (
                        <Badge variant="secondary" className="px-1.5 py-0 font-mono text-[10px]">
                          {deployment.branch}
                        </Badge>
                      )}
                    </div>
                    <p className="mt-0.5 truncate text-foreground text-sm">
                      {deployment.commitMessage || deployment.trigger || 'Deployment'}
                    </p>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </section>

      <Dialog open={connecting !== null} onOpenChange={(open) => !open && setConnecting(null)}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle>Connect {connecting}</DialogTitle>
            <DialogDescription>
              Enter a personal access token for this git account.
            </DialogDescription>
          </DialogHeader>
          <form onSubmit={handleConnect} className="space-y-4">
            <div className="space-y-2">
              <Label>Account name</Label>
              <Input
                value={accountName}
                onChange={(event) => setAccountName(event.target.value)}
                placeholder="Personal"
              />
            </div>
            <div className="space-y-2">
              <Label>Personal access token</Label>
              <Input
                type="password"
                value={accessToken}
                onChange={(event) => setAccessToken(event.target.value)}
                placeholder="Token with repository read access"
                required
                className="font-mono"
              />
            </div>
            <div className="flex justify-end gap-2">
              <Button type="button" variant="ghost" onClick={() => setConnecting(null)}>
                Cancel
              </Button>
              <Button type="submit" disabled={connectMutation.isPending}>
                <Check className="h-3.5 w-3.5" />
                {connectMutation.isPending ? 'Connecting' : 'Connect'}
              </Button>
            </div>
          </form>
        </DialogContent>
      </Dialog>
    </div>
  );
}
