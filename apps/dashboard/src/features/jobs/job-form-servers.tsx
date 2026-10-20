import { Key, Plus, Server, ShieldCheck } from 'lucide-react';
import { Button } from '@codedock/ui/button';
import { Checkbox } from '@codedock/ui/checkbox';
import type { Server as ServerModel } from '#/interfaces/server';
import { Field, type KV, KeyValueEditor, SECTION_TONES, Section } from './job-form-sections';

const inputCls =
  'w-full rounded-xl border border-border/60 bg-background px-3 py-2 text-foreground text-sm outline-none transition-colors focus:border-primary/50';

interface JobEnvironmentSectionProps {
  envRows: KV[];
  setEnvRows: (updater: (prev: KV[]) => KV[]) => void;
  secretRows: KV[];
  setSecretRows: (updater: (prev: KV[]) => KV[]) => void;
}

export function JobEnvironmentSection({
  envRows,
  setEnvRows,
  secretRows,
  setSecretRows,
}: JobEnvironmentSectionProps) {
  return (
    <Section title="Environment" icon={<Key className="size-[18px]" />} tone={SECTION_TONES.environment}>
      <div className="space-y-4">
        <KeyValueEditor rows={envRows} setRows={setEnvRows} addLabel="Add variable" />
        <div className="space-y-2 pt-2">
          <span className="block font-medium text-muted-foreground text-sm">Secrets</span>
          <KeyValueEditor rows={secretRows} setRows={setSecretRows} addLabel="Add secret" secret />
        </div>
      </div>
    </Section>
  );
}

interface JobServersSectionProps {
  servers: ServerModel[];
  serverIds: string[];
  onToggleServer: (id: string) => void;
  onOpenAddServer: () => void;
  isCloud?: boolean;
}

export function JobServersSection({
  servers,
  serverIds,
  onToggleServer,
  onOpenAddServer,
  isCloud = false,
}: JobServersSectionProps) {
  return (
    <Section title="Target servers" icon={<Server className="size-[18px]" />} tone={SECTION_TONES.servers}>
      <p className="text-muted-foreground text-sm">
        {isCloud
          ? 'Use an existing server, or add a dedicated server with its own subscription.'
          : 'Select the servers where this job should execute.'}
      </p>
      <div className="space-y-2">
        {servers.length === 0 ? (
          isCloud ? (
            <button
              type="button"
              onClick={() => onToggleServer('production')}
              className="flex w-full items-center gap-2.5 rounded-xl border border-border/60 bg-muted/20 p-3 text-start transition-colors hover:bg-muted/40"
            >
              <Checkbox checked={serverIds.includes('production')} className="pointer-events-none" />
              <Server className="size-4 text-muted-foreground" />
              <span className="font-medium text-foreground text-sm">Production</span>
            </button>
          ) : (
            <p className="text-muted-foreground text-sm">No connected servers found.</p>
          )
        ) : (
          servers.map((srv) => (
            <button
              key={srv.id}
              type="button"
              onClick={() => onToggleServer(srv.id)}
              className={`flex w-full items-center gap-2.5 rounded-xl border p-3 text-start transition-colors ${
                serverIds.includes(srv.id)
                  ? 'border-primary/50 bg-primary/5 text-foreground'
                  : 'border-border/60 bg-background hover:bg-muted/40 text-muted-foreground'
              }`}
            >
              <Checkbox checked={serverIds.includes(srv.id)} className="pointer-events-none" />
              <Server className="size-4 text-muted-foreground" />
              <span className="font-medium text-foreground text-sm">{srv.name}</span>
            </button>
          ))
        )}
      </div>
      <Button
        type="button"
        variant="secondary"
        size="sm"
        onClick={onOpenAddServer}
        className="gap-1.5 rounded-lg"
      >
        <Plus className="size-3.5" /> {isCloud ? 'Add a dedicated server' : 'Add server'}
      </Button>
      <p className="text-muted-foreground text-xs leading-relaxed">
        {isCloud
          ? "Jobs share the selected server's capacity and usage allowance. Its plan must be active to run."
          : "Jobs share the selected server's capacity and usage allowance."}
      </p>
    </Section>
  );
}

interface JobReliabilitySectionProps {
  timeoutSec: string;
  onTimeoutChange: (val: string) => void;
  maxAttempts: string;
  onMaxAttemptsChange: (val: string) => void;
  backoffSec: string;
  onBackoffChange: (val: string) => void;
}

export function JobReliabilitySection({
  timeoutSec,
  onTimeoutChange,
  maxAttempts,
  onMaxAttemptsChange,
  backoffSec,
  onBackoffChange,
}: JobReliabilitySectionProps) {
  return (
    <Section title="Reliability" icon={<ShieldCheck className="size-[18px]" />} tone={SECTION_TONES.reliability}>
      <div className="grid grid-cols-3 gap-3">
        <Field label="Timeout (s)">
          <input
            value={timeoutSec}
            onChange={(e) => onTimeoutChange(e.target.value)}
            inputMode="numeric"
            placeholder="300"
            className={inputCls}
          />
        </Field>
        <Field label="Max attempts">
          <input
            value={maxAttempts}
            onChange={(e) => onMaxAttemptsChange(e.target.value)}
            inputMode="numeric"
            className={inputCls}
          />
        </Field>
        <Field label="Backoff (s)">
          <input
            value={backoffSec}
            onChange={(e) => onBackoffChange(e.target.value)}
            inputMode="numeric"
            className={inputCls}
          />
        </Field>
      </div>
    </Section>
  );
}
