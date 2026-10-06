import { useEffect, useRef, useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '#/components/ui/dialog';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import { ApiError, apiClient } from '#/lib/api-client';

export function VolumeRestoreDialog({
  recordId,
  onClose,
}: {
  recordId: string | null;
  onClose: () => void;
}) {
  const [target, setTarget] = useState<{ volumeName: string; timeoutSeconds: number } | null>(null);
  const [confirmation, setConfirmation] = useState('');
  const [pending, setPending] = useState(false);
  const [error, setError] = useState('');
  const controller = useRef<AbortController | null>(null);
  useEffect(() => {
    setTarget(null);
    setConfirmation('');
    setError('');
    if (!recordId) return;
    const request = new AbortController();
    apiClient
      .get<{ data: { volumeName: string; timeoutSeconds: number } }>(
        `/backup-records/${recordId}/volume-restore`,
        { signal: request.signal }
      )
      .then((response) => {
        if (!request.signal.aborted) setTarget(response.data);
      })
      .catch((err: Error) => {
        if (!request.signal.aborted) setError(err.message);
      });
    return () => request.abort();
  }, [recordId]);
  useEffect(() => () => controller.current?.abort(), []);
  const restore = async () => {
    if (!recordId || !target) return;
    const request = new AbortController();
    controller.current = request;
    setPending(true);
    setError('');
    try {
      await apiClient.post(
        `/backup-records/${recordId}/volume-restore`,
        { volumeName: confirmation, confirmOverwrite: true },
        { signal: request.signal }
      );
      toast.success('Volume restore completed');
      onClose();
    } catch (err) {
      setError(
        request.signal.aborted
          ? 'Restore interrupted. The volume may contain partially restored data. Keep its services stopped until you verify or restore it again.'
          : (err as Error).message
      );
    } finally {
      controller.current = null;
      setPending(false);
    }
  };
  const interrupt = async () => {
    if (!recordId || !controller.current) return;
    const request = controller.current;
    request.abort();
    try {
      await apiClient.delete(`/backup-records/${recordId}/volume-restore`);
      controller.current?.abort();
    } catch (err) {
      if (!(err instanceof ApiError && err.status === 409)) setError((err as Error).message);
    }
  };
  return (
    <Dialog
      open={Boolean(recordId)}
      onOpenChange={(open) => {
        if (!open && !pending) onClose();
      }}
    >
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Restore named volume</DialogTitle>
          <DialogDescription>
            Stop all containers using this volume first. Files in the archive overwrite matching
            files; other files remain. Interruption does not roll back changes.
          </DialogDescription>
        </DialogHeader>
        {error && (
          <p role="alert" className="text-destructive text-sm">
            {error}
          </p>
        )}
        {!target && !error && <p>Validating target...</p>}
        {target && (
          <>
            <p className="text-sm">
              Target: <strong>{target.volumeName}</strong>. Timeout: {target.timeoutSeconds}{' '}
              seconds.
            </p>
            <Label htmlFor="volume-restore-confirmation">
              Type the volume name to confirm overwrite
            </Label>
            <Input
              id="volume-restore-confirmation"
              value={confirmation}
              disabled={pending}
              onChange={(event) => setConfirmation(event.target.value)}
            />
            <div className="flex justify-end gap-2">
              <Button variant="outline" onClick={pending ? interrupt : onClose}>
                {pending ? 'Interrupt restore' : 'Cancel'}
              </Button>
              <Button
                variant="destructive"
                disabled={pending || confirmation !== target.volumeName}
                onClick={restore}
              >
                {pending ? 'Restoring...' : 'Restore volume'}
              </Button>
            </div>
          </>
        )}
      </DialogContent>
    </Dialog>
  );
}
