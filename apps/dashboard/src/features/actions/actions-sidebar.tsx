import { Link } from '@tanstack/react-router';
import {
  ArrowRight,
  ArrowUpRight,
  BookOpen,
  Check,
  Cloud,
  Play,
  PlayCircle,
  Server,
} from 'lucide-react';
import { Button } from '#/components/ui/button';
import { ActionStatus } from './action-status';
import type { ActionRun, ActionWorkflow } from '#/services/actions';

interface ActionsSidebarProps {
  runnerCount: number;
  readyRunnerCount: number;
  workflowCount: number;
  firstWorkflow?: ActionWorkflow;
  latestRun?: ActionRun;
  onRunners: () => void;
  isCloud?: boolean;
}

export function ActionsSidebar({
  runnerCount,
  readyRunnerCount,
  workflowCount,
  firstWorkflow,
  latestRun,
  onRunners,
  isCloud = false,
}: ActionsSidebarProps) {
  const steps = [
    {
      title: 'Choose a runner',
      description: isCloud
        ? 'Use a temporary Cloud VM or a connected server.'
        : 'Run workflows on your connected servers.',
      icon: Server,
      done: readyRunnerCount > 0,
      href: runnerCount ? undefined : isCloud ? '/billing' : '/actions/runners/new',
      onClick: runnerCount ? onRunners : undefined,
    },
    {
      title: 'Connect a workflow',
      description: 'Load a repository file or paste your YAML.',
      icon: PlayCircle,
      done: workflowCount > 0,
      href: '/actions/new',
    },
    {
      title: 'Follow your first run',
      description: 'See each job, its logs and results.',
      icon: Play,
      done: !!latestRun,
      href: firstWorkflow ? '/actions' : undefined,
    },
  ];

  const current = steps.findIndex((step) => !step.done);

  return (
    <aside className="min-w-0 space-y-4">
      <section className="rounded-2xl bg-card p-5" aria-labelledby="actions-setup-title">
        <h2 id="actions-setup-title" className="font-medium text-base text-foreground">
          {latestRun ? 'Overview' : 'Get started'}
        </h2>

        {latestRun ? (
          <>
            <dl className="mt-5 space-y-3 text-sm">
              <div className="flex items-center justify-between gap-3">
                <dt className="text-muted-foreground">Workflows</dt>
                <dd className="font-medium text-foreground tabular-nums">{workflowCount}</dd>
              </div>
              <div className="flex items-center justify-between gap-3">
                <dt>
                  <button
                    type="button"
                    onClick={onRunners}
                    className="rounded text-muted-foreground underline-offset-4 hover:text-foreground hover:underline focus-visible:outline-2 focus-visible:outline-ring"
                  >
                    Ready runners
                  </button>
                </dt>
                <dd className="font-medium text-foreground tabular-nums">
                  {readyRunnerCount}{' '}
                  <span className="font-normal text-muted-foreground">/ {runnerCount}</span>
                </dd>
              </div>
            </dl>
            <div className="mt-5 border-border/50 border-t pt-4">
              <p className="text-muted-foreground text-xs">Latest run</p>
              <div className="mt-2 flex items-center justify-between gap-3 rounded-lg py-1">
                <span className="min-w-0 truncate font-medium text-foreground text-sm">
                  Run #{latestRun.number}
                </span>
                <ActionStatus status={latestRun.status} />
              </div>
            </div>
          </>
        ) : (
          <ol className="mt-5 space-y-5">
            {steps.map((step, index) => {
              const StepIcon = step.icon;
              const title = (
                <span className="font-medium text-foreground text-sm">{step.title}</span>
              );

              return (
                <li
                  key={step.title}
                  className="relative flex items-start gap-3"
                  aria-current={index === current ? 'step' : undefined}
                >
                  {index < steps.length - 1 && (
                    <span
                      aria-hidden="true"
                      className="-bottom-4 absolute start-4 top-9 w-px bg-border/70"
                    />
                  )}
                  <span
                    aria-hidden="true"
                    className={`relative flex size-8 shrink-0 items-center justify-center rounded-xl ${
                      step.done
                        ? 'bg-emerald-500/10 text-emerald-500'
                        : index === current
                          ? 'bg-primary text-primary-foreground'
                          : 'bg-muted/60 text-muted-foreground'
                    }`}
                  >
                    {step.done ? <Check className="size-4" /> : <StepIcon className="size-4" />}
                  </span>
                  <div className="min-w-0 pt-0.5">
                    {step.href ? (
                      <Link
                        to={step.href as never}
                        className="rounded underline-offset-4 hover:underline focus-visible:outline-2 focus-visible:outline-ring"
                      >
                        {title}
                      </Link>
                    ) : step.onClick ? (
                      <button
                        type="button"
                        onClick={step.onClick}
                        className="rounded text-start underline-offset-4 hover:underline focus-visible:outline-2 focus-visible:outline-ring"
                      >
                        {title}
                      </button>
                    ) : (
                      title
                    )}
                    <p className="mt-1 text-muted-foreground text-xs leading-relaxed">
                      {step.description}
                    </p>
                  </div>
                </li>
              );
            })}
          </ol>
        )}

        <a
          href="https://codedock.io/docs/guides/actions"
          target="_blank"
          rel="noopener noreferrer"
          className="mt-5 flex items-center gap-2 border-border/50 border-t pt-4 font-medium text-muted-foreground text-xs transition-colors hover:text-foreground focus-visible:outline-2 focus-visible:outline-ring"
        >
          <BookOpen className="size-4" />
          Actions guide
          <ArrowUpRight className="ms-auto size-4" />
        </a>
      </section>

      {isCloud && (
        <section className="rounded-2xl bg-card p-5">
          <div className="flex items-center gap-2">
            <Cloud className="size-4 text-muted-foreground" />
            <h2 className="font-medium text-foreground text-sm">Temporary Cloud VM</h2>
          </div>
          <p className="mt-2.5 text-muted-foreground text-xs leading-relaxed">
            Every job gets a temporary Linux VM. VM usage and transfer use your separate Actions balance.
          </p>
          <Button asChild variant="secondary" className="mt-4 w-full rounded-full bg-muted/60 text-foreground hover:bg-muted text-xs">
            <Link to={'/billing' as never}>
              Budget <ArrowRight className="ms-1 size-3.5" />
            </Link>
          </Button>
        </section>
      )}
    </aside>
  );
}
