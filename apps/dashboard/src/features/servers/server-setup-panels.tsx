import {
  Activity,
  ArrowLeft,
  CheckCircle2,
  ChevronRight,
  CircleDashed,
  Download,
  ListChecks,
  Loader2,
  RotateCw,
  XCircle,
  Zap,
} from 'lucide-react';
import { Button } from '#/components/ui/button';
import { cn } from '#/lib/utils';

export type SetupStep = 'choose' | 'checking' | 'results';
export type SetupMode = 'auto' | 'manual';
export type CheckState = 'waiting' | 'running' | 'passed' | 'failed';

export interface VerifyCheck {
  name: string;
  label: string;
  detail: string;
  state: CheckState;
}

export function SetupHeader({
  step,
  serverName,
  onBack,
}: {
  step: SetupStep;
  serverName: string;
  onBack: () => void;
}) {
  const subtitle =
    step === 'choose'
      ? 'Choose how to verify the new node'
      : step === 'checking'
        ? 'Running connection checks'
        : `Verification results for ${serverName}`;

  return (
    <div className="mb-6 flex items-center gap-3">
      <Button variant="ghost" size="icon" onClick={onBack} aria-label="Back">
        <ArrowLeft className="size-4 text-muted-foreground" />
      </Button>
      <div>
        <h1 className="font-medium text-2xl text-foreground/80 tracking-[-0.2px]">Server setup</h1>
        <p className="mt-0.5 text-muted-foreground/70 text-sm">{subtitle}</p>
      </div>
    </div>
  );
}

export function SetupErrorBanner({ message }: { message: string }) {
  return (
    <div className="mb-4 flex items-start gap-3 rounded-xl border border-destructive/30 bg-destructive/5 p-4">
      <XCircle className="mt-0.5 size-5 shrink-0 text-destructive" />
      <div>
        <p className="font-medium text-foreground text-sm">Connection failed</p>
        <p className="mt-0.5 text-muted-foreground text-xs">{message}</p>
      </div>
    </div>
  );
}

export function ChooseMode({
  onSelect,
  onSkip,
  disabled,
}: {
  onSelect: (mode: SetupMode) => void;
  onSkip: () => void;
  disabled?: boolean;
}) {
  const options = [
    {
      mode: 'auto' as const,
      title: 'Verify automatically',
      desc: 'Run every connection check now and open the server when it passes.',
      icon: Zap,
      tone: 'bg-primary/10 text-primary',
    },
    {
      mode: 'manual' as const,
      title: 'Review step by step',
      desc: 'Inspect each check result before continuing to the server.',
      icon: ListChecks,
      tone: 'bg-warning/10 text-warning',
    },
  ];

  return (
    <div className="space-y-3">
      {options.map((option) => (
        <button
          key={option.mode}
          type="button"
          onClick={() => onSelect(option.mode)}
          disabled={disabled}
          className="group w-full rounded-2xl border border-border/50 bg-card p-6 text-start transition-all hover:border-primary/30 hover:bg-primary/[0.02] disabled:opacity-50"
        >
          <div className="flex items-center gap-4">
            <div
              className={cn(
                'flex h-12 w-12 shrink-0 items-center justify-center rounded-xl',
                option.tone
              )}
            >
              <option.icon className="size-6" />
            </div>
            <div className="flex-1">
              <p className="font-semibold text-[15px] text-foreground">{option.title}</p>
              <p className="mt-0.5 text-muted-foreground text-sm">{option.desc}</p>
            </div>
            <ChevronRight className="size-5 text-muted-foreground opacity-0 transition-opacity group-hover:opacity-100" />
          </div>
        </button>
      ))}

      <button
        type="button"
        onClick={onSkip}
        disabled={disabled}
        className="group w-full rounded-2xl border border-border/50 bg-card p-6 text-start transition-all hover:border-primary/30 hover:bg-primary/[0.02] disabled:opacity-50"
      >
        <div className="flex items-center gap-4">
          <div className="flex h-12 w-12 shrink-0 items-center justify-center rounded-xl bg-muted">
            <Activity className="size-6 text-muted-foreground" />
          </div>
          <div className="flex-1">
            <p className="font-semibold text-[15px] text-foreground">Skip verification</p>
            <p className="mt-0.5 text-muted-foreground text-sm">
              Open the server now and verify connectivity later.
            </p>
          </div>
          <ChevronRight className="size-5 text-muted-foreground opacity-0 transition-opacity group-hover:opacity-100" />
        </div>
      </button>
    </div>
  );
}

export function CheckingState() {
  return (
    <div className="rounded-2xl border border-border/50 bg-card p-8">
      <div className="flex flex-col items-center justify-center py-8">
        <Loader2 className="mb-4 size-8 animate-spin text-primary" />
        <p className="font-medium text-foreground text-sm">Connecting to your server</p>
        <p className="mt-1 text-muted-foreground text-xs">Running verification checks</p>
      </div>
    </div>
  );
}

function CheckRow({ check }: { check: VerifyCheck }) {
  return (
    <div className="flex items-center gap-3 rounded-xl px-3 py-2.5 transition-colors hover:bg-muted/30">
      <div className="shrink-0">
        {check.state === 'running' ? (
          <Loader2 className="size-5 animate-spin text-primary" />
        ) : check.state === 'passed' ? (
          <CheckCircle2 className="size-5 text-success" />
        ) : check.state === 'failed' ? (
          <XCircle className="size-5 text-destructive" />
        ) : (
          <CircleDashed className="size-5 text-muted-foreground/40" />
        )}
      </div>
      <div className="min-w-0 flex-1">
        <p className="font-medium text-foreground text-sm">{check.label}</p>
        <p className="mt-0.5 text-muted-foreground text-xs">{check.detail}</p>
      </div>
      <div
        className={cn(
          'rounded-full px-2.5 py-1 font-medium text-xs',
          check.state === 'passed' && 'bg-success/10 text-success',
          check.state === 'failed' && 'bg-destructive/10 text-destructive',
          (check.state === 'waiting' || check.state === 'running') &&
            'bg-muted text-muted-foreground'
        )}
      >
        {check.state === 'passed'
          ? 'Ready'
          : check.state === 'failed'
            ? 'Failed'
            : check.state === 'running'
              ? 'Checking'
              : 'Waiting'}
      </div>
    </div>
  );
}

export function ResultsPanel({
  serverName,
  host,
  checks,
  onRecheck,
  onDone,
  rechecking,
}: {
  serverName: string;
  host: string;
  checks: VerifyCheck[];
  onRecheck: () => void;
  onDone: () => void;
  rechecking: boolean;
}) {
  const failed = checks.filter((check) => check.state === 'failed').length;
  const ready = failed === 0;

  return (
    <div className="space-y-4">
      <div className="rounded-2xl border border-border/50 bg-card">
        <div className="flex items-center gap-3 border-border/50 border-b px-5 py-4">
          <div className="flex h-9 w-9 items-center justify-center rounded-xl bg-primary/10">
            <Download className="size-[18px] text-primary" />
          </div>
          <div>
            <h2 className="font-semibold text-[15px] text-foreground">{serverName}</h2>
            <p className="font-mono text-muted-foreground text-xs">{host}</p>
          </div>
        </div>
        <div className="space-y-1 p-5">
          {checks.map((check) => (
            <CheckRow key={check.name} check={check} />
          ))}
        </div>
      </div>

      <div className="flex items-center gap-3">
        <Button onClick={onDone} className="flex-1 gap-2 py-2.5">
          <CheckCircle2 className="size-4" />
          {ready ? 'Done, open server' : 'Open server anyway'}
        </Button>
        <Button
          variant="secondary"
          onClick={onRecheck}
          disabled={rechecking}
          className="gap-2 py-2.5"
        >
          <RotateCw className={cn('size-4', rechecking && 'animate-spin')} />
          Recheck
        </Button>
      </div>
    </div>
  );
}
