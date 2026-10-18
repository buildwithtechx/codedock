import { Check } from 'lucide-react';
import type { FormEvent } from 'react';
import { Button } from '@codedock/ui/button';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@codedock/ui/dialog';
import { Input } from '@codedock/ui/input';
import { Label } from '@codedock/ui/label';

interface ProjectGitConnectDialogProps {
  connecting: string | null;
  onClose: () => void;
  accountName: string;
  onAccountNameChange: (val: string) => void;
  accessToken: string;
  onAccessTokenChange: (val: string) => void;
  onSubmit: (e: FormEvent) => void;
  isPending: boolean;
}

export function ProjectGitConnectDialog({
  connecting,
  onClose,
  accountName,
  onAccountNameChange,
  accessToken,
  onAccessTokenChange,
  onSubmit,
  isPending,
}: ProjectGitConnectDialogProps) {
  return (
    <Dialog open={connecting !== null} onOpenChange={(open) => !open && onClose()}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>Connect {connecting}</DialogTitle>
          <DialogDescription>Enter a personal access token for this git account.</DialogDescription>
        </DialogHeader>
        <form onSubmit={onSubmit} className="space-y-4">
          <div className="space-y-2">
            <Label>Account name</Label>
            <Input
              value={accountName}
              onChange={(event) => onAccountNameChange(event.target.value)}
              placeholder="Personal"
            />
          </div>
          <div className="space-y-2">
            <Label>Personal access token</Label>
            <Input
              type="password"
              value={accessToken}
              onChange={(event) => onAccessTokenChange(event.target.value)}
              placeholder="Token with repository read access"
              required
              className="font-mono"
            />
          </div>
          <div className="flex justify-end gap-2">
            <Button type="button" variant="ghost" onClick={onClose}>
              Cancel
            </Button>
            <Button type="submit" disabled={isPending}>
              <Check className="h-3.5 w-3.5" />
              {isPending ? 'Connecting' : 'Connect'}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
}
