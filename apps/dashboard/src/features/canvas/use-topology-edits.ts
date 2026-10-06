import { useMutation, useQueryClient } from '@tanstack/react-query';
import type { Connection } from '@xyflow/react';
import { useState } from 'react';
import { apiClient } from '#/lib/api-client';
import { createsCycle, type OperationalCanvas, type TopologyEdge } from './topology-types';

export function useTopologyEdits(canvas: OperationalCanvas) {
  const client = useQueryClient();
  const [draft, setDraft] = useState<{ revision: string; edges: TopologyEdge[] }>();
  const [review, setReview] = useState(false);
  const [error, setError] = useState('');
  const dependencies =
    draft?.edges ??
    canvas.edges.filter(
      (edge) => edge.kind === 'dependency' && !edge.id.startsWith('stack-dependency:')
    );
  const connect = (connection: Connection) => {
    const { source, target } = connection;
    if (!target.startsWith('app-') || (!source.startsWith('app-') && !source.startsWith('db-'))) {
      setError('Connect a service or database to an application that depends on it.');
      return;
    }
    if (dependencies.some((edge) => edge.source === source && edge.target === target)) return;
    if (createsCycle(dependencies, source, target)) {
      setError('This dependency creates a cycle.');
      return;
    }
    const edge: TopologyEdge = {
      id: `dependency:${source}:${target}`,
      source,
      target,
      kind: 'dependency',
      label: 'depends on',
    };
    setDraft({ revision: draft?.revision ?? canvas.revision, edges: [...dependencies, edge] });
    setReview(false);
    setError('');
  };
  const remove = (ids: Set<string>) => {
    setDraft({
      revision: draft?.revision ?? canvas.revision,
      edges: dependencies.filter((edge) => !ids.has(edge.id)),
    });
    setReview(false);
  };
  const apply = useMutation({
    mutationFn: () =>
      apiClient.put(`/environments/${canvas.environment.id}/canvas`, {
        revision: draft?.revision,
        dependencies: draft?.edges,
      }),
    onSuccess: async () => {
      setDraft(undefined);
      setReview(false);
      setError('');
      await client.invalidateQueries({ queryKey: ['canvas'] });
    },
    onError: (err: Error) => setError(err.message),
  });
  return {
    dependencies,
    connect,
    remove,
    draft,
    review,
    setReview,
    error,
    apply,
    discard: () => {
      setDraft(undefined);
      setReview(false);
      setError('');
    },
  };
}
