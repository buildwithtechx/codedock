import { Button } from '#/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '#/components/ui/card';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '#/components/ui/select';
import type { S3Destination } from '#/features/backups/interfaces';

type DatabaseBackupScheduleFormProps = {
  s3DestinationId: string;
  setS3DestinationId: (id: string) => void;
  schedule: string;
  setSchedule: (s: string) => void;
  retentionDays: string;
  setRetentionDays: (d: string) => void;
  destinations: S3Destination[];
  onSave: (e: React.FormEvent) => void;
  onCancel: () => void;
  isPending: boolean;
};

export function DatabaseBackupScheduleForm({
  s3DestinationId,
  setS3DestinationId,
  schedule,
  setSchedule,
  retentionDays,
  setRetentionDays,
  destinations,
  onSave,
  onCancel,
  isPending,
}: DatabaseBackupScheduleFormProps) {
  return (
    <Card className="border-border/60 bg-muted/20">
      <CardHeader className="pb-3">
        <CardTitle className="text-sm">Automated Schedule & Storage Policy</CardTitle>
        <CardDescription className="text-xs">
          Configure snapshot frequency, target storage destination, and retention window.
        </CardDescription>
      </CardHeader>
      <CardContent>
        <form onSubmit={onSave} className="space-y-4">
          <div className="grid gap-4 sm:grid-cols-3">
            <div className="space-y-1.5">
              <Label htmlFor="backup-storage-target" className="text-xs">
                Storage Target
              </Label>
              <Select value={s3DestinationId} onValueChange={setS3DestinationId}>
                <SelectTrigger id="backup-storage-target" className="h-9 text-xs">
                  <SelectValue placeholder="Select target" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="local">Local Server Storage</SelectItem>
                  {destinations.map((d) => (
                    <SelectItem key={d.id} value={d.id}>
                      {d.name} ({d.provider}: {d.bucket})
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>

            <div className="space-y-1.5">
              <Label htmlFor="backup-frequency" className="text-xs">
                Schedule Frequency
              </Label>
              <Select value={schedule} onValueChange={setSchedule}>
                <SelectTrigger id="backup-frequency" className="h-9 text-xs">
                  <SelectValue placeholder="Select schedule" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="0 3 * * *">Daily at 03:00 UTC</SelectItem>
                  <SelectItem value="0 * * * *">Every hour</SelectItem>
                  <SelectItem value="0 */6 * * *">Every 6 hours</SelectItem>
                  <SelectItem value="0 0 * * 0">Weekly (Sunday midnight)</SelectItem>
                  <SelectItem value="manual">Manual only</SelectItem>
                </SelectContent>
              </Select>
            </div>

            <div className="space-y-1.5">
              <Label htmlFor="backup-retention" className="text-xs">
                Retention Window (Days)
              </Label>
              <Input
                id="backup-retention"
                type="number"
                min="1"
                max="365"
                value={retentionDays}
                onChange={(e) => setRetentionDays(e.target.value)}
                className="h-9 text-xs"
                required
              />
            </div>
          </div>

          <div className="flex justify-end gap-2 pt-2">
            <Button type="button" variant="ghost" size="sm" onClick={onCancel}>
              Cancel
            </Button>
            <Button type="submit" size="sm" disabled={isPending}>
              Save Schedule
            </Button>
          </div>
        </form>
      </CardContent>
    </Card>
  );
}
