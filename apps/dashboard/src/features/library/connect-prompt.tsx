import { useNavigate } from '@tanstack/react-router';
import { BookOpen, ExternalLink, KeyRound, LayoutGrid, Loader2 } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import { useConnectProvider } from './hooks';
import { GIT_PROVIDERS } from './types';

const DOCS_URL = 'https://docs.codedock.run';

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
      <path
        d="M118 88a10 10 0 1 0 3 7.3V92h-4v-6h6a4 4 0 1 1 3 4v10h-8"
        stroke="var(--primary)"
        strokeWidth="2.5"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
      <circle cx="30" cy="66" r="3" fill="var(--muted-foreground)" fillOpacity="0.3" />
      <circle cx="222" cy="140" r="4" fill="var(--muted-foreground)" fillOpacity="0.2" />
      <circle cx="64" cy="36" r="3" fill="var(--muted-foreground)" fillOpacity="0.35" />
    </svg>
  );
}

export function ConnectPrompt({ onBrowseApps }: { onBrowseApps: () => void }) {
  const navigate = useNavigate();
  const connect = useConnectProvider();
  const [selected, setSelected] = useState<string | null>(null);
  const [accountName, setAccountName] = useState('');
  const [accessToken, setAccessToken] = useState('');

  const handleConnect = async (provider: string) => {
    if (!accessToken.trim() || connect.isPending) return;
    try {
      await connect.mutateAsync({
        provider,
        accessToken: accessToken.trim(),
        accountName: accountName.trim() || 'Personal',
      });
      setAccessToken('');
      setAccountName('');
      setSelected(null);
      toast.success('Git provider connected');
    } catch {
      toast.error('Failed to connect git provider');
    }
  };

  return (
    <div className="rounded-2xl border border-border/50 bg-card">
      <div className="px-6 pt-10 pb-10 text-center">
        <div className="relative mx-auto w-fit">
          <ConnectIllustration />
        </div>
        <h3 className="mt-2 font-medium text-foreground/85 text-lg">Connect a git provider</h3>
        <p className="mx-auto mt-1.5 mb-7 max-w-md text-muted-foreground text-sm leading-relaxed">
          Connect with a personal access token to browse your repositories and deploy them to
          Codedock.
        </p>
        <div className="mx-auto grid w-full max-w-xl gap-3 text-start sm:grid-cols-2">
          {GIT_PROVIDERS.map((provider) => {
            const open = selected === provider.id;
            return (
              <div
                key={provider.id}
                className={`rounded-xl border p-4 transition-all ${
                  open
                    ? 'border-primary/40 bg-primary/[0.03]'
                    : 'border-border/60 hover:border-primary/40'
                }`}
              >
                <button
                  type="button"
                  onClick={() => setSelected(open ? null : provider.id)}
                  className="flex w-full items-center gap-3 text-start"
                  aria-expanded={open}
                >
                  <span className="flex size-9 items-center justify-center rounded-lg bg-muted">
                    <img
                      src={provider.icon}
                      alt=""
                      aria-hidden="true"
                      className="size-5 object-contain"
                    />
                  </span>
                  <span>
                    <span className="block font-medium text-foreground text-sm">
                      {provider.name}
                    </span>
                    <span className="block text-muted-foreground text-xs">
                      Personal access token
                    </span>
                  </span>
                  <KeyRound className="ml-auto size-4 shrink-0 text-muted-foreground" />
                </button>
                {open && (
                  <div className="mt-3 space-y-3">
                    <div className="space-y-1.5">
                      <Label htmlFor={`connect-account-${provider.id}`}>Account name</Label>
                      <Input
                        id={`connect-account-${provider.id}`}
                        placeholder="Personal"
                        value={accountName}
                        onChange={(event) => setAccountName(event.target.value)}
                      />
                    </div>
                    <div className="space-y-1.5">
                      <Label htmlFor={`connect-token-${provider.id}`}>Access token</Label>
                      <Input
                        id={`connect-token-${provider.id}`}
                        type="password"
                        autoComplete="off"
                        spellCheck={false}
                        placeholder="Token with repo read access"
                        value={accessToken}
                        onChange={(event) => setAccessToken(event.target.value)}
                        onKeyDown={(event) => {
                          if (event.key === 'Enter') void handleConnect(provider.id);
                        }}
                      />
                    </div>
                    <Button
                      className="w-full"
                      disabled={!accessToken.trim() || connect.isPending}
                      onClick={() => void handleConnect(provider.id)}
                    >
                      {connect.isPending && <Loader2 className="size-4 animate-spin" />}
                      Connect {provider.name}
                    </Button>
                  </div>
                )}
              </div>
            );
          })}
        </div>
        <div className="mt-7 flex flex-wrap items-center justify-center gap-3">
          <Button type="button" variant="secondary" onClick={onBrowseApps}>
            <LayoutGrid className="size-4" />
            Deploy an app instead
          </Button>
          <Button
            type="button"
            variant="ghost"
            onClick={() => void navigate({ to: '/settings', search: { tab: 'sources' } as never })}
          >
            Manage in settings
          </Button>
          <Button type="button" variant="secondary" asChild>
            <a href={DOCS_URL} target="_blank" rel="noopener noreferrer">
              <BookOpen className="size-4" />
              Docs
              <ExternalLink className="size-3.5 opacity-60" />
            </a>
          </Button>
        </div>
      </div>
    </div>
  );
}
