import { Database, Loader2 } from 'lucide-react';
import { useId, useState } from 'react';
import { toast } from 'sonner';
import { Button } from '@codedock/ui/button';
import { Input } from '@codedock/ui/input';
import { Label } from '@codedock/ui/label';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@codedock/ui/select';
import { useUpdateServer } from '#/hooks/use-servers';
import type { Server } from '#/interfaces/server';
import {
  BUILD_CACHE_BUDGET_MAX_GIB,
  BUILD_CACHE_BUDGET_MIN_GIB,
  budgetBytesToGib,
  parseBudgetGib,
} from '#/lib/server-build-cache';

type Mode = 'auto' | 'custom';

interface BuildCacheCardProps {
  server: Server;
  onChanged?: (server: Server) => void;
}

export function BuildCacheCard({ server, onChanged }: BuildCacheCardProps) {
  const modeId = useId();
  const sizeId = useId();
  const stored = server.buildCacheBudgetBytes ?? null;
  const [mode, setMode] = useState<Mode>(stored === null ? 'auto' : 'custom');
  const [size, setSize] = useState(budgetBytesToGib(stored));
  const { mutateAsync: updateServer, isPending: saving } = useUpdateServer();

  if (server.buildCacheBudgetBytes === undefined && !server.capabilities?.buildCacheBudget) {
    return null;
  }

  const parsed = parseBudgetGib(size);
  const unchanged = mode === 'custom' && (parsed === stored || size === budgetBytesToGib(stored));
  const status =
    stored === null ? 'Automatic (20% of disk)' : `Custom (${budgetBytesToGib(stored)} GiB)`;

  async function save(budget: number | null, revert: () => void) {
    try {
      const updated = await updateServer({
        id: server.id,
        payload: { buildCacheBudgetBytes: budget },
      });
      onChanged?.(updated);
      toast.success('Build-cache budget updated');
    } catch (err) {
      revert();
      toast.error(err instanceof Error ? err.message : 'Failed to update build-cache budget');
    }
  }

  function handleModeChange(next: Mode) {
    if (next === mode || saving) return;
    const before = mode;
    setMode(next);
    if (next === 'auto' && stored !== null) {
      void save(null, () => setMode(before));
    }
  }

  return (
    <div className="rounded-2xl border border-border/50 bg-card shadow-xs">
      <div className="flex flex-wrap items-center justify-between gap-3 border-border/50 border-b px-5 py-4">
        <div className="flex min-w-0 items-center gap-3">
          <div className="flex size-9 items-center justify-center rounded-xl bg-primary/10">
            <Database className="size-[18px] text-primary" />
          </div>
          <div>
            <h2 className="font-semibold text-[15px] text-foreground">Build Cache Budget</h2>
            <p className="text-muted-foreground text-xs">
              Managed build-cache storage allocation for this server
            </p>
          </div>
        </div>
        <span className="rounded bg-muted/50 px-2 py-1 text-muted-foreground text-xs">
          {status}
        </span>
      </div>

      <div className="space-y-4 px-5 py-4">
        <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between sm:gap-6">
          <div className="min-w-0 space-y-1">
            <div className="flex flex-wrap items-center gap-x-3 gap-y-1">
              <Label htmlFor={modeId} className="font-medium text-foreground text-sm">
                Allocation Mode
              </Label>
              {saving && (
                <span
                  role="status"
                  className="flex items-center gap-1.5 text-muted-foreground text-xs"
                >
                  <Loader2 className="size-3 motion-safe:animate-spin" />
                  Saving...
                </span>
              )}
            </div>
            <p className="text-muted-foreground text-xs leading-relaxed">
              Automatic scales cache retention with Docker disk capacity. Custom sets a fixed limit.
            </p>
          </div>

          <div className="w-full shrink-0 sm:w-56">
            <Select
              value={mode}
              onValueChange={(val) => handleModeChange(val as Mode)}
              disabled={saving}
            >
              <SelectTrigger id={modeId} aria-label="Allocation Mode">
                <SelectValue placeholder="Select allocation mode" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="auto">Automatic (20% of disk)</SelectItem>
                <SelectItem value="custom">Custom size</SelectItem>
              </SelectContent>
            </Select>
          </div>
        </div>

        {mode === 'custom' && (
          <form
            className="flex flex-wrap items-end gap-3"
            onSubmit={(event) => {
              event.preventDefault();
              if (parsed !== null && !unchanged) {
                void save(parsed, () => setSize(budgetBytesToGib(stored)));
              }
            }}
          >
            <div className="min-w-0 flex-1 space-y-1">
              <Label htmlFor={sizeId} className="font-medium text-foreground text-xs">
                Budget in GiB
              </Label>
              <Input
                id={sizeId}
                inputMode="decimal"
                placeholder="e.g. 24"
                value={size}
                onChange={(event) => setSize(event.target.value)}
                disabled={saving}
                aria-invalid={size !== '' && parsed === null}
                aria-describedby={size !== '' && parsed === null ? `${sizeId}-error` : undefined}
              />
            </div>
            <Button
              type="submit"
              className="h-10"
              disabled={saving || parsed === null || unchanged}
            >
              Save Budget
            </Button>
            {size !== '' && parsed === null && (
              <p id={`${sizeId}-error`} role="alert" className="w-full text-destructive text-xs">
                Budget must be between {BUILD_CACHE_BUDGET_MIN_GIB} GiB and{' '}
                {BUILD_CACHE_BUDGET_MAX_GIB} GiB.
              </p>
            )}
          </form>
        )}
      </div>
    </div>
  );
}
