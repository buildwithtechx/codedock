import { Activity, ArrowRightLeft, CheckCircle2, Wifi } from 'lucide-react';
import type { NetworkMember } from '#/services/networks';

interface NetworkPeerDiagnosticsProps {
  members: NetworkMember[];
}

export function NetworkPeerDiagnostics({ members }: NetworkPeerDiagnosticsProps) {
  return (
    <div className="space-y-3" data-testid="network-peer-diagnostics">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2">
          <Activity className="size-4 text-primary" />
          <h4 className="font-semibold text-foreground text-sm">Peer Reachability & Latency</h4>
        </div>
        <span className="text-muted-foreground text-xs">{members.length} peers connected</span>
      </div>

      <div className="divide-y divide-border/50 rounded-xl border border-border/60 bg-card">
        {members.map((member) => (
          <div
            key={member.serverId}
            className="flex flex-wrap items-center justify-between gap-3 p-3 text-xs"
          >
            <div className="flex items-center gap-3">
              <div className="flex h-7 w-7 items-center justify-center rounded-lg bg-muted text-muted-foreground">
                <Wifi className="size-3.5" />
              </div>
              <div>
                <p className="font-semibold text-foreground">
                  {member.serverName || member.serverId}
                </p>
                <p className="font-mono text-[11px] text-muted-foreground">
                  {member.privateIp} ({member.interfaceName})
                </p>
              </div>
            </div>

            <div className="flex items-center gap-4">
              <div className="flex items-center gap-1.5 font-mono text-muted-foreground">
                <ArrowRightLeft className="size-3 text-muted-foreground/70" />
                <span>
                  {typeof member.latencyMs === 'number' ? `${member.latencyMs.toFixed(1)} ms` : '—'}
                </span>
              </div>

              <div className="flex items-center gap-1.5 font-mono text-muted-foreground">
                <span>
                  Loss: {typeof member.packetLoss === 'number' ? `${member.packetLoss}%` : '0%'}
                </span>
              </div>

              <span
                className={`inline-flex items-center gap-1 rounded-full px-2 py-0.5 font-medium text-[11px] ${
                  member.status === 'ready' || member.status === 'verified'
                    ? 'bg-emerald-500/10 text-emerald-500'
                    : 'bg-amber-500/10 text-amber-500'
                }`}
              >
                <CheckCircle2 className="size-2.5" />
                {member.status}
              </span>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
