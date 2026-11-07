import { Badge } from '@codedock/ui/badge';
import { Button } from '@codedock/ui/button';
import { Input } from '@codedock/ui/input';
import { Sheet, SheetContent, SheetDescription, SheetHeader, SheetTitle } from '@codedock/ui/sheet';
import { AlertCircle, Ban, Check, Loader2, RefreshCw, Shield, ShieldAlert, X } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { useBanIP, useFail2banStatus, useUnbanIP } from '#/hooks/use-mail';

interface PendingAction {
  type: 'ban' | 'unban';
  ip: string;
  jail: string;
}

export function Fail2banDrawer({
  serverId,
  serverName,
  isOpen,
  onClose,
}: {
  serverId?: string;
  serverName?: string;
  isOpen: boolean;
  onClose: () => void;
}) {
  const { data, isLoading, isError, error, refetch, isFetching } = useFail2banStatus(
    serverId,
    isOpen
  );
  const banMutation = useBanIP(serverId);
  const unbanMutation = useUnbanIP(serverId);

  const [banInputIp, setBanInputIp] = useState('');
  const [selectedJail, setSelectedJail] = useState('');
  const [pendingAction, setPendingAction] = useState<PendingAction | null>(null);

  const jails = data?.jails || [];
  const details = data?.details || {};
  const canManage = data?.canManage ?? true;

  const activeJail = selectedJail || (jails.length > 0 ? jails[0] : '');

  const handleReviewBan = (e: React.FormEvent) => {
    e.preventDefault();
    const trimmedIp = banInputIp.trim();
    if (!trimmedIp) {
      toast.error('Please enter an IP address');
      return;
    }
    if (!activeJail) {
      toast.error('Please select a jail');
      return;
    }
    setPendingAction({ type: 'ban', ip: trimmedIp, jail: activeJail });
  };

  const handleConfirmAction = async () => {
    if (!pendingAction) return;
    try {
      if (pendingAction.type === 'ban') {
        await banMutation.mutateAsync({ ip: pendingAction.ip, jail: pendingAction.jail });
        toast.success(`Banned ${pendingAction.ip} in ${pendingAction.jail}`);
        setBanInputIp('');
      } else {
        await unbanMutation.mutateAsync({ ip: pendingAction.ip, jail: pendingAction.jail });
        toast.success(`Unbanned ${pendingAction.ip} from ${pendingAction.jail}`);
      }
      setPendingAction(null);
      void refetch();
    } catch (err: any) {
      toast.error(err?.message || `Failed to ${pendingAction.type} IP`);
    }
  };

  return (
    <Sheet open={isOpen} onOpenChange={(open) => !open && onClose()}>
      <SheetContent side="right" className="flex w-full flex-col overflow-hidden p-0 sm:max-w-xl">
        <SheetHeader className="border-border/50 border-b p-6 pb-4">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-2.5">
              <Shield className="size-5 text-primary" />
              <SheetTitle className="text-lg">Fail2ban Brute-force Protection</SheetTitle>
            </div>
            <Button
              size="sm"
              variant="outline"
              onClick={() => void refetch()}
              disabled={isLoading || isFetching}
              className="h-8 gap-1.5 text-xs"
            >
              <RefreshCw className={`size-3.5 ${isFetching ? 'animate-spin' : ''}`} />
              Refresh
            </Button>
          </div>
          <SheetDescription className="mt-1 text-muted-foreground text-xs">
            Active intrusion prevention jails on {serverName || 'mail server'}. Inspect live bans
            and release false positives safely.
          </SheetDescription>
        </SheetHeader>

        <div className="flex-1 space-y-6 overflow-y-auto p-6">
          {isLoading ? (
            <div className="flex items-center justify-center py-12 text-muted-foreground text-sm">
              <Loader2 className="mr-2 size-4 animate-spin" />
              Loading Fail2ban jails...
            </div>
          ) : isError ? (
            <div className="flex items-start gap-3 rounded-xl border border-destructive/30 bg-destructive/10 p-4 text-destructive text-sm">
              <AlertCircle className="mt-0.5 size-5 shrink-0" />
              <div>
                <p className="font-medium">Failed to retrieve Fail2ban status</p>
                <p className="mt-1 text-destructive/80 text-xs">
                  {error instanceof Error ? error.message : 'Unknown error'}
                </p>
              </div>
            </div>
          ) : jails.length === 0 ? (
            <div className="rounded-xl border border-border/60 border-dashed p-8 text-center text-muted-foreground text-sm">
              <ShieldAlert className="mx-auto mb-2 size-8 opacity-50" />
              <p className="font-medium text-foreground">No Fail2ban jails detected</p>
              <p className="mt-1 text-xs">
                Ensure Fail2ban is active on this mail server container.
              </p>
            </div>
          ) : (
            <div className="space-y-6">
              {/* Manual ban form */}
              {canManage && (
                <form
                  onSubmit={handleReviewBan}
                  className="space-y-3 rounded-xl border border-border/60 bg-muted/30 p-4"
                >
                  <p className="font-semibold text-foreground text-xs uppercase tracking-wider">
                    Ban Suspicious IP
                  </p>
                  <div className="flex flex-col gap-2 sm:flex-row">
                    <select
                      value={activeJail}
                      onChange={(e) => setSelectedJail(e.target.value)}
                      className="rounded-lg border border-border/60 bg-background px-3 py-1.5 text-foreground text-xs focus:outline-none focus:ring-1 focus:ring-primary"
                    >
                      {jails.map((j) => (
                        <option key={j} value={j}>
                          {j}
                        </option>
                      ))}
                    </select>
                    <Input
                      dir="ltr"
                      placeholder="e.g. 198.51.100.42 or 2001:db8::1"
                      value={banInputIp}
                      onChange={(e) => setBanInputIp(e.target.value)}
                      className="h-8 flex-1 font-mono text-xs"
                    />
                    <Button
                      type="submit"
                      size="sm"
                      variant="secondary"
                      className="h-8 shrink-0 gap-1 text-xs"
                    >
                      <Ban className="size-3.5" />
                      Review ban
                    </Button>
                  </div>
                </form>
              )}

              {/* Jails and banned IPs */}
              <div className="space-y-4">
                {jails.map((jailName) => {
                  const jail = details[jailName];
                  const bannedIps = jail?.bannedIps || [];

                  return (
                    <div
                      key={jailName}
                      className="space-y-3 rounded-xl border border-border/60 bg-card p-4 shadow-xs"
                    >
                      <div className="flex items-center justify-between border-border/40 border-b pb-2.5">
                        <div className="flex items-center gap-2">
                          <span className="font-semibold text-foreground text-sm">{jailName}</span>
                          <Badge variant="outline" className="text-[10px]">
                            {bannedIps.length} banned
                          </Badge>
                        </div>
                        <div className="flex items-center gap-3 text-muted-foreground text-xs">
                          <span>
                            Failed now:{' '}
                            <strong className="text-foreground">
                              {jail?.currentlyFailed ?? 0}
                            </strong>
                          </span>
                          <span>
                            Total bans:{' '}
                            <strong className="text-foreground">{jail?.totalBanned ?? 0}</strong>
                          </span>
                        </div>
                      </div>

                      {bannedIps.length === 0 ? (
                        <p className="py-1 text-muted-foreground text-xs">
                          No IP addresses are currently banned in this jail.
                        </p>
                      ) : (
                        <ul className="divide-y divide-border/30">
                          {bannedIps.map((ip) => (
                            <li key={ip} className="flex items-center justify-between py-2 text-xs">
                              <bdi dir="ltr" className="font-medium font-mono text-foreground">
                                {ip}
                              </bdi>
                              {canManage && (
                                <Button
                                  size="sm"
                                  variant="ghost"
                                  onClick={() =>
                                    setPendingAction({
                                      type: 'unban',
                                      ip,
                                      jail: jailName,
                                    })
                                  }
                                  className="h-6 px-2 text-muted-foreground text-xs hover:text-foreground"
                                >
                                  Unban
                                </Button>
                              )}
                            </li>
                          ))}
                        </ul>
                      )}
                    </div>
                  );
                })}
              </div>
            </div>
          )}
        </div>

        {/* Sticky two-step confirmation footer */}
        {pendingAction && (
          <div className="flex items-center justify-between gap-3 border-border/60 border-t bg-popover p-4 shadow-xl">
            <div className="min-w-0 flex-1">
              <p className="truncate font-medium text-foreground text-xs">
                {pendingAction.type === 'ban' ? (
                  <>
                    Ban{' '}
                    <bdi dir="ltr" className="font-mono font-semibold">
                      {pendingAction.ip}
                    </bdi>{' '}
                    in <strong className="font-semibold">{pendingAction.jail}</strong>?
                  </>
                ) : (
                  <>
                    Unban{' '}
                    <bdi dir="ltr" className="font-mono font-semibold">
                      {pendingAction.ip}
                    </bdi>{' '}
                    from <strong className="font-semibold">{pendingAction.jail}</strong>?
                  </>
                )}
              </p>
              <p className="text-[11px] text-muted-foreground">
                Action executes immediately on the live server.
              </p>
            </div>
            <div className="flex shrink-0 items-center gap-2">
              <Button
                size="sm"
                variant="outline"
                onClick={() => setPendingAction(null)}
                className="h-8 text-xs"
              >
                <X className="mr-1 size-3.5" />
                Cancel
              </Button>
              <Button
                size="sm"
                variant={pendingAction.type === 'ban' ? 'destructive' : 'default'}
                disabled={banMutation.isPending || unbanMutation.isPending}
                onClick={handleConfirmAction}
                className="h-8 text-xs"
              >
                {banMutation.isPending || unbanMutation.isPending ? (
                  <Loader2 className="mr-1 size-3.5 animate-spin" />
                ) : (
                  <Check className="mr-1 size-3.5" />
                )}
                Confirm
              </Button>
            </div>
          </div>
        )}
      </SheetContent>
    </Sheet>
  );
}
