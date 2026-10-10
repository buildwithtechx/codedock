import { useNavigate } from '@tanstack/react-router';
import { LayoutGrid, Settings } from 'lucide-react';
import { Button } from '#/components/ui/button';
import type { GithubApp } from '#/features/settings';
import { useGetGitApps } from '#/features/settings';
import { GithubIcon } from '#/features/sources/github-app-dialogs';

function GitHubConnectSvg() {
  return (
    <svg className="h-44 w-64" viewBox="0 0 240 156" fill="none" aria-hidden="true">
      <circle
        cx="120"
        cy="74"
        r="56"
        className="stroke-foreground/20 dark:stroke-foreground/25"
        strokeWidth="1.5"
        strokeDasharray="3 7"
      />
      <path
        d="M120 74 L184 40"
        className="stroke-foreground/20 dark:stroke-foreground/25"
        strokeWidth="1.5"
        strokeDasharray="3 4"
      />
      <path
        d="M120 74 L198 94"
        className="stroke-foreground/20 dark:stroke-foreground/25"
        strokeWidth="1.5"
        strokeDasharray="3 4"
      />
      <path
        d="M120 74 L44 100"
        className="stroke-foreground/20 dark:stroke-foreground/25"
        strokeWidth="1.5"
        strokeDasharray="3 4"
      />
      <g>
        <rect
          x="170"
          y="30"
          width="28"
          height="20"
          rx="5"
          className="fill-card stroke-border/80 dark:fill-card dark:stroke-border"
          strokeWidth="1"
        />
        <rect
          x="175"
          y="36"
          width="10"
          height="2.5"
          rx="1.25"
          className="fill-foreground/40 dark:fill-foreground/60"
        />
        <rect
          x="175"
          y="41"
          width="17"
          height="2.5"
          rx="1.25"
          className="fill-foreground/20 dark:fill-foreground/35"
        />
      </g>
      <g>
        <rect
          x="184"
          y="84"
          width="28"
          height="20"
          rx="5"
          className="fill-card stroke-border/80 dark:fill-card dark:stroke-border"
          strokeWidth="1"
        />
        <rect
          x="189"
          y="90"
          width="10"
          height="2.5"
          rx="1.25"
          className="fill-foreground/40 dark:fill-foreground/60"
        />
        <rect
          x="189"
          y="95"
          width="17"
          height="2.5"
          rx="1.25"
          className="fill-foreground/20 dark:fill-foreground/35"
        />
      </g>
      <g>
        <rect
          x="28"
          y="90"
          width="28"
          height="20"
          rx="5"
          className="fill-card stroke-border/80 dark:fill-card dark:stroke-border"
          strokeWidth="1"
        />
        <rect
          x="33"
          y="96"
          width="10"
          height="2.5"
          rx="1.25"
          className="fill-foreground/40 dark:fill-foreground/60"
        />
        <rect
          x="33"
          y="101"
          width="17"
          height="2.5"
          rx="1.25"
          className="fill-foreground/20 dark:fill-foreground/35"
        />
      </g>
      <circle
        cx="120"
        cy="74"
        r="30"
        className="fill-card stroke-border dark:fill-card dark:stroke-border"
        strokeWidth="1.5"
      />
      <g transform="translate(120,74) scale(3) translate(-85,-108)">
        <path
          d="M85 102a6 6 0 0 0-1.9 11.7c.3.05.4-.13.4-.3v-1.05c-1.63.35-1.97-.79-1.97-.79a1.55 1.55 0 0 0-.65-.86c-.53-.36.04-.35.04-.35a1.23 1.23 0 0 1 .9.6 1.25 1.25 0 0 0 1.71.49 1.25 1.25 0 0 1 .37-.78c-1.3-.15-2.67-.65-2.67-2.9a2.27 2.27 0 0 1 .6-1.57 2.1 2.1 0 0 1 .06-1.55s.49-.16 1.6.6a5.5 5.5 0 0 1 2.92 0c1.11-.76 1.6-.6 1.6-.6a2.1 2.1 0 0 1 .06 1.55 2.27 2.27 0 0 1 .6 1.57c0 2.26-1.37 2.75-2.68 2.9a1.4 1.4 0 0 1 .4 1.08v1.6c0 .17.1.35.4.3A6 6 0 0 0 85 102z"
          className="fill-foreground"
        />
      </g>
      <circle cx="26" cy="58" r="3" className="fill-foreground/30 dark:fill-foreground/40" />
      <circle cx="208" cy="124" r="4" className="fill-foreground/25 dark:fill-foreground/35" />
      <circle cx="58" cy="34" r="3" className="fill-foreground/30 dark:fill-foreground/40" />
    </svg>
  );
}

function startAutomaticGithubAppInstall() {
  const origin = window.location.origin;
  const isLocalhost =
    window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1';
  const manifest = {
    name: `codedock-${Math.random().toString(36).substring(7)}`,
    url: origin,
    redirect_url: `${origin}/settings?tab=git`,
    public: false,
    default_permissions: {
      contents: 'read',
      metadata: 'read',
      pull_requests: 'read',
      emails: 'read',
    },
    ...(!isLocalhost
      ? {
          hook_attributes: {
            url: `${origin}/api/webhooks/github/services/generic`,
          },
          default_events: ['push', 'pull_request'],
        }
      : {}),
  };

  const form = document.createElement('form');
  form.method = 'POST';
  form.action = 'https://github.com/settings/apps/new';
  form.target = '_blank';
  form.rel = 'noopener noreferrer';
  const input = document.createElement('input');
  input.type = 'hidden';
  input.name = 'manifest';
  input.value = JSON.stringify(manifest);
  form.appendChild(input);
  document.body.appendChild(form);
  form.submit();
  document.body.removeChild(form);
}

export function ConnectPrompt({ onBrowseApps }: { onBrowseApps: () => void }) {
  const navigate = useNavigate();
  const { data: gitAppsData, isLoading } = useGetGitApps();
  const apps = (gitAppsData?.data as GithubApp[]) || [];
  const primaryApp = apps[0];

  const handleConnect = () => {
    if (primaryApp?.name) {
      window.location.href = `https://github.com/apps/${primaryApp.name}/installations/new`;
      return;
    }
    startAutomaticGithubAppInstall();
  };

  return (
    <div className="rounded-2xl border border-border/50 bg-card">
      <div className="px-6 pt-10 pb-10 text-center">
        <div className="relative mx-auto flex h-44 w-64 items-center justify-center">
          <GitHubConnectSvg />
        </div>
        <h3 className="mt-2 font-medium text-foreground/85 text-lg">Connect GitHub to deploy</h3>
        <p className="mx-auto mt-1.5 mb-7 max-w-md text-muted-foreground text-sm leading-relaxed">
          Connect your GitHub account or install the GitHub App to browse and deploy repositories.
        </p>
        <div className="flex flex-wrap items-center justify-center gap-3">
          <Button type="button" onClick={handleConnect} disabled={isLoading} className="h-11 px-6">
            <GithubIcon className="size-4" />
            {primaryApp?.name ? 'Install GitHub App' : 'Connect GitHub'}
          </Button>
          <Button type="button" variant="secondary" onClick={onBrowseApps} className="h-11 px-6">
            <LayoutGrid className="size-4" aria-hidden="true" />
            Deploy an app instead
          </Button>
        </div>
        <div className="mt-7">
          <button
            type="button"
            onClick={() => void navigate({ to: '/settings', search: { tab: 'git' } as never })}
            className="inline-flex items-center gap-1.5 font-medium text-muted-foreground text-sm transition-colors hover:text-foreground"
          >
            <Settings className="size-3.5" />
            Manage in settings
          </button>
        </div>
      </div>
    </div>
  );
}
