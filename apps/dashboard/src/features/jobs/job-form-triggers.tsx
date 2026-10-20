import { ArrowRight, GitBranch, Bell } from 'lucide-react';
import { Link } from '@tanstack/react-router';
import { Checkbox } from '@codedock/ui/checkbox';
import { Field, Section, SECTION_TONES } from './job-form-sections';

export const EVENT_TRIGGERS = [
  {
    id: 'deployment.succeeded',
    label: 'Deploy succeeded',
    description: 'A production deploy finished successfully.',
  },
  {
    id: 'deployment.failed',
    label: 'Deploy failed',
    description: 'A build or deploy errored out.',
  },
  {
    id: 'backup_run.succeeded',
    label: 'Backup succeeded',
    description: 'A scheduled or manual backup completed.',
  },
  {
    id: 'backup_run.failed',
    label: 'Backup failed',
    description: 'A backup run errored out.',
  },
  {
    id: 'backup_restore.completed',
    label: 'Restore completed',
    description: 'A restore finished.',
  },
  {
    id: 'ssl.renewal_failed',
    label: 'SSL renewal failed',
    description: 'A certificate could not be renewed.',
  },
  {
    id: 'domain.verification_failed',
    label: 'Domain verification failed',
    description: "A domain's DNS check failed.",
  },
] as const;

interface JobTriggersSectionProps {
  triggerIds: string[];
  onToggleTrigger: (id: string) => void;
  otherJobs?: Array<{ key: string; label: string }>;
  dependsOn: string[];
  onToggleDependsOn: (key: string) => void;
}

export function JobTriggersSection({
  triggerIds,
  onToggleTrigger,
  otherJobs = [],
  dependsOn,
  onToggleDependsOn,
}: JobTriggersSectionProps) {
  return (
    <Section
      title="Dependencies & triggers"
      icon={<GitBranch className="size-[18px]" />}
      tone={SECTION_TONES.triggers}
    >
      {otherJobs.length > 0 && (
        <Field label="Run after">
          <div className="grid grid-cols-1 gap-1.5 sm:grid-cols-2">
            {otherJobs.map((o) => (
              <button
                key={o.key}
                type="button"
                onClick={() => onToggleDependsOn(o.key)}
                className={`flex w-full cursor-pointer items-center gap-2 rounded-xl border p-3 text-start text-sm transition-colors ${
                  dependsOn.includes(o.key)
                    ? 'border-primary/50 bg-primary/10 text-foreground'
                    : 'border-border/60 hover:bg-muted/40 text-muted-foreground'
                }`}
              >
                <Checkbox
                  checked={dependsOn.includes(o.key)}
                  className="pointer-events-none"
                />
                <span className="font-medium text-foreground">{o.label}</span>
              </button>
            ))}
          </div>
          <p className="mt-1.5 text-muted-foreground/60 text-xs">
            This job runs after the selected job(s) succeed.
          </p>
        </Field>
      )}

      <Field label="Event triggers">
        <div className="space-y-1">
          {EVENT_TRIGGERS.map((tg) => {
            const checked = triggerIds.includes(tg.id);
            return (
              <button
                key={tg.id}
                type="button"
                onClick={() => onToggleTrigger(tg.id)}
                aria-pressed={checked}
                className="flex w-full cursor-pointer items-start gap-2 rounded-lg px-1 py-1 text-start text-sm hover:bg-muted/40"
              >
                <Checkbox
                  checked={checked}
                  className="pointer-events-none mt-0.5"
                />
                <span className="text-sm">
                  <span className="text-foreground">{tg.label}</span>{' '}
                  <span className="text-muted-foreground/60">— {tg.description}</span>
                </span>
              </button>
            );
          })}
        </div>
      </Field>

      <div className="flex items-center gap-2 rounded-lg px-1 py-1 text-sm opacity-70">
        <Checkbox checked={false} disabled className="pointer-events-none" />
        <span className="text-muted-foreground">Custom trigger</span>
        <span className="rounded-full border border-border/60 bg-muted/40 px-1.5 py-0.5 font-medium text-[10px] text-muted-foreground tracking-wide uppercase">
          SOON
        </span>
      </div>
    </Section>
  );
}

interface JobNotificationsSectionProps {
  channels?: Array<{ id: string; label: string; kind: string }>;
  notifyChannels: string[];
  onToggleChannel: (id: string) => void;
  notifyStates: string[];
  onToggleState: (state: string) => void;
}

export function JobNotificationsSection({
  channels = [],
  notifyChannels,
  onToggleChannel,
  notifyStates,
  onToggleState,
}: JobNotificationsSectionProps) {
  return (
    <Section
      title="Notifications"
      icon={<Bell className="size-[18px]" />}
      tone={SECTION_TONES.notifications}
    >
      {channels.length === 0 ? (
        <div className="space-y-2">
          <p className="text-muted-foreground/60 text-xs">
            No notification channels yet — add one in Monitoring → Notifications.
          </p>
          <Link
            to={'/monitoring' as never}
            search={{ tab: 'notifications' } as never}
            className="inline-flex items-center gap-1.5 font-medium text-primary text-sm hover:underline"
          >
            Set up a channel <ArrowRight className="size-3.5 rtl:rotate-180" />
          </Link>
        </div>
      ) : (
        <>
          <div className="grid grid-cols-1 gap-1.5 sm:grid-cols-2">
            {channels.map((ch) => (
              <button
                key={ch.id}
                type="button"
                onClick={() => onToggleChannel(ch.id)}
                className={`flex w-full cursor-pointer items-center gap-2 rounded-xl border p-3 text-start text-sm transition-colors ${
                  notifyChannels.includes(ch.id)
                    ? 'border-primary/50 bg-primary/10 text-foreground'
                    : 'border-border/60 hover:bg-muted/40 text-muted-foreground'
                }`}
              >
                <Checkbox
                  checked={notifyChannels.includes(ch.id)}
                  className="pointer-events-none"
                />
                <span className="font-medium text-foreground">
                  {ch.label} ({ch.kind})
                </span>
              </button>
            ))}
          </div>

          {notifyChannels.length > 0 && (
            <div className="mt-3 flex gap-4">
              {['running', 'success', 'failed'].map((s) => (
                <button
                  key={s}
                  type="button"
                  onClick={() => onToggleState(s)}
                  aria-pressed={notifyStates.includes(s)}
                  className="flex cursor-pointer items-center gap-1.5 text-muted-foreground text-sm"
                >
                  <Checkbox
                    checked={notifyStates.includes(s)}
                    className="pointer-events-none"
                  />
                  <span className="capitalize">{s}</span>
                </button>
              ))}
            </div>
          )}
          <p className="mt-1.5 text-muted-foreground/60 text-xs">
            Overrides the global notification settings for this job.
          </p>
        </>
      )}
    </Section>
  );
}
