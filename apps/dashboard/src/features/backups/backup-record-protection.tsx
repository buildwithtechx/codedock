import { useMutation, useQueryClient } from '@tanstack/react-query';
import { Button } from '#/components/ui/button';
import { apiClient } from '#/lib/api-client';

export function BackupRecordProtection({
  recordId,
  until = 0,
}: {
  recordId: string;
  until?: number;
}) {
  const client = useQueryClient();
  const protectedRecord = until * 1000 > Date.now();
  const update = useMutation({
    mutationFn: () =>
      apiClient.put(`/backup-records/${recordId}/protection`, {
        protectedUntil: protectedRecord ? 0 : Math.floor(Date.now() / 1000) + 30 * 86400,
      }),
    onSuccess: () =>
      client.invalidateQueries({
        predicate: (query) =>
          query.queryKey.some((key) => typeof key === 'string' && key.includes('backup')),
      }),
  });
  return (
    <div className="space-y-1">
      <Button
        variant="outline"
        size="sm"
        disabled={update.isPending}
        onClick={() => update.mutate()}
      >
        {protectedRecord ? 'Remove protection' : 'Protect 30 days'}
      </Button>
      {protectedRecord && (
        <p className="text-muted-foreground text-xs">
          Protected until {new Date(until * 1000).toLocaleDateString()}
        </p>
      )}
      {update.error && (
        <p role="alert" className="text-destructive text-xs">
          {update.error.message}
        </p>
      )}
    </div>
  );
}
