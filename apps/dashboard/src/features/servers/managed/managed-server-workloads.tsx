import { Layers } from 'lucide-react';
import { useEffect, useState } from 'react';
import { Button } from '@codedock/ui/button';
import { Input } from '@codedock/ui/input';
import { Textarea } from '@codedock/ui/textarea';
import type { ManagedWorkload } from '#/interfaces/server';
import {
  useControlManagedWorkload,
  useCreateManagedWorkload,
  useFetchManagedWorkloadLogs,
  useManagedWorkloads,
} from '#/hooks/use-managed-server';
import { ManagedControlPanel } from './managed-control-panel';
import { ManagedWorkloadRow } from './managed-workload-row';

export function ManagedServerWorkloads({ serverId }: { serverId: string }) {
  const {
    data: workloadsData,
    isLoading: loading,
    error: queryError,
    refetch,
  } = useManagedWorkloads(serverId);
  const createWorkloadMutation = useCreateManagedWorkload();
  const controlWorkloadMutation = useControlManagedWorkload();
  const fetchLogsMutation = useFetchManagedWorkloadLogs();

  const [workloads, setWorkloads] = useState<ManagedWorkload[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [createOpen, setCreateOpen] = useState(false);

  const [name, setName] = useState('');
  const [command, setCommand] = useState('');
  const [directory, setDirectory] = useState('/root');
  const [environment, setEnvironment] = useState('');
  const [policy, setPolicy] = useState('on-failure');
  const [submitting, setSubmitting] = useState(false);

  const [logs, setLogs] = useState<{ id: string; name: string; logs: string } | null>(null);
  const logsLoading = fetchLogsMutation.isPending;

  useEffect(() => {
    if (workloadsData?.workloads) {
      setWorkloads(workloadsData.workloads);
    }
  }, [workloadsData]);

  useEffect(() => {
    if (queryError) {
      setError(queryError instanceof Error ? queryError.message : 'Failed to load workloads');
    }
  }, [queryError]);

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim() || !command.trim()) return;
    setSubmitting(true);
    setError(null);
    try {
      const created = await createWorkloadMutation.mutateAsync({
        serverId,
        payload: {
          name: name.trim(),
          command: command.trim(),
          workingDirectory: directory.trim() || '/root',
          environment: environment.split(/\r?\n/).filter((l) => l.trim()),
          restartPolicy: policy,
          confirm: true,
        },
      });
      setWorkloads((prev) => [...prev, created]);
      setCreateOpen(false);
      setName('');
      setCommand('');
      setEnvironment('');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create workload');
    } finally {
      setSubmitting(false);
    }
  };

  const handleControl = async (
    workloadId: string,
    action: 'start' | 'stop' | 'restart' | 'delete'
  ) => {
    try {
      const res = await controlWorkloadMutation.mutateAsync({
        serverId,
        payload: { workloadId, action, confirm: true },
      });
      if (action === 'delete') {
        setWorkloads((prev) => prev.filter((w) => w.id !== workloadId));
        if (logs?.id === workloadId) setLogs(null);
      } else if (res.workload) {
        setWorkloads((prev) => prev.map((w) => (w.id === workloadId ? res.workload! : w)));
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : `Failed to ${action} workload`);
    }
  };

  const handleViewLogs = async (workload: ManagedWorkload) => {
    try {
      const res = await fetchLogsMutation.mutateAsync({
        serverId,
        workloadId: workload.id,
        tail: 100,
      });
      setLogs({ id: workload.id, name: workload.name, logs: res.logs });
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch logs');
    }
  };

  return (
    <ManagedControlPanel
      title="Workloads"
      description="Manage auxiliary daemon processes and background workers running alongside your containers."
      icon={Layers}
      busy={loading}
      error={error}
      refresh={() => void refetch()}
    >
      <p className="rounded-xl bg-muted/30 p-4 text-muted-foreground text-xs leading-relaxed">
        Background daemons run under systemd supervisor on this server with automated restart
        policies.
      </p>

      <div className="flex items-center justify-between">
        <h3 className="font-medium text-foreground text-sm">Active Daemons ({workloads.length})</h3>
        <Button
          size="sm"
          variant={createOpen ? 'ghost' : 'secondary'}
          onClick={() => setCreateOpen((v) => !v)}
        >
          {createOpen ? 'Cancel' : 'Create workload'}
        </Button>
      </div>

      {createOpen && (
        <form
          onSubmit={(e) => void handleCreate(e)}
          className="space-y-4 rounded-xl bg-muted/20 p-4"
        >
          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <div className="space-y-1.5">
              <label htmlFor="workload-name" className="font-medium text-foreground text-xs">
                Workload Name
              </label>
              <Input
                id="workload-name"
                placeholder="e.g. Queue Consumer"
                value={name}
                onChange={(e) => setName(e.target.value)}
                required
              />
            </div>
            <div className="space-y-1.5">
              <label htmlFor="workload-directory" className="font-medium text-foreground text-xs">
                Working Directory
              </label>
              <Input
                id="workload-directory"
                placeholder="/root"
                value={directory}
                onChange={(e) => setDirectory(e.target.value)}
              />
            </div>
          </div>

          <div className="space-y-1.5">
            <label htmlFor="workload-command" className="font-medium text-foreground text-xs">
              Execution Command
            </label>
            <Input
              id="workload-command"
              placeholder="e.g. node /srv/worker.js"
              value={command}
              onChange={(e) => setCommand(e.target.value)}
              required
            />
          </div>

          <div className="space-y-1.5">
            <label htmlFor="workload-env" className="font-medium text-foreground text-xs">
              Environment Variables (KEY=VALUE)
            </label>
            <Textarea
              id="workload-env"
              rows={3}
              placeholder="PORT=8080&#10;NODE_ENV=production"
              value={environment}
              onChange={(e) => setEnvironment(e.target.value)}
              className="font-mono text-xs"
            />
          </div>

          <div className="flex items-center justify-between pt-2">
            <select
              value={policy}
              onChange={(e) => setPolicy(e.target.value)}
              className="h-9 rounded-lg border border-input bg-background px-3 text-xs"
            >
              <option value="always">Always restart</option>
              <option value="on-failure">Restart on failure</option>
              <option value="no">Never restart</option>
            </select>
            <Button size="sm" type="submit" disabled={submitting}>
              {submitting ? 'Creating...' : 'Create Workload'}
            </Button>
          </div>
        </form>
      )}

      <div className="space-y-3">
        {workloads.map((workload) => (
          <ManagedWorkloadRow
            key={workload.id}
            workload={workload}
            onViewLogs={handleViewLogs}
            onControl={handleControl}
          />
        ))}
      </div>

      {logs && (
        <div className="space-y-2 rounded-xl bg-muted/20 p-4">
          <div className="flex items-center justify-between">
            <h4 className="font-medium text-foreground text-xs">Logs: {logs.name}</h4>
            <Button size="sm" variant="ghost" className="h-7 text-xs" onClick={() => setLogs(null)}>
              Close
            </Button>
          </div>
          <pre className="max-h-64 overflow-auto whitespace-pre-wrap rounded-lg bg-background p-3 font-mono text-muted-foreground text-xs leading-relaxed">
            {logsLoading ? 'Loading logs...' : logs.logs}
          </pre>
        </div>
      )}
    </ManagedControlPanel>
  );
}
