import { Link } from '@tanstack/react-router';
import { RefreshCw, Server } from 'lucide-react';
import type { Server as ServerModel } from '#/interfaces/server';
import type { CloudPlan } from './runner-cloud-plan-modal';
import type { DestinationServer } from './runner-editor';

interface RunnerServerFieldProps {
  isCloud: boolean;
  selectedPlan: CloudPlan | null;
  initialServer?: DestinationServer;
  serverId: string;
  onServerIdChange: (id: string, name?: string) => void;
  servers: ServerModel[];
  onOpenPlanModal: () => void;
  probing: boolean;
  onRecheckCapabilities: () => void;
}

export function RunnerServerField({
  isCloud,
  selectedPlan,
  initialServer,
  serverId,
  onServerIdChange,
  servers,
  onOpenPlanModal,
  probing,
  onRecheckCapabilities,
}: RunnerServerFieldProps) {
  return (
    <div className="space-y-1.5">
      <label className="font-medium text-foreground text-sm">
        {isCloud ? 'Cloud space' : 'Server'}
      </label>
      {isCloud ? (
        <div className="flex items-center justify-between rounded-xl border border-border/60 bg-muted/40 px-3.5 py-2.5 text-sm">
          <span className="font-medium text-foreground">
            {selectedPlan ? `${selectedPlan.name} (${selectedPlan.price}/mo)` : 'Cloud space (Starter)'}
          </span>
          <span className="text-muted-foreground text-xs">Managed</span>
        </div>
      ) : initialServer ? (
        <div className="flex items-center gap-2.5 rounded-xl border border-border/60 bg-muted/40 px-3.5 py-2.5 text-sm">
          <Server className="size-4 text-muted-foreground" />
          <span className="font-medium text-foreground">{initialServer.name}</span>
        </div>
      ) : (
        <select
          value={serverId}
          onChange={(e) => {
            const nextId = e.target.value;
            const s = servers.find((srv) => srv.id === nextId);
            onServerIdChange(nextId, s?.name);
          }}
          required
          className="h-10 w-full rounded-xl border border-border/60 bg-muted/40 px-3 text-foreground text-sm focus:outline-none focus:ring-2 focus:ring-ring"
        >
          {servers.map((s) => (
            <option key={s.id} value={s.id}>
              {s.name} ({s.ipAddress || 'Connected'})
            </option>
          ))}
        </select>
      )}
      <div className="flex items-center justify-between text-muted-foreground text-xs">
        <span>Uses the existing plan.</span>
        {isCloud ? (
          <button
            type="button"
            onClick={onOpenPlanModal}
            className="font-medium text-primary hover:underline"
          >
            + New Cloud space
          </button>
        ) : (
          <Link to={'/servers' as never} className="text-primary hover:underline">
            + Connect server
          </Link>
        )}
      </div>

      {!isCloud && serverId && (
        <div className="mt-2 flex items-center justify-between rounded-xl bg-muted/30 px-3 py-2 text-muted-foreground text-xs">
          <span>
            {probing ? 'Checking destination capabilities...' : 'Linux · x64 · Docker available'}
          </span>
          <button
            type="button"
            onClick={onRecheckCapabilities}
            disabled={probing}
            className="hover:text-foreground"
          >
            <RefreshCw className={`size-3.5 ${probing ? 'animate-spin' : ''}`} />
          </button>
        </div>
      )}
    </div>
  );
}
