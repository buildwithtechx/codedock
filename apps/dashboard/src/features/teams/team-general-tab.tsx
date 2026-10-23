import { Copy, Check, Building2, Calendar, Hash, Save } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Button } from '@codedock/ui/button';
import { Input } from '@codedock/ui/input';
import { Label } from '@codedock/ui/label';
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

  const handleSave = (e: React.FormEvent) => {
    e.preventDefault();
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

  return (
    <div className="space-y-6">
      <div className="rounded-2xl border border-border/60 bg-card p-6 shadow-sm">
        <div className="flex items-start gap-4">
          <div className="flex size-12 shrink-0 items-center justify-center rounded-2xl bg-primary/10 text-primary">
            <Building2 className="size-6" />
          </div>
          <div className="min-w-0 flex-1">
            <h2 className="font-semibold text-foreground text-lg">{organization.name}</h2>
            <p className="text-muted-foreground text-sm">
              General settings, identity, and properties for this workspace.
            </p>
          </div>
        </div>

        <form onSubmit={handleSave} className="mt-6 space-y-4 border-t border-border/50 pt-5">
          <div className="space-y-2">
            <Label htmlFor="team-name" className="text-sm font-medium">
              Workspace Name
            </Label>
            <div className="flex max-w-md items-center gap-3">
              <Input
                id="team-name"
                value={name}
                onChange={(e) => setName(e.target.value)}
                placeholder="My Team"
                className="h-10"
              />
              <Button
                type="submit"
                size="sm"
                disabled={updateMutation.isPending || name.trim() === organization.name}
                className="h-10 gap-1.5 shrink-0"
              >
                <Save className="size-4" />
                {updateMutation.isPending ? 'Saving...' : 'Save'}
              </Button>
            </div>
          </div>
        </form>
      </div>

      <div className="grid gap-4 sm:grid-cols-2">
        <div className="rounded-2xl border border-border/60 bg-card p-5">
          <div className="flex items-center gap-3">
            <div className="flex size-9 items-center justify-center rounded-xl bg-muted text-muted-foreground">
              <Hash className="size-4" />
            </div>
            <div className="min-w-0 flex-1">
              <p className="text-xs font-medium text-muted-foreground uppercase tracking-wider">
                Workspace ID
              </p>
              <div className="mt-1 flex items-center gap-2">
                <code className="truncate font-mono text-xs text-foreground bg-muted/60 px-2 py-0.5 rounded">
                  {organization.id}
                </code>
                <button
                  type="button"
                  onClick={handleCopyId}
                  className="rounded p-1 text-muted-foreground hover:bg-muted hover:text-foreground transition-colors"
                  aria-label="Copy workspace ID"
                >
                  {copied ? (
                    <Check className="size-3.5 text-emerald-500" />
                  ) : (
                    <Copy className="size-3.5" />
                  )}
                </button>
              </div>
            </div>
          </div>
        </div>

        <div className="rounded-2xl border border-border/60 bg-card p-5">
          <div className="flex items-center gap-3">
            <div className="flex size-9 items-center justify-center rounded-xl bg-muted text-muted-foreground">
              <Calendar className="size-4" />
            </div>
            <div className="min-w-0 flex-1">
              <p className="text-xs font-medium text-muted-foreground uppercase tracking-wider">
                Created At
              </p>
              <p className="mt-1 text-sm font-medium text-foreground">
                {organization.createdAt
                  ? new Date(organization.createdAt).toLocaleDateString(undefined, {
                      year: 'numeric',
                      month: 'long',
                      day: 'numeric',
                    })
                  : 'Active workspace'}
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
