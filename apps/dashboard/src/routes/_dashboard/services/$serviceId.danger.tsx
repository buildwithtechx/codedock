import { createFileRoute, useNavigate } from '@tanstack/react-router';
import { AlertTriangle, PowerOff, RefreshCw, Trash2 } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Button } from '@codedock/ui/button';
import { Dialog, DialogContent, DialogTitle } from '@codedock/ui/dialog';
import { Input } from '@codedock/ui/input';
import { Label } from '@codedock/ui/label';
import { SettingsSection } from '#/features/settings/settings-section';
import { useDeleteApp, useGetApp, useRestartApp, useStopApp } from '#/hooks/use-apps';
import {
  useDeleteDatabase,
  useGetDatabase,
  useRestartDatabase,
  useStopDatabase,
} from '#/hooks/use-databases';

export const Route = createFileRoute('/_dashboard/services/$serviceId/danger')({
  component: DangerZoneRoute,
});

function DangerZoneRoute() {
  const { serviceId } = Route.useParams();
  const navigate = useNavigate();

  const { data: appData } = useGetApp(serviceId);
  const { data: dbData } = useGetDatabase(serviceId);
  const deleteApp = useDeleteApp();
  const deleteDb = useDeleteDatabase();
  const restartApp = useRestartApp();
  const stopApp = useStopApp();
  const restartDb = useRestartDatabase();
  const stopDb = useStopDatabase();

  const [openDelete, setOpenDelete] = useState(false);
  const [confirmName, setConfirmName] = useState('');

  const app = appData?.data;
  const db = dbData?.data;
  const isDatabase = !!db;
  const resourceName = app?.name || db?.name || 'service';
  const projectId = app?.projectId || db?.projectId;
  const isMatched = confirmName.trim() === resourceName.trim();
  const isDeleting = deleteApp.isPending || deleteDb.isPending;

  const handleDelete = async () => {
    if (!isMatched) {
      toast.error('Service name does not match');
      return;
    }

    try {
      if (isDatabase) {
        await deleteDb.mutateAsync({ id: serviceId });
      } else {
        await deleteApp.mutateAsync({ appId: serviceId });
      }
      toast.success(`"${resourceName}" was deleted`);
      setOpenDelete(false);
      void navigate({ to: `/projects/${projectId}` });
    } catch {
      toast.error('Failed to delete service');
    }
  };

  const handleRestart = async () => {
    try {
      if (isDatabase) {
        await restartDb.mutateAsync({ id: serviceId });
      } else {
        await restartApp.mutateAsync({ appId: serviceId });
      }
      toast.success('Service restart initiated');
    } catch {
      toast.error('Failed to restart service');
    }
  };

  const handleStop = async () => {
    try {
      if (isDatabase) {
        await stopDb.mutateAsync({ id: serviceId });
      } else {
        await stopApp.mutateAsync({ appId: serviceId });
      }
      toast.success('Service stop initiated');
    } catch {
      toast.error('Failed to stop service');
    }
  };

  return (
    <div className="space-y-6">
      <SettingsSection
        icon={<AlertTriangle className="size-4 text-destructive" />}
        title="Danger Zone"
        description={`Runtime and irreversible destructive operations for ${resourceName}.`}
        iconBg="bg-destructive/10"
        iconColor="text-destructive"
      >
        <div className="divide-y divide-border/50">
          <div className="flex flex-col items-start justify-between gap-4 pb-4 sm:flex-row sm:items-center">
            <div className="space-y-1">
              <h3 className="font-medium text-foreground text-sm">Restart service</h3>
              <p className="max-w-xl text-muted-foreground text-xs leading-relaxed">
                Gracefully restart the container for{' '}
                <span className="font-semibold text-foreground">{resourceName}</span>. This may
                cause a brief period of downtime.
              </p>
            </div>
            <Button
              variant="outline"
              size="sm"
              onClick={() => void handleRestart()}
              disabled={restartApp.isPending || restartDb.isPending}
              className="shrink-0 gap-2"
            >
              <RefreshCw className="size-4" />
              Restart container
            </Button>
          </div>

          <div className="flex flex-col items-start justify-between gap-4 py-4 sm:flex-row sm:items-center">
            <div className="space-y-1">
              <h3 className="font-medium text-foreground text-sm">Stop service</h3>
              <p className="max-w-xl text-muted-foreground text-xs leading-relaxed">
                Stop the running container. The service will remain unavailable until started again.
              </p>
            </div>
            <Button
              variant="outline"
              size="sm"
              onClick={() => void handleStop()}
              disabled={stopApp.isPending || stopDb.isPending}
              className="shrink-0 gap-2"
            >
              <PowerOff className="size-4" />
              Stop container
            </Button>
          </div>

          <div className="flex flex-col items-start justify-between gap-4 pt-4 sm:flex-row sm:items-center">
            <div className="space-y-1">
              <h3 className="font-medium text-foreground text-sm">
                Delete this {isDatabase ? 'database' : 'application'}
              </h3>
              <p className="max-w-xl text-muted-foreground text-xs leading-relaxed">
                Permanently removes{' '}
                <span className="font-semibold text-foreground">{resourceName}</span>, its volumes,
                environment variables, and routing configuration. This action cannot be undone.
              </p>
            </div>
            <Button
              variant="destructive"
              size="sm"
              onClick={() => {
                setConfirmName('');
                setOpenDelete(true);
              }}
              disabled={isDeleting}
              className="shrink-0 gap-2"
            >
              <Trash2 className="size-4" />
              Delete {isDatabase ? 'database' : 'application'}
            </Button>
          </div>
        </div>
      </SettingsSection>

      <Dialog open={openDelete} onOpenChange={(open) => !open && setOpenDelete(false)}>
        <DialogContent className="border/50 gap-0 bg-card/95 p-0 backdrop-blur-xl sm:max-w-md [&>button]:hidden">
          <div className="p-5">
            <DialogTitle className="flex items-center gap-2 font-bold text-foreground text-xl tracking-tight">
              <Trash2 className="h-5 w-5 text-destructive" />
              Delete {isDatabase ? 'Database' : 'Application'}
            </DialogTitle>

            <p className="mt-3 text-muted-foreground text-sm leading-relaxed">
              Are you sure you want to delete{' '}
              <strong className="text-foreground">{resourceName}</strong>? All associated volumes,
              variables, and deployments will be immediately removed.
            </p>

            <div className="mt-4 rounded-xl border border-destructive/30 bg-destructive/5 p-3 text-destructive text-xs">
              This action is permanent and cannot be reversed.
            </div>

            <div className="mt-4 space-y-2">
              <Label htmlFor="confirm-service-input" className="text-muted-foreground text-xs">
                To confirm, type{' '}
                <strong className="font-mono text-foreground">{resourceName}</strong>:
              </Label>
              <Input
                id="confirm-service-input"
                value={confirmName}
                onChange={(e) => setConfirmName(e.target.value)}
                placeholder={resourceName}
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
              {isDeleting ? 'Deleting...' : `Delete ${isDatabase ? 'Database' : 'Application'}`}
            </Button>
          </div>
        </DialogContent>
      </Dialog>
    </div>
  );
}
