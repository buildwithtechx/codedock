import { useMutation, useQuery } from '@tanstack/react-query';
import { useState } from 'react';
import { Button } from '#/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '#/components/ui/dialog';
import type { BaseResponse } from '#/interfaces/base';
import type { Cluster } from '#/interfaces/clusters';
import type {
  ClusterDatabase,
  ClusterDatabaseSpec,
  DatabaseAction,
  DatabaseOperation,
  DatabaseReview,
} from '#/interfaces/databases';
import { apiClient } from '#/lib/api-client';
import { ClusterDatabaseForm } from './cluster-database-form';

export function ClusterDatabaseDialog({
  projectId,
  cluster,
  action,
  record,
  onClose,
}: {
  projectId: string;
  cluster: Cluster;
  action: DatabaseAction;
  record?: ClusterDatabase;
  onClose: () => void;
}) {
  const [spec, setSpec] = useState<ClusterDatabaseSpec>(
    action === 'recover' || action === 'backup'
      ? record!.spec
      : {
          id: crypto.randomUUID(),
          name: record ? `${record.spec.name.slice(0, 30)}-restore` : '',
          environmentId: record?.spec.environmentId ?? '',
          engine: record?.spec.engine ?? 'postgres',
          image: record?.spec.image ?? 'ghcr.io/cloudnative-pg/postgresql:17.6',
          instances: record?.spec.instances ?? 1,
          storageGiB: record?.spec.storageGiB ?? 10,
          storageClass: record?.spec.storageClass ?? '',
          s3DestinationId: record?.spec.s3DestinationId ?? '',
        }
  );
  const [versions, setVersions] = useState({
    operatorVersion: '1.28.4',
    barmanVersion: '0.15.1',
    certManagerVersion: '1.21.2',
  });
  const [restoreTime, setRestoreTime] = useState(new Date().toISOString());
  const [review, setReview] = useState<DatabaseReview>();
  const [approved, setApproved] = useState(false);
  const environments = useQuery({
    queryKey: ['project-environments', projectId],
    queryFn: () =>
      apiClient.get<BaseResponse<{ id: string; name: string }[]>>(
        `/projects/${projectId}/environments`
      ),
  });
  const destinations = useQuery({
    queryKey: ['s3-destinations'],
    queryFn: () => apiClient.get<BaseResponse<{ id: string; name: string }[]>>('/s3-destinations'),
  });
  const prepare = useMutation({
    mutationFn: () =>
      apiClient.post<BaseResponse<DatabaseReview>>(
        `/projects/${projectId}/clusters/${cluster.id}/databases/review`,
        {
          action,
          spec,
          ...versions,
          sourceId: action === 'restore' ? record?.id : undefined,
          restoreTime,
        }
      ),
    onSuccess: (response) => {
      setReview(response.data);
      setApproved(false);
    },
  });
  const operation = useQuery({
    queryKey: ['cluster-data-operation', review?.operation.id],
    queryFn: () =>
      apiClient.get<BaseResponse<DatabaseOperation>>(
        `/cluster-data-operations/${review?.operation.id}`
      ),
    enabled: !!review,
    refetchInterval: 1500,
  });
  const current = operation.data?.data ?? review?.operation;
  const apply = useMutation({
    mutationFn: () =>
      apiClient.post(`/cluster-data-operations/${review?.operation.id}/apply`, {
        confirmation: review?.confirmation,
      }),
    onSettled: () => operation.refetch(),
  });
  const cancel = useMutation({
    mutationFn: () => apiClient.post(`/cluster-data-operations/${review?.operation.id}/cancel`),
    onSettled: () => operation.refetch(),
  });
  const active = current?.status === 'RUNNING' || current?.status === 'CANCELLING';
  const error =
    prepare.error ??
    operation.error ??
    apply.error ??
    cancel.error ??
    environments.error ??
    destinations.error;

  return (
    <Dialog
      open
      onOpenChange={(open) => {
        if (!open) onClose();
      }}
    >
      <DialogContent className="max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>Review {action}</DialogTitle>
          <DialogDescription>
            Confirm the exact database and storage effects before execution. Persistent data remains
            during interruption and recovery.
          </DialogDescription>
        </DialogHeader>
        {error && (
          <p role="alert" className="text-destructive">
            {error.message}
          </p>
        )}
        {!review && (
          <ClusterDatabaseForm
            action={action}
            spec={spec}
            setSpec={setSpec}
            versions={versions}
            setVersions={setVersions}
            environments={environments.data?.data ?? []}
            destinations={destinations.data?.data ?? []}
            restoreTime={restoreTime}
            setRestoreTime={setRestoreTime}
            onPrepare={() => prepare.mutate()}
            isPreparing={prepare.isPending}
          />
        )}
        {current && (
          <div className="space-y-3">
            <p>{current.effects}</p>
            <p role="status">
              {current.status} ? {current.phase}
            </p>
            {current.logs && (
              <pre className="max-h-48 overflow-auto whitespace-pre-wrap text-xs">
                {current.logs}
              </pre>
            )}
            {current.error && <p role="alert">{current.error}</p>}
          </div>
        )}
        {current?.status === 'REVIEWED' && (
          <>
            <label className="flex gap-2 text-sm">
              <input
                type="checkbox"
                checked={approved}
                onChange={(event) => setApproved(event.target.checked)}
              />
              I approve the reviewed effects.
            </label>
            <Button disabled={!approved || apply.isPending} onClick={() => apply.mutate()}>
              Apply reviewed operation
            </Button>
          </>
        )}
        {active && (
          <Button variant="destructive" disabled={cancel.isPending} onClick={() => cancel.mutate()}>
            Interrupt execution
          </Button>
        )}
        <Button variant="outline" onClick={onClose}>
          Close
        </Button>
      </DialogContent>
    </Dialog>
  );
}
