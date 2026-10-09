import { useMutation, useQueryClient } from '@tanstack/react-query';
import { Trash2 } from 'lucide-react';
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
import type { Server } from '#/interfaces/server';
import { removeServer } from './api';

interface ServerDeleteDialogProps {
  server: Server | null;
  onClose: () => void;
  onRemoved?: (server: Server) => void;
}

export function ServerDeleteDialog({ server, onClose, onRemoved }: ServerDeleteDialogProps) {
  const queryClient = useQueryClient();
  const deletion = useMutation({
    mutationFn: (id: string) => removeServer(id),
    onSuccess: (_, id) => {
      queryClient.invalidateQueries({ queryKey: ['servers'] });
      queryClient.invalidateQueries({ queryKey: ['servers', id] });
    },
  });

  const handleDelete = async () => {
    if (!server) return;
    try {
      await deletion.mutateAsync(server.id);
      toast.success(`Server ${server.name} removed`);
      onRemoved?.(server);
      onClose();
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Failed to remove server');
    }
  };

  return (
    <Dialog open={server !== null} onOpenChange={(open) => !open && onClose()}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Remove {server?.name ?? 'server'}?</DialogTitle>
          <DialogDescription>
            This removes the server from Codedock. Running workloads on the node are not affected.
          </DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <Button variant="outline" onClick={onClose}>
            Cancel
          </Button>
          <Button
            variant="destructive"
            onClick={handleDelete}
            disabled={deletion.isPending}
            className="gap-1.5"
          >
            <Trash2 className="h-4 w-4" />
            {deletion.isPending ? 'Removing…' : 'Remove server'}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
