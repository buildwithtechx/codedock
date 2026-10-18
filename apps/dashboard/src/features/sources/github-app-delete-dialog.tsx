import { Trash } from 'lucide-react';
import { Button } from '@codedock/ui/button';
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogTitle,
} from '@codedock/ui/dialog';
interface GithubAppDeleteDialogProps {
  deletingApp: string | null;
  onClose: () => void;
  onConfirm: () => void;
  isDeleting: boolean;
}

export function GithubAppDeleteDialog({
  deletingApp,
  onClose,
  onConfirm,
  isDeleting,
}: GithubAppDeleteDialogProps) {
  return (
    <Dialog open={!!deletingApp} onOpenChange={(open) => !open && onClose()}>
      <DialogContent className="border/50 gap-0 bg-card/95 p-0 backdrop-blur-xl sm:max-w-md [&>button]:hidden">
        <div className="p-5">
          <div className="flex items-start justify-between">
            <div className="flex flex-col">
              <DialogTitle className="flex items-center gap-2 font-bold text-destructive text-xl tracking-tight">
                <Trash className="h-5 w-5" />
                Remove GitHub App
              </DialogTitle>
              <DialogDescription>This will break existing deployments</DialogDescription>
            </div>
            <DialogClose asChild>
              <Button
                variant="ghost"
                className="font-medium text-foreground/80 text-sm hover:bg-transparent hover:text-foreground"
              >
                CLOSE
              </Button>
            </DialogClose>
          </div>
        </div>

        <div className="flex items-center justify-end gap-3 p-5 pt-0">
          <Button
            variant="ghost"
            onClick={onClose}
            className="h-9 font-mono font-semibold text-[11px] uppercase tracking-wider"
          >
            Cancel
          </Button>
          <Button
            onClick={(e) => {
              e.preventDefault();
              onConfirm();
            }}
            disabled={isDeleting}
            variant="destructive"
            className="h-9 gap-2 font-mono font-semibold text-[11px] uppercase tracking-wider"
          >
            <Trash className="h-3.5 w-3.5" />
            {isDeleting ? 'Removing...' : 'Remove App'}
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  );
}
