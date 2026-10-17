import { X } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import { Dialog, DialogContent } from '#/components/ui/dialog';
import { Input } from '#/components/ui/input';
import { serverService } from '#/services/servers';
import type { Server as ServerModel } from '#/interfaces/server';
import { AddServerCloudPlans, type CloudPlan } from './add-server-cloud-plans';

export interface AddServerModalProps {
  open: boolean;
  onClose: () => void;
  initialMode?: 'managed' | 'connected';
  onCreated?: (server: ServerModel) => void;
  onSelectPlan?: (plan: CloudPlan) => void;
}

export function AddServerModal({
  open,
  onClose,
  initialMode = 'managed',
  onCreated,
  onSelectPlan,
}: AddServerModalProps) {
  const [mode, setMode] = useState<'managed' | 'connected'>(initialMode);
  const [tab, setTab] = useState<'monthly' | 'payg'>('monthly');

  // Connected SSH form state
  const [serverName, setServerName] = useState('');
  const [ipAddress, setIpAddress] = useState('');
  const [sshPort, setSshPort] = useState(22);
  const [sshUser, setSshUser] = useState('root');
  const [connecting, setConnecting] = useState(false);

  const handleConnectServer = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!ipAddress.trim()) {
      toast.error('Host IP is required');
      return;
    }
    setConnecting(true);
    try {
      const server = await serverService.create({
        name: serverName.trim() || ipAddress.trim(),
        ipAddress: ipAddress.trim(),
        sshPort,
        sshUser: sshUser.trim() || 'root',
        sshAuthMethod: 'password',
      });
      toast.success('Server connected successfully');
      onCreated?.(server);
      onClose();
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Failed to connect server');
    } finally {
      setConnecting(false);
    }
  };

  return (
    <Dialog open={open} onOpenChange={(val) => !val && onClose()}>
      <DialogContent
        className={`overflow-hidden p-6 sm:p-7 ${
          mode === 'managed' ? 'max-w-5xl' : 'max-w-xl'
        }`}
        showCloseButton={false}
      >
        <div className="space-y-6">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-3">
              <h2 className="font-semibold text-foreground text-xl tracking-tight">Add Server</h2>
              <div className="flex items-center rounded-lg bg-muted/60 p-0.5 text-xs">
                <button
                  type="button"
                  onClick={() => setMode('managed')}
                  className={`rounded-md px-2.5 py-1 font-medium transition-colors ${
                    mode === 'managed'
                      ? 'bg-background text-foreground shadow-xs'
                      : 'text-muted-foreground hover:text-foreground'
                  }`}
                >
                  Cloud space
                </button>
                <button
                  type="button"
                  onClick={() => setMode('connected')}
                  className={`rounded-md px-2.5 py-1 font-medium transition-colors ${
                    mode === 'connected'
                      ? 'bg-background text-foreground shadow-xs'
                      : 'text-muted-foreground hover:text-foreground'
                  }`}
                >
                  Self-hosted
                </button>
              </div>
            </div>
            <button
              type="button"
              onClick={onClose}
              className="rounded-lg p-1.5 text-muted-foreground hover:bg-muted hover:text-foreground"
            >
              <X className="size-4" />
            </button>
          </div>

          {mode === 'managed' ? (
            <>
              <div className="flex items-center justify-between border-border/60 border-b pb-3">
                <div className="flex items-center gap-6">
                  <button
                    type="button"
                    onClick={() => setTab('monthly')}
                    className={`font-medium text-sm transition-colors ${
                      tab === 'monthly'
                        ? 'border-foreground border-b-2 pb-3 text-foreground'
                        : 'pb-3 text-muted-foreground hover:text-foreground'
                    }`}
                  >
                    Monthly plans
                  </button>
                  <button
                    type="button"
                    onClick={() => setTab('payg')}
                    className={`font-medium text-sm transition-colors ${
                      tab === 'payg'
                        ? 'border-foreground border-b-2 pb-3 text-foreground'
                        : 'pb-3 text-muted-foreground hover:text-foreground'
                    }`}
                  >
                    Prepaid PAYG
                  </button>
                </div>
                <div className="flex items-center rounded-lg bg-muted/60 p-0.5 text-xs">
                  <span className="rounded-md bg-background px-2.5 py-1 font-medium text-foreground shadow-xs">
                    Plans
                  </span>
                  <span className="px-2.5 py-1 text-muted-foreground">Custom</span>
                </div>
              </div>

              <AddServerCloudPlans onSelectPlan={onSelectPlan} onClose={onClose} />
            </>
          ) : (
            <form onSubmit={handleConnectServer} className="space-y-4">
              <div className="space-y-1.5">
                <label className="font-medium text-foreground text-sm">Server name</label>
                <Input
                  value={serverName}
                  onChange={(e) => setServerName(e.target.value)}
                  placeholder="e.g. production-vps"
                  className="rounded-xl bg-muted/40"
                />
              </div>
              <div className="grid gap-3 sm:grid-cols-3">
                <div className="space-y-1.5 sm:col-span-2">
                  <label className="font-medium text-foreground text-sm">Host IP / Domain</label>
                  <Input
                    value={ipAddress}
                    onChange={(e) => setIpAddress(e.target.value)}
                    placeholder="192.168.1.100"
                    required
                    className="rounded-xl bg-muted/40"
                  />
                </div>
                <div className="space-y-1.5">
                  <label className="font-medium text-foreground text-sm">SSH Port</label>
                  <Input
                    type="number"
                    value={sshPort}
                    onChange={(e) => setSshPort(Number(e.target.value))}
                    className="rounded-xl bg-muted/40"
                  />
                </div>
              </div>
              <div className="space-y-1.5">
                <label className="font-medium text-foreground text-sm">SSH User</label>
                <Input
                  value={sshUser}
                  onChange={(e) => setSshUser(e.target.value)}
                  placeholder="root"
                  className="rounded-xl bg-muted/40"
                />
              </div>
              <div className="pt-2">
                <Button type="submit" disabled={connecting} className="w-full rounded-xl">
                  {connecting ? 'Connecting...' : 'Connect Server'}
                </Button>
              </div>
            </form>
          )}
        </div>
      </DialogContent>
    </Dialog>
  );
}

export function useAddServerModal({
  initialMode = 'managed',
}: {
  initialMode?: 'managed' | 'connected';
} = {}) {
  const [open, setOpen] = useState(false);
  const [callback, setCallback] = useState<((server: ServerModel) => void) | undefined>();

  const openModal = (onCreated?: (server: ServerModel) => void) => {
    setCallback(() => onCreated);
    setOpen(true);
  };

  const modal = (
    <AddServerModal
      open={open}
      initialMode={initialMode}
      onClose={() => setOpen(false)}
      onCreated={(server) => {
        callback?.(server);
        setOpen(false);
      }}
    />
  );

  return { openModal, modal };
}
