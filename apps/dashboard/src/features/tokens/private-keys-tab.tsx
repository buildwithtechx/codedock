import { Check, Copy, KeyRound, Plus, Trash2 } from 'lucide-react';
import { useEffect, useState } from 'react';
import { toast } from 'sonner';
import { Button } from '@codedock/ui/button';
import { Input } from '@codedock/ui/input';
import { Label } from '@codedock/ui/label';
import { Textarea } from '@codedock/ui/textarea';
import { SettingsSection } from '#/features/settings/settings-section';

export interface SSHPrivateKey {
  id: string;
  name: string;
  description?: string;
  publicKey: string;
  fingerprint: string;
  createdAt: string;
}

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

  const [showForm, setShowForm] = useState(false);
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [privateKey, setPrivateKey] = useState('');
  const [publicKey, setPublicKey] = useState('');
  const [copiedId, setCopiedId] = useState<string | null>(null);

  useEffect(() => {
    try {
      localStorage.setItem(STORAGE_KEY, JSON.stringify(keys));
    } catch {
      // Storage unavailable
    }
  }, [keys]);

  const resetForm = () => {
    setShowForm(false);
    setName('');
    setDescription('');
    setPrivateKey('');
    setPublicKey('');
  };

  const handleGenerate = () => {
    const mockHash = Array.from({ length: 8 }, () =>
      Math.floor(Math.random() * 256)
        .toString(16)
        .padStart(2, '0')
    ).join(':');
    const mockPub = `ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAI${btoa(mockHash).replace(/=/g, '')} codedock@instance`;
    setPublicKey(mockPub);
    setPrivateKey(
      '-----BEGIN OPENSSH PRIVATE KEY-----\nb3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAMwAAAAtzc2gt\nZWQyNTUxOQAAAC...[AUTOMATICALLY GENERATED DEPLOY KEY]...\n-----END OPENSSH PRIVATE KEY-----'
    );
    toast.success('Generated new key pair');
  };

  const handleCreate = (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim()) {
      toast.error('Key name is required');
      return;
    }
    if (!privateKey.trim()) {
      toast.error('Private key content is required');
      return;
    }

    const fingerprint = publicKey.trim()
      ? `SHA256:${btoa(name).slice(0, 16)}...`
      : 'SHA256:imported_key';

    const derivedPublic =
      publicKey.trim() || `ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAA... ${name}@codedock`;

    const newKey: SSHPrivateKey = {
      id: `key-${Date.now()}`,
      name: name.trim(),
      description: description.trim() || undefined,
      publicKey: derivedPublic,
      fingerprint,
      createdAt: new Date().toISOString(),
    };

    setKeys((prev) => [newKey, ...prev]);
    toast.success(`Private key "${name}" added`);
    resetForm();
  };

  const handleDeleteKey = (id: string, keyName: string) => {
    setKeys((prev) => prev.filter((k) => k.id !== id));
    toast.success(`Removed private key "${keyName}"`);
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
      description="SSH private keys used for server connection authentication and private Git repository access."
      action={
        !showForm ? (
          <Button size="sm" onClick={() => setShowForm(true)} className="gap-1.5">
            <Plus className="size-3.5" />
            <span>New private key</span>
          </Button>
        ) : undefined
      }
    >
      {showForm && (
        <form
          onSubmit={handleCreate}
          className="mb-5 space-y-4 rounded-xl border border-border/50 bg-muted/20 p-4"
        >
          <div className="flex items-center justify-between border-b border-border/40 pb-2">
            <p className="font-semibold text-xs text-foreground">Add New Private Key</p>
            <button
              type="button"
              onClick={handleGenerate}
              className="font-medium text-xs text-primary hover:underline"
            >
              Auto-generate key pair
            </button>
          </div>

          <div className="grid gap-3 sm:grid-cols-2">
            <div className="space-y-1">
              <Label htmlFor="inline-key-name" className="text-xs">
                Name
              </Label>
              <Input
                id="inline-key-name"
                value={name}
                onChange={(e) => setName(e.target.value)}
                placeholder="e.g. Production Hetzner Key"
                className="bg-muted/40 text-xs"
                required
              />
            </div>
            <div className="space-y-1">
              <Label htmlFor="inline-key-desc" className="text-xs">
                Description (optional)
              </Label>
              <Input
                id="inline-key-desc"
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                placeholder="e.g. Used for root ssh access"
                className="bg-muted/40 text-xs"
              />
            </div>
          </div>

          <div className="space-y-1">
            <Label htmlFor="inline-priv-key" className="text-xs">
              Private Key (OpenSSH format)
            </Label>
            <Textarea
              id="inline-priv-key"
              rows={4}
              value={privateKey}
              onChange={(e) => setPrivateKey(e.target.value)}
              placeholder="-----BEGIN OPENSSH PRIVATE KEY-----&#10;...&#10;-----END OPENSSH PRIVATE KEY-----"
              className="bg-muted/40 font-mono text-xs"
              required
            />
          </div>

          <div className="space-y-1">
            <Label htmlFor="inline-pub-key" className="text-xs">
              Public Key (derived or custom)
            </Label>
            <Input
              id="inline-pub-key"
              value={publicKey}
              onChange={(e) => setPublicKey(e.target.value)}
              placeholder="ssh-ed25519 AAAAC3NzaC1... or ssh-rsa AAAAB3..."
              className="bg-muted/40 font-mono text-xs"
            />
          </div>

          <div className="flex justify-end gap-2 pt-1 border-t border-border/40">
            <Button type="button" variant="ghost" size="sm" onClick={resetForm}>
              Cancel
            </Button>
            <Button type="submit" size="sm">
              Save Private Key
            </Button>
          </div>
        </form>
      )}

      {keys.length === 0 ? (
        <div className="rounded-xl border border-border/60 border-dashed py-8 text-center text-muted-foreground text-xs">
          No private keys found. Add an SSH private key to connect to destination servers.
        </div>
      ) : (
        <div className="divide-y divide-border/50 rounded-xl border border-border/50">
          {keys.map((key) => {
            const isCopied = copiedId === key.id;
            return (
              <div
                key={key.id}
                className="flex flex-col gap-3 p-3.5 sm:flex-row sm:items-center sm:justify-between"
              >
                <div className="min-w-0 flex-1 space-y-1">
                  <div className="flex items-center gap-2">
                    <span className="font-semibold text-foreground text-xs">{key.name}</span>
                    <span className="rounded bg-muted px-1.5 py-0.5 font-mono text-[10px] text-muted-foreground">
                      {key.fingerprint}
                    </span>
                  </div>
                  {key.description && (
                    <p className="text-[11px] text-muted-foreground">{key.description}</p>
                  )}
                  <div className="flex items-center gap-2 pt-0.5">
                    <code className="truncate max-w-sm font-mono text-[10px] text-muted-foreground bg-muted/40 px-1.5 py-0.5 rounded">
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
                      <span>{isCopied ? 'Copied' : 'Copy'}</span>
                    </button>
                  </div>
                </div>

                <Button
                  variant="ghost"
                  size="icon"
                  onClick={() => handleDeleteKey(key.id, key.name)}
                  className="size-7 text-muted-foreground hover:text-destructive shrink-0"
                  aria-label={`Delete ${key.name}`}
                >
                  <Trash2 className="size-3.5" />
                </Button>
              </div>
            );
          })}
        </div>
      )}
    </SettingsSection>
  );
}
