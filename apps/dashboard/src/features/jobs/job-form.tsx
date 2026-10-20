import { Check, Loader2 } from 'lucide-react';
import { useMemo, useState } from 'react';
import { toast } from 'sonner';
import { Button } from '@codedock/ui/button';
import { AddServerModal } from '#/features/servers/add-server-modal';
import { useWorkflows } from '#/hooks/use-actions';
import { useListByProject } from '#/hooks/use-apps';
import { useCreateJob, useUpdateJob } from '#/hooks/use-jobs';
import { useListAllProjects } from '#/hooks/use-projects';
import { useListServers } from '#/hooks/use-servers';
import type { AppService, Job } from '#/interfaces/services';
import { JobBasicsSection, JobScheduleSection } from './job-form-basics';
import type { KV } from './job-form-sections';
import {
  JobEnvironmentSection,
  JobReliabilitySection,
  JobServersSection,
} from './job-form-servers';
import { JobSummaryRail } from './job-form-summary';
import { JobNotificationsSection, JobTriggersSection } from './job-form-triggers';

export function JobForm({
  job,
  initialWorkflowId,
  onCancel,
  onSaved,
}: {
  job?: Job;
  initialWorkflowId?: string;
  onCancel: () => void;
  onSaved: (job: Job) => void;
}) {
  const editing = Boolean(job);
  const { data: workflows = [] } = useWorkflows();
  const { data: servers = [] } = useListServers();
  const { data: projects = [] } = useListAllProjects();

  const [mode, setMode] = useState<'command' | 'workflow'>(
    initialWorkflowId ? 'workflow' : 'command'
  );
  const [label, setLabel] = useState(job?.name ?? '');
  const [command, setCommand] = useState(job?.command ?? '');
  const [workflowId, setWorkflowId] = useState(initialWorkflowId ?? '');
  const [scheduleType, setScheduleType] = useState<'recurring' | 'once' | 'manual'>('recurring');
  const [cron, setCron] = useState(job?.schedule ?? '0 3 * * *');
  const [runAt, setRunAt] = useState('');
  const [serverIds, setServerIds] = useState<string[]>(
    servers.length > 0 ? [servers[0].id] : []
  );
  const [envRows, setEnvRows] = useState<KV[]>([]);
  const [secretRows, setSecretRows] = useState<KV[]>([]);
  const [timeoutSec, setTimeoutSec] = useState('300');
  const [maxAttempts, setMaxAttempts] = useState('1');
  const [backoffSec, setBackoffSec] = useState('0');
  const [triggerIds, setTriggerIds] = useState<string[]>([]);
  const [dependsOn, setDependsOn] = useState<string[]>([]);
  const [notifyChannels, setNotifyChannels] = useState<string[]>([]);
  const [notifyStates, setNotifyStates] = useState<string[]>(['failed']);
  const [addServerOpen, setAddServerOpen] = useState(false);

  const [projectId] = useState(job?.projectId ?? (projects[0]?.id || ''));
  const servicesQuery = useListByProject(projectId, undefined, Boolean(projectId));
  const services = (servicesQuery.data?.data as AppService[] | undefined) ?? [];
  const [serviceId] = useState(job?.serviceId ?? (services[0]?.id || ''));

  const createJobMutation = useCreateJob();
  const updateJobMutation = useUpdateJob();
  const saving = createJobMutation.isPending || updateJobMutation.isPending;

  const toggleItem = (list: string[], item: string, setter: (val: string[]) => void) => {
    setter(list.includes(item) ? list.filter((x) => x !== item) : [...list, item]);
  };

  const selectedWorkflow = workflows.find((w) => w.id === workflowId);
  const canSave = useMemo(() => {
    if (!label.trim() || saving) return false;
    if (mode === 'workflow' && !workflowId) return false;
    if (mode === 'command' && !command.trim()) return false;
    if (scheduleType === 'recurring' && !cron.trim()) return false;
    if (scheduleType === 'once' && !runAt) return false;
    return true;
  }, [label, saving, mode, workflowId, command, scheduleType, cron, runAt]);

  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault();
    if (!canSave) return;

    try {
      const resolvedServiceId = serviceId || (services[0]?.id ?? 'default-service');
      const resolvedSchedule =
        scheduleType === 'recurring' ? cron.trim() : scheduleType === 'once' ? `@once ${runAt}` : '@manual';

      if (editing && job) {
        const res = await updateJobMutation.mutateAsync({
          id: job.id,
          payload: {
            name: label.trim(),
            schedule: resolvedSchedule,
            command: mode === 'workflow' ? `workflow:${workflowId}` : command.trim(),
          },
        });
        toast.success(`Job ${label.trim()} saved`);
        onSaved(res.data);
      } else {
        const res = await createJobMutation.mutateAsync({
          serviceId: resolvedServiceId,
          name: label.trim(),
          schedule: resolvedSchedule,
          command: mode === 'workflow' ? `workflow:${workflowId}` : command.trim(),
        });
        toast.success(`Job ${label.trim()} created`);
        onSaved(res.data);
      }
    } catch {
      toast.error(editing ? 'Failed to save job' : 'Failed to create job');
    }
  };

  const scheduleRecap =
    scheduleType === 'recurring' ? cron : scheduleType === 'once' ? runAt || '—' : 'Manual';
  const retryRecap =
    parseInt(maxAttempts, 10) > 1 ? `${maxAttempts}× · ${backoffSec}s` : 'None';

  return (
    <form onSubmit={handleSubmit} className="grid grid-cols-1 items-start gap-6 pb-16 lg:grid-cols-[minmax(0,1fr)_340px]">
      <div className="min-w-0 space-y-5">
        <JobBasicsSection
          label={label}
          onLabelChange={setLabel}
          mode={mode}
          onModeChange={setMode}
          command={command}
          onCommandChange={setCommand}
          workflowId={workflowId}
          onWorkflowIdChange={setWorkflowId}
          workflows={workflows}
          editing={editing}
        />

        <JobScheduleSection
          scheduleType={scheduleType}
          onScheduleTypeChange={setScheduleType}
          cron={cron}
          onCronChange={setCron}
          runAt={runAt}
          onRunAtChange={setRunAt}
        />

        {mode === 'command' && (
          <>
            <JobEnvironmentSection
              envRows={envRows}
              setEnvRows={setEnvRows}
              secretRows={secretRows}
              setSecretRows={setSecretRows}
            />

            <JobServersSection
              servers={servers}
              serverIds={serverIds}
              onToggleServer={(id) => toggleItem(serverIds, id, setServerIds)}
              onOpenAddServer={() => setAddServerOpen(true)}
            />

            <JobReliabilitySection
              timeoutSec={timeoutSec}
              onTimeoutChange={setTimeoutSec}
              maxAttempts={maxAttempts}
              onMaxAttemptsChange={setMaxAttempts}
              backoffSec={backoffSec}
              onBackoffChange={setBackoffSec}
            />
          </>
        )}

        <JobTriggersSection
          triggerIds={triggerIds}
          onToggleTrigger={(id) => toggleItem(triggerIds, id, setTriggerIds)}
          dependsOn={dependsOn}
          onToggleDependsOn={(key) => toggleItem(dependsOn, key, setDependsOn)}
        />

        <JobNotificationsSection
          notifyChannels={notifyChannels}
          onToggleChannel={(id) => toggleItem(notifyChannels, id, setNotifyChannels)}
          notifyStates={notifyStates}
          onToggleState={(s) => toggleItem(notifyStates, s, setNotifyStates)}
        />

        <div className="flex items-center justify-end gap-2 pt-2">
          <Button type="button" variant="ghost" onClick={onCancel} className="rounded-xl">
            Cancel
          </Button>
          <Button type="submit" disabled={!canSave} className="gap-2 rounded-xl px-6">
            {saving ? <Loader2 className="size-4 animate-spin" /> : <Check className="size-4" />}
            {saving ? 'Saving…' : editing ? 'Save changes' : 'Create job'}
          </Button>
        </div>
      </div>

      <JobSummaryRail
        mode={mode}
        workflowName={selectedWorkflow?.name}
        label={label}
        command={command}
        scheduleRecap={scheduleRecap}
        isRecurring={scheduleType === 'recurring'}
        targetsCount={serverIds.length}
        retryRecap={retryRecap}
        timeoutSec={timeoutSec}
        notifyChannelsCount={notifyChannels.length}
      />

      <AddServerModal
        open={addServerOpen}
        initialMode="managed"
        onClose={() => setAddServerOpen(false)}
        onCreated={(srv) => {
          setServerIds((prev) => (prev.includes(srv.id) ? prev : [...prev, srv.id]));
          setAddServerOpen(false);
        }}
      />
    </form>
  );
}
