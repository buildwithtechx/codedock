import { CheckCircle2, KeyRound, Loader2, Plus } from 'lucide-react';
import { useState } from 'react';
import { Button } from '@codedock/ui/button';
import { Input } from '@codedock/ui/input';
import { Label } from '@codedock/ui/label';
import type { ManagedCredential } from '#/interfaces/managed';
import { cn } from '#/lib/utils';

export function HetznerCredentialStep({
  credentials,
  selectedId,
  loading,
  creating,
  createError,
  onSelect,
  onCreate,
}: {
  credentials: ManagedCredential[];
  selectedId: string | null;
  loading: boolean;
  creating: boolean;
  createError: string | null;
  onSelect: (credentialId: string) => void;
  onCreate: (label: string, token: string) => void;
}) {
  const [showForm, setShowForm] = useState(false);
  const [label, setLabel] = useState('');
  const [token, setToken] = useState('');

  if (loading) {
    return (
      <div className="flex items-center gap-3 rounded-xl border border-border/60 p-6 text-muted-foreground text-sm">
        <Loader2 className="size-4 animate-spin" />
        Loading Hetzner credentials…
      </div>
    );
  }

  return (
    <div className="space-y-4">
      <div>
        <h3 className="font-medium text-foreground text-sm">Hetzner API credential</h3>
        <p className="mt-1 text-muted-foreground text-xs">
          Pick the token Codedock uses against Hetzner Cloud. New tokens are verified with a live
          API probe before they are stored.
        </p>
      </div>

      {credentials.length > 0 && (
        <div className="space-y-2">
          {credentials.map((credential) => {
            const selected = credential.id === selectedId;
            return (
              <button
                key={credential.id}
                type="button"
                onClick={() => onSelect(credential.id)}
                className={cn(
                  'flex w-full items-center gap-3 rounded-xl border px-4 py-3 text-left transition-colors',
                  selected
                    ? 'border-primary bg-primary/5'
                    : 'border-border/70 hover:border-primary/30 hover:bg-muted/30'
                )}
              >
                <KeyRound
                  className={cn(
                    'size-4 shrink-0',
                    selected ? 'text-primary' : 'text-muted-foreground'
                  )}
                />
                <span className="min-w-0 flex-1">
                  <span className="block truncate font-medium text-foreground text-sm">
                    {credential.label}
                  </span>
                  <span className="block text-muted-foreground text-xs">
                    {credential.provider} · added {credential.createdAt.slice(0, 10)}
                  </span>
                </span>
                {selected && <CheckCircle2 className="size-4 shrink-0 text-primary" />}
              </button>
            );
          })}
        </div>
      )}

      {showForm ? (
        <div className="space-y-4 rounded-xl border border-border/70 p-4">
          <div className="space-y-1.5">
            <Label htmlFor="hetzner-credential-label" className="text-xs">
              Label
            </Label>
            <Input
              id="hetzner-credential-label"
              value={label}
              onChange={(e) => setLabel(e.target.value)}
              placeholder="production-hetzner"
              disabled={creating}
            />
          </div>
          <div className="space-y-1.5">
            <Label htmlFor="hetzner-credential-token" className="text-xs">
              Hetzner Cloud API token
            </Label>
            <Input
              id="hetzner-credential-token"
              type="password"
              value={token}
              onChange={(e) => setToken(e.target.value)}
              placeholder="Paste a read/write Cloud API token"
              disabled={creating}
              autoComplete="off"
            />
          </div>
          {createError && <p className="text-destructive text-xs">{createError}</p>}
          <div className="flex items-center justify-end gap-2">
            <Button
              type="button"
              variant="ghost"
              size="sm"
              onClick={() => setShowForm(false)}
              disabled={creating}
            >
              Cancel
            </Button>
            <Button
              type="button"
              size="sm"
              onClick={() => onCreate(label.trim(), token.trim())}
              disabled={creating || !label.trim() || !token.trim()}
            >
              {creating && <Loader2 className="size-3.5 animate-spin" />}
              Verify & save
            </Button>
          </div>
        </div>
      ) : (
        <Button
          type="button"
          variant="outline"
          size="sm"
          onClick={() => setShowForm(true)}
          className="gap-2"
        >
          <Plus className="size-3.5" />
          {credentials.length === 0 ? 'Add Hetzner credential' : 'Use a different token'}
        </Button>
      )}
    </div>
  );
}
