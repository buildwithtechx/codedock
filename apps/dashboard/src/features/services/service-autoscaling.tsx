import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '#/components/ui/card';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import { Switch } from '#/components/ui/switch';
import { apiClient } from '#/lib/api-client';

type ScalingPolicy = {
  supported: boolean;
  serviceId: string;
  enabled: boolean;
  minReplicas: number;
  maxReplicas: number;
  scaleUpCpu: number;
  scaleDownCpu: number;
  cooldownSeconds: number;
  lastCpu: number;
  lastDecision: string;
  lastEvaluatedAt: string;
  lastScaledAt: string;
};

export function ServiceAutoscaling({ serviceId }: { serviceId: string }) {
  const queryClient = useQueryClient();
  const [draft, setDraft] = useState<ScalingPolicy | null>(null);
  const query = useQuery({
    queryKey: ['autoscaling', serviceId],
    queryFn: () => apiClient.get<{ data: ScalingPolicy }>(`/services/${serviceId}/autoscaling`),
    refetchInterval: 10000,
  });
  const save = useMutation({
    mutationFn: (policy: ScalingPolicy) =>
      apiClient.put<{ data: ScalingPolicy }>(`/services/${serviceId}/autoscaling`, policy),
    onSuccess: async () => {
      setDraft(null);
      await queryClient.invalidateQueries({ queryKey: ['autoscaling', serviceId] });
      toast.success('Autoscaling policy saved');
    },
  });
  const policy = draft ?? query.data?.data;
  const fields = [
    ['minReplicas', 'Minimum replicas', 1, 10],
    ['maxReplicas', 'Maximum replicas', 1, 10],
    ['scaleDownCpu', 'Scale down below CPU (%)', 0, 99],
    ['scaleUpCpu', 'Scale up above CPU (%)', 1, 100],
    ['cooldownSeconds', 'Cooldown (seconds)', 120, 86400],
  ] as const;
  return (
    <Card>
      <CardHeader>
        <CardTitle>Autoscaling</CardTitle>
        <CardDescription>
          Opt in to CPU-based replica changes on local Docker. The primary replica supplies the CPU
          sample. Decisions are checked every two minutes and create deployments.
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        {query.isLoading && <p>Loading policy...</p>}
        {query.isError && (
          <p role="alert">
            Could not load autoscaling policy.{' '}
            <Button variant="outline" onClick={() => query.refetch()}>
              Retry
            </Button>
          </p>
        )}
        {policy && (
          <>
            {!policy.supported && (
              <p role="status">
                Autoscaling is unavailable for this target. Choose a local Docker target to enable
                it.
              </p>
            )}
            <div className="flex items-center gap-3">
              <Switch
                id="autoscaling-enabled"
                checked={policy.enabled}
                disabled={save.isPending || (!policy.supported && !policy.enabled)}
                onCheckedChange={(enabled) => setDraft({ ...policy, enabled })}
              />
              <Label htmlFor="autoscaling-enabled">Enable autoscaling</Label>
            </div>
            <div className="grid gap-4 sm:grid-cols-2">
              {fields.map(([key, label, min, max]) => (
                <div key={key} className="space-y-2">
                  <Label htmlFor={`autoscaling-${key}`}>{label}</Label>
                  <Input
                    id={`autoscaling-${key}`}
                    type="number"
                    min={min}
                    max={max}
                    value={policy[key]}
                    disabled={save.isPending}
                    onChange={(event) => setDraft({ ...policy, [key]: Number(event.target.value) })}
                  />
                </div>
              ))}
            </div>
            {save.error && (
              <p role="alert" className="text-destructive text-sm">
                {save.error.message}
              </p>
            )}
            <div className="flex gap-2">
              <Button disabled={save.isPending || !draft} onClick={() => save.mutate(policy)}>
                {save.isPending ? 'Saving...' : 'Save policy'}
              </Button>
              <Button
                variant="outline"
                disabled={save.isPending || !draft}
                onClick={() => setDraft(null)}
              >
                Discard changes
              </Button>
            </div>
            <div className="rounded-lg bg-muted p-3 text-sm">
              <p>Last decision: {query.data?.data.lastDecision || 'Not evaluated yet'}</p>
              <p>
                Last CPU sample:{' '}
                {query.data?.data.lastEvaluatedAt
                  ? `${query.data.data.lastCpu.toFixed(2)}%`
                  : 'Unavailable'}
              </p>
              <p>
                Evaluated:{' '}
                {query.data?.data.lastEvaluatedAt
                  ? new Date(query.data.data.lastEvaluatedAt).toLocaleString()
                  : 'Never'}
              </p>
              <p>
                Last scale request:{' '}
                {query.data?.data.lastScaledAt
                  ? new Date(query.data.data.lastScaledAt).toLocaleString()
                  : 'Never'}
              </p>
            </div>
          </>
        )}
      </CardContent>
    </Card>
  );
}
