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
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import type { BaseResponse } from '#/interfaces/base';
import { apiClient } from '#/lib/api-client';
import { ClusterNodeFields } from './cluster-node-fields';
import type { Cluster, ClusterOperation, ClusterReview } from './cluster-types';

export function ClusterReviewDialog({
  projectId,
  cluster,
  action,
  onClose,
}: {
  projectId: string;
  cluster?: Cluster;
  action: 'install' | 'join' | 'upgrade' | 'remove' | 'recover';
  onClose: () => void;
}) {
  const client = useQueryClient();
  const [spec, setSpec] = useState<Cluster>(
    cluster ?? {
      id: '',
      name: '',
      version: '',
      nodes: [{ serverId: '', privateIp: '', interface: '', fingerprint: '' }],
      revision: 0,
      status: '',
      error: '',
    }
  );
  const [review, setReview] = useState<ClusterReview>();
  const [approved, setApproved] = useState(false);
  const prepare = useMutation({
    mutationFn: () =>
      apiClient.post<BaseResponse<ClusterReview>>(`/projects/${projectId}/clusters/review`, {
        cluster: spec,
        action,
      }),
    onSuccess: (response) => {
      setReview(response.data);
      setApproved(false);
    },
  });
  const operation = useQuery({
    queryKey: ['cluster-operation', review?.operation.id],
    queryFn: () =>
      apiClient.get<BaseResponse<ClusterOperation>>(`/cluster-operations/${review?.operation.id}`),
    enabled: !!review,
    refetchInterval: 1500,
  });
  const current = operation.data?.data ?? review?.operation;
  const apply = useMutation({
    mutationFn: () =>
      apiClient.post(`/cluster-operations/${review?.operation.id}/apply`, {
        confirmation: review?.confirmation,
      }),
    onSettled: () => client.invalidateQueries({ queryKey: ['cluster-operation'] }),
  });
  const cancel = useMutation({
    mutationFn: () => apiClient.post(`/cluster-operations/${review?.operation.id}/cancel`),
    onSettled: () => client.invalidateQueries({ queryKey: ['cluster-operation'] }),
  });
  const active = current?.status === 'RUNNING' || current?.status === 'CANCELLING';
  const pending = active || prepare.isPending || apply.isPending || cancel.isPending;
  const error = prepare.error ?? apply.error ?? cancel.error ?? operation.error;
  return (
    <Dialog
      open
      onOpenChange={(open) => {
        if (!open && !pending) onClose();
      }}
    >
      <DialogContent className="max-h-[85vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>Review cluster {action}</DialogTitle>
          <DialogDescription>
            Use registered Linux SSH servers on a restricted private network. One control plane
            provides no automatic failover. Check host fingerprints through your server provider.
          </DialogDescription>
        </DialogHeader>
        {!review && (
          <div className="space-y-3">
            <Label htmlFor="cluster-name">Cluster name</Label>
            <Input
              id="cluster-name"
              disabled={pending || action === 'remove'}
              value={spec.name}
              onChange={(event) => setSpec({ ...spec, name: event.target.value })}
            />
            <Label htmlFor="cluster-version">Exact K3s release</Label>
            <Input
              id="cluster-version"
              disabled={pending || action === 'join' || action === 'remove'}
              value={spec.version}
              placeholder="v1.x.y+k3s1"
              onChange={(event) => setSpec({ ...spec, version: event.target.value })}
            />
            {action === 'install' || action === 'join' ? (
              <ClusterNodeFields
                nodes={spec.nodes}
                locked={cluster?.nodes.length ?? 0}
                onChange={(nodes) => setSpec({ ...spec, nodes })}
              />
            ) : (
              <p>{spec.nodes.length} saved nodes</p>
            )}
            <Button disabled={pending} onClick={() => prepare.mutate()}>
              Check nodes and prepare
            </Button>
          </div>
        )}
        {current && (
          <div className="space-y-3">
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
              I approve these server and storage effects.
            </label>
            <Button
              variant="destructive"
              disabled={!approved || pending || current.expiresAt * 1000 <= Date.now()}
              onClick={() => apply.mutate()}
            >
              Apply cluster {action}
            </Button>
          </>
        )}
        {active && (
          <Button
            variant="destructive"
            disabled={cancel.isPending || current?.status === 'CANCELLING'}
            onClick={() => cancel.mutate()}
          >
            Interrupt operation
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
