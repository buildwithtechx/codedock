import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { useRef, useState } from 'react';
import type { CreateAppServiceRequest } from '#/features/services';
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
    if (project.data.data.serverId)
      throw new Error(
        'This setup requires a local Docker project. SSH worker setup is not supported here yet.'
      );
    parseSetupVariables(variables);
  };
  const apply = useMutation({
    mutationFn: async (deploy: boolean) => {
      validate();
      if (!created.current) {
        const app = await appsService.createApp(environmentId, payload);
        created.current = app.data.id;
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
