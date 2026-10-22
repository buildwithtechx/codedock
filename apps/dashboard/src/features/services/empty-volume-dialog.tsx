import { HardDrive, Loader2, RotateCcw } from 'lucide-react';
import { useId, useState } from 'react';
import { toast } from 'sonner';
import { Button } from '@codedock/ui/button';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@codedock/ui/dialog';
import { Input } from '@codedock/ui/input';
import { Label } from '@codedock/ui/label';
import { apiClient } from '#/lib/api-client';

export interface EmptyVolumeDialogProps {
  serviceId: string;
  serviceName: string;
  volumeId: string;
  hostPath: string;
  containerPath: string;
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSuccess?: () => void;
}

export function EmptyVolumeDialog({
  serviceId,
  serviceName,
  volumeId,
  hostPath,
  containerPath,
  open,
  onOpenChange,
  onSuccess,
}: EmptyVolumeDialogProps) {
  const [confirmation, setConfirmation] = useState('');
  const [busy, setBusy] = useState(false);
  const inputId = useId();

  const isConfirmed = confirmation.trim() === serviceName.trim();

  const handleReset = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!isConfirmed || busy) return;

    setBusy(true);
    try {
      await apiClient.post(`/apps/${serviceId}/volumes/${volumeId}/empty`);
      toast.success('Volume reset to fresh storage. Previous volume retained safely.');
      onOpenChange(false);
      setConfirmation('');
      onSuccess?.();
    } catch (err: any) {
      toast.error(err?.response?.data?.message || err?.message || 'Failed to empty storage volume');
    } finally {
      setBusy(false);
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-md">
        <form onSubmit={handleReset} className="space-y-4">
          <DialogHeader className="space-y-2">
            <div className="flex items-center gap-2">
              <span className="flex size-9 items-center justify-center rounded-xl bg-amber-500/15 text-amber-600 dark:text-amber-400">
                <HardDrive className="size-5" />
              </span>
              <div>
                <DialogTitle className="text-base font-semibold">
                  Start with empty storage
                </DialogTitle>
                <DialogDescription className="text-xs">
                  Attach fresh empty storage while retaining the previous volume.
                </DialogDescription>
              </div>
            </div>
          </DialogHeader>

          <div className="space-y-2 rounded-xl bg-muted/40 p-3.5 text-xs">
            <div className="flex justify-between">
              <span className="text-muted-foreground">Mount path:</span>
              <span className="font-mono font-medium">{containerPath}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-muted-foreground">Retained volume:</span>
              <span className="font-mono font-medium">{hostPath}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-muted-foreground">New volume:</span>
              <span className="font-mono text-emerald-600 dark:text-emerald-400">
                codedock-fresh-new
              </span>
            </div>
          </div>

          <p className="text-muted-foreground text-xs leading-relaxed">
            Your previous volume will not be deleted. It remains safely archived in the host storage
            pool so you can roll back or reattach it anytime.
          </p>

          <div className="space-y-1.5">
            <Label htmlFor={inputId} className="text-xs font-medium">
              Type <span className="font-bold font-mono text-foreground">{serviceName}</span> to
              confirm
            </Label>
            <Input
              id={inputId}
              value={confirmation}
              onChange={(e) => setConfirmation(e.target.value)}
              placeholder={serviceName}
              disabled={busy}
              autoComplete="off"
              className="font-mono text-sm"
            />
          </div>

          <DialogFooter className="flex gap-2 pt-2 sm:justify-end">
            <Button
              type="button"
              variant="outline"
              disabled={busy}
              onClick={() => onOpenChange(false)}
            >
              Cancel
            </Button>
            <Button
              type="submit"
              variant="default"
              disabled={!isConfirmed || busy}
              className="gap-1.5 bg-amber-600 hover:bg-amber-700 text-white"
            >
              {busy ? (
                <Loader2 className="size-4 animate-spin" />
              ) : (
                <RotateCcw className="size-4" />
              )}
              Empty storage
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
