import { useMutation, useQuery } from '@tanstack/react-query';
import { useRef, useState } from 'react';
import { Button } from '#/components/ui/button';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import { useListByProject } from '#/hooks/use-environments';
import type { BaseResponse } from '#/interfaces/base';
import { apiClient } from '#/lib/api-client';
import { parseSetupVariables } from '../sources/application-setup-types';
import { ComposeStackCard } from './compose-stack-card';
import type { ComposeStack, ComposeStackReview } from './compose-stack-types';

export function ComposeDeployForm({ projectId }: { projectId: string }) {
  const [id, setId] = useState<string>(() => crypto.randomUUID());
  const [revision, setRevision] = useState(0);
  const [repositoryUrl, setRepositoryUrl] = useState('');
  const [branch, setBranch] = useState('main');
  const [rootDirectory, setRootDirectory] = useState('/');
  const [name, setName] = useState('application-stack');
  const [environment, setEnvironment] = useState('');
  const [content, setContent] = useState(
    'services:\n  web:\n    image: nginx:alpine\n    ports:\n      - "8080:80"\n'
  );
  const [variables, setVariables] = useState('');
  const [review, setReview] = useState<{ source: string; result: ComposeStackReview }>();
  const [error, setError] = useState('');
  const sequence = useRef(0);
  const environments = useListByProject(projectId);
  const environmentId =
    environments.data?.data?.find((item) => item.id === environment)?.id ??
    environments.data?.data?.find((item) => item.isDefault)?.id ??
    environments.data?.data?.[0]?.id ??
    '';
  const stacks = useQuery({
    queryKey: ['compose-stacks', projectId],
    queryFn: () => apiClient.get<BaseResponse<ComposeStack[]>>(`/projects/${projectId}/stacks`),
    refetchInterval: 5000,
  });
  const source = JSON.stringify({
    projectId,
    id,
    environmentId,
    name,
    content,
    variables,
    revision,
    repositoryUrl,
    branch,
    rootDirectory,
  });
  const currentReview = review?.source === source ? review.result : undefined;
  const request = () => ({
    id,
    environmentId,
    name,
    content,
    repositoryUrl,
    branch,
    rootDirectory,
    variables: Object.fromEntries(
      parseSetupVariables(variables).map((item) => [item.key, item.value])
    ),
    revision,
  });
  const inspect = useMutation({
    mutationFn: async () => {
      if (!environmentId) throw new Error('Select an environment.');
      const current = ++sequence.current;
      const result = await apiClient.post<BaseResponse<ComposeStackReview>>(
        `/projects/${projectId}/stacks/review`,
        request()
      );
      if (sequence.current === current) setReview({ source, result: result.data });
    },
    onError: (err: Error) => setError(err.message),
  });
  const save = useMutation({
    mutationFn: async (deploy: boolean) => {
      if (!currentReview) throw new Error('Review this exact configuration before saving.');
      const result = await apiClient.post<BaseResponse<ComposeStack>>(
        `/projects/${projectId}/stacks`,
        { ...request(), digest: currentReview.digest }
      );
      setReview(undefined);
      setRevision(result.data.revision);
      await stacks.refetch();
      if (deploy) {
        await apiClient.post(`/projects/${projectId}/stacks/${result.data.id}/deploy`);
        await stacks.refetch();
      }
    },
    onError: (err: Error) => setError(err.message),
  });
  const change = (setter: (value: string) => void, value: string) => {
    sequence.current++;
    setter(value);
    setReview(undefined);
    setError('');
  };
  return (
    <div className="grid gap-6 lg:grid-cols-2">
      <section className="space-y-4 rounded-xl border p-5">
        <h2 className="font-semibold text-lg">Deploy a Compose stack</h2>
        <Button
          variant="outline"
          size="sm"
          onClick={() => {
            sequence.current++;
            setId(crypto.randomUUID());
            setRevision(0);
            setName('application-stack');
            setReview(undefined);
            setError('');
            save.reset();
            inspect.reset();
          }}
        >
          New stack
        </Button>
        <p className="text-muted-foreground text-sm">
          Review the resolved configuration before saving. Services, variables, ports, named
          volumes, networks, health checks and dependencies run through Docker Compose. Unsupported
          fields are rejected. Set a Git repository to resolve relative build contexts. Host bind
          mounts, file secrets and external resources are rejected before apply.
        </p>
        <fieldset disabled={save.isPending || save.isSuccess} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="stack-name">Stack name</Label>
            <Input
              id="stack-name"
              value={name}
              onChange={(event) => change(setName, event.target.value)}
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="stack-environment">Environment</Label>
            <select
              id="stack-environment"
              className="w-full rounded-md border bg-background p-2"
              value={environmentId}
              onChange={(event) => change(setEnvironment, event.target.value)}
            >
              <option disabled value="">
                Select environment
              </option>
              {environments.data?.data?.map((item) => (
                <option key={item.id} value={item.id}>
                  {item.name}
                </option>
              ))}
            </select>
          </div>
          <div className="space-y-2">
            <Label htmlFor="stack-repository">
              Git repository for relative build contexts (optional)
            </Label>
            <Input
              id="stack-repository"
              value={repositoryUrl}
              onChange={(event) => change(setRepositoryUrl, event.target.value)}
            />
            <Label htmlFor="stack-branch">Branch or full commit SHA</Label>
            <Input
              id="stack-branch"
              value={branch}
              onChange={(event) => change(setBranch, event.target.value)}
            />
            <Label htmlFor="stack-root">Compose directory in repository</Label>
            <Input
              id="stack-root"
              value={rootDirectory}
              onChange={(event) => change(setRootDirectory, event.target.value)}
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="stack-content">Compose configuration</Label>
            <textarea
              id="stack-content"
              className="min-h-80 w-full rounded-md border bg-background p-3 font-mono text-sm"
              value={content}
              onChange={(event) => change(setContent, event.target.value)}
              spellCheck={false}
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="stack-variables">Interpolation variables (KEY=value)</Label>
            <textarea
              id="stack-variables"
              className="min-h-20 w-full rounded-md border bg-background p-3 font-mono text-sm"
              value={variables}
              onChange={(event) => change(setVariables, event.target.value)}
              spellCheck={false}
            />
          </div>
          <Button
            variant="outline"
            disabled={inspect.isPending || !content.trim()}
            onClick={() => inspect.mutate()}
          >
            {inspect.isPending ? 'Validating?' : 'Validate and review'}
          </Button>
        </fieldset>
        {error && (
          <p role="alert" className="text-destructive text-sm">
            {error}
          </p>
        )}
        {currentReview && (
          <div className="space-y-3">
            <h3 className="font-medium">Deployment review</h3>
            <ul className="list-disc space-y-1 pl-5 text-sm">
              {currentReview.effects.map((effect) => (
                <li key={effect}>{effect}</li>
              ))}
            </ul>
            <pre className="max-h-72 overflow-auto rounded-md bg-muted p-3 text-xs">
              {currentReview.config}
            </pre>
            <div className="flex gap-2">
              <Button
                variant="outline"
                disabled={save.isPending}
                onClick={() => save.mutate(false)}
              >
                Save stack
              </Button>
              <Button disabled={save.isPending} onClick={() => save.mutate(true)}>
                Save and deploy
              </Button>
            </div>
          </div>
        )}
      </section>
      <section className="space-y-4">
        <h2 className="font-semibold text-lg">Saved stacks</h2>
        {stacks.isError && (
          <p role="alert" className="text-destructive">
            Saved stacks unavailable.{' '}
            <Button onClick={() => stacks.refetch()} variant="outline">
              Retry
            </Button>
          </p>
        )}
        {stacks.data?.data?.map((stack) => (
          <ComposeStackCard
            key={stack.id}
            stack={stack}
            onEdit={(saved, config) => {
              sequence.current++;
              setId(saved.id);
              setRevision(saved.revision);
              setName(saved.name);
              setEnvironment(saved.environmentId);
              setContent(config);
              setRepositoryUrl('');
              setVariables('');
              setReview(undefined);
              setError('');
              save.reset();
              inspect.reset();
            }}
            onRefresh={() => {
              void stacks.refetch();
            }}
          />
        ))}
      </section>
    </div>
  );
}
