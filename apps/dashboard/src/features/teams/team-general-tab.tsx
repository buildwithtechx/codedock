import { Building2, Calendar, Check, Copy, Hash, Save } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Button } from '@codedock/ui/button';
import { Input } from '@codedock/ui/input';
import { Label } from '@codedock/ui/label';
import { SettingsSection } from '#/features/settings/settings-section';
import { useUpdateOrganization } from '#/hooks/use-organizations';
import type { Organization } from '#/interfaces/organizations';

interface TeamGeneralTabProps {
  organization: Organization;
}

export function TeamGeneralTab({ organization }: TeamGeneralTabProps) {
  const [name, setName] = useState(organization.name);
  const [copied, setCopied] = useState(false);
  const updateMutation = useUpdateOrganization();

  const handleCopyId = () => {
    void navigator.clipboard.writeText(organization.id);
    setCopied(true);
    toast.success('Workspace ID copied to clipboard');
    setTimeout(() => setCopied(false), 2000);
  };

  const handleSave = (e?: React.FormEvent) => {
    e?.preventDefault();
    if (!name.trim()) {
      toast.error('Workspace name cannot be empty');
      return;
    }
    updateMutation.mutate(
      { id: organization.id, name: name.trim() },
      {
        onSuccess: () => toast.success('Workspace name updated successfully'),
        onError: (err) =>
          toast.error(err instanceof Error ? err.message : 'Failed to update workspace name'),
      }
    );
  };

  const isChanged = name.trim() !== organization.name && name.trim().length > 0;

  return (
    <div className="space-y-6">
      <SettingsSection
        icon={<Building2 className="size-4 text-primary" />}
        title="Workspace Profile"
        description="The display name and identity for your team workspace."
        action={
          <Button
            size="sm"
            onClick={() => handleSave()}
            disabled={!isChanged || updateMutation.isPending}
            className="gap-1.5"
          >
            <Save className="size-3.5" />
            <span>{updateMutation.isPending ? 'Saving...' : 'Save changes'}</span>
          </Button>
        }
      >
        <form onSubmit={handleSave} className="space-y-4">
          <div className="space-y-1.5 max-w-md">
            <Label htmlFor="workspace-name" className="text-xs font-medium">
              Workspace Name
            </Label>
            <Input
              id="workspace-name"
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="My Team"
              className="bg-muted/30"
            />
            <p className="text-[11px] text-muted-foreground">
              Shown across the dashboard, projects, and notifications.
            </p>
          </div>
        </form>
      </SettingsSection>

      <SettingsSection
        icon={<Hash className="size-4 text-primary" />}
        title="Workspace Properties"
        description="Technical identifiers, unique slug, and registration details."
      >
        <div className="space-y-4 divide-y divide-border/40">
          <div className="flex flex-col gap-2 pt-1 sm:flex-row sm:items-center sm:justify-between">
            <div>
              <p className="font-medium text-foreground text-xs">Workspace ID</p>
              <p className="text-muted-foreground text-[11px]">
                Unique database identifier used for API requests and CLI scoping.
              </p>
            </div>
            <div className="flex items-center gap-2">
              <code className="rounded bg-muted/60 px-2.5 py-1 font-mono text-xs text-foreground">
                {organization.id}
              </code>
              <Button
                variant="outline"
                size="sm"
                onClick={handleCopyId}
                className="gap-1.5 h-8 text-xs"
              >
                {copied ? (
                  <Check className="size-3.5 text-emerald-500" />
                ) : (
                  <Copy className="size-3.5" />
                )}
                <span>{copied ? 'Copied' : 'Copy'}</span>
              </Button>
            </div>
          </div>

          <div className="flex flex-col gap-2 pt-4 sm:flex-row sm:items-center sm:justify-between">
            <div className="flex items-center gap-2">
              <Calendar className="size-4 text-muted-foreground" />
              <div>
                <p className="font-medium text-foreground text-xs">Created At</p>
                <p className="text-muted-foreground text-[11px]">
                  Date this workspace was originally provisioned.
                </p>
              </div>
            </div>
            <span className="text-xs text-muted-foreground font-medium">
              {organization.createdAt
                ? new Date(organization.createdAt).toLocaleDateString(undefined, {
                    year: 'numeric',
                    month: 'long',
                    day: 'numeric',
                  })
                : 'Active workspace'}
            </span>
          </div>
        </div>
      </SettingsSection>
    </div>
  );
}
