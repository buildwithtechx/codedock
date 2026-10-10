import { Cloud, GitBranch, Globe, Loader2, MoreHorizontal, Plus, Rocket } from 'lucide-react';
import { useState } from 'react';
import { Button } from '#/components/ui/button';
import { Input } from '#/components/ui/input';
import { useGetSettings } from '#/features/settings';

interface DeploySidebarProps {
  repoName: string;
  branch: string;
  onBranchChange: (b: string) => void;
  projectName: string;
  framework: string;
  runtimeMode: 'web' | 'worker' | 'static';
  installCommand: string;
  buildCommand: string;
  outputDirectory: string;
  subdomain: string;
  onSubdomainChange: (s: string) => void;
  isDeploying: boolean;
  onDeploy: () => void;
}

export function DeploySidebar({
  repoName,
  branch,
  onBranchChange,
  projectName,
  framework,
  runtimeMode,
  installCommand,
  buildCommand,
  outputDirectory,
  subdomain,
  onSubdomainChange,
  isDeploying,
  onDeploy,
}: DeploySidebarProps) {
  const { data: settingsRes } = useGetSettings();
  const wildcardDomain = settingsRes?.data?.defaultWildcardDomain?.trim() || 'codedock.run';
  const serverIp = settingsRes?.data?.publicIpv4 || settingsRes?.data?.traefikWildcardIp || '';
  const cleanIp = serverIp ? serverIp.replace(/\./g, '-') : '';
  const sslipSuffix = cleanIp ? `${cleanIp}.sslip.io` : 'sslip.io';

  const [domainMode, setDomainMode] = useState<'free' | 'custom' | 'none'>('free');
  const [freeProvider, setFreeProvider] = useState<'wildcard' | 'sslip'>('wildcard');

  const activeSuffix = freeProvider === 'wildcard' ? wildcardDomain : sslipSuffix;

  return (
    <div className="space-y-4 lg:sticky lg:top-6 lg:self-start">
      <div className="rounded-2xl border border-border/60 bg-card p-4">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            <GitBranch className="size-4 text-muted-foreground" />
            <span className="max-w-[200px] truncate font-semibold text-foreground text-xs">
              {repoName || 'my-project'}
            </span>
          </div>
          <button type="button" className="text-muted-foreground hover:text-foreground">
            <MoreHorizontal className="size-4" />
          </button>
        </div>
        <div className="mt-3">
          <select
            value={branch}
            onChange={(e) => onBranchChange(e.target.value)}
            className="w-full rounded-lg border border-border/60 bg-background px-3 py-1.5 font-mono text-foreground text-xs"
          >
            <option value="main">main</option>
            <option value="master">master</option>
            <option value="staging">staging</option>
          </select>
        </div>
      </div>

      <div className="rounded-2xl border border-border/60 bg-card p-4">
        <div className="flex rounded-lg bg-muted/40 p-1">
          {(['free', 'custom', 'none'] as const).map((mode) => (
            <button
              key={mode}
              type="button"
              onClick={() => setDomainMode(mode)}
              className={`flex-1 rounded-md py-1 font-medium text-xs capitalize transition-colors ${
                domainMode === mode
                  ? 'bg-card text-foreground shadow-sm'
                  : 'text-muted-foreground hover:text-foreground'
              }`}
            >
              {mode}
            </button>
          ))}
        </div>

        {domainMode !== 'none' && (
          <div className="mt-4 space-y-3">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-2">
                <Globe className="size-4 text-muted-foreground" />
                <div>
                  <p className="font-medium text-foreground text-xs">Domain</p>
                  <p className="text-[10px] text-muted-foreground">
                    Where your site will be accessible
                  </p>
                </div>
              </div>
              <button
                type="button"
                className="grid size-6 place-items-center rounded-lg border border-border/60 text-muted-foreground hover:bg-muted/50"
              >
                <Plus className="size-3.5" />
              </button>
            </div>

            {domainMode === 'free' && (
              <div className="flex items-center gap-1.5 pb-1">
                <button
                  type="button"
                  onClick={() => setFreeProvider('wildcard')}
                  className={`rounded-md border px-2 py-0.5 font-medium text-[10px] transition-colors ${
                    freeProvider === 'wildcard'
                      ? 'border-primary/40 bg-primary/10 text-primary'
                      : 'border-transparent text-muted-foreground hover:bg-muted/50'
                  }`}
                >
                  .{wildcardDomain}
                </button>
                <button
                  type="button"
                  onClick={() => setFreeProvider('sslip')}
                  className={`rounded-md border px-2 py-0.5 font-medium text-[10px] transition-colors ${
                    freeProvider === 'sslip'
                      ? 'border-primary/40 bg-primary/10 text-primary'
                      : 'border-transparent text-muted-foreground hover:bg-muted/50'
                  }`}
                >
                  .{sslipSuffix}
                </button>
              </div>
            )}

            <div className="flex items-center rounded-lg border border-border/60 bg-background px-3 py-1.5">
              <Input
                value={subdomain}
                onChange={(e) => onSubdomainChange(e.target.value)}
                placeholder="my-app"
                className="h-6 border-0 p-0 text-xs shadow-none focus-visible:ring-0"
              />
              <span className="font-mono text-muted-foreground text-xs">
                {domainMode === 'free' ? `.${activeSuffix}` : ''}
              </span>
            </div>

            {runtimeMode === 'static' && (
              <div className="flex items-center justify-between text-muted-foreground text-xs">
                <span>Static path</span>
                <code className="font-mono text-foreground">/</code>
              </div>
            )}
          </div>
        )}
      </div>

      <Button
        className="w-full gap-2 py-6 font-semibold text-sm shadow-md"
        size="lg"
        onClick={onDeploy}
        disabled={isDeploying || !projectName.trim()}
      >
        {isDeploying ? <Loader2 className="size-4 animate-spin" /> : <Rocket className="size-4" />}
        {isDeploying ? 'Deploying...' : 'Deploy'}
      </Button>

      <div className="rounded-2xl border border-border/60 bg-card p-4">
        <h4 className="font-semibold text-[11px] text-muted-foreground uppercase tracking-wider">
          Deploy Summary
        </h4>
        <div className="mt-3 space-y-2.5 text-xs">
          <div className="flex items-start justify-between gap-2">
            <span className="flex items-center gap-1.5 text-muted-foreground">
              <Globe className="size-3.5" />
              Domain
            </span>
            <span className="truncate font-mono text-foreground">
              {domainMode === 'none'
                ? 'None'
                : subdomain
                  ? domainMode === 'free'
                    ? `${subdomain}.${activeSuffix}`
                    : subdomain
                  : 'Auto-generated'}
            </span>
          </div>

          <div className="flex items-start justify-between gap-2">
            <span className="flex items-center gap-1.5 text-muted-foreground">
              <Cloud className="size-3.5" />
              Build Location
            </span>
            <span className="text-foreground">Codedock Cloud</span>
          </div>

          <div className="flex items-start justify-between gap-2">
            <span className="text-muted-foreground">Framework</span>
            <span className="font-medium text-foreground">{framework}</span>
          </div>

          <div className="flex items-start justify-between gap-2">
            <span className="text-muted-foreground">Install</span>
            <span className="max-w-[140px] truncate font-mono text-[11px] text-foreground">
              {installCommand || 'Default'}
            </span>
          </div>

          <div className="flex items-start justify-between gap-2">
            <span className="text-muted-foreground">Build</span>
            <span className="max-w-[140px] truncate font-mono text-[11px] text-foreground">
              {buildCommand || 'None'}
            </span>
          </div>

          <div className="flex items-start justify-between gap-2">
            <span className="text-muted-foreground">Output</span>
            <span className="font-mono text-[11px] text-foreground">
              {outputDirectory || 'dist'}
            </span>
          </div>
        </div>
      </div>
    </div>
  );
}
