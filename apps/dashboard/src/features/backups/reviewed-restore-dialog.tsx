import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { useState } from 'react';
import { Button } from '#/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '#/components/ui/dialog';
import { Label } from '#/components/ui/label';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '#/components/ui/select';
import { databasesService } from '#/features/databases/api';
import type { BaseResponse } from '#/interfaces/base';
import { apiClient } from '#/lib/api-client';

type Operation = {
  id: string;
  status: string;
  phase: string;
  effects: string;
  error: string;
  logs: string;
  expiresAt: number;
};
type Review = { operation: Operation; confirmation: string };

export function ReviewedRestoreDialog({
  recordId,
  onClose,
  allowDatabaseTarget = false,
  sourceDatabaseId,
}: {
  recordId: string;
  onClose: () => void;
  allowDatabaseTarget?: boolean;
  sourceDatabaseId?: string;
}) {
  const client = useQueryClient();
  const [target, setTarget] = useState('');
  const [review, setReview] = useState<Review>();
  const [approved, setApproved] = useState(false);
  const source = useQuery({
    queryKey: ['restore-source', sourceDatabaseId],
    queryFn: () => databasesService.getDatabase(sourceDatabaseId!),
    enabled: allowDatabaseTarget && !!sourceDatabaseId,
  });
  const destinations = useQuery({
    queryKey: ['restore-destinations', source.data?.data.projectId],
    queryFn: () => databasesService.getDatabases(source.data!.data.projectId),
    enabled: !!source.data?.data.projectId,
  });
  const candidates = (destinations.data?.data ?? []).filter(
    (database) =>
      database.engine === source.data?.data.engine &&
      database.status === 'running' &&
      database.id !== sourceDatabaseId
  );

  const prepare = useMutation({
    mutationFn: () =>
      apiClient.post<BaseResponse<Review>>(`/backup-records/${recordId}/restore/review`, {
        targetDatabaseId: target,
      }),
    onSuccess: (response) => {
      setReview(response.data);
      setApproved(false);
    },
  });
  const operation = useQuery({
    queryKey: ['operations', review?.operation.id],
    queryFn: () => apiClient.get<BaseResponse<Operation>>(`/operations/${review?.operation.id}`),
    enabled: !!review,
    refetchInterval: 1500,
  });
  const apply = useMutation({
    mutationFn: () =>
      apiClient.post(`/backup-operations/${review?.operation.id}/apply`, {
        confirmation: review?.confirmation,
      }),
    onSettled: () => client.invalidateQueries({ queryKey: ['operations'] }),
  });
  const cancel = useMutation({
    mutationFn: () => apiClient.post(`/operations/${review?.operation.id}/cancel`),
    onSettled: () => client.invalidateQueries({ queryKey: ['operations'] }),
  });
  const current = operation.data?.data ?? review?.operation;
  const active = current?.status === 'RUNNING' || current?.status === 'CANCELLING';
  const pending = active || apply.isPending || cancel.isPending || prepare.isPending;
  const error = prepare.error ?? apply.error ?? cancel.error ?? operation.error;
  return (
    <Dialog
      open
      onOpenChange={(open) => {
        if (!open && !pending) onClose();
      }}
    >
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Review restore</DialogTitle>
          <DialogDescription>
            Review the archive and target before changing data. Cancellation may leave partial data.
          </DialogDescription>
        </DialogHeader>
        {allowDatabaseTarget && (
          <div className="space-y-2">
            <Label htmlFor="restore-database-target">Restore target</Label>
            <Select
              value={target || 'original'}
              disabled={pending}
              onValueChange={(value) => {
                setTarget(value === 'original' ? '' : value);
                setReview(undefined);
                setApproved(false);
              }}
            >
              <SelectTrigger id="restore-database-target">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="original">Original database</SelectItem>
                {candidates.map((database) => (
                  <SelectItem key={database.id} value={database.id}>
                    {database.name}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
            <p className="text-muted-foreground text-sm">
              Create an empty database with the same engine in this project to restore into a
              separate target.
            </p>
          </div>
        )}
        {!review && (
          <Button disabled={prepare.isPending} onClick={() => prepare.mutate()}>
            Prepare restore
          </Button>
        )}
        {current && (
          <>
            <p>{current.effects}</p>
            <p role="status">
              {current.status} · {current.phase}
            </p>
            {current.logs && (
              <pre className="max-h-48 overflow-auto whitespace-pre-wrap text-xs">
                {current.logs}
              </pre>
            )}
            {current.error && (
              <p role="alert" className="text-destructive">
                {current.error}
              </p>
            )}
          </>
        )}
        {current?.status === 'REVIEWED' && (
          <>
            <label className="flex gap-2 text-sm">
              <input
                type="checkbox"
                checked={approved}
                onChange={(event) => setApproved(event.target.checked)}
              />
              I approve the reviewed target and overwrite effects.
            </label>
            <Button
              variant="destructive"
              disabled={!approved || pending || current.expiresAt * 1000 <= Date.now()}
              onClick={() => apply.mutate()}
            >
              Apply restore
            </Button>
          </>
        )}
        {active && (
          <Button
            variant="destructive"
            disabled={cancel.isPending || current?.status === 'CANCELLING'}
            onClick={() => cancel.mutate()}
          >
            Interrupt restore
          </Button>
        )}
        {current && !active && current.status !== 'REVIEWED' && current.status !== 'COMPLETED' && (
          <Button
            disabled={pending}
            onClick={() => {
              setReview(undefined);
              setApproved(false);
              apply.reset();
              cancel.reset();
              prepare.mutate();
            }}
          >
            Prepare retry
          </Button>
        )}
        {error && (
          <p role="alert" className="text-destructive">
            {error.message}
          </p>
        )}
        <Button variant="outline" disabled={pending} onClick={onClose}>
          Close
        </Button>
      </DialogContent>
    </Dialog>
  );
}
