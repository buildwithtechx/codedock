import { CheckCircle2, Database, Loader2, ShieldCheck, TrendingUp } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Badge } from '@codedock/ui/badge';
import { Button } from '@codedock/ui/button';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@codedock/ui/dialog';
import type { ClusterOperation } from '#/interfaces/clusters';
import { useApplyDatabaseScale, useReviewDatabaseScale } from '#/hooks/use-clusters';

interface ClusterDatabaseScalingProps {
  projectId: string;
  clusterId?: string;
  databaseId: string;
  currentInstances?: number;
  engine?: string;
  onOpenSetup?: () => void;
}

export function ClusterDatabaseScaling({
  projectId,
  clusterId,
  databaseId,
  currentInstances = 1,
  engine = 'postgres',
  onOpenSetup,
}: ClusterDatabaseScalingProps) {
  const [targetInstances, setTargetInstances] = useState(currentInstances);
  const [reviewOpen, setReviewOpen] = useState(false);
  const [review, setReview] = useState<ClusterOperation | null>(null);

  const reviewScale = useReviewDatabaseScale();
  const applyScale = useApplyDatabaseScale();

  const isReviewing = reviewScale.isPending;
  const isApplying = applyScale.isPending;
  const isScaled = currentInstances > 1;

  const handleStartReview = async (instances: number) => {
    if (!clusterId) {
      onOpenSetup?.();
      return;
    }

    setTargetInstances(instances);
    try {
      const data = await reviewScale.mutateAsync({
        projectId,
        clusterId,
        spec: {
          id: databaseId,
          instances,
        },
      });
      setReview(data);
      setReviewOpen(true);
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Failed to prepare scaling review');
    }
  };

  const handleApply = async () => {
    if (!review || !clusterId) return;
    try {
      await applyScale.mutateAsync({
        projectId,
        clusterId,
        operationId: review.id,
        confirmation: 'CONFIRM',
      });
      toast.success(`Scaling requested: targeting ${targetInstances} instances`);
      setReviewOpen(false);
      setReview(null);
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Failed to apply scaling mutation');
    }
  };

  return (
    <div className="space-y-3 rounded-xl border border-border/60 bg-muted/20 p-3.5">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2">
          <Database className="size-4 text-primary" />
          <span className="font-semibold text-foreground text-xs">Database Topology Scaling</span>
        </div>
        <Badge
          variant="outline"
          className={`text-[10px] ${
            isScaled ? 'border-primary/40 bg-primary/10 text-primary' : 'text-muted-foreground'
          }`}
        >
          {currentInstances} {currentInstances === 1 ? 'instance' : 'replicas'}
        </Badge>
      </div>

      <p className="text-[11px] text-muted-foreground">
        {engine === 'postgres'
          ? 'CloudNativePG handles automated streaming replication, primary failover, and read-replica DNS routing on K3s.'
          : 'Redis replication with durable append-only persistence.'}
      </p>

      <div className="flex items-center gap-2 pt-1">
        {[1, 2, 3].map((count) => {
          const isCurrent = currentInstances === count;
          return (
            <Button
              key={count}
              type="button"
              variant={isCurrent ? 'secondary' : 'outline'}
              size="sm"
              disabled={isCurrent || isReviewing}
              onClick={() => handleStartReview(count)}
              className={`h-7 flex-1 rounded-lg text-xs ${
                isCurrent ? 'border-primary/40 font-semibold text-primary' : ''
              }`}
            >
              {count === 1 ? '1 (Solo)' : count === 2 ? '2 (Pair)' : '3 (HA)'}
            </Button>
          );
        })}
      </div>

      <Dialog open={reviewOpen} onOpenChange={setReviewOpen}>
        <DialogContent className="max-w-md rounded-2xl p-6">
          <DialogHeader icon={TrendingUp}>
            <DialogTitle className="text-base font-semibold">Review Database Scaling</DialogTitle>
            <DialogDescription className="text-xs text-muted-foreground">
              Verify topology replica changes and environment bindings before applying.
            </DialogDescription>
          </DialogHeader>

          <div className="space-y-4 pt-2 text-xs">
            <div className="flex items-center justify-between rounded-xl border border-border/60 bg-muted/20 p-3">
              <span className="text-muted-foreground">Scale Transition</span>
              <span className="font-semibold font-mono text-foreground">
                {currentInstances} → {targetInstances} instances
              </span>
            </div>

            {targetInstances > 1 && engine === 'postgres' && (
              <div className="rounded-xl border border-primary/20 bg-primary/5 p-3 text-[11px] text-primary">
                <div className="flex items-center gap-1.5 font-semibold">
                  <ShieldCheck className="size-3.5" />
                  Read Replica Routing Enabled
                </div>
                <p className="mt-1 text-muted-foreground">
                  Applications bound to this cluster database will automatically receive
                  <span className="font-mono text-foreground"> DATABASE_READ_URL</span> pointing to
                  the replica service.
                </p>
              </div>
            )}

            {review?.effects && (
              <div className="space-y-1">
                <span className="font-medium text-muted-foreground">Planned Effects:</span>
                <p className="rounded-lg bg-background p-2.5 font-mono text-[11px] leading-relaxed text-foreground">
                  {review.effects}
                </p>
              </div>
            )}
          </div>

          <DialogFooter className="gap-2 sm:justify-end">
            <Button
              type="button"
              variant="ghost"
              size="sm"
              onClick={() => setReviewOpen(false)}
              disabled={isApplying}
              className="rounded-xl text-xs"
            >
              Cancel
            </Button>
            <Button
              type="button"
              size="sm"
              onClick={handleApply}
              disabled={isApplying}
              className="gap-1.5 rounded-xl text-xs"
            >
              {isApplying ? (
                <Loader2 className="size-3.5 animate-spin" />
              ) : (
                <CheckCircle2 className="size-3.5" />
              )}
              Confirm & Scale
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
