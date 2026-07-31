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
import { useCreateApp } from '#/hooks/use-apps';
import { useTrigger as useTriggerDeployment } from '#/hooks/use-deployments';
import { useListByProject as useListEnvironments } from '#/hooks/use-environments';

const schema = z.object({
  name: z.string().min(1, 'Name is required'),
  imageRef: z.string().min(1, 'Image reference is required (e.g. nginx:latest)'),
  internalPort: z.number(),
  startCommand: z.string().optional(),
});

type FormData = z.infer<typeof schema>;

interface Props {
  isOpen: boolean;
  onOpenChange: (open: boolean) => void;
  projectId: string;
}

export function CreateDockerImageModal({ isOpen, onOpenChange, projectId }: Props) {
  const { data: environments } = useListEnvironments(projectId);
  const environmentId = environments?.data?.[0]?.id || '';

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<FormData>({
    resolver: zodResolver(schema),
    defaultValues: {
      name: '',
      imageRef: '',
      internalPort: 80,
      startCommand: '',
    },
  });

  const createApp = useCreateApp();
  const triggerDeployment = useTriggerDeployment();

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
          imageRef: data.imageRef,
          repositoryUrl: '',
          branch: '',
          rootDirectory: '/',
          buildEngine: 'dockerfile',
          dockerfilePath: '',
          internalPort: data.internalPort,
          installCommand: '',
          buildCommand: '',
          startCommand: data.startCommand || '',
          runtimeMode: 'web',
          domain: '',
          staticOutput: '',
          healthCheckPath: '',
        },
      },
      {
        onSuccess: (res) => {
          const app = (res as { data?: { id: string } })?.data || res;
          if (app && 'id' in app && typeof app.id === 'string') {
            triggerDeployment.mutate(
              { serviceId: app.id },
              {
                onSuccess: () => {
                  toast.success('App created and deployment triggered successfully');
                  onOpenChange(false);
                },
                onError: () => {
                  toast.success('App created successfully (deploy container manually)');
                  onOpenChange(false);
                },
              }
            );
          } else {
            toast.success('App created successfully');
            onOpenChange(false);
          }
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
          <DialogTitle>Deploy Docker Image</DialogTitle>
          <DialogDescription>Deploy a pre-built Docker image from a registry.</DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit(onSubmit)} className="mt-4 space-y-4">
          <div className="space-y-2">
            <Label>App Name</Label>
            <Input {...register('name')} placeholder="my-docker-app" />
            {errors.name && <p className="text-destructive text-sm">{errors.name.message}</p>}
          </div>

          <div className="space-y-2">
            <Label>Image Name</Label>
            <Input {...register('imageRef')} placeholder="nginx:latest or username/repo:tag" />
            {errors.imageRef && (
              <p className="text-destructive text-sm">{errors.imageRef.message}</p>
            )}
          </div>

          <div className="space-y-2">
            <Label>Internal Port</Label>
            <Input type="number" {...register('internalPort', { valueAsNumber: true })} />
          </div>

          <div className="space-y-2">
            <Label>Start Command (optional)</Label>
            <Input {...register('startCommand')} placeholder="Overrides the container CMD" />
          </div>

          <div className="flex justify-end gap-2 pt-4">
            <Button variant="outline" type="button" onClick={() => onOpenChange(false)}>
              Cancel
            </Button>
            <Button type="submit" disabled={createApp.isPending || triggerDeployment.isPending}>
              {createApp.isPending || triggerDeployment.isPending ? 'Deploying…' : 'Deploy Image'}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
}
