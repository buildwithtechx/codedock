import { useState } from 'react';
import { Button } from '#/components/ui/button';
import { apiClient } from '#/lib/api-client';
import { inspectCommandEnvironment } from '#/lib/command-environment';

type Commands = Partial<Record<'installCommand' | 'buildCommand' | 'startCommand', string | null>>;

interface CommandEnvironmentNoticeProps {
  commands: Commands;
  projectId?: string;
  serviceId?: string;
  onMigrated?: (commands: Record<string, string>) => void | Promise<void>;
  migrateDraft?: (commands: Record<string, string>, variables: Record<string, string>) => void;
}

export function CommandEnvironmentNotice({
  commands,
  projectId,
  serviceId,
  onMigrated,
  migrateDraft,
}: CommandEnvironmentNoticeProps) {
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState('');

  const entries = (['installCommand', 'buildCommand', 'startCommand'] as const)
    .map((field) => ({
      field,
      command: commands[field] ?? '',
      parsed: inspectCommandEnvironment(commands[field] ?? ''),
    }))
    .filter((row) => row.parsed.keys.length > 0);

  if (!entries.length) return null;

  const migrate = async () => {
    setBusy(true);
    setError('');
    try {
      const selected = Object.fromEntries(entries.map((row) => [row.field, row.command]));
      if (migrateDraft) {
        const variables: Record<string, string> = {};
        for (const row of entries) {
          for (const [key, value] of Object.entries(row.parsed.variables)) {
            if (Object.hasOwn(variables, key) && variables[key] !== value) {
              throw new Error(
                `Commands use different values for ${key}. Resolve the difference first.`
              );
            }
            variables[key] = value;
          }
        }
        migrateDraft(
          Object.fromEntries(entries.map((row) => [row.field, row.parsed.command])),
          variables
        );
      } else if (projectId) {
        const response = await apiClient.post<{ data: { commands: Record<string, string> } }>(
          `/projects/${encodeURIComponent(projectId)}/env/migrate-commands`,
          { environment: 'production', serviceId, commands: selected }
        );
        const updatedCommands = response?.data?.commands || {};
        await onMigrated?.(updatedCommands);
      }
    } catch (cause) {
      const message = cause instanceof Error ? cause.message : 'Could not move command variables.';
      setError(message);
    } finally {
      setBusy(false);
    }
  };

  const allMigratable = entries.every((row) => row.parsed.migratable);
  const uniqueKeys = [...new Set(entries.flatMap((row) => row.parsed.keys))];

  return (
    <div className="rounded-xl border border-amber-500/30 bg-amber-500/10 p-4 text-sm" role="status">
      <p className="font-medium text-foreground">Manage command variables in Environment</p>
      <p className="mt-1 text-muted-foreground text-xs">
        Embedded variables: {uniqueKeys.join(', ')}. Saved variables are available during builds and at runtime.
      </p>
      {allMigratable && (projectId || migrateDraft) ? (
        <Button
          type="button"
          variant="outline"
          size="sm"
          className="mt-3 text-xs"
          disabled={busy}
          onClick={migrate}
        >
          {busy ? 'Moving variables…' : 'Move to Environment'}
        </Button>
      ) : (
        <p className="mt-2 text-muted-foreground text-xs">
          Review this command manually before moving these variables.
        </p>
      )}
      {error && (
        <p role="alert" className="mt-2 text-destructive text-xs font-medium">
          {error}
        </p>
      )}
    </div>
  );
}
