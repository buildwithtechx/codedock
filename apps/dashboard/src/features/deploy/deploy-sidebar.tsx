import {
  Container,
  Flame,
  GitBranch,
  Github,
  Globe,
  Loader2,
  MoreHorizontal,
  Plus,
  Rocket,
} from 'lucide-react';
import { useState } from 'react';
import { Button } from '#/components/ui/button';
import { Input } from '#/components/ui/input';
import { useGetPublicSettings, useGetSettings } from '#/features/settings';

interface DeploySidebarProps {
  repoName: string;
  branch: string;
  onBranchChange: (b: string) => void;
  projectName: string;
  framework: string;
  subdomain: string;
  onSubdomainChange: (s: string) => void;
  exposedPort?: number;
  isDeploying: boolean;
  onDeploy: () => void;
}

export function DeploySidebar({
  repoName,
  branch,
  onBranchChange,
  projectName,
  framework,
  subdomain,
  onSubdomainChange,
  exposedPort = 8080,
  isDeploying,
  onDeploy,
}: DeploySidebarProps) {
  const { data: publicRes } = useGetPublicSettings();
  const { data: settingsRes } = useGetSettings();
  const isCloud = Boolean(publicRes?.data?.cloudMode);

  const serverIp = settingsRes?.data?.publicIpv4 || settingsRes?.data?.traefikWildcardIp || '';
  const cleanIp = serverIp ? serverIp.replace(/\./g, '-') : '';
  const sslipSuffix = cleanIp ? `${cleanIp}.sslip.io` : 'sslip.io';

  const configuredWildcard = settingsRes?.data?.defaultWildcardDomain?.trim();
  const defaultFreeSuffix = isCloud ? 'codedock.run' : configuredWildcard || sslipSuffix;

  const [domainMode, setDomainMode] = useState<'free' | 'custom' | 'none'>('free');
  const [customDomain, setCustomDomain] = useState('');

  const activeDomain =
    domainMode === 'none'
      ? null
      : domainMode === 'custom'
        ? customDomain || null
        : subdomain
          ? `${subdomain}.${defaultFreeSuffix}`
          : null;

  const isDocker = framework.toLowerCase() === 'dockerfile';

  return (
    <div className="space-y-4 lg:sticky lg:top-6 lg:self-start">
      <div className="rounded-2xl border border-border/60 bg-card p-4">
        <div className="flex items-center gap-1.5 pb-3">
          <span className="size-2 rounded-full bg-muted-foreground/30" />
          <span className="size-2 rounded-full bg-muted-foreground/30" />
          <span className="size-2 rounded-full bg-muted-foreground/30" />
        </div>

        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            <Github className="size-4 text-muted-foreground" />
            <span className="max-w-[200px] truncate font-medium text-foreground text-xs">
              {repoName || 'my-project'}
            </span>
          </div>
          <button type="button" className="text-muted-foreground hover:text-foreground">
            <MoreHorizontal className="size-4" />
          </button>
        </div>

        <div className="mt-3">
          <div className="flex items-center gap-2 rounded-xl border border-border/60 bg-background/50 px-3 py-1.5 text-xs">
            <GitBranch className="size-3.5 text-muted-foreground" />
            <select
              value={branch}
              onChange={(e) => onBranchChange(e.target.value)}
              className="w-full bg-transparent font-mono text-foreground focus:outline-none"
            >
              <option value="main">main</option>
              <option value="master">master</option>
              <option value="staging">staging</option>
            </select>
          </div>
        </div>
      </div>

      <div className="flex rounded-xl bg-muted/40 p-1">
        {(['free', 'custom', 'none'] as const).map((mode) => (
          <button
            key={mode}
            type="button"
            onClick={() => setDomainMode(mode)}
            className={`flex-1 rounded-lg py-1.5 font-medium text-xs capitalize transition-colors ${
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
        <div className="rounded-2xl border border-border/60 bg-card p-4">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-2.5">
              <div className="flex size-7 items-center justify-center rounded-lg bg-muted text-muted-foreground">
                <Globe className="size-4" />
              </div>
              <div>
                <p className="font-semibold text-foreground text-xs">Domain</p>
                <p className="text-[10px] text-muted-foreground">
                  Where your site will be accessible
                </p>
              </div>
            </div>
            <button
              type="button"
              className="flex size-6 items-center justify-center rounded-full bg-muted text-muted-foreground hover:bg-muted/80"
            >
              <Plus className="size-3.5" />
            </button>
          </div>

          <div className="mt-4 flex items-center rounded-xl border border-border/60 bg-background/50 px-3 py-2">
            {domainMode === 'free' ? (
              <>
                <Input
                  value={subdomain}
                  onChange={(e) => onSubdomainChange(e.target.value)}
                  placeholder="send-matrix"
                  className="h-6 border-0 p-0 font-mono text-xs shadow-none focus-visible:ring-0"
                />
                <span className="font-mono text-muted-foreground text-xs">
                  .{defaultFreeSuffix}
                </span>
              </>
            ) : (
              <Input
                value={customDomain}
                onChange={(e) => setCustomDomain(e.target.value)}
                placeholder="app.example.com"
                className="h-6 border-0 p-0 font-mono text-xs shadow-none focus-visible:ring-0"
              />
            )}
          </div>

          <div className="mt-3 flex items-center justify-between text-xs">
            <span className="flex items-center gap-1 text-muted-foreground text-xs">
              <span className="font-mono text-muted-foreground">#</span>
              Exposed port
            </span>
            <span className="font-medium font-mono text-foreground text-xs">{exposedPort}</span>
          </div>
        </div>
      )}

      <Button
        className="w-full gap-2 rounded-xl bg-foreground py-3 font-semibold text-background text-sm shadow-md hover:bg-foreground/90"
        size="lg"
        onClick={onDeploy}
        disabled={isDeploying || !projectName.trim()}
      >
        {isDeploying ? <Loader2 className="size-4 animate-spin" /> : <Rocket className="size-4" />}
        {isDeploying ? 'Deploying...' : 'Deploy'}
      </Button>

      <div className="rounded-2xl border border-border/60 bg-card p-4">
        <h3 className="font-semibold text-[11px] text-muted-foreground uppercase tracking-wider">
          Deploy Summary
        </h3>
        <div className="mt-3.5 space-y-3 text-xs">
          <div className="flex items-start justify-between gap-2">
            <span className="flex items-center gap-1.5 text-muted-foreground">
              <Globe className="size-3.5" />
              Domain
            </span>
            <span className="truncate font-mono text-foreground">{activeDomain || '—'}</span>
          </div>

          <div className="flex items-start justify-between gap-2">
            <span className="flex items-center gap-1.5 text-muted-foreground">
              {isDocker ? <Container className="size-3.5" /> : <Flame className="size-3.5" />}
              Runtime
            </span>
            <span className="font-medium text-foreground capitalize">{framework}</span>
          </div>

          <div className="flex items-start justify-between gap-2">
            <span className="flex items-center gap-1.5 text-muted-foreground">
              <span className="font-mono">#</span>
              Port
            </span>
            <span className="font-mono text-foreground">{exposedPort}</span>
          </div>
        </div>
      </div>
    </div>
  );
}
