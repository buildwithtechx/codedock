import { CalendarClock, Terminal } from 'lucide-react';
import type { ActionWorkflow } from '#/interfaces/actions';
import { Field, SECTION_TONES, Section } from './job-form-sections';

const inputCls =
  'w-full rounded-xl border border-border/60 bg-background px-3 py-2 text-foreground text-sm outline-none transition-colors focus:border-primary/50';

interface JobBasicsSectionProps {
  label: string;
  onLabelChange: (val: string) => void;
  mode: 'command' | 'workflow';
  onModeChange: (val: 'command' | 'workflow') => void;
  command: string;
  onCommandChange: (val: string) => void;
  workflowId: string;
  onWorkflowIdChange: (val: string) => void;
  workflows: ActionWorkflow[];
  editing?: boolean;
}

export function JobBasicsSection({
  label,
  onLabelChange,
  mode,
  onModeChange,
  command,
  onCommandChange,
  workflowId,
  onWorkflowIdChange,
  workflows,
  editing,
}: JobBasicsSectionProps) {
  return (
    <Section title="Basics" icon={<Terminal className="size-[18px]" />} tone={SECTION_TONES.basics}>
      <Field label="Name">
        <input
          value={label}
          onChange={(e) => onLabelChange(e.target.value)}
          placeholder="Nightly DB vacuum"
          className={inputCls}
          autoFocus
        />
      </Field>

      {!editing && (
        <div className="flex border-b border-border/50">
          <button
            type="button"
            onClick={() => onModeChange('command')}
            className={`border-b-2 px-4 py-2 font-medium text-sm transition-colors ${
              mode === 'command'
                ? 'border-primary text-foreground'
                : 'border-transparent text-muted-foreground hover:text-foreground'
            }`}
          >
            Run command
          </button>
          <button
            type="button"
            onClick={() => onModeChange('workflow')}
            className={`border-b-2 px-4 py-2 font-medium text-sm transition-colors ${
              mode === 'workflow'
                ? 'border-primary text-foreground'
                : 'border-transparent text-muted-foreground hover:text-foreground'
            }`}
          >
            Run workflow
          </button>
        </div>
      )}

      {mode === 'workflow' ? (
        <div className="space-y-3">
          <Field label="Workflow">
            <select
              value={workflowId}
              onChange={(e) => {
                onWorkflowIdChange(e.target.value);
                const wf = workflows.find((w) => w.id === e.target.value);
                if (wf && !label.trim()) onLabelChange(wf.name);
              }}
              className={`${inputCls} h-10`}
            >
              <option value="" disabled>Select workflow</option>
              {workflows.map((wf) => (
                <option key={wf.id} value={wf.id}>
                  {wf.name}
                </option>
              ))}
            </select>
          </Field>
          <p className="text-muted-foreground text-xs">
            Trigger a GitHub Actions-compatible workflow run on a defined schedule.
          </p>
        </div>
      ) : (
        <Field label="Command">
          <textarea
            value={command}
            onChange={(e) => onCommandChange(e.target.value)}
            placeholder="pg_dump ... > /backups/db.sql"
            rows={3}
            spellCheck={false}
            className={`${inputCls} resize-y font-mono text-sm`}
          />
          <p className="mt-1.5 text-muted-foreground/60 text-xs">
            Runs on each selected server. Use docker run --rm for a container task.
          </p>
        </Field>
      )}
    </Section>
  );
}

interface JobScheduleSectionProps {
  scheduleType: 'recurring' | 'once' | 'manual';
  onScheduleTypeChange: (val: 'recurring' | 'once' | 'manual') => void;
  cron: string;
  onCronChange: (val: string) => void;
  runAt: string;
  onRunAtChange: (val: string) => void;
}

export function JobScheduleSection({
  scheduleType,
  onScheduleTypeChange,
  cron,
  onCronChange,
  runAt,
  onRunAtChange,
}: JobScheduleSectionProps) {
  return (
    <Section title="Schedule" icon={<CalendarClock className="size-[18px]" />} tone={SECTION_TONES.schedule}>
      <div className="grid grid-cols-3 gap-2">
        {[
          { id: 'recurring', label: 'On a schedule' },
          { id: 'once', label: 'Once' },
          { id: 'manual', label: 'Manual' },
        ].map((st) => (
          <button
            key={st.id}
            type="button"
            onClick={() => onScheduleTypeChange(st.id as any)}
            className={`rounded-xl border px-3 py-2.5 font-medium text-sm transition-colors ${
              scheduleType === st.id
                ? 'border-primary/60 bg-primary/10 text-foreground'
                : 'border-border/60 text-muted-foreground hover:bg-muted/50'
            }`}
          >
            {st.label}
          </button>
        ))}
      </div>

      {scheduleType === 'recurring' && (
        <Field label="Schedule (cron)">
          <input
            value={cron}
            onChange={(e) => onCronChange(e.target.value)}
            spellCheck={false}
            placeholder="0 3 * * *"
            className={`${inputCls} font-mono`}
          />
          <p className="mt-1.5 text-muted-foreground/60 text-xs">
            Standard 5-field cron (minute hour day month weekday), UTC.
          </p>
        </Field>
      )}

      {scheduleType === 'once' && (
        <Field label="Run at">
          <input
            type="datetime-local"
            value={runAt}
            onChange={(e) => onRunAtChange(e.target.value)}
            className={inputCls}
          />
        </Field>
      )}

      {scheduleType === 'manual' && (
        <p className="text-muted-foreground/70 text-sm">
          Runs only when triggered — Run Now, a dependency, or an event.
        </p>
      )}
    </Section>
  );
}
