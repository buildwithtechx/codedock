import { Clock, Database, Loader2, Play, Shield, Trash2 } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import { Checkbox } from '#/components/ui/checkbox';
import { Switch } from '#/components/ui/switch';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '#/components/ui/table';
import { useDelete, useTrigger, useUpdate } from '#/features/backups';
import type { BackupConfig } from './interfaces';

type BackupPoliciesProps = {
  configs: BackupConfig[];
  isLoading: boolean;
};

export function BackupPolicies({ configs, isLoading }: BackupPoliciesProps) {
  const triggerBackup = useTrigger();
  const updateBackup = useUpdate();
  const deleteBackup = useDelete();

  const [triggeringId, setTriggeringId] = useState<string | null>(null);

  const handleTrigger = async (id: string, name: string) => {
    setTriggeringId(id);
    try {
      await triggerBackup.mutateAsync({ id });
      toast.success(`Snapshot started for "${name}"`);
    } catch {
      toast.error('Failed to trigger snapshot');
    } finally {
      setTriggeringId(null);
    }
  };

  const handleToggleEnabled = async (config: BackupConfig) => {
    try {
      await updateBackup.mutateAsync({
        id: config.id,
        payload: {
          ...config,
          backupEnabled: !config.backupEnabled,
        },
      });
      toast.success(
        `Policy ${!config.backupEnabled ? 'enabled' : 'disabled'} for "${config.name}"`
      );
    } catch {
      toast.error('Failed to update policy status');
    }
  };

  const handleDelete = async (id: string, name: string) => {
    if (!window.confirm(`Are you sure you want to delete policy "${name}"?`)) {
      return;
    }
    try {
      await deleteBackup.mutateAsync({ id });
      toast.success('Policy deleted');
    } catch {
      toast.error('Failed to delete policy');
    }
  };

  return (
    <section className="rounded-2xl border border-border/60 bg-card p-5 shadow-sm">
      <div className="mb-4 flex items-center justify-between">
        <div className="flex items-center gap-2">
          <Clock className="h-4 w-4 text-primary" />
          <h2 className="font-semibold text-foreground text-sm">Backup policies & schedules</h2>
        </div>
        <span className="text-muted-foreground text-xs">{configs.length} configured</span>
      </div>

      {isLoading ? (
        <div className="flex min-h-40 items-center justify-center">
          <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
        </div>
      ) : configs.length === 0 ? (
        <div className="flex flex-col items-center justify-center py-12 text-center text-muted-foreground">
          <Shield className="mb-3 h-8 w-8 opacity-20" />
          <p className="font-medium text-sm">No backup policies configured</p>
          <p className="mt-1 text-xs">
            Policies define automated cron schedules and retention periods for databases.
          </p>
        </div>
      ) : (
        <div className="overflow-x-auto">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Target</TableHead>
                <TableHead>Schedule</TableHead>
                <TableHead>Retention</TableHead>
                <TableHead>Active</TableHead>
                <TableHead className="text-right">Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {configs.map((config) => (
                <TableRow key={config.id} className="hover:bg-muted/40">
                  <TableCell>
                    <div className="flex items-center gap-2.5">
                      <div className="flex h-7 w-7 items-center justify-center rounded-lg bg-muted text-muted-foreground">
                        <Database className="h-3.5 w-3.5" />
                      </div>
                      <div>
                        <p className="font-medium text-foreground text-sm">{config.name}</p>
                        <p className="text-muted-foreground text-xs">
                          {config.description || 'Database snapshot'}
                        </p>
                      </div>
                    </div>
                  </TableCell>
                  <TableCell>
                    <div className="font-mono text-xs">{config.schedule}</div>
                    <div className="text-muted-foreground text-xs">{config.timezone || 'UTC'}</div>
                  </TableCell>
                  <TableCell>
                    <div className="text-xs">
                      {config.retentionDays ? `${config.retentionDays} days` : 'Permanent'}
                    </div>
                    {config.maxBackups > 0 && (
                      <div className="text-muted-foreground text-xs">
                        max {config.maxBackups} copies
                      </div>
                    )}
                  </TableCell>
                  <TableCell>
                    <Switch
                      checked={config.backupEnabled}
                      onCheckedChange={() => handleToggleEnabled(config)}
                      aria-label="Toggle policy"
                    />
                    <label
                      htmlFor={`pre-deploy-${config.id}`}
                      className="mt-2 flex items-center gap-2 text-xs"
                    >
                      <Checkbox
                        id={`pre-deploy-${config.id}`}
                        checked={config.preDeployment ?? false}
                        disabled={!config.projectId || updateBackup.isPending}
                        onCheckedChange={(checked) =>
                          updateBackup.mutate({
                            id: config.id,
                            payload: { ...config, preDeployment: checked === true },
                          })
                        }
                        aria-label="Require verified backup before deployment"
                      />
                      Require verified backup before deployment
                    </label>
                  </TableCell>
                  <TableCell className="text-right">
                    <div className="flex items-center justify-end gap-1.5">
                      <Button
                        variant="outline"
                        size="sm"
                        className="h-7 text-xs"
                        onClick={() => handleTrigger(config.id, config.name)}
                        disabled={triggeringId !== null}
                      >
                        {triggeringId === config.id ? (
                          <Loader2 className="mr-1 h-3 w-3 animate-spin" />
                        ) : (
                          <Play className="mr-1 h-3 w-3" />
                        )}
                        Snapshot now
                      </Button>
                      <Button
                        variant="ghost"
                        size="icon"
                        className="h-7 w-7 text-destructive hover:bg-destructive/10 hover:text-destructive"
                        onClick={() => handleDelete(config.id, config.name)}
                      >
                        <Trash2 className="h-3.5 w-3.5" />
                      </Button>
                    </div>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </div>
      )}
    </section>
  );
}
