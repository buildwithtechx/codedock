import { useNavigate } from '@tanstack/react-router';
import { Trash2 } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Button } from '@codedock/ui/button';
import { Dialog, DialogContent, DialogTitle } from '@codedock/ui/dialog';
import { Input } from '@codedock/ui/input';
import { Label } from '@codedock/ui/label';
import { useDeleteProject } from '#/hooks/use-projects';

export interface ProjectDeleteDialogProps {
  isOpen: boolean;
  onOpenChange: (open: boolean) => void;
  project: {
    id: string;
    name: string;
  };
  onDeleted?: () => void;
  redirectToProjectsOnDelete?: boolean;
}

export function ProjectDeleteDialog({
  isOpen,
  onOpenChange,
  project,
  onDeleted,
  redirectToProjectsOnDelete = false,
}: ProjectDeleteDialogProps) {
  const navigate = useNavigate();
  const [confirmName, setConfirmName] = useState('');
  const deleteProjectMutation = useDeleteProject();

  const isConfirmed = confirmName.trim() === project.name.trim();

  const handleClose = () => {
    if (deleteProjectMutation.isPending) return;
    setConfirmName('');
    onOpenChange(false);
  };

  const handleDelete = async () => {
    if (!isConfirmed) return;
    try {
      await deleteProjectMutation.mutateAsync({ id: project.id });
      toast.success(`Project "${project.name}" was deleted`);
      setConfirmName('');
      onOpenChange(false);
      onDeleted?.();
      if (redirectToProjectsOnDelete) {
        void navigate({ to: '/projects' });
      }
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : 'Failed to delete project';
      toast.error(message);
    }
  };

  return (
    <Dialog open={isOpen} onOpenChange={(open) => !open && handleClose()}>
      <DialogContent className="border/50 gap-0 bg-card/95 p-0 backdrop-blur-xl sm:max-w-md [&>button]:hidden">
        <div className="p-5">
          <DialogTitle className="flex items-center gap-2 font-bold text-foreground text-xl tracking-tight">
            <Trash2 className="h-5 w-5 text-destructive" />
            Delete Project
          </DialogTitle>

          <p className="mt-3 text-muted-foreground text-sm leading-relaxed">
            Are you sure you want to delete{' '}
            <strong className="text-foreground">{project.name}</strong>? All associated services,
            containers, logs, and configurations will be immediately removed.
          </p>

          <div className="mt-4 rounded-xl border border-destructive/30 bg-destructive/5 p-3 text-destructive text-xs">
            This action is permanent and cannot be reversed.
          </div>

          <div className="mt-4 space-y-2">
            <Label htmlFor="delete-confirm-input" className="text-muted-foreground text-xs">
              To confirm, type <strong className="font-mono text-foreground">{project.name}</strong>
              :
            </Label>
            <Input
              id="delete-confirm-input"
              value={confirmName}
              onChange={(e) => setConfirmName(e.target.value)}
              placeholder={project.name}
              disabled={deleteProjectMutation.isPending}
              className="bg-muted/40 font-mono text-xs"
              autoComplete="off"
              autoFocus
            />
          </div>
        </div>

        <div className="flex items-center justify-end gap-3 p-5 pt-0">
          <Button
            type="button"
            variant="ghost"
            size="sm"
            onClick={handleClose}
            disabled={deleteProjectMutation.isPending}
          >
            Cancel
          </Button>
          <Button
            type="button"
            variant="destructive"
            size="sm"
            onClick={() => void handleDelete()}
            disabled={!isConfirmed || deleteProjectMutation.isPending}
          >
            {deleteProjectMutation.isPending ? 'Deleting...' : 'Delete Project'}
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  );
}
