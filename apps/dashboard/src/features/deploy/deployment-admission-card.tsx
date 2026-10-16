import { CheckCircle2, Clock, Cpu, RefreshCw } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Badge } from '#/components/ui/badge';
import { Button } from '#/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '#/components/ui/card';
import { Progress } from '#/components/ui/progress';
import { checksService, type BuildAdmissionStatus } from '#/services/checks';

interface DeploymentAdmissionCardProps {
  projectId: string;
  initialAdmission?: BuildAdmissionStatus | null;
}

export function DeploymentAdmissionCard({
  projectId,
  initialAdmission,
}: DeploymentAdmissionCardProps) {
  const [admission, setAdmission] = useState<BuildAdmissionStatus | null>(
    initialAdmission ?? {
      canAdmit: true,
      activeBuilds: 1,
      maxConcurrentBuilds: 4,
      queuePosition: 0,
      estimatedWaitSeconds: 0,
      reason: 'Capacity available for immediate execution',
    }
  );
  const [isRefreshing, setIsRefreshing] = useState(false);

  const handleRefresh = async () => {
    setIsRefreshing(true);
    try {
      const updated = await checksService.getAdmission(projectId);
      setAdmission(updated);
      toast.success('Admission quota refreshed');
    } catch {
      toast.error('Failed to refresh build admission status');
    } finally {
      setIsRefreshing(false);
    }
  };

  if (!admission) return null;

  const usagePercent = Math.min(
    100,
    Math.round((admission.activeBuilds / Math.max(1, admission.maxConcurrentBuilds)) * 100)
  );

  return (
    <Card className="border-border/70 shadow-sm">
      <CardHeader className="flex flex-row items-center justify-between pb-3">
        <div className="space-y-1">
          <div className="flex items-center gap-2">
            <Cpu className="size-4 text-primary" />
            <CardTitle className="text-base font-semibold">Build Admission & Concurrency</CardTitle>
          </div>
          <CardDescription className="text-xs">
            Guarantees fair execution slots and prevents container host memory exhaustion.
          </CardDescription>
        </div>
        <div className="flex items-center gap-2">
          <Badge variant={admission.canAdmit ? 'success' : 'secondary'} className="gap-1 text-xs">
            {admission.canAdmit ? (
              <CheckCircle2 className="size-3 text-emerald-500" />
            ) : (
              <Clock className="size-3 text-amber-500" />
            )}
            <span>{admission.canAdmit ? 'Admission Ready' : `Queued (#${admission.queuePosition})`}</span>
          </Badge>
          <Button
            type="button"
            variant="ghost"
            size="icon"
            onClick={handleRefresh}
            disabled={isRefreshing}
            className="size-7"
            aria-label="Refresh admission"
          >
            <RefreshCw className={`size-3.5 ${isRefreshing ? 'animate-spin' : ''}`} />
          </Button>
        </div>
      </CardHeader>

      <CardContent className="space-y-4 pt-1">
        <div className="space-y-1.5">
          <div className="flex justify-between text-xs text-muted-foreground">
            <span>Concurrent Slot Utilization</span>
            <span className="font-medium text-foreground">
              {admission.activeBuilds} of {admission.maxConcurrentBuilds} active
            </span>
          </div>
          <Progress value={usagePercent} className="h-2" />
        </div>

        <div className="flex flex-wrap items-center justify-between gap-2 rounded-lg border border-border/50 bg-muted/30 p-2.5 text-xs text-muted-foreground">
          <span>{admission.reason}</span>
          {admission.estimatedWaitSeconds > 0 && (
            <span className="font-mono text-foreground">~{admission.estimatedWaitSeconds}s wait</span>
          )}
        </div>
      </CardContent>
    </Card>
  );
}
