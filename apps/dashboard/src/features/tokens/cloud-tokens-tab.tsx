import { Cloud, Plus, Trash2 } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Badge } from '@codedock/ui/badge';
import { Button } from '@codedock/ui/button';
import { Input } from '@codedock/ui/input';
import { Label } from '@codedock/ui/label';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@codedock/ui/select';
import { SettingsSection } from '#/features/settings/settings-section';
import {
  useDeleteManagedCredential,
  useManagedCredentials,
  useCreateManagedCredential,
} from '#/hooks/use-managed';
import { useListOrganizations } from '#/hooks/use-organizations';
import { useOrganizationStore } from '#/stores/organization-store';

const CLOUD_PROVIDERS = [
  { id: 'hetzner', name: 'Hetzner Cloud', placeholder: 'Hetzner Cloud API token' },
  { id: 'aws', name: 'Amazon Web Services (AWS)', placeholder: 'AWS Access Key / Secret token' },
  { id: 'digitalocean', name: 'DigitalOcean', placeholder: 'DigitalOcean Personal Access Token' },
  { id: 'scaleway', name: 'Scaleway', placeholder: 'Scaleway API secret key' },
];

export function CloudTokensTab() {
  const activeOrgId = useOrganizationStore((s) => s.activeOrganizationId);
  const { data: orgs } = useListOrganizations();
  const organizationId = activeOrgId || orgs?.[0]?.id || '';

  const { data: credentials = [], isLoading } = useManagedCredentials(organizationId);
  const createMutation = useCreateManagedCredential();
  const deleteMutation = useDeleteManagedCredential();

  const [showForm, setShowForm] = useState(false);
  const [provider, setProvider] = useState('hetzner');
  const [label, setLabel] = useState('');
  const [token, setToken] = useState('');

  const activeProviderMeta = CLOUD_PROVIDERS.find((p) => p.id === provider) || CLOUD_PROVIDERS[0];

  const resetForm = () => {
    setShowForm(false);
    setLabel('');
    setToken('');
  };

  const handleCreate = (e: React.FormEvent) => {
    e.preventDefault();
    if (!label.trim()) {
      toast.error('Token label is required');
      return;
    }
    if (!token.trim()) {
      toast.error('API token is required');
      return;
    }
    if (!organizationId) {
      toast.error('No organization selected');
      return;
    }

    createMutation.mutate(
      {
        organizationId,
        payload: {
          provider,
          label: label.trim(),
          token: token.trim(),
        },
      },
      {
        onSuccess: () => {
          toast.success(`Cloud token "${label}" added successfully`);
          resetForm();
        },
        onError: (err) => {
          toast.error(err instanceof Error ? err.message : 'Failed to save cloud token');
        },
      }
    );
  };

  const handleDelete = (credentialId: string, itemLabel: string) => {
    if (!organizationId) return;
    deleteMutation.mutate(
      { organizationId, credentialId },
      {
        onSuccess: () => toast.success(`Cloud token "${itemLabel}" removed`),
        onError: (err) =>
          toast.error(err instanceof Error ? err.message : 'Failed to delete token'),
      }
    );
  };

  return (
    <SettingsSection
      icon={<Cloud className="size-4 text-primary" />}
      title="Cloud Tokens"
      description="Provider credentials used to auto-provision, resize, and manage cloud instances with Hetzner, AWS, and DigitalOcean."
      action={
        !showForm ? (
          <Button size="sm" onClick={() => setShowForm(true)} className="gap-1.5">
            <Plus className="size-3.5" />
            <span>New cloud token</span>
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
            Add Provider Credential
          </p>

          <div className="grid gap-3 sm:grid-cols-2">
            <div className="space-y-1">
              <Label htmlFor="inline-cloud-provider" className="text-xs">
                Provider
              </Label>
              <Select value={provider} onValueChange={(val) => setProvider(val)}>
                <SelectTrigger id="inline-cloud-provider" className="bg-muted/40 text-xs h-9">
                  <SelectValue placeholder="Select provider" />
                </SelectTrigger>
                <SelectContent>
                  {CLOUD_PROVIDERS.map((p) => (
                    <SelectItem key={p.id} value={p.id}>
                      {p.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>

            <div className="space-y-1">
              <Label htmlFor="inline-cloud-label" className="text-xs">
                Label / Account Name
              </Label>
              <Input
                id="inline-cloud-label"
                value={label}
                onChange={(e) => setLabel(e.target.value)}
                placeholder="e.g. Production Hetzner"
                className="bg-muted/40 text-xs h-9"
                required
              />
            </div>
          </div>

          <div className="space-y-1">
            <Label htmlFor="inline-cloud-token" className="text-xs">
              API Token / Secret Key
            </Label>
            <Input
              id="inline-cloud-token"
              type="password"
              value={token}
              onChange={(e) => setToken(e.target.value)}
              placeholder={activeProviderMeta.placeholder}
              className="bg-muted/40 font-mono text-xs h-9"
              required
            />
            <p className="text-[11px] text-muted-foreground">
              Encrypted at rest using AES-256 GCM in the backend database.
            </p>
          </div>

          <div className="flex justify-end gap-2 pt-1 border-t border-border/40">
            <Button type="button" variant="ghost" size="sm" onClick={resetForm}>
              Cancel
            </Button>
            <Button type="submit" size="sm" disabled={createMutation.isPending}>
              {createMutation.isPending ? 'Adding...' : 'Add Cloud Token'}
            </Button>
          </div>
        </form>
      )}

      {isLoading ? (
        <div className="space-y-2">
          {[1, 2].map((i) => (
            <div key={i} className="h-14 animate-pulse rounded-xl bg-muted/40" />
          ))}
        </div>
      ) : credentials.length === 0 ? (
        <div className="rounded-xl border border-border/60 border-dashed py-8 text-center text-muted-foreground text-xs">
          No cloud provider tokens yet. Add a provider token to start auto-provisioning managed
          servers.
        </div>
      ) : (
        <div className="divide-y divide-border/50 rounded-xl border border-border/50">
          {credentials.map((cred) => (
            <div key={cred.id} className="flex items-center justify-between p-3.5">
              <div className="min-w-0 flex-1 space-y-1">
                <div className="flex items-center gap-2">
                  <span className="font-semibold text-foreground text-xs">{cred.label}</span>
                  <Badge variant="outline" className="text-[10px] capitalize">
                    {cred.provider}
                  </Badge>
                </div>
                <div className="flex items-center gap-3 text-xs text-muted-foreground">
                  <span>
                    Token: <code className="font-mono">••••••••••••••••</code>
                  </span>
                  {cred.createdAt && (
                    <span>Added {new Date(cred.createdAt).toLocaleDateString()}</span>
                  )}
                </div>
              </div>

              <Button
                variant="ghost"
                size="icon"
                onClick={() => handleDelete(cred.id, cred.label)}
                disabled={deleteMutation.isPending}
                className="size-7 text-muted-foreground hover:text-destructive shrink-0"
                aria-label={`Delete ${cred.label}`}
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
