import { useNavigate } from '@tanstack/react-router';
import { LayoutGrid, Loader2, Settings } from 'lucide-react';
import { type FormEvent, useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import { Dialog, DialogContent, DialogDescription, DialogTitle } from '#/components/ui/dialog';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import type { GithubApp } from '#/features/settings';
import { useGetGitApps } from '#/features/settings';
import { GithubIcon } from '#/features/sources/github-app-dialogs';
import { useConnectProvider } from './hooks';

function ConnectIllustration() {
  return (
    <svg className="h-44 w-64" viewBox="0 0 256 176" fill="none" aria-hidden="true">
      <circle
        cx="128"
        cy="88"
        r="62"
        stroke="var(--border)"
        strokeWidth="1.5"
        strokeDasharray="4 8"
      />
      <path
        d="M128 88 196 48"
        stroke="var(--muted-foreground)"
        strokeOpacity="0.3"
        strokeWidth="1.5"
        strokeDasharray="4 4"
      />
      <path
        d="M128 88 208 108"
        stroke="var(--muted-foreground)"
        strokeOpacity="0.3"
        strokeWidth="1.5"
        strokeDasharray="4 4"
      />
      <path
        d="M128 88 52 116"
        stroke="var(--muted-foreground)"
        strokeOpacity="0.3"
        strokeWidth="1.5"
        strokeDasharray="4 4"
      />
      <g>
        <rect
          x="182"
          y="38"
          width="30"
          height="22"
          rx="6"
          fill="var(--muted)"
          stroke="var(--border)"
        />
        <rect
          x="187"
          y="44"
          width="11"
          height="3"
          rx="1.5"
          fill="var(--muted-foreground)"
          fillOpacity="0.5"
        />
        <rect
          x="187"
          y="50"
          width="18"
          height="3"
          rx="1.5"
          fill="var(--muted-foreground)"
          fillOpacity="0.25"
        />
      </g>
      <g>
        <rect
          x="194"
          y="98"
          width="30"
          height="22"
          rx="6"
          fill="var(--muted)"
          stroke="var(--border)"
        />
        <rect
          x="199"
          y="104"
          width="11"
          height="3"
          rx="1.5"
          fill="var(--muted-foreground)"
          fillOpacity="0.5"
        />
        <rect
          x="199"
          y="110"
          width="18"
          height="3"
          rx="1.5"
          fill="var(--muted-foreground)"
          fillOpacity="0.25"
        />
      </g>
      <g>
        <rect
          x="36"
          y="106"
          width="30"
          height="22"
          rx="6"
          fill="var(--muted)"
          stroke="var(--border)"
        />
        <rect
          x="41"
          y="112"
          width="11"
          height="3"
          rx="1.5"
          fill="var(--muted-foreground)"
          fillOpacity="0.5"
        />
        <rect
          x="41"
          y="118"
          width="18"
          height="3"
          rx="1.5"
          fill="var(--muted-foreground)"
          fillOpacity="0.25"
        />
      </g>
      <circle cx="128" cy="88" r="32" fill="var(--card)" stroke="var(--border)" strokeWidth="1.5" />
      <g transform="translate(128,88) scale(2.8) translate(-12,-12)">
        <path
          d="M15 22v-4a4.8 4.8 0 0 0-1-3.5c3 0 6-2 6-5.5.08-1.25-.27-2.48-1-3.5.28-1.15.28-2.35 0-3.5 0 0-1 0-3 1.5-2.64-.5-5.36-.5-8 0C6 2 5 2 5 2c-.3 1.15-.3 2.35 0 3.5A5.403 5.403 0 0 0 4 9c0 3.5 3 5.5 6 5.5-.39.49-.68 1.05-.85 1.65-.17.6-.22 1.23-.15 1.85v4"
          stroke="var(--foreground)"
          strokeWidth="1.5"
          strokeLinecap="round"
          strokeLinejoin="round"
        />
        <path
          d="M9 18c-4.51 2-5-2-7-2"
          stroke="var(--foreground)"
          strokeWidth="1.5"
          strokeLinecap="round"
          strokeLinejoin="round"
        />
      </g>
      <circle cx="30" cy="66" r="3" fill="var(--muted-foreground)" fillOpacity="0.3" />
      <circle cx="222" cy="140" r="4" fill="var(--muted-foreground)" fillOpacity="0.2" />
      <circle cx="64" cy="36" r="3" fill="var(--muted-foreground)" fillOpacity="0.35" />
    </svg>
  );
}

export function ConnectPrompt({ onBrowseApps }: { onBrowseApps: () => void }) {
  const navigate = useNavigate();
  const connect = useConnectProvider();
  const { data: gitAppsData } = useGetGitApps();
  const apps = (gitAppsData?.data as GithubApp[]) || [];
  const primaryApp = apps[0];

  const [modalOpen, setModalOpen] = useState(false);
  const [accountName, setAccountName] = useState('');
  const [accessToken, setAccessToken] = useState('');

  const handleConnectToken = async (event: FormEvent) => {
    event.preventDefault();
    if (!accessToken.trim() || connect.isPending) return;
    try {
      await connect.mutateAsync({
        provider: 'github',
        accessToken: accessToken.trim(),
        accountName: accountName.trim() || 'Personal',
      });
      setAccessToken('');
      setAccountName('');
      setModalOpen(false);
      toast.success('GitHub connected successfully');
    } catch {
      toast.error('Failed to connect GitHub');
    }
  };

  const handlePrimaryClick = () => {
    if (primaryApp?.name) {
      window.location.href = `https://github.com/apps/${primaryApp.name}/installations/new`;
      return;
    }
    setModalOpen(true);
  };

  return (
    <div className="rounded-2xl border border-border/50 bg-card">
      <div className="px-6 pt-10 pb-10 text-center">
        <div className="relative mx-auto h-44 w-64">
          <ConnectIllustration />
        </div>
        <h3 className="mt-2 font-medium text-foreground/85 text-lg">Connect GitHub to deploy</h3>
        <p className="mx-auto mt-1.5 mb-7 max-w-md text-muted-foreground text-sm leading-relaxed">
          Connect your GitHub account or install the GitHub App to browse and deploy repositories.
        </p>
        <div className="flex flex-wrap items-center justify-center gap-3">
          <Button
            type="button"
            onClick={handlePrimaryClick}
            disabled={connect.isPending}
            className="h-11 px-6"
          >
            {connect.isPending ? (
              <>
                <Loader2 className="size-4 animate-spin" />
                Connecting...
              </>
            ) : (
              <>
                <GithubIcon className="size-4" />
                {primaryApp?.name ? 'Install GitHub App' : 'Connect GitHub'}
              </>
            )}
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

      <Dialog open={modalOpen} onOpenChange={setModalOpen}>
        <DialogContent className="sm:max-w-md">
          <DialogTitle>Connect GitHub</DialogTitle>
          <DialogDescription>
            Enter a personal access token with repository read permissions to browse and deploy
            repositories.
          </DialogDescription>
          <form onSubmit={handleConnectToken} className="mt-4 space-y-4">
            <div className="space-y-1.5">
              <Label htmlFor="connect-account-name">Account name (optional)</Label>
              <Input
                id="connect-account-name"
                placeholder="Personal"
                value={accountName}
                onChange={(event) => setAccountName(event.target.value)}
              />
            </div>
            <div className="space-y-1.5">
              <Label htmlFor="connect-access-token">Personal access token</Label>
              <Input
                id="connect-access-token"
                type="password"
                placeholder="ghp_..."
                autoComplete="off"
                spellCheck={false}
                value={accessToken}
                onChange={(event) => setAccessToken(event.target.value)}
              />
            </div>
            <div className="flex items-center justify-between pt-2">
              <button
                type="button"
                onClick={() => {
                  setModalOpen(false);
                  void navigate({ to: '/settings', search: { tab: 'git' } as never });
                }}
                className="font-medium text-muted-foreground text-xs hover:text-foreground hover:underline"
              >
                Or configure GitHub App in settings
              </button>
              <Button type="submit" disabled={!accessToken.trim() || connect.isPending}>
                {connect.isPending && <Loader2 className="mr-2 size-4 animate-spin" />}
                Connect
              </Button>
            </div>
          </form>
        </DialogContent>
      </Dialog>
    </div>
  );
}
