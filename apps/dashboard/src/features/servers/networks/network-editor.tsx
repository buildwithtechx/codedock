import { Check, Loader2, Network as NetworkIcon, Plus, Shield, Trash2 } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import {
  networksService,
  type CreateNetworkInput,
  type NetworkMember,
  type ServerNetwork,
} from '#/services/networks';
import { NetworkPeerDiagnostics } from './network-peer-diagnostics';
import { NetworkRollbackWatchdog } from './network-rollback-watchdog';

interface NetworkEditorProps {
  initialNetwork?: ServerNetwork | null;
  onSaved?: (network: ServerNetwork) => void;
}

export function NetworkEditor({ initialNetwork, onSaved }: NetworkEditorProps) {
  const [network, setNetwork] = useState<ServerNetwork | null>(initialNetwork ?? null);
  const [name, setName] = useState(initialNetwork?.name ?? 'alpha-mesh');
  const [mode, setMode] = useState<'wireguard' | 'native'>(initialNetwork?.mode ?? 'wireguard');
  const [subnetCidr, setSubnetCidr] = useState(initialNetwork?.subnetCidr ?? '10.42.0.0/16');
  const [mtu, setMtu] = useState(initialNetwork?.mtu ?? 1420);
  const [verificationPort, setVerificationPort] = useState(initialNetwork?.verificationPort ?? 51821);
  const [members, setMembers] = useState<NetworkMember[]>(
    initialNetwork?.members ?? [
      {
        serverId: 'srv-node-1',
        serverName: 'Control Node 1',
        privateIp: '10.42.0.1',
        interfaceName: 'wg0',
        status: 'ready',
        latencyMs: 1.4,
      },
      {
        serverId: 'srv-node-2',
        serverName: 'Worker Node 2',
        privateIp: '10.42.0.2',
        interfaceName: 'wg0',
        status: 'ready',
        latencyMs: 2.1,
      },
    ]
  );
  const [isSubmitting, setIsSubmitting] = useState(false);

  const addMember = () => {
    const nextIndex = members.length + 1;
    setMembers([
      ...members,
      {
        serverId: `srv-node-${nextIndex}`,
        serverName: `Worker Node ${nextIndex}`,
        privateIp: `10.42.0.${nextIndex}`,
        interfaceName: 'wg0',
        status: 'ready',
        latencyMs: 1.8,
      },
    ]);
  };

  const removeMember = (index: number) => {
    if (members.length <= 2) {
      toast.error('A private network mesh requires at least 2 servers');
      return;
    }
    setMembers(members.filter((_, i) => i !== index));
  };

  const handleApply = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim()) {
      toast.error('Network name is required');
      return;
    }
    if (members.length < 2) {
      toast.error('At least 2 member servers are required');
      return;
    }

    setIsSubmitting(true);
    try {
      const payload: CreateNetworkInput = {
        name: name.trim(),
        mode,
        subnetCidr: subnetCidr.trim(),
        mtu: Number(mtu) || 1420,
        verificationPort: Number(verificationPort) || 51821,
        members,
        accessPolicy: { mode: 'full_mesh' },
      };

      const created = await networksService.create(payload);
      setNetwork(created);
      toast.success('Private network created and 120s rollback watchdog armed');
      onSaved?.(created);
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Failed to apply network');
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="space-y-6" data-testid="network-editor">
      {network?.watchdog && (
        <NetworkRollbackWatchdog
          watchdog={network.watchdog}
          onCommitted={() => {
            setNetwork((prev) =>
              prev
                ? {
                    ...prev,
                    status: 'verified',
                    watchdog: prev.watchdog ? { ...prev.watchdog, status: 'committed', armed: false } : undefined,
                  }
                : null
            );
          }}
          onRolledBack={() => {
            setNetwork((prev) =>
              prev
                ? {
                    ...prev,
                    status: 'rolled_back',
                    watchdog: prev.watchdog ? { ...prev.watchdog, status: 'rolled_back', armed: false } : undefined,
                  }
                : null
            );
          }}
        />
      )}

      <form onSubmit={handleApply} className="rounded-2xl bg-card p-5 space-y-6">
        <div className="flex items-center justify-between border-border/50 border-b pb-4">
          <div className="flex items-center gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-primary/10">
              <NetworkIcon className="size-5 text-primary" />
            </div>
            <div>
              <h3 className="font-semibold text-base text-foreground">Private Network Mesh</h3>
              <p className="text-muted-foreground text-xs">
                Encrypted WireGuard peer-to-peer overlay with rollback safety
              </p>
            </div>
          </div>

          <div className="flex items-center gap-2">
            <Button
              type="button"
              variant={mode === 'wireguard' ? 'default' : 'outline'}
              size="sm"
              onClick={() => setMode('wireguard')}
            >
              Managed WireGuard
            </Button>
            <Button
              type="button"
              variant={mode === 'native' ? 'default' : 'outline'}
              size="sm"
              onClick={() => setMode('native')}
            >
              Native Adoption
            </Button>
          </div>
        </div>

        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
          <div className="space-y-2">
            <Label htmlFor="net-name">Network Name</Label>
            <Input
              id="net-name"
              value={name}
              onChange={(e) => setName(e.target.value)}
              disabled={isSubmitting}
              required
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="net-subnet">Subnet CIDR</Label>
            <Input
              id="net-subnet"
              value={subnetCidr}
              onChange={(e) => setSubnetCidr(e.target.value)}
              disabled={isSubmitting}
              required
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="net-mtu">Tunnel MTU</Label>
            <Input
              id="net-mtu"
              type="number"
              value={mtu}
              onChange={(e) => setMtu(Number(e.target.value))}
              disabled={isSubmitting}
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="net-port">Verification UDP Port</Label>
            <Input
              id="net-port"
              type="number"
              value={verificationPort}
              onChange={(e) => setVerificationPort(Number(e.target.value))}
              disabled={isSubmitting}
            />
          </div>
        </div>

        <div className="space-y-3 pt-2">
          <div className="flex items-center justify-between">
            <Label className="font-semibold text-foreground text-sm">
              Member Nodes ({members.length})
            </Label>
            <Button type="button" variant="outline" size="sm" onClick={addMember} className="gap-1.5">
              <Plus className="size-3.5" />
              Add Server
            </Button>
          </div>

          <div className="divide-y divide-border/50 rounded-xl border border-border/60 bg-muted/20">
            {members.map((member, idx) => (
              <div key={member.serverId} className="flex flex-wrap items-center justify-between gap-3 p-3">
                <div className="flex items-center gap-3">
                  <span className="flex h-6 w-6 items-center justify-center rounded-md bg-muted font-mono text-muted-foreground text-xs">
                    {idx + 1}
                  </span>
                  <div>
                    <p className="font-medium text-foreground text-xs">{member.serverName || member.serverId}</p>
                    <p className="font-mono text-muted-foreground text-[11px]">{member.privateIp}</p>
                  </div>
                </div>

                <div className="flex items-center gap-2">
                  <Input
                    className="h-8 w-32 font-mono text-xs"
                    value={member.privateIp}
                    onChange={(e) => {
                      const updated = [...members];
                      updated[idx].privateIp = e.target.value;
                      setMembers(updated);
                    }}
                    placeholder="Private IP"
                  />
                  <Button
                    type="button"
                    variant="ghost"
                    size="icon"
                    className="h-8 w-8 text-destructive hover:bg-destructive/10"
                    onClick={() => removeMember(idx)}
                    disabled={members.length <= 2}
                  >
                    <Trash2 className="size-3.5" />
                  </Button>
                </div>
              </div>
            ))}
          </div>
        </div>

        {network && network.members.length > 0 && (
          <NetworkPeerDiagnostics members={network.members} />
        )}

        <div className="flex items-center justify-between border-border/50 border-t pt-4">
          <div className="flex items-center gap-2 text-muted-foreground text-xs">
            <Shield className="size-4 text-emerald-500" />
            <span>120s rollback watchdog will automatically protect against host lock-out</span>
          </div>

          <Button type="submit" disabled={isSubmitting} className="gap-2">
            {isSubmitting && <Loader2 className="size-4 animate-spin" />}
            Apply & Arm Watchdog
          </Button>
        </div>
      </form>
    </div>
  );
}
