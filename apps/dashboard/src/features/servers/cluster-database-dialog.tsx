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
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import type { BaseResponse } from '#/interfaces/base';
import { apiClient } from '#/lib/api-client';
import type {
  ClusterDatabase,
  ClusterDatabaseSpec,
  DatabaseAction,
  DatabaseOperation,
  DatabaseReview,
} from './cluster-database-types';
import type { Cluster } from './cluster-types';
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
          <div className="space-y-3">
            {action === 'operators'
              ? Object.entries(versions).map(([key, value]) => (
                  <div key={key}>
                    <Label htmlFor={key}>{key.replace(/Version$/, ' version')}</Label>
                    <Input
                      id={key}
                      value={value}
                      onChange={(event) => setVersions({ ...versions, [key]: event.target.value })}
                    />
                  </div>
                ))
              : (action === 'create' || action === 'restore') && (
                  <>
                    <Label htmlFor="cluster-db-name">Database name</Label>
                    <Input
                      id="cluster-db-name"
                      value={spec.name}
                      onChange={(event) => setSpec({ ...spec, name: event.target.value })}
                    />
                    <Label htmlFor="cluster-db-env">Environment</Label>
                    <select
                      id="cluster-db-env"
                      className="w-full rounded border bg-background p-2"
                      value={spec.environmentId}
                      onChange={(event) => setSpec({ ...spec, environmentId: event.target.value })}
                    >
                      <option value="">Select an environment</option>
                      {environments.data?.data.map((environment) => (
                        <option key={environment.id} value={environment.id}>
                          {environment.name}
                        </option>
                      ))}
                    </select>
                    <Label htmlFor="cluster-db-engine">Engine</Label>
                    <select
                      id="cluster-db-engine"
                      className="w-full rounded border bg-background p-2"
                      disabled={action === 'restore'}
                      value={spec.engine}
                      onChange={(event) => {
                        const engine = event.target.value as 'postgres' | 'redis';
                        setSpec({
                          ...spec,
                          engine,
                          image:
                            engine === 'redis'
                              ? 'redis:7.4.2'
                              : 'ghcr.io/cloudnative-pg/postgresql:17.6',
                          s3DestinationId: '',
                        });
                      }}
                    >
                      <option value="postgres">PostgreSQL</option>
                      <option value="redis">Redis</option>
                    </select>
                    <Label htmlFor="cluster-db-image">Exact engine image tag</Label>
                    <Input
                      id="cluster-db-image"
                      value={spec.image}
                      onChange={(event) => setSpec({ ...spec, image: event.target.value })}
                    />
                    <Label htmlFor="cluster-db-instances">Database mode</Label>
                    <select
                      id="cluster-db-instances"
                      className="w-full rounded border bg-background p-2"
                      value={spec.instances}
                      onChange={(event) =>
                        setSpec({ ...spec, instances: Number(event.target.value) })
                      }
                    >
                      <option value={1}>Standalone</option>
                      <option value={3}>3 replicated instances</option>
                    </select>
                    <Label htmlFor="cluster-db-storage">Storage per instance (GiB)</Label>
                    <Input
                      id="cluster-db-storage"
                      type="number"
                      min={1}
                      value={spec.storageGiB}
                      onChange={(event) =>
                        setSpec({ ...spec, storageGiB: Number(event.target.value) })
                      }
                    />
                    <Label htmlFor="cluster-db-class">
                      Storage class (empty uses cluster default)
                    </Label>
                    <Input
                      id="cluster-db-class"
                      value={spec.storageClass}
                      onChange={(event) => setSpec({ ...spec, storageClass: event.target.value })}
                    />
                    {spec.engine === 'postgres' && (
                      <>
                        <Label htmlFor="cluster-db-s3">Object-store backups</Label>
                        <select
                          id="cluster-db-s3"
                          className="w-full rounded border bg-background p-2"
                          value={spec.s3DestinationId}
                          onChange={(event) =>
                            setSpec({ ...spec, s3DestinationId: event.target.value })
                          }
                        >
                          <option value="">No object-store backups</option>
                          {destinations.data?.data.map((destination) => (
                            <option key={destination.id} value={destination.id}>
                              {destination.name}
                            </option>
                          ))}
                        </select>
                        <p className="text-muted-foreground text-xs">
                          Daily base backups and WAL archiving use a separate database prefix with
                          30-day retention. Generated credentials stay encrypted in Codedock and
                          Kubernetes Secrets.
                        </p>
                      </>
                    )}
                    {action === 'restore' && (
                      <>
                        <Label htmlFor="cluster-db-time">Recovery time (UTC ISO timestamp)</Label>
                        <Input
                          id="cluster-db-time"
                          value={restoreTime}
                          onChange={(event) => setRestoreTime(event.target.value)}
                        />
                        <p className="text-sm">
                          The source remains unchanged. The source must have a completed base backup
                          and continuous WAL coverage for this time.
                        </p>
                      </>
                    )}
                  </>
                )}
            <Button disabled={prepare.isPending} onClick={() => prepare.mutate()}>
              Prepare review
            </Button>
          </div>
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
