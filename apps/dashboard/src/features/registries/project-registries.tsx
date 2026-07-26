import { Loader2, Plus, Trash2 } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '#/components/ui/dialog';
import { Input } from '#/components/ui/input';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '#/components/ui/table';
import {
  useCreateRegistry,
  useDeleteRegistry,
  useListRegistries,
} from '#/features/registries/hooks';
import type { Registry } from '#/features/registries/types';

export function ProjectRegistries({ projectId }: { projectId: string }) {
  const { data: registries, isLoading } = useListRegistries(projectId);
  const [createOpen, setCreateOpen] = useState(false);

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <div>
          <h3 className="font-semibold text-lg">Private Docker Registries</h3>
          <p className="text-muted-foreground text-sm">
            Add private registries (e.g., Docker Hub, GHCR) to deploy private images.
          </p>
        </div>
        <Button onClick={() => setCreateOpen(true)} size="sm">
          <Plus className="mr-2 h-4 w-4" />
          Add Registry
        </Button>
      </div>

      <div className="rounded-md border">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Name</TableHead>
              <TableHead>Registry URL</TableHead>
              <TableHead>Username</TableHead>
              <TableHead className="w-[100px] text-right">Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {isLoading ? (
              <TableRow>
                <TableCell colSpan={4} className="h-24 text-center">
                  <Loader2 className="mx-auto h-6 w-6 animate-spin text-muted-foreground" />
                </TableCell>
              </TableRow>
            ) : registries?.length === 0 ? (
              <TableRow>
                <TableCell colSpan={4} className="h-24 text-center text-muted-foreground">
                  No private registries added.
                </TableCell>
              </TableRow>
            ) : (
              registries?.map((registry) => (
                <RegistryRow key={registry.id} registry={registry} projectId={projectId} />
              ))
            )}
          </TableBody>
        </Table>
      </div>

      <CreateRegistryModal projectId={projectId} open={createOpen} onOpenChange={setCreateOpen} />
    </div>
  );
}

function RegistryRow({ registry, projectId }: { registry: Registry; projectId: string }) {
  const { mutateAsync: deleteRegistry, isPending } = useDeleteRegistry(projectId);

  const handleDelete = async () => {
    if (!confirm('Are you sure you want to delete this registry?')) return;
    try {
      await deleteRegistry(registry.id);
      toast.success('Registry deleted successfully');
    } catch (err: any) {
      toast.error(err.message || 'Failed to delete registry');
    }
  };

  return (
    <TableRow>
      <TableCell className="font-medium">{registry.name}</TableCell>
      <TableCell>{registry.registryUrl}</TableCell>
      <TableCell>{registry.username}</TableCell>
      <TableCell className="text-right">
        <Button
          variant="ghost"
          size="icon"
          onClick={handleDelete}
          disabled={isPending}
          className="text-destructive hover:bg-destructive/10 hover:text-destructive"
        >
          {isPending ? (
            <Loader2 className="h-4 w-4 animate-spin" />
          ) : (
            <Trash2 className="h-4 w-4" />
          )}
        </Button>
      </TableCell>
    </TableRow>
  );
}

function CreateRegistryModal({
  projectId,
  open,
  onOpenChange,
}: {
  projectId: string;
  open: boolean;
  onOpenChange: (open: boolean) => void;
}) {
  const [name, setName] = useState('');
  const [registryUrl, setRegistryUrl] = useState('');
  const [username, setUsername] = useState('');
  const [passwordToken, setPasswordToken] = useState('');

  const { mutateAsync: createRegistry, isPending } = useCreateRegistry(projectId);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name || !registryUrl || !username || !passwordToken) {
      toast.error('All fields are required');
      return;
    }

    try {
      await createRegistry({ name, registryUrl, username, passwordToken });
      toast.success('Registry added successfully');
      onOpenChange(false);
      setName('');
      setRegistryUrl('');
      setUsername('');
      setPasswordToken('');
    } catch (err: any) {
      toast.error(err.message || 'Failed to add registry');
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <form onSubmit={handleSubmit}>
          <DialogHeader>
            <DialogTitle>Add Private Registry</DialogTitle>
            <DialogDescription>
              Connect a private Docker registry to deploy private images.
            </DialogDescription>
          </DialogHeader>
          <div className="space-y-4 py-4">
            <div className="space-y-2">
              <label htmlFor="name" className="font-medium text-sm">
                Name
              </label>
              <Input
                id="name"
                value={name}
                onChange={(e) => setName(e.target.value)}
                placeholder="e.g. My GitHub CR"
                disabled={isPending}
              />
            </div>
            <div className="space-y-2">
              <label htmlFor="url" className="font-medium text-sm">
                Registry URL
              </label>
              <Input
                id="url"
                value={registryUrl}
                onChange={(e) => setRegistryUrl(e.target.value)}
                placeholder="e.g. ghcr.io or docker.io"
                disabled={isPending}
              />
            </div>
            <div className="space-y-2">
              <label htmlFor="username" className="font-medium text-sm">
                Username
              </label>
              <Input
                id="username"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                placeholder="e.g. solomonolatunji"
                disabled={isPending}
              />
            </div>
            <div className="space-y-2">
              <label htmlFor="password" className="font-medium text-sm">
                Password / Access Token
              </label>
              <Input
                id="password"
                type="password"
                value={passwordToken}
                onChange={(e) => setPasswordToken(e.target.value)}
                placeholder="Required for authentication"
                disabled={isPending}
              />
            </div>
          </div>
          <DialogFooter>
            <Button
              type="button"
              variant="outline"
              onClick={() => onOpenChange(false)}
              disabled={isPending}
            >
              Cancel
            </Button>
            <Button type="submit" disabled={isPending}>
              {isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              Save Registry
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
