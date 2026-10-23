import { Blocks, Plus, Trash2 } from 'lucide-react';
import { useState, useEffect } from 'react';
import { toast } from 'sonner';
import { Badge } from '@codedock/ui/badge';
import { Button } from '@codedock/ui/button';
import { SettingsSection } from '#/features/settings/settings-section';
import { NewIntegrationTokenDialog, type IntegrationToken } from './new-integration-token-dialog';

const STORAGE_KEY = 'codedock_integration_tokens_store';

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

  const [dialogOpen, setDialogOpen] = useState(false);

  useEffect(() => {
    try {
      localStorage.setItem(STORAGE_KEY, JSON.stringify(tokens));
    } catch {
      // Storage unavailable
    }
  }, [tokens]);

  const handleAdd = (newToken: IntegrationToken) => {
    setTokens((prev) => [newToken, ...prev]);
  };

  const handleDelete = (id: string, name: string) => {
    setTokens((prev) => prev.filter((t) => t.id !== id));
    toast.success(`Integration token "${name}" deleted`);
  };

  return (
    <SettingsSection
      icon={<Blocks className="size-4 text-primary" />}
      title="Integration Tokens"
      description="Tokens and credentials used to authenticate with external registries, CI pipelines, and notification channels."
      action={
        <Button size="sm" onClick={() => setDialogOpen(true)} className="gap-1.5">
          <Plus className="size-3.5" />
          New integration token
        </Button>
      }
    >
      {tokens.length === 0 ? (
        <div className="rounded-xl border border-border/60 border-dashed py-10 text-center">
          <Blocks className="mx-auto size-8 text-muted-foreground/60" />
          <h3 className="mt-3 font-semibold text-sm">No integration tokens</h3>
          <p className="mt-1 text-xs text-muted-foreground">
            Connect external Docker registries, webhook alerts, or monitoring providers.
          </p>
        </div>
      ) : (
        <div className="divide-y divide-border/50 rounded-xl border border-border/50">
          {tokens.map((token) => (
            <div key={token.id} className="flex items-center justify-between p-4">
              <div className="min-w-0 flex-1 space-y-1">
                <div className="flex items-center gap-2.5">
                  <span className="font-semibold text-foreground text-sm">{token.name}</span>
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
                className="size-8 text-muted-foreground hover:text-destructive shrink-0"
                aria-label={`Delete ${token.name}`}
              >
                <Trash2 className="size-4" />
              </Button>
            </div>
          ))}
        </div>
      )}

      <NewIntegrationTokenDialog
        open={dialogOpen}
        onOpenChange={setDialogOpen}
        onAddToken={handleAdd}
      />
    </SettingsSection>
  );
}
