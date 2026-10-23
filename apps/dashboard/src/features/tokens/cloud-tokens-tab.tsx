import { Cloud, Plus, Trash2 } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Badge } from '@codedock/ui/badge';
import { Button } from '@codedock/ui/button';
import { SettingsSection } from '#/features/settings/settings-section';
import { useDeleteManagedCredential, useManagedCredentials } from '#/hooks/use-managed';
import { useListOrganizations } from '#/hooks/use-organizations';
import { useOrganizationStore } from '#/stores/organization-store';
import { NewCloudTokenDialog } from './new-cloud-token-dialog';

export function CloudTokensTab() {
  const activeOrgId = useOrganizationStore((s) => s.activeOrganizationId);
  const { data: orgs } = useListOrganizations();
  const organizationId = activeOrgId || orgs?.[0]?.id || '';

  const { data: credentials = [], isLoading } = useManagedCredentials(organizationId);
  const deleteMutation = useDeleteManagedCredential();
  const [dialogOpen, setDialogOpen] = useState(false);

  const handleDelete = (credentialId: string, label: string) => {
    if (!organizationId) return;
    deleteMutation.mutate(
      { organizationId, credentialId },
      {
        onSuccess: () => toast.success(`Cloud token "${label}" removed`),
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
        <Button size="sm" onClick={() => setDialogOpen(true)} className="gap-1.5">
          <Plus className="size-3.5" />
          New cloud token
        </Button>
      }
    >
      {isLoading ? (
        <div className="space-y-2">
          {[1, 2].map((i) => (
            <div key={i} className="h-16 animate-pulse rounded-xl bg-muted/40" />
          ))}
        </div>
      ) : credentials.length === 0 ? (
        <div className="rounded-xl border border-border/60 border-dashed py-10 text-center">
          <Cloud className="mx-auto size-8 text-muted-foreground/60" />
          <h3 className="mt-3 font-semibold text-sm">No cloud provider tokens</h3>
          <p className="mt-1 text-xs text-muted-foreground">
            Add a provider token to start auto-provisioning managed cloud servers.
          </p>
        </div>
      ) : (
        <div className="divide-y divide-border/50 rounded-xl border border-border/50">
          {credentials.map((cred) => (
            <div key={cred.id} className="flex items-center justify-between p-4">
              <div className="min-w-0 flex-1 space-y-1">
                <div className="flex items-center gap-2.5">
                  <span className="font-semibold text-foreground text-sm">{cred.label}</span>
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
                className="size-8 text-muted-foreground hover:text-destructive shrink-0"
                aria-label={`Delete ${cred.label}`}
              >
                <Trash2 className="size-4" />
              </Button>
            </div>
          ))}
        </div>
      )}

      <NewCloudTokenDialog
        open={dialogOpen}
        onOpenChange={setDialogOpen}
        organizationId={organizationId}
      />
    </SettingsSection>
  );
}
