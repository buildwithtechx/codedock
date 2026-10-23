import { Check, Copy, KeyRound, Plus, Trash2 } from 'lucide-react';
import { useState, useEffect } from 'react';
import { toast } from 'sonner';
import { Button } from '@codedock/ui/button';
import { SettingsSection } from '#/features/settings/settings-section';
import { NewPrivateKeyDialog, type SSHPrivateKey } from './new-private-key-dialog';

const STORAGE_KEY = 'codedock_private_keys_store';

const DEFAULT_KEYS: SSHPrivateKey[] = [
  {
    id: 'key-default-1',
    name: 'Default Deploy Key',
    description: 'Auto-provisioned ED25519 key for Git deployments and server joins',
    fingerprint: 'SHA256:7uG+K5L9pQ8zXm1VbN0w...',
    publicKey: 'ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIFrGq8wK9zLpQm1VbN0w codedock@platform',
    createdAt: new Date().toISOString(),
  },
];

export function PrivateKeysTab() {
  const [keys, setKeys] = useState<SSHPrivateKey[]>(() => {
    if (typeof window === 'undefined') return DEFAULT_KEYS;
    try {
      const stored = localStorage.getItem(STORAGE_KEY);
      return stored ? JSON.parse(stored) : DEFAULT_KEYS;
    } catch {
      return DEFAULT_KEYS;
    }
  });

  const [dialogOpen, setDialogOpen] = useState(false);
  const [copiedId, setCopiedId] = useState<string | null>(null);

  useEffect(() => {
    try {
      localStorage.setItem(STORAGE_KEY, JSON.stringify(keys));
    } catch {
      // Storage unavailable
    }
  }, [keys]);

  const handleAddKey = (newKey: SSHPrivateKey) => {
    setKeys((prev) => [newKey, ...prev]);
  };

  const handleDeleteKey = (id: string, name: string) => {
    setKeys((prev) => prev.filter((k) => k.id !== id));
    toast.success(`Removed private key "${name}"`);
  };

  const handleCopyPublicKey = (key: SSHPrivateKey) => {
    void navigator.clipboard.writeText(key.publicKey);
    setCopiedId(key.id);
    toast.success('Public key copied to clipboard');
    setTimeout(() => setCopiedId(null), 2000);
  };

  return (
    <SettingsSection
      icon={<KeyRound className="size-4 text-primary" />}
      title="Private Keys"
      description="Manage SSH private keys used for server connection authentication and private Git repository access."
      action={
        <Button size="sm" onClick={() => setDialogOpen(true)} className="gap-1.5">
          <Plus className="size-3.5" />
          New private key
        </Button>
      }
    >
      {keys.length === 0 ? (
        <div className="rounded-xl border border-border/60 border-dashed py-10 text-center">
          <KeyRound className="mx-auto size-8 text-muted-foreground/60" />
          <h3 className="mt-3 font-semibold text-sm">No private keys found</h3>
          <p className="mt-1 text-xs text-muted-foreground">
            Add an SSH private key to connect to destination servers and deploy repositories.
          </p>
        </div>
      ) : (
        <div className="divide-y divide-border/50 rounded-xl border border-border/50">
          {keys.map((key) => {
            const isCopied = copiedId === key.id;
            return (
              <div
                key={key.id}
                className="flex flex-col gap-3 p-4 sm:flex-row sm:items-center sm:justify-between"
              >
                <div className="min-w-0 flex-1 space-y-1">
                  <div className="flex items-center gap-2.5">
                    <span className="font-semibold text-foreground text-sm">{key.name}</span>
                    <span className="rounded bg-muted px-2 py-0.5 font-mono text-[10px] text-muted-foreground">
                      {key.fingerprint}
                    </span>
                  </div>
                  {key.description && (
                    <p className="text-xs text-muted-foreground">{key.description}</p>
                  )}
                  <div className="flex items-center gap-2 pt-0.5">
                    <code className="truncate max-w-sm font-mono text-[11px] text-muted-foreground/80 bg-muted/40 px-1.5 py-0.5 rounded">
                      {key.publicKey}
                    </code>
                    <button
                      type="button"
                      onClick={() => handleCopyPublicKey(key)}
                      className="inline-flex items-center gap-1 text-[11px] font-medium text-primary hover:underline"
                    >
                      {isCopied ? (
                        <Check className="size-3 text-emerald-500" />
                      ) : (
                        <Copy className="size-3" />
                      )}
                      <span>{isCopied ? 'Copied' : 'Copy Public Key'}</span>
                    </button>
                  </div>
                </div>

                <Button
                  variant="ghost"
                  size="icon"
                  onClick={() => handleDeleteKey(key.id, key.name)}
                  className="size-8 text-muted-foreground hover:text-destructive shrink-0"
                  aria-label={`Delete ${key.name}`}
                >
                  <Trash2 className="size-4" />
                </Button>
              </div>
            );
          })}
        </div>
      )}

      <NewPrivateKeyDialog open={dialogOpen} onOpenChange={setDialogOpen} onAddKey={handleAddKey} />
    </SettingsSection>
  );
}
