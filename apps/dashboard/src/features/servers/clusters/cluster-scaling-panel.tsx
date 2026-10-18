import {
  AlertCircle,
  CheckCircle2,
  ChevronDown,
  ChevronUp,
  Cpu,
  Loader2,
  Play,
  RotateCcw,
  Trash2,
} from 'lucide-react';
import { useEffect, useState } from 'react';
import { toast } from 'sonner';
import { Button } from '@codedock/ui/button';
import {
  useClusterScaling,
  useRemoveClusterScaling,
  useRetryClusterScaling,
  useSetupClusterScaling,
} from '#/hooks/use-clusters';

interface ClusterScalingPanelProps {
  clusterId: string;
  canManage?: boolean;
  onReadyChange?: (ready: boolean) => void;
}

export function ClusterScalingPanel({
  clusterId,
  canManage = true,
  onReadyChange,
}: ClusterScalingPanelProps) {
  const [detailsOpen, setDetailsOpen] = useState(false);

  const { data: scaling, isLoading, isRefetching } = useClusterScaling(clusterId);
  const setupMutation = useSetupClusterScaling();
  const retryMutation = useRetryClusterScaling();
  const removeMutation = useRemoveClusterScaling();

  useEffect(() => {
    if (scaling) {
      onReadyChange?.(scaling.status === 'ready');
    }
  }, [scaling, onReadyChange]);

  const enableScaling = () => {
    setupMutation.mutate(clusterId, {
      onSuccess: (data) => {
        toast.success('Cluster scaling enabled successfully');
        onReadyChange?.(data.status === 'ready');
      },
      onError: (err) => {
        toast.error(err instanceof Error ? err.message : 'Failed to enable scaling');
      },
    });
  };

  const retryScaling = () => {
    retryMutation.mutate(clusterId, {
      onSuccess: (data) => {
        toast.success('Retrying scaling setup');
        onReadyChange?.(data.status === 'ready');
      },
      onError: (err) => {
        toast.error(err instanceof Error ? err.message : 'Failed to retry scaling');
      },
    });
  };

  const removeScaling = () => {
    removeMutation.mutate(clusterId, {
      onSuccess: () => {
        toast.warning('Cluster scaling runtime removed');
        onReadyChange?.(false);
      },
      onError: (err) => {
        toast.error(err instanceof Error ? err.message : 'Failed to remove scaling runtime');
      },
    });
  };

  const isSettingUp = setupMutation.isPending;
  const isRetrying = retryMutation.isPending;
  const isRemoving = removeMutation.isPending;
  const status = scaling?.status ?? 'not_enabled';
  const isBusy = isSettingUp || isRetrying || isRemoving || isRefetching;

  return (
    <div className="rounded-2xl bg-card p-5 sm:p-6" data-testid="cluster-scaling-panel">
      <div className="flex flex-wrap items-center justify-between gap-4 border-border/50 border-b pb-4">
        <div className="flex items-center gap-3">
          <div
            className={`flex h-10 w-10 items-center justify-center rounded-xl ${
              status === 'ready'
                ? 'bg-emerald-500/10 text-emerald-500'
                : status === 'setting_up' || status === 'removing'
                  ? 'bg-primary/10 text-primary'
                  : status === 'failed' || status === 'interrupted'
                    ? 'bg-destructive/10 text-destructive'
                    : 'bg-muted text-muted-foreground'
            }`}
          >
            {isBusy ? (
              <Loader2 className="size-5 animate-spin" />
            ) : status === 'ready' ? (
              <CheckCircle2 className="size-5" />
            ) : status === 'failed' || status === 'interrupted' ? (
              <AlertCircle className="size-5" />
            ) : (
              <Cpu className="size-5" />
            )}
          </div>
          <div>
            <div className="flex items-center gap-2">
              <h3 className="font-semibold text-base text-foreground">Cluster Scaling Runtime</h3>
              <span
                className={`inline-flex items-center rounded-full px-2 py-0.5 font-medium text-xs ${
                  status === 'ready'
                    ? 'bg-emerald-500/10 text-emerald-500'
                    : status === 'failed' || status === 'interrupted'
                      ? 'bg-destructive/10 text-destructive'
                      : 'bg-muted text-muted-foreground'
                }`}
              >
                {status.replace('_', ' ')}
              </span>
            </div>
            <p className="text-muted-foreground text-xs">
              {status === 'ready'
                ? 'K3s lightweight cluster runtime verified and ready to run scaled workloads'
                : 'Automated K3s join and multi-node pod orchestration'}
            </p>
          </div>
        </div>

        {canManage && (
          <div className="flex items-center gap-2">
            {(status === 'not_enabled' || status === 'removing') && (
              <Button
                size="sm"
                disabled={isBusy || isLoading}
                onClick={() => enableScaling()}
                className="gap-1.5"
              >
                <Play className="size-3.5" />
                Enable scaling
              </Button>
            )}

            {(status === 'failed' || status === 'interrupted') && (
              <Button
                size="sm"
                variant="outline"
                disabled={isBusy}
                onClick={() => retryScaling()}
                className="gap-1.5"
              >
                <RotateCcw className="size-3.5" />
                Retry setup
              </Button>
            )}

            {status !== 'not_enabled' && status !== 'removing' && (
              <Button
                variant="ghost"
                size="sm"
                disabled={isBusy}
                onClick={() => removeScaling()}
                className="text-destructive hover:bg-destructive/10 hover:text-destructive"
              >
                <Trash2 className="size-3.5" />
              </Button>
            )}
          </div>
        )}
      </div>

      <div className="mt-5 space-y-4">
        <div className="grid grid-cols-2 gap-3 sm:grid-cols-4">
          <div className="rounded-xl bg-muted/30 p-3">
            <span className="text-muted-foreground text-xs">Control Nodes</span>
            <p className="font-semibold text-foreground text-sm">{scaling?.controlNodes ?? 1}</p>
          </div>
          <div className="rounded-xl bg-muted/30 p-3">
            <span className="text-muted-foreground text-xs">Worker Agents</span>
            <p className="font-semibold text-foreground text-sm">{scaling?.workerNodes ?? 0}</p>
          </div>
          <div className="rounded-xl bg-muted/30 p-3">
            <span className="text-muted-foreground text-xs">Orchestrator</span>
            <p className="font-semibold text-foreground text-sm">K3s v1.31.1</p>
          </div>
          <div className="rounded-xl bg-muted/30 p-3">
            <span className="text-muted-foreground text-xs">Network Mode</span>
            <p className="font-semibold text-foreground text-sm">WireGuard Flannel</p>
          </div>
        </div>

        <div>
          <button
            type="button"
            onClick={() => setDetailsOpen(!detailsOpen)}
            className="flex items-center gap-1.5 text-muted-foreground text-xs hover:text-foreground"
          >
            {detailsOpen ? (
              <ChevronUp className="size-3.5" />
            ) : (
              <ChevronDown className="size-3.5" />
            )}
            {detailsOpen ? 'Hide technical details' : 'Show technical details'}
          </button>

          {detailsOpen && (
            <div className="mt-3 space-y-2 rounded-xl border border-border/60 bg-muted/20 p-4 font-mono text-xs">
              <div className="flex justify-between">
                <span className="text-muted-foreground">K3s Release:</span>
                <span className="text-foreground">{scaling?.k3sVersion ?? 'v1.31.1+k3s1'}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-muted-foreground">Pod CIDR:</span>
                <span className="text-foreground">{scaling?.podCidr ?? '10.42.0.0/16'}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-muted-foreground">Service CIDR:</span>
                <span className="text-foreground">{scaling?.serviceCidr ?? '10.43.0.0/16'}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-muted-foreground">Cluster Domain:</span>
                <span className="text-foreground">{scaling?.clusterDomain ?? 'cluster.local'}</span>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
