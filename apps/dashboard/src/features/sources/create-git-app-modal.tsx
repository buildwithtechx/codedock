import { zodResolver } from '@hookform/resolvers/zod';
import { useForm } from 'react-hook-form';
import { toast } from 'sonner';
import { z } from 'zod';
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
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '#/components/ui/select';
import type { BuildEngine } from '#/features/services/interfaces';
import { useCreateApp } from '#/hooks/use-apps';
import { useListByProject as useListEnvironments } from '#/hooks/use-environments';
import { useGitStatus, useListGitRepos } from '#/hooks/use-git';

const schema = z.object({
  name: z.string().min(1, 'Name is required'),
  provider: z.string().min(1, 'Provider is required'),
  repositoryUrl: z.string().min(1, 'Repository URL is required'),
  branch: z.string().min(1, 'Branch is required'),
  rootDirectory: z.string(),
  buildEngine: z.string(),
  dockerfilePath: z.string(),
  internalPort: z.number(),
  installCommand: z.string().optional(),
  buildCommand: z.string().optional(),
  startCommand: z.string().optional(),
});

type FormData = z.infer<typeof schema>;

interface Props {
  isOpen: boolean;
  onOpenChange: (open: boolean) => void;
  projectId: string;
}

const PROVIDER_LABELS: Record<string, string> = {
  github: 'GitHub',
  gitlab: 'GitLab',
  bitbucket: 'Bitbucket',
  gitea: 'Gitea',
};

export function CreateGitAppModal({ isOpen, onOpenChange, projectId }: Props) {
  const { data: environments } = useListEnvironments(projectId);
  const environmentId = environments?.data?.[0]?.id || '';

  const { data: gitStatus } = useGitStatus();
  const connectedProviders = gitStatus?.data || [];

  const {
    register,
    handleSubmit,
    setValue,
    watch,
    formState: { errors },
  } = useForm<FormData>({
    resolver: zodResolver(schema),
    defaultValues: {
      name: '',
      provider: '',
      repositoryUrl: '',
      branch: 'main',
      rootDirectory: '/',
      buildEngine: 'nixpacks',
      dockerfilePath: 'Dockerfile',
      internalPort: 3000,
      installCommand: '',
      buildCommand: '',
      startCommand: '',
    },
  });

  const provider = watch('provider');
  const buildEngine = watch('buildEngine');
  const repositoryUrl = watch('repositoryUrl');

  const { data: reposData, isLoading: reposLoading } = useListGitRepos(provider);
  const repos = reposData?.data || [];

  const createApp = useCreateApp();

  const onSubmit = (data: FormData) => {
    if (!environmentId) {
      toast.error('No environment found for this project');
      return;
    }

    createApp.mutate(
      {
        environmentId,
        payload: {
          projectId,
          name: data.name,
          repositoryUrl: data.repositoryUrl,
          branch: data.branch,
          rootDirectory: data.rootDirectory || '/',
          buildEngine: data.buildEngine as BuildEngine,
          dockerfilePath: data.dockerfilePath || 'Dockerfile',
          internalPort: data.internalPort || 3000,
          installCommand: data.installCommand || '',
          buildCommand: data.buildCommand || '',
          startCommand: data.startCommand || '',
          runtimeMode: 'web',
          domain: '',
          staticOutput: '',
          healthCheckPath: '',
        },
      },
      {
        onSuccess: () => {
          toast.success('App created successfully');
          onOpenChange(false);
        },
        onError: (err: Error) => {
          toast.error(err.message || 'Failed to create app');
        },
      }
    );
  };

  return (
    <Dialog open={isOpen} onOpenChange={onOpenChange}>
      <DialogContent className="max-h-[80vh] overflow-y-auto sm:max-w-md">
        <DialogHeader>
          <DialogTitle>Deploy Git Repository</DialogTitle>
          <DialogDescription>Deploy an application from a Git repository.</DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit((d) => onSubmit(d))} className="mt-4 space-y-4">
          <div className="space-y-2">
            <Label>App Name</Label>
            <Input {...register('name')} placeholder="my-awesome-app" />
            {errors.name && <p className="text-destructive text-sm">{errors.name.message}</p>}
          </div>

          <div className="space-y-2">
            <Label>Git Provider</Label>
            <Select value={provider} onValueChange={(val) => setValue('provider', val)}>
              <SelectTrigger>
                <SelectValue placeholder="Select a Git Provider" />
              </SelectTrigger>
              <SelectContent>
                {connectedProviders.map((p) => (
                  <SelectItem key={p.provider} value={p.provider}>
                    {PROVIDER_LABELS[p.provider] ?? p.provider}
                  </SelectItem>
                ))}
                <SelectItem value="public">Public Repository URL</SelectItem>
              </SelectContent>
            </Select>
            {errors.provider && (
              <p className="text-destructive text-sm">{errors.provider.message}</p>
            )}
          </div>

          {provider === 'public' ? (
            <div className="space-y-2">
              <Label>Repository URL</Label>
              <Input {...register('repositoryUrl')} placeholder="https://github.com/user/repo" />
              {errors.repositoryUrl && (
                <p className="text-destructive text-sm">{errors.repositoryUrl.message}</p>
              )}
            </div>
          ) : provider ? (
            <div className="space-y-2">
              <Label>Repository</Label>
              <Select value={repositoryUrl} onValueChange={(val) => setValue('repositoryUrl', val)}>
                <SelectTrigger>
                  <SelectValue placeholder={reposLoading ? 'Loading…' : 'Select a repository'} />
                </SelectTrigger>
                <SelectContent>
                  {repos.map((r) => (
                    <SelectItem key={r.cloneUrl} value={r.cloneUrl}>
                      {r.fullName}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
              {errors.repositoryUrl && (
                <p className="text-destructive text-sm">{errors.repositoryUrl.message}</p>
              )}
            </div>
          ) : null}

          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label>Branch</Label>
              <Input {...register('branch')} placeholder="main" />
            </div>
            <div className="space-y-2">
              <Label>Root Directory</Label>
              <Input {...register('rootDirectory')} placeholder="/" />
            </div>
          </div>

          <div className="space-y-2">
            <Label>Build Engine</Label>
            <Select value={buildEngine} onValueChange={(val) => setValue('buildEngine', val)}>
              <SelectTrigger>
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="nixpacks">Nixpacks (Auto-detect)</SelectItem>
                <SelectItem value="dockerfile">Dockerfile</SelectItem>
                <SelectItem value="buildpacks">Heroku Buildpacks</SelectItem>
                <SelectItem value="static">Static</SelectItem>
              </SelectContent>
            </Select>
          </div>

          {buildEngine === 'dockerfile' && (
            <div className="space-y-2">
              <Label>Dockerfile Path</Label>
              <Input {...register('dockerfilePath')} placeholder="Dockerfile" />
            </div>
          )}

          <div className="space-y-2">
            <Label>Internal Port</Label>
            <Input type="number" {...register('internalPort', { valueAsNumber: true })} />
          </div>

          {buildEngine !== 'dockerfile' && (
            <>
              <div className="space-y-2">
                <Label>Install Command (optional)</Label>
                <Input {...register('installCommand')} placeholder="npm install" />
              </div>
              <div className="space-y-2">
                <Label>Build Command (optional)</Label>
                <Input {...register('buildCommand')} placeholder="npm run build" />
              </div>
              <div className="space-y-2">
                <Label>Start Command (optional)</Label>
                <Input {...register('startCommand')} placeholder="npm start" />
              </div>
            </>
          )}

          <div className="flex justify-end gap-2 pt-4">
            <Button variant="outline" type="button" onClick={() => onOpenChange(false)}>
              Cancel
            </Button>
            <Button type="submit" disabled={createApp.isPending}>
              {createApp.isPending ? 'Deploying…' : 'Deploy App'}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
}
