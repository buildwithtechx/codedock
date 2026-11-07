import { Blocks, Plus, Trash2 } from 'lucide-react';
import { useEffect, useState } from 'react';
import { toast } from 'sonner';
import { Badge } from '@codedock/ui/badge';
import { Button } from '@codedock/ui/button';
import { Input } from '@codedock/ui/input';
import { Label } from '@codedock/ui/label';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@codedock/ui/select';
import { SettingsSection } from '#/features/settings/settings-section';

export interface IntegrationToken {
  id: string;
  name: string;
  provider: string;
  maskedToken: string;
  createdAt: string;
}

const STORAGE_KEY = 'codedock_integration_tokens_store';

const INTEGRATION_PROVIDERS = [
  { id: 'docker', name: 'Docker Hub / Container Registry', placeholder: 'dckr_pat_...' },
  { id: 'github', name: 'GitHub Personal Access Token', placeholder: 'ghp_...' },
  { id: 'gitlab', name: 'GitLab Personal Access Token', placeholder: 'glpat-...' },
  { id: 'sentry', name: 'Sentry Auth Token', placeholder: 'sntrys_...' },
  {
    id: 'webhook',
    name: 'Slack / Discord Webhook URL',
    placeholder: 'https://hooks.slack.com/services/...',
  },
];

const DEFAULT_INTEGRATIONS: IntegrationToken[] = [
  {
    id: 'int-docker-1',
    name: 'Docker Hub Registry',
    provider: 'Docker Hub / Container Registry',
    maskedToken: 'dckr••••91a2',
    createdAt: new Date().toISOString(),
  },
];

export function IntegrationTokensTab() {
  const [tokens, setTokens] = useState<IntegrationToken[]>(() => {
    if (typeof window === 'undefined') return DEFAULT_INTEGRATIONS;
    try {
      const stored = localStorage.getItem(STORAGE_KEY);
      return stored ? JSON.parse(stored) : DEFAULT_INTEGRATIONS;
    } catch {
      return DEFAULT_INTEGRATIONS;
    }
  });

  const [showForm, setShowForm] = useState(false);
  const [provider, setProvider] = useState('docker');
  const [name, setName] = useState('');
  const [tokenValue, setTokenValue] = useState('');

  const activeProvider =
    INTEGRATION_PROVIDERS.find((p) => p.id === provider) || INTEGRATION_PROVIDERS[0];

  useEffect(() => {
    try {
      localStorage.setItem(STORAGE_KEY, JSON.stringify(tokens));
    } catch {
      // Storage unavailable
    }
  }, [tokens]);

  const resetForm = () => {
    setShowForm(false);
    setName('');
    setTokenValue('');
  };

  const handleCreate = (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim()) {
      toast.error('Token name is required');
      return;
    }
    if (!tokenValue.trim()) {
      toast.error('Token or credential value is required');
      return;
    }

    const masked =
      tokenValue.length > 8 ? `${tokenValue.slice(0, 4)}••••${tokenValue.slice(-4)}` : '••••••••';

    const newItem: IntegrationToken = {
      id: `int-${Date.now()}`,
      name: name.trim(),
      provider: activeProvider.name,
      maskedToken: masked,
      createdAt: new Date().toISOString(),
    };

    setTokens((prev) => [newItem, ...prev]);
    toast.success(`Integration token "${name}" saved`);
    resetForm();
  };

  const handleDelete = (id: string, tokenName: string) => {
    setTokens((prev) => prev.filter((t) => t.id !== id));
    toast.success(`Integration token "${tokenName}" deleted`);
  };

  return (
    <SettingsSection
      icon={<Blocks className="size-4 text-primary" />}
      title="Integration Tokens"
      description="Tokens and credentials used to authenticate with external registries, CI pipelines, and notification channels."
      action={
        !showForm ? (
          <Button size="sm" onClick={() => setShowForm(true)} className="gap-1.5">
            <Plus className="size-3.5" />
            <span>New integration token</span>
          </Button>
        ) : undefined
      }
    >
      {showForm && (
        <form
          onSubmit={handleCreate}
          className="mb-5 space-y-4 rounded-xl border border-border/50 bg-muted/20 p-4"
        >
          <p className="font-semibold text-xs text-foreground border-b border-border/40 pb-2">
            Add Integration Credential
          </p>

          <div className="grid gap-3 sm:grid-cols-2">
            <div className="space-y-1">
              <Label htmlFor="inline-int-provider" className="text-xs">
                Integration Type
              </Label>
              <Select value={provider} onValueChange={setProvider}>
                <SelectTrigger id="inline-int-provider" className="bg-muted/40 text-xs h-9">
                  <SelectValue placeholder="Select type" />
                </SelectTrigger>
                <SelectContent>
                  {INTEGRATION_PROVIDERS.map((p) => (
                    <SelectItem key={p.id} value={p.id}>
                      {p.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>

            <div className="space-y-1">
              <Label htmlFor="inline-int-name" className="text-xs">
                Name
              </Label>
              <Input
                id="inline-int-name"
                value={name}
                onChange={(e) => setName(e.target.value)}
                placeholder="e.g. Production Docker Hub Registry"
                className="bg-muted/40 text-xs h-9"
                required
              />
            </div>
          </div>

          <div className="space-y-1">
            <Label htmlFor="inline-int-val" className="text-xs">
              Secret / Access Token
            </Label>
            <Input
              id="inline-int-val"
              type="password"
              value={tokenValue}
              onChange={(e) => setTokenValue(e.target.value)}
              placeholder={activeProvider.placeholder}
              className="bg-muted/40 font-mono text-xs h-9"
              required
            />
          </div>

          <div className="flex justify-end gap-2 pt-1 border-t border-border/40">
            <Button type="button" variant="ghost" size="sm" onClick={resetForm}>
              Cancel
            </Button>
            <Button type="submit" size="sm">
              Save Integration Token
            </Button>
          </div>
        </form>
      )}

      {tokens.length === 0 ? (
        <div className="rounded-xl border border-border/60 border-dashed py-8 text-center text-muted-foreground text-xs">
          No integration tokens yet. Connect external Docker registries, webhook alerts, or
          monitoring providers.
        </div>
      ) : (
        <div className="divide-y divide-border/50 rounded-xl border border-border/50">
          {tokens.map((token) => (
            <div key={token.id} className="flex items-center justify-between p-3.5">
              <div className="min-w-0 flex-1 space-y-1">
                <div className="flex items-center gap-2">
                  <span className="font-semibold text-foreground text-xs">{token.name}</span>
                  <Badge variant="secondary" className="text-[10px]">
                    {token.provider}
                  </Badge>
                </div>
                <div className="flex items-center gap-3 text-xs text-muted-foreground">
                  <span>
                    Secret: <code className="font-mono">{token.maskedToken}</code>
                  </span>
                  <span>Added {new Date(token.createdAt).toLocaleDateString()}</span>
                </div>
              </div>

              <Button
                variant="ghost"
                size="icon"
                onClick={() => handleDelete(token.id, token.name)}
                className="size-7 text-muted-foreground hover:text-destructive shrink-0"
                aria-label={`Delete ${token.name}`}
              >
                <Trash2 className="size-3.5" />
              </Button>
            </div>
          ))}
        </div>
      )}
    </SettingsSection>
  );
}
