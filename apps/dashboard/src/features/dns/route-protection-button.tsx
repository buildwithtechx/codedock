import {
  AlertTriangle,
  CheckCircle2,
  KeyRound,
  Loader2,
  Lock,
  LogOut,
  RefreshCw,
  Shield,
  Unlock,
} from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import type { RouteProtectionAction, RouteProtectionDraft } from '@codedock/contracts';
import { Badge } from '@codedock/ui/badge';
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
import { useChangeRouteProtection, useRouteProtection } from '#/hooks/use-route-protection';

export type { RouteProtectionDraft };

interface RouteProtectionButtonProps {
  hostname: string;
  projectId?: string;
  unsupported?: boolean;
  draft?: RouteProtectionDraft;
  onDraftChange?: (draft: RouteProtectionDraft | undefined) => void;
}

export function RouteProtectionButton({
  hostname,
  projectId,
  unsupported,
  draft,
  onDraftChange,
}: RouteProtectionButtonProps) {
  const [open, setOpen] = useState(false);
  const [password, setPassword] = useState('');
  const [confirmUnlock, setConfirmUnlock] = useState(false);

  const { data: protectionList, isLoading } = useRouteProtection(projectId);
  const changeProtection = useChangeRouteProtection();

  const currentProtection = protectionList?.find((r) => r.hostname === hostname);
  const isEnabled = Boolean(draft || currentProtection?.enabled);
  const isApplied = Boolean(
    (!draft && currentProtection?.enabled && currentProtection.status === 'applied') ||
      (draft && !projectId)
  );
  const isPending = Boolean(
    !draft && currentProtection?.enabled && currentProtection.status === 'pending'
  );

  const handleAction = async (action: RouteProtectionAction) => {
    if (!projectId) {
      if (action === 'protect') {
        if (!password || password.length < 12) {
          toast.error('Password must be at least 12 characters long');
          return;
        }
        onDraftChange?.({ hostname, action: 'protect', password });
        toast.success('Route protection staged for deployment');
      } else if (action === 'unlock') {
        onDraftChange?.(undefined);
        toast.info('Route protection draft removed');
      }
      setOpen(false);
      setPassword('');
      setConfirmUnlock(false);
      return;
    }

    try {
      await changeProtection.mutateAsync({
        projectId,
        payload: {
          hostname,
          action,
          password: action === 'protect' ? password : undefined,
        },
      });

      if (action === 'protect') {
        toast.success('Route password protection updated');
      } else if (action === 'unlock') {
        toast.success('Route protection removed');
      } else if (action === 'revoke') {
        toast.success('Active browser sessions revoked across edge');
      } else if (action === 'retry') {
        toast.success('Edge protection synchronization retried');
      }

      setPassword('');
      setConfirmUnlock(false);
      if (action === 'unlock') {
        setOpen(false);
      }
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Failed to update route protection');
    }
  };

  return (
    <>
      <Button
        type="button"
        variant="ghost"
        size="icon"
        onClick={() => setOpen(true)}
        disabled={unsupported}
        title="Route Password Protection"
        className={`h-8 w-8 rounded-xl ${
          isApplied
            ? 'text-primary hover:bg-primary/10 hover:text-primary'
            : isPending
              ? 'text-amber-500 hover:bg-amber-500/10 hover:text-amber-500'
              : 'text-muted-foreground hover:bg-muted/50 hover:text-foreground'
        }`}
        aria-label={`Route Protection: ${hostname}`}
      >
        {isLoading ? (
          <Loader2 className="size-4 animate-spin" />
        ) : isApplied ? (
          <Lock className="size-4" />
        ) : (
          <Unlock className="size-4" />
        )}
      </Button>

      <Dialog open={open} onOpenChange={setOpen}>
        <DialogContent className="max-w-md rounded-2xl p-6">
          <DialogHeader icon={Shield}>
            <DialogTitle className="text-base font-semibold">Route Password Protection</DialogTitle>
            <DialogDescription className="break-all font-mono text-xs text-muted-foreground">
              {hostname}
            </DialogDescription>
          </DialogHeader>

          <div className="space-y-4 pt-2">
            <div className="flex items-center justify-between rounded-xl border border-border/60 bg-muted/20 p-3">
              <span className="text-xs font-medium text-foreground">Current Edge Status</span>
              {draft ? (
                <Badge
                  variant="outline"
                  className="gap-1 border-primary/40 bg-primary/10 text-primary"
                >
                  <CheckCircle2 className="size-3" /> Staged
                </Badge>
              ) : isApplied ? (
                <Badge
                  variant="outline"
                  className="gap-1 border-primary/40 bg-primary/10 text-primary"
                >
                  <Lock className="size-3" /> Protected
                </Badge>
              ) : isPending ? (
                <Badge
                  variant="outline"
                  className="gap-1 border-amber-500/40 bg-amber-500/10 text-amber-500"
                >
                  <AlertTriangle className="size-3" /> Syncing Edge
                </Badge>
              ) : (
                <Badge variant="outline" className="gap-1 border-border/70 text-muted-foreground">
                  <Unlock className="size-3" /> Public Route
                </Badge>
              )}
            </div>

            <p className="text-xs text-muted-foreground">
              Gate incoming requests behind a secure Traefik edge challenge with PBKDF2-HMAC-SHA256
              session encryption.
            </p>

            <form
              onSubmit={(e) => {
                e.preventDefault();
                void handleAction('protect');
              }}
              className="space-y-3"
            >
              <div className="space-y-1.5">
                <Label htmlFor={`pwd-${hostname}`} className="text-xs font-medium text-foreground">
                  {isEnabled ? 'Change Route Password' : 'Set Route Password'}
                </Label>
                <Input
                  id={`pwd-${hostname}`}
                  type="password"
                  placeholder="Minimum 12 characters"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  minLength={12}
                  maxLength={256}
                  disabled={changeProtection.isPending}
                  className="rounded-xl font-mono text-xs"
                />
                <p className="text-[11px] text-muted-foreground">
                  Must be at least 12 characters. Active sessions remain authenticated for 24 hours.
                </p>
              </div>

              <Button
                type="submit"
                disabled={changeProtection.isPending || password.length < 12}
                className="h-9 w-full gap-2 rounded-xl text-xs"
              >
                {changeProtection.isPending ? (
                  <Loader2 className="size-3.5 animate-spin" />
                ) : (
                  <KeyRound className="size-3.5" />
                )}
                {isEnabled ? 'Update Password' : 'Enable Route Protection'}
              </Button>
            </form>

            {isPending && projectId && (
              <Button
                type="button"
                variant="outline"
                size="sm"
                onClick={() => void handleAction('retry')}
                disabled={changeProtection.isPending}
                className="h-8 w-full gap-1.5 rounded-xl text-xs"
              >
                <RefreshCw className="size-3" /> Retry Edge Synchronization
              </Button>
            )}

            {isEnabled && (
              <div className="space-y-2 border-t border-border/40 pt-3">
                {projectId && (
                  <Button
                    type="button"
                    variant="outline"
                    size="sm"
                    onClick={() => void handleAction('revoke')}
                    disabled={changeProtection.isPending}
                    className="h-8 w-full gap-1.5 rounded-xl text-xs text-muted-foreground hover:text-foreground"
                  >
                    <LogOut className="size-3" /> Revoke All Active Browser Sessions
                  </Button>
                )}

                {!confirmUnlock ? (
                  <Button
                    type="button"
                    variant="ghost"
                    size="sm"
                    onClick={() => setConfirmUnlock(true)}
                    disabled={changeProtection.isPending}
                    className="h-8 w-full gap-1.5 rounded-xl text-xs text-destructive hover:bg-destructive/10 hover:text-destructive"
                  >
                    <Unlock className="size-3" /> Remove Protection
                  </Button>
                ) : (
                  <div className="space-y-2 rounded-xl border border-destructive/30 bg-destructive/5 p-3">
                    <p className="text-xs text-destructive">
                      Make this route public to all incoming internet traffic?
                    </p>
                    <div className="flex gap-2">
                      <Button
                        type="button"
                        variant="destructive"
                        size="sm"
                        onClick={() => void handleAction('unlock')}
                        disabled={changeProtection.isPending}
                        className="h-8 flex-1 rounded-xl text-xs"
                      >
                        Yes, Unlock Route
                      </Button>
                      <Button
                        type="button"
                        variant="ghost"
                        size="sm"
                        onClick={() => setConfirmUnlock(false)}
                        disabled={changeProtection.isPending}
                        className="h-8 flex-1 rounded-xl text-xs"
                      >
                        Cancel
                      </Button>
                    </div>
                  </div>
                )}
              </div>
            )}
          </div>
        </DialogContent>
      </Dialog>
    </>
  );
}
