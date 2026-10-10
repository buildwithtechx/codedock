import { useNavigate } from '@tanstack/react-router';
import { AlertTriangle, Loader2 } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '#/components/ui/dialog';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import { useDeleteProject } from './hooks';

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
      toast.success(`Project "${project.name}" was deleted successfully`);
      setConfirmName('');
      onOpenChange(false);
      if (onDeleted) {
        onDeleted();
      }
      if (redirectToProjectsOnDelete) {
        void navigate({ to: '/projects' });
      }
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : 'Failed to delete project';
      toast.error(message);
    }
  };

  return (
    <Dialog open={isOpen} onOpenChange={handleClose}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader icon={AlertTriangle}>
          <DialogTitle>Delete project</DialogTitle>
          <DialogDescription>
            This action cannot be undone. All services, containers, logs, and configurations
            associated with <span className="font-semibold text-foreground">{project.name}</span>{' '}
            will be permanently removed.
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-4 py-2">
          <div className="rounded-xl border border-destructive/20 bg-destructive/10 p-3 text-destructive text-xs leading-relaxed">
            Please be certain. Once deleted, this project cannot be recovered.
          </div>

          <div className="space-y-2">
            <Label htmlFor="delete-confirm-input" className="text-xs">
              To confirm, type <span className="font-semibold text-foreground">{project.name}</span>{' '}
              below:
            </Label>
            <Input
              id="delete-confirm-input"
              value={confirmName}
              onChange={(e) => setConfirmName(e.target.value)}
              placeholder={project.name}
              disabled={deleteProjectMutation.isPending}
              className="text-sm"
              autoComplete="off"
            />
          </div>
        </div>

        <DialogFooter className="gap-2 sm:gap-0">
          <Button
            type="button"
            variant="outline"
            onClick={handleClose}
            disabled={deleteProjectMutation.isPending}
          >
            Cancel
          </Button>
          <Button
            type="button"
            variant="destructive"
            onClick={() => void handleDelete()}
            disabled={!isConfirmed || deleteProjectMutation.isPending}
            className="gap-2"
          >
            {deleteProjectMutation.isPending ? (
              <>
                <Loader2 className="size-4 animate-spin" />
                Deleting…
              </>
            ) : (
              'Delete project'
            )}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
