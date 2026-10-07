import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { useRef, useState } from 'react';
import type { Cluster } from '#/features/servers/cluster-types';
import type { CreateAppServiceRequest } from '#/features/services';
import type { RuntimeReview, RuntimeTarget } from '#/features/services/runtime-types';
import { useListByProject } from '#/hooks/use-environments';
import type { BaseResponse } from '#/interfaces/base';
import { apiClient } from '#/lib/api-client';
import { appsService } from '#/services/apps';
import { canvasService } from '#/services/canvas';
import { deploymentsService } from '#/services/deployments';
import { serviceVariablesService } from '#/services/service-variables';
import { useOrganizationStore } from '#/stores/organization-store';
import {
  type ApplicationSetupProps,
  parseSetupVariables,
  type RepositoryInspection,
  setupDefaults,
} from './application-setup-types';
import {
  describeDestination,
  emptyDestination,
  type SetupDestination,
  toRuntimeTarget,
} from './destination-picker';

export function useApplicationSetup(props: ApplicationSetupProps) {
  const client = useQueryClient();
  const organizationId = useOrganizationStore((state) => state.activeOrganizationId);
  const [projectId, setProjectId] = useState(props.projectId);
  const [selectedEnvironment, setEnvironment] = useState('');
  const [source, setSource] = useState(props.initialSource ?? 'git');
  const [draft, setDraft] = useState<CreateAppServiceRequest>({ ...setupDefaults, projectId });
  const [variables, setVariables] = useState('');
  const [review, setReview] = useState(false);
  const [error, setError] = useState('');
  const [detection, setDetection] = useState<RepositoryInspection>();
  const [deploymentId, setDeploymentId] = useState('');
  const [destination, setDestination] = useState<SetupDestination>(emptyDestination);
  const created = useRef('');
  const creationId = useRef(crypto.randomUUID());
  const savedVariables = useRef(new Set<string>());
  const requestSequence = useRef(0);
  const environments = useListByProject(projectId);
  const environmentId =
    environments.data?.data?.find((env) => env.id === selectedEnvironment)?.id ??
    environments.data?.data?.find((env) => env.isDefault)?.id ??
    environments.data?.data?.[0]?.id ??
    '';
  const projects = useQuery({
    queryKey: ['canvas', 'setup-projects', organizationId],
    queryFn: () => canvasService.listCanvasSummaries(organizationId),
    enabled: !!organizationId && props.isOpen,
  });
  const project = useQuery({
    queryKey: ['setup-target', projectId],
    queryFn: () => apiClient.get<BaseResponse<{ serverId?: string }>>(`/projects/${projectId}`),
    enabled: !!projectId && props.isOpen,
  });
  const serverId = project.data?.data.serverId ?? '';
  const clusters = useQuery({
    queryKey: ['setup-clusters', projectId],
    queryFn: () =>
      apiClient.get<BaseResponse<Cluster[]>>(`/projects/${projectId}/runtime-clusters`),
    enabled: !!projectId && props.isOpen && destination.kind === 'kubernetes',
  });
  const registries = useQuery({
    queryKey: ['setup-registries', projectId],
    queryFn: () =>
      apiClient.get<BaseResponse<{ id: string; registryUrl: string }[]>>(
        `/projects/${projectId}/registries`
      ),
    enabled: !!projectId && props.isOpen && destination.kind === 'kubernetes',
  });
  const payload = {
    ...draft,
    id: creationId.current,
    projectId,
    environmentId,
    repositoryUrl: source === 'git' ? draft.repositoryUrl : '',
    imageRef: source === 'image' ? draft.imageRef : '',
  };
  const update = <K extends keyof CreateAppServiceRequest>(
    key: K,
    value: CreateAppServiceRequest[K]
  ) => {
    requestSequence.current++;
    setDraft((previous) => ({ ...previous, [key]: value }));
    setReview(false);
    setError('');
  };
  const inspect = useMutation({
    mutationFn: async () => {
      const sequence = ++requestSequence.current;
      const result = await apiClient.post<BaseResponse<RepositoryInspection>>(
        `/projects/${projectId}/repository-inspection`,
        {
          repositoryUrl: draft.repositoryUrl,
          branch: draft.branch,
          rootDirectory: draft.rootDirectory,
        }
      );
      if (sequence !== requestSequence.current) return;
      setDetection(result.data);
      const {
        framework: _framework,
        packageManager: _manager,
        warnings: _warnings,
        ...defaults
      } = result.data;
      setDraft((previous) => ({ ...previous, ...defaults }));
      setReview(false);
    },
    onError: (err: Error) => setError(err.message),
  });
  const validate = () => {
    if (!environmentId) throw new Error('Select an available environment.');
    if (!draft.name.trim()) throw new Error('Enter an application name.');
    if (source === 'git' && !draft.repositoryUrl.trim()) throw new Error('Enter a repository URL.');
    if (source === 'image' && !draft.imageRef?.trim()) throw new Error('Enter an image reference.');
    if (
      !Number.isInteger(draft.internalPort) ||
      draft.internalPort < 1 ||
      draft.internalPort > 65535
    )
      throw new Error('Port must be between 1 and 65535.');
    if (
      !Number.isFinite(draft.cpuLimit ?? 0) ||
      (draft.cpuLimit ?? 0) < 0 ||
      !Number.isInteger(draft.memoryLimit ?? 0) ||
      (draft.memoryLimit ?? 0) < 0
    )
      throw new Error('CPU and memory limits must be non-negative numbers.');
    if (project.isError || !project.data)
      throw new Error('Deployment target could not be verified.');
    if (destination.kind === 'kubernetes' && !destination.clusterId)
      throw new Error('Select a ready cluster for the Kubernetes destination.');
    if (destination.kind === 'bare') {
      if (!destination.bareNode.serverId.trim()) throw new Error('Select a native server.');
      if (!destination.bareReleaseUrl.startsWith('https://'))
        throw new Error('Native artifact URL must use https.');
      if (!/^[0-9a-f]{64}$/i.test(destination.bareSha256.trim()))
        throw new Error('Native artifact SHA256 must be 64 hex characters.');
    }
    parseSetupVariables(variables);
  };
  const apply = useMutation({
    mutationFn: async (deploy: boolean) => {
      validate();
      if (!created.current) {
        const app = await appsService.createApp(environmentId, payload);
        created.current = app.data.id;
      }
      const target: RuntimeTarget = toRuntimeTarget(destination);
      if (target.kind !== 'docker') {
        const review = await apiClient.post<BaseResponse<RuntimeReview>>(
          `/apps/${created.current}/runtime/review`,
          { target, revision: 0 }
        );
        await apiClient.post(`/apps/${created.current}/runtime/apply`, {
          operationId: review.data.operation.id,
          confirmation: review.data.confirmation,
        });
      }
      for (const variable of parseSetupVariables(variables)) {
        if (savedVariables.current.has(variable.key)) continue;
        await serviceVariablesService.create(created.current, variable);
        savedVariables.current.add(variable.key);
      }
      await client.invalidateQueries({ queryKey: ['apps'] });
      await client.invalidateQueries({ queryKey: ['canvas'] });
      if (deploy) {
        const result = await deploymentsService.trigger(created.current);
        setDeploymentId(result.data.id);
      } else props.onOpenChange(false);
    },
    onError: (err: Error) =>
      setError(
        `${created.current ? 'Application saved. Retry to finish configuring this same application. ' : ''}${err.message}`
      ),
  });
  return {
    projectId,
    environmentId,
    projects,
    environments,
    project,
    serverId,
    clusters,
    registries,
    destination,
    destinationSummary: describeDestination(destination, serverId),
    runtimeTarget: toRuntimeTarget(destination),
    source,
    draft,
    variables,
    review,
    error,
    detection,
    deploymentId,
    inspect,
    apply,
    payload,
    locked: !!created.current || apply.isPending,
    update,
    setVariables: (value: string) => {
      setVariables(value);
      setReview(false);
    },
    setSource: (value: 'git' | 'image') => {
      requestSequence.current++;
      setSource(value);
      setReview(false);
    },
    setProjectId: (id: string) => {
      requestSequence.current++;
      setProjectId(id);
      setEnvironment('');
      setReview(false);
    },
    setDestination: (next: SetupDestination) => {
      requestSequence.current++;
      setDestination(next);
      setReview(false);
    },
    setEnvironment: (id: string) => {
      setEnvironment(id);
      setReview(false);
    },
    edit: () => setReview(false),
    prepareReview: () => {
      try {
        validate();
        setError('');
        setReview(true);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Invalid setup');
      }
    },
  };
}
