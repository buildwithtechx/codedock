import { ArrowUpRight, Check, Copy, GitBranch, KeyRound, RefreshCw } from 'lucide-react';
import { useMemo, useState } from 'react';
import { toast } from 'sonner';
import { v4 as uuidv4 } from 'uuid';
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
import { ServiceIcon } from '#/components/ui/service-icon';
import type { AppService } from '#/features/services';
import { useUpdateApp } from '#/hooks/use-apps';
import { useListGitBranches } from '#/hooks/use-git';
import type { GitRepo } from '#/interfaces/git';
import { ProjectGitRepoDialog } from './project-git-repo-dialog';

const PROVIDER_HOSTS: Record<string, string> = {
  'github.com': 'github',
  'gitlab.com': 'gitlab',
  'bitbucket.org': 'bitbucket',
};

function parseRepository(repositoryUrl: string): { provider: string; slug: string } {
  try {
    const parsed = new URL(repositoryUrl);
    const provider = PROVIDER_HOSTS[parsed.hostname.toLowerCase()] ?? '';
    const slug = parsed.pathname.replace(/^\//, '').replace(/\.git$/, '');
    return { provider, slug };
  } catch {
    return { provider: '', slug: '' };
  }
}

export function ProjectGitServiceSource({ service }: { service: AppService }) {
  const { mutateAsync: updateApp, isPending } = useUpdateApp();
  const [branchDraft, setBranchDraft] = useState(service.branch);
  const [repoDialogOpen, setRepoDialogOpen] = useState(false);
  const [showToken, setShowToken] = useState(false);

  const parsed = useMemo(() => parseRepository(service.repositoryUrl), [service.repositoryUrl]);
  const { data: branchesRes, isLoading: branchesLoading } = useListGitBranches(
    parsed.provider,
    parsed.slug
  );
  const branches = branchesRes?.data ?? [];
  const branchDirty = branchDraft.trim() !== '' && branchDraft !== service.branch;
  const webhookUrl =
    typeof window !== 'undefined'
      ? `${window.location.origin}/api/webhooks/git/services/${service.id}?token=${service.deployToken || '<token>'}`
      : '';

  const saveSource = async (payload: Partial<AppService>) => {
    try {
      await updateApp({ appId: service.id, payload: { ...service, ...payload } });
      toast.success(`Updated ${service.name}`);
    } catch (error) {
      toast.error(error instanceof Error ? error.message : 'Failed to update service');
    }
  };

  const handleSelectRepository = async (repository: GitRepo) => {
    setRepoDialogOpen(false);
    setBranchDraft(repository.defaultBranch);
    await saveSource({ repositoryUrl: repository.cloneUrl, branch: repository.defaultBranch });
  };

  const handleCopyWebhook = async () => {
    try {
      await navigator.clipboard.writeText(webhookUrl);
      toast.success('Webhook URL copied');
    } catch {
      toast.error('Could not copy the webhook URL');
    }
  };

  const handleRegenerateToken = async () => {
    await saveSource({ deployToken: uuidv4() });
  };

  const handleRevokeToken = async () => {
    await saveSource({ deployToken: '' });
  };

  return (
    <div className="rounded-xl border border-border/60 bg-muted/20 px-4 py-3.5">
      <div className="flex flex-wrap items-start justify-between gap-4">
        <div className="min-w-0 grow basis-64">
          <div className="flex items-center gap-2">
            {service.icon ? (
              <ServiceIcon icon={service.icon} className="h-4 w-4" />
            ) : (
              <GitBranch className="h-4 w-4 text-muted-foreground" />
            )}
            <p className="truncate font-semibold text-sm">{service.name}</p>
          </div>
          {service.repositoryUrl ? (
            <a
              href={service.repositoryUrl}
              target="_blank"
              rel="noopener noreferrer"
              className="group mt-1.5 inline-flex max-w-full items-center gap-1.5"
            >
              <span className="truncate text-muted-foreground text-sm group-hover:text-primary group-hover:underline">
                {parsed.slug || service.repositoryUrl}
              </span>
              <ArrowUpRight className="h-3.5 w-3.5 shrink-0 text-muted-foreground group-hover:text-primary" />
            </a>
          ) : (
            <p className="mt-1.5 text-muted-foreground text-sm">Deploys from a container image.</p>
          )}
          <div className="mt-2 flex items-center gap-1.5 text-muted-foreground text-xs">
            <GitBranch className="h-3.5 w-3.5" />
            <span className="font-mono">{service.branch || 'main'}</span>
            {service.rootDirectory && <span className="truncate">· {service.rootDirectory}</span>}
          </div>
        </div>
        <Button variant="outline" size="sm" onClick={() => setRepoDialogOpen(true)}>
          Change repository
        </Button>
      </div>

      <div className="mt-3 grid gap-3 border-border/50 border-t pt-3 sm:grid-cols-[minmax(0,1fr)_auto] sm:items-end">
        <div className="space-y-2">
          <Label>Deploy branch</Label>
          {parsed.provider && (branchesLoading || branches.length > 0) ? (
            <Select value={branchDraft} onValueChange={setBranchDraft}>
              <SelectTrigger className="font-mono">
                <SelectValue
                  placeholder={branchesLoading ? 'Loading branches' : 'Select a branch'}
                />
              </SelectTrigger>
              <SelectContent>
                {branches.map((branch) => (
                  <SelectItem key={branch.name} value={branch.name} className="font-mono">
                    {branch.name}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          ) : (
            <Input
              value={branchDraft}
              onChange={(event) => setBranchDraft(event.target.value)}
              placeholder="main"
              className="font-mono"
            />
          )}
        </div>
        <Button
          size="sm"
          disabled={!branchDirty || isPending}
          onClick={() => saveSource({ branch: branchDraft })}
        >
          <Check className="h-3.5 w-3.5" />
          {isPending ? 'Saving' : 'Save branch'}
        </Button>
      </div>

      <div className="mt-3 rounded-lg border border-border/50 bg-background/60 p-3">
        <div className="flex flex-wrap items-center justify-between gap-2">
          <p className="font-medium text-sm">Deploy on push</p>
          <span
            className={`rounded-full px-2 py-0.5 font-medium text-[11px] ${
              service.deployToken
                ? 'bg-emerald-500/10 text-emerald-600 dark:text-emerald-400'
                : 'bg-muted text-muted-foreground'
            }`}
          >
            {service.deployToken ? 'Enabled' : 'Disabled'}
          </span>
        </div>
        <p className="mt-1 text-muted-foreground text-xs leading-5">
          Point a repository webhook at this URL and pushes to the deploy branch trigger a new
          deployment.
        </p>
        {service.deployToken ? (
          <div className="mt-2 flex flex-col gap-2">
            <div className="flex items-center gap-2">
              <Input
                readOnly
                type={showToken ? 'text' : 'password'}
                value={webhookUrl}
                className="font-mono text-xs"
                onFocus={(event) => event.target.select()}
              />
              <Button
                variant="outline"
                size="icon-sm"
                onClick={handleCopyWebhook}
                aria-label="Copy webhook URL"
              >
                <Copy className="h-3.5 w-3.5" />
              </Button>
            </div>
            <div className="flex flex-wrap gap-2">
              <Button variant="ghost" size="sm" onClick={() => setShowToken((value) => !value)}>
                {showToken ? 'Hide' : 'Reveal'}
              </Button>
              <Button
                variant="ghost"
                size="sm"
                onClick={handleRegenerateToken}
                disabled={isPending}
              >
                <RefreshCw className="h-3.5 w-3.5" />
                Rotate token
              </Button>
              <Button
                variant="ghost"
                size="sm"
                onClick={handleRevokeToken}
                disabled={isPending}
                className="text-destructive hover:text-destructive"
              >
                Disable
              </Button>
            </div>
          </div>
        ) : (
          <Button
            variant="outline"
            size="sm"
            className="mt-2"
            onClick={handleRegenerateToken}
            disabled={isPending}
          >
            <KeyRound className="h-3.5 w-3.5" />
            Generate webhook URL
          </Button>
        )}
      </div>

      <ProjectGitRepoDialog
        open={repoDialogOpen}
        onOpenChange={setRepoDialogOpen}
        onSelect={(repository) => void handleSelectRepository(repository)}
        title={`Change repository for ${service.name}`}
        description="Pick a repository from a connected git account. The deploy branch resets to its default."
      />
    </div>
  );
}
