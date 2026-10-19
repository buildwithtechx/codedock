import { Loader2 } from 'lucide-react';
import { Button } from '@codedock/ui/button';
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogTitle,
} from '@codedock/ui/dialog';
import { Input } from '@codedock/ui/input';
import type { Passkey } from '#/interfaces/passkey';

interface AddDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  name: string;
  onNameChange: (val: string) => void;
  onSubmit: () => void;
  isPending: boolean;
}

export function PasskeyAddDialog({
  open,
  onOpenChange,
  name,
  onNameChange,
  onSubmit,
  isPending,
}: AddDialogProps) {
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="border/50 gap-4 bg-card/95 backdrop-blur-xl sm:max-w-md">
        <div className="flex items-start justify-between">
          <div>
            <DialogTitle className="font-bold text-lg">Add New Passkey</DialogTitle>
            <DialogDescription className="text-xs">
              Give your passkey a name to identify this device.
            </DialogDescription>
          </div>
          <DialogClose asChild>
            <Button type="button" variant="ghost" size="sm">
              Close
            </Button>
          </DialogClose>
        </div>
        <Input
          placeholder="e.g. MacBook Pro, YubiKey 5C"
          value={name}
          onChange={(e) => onNameChange(e.target.value)}
          onKeyDown={(e) => e.key === 'Enter' && onSubmit()}
        />
        <div className="flex justify-end gap-2">
          <Button type="button" variant="outline" size="sm" onClick={() => onOpenChange(false)}>
            Cancel
          </Button>
          <Button type="button" size="sm" onClick={onSubmit} disabled={isPending}>
            {isPending && <Loader2 className="mr-1.5 h-3.5 w-3.5 animate-spin" />}
            Continue
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  );
}

interface RenameDialogProps {
  target: Passkey | null;
  onClose: () => void;
  name: string;
  onNameChange: (val: string) => void;
  onSubmit: () => void;
  isPending: boolean;
}

export function PasskeyRenameDialog({
  target,
  onClose,
  name,
  onNameChange,
  onSubmit,
  isPending,
}: RenameDialogProps) {
  return (
    <Dialog open={!!target} onOpenChange={(open) => !open && onClose()}>
      <DialogContent className="border/50 gap-4 bg-card/95 backdrop-blur-xl sm:max-w-md">
        <div className="flex items-start justify-between">
          <div>
            <DialogTitle className="font-bold text-lg">Rename Passkey</DialogTitle>
            <DialogDescription className="text-xs">
              Update the friendly label for this passkey.
            </DialogDescription>
          </div>
          <DialogClose asChild>
            <Button type="button" variant="ghost" size="sm">
              Close
            </Button>
          </DialogClose>
        </div>
        <Input
          value={name}
          onChange={(e) => onNameChange(e.target.value)}
          onKeyDown={(e) => e.key === 'Enter' && onSubmit()}
        />
        <div className="flex justify-end gap-2">
          <Button type="button" variant="outline" size="sm" onClick={onClose}>
            Cancel
          </Button>
          <Button type="button" size="sm" onClick={onSubmit} disabled={isPending || !name.trim()}>
            {isPending && <Loader2 className="mr-1.5 h-3.5 w-3.5 animate-spin" />}
            Save
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  );
}

interface DeleteDialogProps {
  target: Passkey | null;
  onClose: () => void;
  onSubmit: () => void;
  isPending: boolean;
}

export function PasskeyDeleteDialog({ target, onClose, onSubmit, isPending }: DeleteDialogProps) {
  return (
    <Dialog open={!!target} onOpenChange={(open) => !open && onClose()}>
      <DialogContent className="border/50 gap-4 bg-card/95 backdrop-blur-xl sm:max-w-md">
        <div>
          <DialogTitle className="font-bold text-destructive text-lg">Delete Passkey</DialogTitle>
          <DialogDescription className="mt-1 text-xs">
            Are you sure you want to remove &quot;{target?.name}&quot;? You will no longer be able
            to sign in with this passkey.
          </DialogDescription>
        </div>
        <div className="flex justify-end gap-2 pt-2">
          <Button type="button" variant="outline" size="sm" onClick={onClose}>
            Cancel
          </Button>
          <Button
            type="button"
            variant="destructive"
            size="sm"
            onClick={onSubmit}
            disabled={isPending}
          >
            {isPending && <Loader2 className="mr-1.5 h-3.5 w-3.5 animate-spin" />}
            Delete
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  );
}
