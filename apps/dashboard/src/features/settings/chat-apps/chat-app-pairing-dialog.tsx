import {
  Button,
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@codedock/ui';
import { Check, Copy, ShieldCheck } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { usePairingChallenge, useVerifyPairing } from '#/hooks/use-chat-apps';

interface ChatAppPairingDialogProps {
  token: string | null;
  onClose: () => void;
}

export function ChatAppPairingDialog({ token, onClose }: ChatAppPairingDialogProps) {
  const { data: challenge, isLoading } = usePairingChallenge(token || undefined);
  const verifyMutation = useVerifyPairing();
  const [copied, setCopied] = useState(false);
  const [verifiedSuccess, setVerifiedSuccess] = useState(false);

  const handleAuthorize = async () => {
    if (!challenge) return;
    try {
      await verifyMutation.mutateAsync({
        token: challenge.token,
        code: challenge.code,
      });
      setVerifiedSuccess(true);
    } catch {
      // error handled in hook
    }
  };

  const handleCopySnippet = () => {
    if (!challenge) return;
    navigator.clipboard.writeText(`/codedock verify ${challenge.code}`);
    setCopied(true);
    toast.success('Copied verification command to clipboard');
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <Dialog open={Boolean(token)} onOpenChange={(open) => !open && onClose()}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <div className="flex items-center gap-2">
            <div className="flex size-8 items-center justify-center rounded-lg bg-primary/10 text-primary">
              <ShieldCheck className="size-5" />
            </div>
            <div>
              <DialogTitle>Pair Chat Account</DialogTitle>
              <DialogDescription>
                Authorize incoming chat command dispatch to your Codedock workspace.
              </DialogDescription>
            </div>
          </div>
        </DialogHeader>

        {isLoading ? (
          <div className="flex h-32 items-center justify-center text-muted-foreground text-xs">
            Loading pairing challenge...
          </div>
        ) : !challenge ? (
          <div className="py-6 text-center text-muted-foreground text-xs">
            Pairing challenge is invalid or has expired. Please run{' '}
            <code className="rounded bg-muted px-1.5 py-0.5 font-mono text-foreground">
              /codedock link
            </code>{' '}
            in chat to generate a new challenge.
          </div>
        ) : verifiedSuccess ? (
          <div className="space-y-4 py-4 text-center">
            <div className="mx-auto flex size-12 items-center justify-center rounded-full bg-emerald-500/10 text-emerald-500">
              <Check className="size-6" />
            </div>
            <div className="space-y-1">
              <p className="font-semibold text-foreground text-sm">Account Paired Successfully!</p>
              <p className="text-muted-foreground text-xs">
                Your {challenge.provider} user (
                <code className="font-mono text-foreground">{challenge.externalUserId}</code>) is
                now authorized to run operator commands.
              </p>
            </div>
            <DialogFooter>
              <Button type="button" size="sm" onClick={onClose} className="w-full">
                Done
              </Button>
            </DialogFooter>
          </div>
        ) : (
          <div className="space-y-4 py-2">
            <div className="space-y-2 rounded-lg border border-border/80 bg-muted/30 p-3 text-xs">
              <div className="flex items-center justify-between">
                <span className="text-muted-foreground">Provider:</span>
                <span className="font-semibold text-foreground uppercase tracking-wider">
                  {challenge.provider}
                </span>
              </div>
              <div className="flex items-center justify-between">
                <span className="text-muted-foreground">External User:</span>
                <span className="font-mono text-foreground">{challenge.externalUserId}</span>
              </div>
              <div className="flex items-center justify-between">
                <span className="text-muted-foreground">Default Role:</span>
                <span className="font-medium text-foreground">Operator (Scoped)</span>
              </div>
              <div className="flex items-center justify-between">
                <span className="text-muted-foreground">Expires:</span>
                <span className="text-muted-foreground">
                  {new Date(challenge.expiresAt).toLocaleTimeString()}
                </span>
              </div>
            </div>

            <div className="space-y-2 rounded-lg border border-primary/20 bg-primary/5 p-3">
              <p className="font-medium text-foreground text-xs">
                Enter this verification code in your chat channel:
              </p>
              <div className="flex items-center justify-between rounded border border-border bg-background px-3 py-2 font-mono text-sm">
                <span>/codedock verify {challenge.code}</span>
                <Button
                  type="button"
                  variant="ghost"
                  size="icon"
                  className="size-7"
                  onClick={handleCopySnippet}
                  aria-label="Copy verification command"
                >
                  {copied ? (
                    <Check className="size-3.5 text-emerald-500" />
                  ) : (
                    <Copy className="size-3.5" />
                  )}
                </Button>
              </div>
            </div>

            <DialogFooter className="gap-2 sm:gap-0">
              <Button type="button" variant="outline" size="sm" onClick={onClose}>
                Cancel
              </Button>
              <Button
                type="button"
                size="sm"
                onClick={handleAuthorize}
                disabled={verifyMutation.isPending}
              >
                {verifyMutation.isPending ? 'Authorizing...' : 'Authorize Directly'}
              </Button>
            </DialogFooter>
          </div>
        )}
      </DialogContent>
    </Dialog>
  );
}
