import { ArrowRight, Plus, Server, Sparkles } from 'lucide-react';
import { useState, type ReactNode } from 'react';
import { Button } from '#/components/ui/button';
import { AddServerModal } from '#/features/servers/add-server-modal';
import type { Server as ServerModel } from '#/interfaces/server';

export type DestinationServer = { id: string; name: string };

interface RunnerDestinationProps {
  onCancel: () => void;
  selfHosted?: boolean;
  renderForm: (server: DestinationServer | undefined, onBack?: () => void) => ReactNode;
}

export function RunnerDestination({
  onCancel,
  selfHosted = true,
  renderForm,
}: RunnerDestinationProps) {
  const [mode, setMode] = useState<'existing' | 'cloud' | null>(null);
  const [selectedServer, setSelectedServer] = useState<DestinationServer>();
  const [addServerOpen, setAddServerOpen] = useState(false);

  if (!selfHosted || mode === 'cloud') {
    return renderForm(undefined, () => setMode(null));
  }

  if (mode === 'existing') {
    return renderForm(selectedServer, () => {
      setSelectedServer(undefined);
      setMode(null);
    });
  }

  return (
    <div className="space-y-5">
      <header className="flex items-start justify-between gap-4">
        <div>
          <h1 className="font-medium text-2xl text-foreground tracking-tight">Add runner</h1>
          <p className="mt-1 text-muted-foreground text-sm">Choose how to host this runner</p>
        </div>
        <Button variant="ghost" onClick={onCancel} className="rounded-xl">
          Cancel
        </Button>
      </header>

      <div className="grid gap-3 sm:grid-cols-2">
        <button
          type="button"
          onClick={() => setMode('existing')}
          className="flex flex-col items-start rounded-2xl border border-border/60 bg-card p-5 text-start transition-all hover:bg-muted/40"
        >
          <div className="flex size-10 items-center justify-center rounded-xl bg-muted/60 text-foreground">
            <Server className="size-5" />
          </div>
          <h2 className="mt-3 font-medium text-foreground text-base">Existing server</h2>
          <p className="mt-1 text-muted-foreground text-xs leading-relaxed">
            Run inside a container or natively on one of your connected servers.
          </p>
        </button>

        <button
          type="button"
          onClick={() => setAddServerOpen(true)}
          className="flex flex-col items-start rounded-2xl border border-border/60 bg-card p-5 text-start transition-all hover:bg-muted/40"
        >
          <div className="flex size-10 items-center justify-center rounded-xl bg-muted/60 text-foreground">
            <Plus className="size-5" />
          </div>
          <h2 className="mt-3 font-medium text-foreground text-base">Connect dedicated server</h2>
          <p className="mt-1 text-muted-foreground text-xs leading-relaxed">
            Provision or pair a separate server reserved for running Actions.
          </p>
        </button>
      </div>

      <div className="pt-2">
        <Button
          type="button"
          variant="ghost"
          size="sm"
          onClick={() => setMode('cloud')}
          className="gap-1.5 text-muted-foreground text-xs hover:text-foreground"
        >
          <Sparkles className="size-3.5 text-blue-500" />
          Or configure a Managed Cloud space runner
          <ArrowRight className="size-3.5 rtl:rotate-180" />
        </Button>
      </div>

      <AddServerModal
        open={addServerOpen}
        initialMode="connected"
        onClose={() => setAddServerOpen(false)}
        onCreated={(created: ServerModel) => {
          setSelectedServer({ id: created.id, name: created.name || created.ipAddress });
          setMode('existing');
        }}
      />
    </div>
  );
}
