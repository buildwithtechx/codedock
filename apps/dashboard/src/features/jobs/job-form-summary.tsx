import { ArrowRight, Book, ListCheck } from 'lucide-react';
import { SummaryRow } from './job-form-sections';

const DOCS_URL = 'https://docs.codedock.io/guides/jobs';

interface JobSummaryRailProps {
  mode: 'command' | 'workflow';
  workflowName?: string;
  label: string;
  command: string;
  scheduleRecap: string;
  isRecurring: boolean;
  targetsCount: number;
  retryRecap: string;
  timeoutSec: string;
  notifyChannelsCount: number;
}

export function JobSummaryRail({
  mode,
  workflowName,
  label,
  command,
  scheduleRecap,
  isRecurring,
  targetsCount,
  retryRecap,
  timeoutSec,
  notifyChannelsCount,
}: JobSummaryRailProps) {
  const displayCommand =
    mode === 'workflow'
      ? workflowName || '—'
      : label.trim() || command.trim() || '—';

  return (
    <div className="space-y-4 lg:sticky lg:top-6 lg:self-start">
      <div className="rounded-2xl border border-border/60 bg-card p-5">
        <div className="mb-4 flex items-center gap-2">
          <ListCheck className="size-4 text-muted-foreground/70" />
          <h3 className="font-medium text-foreground text-sm">Summary</h3>
        </div>
        <div className="space-y-2.5">
          <SummaryRow
            label={mode === 'workflow' ? 'Workflow' : 'Command'}
            value={displayCommand}
            mono={mode === 'command' && !label.trim() && Boolean(command.trim())}
          />
          <SummaryRow
            label="Schedule"
            value={scheduleRecap}
            mono={isRecurring}
          />
          {mode === 'command' && (
            <>
              <SummaryRow label="Targets" value={String(targetsCount)} />
              <SummaryRow label="Retry" value={retryRecap} />
              <SummaryRow label="Timeout" value={timeoutSec.trim() ? `${timeoutSec}s` : '—'} />
            </>
          )}
          <SummaryRow
            label="Notifications"
            value={notifyChannelsCount > 0 ? String(notifyChannelsCount) : '—'}
          />
        </div>
      </div>

      <a
        href={DOCS_URL}
        target="_blank"
        rel="noreferrer"
        className="flex items-center gap-3 rounded-2xl border border-border/60 bg-card p-4 transition-colors hover:border-border"
      >
        <div className="flex size-9 shrink-0 items-center justify-center rounded-xl bg-muted">
          <Book className="size-[18px] text-muted-foreground" />
        </div>
        <div className="min-w-0 flex-1">
          <p className="font-medium text-foreground text-sm">Learn about jobs</p>
          <p className="truncate text-muted-foreground/70 text-xs">
            Schedules, retries, secrets & notifications
          </p>
        </div>
        <ArrowRight className="size-4 shrink-0 text-muted-foreground/40 rtl:rotate-180" />
      </a>
    </div>
  );
}
