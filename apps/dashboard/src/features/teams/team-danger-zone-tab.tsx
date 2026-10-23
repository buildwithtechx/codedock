import { AlertTriangle, Trash2 } from 'lucide-react';
import { useState } from 'react';
import { useNavigate } from '@tanstack/react-router';
import { toast } from 'sonner';
import { Button } from '@codedock/ui/button';
import { Input } from '@codedock/ui/input';
import { Label } from '@codedock/ui/label';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@codedock/ui/dialog';
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

  return (
    <div className="space-y-6">
      <div className="rounded-2xl border border-destructive/30 bg-destructive/5 p-6">
        <div className="flex items-start gap-4">
          <div className="flex size-11 shrink-0 items-center justify-center rounded-2xl bg-destructive/15 text-destructive">
            <AlertTriangle className="size-5" />
          </div>
          <div className="min-w-0 flex-1">
            <h2 className="font-semibold text-foreground text-lg">Danger Zone</h2>
            <p className="text-muted-foreground text-sm">
              Irreversible destructive actions for the workspace{' '}
              <span className="font-semibold text-foreground">{organization.name}</span>.
            </p>
          </div>
        </div>

        <div className="mt-6 flex flex-col items-start justify-between gap-4 border-t border-destructive/20 pt-5 sm:flex-row sm:items-center">
          <div>
            <h3 className="font-medium text-foreground text-sm">Delete this workspace</h3>
            <p className="max-w-xl text-muted-foreground text-xs">
              Permanently removes this workspace, projects, associated credentials, and team
              members. This action cannot be undone.
            </p>
          </div>
          <Button
            variant="destructive"
            size="sm"
            onClick={() => setOpenDelete(true)}
            className="gap-2 shrink-0"
          >
            <Trash2 className="size-4" />
            Delete Workspace
          </Button>
        </div>
      </div>

      <Dialog open={openDelete} onOpenChange={setOpenDelete}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle className="text-destructive flex items-center gap-2">
              <AlertTriangle className="size-5" />
              Delete Workspace
            </DialogTitle>
            <DialogDescription className="space-y-2 pt-2">
              <span>
                Are you sure you want to delete{' '}
                <strong className="text-foreground">{organization.name}</strong>? All associated
                services, domains, and team members will be removed.
              </span>
            </DialogDescription>
          </DialogHeader>

          <div className="space-y-2 py-3">
            <Label htmlFor="confirm-org-name" className="text-xs">
              Type <strong className="text-foreground">{organization.name}</strong> to confirm:
            </Label>
            <Input
              id="confirm-org-name"
              value={confirmName}
              onChange={(e) => setConfirmName(e.target.value)}
              placeholder={organization.name}
              className="font-mono text-sm"
            />
          </div>

          <DialogFooter className="gap-2 sm:gap-0">
            <Button
              type="button"
              variant="outline"
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
              disabled={confirmName !== organization.name || deleteMutation.isPending}
              onClick={handleDelete}
            >
              {deleteMutation.isPending ? 'Deleting...' : 'Delete Workspace'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
