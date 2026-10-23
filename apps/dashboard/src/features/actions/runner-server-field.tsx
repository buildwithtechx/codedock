import { ChevronDown, Plus, RefreshCw, Server as ServerIcon } from 'lucide-react';
import { useState } from 'react';
import { Button } from '@codedock/ui/button';
import { AddServerModal } from '#/features/servers/add-server-modal';
import type { Server as ServerModel } from '#/interfaces/server';
import type { DestinationServer } from './runner-editor';

interface RunnerServerFieldProps {
  isCloud?: boolean;
  initialServer?: DestinationServer;
  serverId: string;
  onServerIdChange: (id: string, name?: string) => void;
  servers: ServerModel[];
  probing: boolean;
  onRecheckCapabilities: () => void;
}

export function RunnerServerField({
  isCloud = false,
  initialServer,
  serverId,
  onServerIdChange,
  servers,
  probing,
  onRecheckCapabilities,
}: RunnerServerFieldProps) {
  const [dropdownOpen, setDropdownOpen] = useState(false);
  const [addModalOpen, setAddModalOpen] = useState(false);

  const selectedServer =
    servers.find((s) => s.id === serverId) ??
    (initialServer ? { id: initialServer.id, name: initialServer.name } : undefined);

  const fieldLabel = isCloud ? 'Cloud space' : 'Server';
  const placeholder = selectedServer
    ? selectedServer.name
    : isCloud
      ? 'Cloud space'
      : 'Select server';
  const helperText = isCloud ? 'Uses the existing plan.' : 'Target execution host';
  const addLabel = isCloud ? 'New Cloud space' : 'Connect server';

  return (
    <div className="space-y-2">
      <label className="font-medium text-foreground text-sm">{fieldLabel}</label>

      <div className="relative">
        <button
          type="button"
          onClick={() => setDropdownOpen((prev) => !prev)}
          className="flex h-11 w-full items-center justify-between rounded-xl bg-muted/40 px-3.5 text-left text-sm transition-colors hover:bg-muted/60"
        >
          <span className={selectedServer ? 'text-foreground' : 'text-muted-foreground'}>
            {placeholder}
          </span>
          <ChevronDown className="size-4 text-muted-foreground" />
        </button>

        {dropdownOpen && (
          <div className="absolute z-20 mt-1 max-h-60 w-full overflow-auto rounded-xl border border-border/80 bg-card p-1 shadow-lg">
            {servers.length === 0 ? (
              <div className="p-3 text-center text-muted-foreground text-xs">
                {isCloud
                  ? 'No Cloud spaces available. Click + New Cloud space below.'
                  : 'No servers connected. Click + Connect server below.'}
              </div>
            ) : (
              servers.map((s) => (
                <button
                  key={s.id}
                  type="button"
                  onClick={() => {
                    onServerIdChange(s.id, s.name);
                    setDropdownOpen(false);
                  }}
                  className={`flex w-full items-center gap-2.5 rounded-lg px-3 py-2 text-left text-sm transition-colors hover:bg-muted/60 ${
                    s.id === serverId
                      ? 'bg-muted/40 font-medium text-foreground'
                      : 'text-muted-foreground'
                  }`}
                >
                  <ServerIcon className="size-3.5 shrink-0" />
                  <span className="truncate">{s.name}</span>
                </button>
              ))
            )}
          </div>
        )}
      </div>

      <div className="flex items-center justify-between pt-0.5 text-xs">
        <span className="text-muted-foreground">{helperText}</span>
        <Button
          type="button"
          variant="outline"
          size="sm"
          onClick={() => setAddModalOpen(true)}
          className="h-7 gap-1 px-2.5 text-xs text-muted-foreground hover:text-foreground"
        >
          <Plus className="size-3" />
          {addLabel}
        </Button>
      </div>

      {serverId && (
        <div className="flex items-center justify-between rounded-xl bg-muted/30 px-3 py-2 text-muted-foreground text-xs">
          <span>
            {probing ? 'Checking destination capabilities...' : 'Linux · x64 · Docker available'}
          </span>
          <button
            type="button"
            onClick={onRecheckCapabilities}
            disabled={probing}
            className="hover:text-foreground"
            aria-label="Probe server capabilities"
          >
            <RefreshCw className={`size-3.5 ${probing ? 'animate-spin' : ''}`} />
          </button>
        </div>
      )}

      <AddServerModal
        open={addModalOpen}
        initialMode={isCloud ? 'managed' : 'connected'}
        onClose={() => setAddModalOpen(false)}
        onCreated={(created) => {
          onServerIdChange(created.id, created.name);
          setAddModalOpen(false);
        }}
      />
    </div>
  );
}
