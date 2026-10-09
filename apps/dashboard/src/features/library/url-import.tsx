import { Link2 } from 'lucide-react';
import { useId, useState } from 'react';
import { Button } from '#/components/ui/button';
import { Input } from '#/components/ui/input';
import type { PendingImport } from './types';

export function parseGitUrl(raw: string): PendingImport | null {
  const value = raw.trim();
  if (!value) return null;
  const scp = value.match(/^git@([^:]+):([^/]+\/[^/]+?)(\.git)?\/?$/);
  if (scp) {
    const name = scp[2].split('/')[1] ?? scp[2];
    return { repositoryUrl: value, branch: 'main', name };
  }
  try {
    const parsed = new URL(value);
    if (!['https:', 'http:', 'ssh:'].includes(parsed.protocol)) return null;
    if (parsed.username || parsed.password || parsed.port) return null;
    const parts = parsed.pathname.replace(/\/+$/, '').slice(1).split('/');
    if (parts.length < 2) return null;
    const [owner, repoRaw] = parts.slice(-2);
    const repo = (repoRaw ?? '').replace(/\.git$/, '');
    if (!owner || !repo) return null;
    if (!/^[A-Za-z0-9_.-]+$/.test(owner) || !/^[A-Za-z0-9_.-]+$/.test(repo)) return null;
    return { repositoryUrl: value, branch: 'main', name: repo };
  } catch {
    return null;
  }
}

export function UrlImport({ onImport }: { onImport: (pending: PendingImport) => void }) {
  const [url, setUrl] = useState('');
  const [error, setError] = useState('');
  const errorId = useId();

  const handleSubmit = (event: React.FormEvent) => {
    event.preventDefault();
    const parsed = parseGitUrl(url);
    if (!parsed) {
      setError('Enter a valid repository URL, e.g. https://github.com/owner/repo');
      return;
    }
    setError('');
    onImport(parsed);
  };

  return (
    <div className="rounded-2xl border border-border/50 bg-card">
      <div className="p-8">
        <div className="mx-auto max-w-lg">
          <div className="mx-auto mb-4 flex h-14 w-14 items-center justify-center rounded-2xl bg-foreground/6">
            <Link2 className="size-7 text-muted-foreground" />
          </div>
          <h3 className="mb-1.5 text-center font-semibold text-base text-foreground">
            Import from URL
          </h3>
          <p className="mb-6 text-center text-muted-foreground text-sm leading-relaxed">
            Deploy any public repository, or a private one your connected providers can read.
          </p>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <Input
                type="url"
                autoFocus
                aria-label="Repository URL"
                aria-invalid={!!error}
                aria-describedby={error ? errorId : undefined}
                value={url}
                onChange={(event) => {
                  setUrl(event.target.value);
                  setError('');
                }}
                placeholder="https://github.com/owner/repository"
                className={error ? 'ring-2 ring-destructive/40' : undefined}
              />
              {error && (
                <p id={errorId} role="alert" className="mt-1.5 text-destructive text-xs">
                  {error}
                </p>
              )}
            </div>
            <Button type="submit" disabled={!url.trim()} className="h-11 w-full">
              Continue
            </Button>
          </form>
        </div>
      </div>
    </div>
  );
}
