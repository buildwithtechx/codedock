import { Link, useNavigate } from '@tanstack/react-router';
import {
  ArrowLeft,
  CheckCircle2,
  Cpu,
  Info,
  KeyRound,
  Network,
  Server as ServerIcon,
} from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import { Card } from '#/components/ui/card';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import type { Server, TestSSHRequest } from '#/interfaces/server';
import { cn } from '#/lib/utils';
import { createServer, testServerConnection } from './api';
import {
  CheckingState,
  ChooseMode,
  ResultsPanel,
  SetupErrorBanner,
  SetupHeader,
  type SetupMode,
  type SetupStep,
  type VerifyCheck,
} from './server-setup-panels';
import { ServerSshForm } from './server-ssh-form';

type FormMode = 'remote' | 'local';

const guidance = [
  {
    icon: Network,
    text: 'Any host with SSH access. Codedock connects over port 22 by default.',
  },
  {
    icon: KeyRound,
    text: 'An SSH private key or password for the connection user.',
  },
  {
    icon: ServerIcon,
    text: 'Saved servers are verified, then appear in your fleet.',
  },
];

export function ServerCreationPage() {
  const navigate = useNavigate();
  const [formMode, setFormMode] = useState<FormMode>('remote');
  const [localName, setLocalName] = useState('Local Daemon Node');
  const [creatingLocal, setCreatingLocal] = useState(false);

  const [step, setStep] = useState<SetupStep | null>(null);
  const [mode, setMode] = useState<SetupMode | null>(null);
  const [savedServer, setSavedServer] = useState<Server | null>(null);
  const [credentials, setCredentials] = useState<TestSSHRequest | null>(null);
  const [checks, setChecks] = useState<VerifyCheck[]>([]);
  const [setupError, setSetupError] = useState<string | null>(null);
  const [rechecking, setRechecking] = useState(false);

  const goToServer = (id: string) =>
    navigate({ to: '/servers/$serverId', params: { serverId: id } });

  const handleSaved = (server: Server, creds: TestSSHRequest | null) => {
    setSavedServer(server);
    setCredentials(creds);
    setSetupError(null);
    setChecks([]);
    setMode(null);
    setStep('choose');
  };

  const runChecks = async (selectedMode: SetupMode, stayOnResults: boolean) => {
    if (!savedServer) return;
    setMode(selectedMode);
    setSetupError(null);
    if (!stayOnResults) setStep('checking');
    setRechecking(stayOnResults);

    const endpoint = savedServer.isLocal
      ? 'Local Docker socket'
      : `${credentials?.sshUser ?? 'root'}@${credentials?.sshHost ?? savedServer.ipAddress}:${credentials?.sshPort ?? 22}`;

    if (savedServer.isLocal) {
      const localChecks: VerifyCheck[] = [
        {
          name: 'registered',
          label: 'Local node registered',
          detail: `${savedServer.name} was added to your fleet.`,
          state: 'passed',
        },
        {
          name: 'telemetry',
          label: 'Telemetry stream',
          detail: 'Live metrics start flowing once the node reports its first heartbeat.',
          state: 'passed',
        },
      ];
      setChecks(localChecks);
      setRechecking(false);
      if (selectedMode === 'auto' && !stayOnResults) {
        toast.success('Local node connected');
        goToServer(savedServer.id);
        return;
      }
      setStep('results');
      return;
    }

    setChecks([
      { name: 'connection', label: 'SSH reachable', detail: endpoint, state: 'running' },
      { name: 'auth', label: 'Credentials accepted', detail: endpoint, state: 'waiting' },
      {
        name: 'registered',
        label: 'Server registered',
        detail: `${savedServer.name} was saved to your fleet.`,
        state: 'waiting',
      },
    ]);

    const verdict = credentials
      ? await testServerConnection(credentials)
      : { ok: false, message: 'No credentials captured for verification' };

    const next: VerifyCheck[] = [
      {
        name: 'connection',
        label: 'SSH reachable',
        detail: verdict.ok ? `${endpoint} answered.` : verdict.message,
        state: verdict.ok ? 'passed' : 'failed',
      },
      {
        name: 'auth',
        label: 'Credentials accepted',
        detail: verdict.ok
          ? savedServer.sshAuthMethod === 'password'
            ? 'Password accepted by the host.'
            : 'SSH key accepted by the host.'
          : verdict.message,
        state: verdict.ok ? 'passed' : 'failed',
      },
      {
        name: 'registered',
        label: 'Server registered',
        detail: `${savedServer.name} was saved to your fleet.`,
        state: 'passed',
      },
    ];
    setChecks(next);
    setRechecking(false);

    if (!verdict.ok) {
      setSetupError(verdict.message);
      setStep('results');
      return;
    }
    if (selectedMode === 'auto' && !stayOnResults) {
      toast.success('Server verified and ready');
      goToServer(savedServer.id);
      return;
    }
    setStep('results');
  };

  const handleSetupBack = () => {
    if (step === 'checking') return;
    if (step === 'results') {
      setStep('choose');
      return;
    }
    setStep(null);
    setMode(null);
    setSetupError(null);
  };

  const handleConnectLocal = async () => {
    setCreatingLocal(true);
    try {
      const server = await createServer({
        name: localName.trim() || 'Local Daemon Node',
        isLocal: true,
        ipAddress: '127.0.0.1',
      });
      toast.success('Local Docker node registered');
      handleSaved(server, null);
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Failed to register local node');
    } finally {
      setCreatingLocal(false);
    }
  };

  if (step && savedServer) {
    return (
      <div className="mx-auto w-full max-w-[1180px]">
        <SetupHeader step={step} serverName={savedServer.name} onBack={handleSetupBack} />
        {setupError && <SetupErrorBanner message={setupError} />}
        {step === 'choose' && (
          <ChooseMode
            onSelect={(selected) => void runChecks(selected, false)}
            onSkip={() => goToServer(savedServer.id)}
          />
        )}
        {step === 'checking' && <CheckingState />}
        {step === 'results' && (
          <ResultsPanel
            serverName={savedServer.name}
            host={savedServer.sshHost || savedServer.ipAddress}
            checks={checks}
            onRecheck={() => void runChecks(mode ?? 'manual', true)}
            onDone={() => goToServer(savedServer.id)}
            rechecking={rechecking}
          />
        )}
      </div>
    );
  }

  return (
    <div>
      <div className="mb-6 flex items-center gap-3">
        <Link
          to="/servers"
          className="flex h-8 w-8 items-center justify-center rounded-lg transition-colors hover:bg-muted"
          aria-label="Back to servers"
        >
          <ArrowLeft className="h-4 w-4 text-muted-foreground" />
        </Link>
        <div>
          <h1 className="font-medium text-2xl text-foreground/80 tracking-[-0.2px]">Add server</h1>
          <p className="mt-0.5 text-muted-foreground/70 text-sm">
            Enter the connection details for the new node
          </p>
        </div>
      </div>

      <div className="grid grid-cols-1 items-start gap-6 xl:grid-cols-[minmax(0,1fr)_340px]">
        <div className="min-w-0 xl:col-start-1 xl:row-start-1">
          <div className="flex items-center gap-2">
            <button
              type="button"
              onClick={() => setFormMode('remote')}
              className={cn(
                'flex items-center gap-2 rounded-xl border px-4 py-2.5 font-medium text-xs transition-colors',
                formMode === 'remote'
                  ? 'border-primary bg-primary/10 text-primary'
                  : 'border-border/80 bg-card text-muted-foreground hover:bg-muted/40 hover:text-foreground'
              )}
            >
              <ServerIcon className="h-4 w-4" />
              Remote SSH Server
            </button>
            <button
              type="button"
              onClick={() => setFormMode('local')}
              className={cn(
                'flex items-center gap-2 rounded-xl border px-4 py-2.5 font-medium text-xs transition-colors',
                formMode === 'local'
                  ? 'border-primary bg-primary/10 text-primary'
                  : 'border-border/80 bg-card text-muted-foreground hover:bg-muted/40 hover:text-foreground'
              )}
            >
              <Cpu className="h-4 w-4" />
              Local Docker Node
            </button>
          </div>

          <Card className="mt-6 p-6">
            {formMode === 'remote' ? (
              <ServerSshForm
                onSuccess={handleSaved}
                onCancel={() => navigate({ to: '/servers' })}
              />
            ) : (
              <div className="space-y-5">
                <div className="flex items-start gap-3 rounded-xl border border-primary/20 bg-primary/5 p-4">
                  <CheckCircle2 className="mt-0.5 h-4 w-4 shrink-0 text-primary" />
                  <div className="text-xs">
                    <p className="font-semibold text-foreground">Zero-Config Local Daemon</p>
                    <p className="mt-0.5 text-muted-foreground leading-relaxed">
                      Connects directly to the local Docker socket on this system. Ideal for local
                      development, single-node installations, and edge workstations.
                    </p>
                  </div>
                </div>

                <div className="space-y-1.5">
                  <Label htmlFor="local-node-name" className="text-xs">
                    Node Name
                  </Label>
                  <Input
                    id="local-node-name"
                    value={localName}
                    onChange={(e) => setLocalName(e.target.value)}
                    placeholder="Local Daemon Node"
                  />
                </div>

                <div className="flex items-center justify-between border-border/60 border-t pt-4">
                  <Button
                    type="button"
                    variant="ghost"
                    onClick={() => navigate({ to: '/servers' })}
                    className="text-xs"
                  >
                    Cancel
                  </Button>
                  <Button
                    type="button"
                    onClick={handleConnectLocal}
                    disabled={creatingLocal}
                    className="text-xs"
                  >
                    {creatingLocal ? 'Connecting…' : 'Connect Local Node'}
                  </Button>
                </div>
              </div>
            )}
          </Card>
        </div>

        <aside className="xl:sticky xl:top-6 xl:col-start-2 xl:row-start-1">
          <div className="rounded-2xl border border-border/50 bg-card">
            <div className="flex items-center gap-3 border-border/50 border-b px-5 py-4">
              <div className="flex h-9 w-9 items-center justify-center rounded-xl bg-warning/10">
                <Info className="size-[18px] text-warning" />
              </div>
              <div>
                <h2 className="font-semibold text-[15px] text-foreground">Getting started</h2>
                <p className="text-muted-foreground text-xs">What you need</p>
              </div>
            </div>
            <div className="p-5">
              <ul className="space-y-3 text-muted-foreground text-sm">
                {guidance.map((item) => (
                  <li key={item.text} className="flex items-start gap-2">
                    <item.icon className="mt-0.5 size-4 shrink-0 text-muted-foreground" />
                    <span>{item.text}</span>
                  </li>
                ))}
              </ul>
            </div>
          </div>
        </aside>
      </div>
    </div>
  );
}
