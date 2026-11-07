import { AlertTriangle, Trash2 } from 'lucide-react';
import { useState } from 'react';
import { useNavigate } from '@tanstack/react-router';
import { toast } from 'sonner';
import { Button } from '@codedock/ui/button';
import { Input } from '@codedock/ui/input';
import { Label } from '@codedock/ui/label';
import { Dialog, DialogContent, DialogTitle } from '@codedock/ui/dialog';
import { SettingsSection } from '#/features/settings/settings-section';
import { useDeleteOrganization } from '#/hooks/use-organizations';
import type { Organization } from '#/interfaces/organizations';

interface TeamDangerZoneTabProps {
  organization: Organization;
}

export function TeamDangerZoneTab({ organization }: TeamDangerZoneTabProps) {
  const [openDelete, setOpenDelete] = useState(false);
  const [confirmName, setConfirmName] = useState('');
  const deleteMutation = useDeleteOrganization();
  const navigate = useNavigate();

  const handleDelete = () => {
    if (confirmName !== organization.name) {
      toast.error('Workspace name does not match');
      return;
    }
    deleteMutation.mutate(organization.id, {
      onSuccess: () => {
        toast.success(`Workspace "${organization.name}" was deleted`);
        setOpenDelete(false);
        void navigate({ to: '/', replace: true });
      },
      onError: (err) => {
        toast.error(err instanceof Error ? err.message : 'Failed to delete workspace');
      },
    });
  };

  const isMatched = confirmName.trim() === organization.name;

  return (
    <div className="space-y-6">
      <SettingsSection
        icon={<AlertTriangle className="size-4 text-destructive" />}
        title="Danger Zone"
        description="Irreversible destructive actions for this workspace."
        iconBg="bg-destructive/10"
        iconColor="text-destructive"
      >
        <div className="flex flex-col items-start justify-between gap-4 sm:flex-row sm:items-center">
          <div className="space-y-1">
            <h3 className="font-medium text-foreground text-sm">Delete this workspace</h3>
            <p className="max-w-xl text-muted-foreground text-xs leading-relaxed">
              Permanently removes{' '}
              <span className="font-semibold text-foreground">{organization.name}</span>, associated
              servers, projects, domain routing, and teammate access. This action cannot be undone.
            </p>
          </div>
          <Button
            variant="destructive"
            size="sm"
            onClick={() => {
              setConfirmName('');
              setOpenDelete(true);
            }}
            className="gap-2 shrink-0"
          >
            <Trash2 className="size-4" />
            Delete workspace
          </Button>
        </div>
      </SettingsSection>

      <Dialog open={openDelete} onOpenChange={(open) => !open && setOpenDelete(false)}>
        <DialogContent className="border/50 gap-0 bg-card/95 p-0 backdrop-blur-xl sm:max-w-md [&>button]:hidden">
          <div className="p-5">
            <div className="flex items-start justify-between">
              <div className="flex flex-col">
                <DialogTitle className="flex items-center gap-2 font-bold text-foreground text-xl tracking-tight">
                  <Trash2 className="h-5 w-5 text-destructive" />
                  Delete Workspace
                </DialogTitle>
              </div>
            </div>

            <p className="mt-3 text-muted-foreground text-sm leading-relaxed">
              Are you sure you want to delete{' '}
              <strong className="text-foreground">{organization.name}</strong>? All associated
              services, deployments, and databases will be immediately removed.
            </p>

            <div className="mt-4 rounded-xl border border-destructive/30 bg-destructive/5 p-3 text-xs text-destructive">
              This action is permanent and cannot be reversed.
            </div>

            <div className="mt-4 space-y-2">
              <Label htmlFor="confirm-workspace-input" className="text-xs text-muted-foreground">
                To confirm, type{' '}
                <strong className="font-mono text-foreground">{organization.name}</strong>:
              </Label>
              <Input
                id="confirm-workspace-input"
                value={confirmName}
                onChange={(e) => setConfirmName(e.target.value)}
                placeholder={organization.name}
                className="font-mono text-xs bg-muted/40"
                autoFocus
              />
            </div>
          </div>

          <div className="flex items-center justify-end gap-3 p-5 pt-0">
            <Button
              type="button"
              variant="ghost"
              size="sm"
              onClick={() => {
                setOpenDelete(false);
                setConfirmName('');
              }}
            >
              Cancel
            </Button>
            <Button
              type="button"
              variant="destructive"
              size="sm"
              disabled={!isMatched || deleteMutation.isPending}
              onClick={handleDelete}
            >
              {deleteMutation.isPending ? 'Deleting...' : 'Delete Workspace'}
            </Button>
          </div>
        </DialogContent>
      </Dialog>
    </div>
  );
}
