import { Building2, Calendar, Check, Copy, Hash } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Button } from '@codedock/ui/button';
import { Input } from '@codedock/ui/input';
import { SettingsRow } from '#/features/settings/settings-row';
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
  const createdDate = organization.createdAt
    ? new Date(organization.createdAt).toLocaleDateString(undefined, {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
      })
    : 'Active workspace';

  return (
    <div className="space-y-6">
      {/* Workspace Profile (styled exactly like Domains & Routing, open permanently) */}
      <SettingsSection
        icon={<Building2 className="h-4 w-4" />}
        title="Workspace Profile"
        description="Configure your team workspace name and display settings."
      >
        <SettingsRow
          label="Workspace Name"
          description="Shown across the dashboard, projects, and notifications."
          bordered={false}
        >
          <Input
            id="workspace-name"
            aria-label="Workspace Name"
            value={name}
            onChange={(e) => setName(e.target.value)}
            placeholder="Workspace Name"
            className="text-sm"
          />
        </SettingsRow>

        <div className="flex justify-end border-border/40 border-t pt-4">
          <Button
            size="sm"
            onClick={handleSave}
            disabled={!isChanged || updateMutation.isPending}
            className="gap-1.5"
          >
            <Check className="h-3.5 w-3.5" />
            <span>{updateMutation.isPending ? 'Saving...' : 'Save Workspace Profile'}</span>
          </Button>
        </div>
      </SettingsSection>

      {/* Workspace Properties (styled like Instance Overview) */}
      <SettingsSection
        icon={<Hash className="size-4 text-violet-500" />}
        iconBg="bg-violet-500/10"
        iconColor="text-violet-500"
        title="Workspace Properties"
        description="Technical identifiers, unique slug, and registration details."
      >
        <div className="grid grid-cols-1 gap-3 sm:grid-cols-2">
          <div className="flex items-center justify-between gap-3 rounded-xl border border-border/60 bg-muted/20 p-4">
            <div className="flex min-w-0 items-center gap-3">
              <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-muted">
                <Hash className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="min-w-0">
                <p className="truncate font-medium text-foreground text-sm font-mono">
                  {organization.id}
                </p>
                <p className="text-muted-foreground text-xs">Workspace ID</p>
              </div>
            </div>
            <Button
              variant="outline"
              size="sm"
              onClick={handleCopyId}
              className="gap-1.5 h-8 text-xs shrink-0"
            >
              {copied ? (
                <Check className="size-3.5 text-emerald-500" />
              ) : (
                <Copy className="size-3.5" />
              )}
              <span>{copied ? 'Copied' : 'Copy'}</span>
            </Button>
          </div>

          <div className="flex items-center justify-between gap-3 rounded-xl border border-border/60 bg-muted/20 p-4">
            <div className="flex min-w-0 items-center gap-3">
              <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-muted">
                <Calendar className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="min-w-0">
                <p className="font-medium text-foreground text-sm">{createdDate}</p>
                <p className="text-muted-foreground text-xs">Creation Date</p>
              </div>
            </div>
            <span className="inline-flex shrink-0 items-center gap-1 rounded-full bg-emerald-500/10 px-2 py-0.5 font-medium text-[11px] text-emerald-600 dark:text-emerald-400">
              Active
            </span>
          </div>
        </div>
      </SettingsSection>
    </div>
  );
}
