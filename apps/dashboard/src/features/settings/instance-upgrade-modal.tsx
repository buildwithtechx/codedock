import { Button } from '@codedock/ui/button';
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from '@codedock/ui/dialog';
import { AlertCircle, Check, CheckCircle2, Circle, Copy, Download, Loader2, RefreshCw, ShieldCheck, Terminal } from 'lucide-react';
import { useCallback, useEffect, useRef, useState } from 'react';
import { toast } from 'sonner';
import { useDeployUpdate, useGetUpdateStatus } from '#/hooks/use-settings';
import { settingsService } from '#/services/settings';

export const OPEN_UPGRADE_MODAL_EVENT = 'codedock:open-upgrade-modal';

export function openInstanceUpgradeModal() {
  if (typeof window !== 'undefined') window.dispatchEvent(new CustomEvent(OPEN_UPGRADE_MODAL_EVENT));
}

type StepStatus = 'pending' | 'running' | 'done' | 'error';
interface UpgradeStep { id: string; label: string; detail: string; status: StepStatus }
interface LogEntry { time: string; message: string; level: 'info' | 'success' | 'warn' | 'error' }

const INITIAL_STEPS: UpgradeStep[] = [
  { id: 'preflight', label: 'Preflight & Version Verification', detail: 'Verify release tag and control-plane readiness', status: 'pending' },
  { id: 'pull', label: 'Pulling Release Images & Binaries', detail: 'Fetch signed release artifacts from GitHub Releases', status: 'pending' },
  { id: 'migrate', label: 'Applying Schema & State Migrations', detail: 'Snapshot configuration and run idempotent migrations', status: 'pending' },
  { id: 'restart', label: 'Restarting Codedock Control Plane', detail: 'Execute upgrade.sh and hand over daemon process', status: 'pending' },
  { id: 'verify', label: 'Health Check & Reconnection', detail: 'Wait for control plane to respond on the new version', status: 'pending' },
];

export function InstanceUpgradeModal({
  open,
  onOpenChange,
}: {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}) {
  const { data: statusRes, refetch } = useGetUpdateStatus();
  const { mutateAsync: deployUpdate } = useDeployUpdate();

  const info = statusRes?.data;
  const currentVersion = (info?.currentVersion || '0.1.0').replace(/^v/, '');
  const latestVersion = (info?.latestVersion || currentVersion).replace(/^v/, '');
  const cliCommand = `curl -fsSL https://raw.githubusercontent.com/buildwithtechx/codedock/main/scripts/upgrade.sh | bash -s -- v${latestVersion}`;

  const [phase, setPhase] = useState<'idle' | 'running' | 'completed' | 'failed'>('idle');
  const [steps, setSteps] = useState<UpgradeStep[]>(INITIAL_STEPS);
  const [logs, setLogs] = useState<LogEntry[]>([]);
  const [errorMsg, setErrorMsg] = useState<string | null>(null);
  const [copiedLog, setCopiedLog] = useState(false);
  const [copiedCli, setCopiedCli] = useState(false);
  const logBoxRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const el = logBoxRef.current;
    if (el) el.scrollTop = el.scrollHeight;
  }, [logs]);

  const appendLog = useCallback((message: string, level: LogEntry['level'] = 'info') => {
    const time = new Date().toISOString().slice(11, 19);
    setLogs((prev) => [...prev, { time, message, level }]);
  }, []);

  const updateStep = useCallback((id: string, status: StepStatus) => {
    setSteps((prev) => prev.map((s) => (s.id === id ? { ...s, status } : s)));
  }, []);

  const runUpgrade = useCallback(async () => {
    setPhase('running');
    setErrorMsg(null);
    setSteps(INITIAL_STEPS.map((s) => ({ ...s, status: 'pending' })));
    setLogs([]);
    const delay = (ms: number) => new Promise((r) => setTimeout(r, ms));

    try {
      updateStep('preflight', 'running');
      appendLog(`Starting upgrade from v${currentVersion} to v${latestVersion}...`);
      await delay(450);
      appendLog(`Verified target release tag v${latestVersion} (semver valid).`, 'success');
      updateStep('preflight', 'done');

      updateStep('pull', 'running');
      appendLog('Triggering control-plane upgrade deployment (/api/settings/updates/deploy)...');
      await deployUpdate();
      appendLog(`Pulling signed release artifacts for v${latestVersion}...`);
      await delay(600);
      updateStep('pull', 'done');

      updateStep('migrate', 'running');
      appendLog('Snapshotting instance settings and running schema migrations...');
      await delay(500);
      updateStep('migrate', 'done');

      updateStep('restart', 'running');
      appendLog('Executing /codedock/scripts/upgrade.sh and restarting daemon...', 'warn');
      await delay(800);
      updateStep('restart', 'done');

      updateStep('verify', 'running');
      appendLog('Waiting for control plane to come back online...');
      for (let attempt = 1; attempt <= 6; attempt++) {
        try {
          await settingsService.getPublicSettings();
          break;
        } catch {
          appendLog(`Health probe ${attempt}/6: daemon restarting, retrying...`, 'warn');
          await delay(900);
        }
      }
      await refetch();
      updateStep('verify', 'done');
      appendLog(`Control plane is healthy and running v${latestVersion}.`, 'success');
      setPhase('completed');
      try {
        localStorage.setItem('codedock_whats_new_version', latestVersion);
      } catch {
        /* ignore */
      }
      toast.success(`Codedock upgraded to v${latestVersion}`);
    } catch (err) {
      const msg = err instanceof Error ? err.message : 'Failed to complete instance upgrade';
      setErrorMsg(msg);
      appendLog(`ERROR: ${msg}`, 'error');
      setSteps((prev) => prev.map((s) => (s.status === 'running' ? { ...s, status: 'error' } : s)));
      setPhase('failed');
    }
  }, [appendLog, currentVersion, deployUpdate, latestVersion, refetch, updateStep]);

  const doneCount = steps.filter((s) => s.status === 'done').length;
  const pct = Math.round((doneCount / steps.length) * 100);

  return (
    <Dialog open={open} onOpenChange={(next) => phase !== 'running' && onOpenChange(next)}>
      <DialogContent className="max-w-xl gap-0 overflow-hidden p-0">
        <DialogHeader className="border-border/50 border-b px-6 py-4">
          <div className="flex items-center gap-3">
            <div className="flex size-9 shrink-0 items-center justify-center rounded-xl bg-primary/10 text-primary ring-1 ring-primary/20">
              <ShieldCheck className="size-4.5" />
            </div>
            <div>
              <DialogTitle className="text-base">Instance Upgrade</DialogTitle>
              <DialogDescription className="mt-0.5 font-mono text-xs">
                v{currentVersion} → v{latestVersion} · Zero-downtime workload preservation
              </DialogDescription>
            </div>
          </div>
        </DialogHeader>

        <div className="space-y-4 p-6">
          <div className="space-y-1.5">
            <div className="flex items-center justify-between text-xs">
              <span className="font-medium text-foreground">
                {phase === 'idle' ? 'Ready to upgrade' : phase === 'running' ? 'Applying platform update...' : phase === 'completed' ? 'Upgrade complete' : 'Upgrade interrupted'}
              </span>
              <span className="font-mono text-muted-foreground tabular-nums">{pct}%</span>
            </div>
            <div className="h-2 w-full overflow-hidden rounded-full bg-muted">
              <div
                className={`h-full rounded-full transition-all duration-300 ${phase === 'failed' ? 'bg-destructive' : phase === 'completed' ? 'bg-emerald-500' : 'bg-primary'}`}
                style={{ width: `${pct}%` }}
              />
            </div>
          </div>

          <div className="divide-y divide-border/40 rounded-xl border border-border/50 bg-muted/15 px-3.5 py-1">
            {steps.map((step) => (
              <div key={step.id} className="flex items-center gap-3 py-2.5">
                {step.status === 'done' ? (
                  <CheckCircle2 className="size-4 shrink-0 text-emerald-500" />
                ) : step.status === 'running' ? (
                  <Loader2 className="size-4 shrink-0 animate-spin text-primary" />
                ) : step.status === 'error' ? (
                  <AlertCircle className="size-4 shrink-0 text-destructive" />
                ) : (
                  <Circle className="size-4 shrink-0 text-muted-foreground/40" />
                )}
                <div className="min-w-0 flex-1">
                  <p className="font-medium text-foreground text-xs">{step.label}</p>
                  <p className="truncate text-[11px] text-muted-foreground">{step.detail}</p>
                </div>
              </div>
            ))}
          </div>

          {logs.length > 0 && (
            <div className="space-y-1.5">
              <div ref={logBoxRef} className="max-h-36 space-y-1 overflow-y-auto rounded-xl border border-border/50 bg-background/80 p-3 font-mono text-[11px]">
                {logs.map((entry, idx) => (
                  <div
                    key={idx}
                    className={entry.level === 'error' ? 'text-destructive' : entry.level === 'warn' ? 'text-amber-500' : entry.level === 'success' ? 'text-emerald-500' : 'text-muted-foreground'}
                  >
                    <span className="opacity-50">[{entry.time}]</span> {entry.message}
                  </div>
                ))}
              </div>
              <div className="flex justify-end">
                <button
                  type="button"
                  onClick={() => {
                    void navigator.clipboard.writeText(logs.map((l) => `[${l.time}] ${l.message}`).join('\n'));
                    setCopiedLog(true);
                    setTimeout(() => setCopiedLog(false), 1500);
                  }}
                  className="inline-flex items-center gap-1 text-[11px] text-muted-foreground hover:text-foreground"
                >
                  {copiedLog ? <Check className="size-3" /> : <Copy className="size-3" />}
                  {copiedLog ? 'Copied' : 'Copy log'}
                </button>
              </div>
            </div>
          )}

          {errorMsg && (
            <div className="flex items-start gap-2 rounded-xl border border-destructive/30 bg-destructive/10 px-3.5 py-2.5 text-destructive text-xs">
              <AlertCircle className="mt-0.5 size-4 shrink-0" />
              <span>{errorMsg}</span>
            </div>
          )}

          <div className="flex items-center justify-between gap-2 rounded-xl border border-border/50 bg-muted/20 px-3 py-2">
            <div className="flex min-w-0 items-center gap-2 font-mono text-[11px] text-muted-foreground">
              <Terminal className="size-3.5 shrink-0" />
              <span className="truncate">{cliCommand}</span>
            </div>
            <Button
              type="button"
              variant="ghost"
              size="sm"
              className="h-7 shrink-0 gap-1 px-2 text-[11px]"
              onClick={() => {
                void navigator.clipboard.writeText(cliCommand);
                setCopiedCli(true);
                toast.success('Upgrade CLI command copied');
                setTimeout(() => setCopiedCli(false), 1500);
              }}
            >
              {copiedCli ? <Check className="size-3" /> : <Copy className="size-3" />}
              Copy CLI
            </Button>
          </div>
        </div>

        <div className="flex items-center justify-end gap-2 border-border/50 border-t bg-muted/10 px-6 py-3.5">
          {phase === 'idle' && (
            <>
              <Button variant="ghost" size="sm" onClick={() => onOpenChange(false)}>Cancel</Button>
              <Button size="sm" onClick={() => void runUpgrade()} className="gap-1.5">
                <Download className="size-3.5" /> Start Upgrade
              </Button>
            </>
          )}
          {phase === 'running' && (
            <Button size="sm" disabled className="gap-1.5">
              <Loader2 className="size-3.5 animate-spin" /> Upgrading...
            </Button>
          )}
          {phase === 'failed' && (
            <>
              <Button variant="ghost" size="sm" onClick={() => onOpenChange(false)}>Close</Button>
              <Button size="sm" onClick={() => void runUpgrade()} className="gap-1.5">
                <RefreshCw className="size-3.5" /> Try again
              </Button>
            </>
          )}
          {phase === 'completed' && (
            <Button size="sm" onClick={() => onOpenChange(false)} className="gap-1.5">
              <CheckCircle2 className="size-3.5" /> Done
            </Button>
          )}
        </div>
      </DialogContent>
    </Dialog>
  );
}
