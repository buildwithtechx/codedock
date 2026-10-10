import {
  ArrowRight,
  BookOpen,
  CalendarClock,
  Check,
  ListCheck,
  Loader2,
  Server,
  Terminal,
} from 'lucide-react';
import type { ReactNode } from 'react';
import { useEffect, useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '#/components/ui/select';
import { Textarea } from '#/components/ui/textarea';
import { useListAllProjects } from '#/features/projects';
import type { AppService, Job } from '#/features/services';
import { appsService } from '#/services/apps';
import { jobsApi } from './api';

const SCHEDULE_PRESETS = [
  { label: 'Every 5 minutes', value: '*/5 * * * *' },
  { label: 'Hourly', value: '0 * * * *' },
  { label: 'Daily at midnight', value: '0 0 * * *' },
  { label: 'Weekly on Sunday', value: '0 0 * * 0' },
];

function Section({
  title,
  icon,
  children,
}: {
  title: string;
  icon?: ReactNode;
  children: ReactNode;
}) {
  return (
    <section className="rounded-2xl border border-border/50 bg-card">
      <div className="flex items-center gap-3 border-border/50 border-b px-5 py-4">
        {icon && (
          <div className="flex size-9 shrink-0 items-center justify-center rounded-xl bg-primary/10 text-primary">
            {icon}
          </div>
        )}
        <h2 className="font-semibold text-[15px] text-foreground">{title}</h2>
      </div>
      <div className="space-y-4 p-5">{children}</div>
    </section>
  );
}

function SummaryRow({ label, value, mono }: { label: string; value: string; mono?: boolean }) {
  return (
    <div className="flex items-center justify-between gap-3">
      <span className="shrink-0 text-muted-foreground text-sm">{label}</span>
      <span
        className={`min-w-0 truncate text-right text-foreground text-sm ${mono ? 'font-mono' : ''}`}
      >
        {value}
      </span>
    </div>
  );
}

export function JobForm({
  job,
  onCancel,
  onSaved,
}: {
  job?: Job;
  onCancel: () => void;
  onSaved: (job: Job) => void;
}) {
  const editing = !!job;
  const { data: projects = [], isLoading: projectsLoading } = useListAllProjects();
  const [projectId, setProjectId] = useState(job?.projectId ?? '');
  const [services, setServices] = useState<AppService[]>([]);
  const [servicesLoading, setServicesLoading] = useState(false);
  const [serviceId, setServiceId] = useState(job?.serviceId ?? '');
  const [name, setName] = useState(job?.name ?? '');
  const [schedule, setSchedule] = useState(job?.schedule ?? '0 * * * *');
  const [command, setCommand] = useState(job?.command ?? '');
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    if (!projectId) {
      setServices([]);
      if (!editing) setServiceId('');
      return;
    }
    let cancelled = false;
    setServicesLoading(true);
    appsService
      .listByProject(projectId)
      .then((response) => {
        if (cancelled) return;
        setServices(response.data ?? []);
      })
      .catch(() => {
        if (!cancelled) toast.error('Could not load services for this project');
      })
      .finally(() => {
        if (!cancelled) setServicesLoading(false);
      });
    return () => {
      cancelled = true;
    };
  }, [projectId, editing]);

  const valid =
    serviceId !== '' && name.trim() !== '' && schedule.trim() !== '' && command.trim() !== '';

  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault();
    if (!valid || saving) return;
    setSaving(true);
    try {
      if (editing && job) {
        const response = await jobsApi.update(job.id, {
          name: name.trim(),
          schedule: schedule.trim(),
          command: command.trim(),
        });
        toast.success(`Job ${name.trim()} saved`);
        onSaved(response.data);
      } else {
        const response = await jobsApi.create({
          serviceId,
          name: name.trim(),
          schedule: schedule.trim(),
          command: command.trim(),
        });
        toast.success(`Job ${name.trim()} created`);
        onSaved(response.data);
      }
    } catch {
      toast.error(editing ? 'Failed to save job' : 'Failed to create job');
    } finally {
      setSaving(false);
    }
  };

  const selectedProject = projects.find((p: { id: string; name: string }) => p.id === projectId);
  const selectedService = services.find((s) => s.id === serviceId);

  return (
    <form onSubmit={handleSubmit}>
      <div className="grid grid-cols-1 items-start gap-6 lg:grid-cols-[minmax(0,1fr)_340px]">
        <div className="space-y-6">
          <Section title="Basics" icon={<Terminal className="size-4" />}>
            <div className="space-y-1.5">
              <Label htmlFor="job-name">Name</Label>
              <Input
                id="job-name"
                placeholder="nightly-cleanup"
                value={name}
                onChange={(event) => setName(event.target.value)}
              />
            </div>
            <div className="space-y-1.5">
              <Label htmlFor="job-command">Command</Label>
              <Textarea
                id="job-command"
                placeholder="./scripts/cleanup.sh --prune"
                value={command}
                onChange={(event) => setCommand(event.target.value)}
                rows={3}
                spellCheck={false}
                className="resize-y font-mono text-sm"
              />
              <p className="text-muted-foreground/60 text-xs">
                Runs inside the target service container on every tick of the schedule.
              </p>
            </div>
          </Section>

          <Section title="Schedule" icon={<CalendarClock className="size-4" />}>
            <div className="space-y-1.5">
              <Label htmlFor="job-schedule">Cron expression</Label>
              <Input
                id="job-schedule"
                placeholder="0 * * * *"
                value={schedule}
                spellCheck={false}
                className="font-mono"
                onChange={(event) => setSchedule(event.target.value)}
              />
              <div className="flex flex-wrap gap-2 pt-1">
                {SCHEDULE_PRESETS.map((preset) => (
                  <Button
                    key={preset.value}
                    type="button"
                    variant={schedule === preset.value ? 'secondary' : 'outline'}
                    size="sm"
                    onClick={() => setSchedule(preset.value)}
                  >
                    {preset.label}
                  </Button>
                ))}
              </div>
            </div>
          </Section>

          <Section title="Target service" icon={<Server className="size-4" />}>
            <div className="grid gap-5 sm:grid-cols-2">
              <div className="space-y-1.5">
                <Label>Project</Label>
                <Select value={projectId} onValueChange={setProjectId} disabled={projectsLoading}>
                  <SelectTrigger>
                    <SelectValue placeholder="Select a project" />
                  </SelectTrigger>
                  <SelectContent>
                    {projects.map((project: { id: string; name: string }) => (
                      <SelectItem key={project.id} value={project.id}>
                        {project.name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
              <div className="space-y-1.5">
                <Label>Service</Label>
                <Select
                  value={serviceId}
                  onValueChange={setServiceId}
                  disabled={!projectId || editing}
                >
                  <SelectTrigger>
                    <SelectValue
                      placeholder={
                        servicesLoading
                          ? 'Loading services…'
                          : editing
                            ? 'Service cannot be changed'
                            : 'Select a service'
                      }
                    />
                  </SelectTrigger>
                  <SelectContent>
                    {services.map((service) => (
                      <SelectItem key={service.id} value={service.id}>
                        {service.name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
            </div>
            {editing && (
              <p className="text-muted-foreground/60 text-xs">
                A job stays attached to its service; create a new job to run elsewhere.
              </p>
            )}
          </Section>

          <div className="flex items-center gap-3">
            <Button type="submit" disabled={!valid || saving} className="gap-2">
              {saving ? <Loader2 className="size-4 animate-spin" /> : <Check className="size-4" />}
              {saving
                ? editing
                  ? 'Saving…'
                  : 'Creating…'
                : editing
                  ? 'Save changes'
                  : 'Create job'}
            </Button>
            <Button type="button" variant="outline" onClick={onCancel}>
              Cancel
            </Button>
          </div>
        </div>

        <div className="space-y-4 lg:sticky lg:top-6 lg:self-start">
          <div className="rounded-2xl border border-border/60 bg-card p-5">
            <div className="mb-4 flex items-center gap-2">
              <ListCheck className="size-4 text-muted-foreground/70" />
              <h3 className="font-medium text-[14px] text-foreground">Job Summary</h3>
            </div>
            <div className="space-y-2.5">
              <SummaryRow
                label="Command"
                value={name.trim() || command.trim() || 'None'}
                mono={!name.trim() && Boolean(command.trim())}
              />
              <SummaryRow label="Schedule" value={schedule || 'None'} mono />
              <SummaryRow label="Target" value={selectedService?.name || 'None'} />
              <SummaryRow label="Project" value={selectedProject?.name || 'None'} />
            </div>
          </div>

          <a
            href="https://codedock.io/docs/jobs"
            target="_blank"
            rel="noreferrer"
            className="flex items-center gap-3 rounded-2xl border border-border/60 bg-card p-4 transition-colors hover:border-border"
          >
            <div className="flex size-9 shrink-0 items-center justify-center rounded-xl bg-muted">
              <BookOpen className="size-[18px] text-muted-foreground" />
            </div>
            <div className="min-w-0 flex-1">
              <p className="font-medium text-foreground text-sm">Cron job documentation</p>
              <p className="truncate text-muted-foreground/70 text-xs">
                Learn about cron expressions, timeouts, and execution logs
              </p>
            </div>
            <ArrowRight className="size-4 shrink-0 text-muted-foreground/40" />
          </a>
        </div>
      </div>
    </form>
  );
}
