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
import { Label } from '#/components/ui/label';
import { Switch } from '#/components/ui/switch';
import { useCreate } from '#/hooks/use-backups';

export function VolumeBackupModal({
  serviceId,
  projectId,
  volumeName,
  open,
  onOpenChange,
}: {
  serviceId: string;
  projectId: string;
  volumeName: string;
  open: boolean;
  onOpenChange: (open: boolean) => void;
}) {
  const [name, setName] = useState(`${volumeName}-backup`);
  const [schedule, setSchedule] = useState('0 2 * * *');
  const [s3Enabled, setS3Enabled] = useState(false);
  const [disableLocal, setDisableLocal] = useState(false);
  const [retentionDays, setRetentionDays] = useState('7');
  const [maxBackups, setMaxBackups] = useState('0');

  const { mutateAsync: createBackup, isPending } = useCreate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await createBackup({
        payload: {
          projectId,
          serviceId,
          volumeName,
          name,
          description: `Backup for volume ${volumeName}`,
          dbUser: '',
          backupEnabled: true,
          s3Enabled,
          disableLocal,
          schedule,
          timezone: 'UTC',
          timeout: 3600,
          retentionDays: parseInt(retentionDays, 10),
          maxBackups: parseInt(maxBackups, 10),
          maxStorageGb: 0,
        },
      });
      toast.success('Volume backup configured successfully');
      onOpenChange(false);
    } catch {
      toast.error('Failed to configure volume backup');
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <form onSubmit={handleSubmit}>
          <DialogHeader>
            <DialogTitle>Configure Volume Backup</DialogTitle>
            <DialogDescription>
              Set up automated backups for volume:{' '}
              <code className="rounded bg-muted px-1 py-0.5 text-xs">{volumeName}</code>
            </DialogDescription>
          </DialogHeader>

          <div className="space-y-4 py-4">
            <div className="space-y-2">
              <Label>Configuration Name</Label>
              <Input value={name} onChange={(e) => setName(e.target.value)} required />
            </div>

            <div className="space-y-2">
              <Label>Cron Schedule</Label>
              <Input
                value={schedule}
                onChange={(e) => setSchedule(e.target.value)}
                placeholder="0 2 * * *"
                required
              />
              <p className="text-muted-foreground text-xs">Default: Daily at 2:00 AM UTC</p>
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label>Retention Days</Label>
                <Input
                  type="number"
                  value={retentionDays}
                  onChange={(e) => setRetentionDays(e.target.value)}
                  min="0"
                />
              </div>
              <div className="space-y-2">
                <Label>Max Backups</Label>
                <Input
                  type="number"
                  value={maxBackups}
                  onChange={(e) => setMaxBackups(e.target.value)}
                  min="0"
                />
              </div>
            </div>

            <div className="flex items-center justify-between">
              <div className="space-y-0.5">
                <Label>Enable S3 Upload</Label>
                <p className="text-muted-foreground text-xs">
                  Upload backups to configured S3 storage
                </p>
              </div>
              <Switch checked={s3Enabled} onCheckedChange={setS3Enabled} />
            </div>

            <div className="flex items-center justify-between">
              <div className="space-y-0.5">
                <Label>Disable Local Storage</Label>
                <p className="text-muted-foreground text-xs">
                  Do not save backups on the server disk
                </p>
              </div>
              <Switch checked={disableLocal} onCheckedChange={setDisableLocal} />
            </div>
          </div>

          <DialogFooter>
            <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>
              Cancel
            </Button>
            <Button type="submit" disabled={isPending}>
              {isPending ? 'Saving...' : 'Save Configuration'}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
