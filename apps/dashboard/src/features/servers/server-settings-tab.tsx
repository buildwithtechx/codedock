import { useNavigate } from '@tanstack/react-router';
import { AlertTriangle, Save, Server as ServerIcon, Trash2 } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Button } from '@codedock/ui/button';
import { Dialog, DialogContent, DialogTitle } from '@codedock/ui/dialog';
import { Input } from '@codedock/ui/input';
import { Label } from '@codedock/ui/label';
import { SettingsSection } from '#/features/settings/settings-section';
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
  const [openDelete, setOpenDelete] = useState(false);
  const [confirmName, setConfirmName] = useState('');

  const { mutateAsync: updateServer, isPending: isUpdating } = useUpdateServer();
  const { mutateAsync: deleteServer, isPending: isDeleting } = useDeleteServer();

  const isMatched = confirmName.trim() === server.name.trim();

  const handleUpdate = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim()) {
      toast.error('Server name cannot be empty');
      return;
    }
    const port = Number(sshPort);
    if (!Number.isInteger(port) || port < 1 || port > 65535) {
      toast.error('SSH port must be between 1 and 65535');
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
    if (!isMatched) {
      toast.error('Server name does not match');
      return;
    }
    try {
      await deleteServer(server.id);
      toast.success(`Server "${server.name}" was deleted`);
      setOpenDelete(false);
      void navigate({ to: '/servers' });
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Failed to delete server');
    }
  };

  return (
    <div className="space-y-6">
      <SettingsSection
        icon={<ServerIcon className="size-4 text-primary" />}
        title="Server Configuration"
        description="Modify the identification and connection parameters for this node."
      >
        <form onSubmit={handleUpdate} className="space-y-4">
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
                  min={1}
                  max={65535}
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
      </SettingsSection>

      <SettingsSection
        icon={<AlertTriangle className="size-4 text-destructive" />}
        title="Danger Zone"
        description="Irreversible destructive actions for this server node."
        iconBg="bg-destructive/10"
        iconColor="text-destructive"
      >
        <div className="flex flex-col items-start justify-between gap-4 sm:flex-row sm:items-center">
          <div className="space-y-1">
            <h3 className="font-medium text-foreground text-sm">Delete this server</h3>
            <p className="max-w-xl text-muted-foreground text-xs leading-relaxed">
              Permanently removes{' '}
              <span className="font-semibold text-foreground">{server.name}</span> from your fleet.
              Any containers scheduled to this node will stop receiving updates and health checks.
            </p>
          </div>
          <Button
            variant="destructive"
            size="sm"
            onClick={() => {
              setConfirmName('');
              setOpenDelete(true);
            }}
            className="shrink-0 gap-2"
          >
            <Trash2 className="size-4" />
            Delete server
          </Button>
        </div>
      </SettingsSection>

      <Dialog open={openDelete} onOpenChange={(open) => !open && setOpenDelete(false)}>
        <DialogContent className="border/50 gap-0 bg-card/95 p-0 backdrop-blur-xl sm:max-w-md [&>button]:hidden">
          <div className="p-5">
            <DialogTitle className="flex items-center gap-2 font-bold text-foreground text-xl tracking-tight">
              <Trash2 className="h-5 w-5 text-destructive" />
              Delete Server
            </DialogTitle>

            <p className="mt-3 text-muted-foreground text-sm leading-relaxed">
              Are you sure you want to delete{' '}
              <strong className="text-foreground">{server.name}</strong>? All worker associations
              and node metrics for this host will be immediately removed.
            </p>

            <div className="mt-4 rounded-xl border border-destructive/30 bg-destructive/5 p-3 text-destructive text-xs">
              This action is permanent and cannot be reversed.
            </div>

            <div className="mt-4 space-y-2">
              <Label htmlFor="confirm-server-input" className="text-muted-foreground text-xs">
                To confirm, type{' '}
                <strong className="font-mono text-foreground">{server.name}</strong>:
              </Label>
              <Input
                id="confirm-server-input"
                value={confirmName}
                onChange={(e) => setConfirmName(e.target.value)}
                placeholder={server.name}
                className="bg-muted/40 font-mono text-xs"
                autoFocus
              />
            </div>
          </div>

          <div className="flex items-center justify-end gap-3 p-5 pt-0">
            <Button
              type="button"
              variant="ghost"
              size="sm"
              onClick={() => {
                setOpenDelete(false);
                setConfirmName('');
              }}
            >
              Cancel
            </Button>
            <Button
              type="button"
              variant="destructive"
              size="sm"
              disabled={!isMatched || isDeleting}
              onClick={() => void handleDelete()}
            >
              {isDeleting ? 'Deleting...' : 'Delete Server'}
            </Button>
          </div>
        </DialogContent>
      </Dialog>
    </div>
  );
}
