import { Link, useNavigate } from '@tanstack/react-router';
import { ArrowLeft, CheckCircle2, Cpu, HardDrive, ServerIcon, ShieldCheck } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import { Card } from '#/components/ui/card';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import { ServerSshForm } from '#/features/servers/server-ssh-form';
import { useCreateServer } from '#/hooks/use-servers';
import type { Server } from '#/interfaces/server';

export function ServerCreationPage() {
  const navigate = useNavigate();
  const [mode, setMode] = useState<'remote' | 'local'>('remote');
  const [localName, setLocalName] = useState('Local Daemon Node');
  const { mutateAsync: createServer, isPending: isCreatingLocal } = useCreateServer();

  const handleSuccess = (server: Server) => {
    navigate({ to: '/servers/$serverId', params: { serverId: server.id } });
  };

  const handleConnectLocal = async () => {
    try {
      const server = await createServer({
        name: localName.trim() || 'Local Daemon Node',
        isLocal: true,
        ipAddress: '127.0.0.1',
      });
      toast.success('Local Docker node registered');
      handleSuccess(server);
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Failed to register local node');
    }
  };

  return (
    <div className="grid min-w-0 gap-8 xl:grid-cols-[minmax(0,1fr)_20rem]">
      <main className="min-w-0">
        <Link
          to="/servers"
          className="inline-flex items-center gap-2 text-muted-foreground text-xs transition-colors hover:text-foreground"
        >
          <ArrowLeft className="h-3.5 w-3.5" />
          Servers
        </Link>
        <header className="mt-4">
          <p className="font-medium text-muted-foreground text-xs uppercase tracking-wider">
            Infrastructure
          </p>
          <h1 className="mt-1 font-semibold text-2xl tracking-tight">Add a Server</h1>
          <p className="mt-1 text-muted-foreground text-xs">
            Connect an external host over SSH or connect the local Docker engine node.
          </p>
        </header>

        <div className="mt-6 flex items-center gap-2">
          <button
            type="button"
            onClick={() => setMode('remote')}
            className={`flex items-center gap-2 rounded-xl border px-4 py-2.5 font-medium text-xs transition-colors ${
              mode === 'remote'
                ? 'border-primary bg-primary/10 text-primary'
                : 'border-border/80 bg-card text-muted-foreground hover:bg-muted/40 hover:text-foreground'
            }`}
          >
            <ServerIcon className="h-4 w-4" />
            Remote SSH Server
          </button>
          <button
            type="button"
            onClick={() => setMode('local')}
            className={`flex items-center gap-2 rounded-xl border px-4 py-2.5 font-medium text-xs transition-colors ${
              mode === 'local'
                ? 'border-primary bg-primary/10 text-primary'
                : 'border-border/80 bg-card text-muted-foreground hover:bg-muted/40 hover:text-foreground'
            }`}
          >
            <Cpu className="h-4 w-4" />
            Local Docker Node
          </button>
        </div>

        <Card className="mt-6 p-6">
          {mode === 'remote' ? (
            <ServerSshForm
              onSuccess={handleSuccess}
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
                  disabled={isCreatingLocal}
                  className="text-xs"
                >
                  {isCreatingLocal ? 'Connecting...' : 'Connect Local Node'}
                </Button>
              </div>
            </div>
          )}
        </Card>
      </main>

      <aside className="hidden xl:sticky xl:top-6 xl:block xl:self-start">
        <section className="space-y-4 rounded-2xl border border-border bg-card p-5">
          <div className="flex items-center gap-2">
            <ShieldCheck className="h-4 w-4 text-primary" />
            <h2 className="font-semibold text-sm">Security & Provisioning</h2>
          </div>
          <p className="text-muted-foreground text-xs leading-relaxed">
            All SSH credentials and private keys are encrypted at rest with hardware or AES vault
            encryption. They are only utilized during container lifecycle operations and telemetry.
          </p>

          <div className="space-y-2 border-border/60 border-t pt-3">
            <div className="flex items-center gap-2 text-muted-foreground text-xs">
              <HardDrive className="h-3.5 w-3.5 text-primary" />
              <span>Direct Docker Engine management</span>
            </div>
            <div className="flex items-center gap-2 text-muted-foreground text-xs">
              <ServerIcon className="h-3.5 w-3.5 text-primary" />
              <span>Real-time CPU & memory telemetry</span>
            </div>
          </div>
        </section>
      </aside>
    </div>
  );
}
