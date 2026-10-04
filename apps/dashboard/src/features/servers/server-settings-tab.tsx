import { useNavigate } from '@tanstack/react-router';
import { AlertTriangle, Save, Trash2 } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import { Card } from '#/components/ui/card';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import { useDeleteServer, useUpdateServer } from '#/hooks/use-servers';
import type { Server } from '#/interfaces/server';

interface ServerSettingsTabProps {
  server: Server;
}

export function ServerSettingsTab({ server }: ServerSettingsTabProps) {
  const navigate = useNavigate();
  const [name, setName] = useState(server.name);
  const [sshHost, setSshHost] = useState(server.sshHost || server.ipAddress);
  const [sshPort, setSshPort] = useState(String(server.sshPort || 22));
  const [sshUser, setSshUser] = useState(server.sshUser || 'root');
  const [confirmDelete, setConfirmDelete] = useState(false);

  const { mutateAsync: updateServer, isPending: isUpdating } = useUpdateServer();
  const { mutateAsync: deleteServer, isPending: isDeleting } = useDeleteServer();

  const handleUpdate = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim()) {
      toast.error('Server name cannot be empty');
      return;
    }
    try {
      await updateServer({
        id: server.id,
        payload: {
          name: name.trim(),
          sshHost: sshHost.trim(),
          sshPort: Number.parseInt(sshPort, 10) || 22,
          sshUser: sshUser.trim(),
        },
      });
      toast.success('Server settings updated');
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Failed to update server');
    }
  };

  const handleDelete = async () => {
    try {
      await deleteServer(server.id);
      toast.success('Server deleted');
      navigate({ to: '/servers' });
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Failed to delete server');
    }
  };

  return (
    <div className="space-y-6">
      <Card className="p-6">
        <h3 className="font-semibold text-sm">Server Configuration</h3>
        <p className="mt-0.5 text-muted-foreground text-xs">
          Modify the identification and connection parameters for this node.
        </p>

        <form onSubmit={handleUpdate} className="mt-5 space-y-4">
          <div className="space-y-1.5">
            <Label htmlFor="edit-srv-name" className="text-xs">
              Server Name
            </Label>
            <Input
              id="edit-srv-name"
              value={name}
              onChange={(e) => setName(e.target.value)}
              required
            />
          </div>

          {!server.isLocal && (
            <div className="grid gap-4 sm:grid-cols-3">
              <div className="space-y-1.5">
                <Label htmlFor="edit-srv-host" className="text-xs">
                  Host / IP Address
                </Label>
                <Input
                  id="edit-srv-host"
                  value={sshHost}
                  onChange={(e) => setSshHost(e.target.value)}
                  required
                />
              </div>
              <div className="space-y-1.5">
                <Label htmlFor="edit-srv-port" className="text-xs">
                  SSH Port
                </Label>
                <Input
                  id="edit-srv-port"
                  type="number"
                  value={sshPort}
                  onChange={(e) => setSshPort(e.target.value)}
                />
              </div>
              <div className="space-y-1.5">
                <Label htmlFor="edit-srv-user" className="text-xs">
                  SSH User
                </Label>
                <Input
                  id="edit-srv-user"
                  value={sshUser}
                  onChange={(e) => setSshUser(e.target.value)}
                />
              </div>
            </div>
          )}

          <div className="flex justify-end pt-2">
            <Button type="submit" size="sm" disabled={isUpdating} className="gap-1.5 text-xs">
              <Save className="h-3.5 w-3.5" />
              {isUpdating ? 'Saving...' : 'Save Changes'}
            </Button>
          </div>
        </form>
      </Card>

      <Card className="border-destructive/30 bg-destructive/5 p-6">
        <div className="flex items-start gap-3">
          <div className="flex h-9 w-9 shrink-0 items-center justify-center rounded-xl bg-destructive/10 text-destructive">
            <AlertTriangle className="h-4 w-4" />
          </div>
          <div className="flex-1">
            <h3 className="font-semibold text-destructive text-sm">Danger Zone</h3>
            <p className="mt-1 text-muted-foreground text-xs leading-relaxed">
              Deleting this server removes the worker association. Any containers scheduled to this
              node will stop receiving updates and health checks.
            </p>

            {confirmDelete ? (
              <div className="mt-4 flex flex-wrap items-center gap-3">
                <Button
                  variant="destructive"
                  size="sm"
                  onClick={handleDelete}
                  disabled={isDeleting}
                  className="gap-1.5 text-xs"
                >
                  <Trash2 className="h-3.5 w-3.5" />
                  {isDeleting ? 'Deleting...' : 'Confirm Deletion'}
                </Button>
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => setConfirmDelete(false)}
                  className="text-xs"
                >
                  Cancel
                </Button>
              </div>
            ) : (
              <div className="mt-4">
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => setConfirmDelete(true)}
                  className="border-destructive/30 text-destructive text-xs hover:bg-destructive hover:text-destructive-foreground"
                >
                  Delete Server
                </Button>
              </div>
            )}
          </div>
        </div>
      </Card>
    </div>
  );
}
