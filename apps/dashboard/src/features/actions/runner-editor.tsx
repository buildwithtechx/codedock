import { Link, useNavigate } from '@tanstack/react-router';
import { ArrowLeft, Grid, PlayCircle, Server } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import { Input } from '#/components/ui/input';
import { useListServers } from '#/hooks/use-servers';
import type { Server as ServerModel } from '#/interfaces/server';
import { useCreateRunner } from './use-actions';
import {
  RunnerCatalogModal,
  type RunnerPreset,
  RUNNER_PRESETS,
} from './runner-catalog-modal';
import { RunnerAdvancedFields } from './runner-advanced-fields';
export type DestinationServer = { id: string; name: string };

interface RunnerEditorProps {
  initialServer?: DestinationServer;
  onBack?: () => void;
  open?: boolean;
  onClose?: () => void;
  onSuccess?: (runner: any) => void;
}

export function RunnerEditor({ initialServer, onBack, onClose, onSuccess }: RunnerEditorProps) {
  const navigate = useNavigate();
  const { data: serversData } = useListServers();
  const servers: ServerModel[] = serversData ?? [];
  const { mutateAsync: createRunner, isPending } = useCreateRunner();

  const [serverId, setServerId] = useState<string>(
    initialServer?.id ?? (servers.length === 1 ? servers[0].id : '')
  );

  const [preset, setPreset] = useState<RunnerPreset>('ubuntu');
  const [catalogOpen, setCatalogOpen] = useState(false);
  const [name, setName] = useState(initialServer?.name ?? (servers.length === 1 ? servers[0].name : ''));
  const [enabled, setEnabled] = useState(true);

  // Advanced settings state
  const [image, setImage] = useState('catthehacker/ubuntu:act-22.04');
  const [labels, setLabels] = useState('self-hosted, ubuntu, x64');
  const [cpu, setCpu] = useState(1);
  const [memoryMb, setMemoryMb] = useState(2048);
  const [maxParallel, setMaxParallel] = useState(1);
  const [allowDockerSocket, setAllowDockerSocket] = useState(false);

  const handleSelectPreset = (nextPreset: RunnerPreset) => {
    setPreset(nextPreset);
    const defaults: Record<RunnerPreset, { image: string; labels: string }> = {
      ubuntu: { image: 'catthehacker/ubuntu:act-22.04', labels: 'self-hosted, ubuntu, x64' },
      macos: { image: '', labels: 'self-hosted, macos, arm64' },
      linux: { image: '', labels: 'self-hosted, linux, x64' },
      custom: { image: '', labels: 'self-hosted, custom' },
    };
    if (defaults[nextPreset]) {
      setImage(defaults[nextPreset].image);
      setLabels(defaults[nextPreset].labels);
    }
  };

  const isContainer = preset === 'ubuntu' || preset === 'custom';
  const presetInfo = RUNNER_PRESETS.find((p) => p.id === preset) ?? RUNNER_PRESETS[0];
  const PresetIcon = presetInfo.icon;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!serverId) {
      toast.error('Please select a destination server');
      return;
    }
    if (!name.trim()) {
      toast.error('Please enter a runner name');
      return;
    }

    try {
      const created = await createRunner({
        serverId,
        name: name.trim(),
        enabled,
        config: {
          mode: isContainer ? 'container' : 'native',
          image: isContainer ? image.trim() : null,
          labels: labels.split(',').map((l) => l.trim()).filter(Boolean),
          maxParallel,
          cpu,
          memoryMb,
          allowDockerSocket,
        },
      });
      toast.success('Runner created successfully');
      onSuccess?.(created);
      if (onClose) onClose();
      else navigate({ to: '/actions' as never });
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Failed to create runner');
    }
  };

  return (
    <div className="space-y-6">
      <header className="flex flex-wrap items-center justify-between gap-4">
        <div>
          {onBack ? (
            <Button
              type="button"
              variant="ghost"
              size="sm"
              onClick={onBack}
              className="-ms-3 mb-2 gap-1.5 rounded-lg text-muted-foreground hover:text-foreground"
            >
              <ArrowLeft className="size-4 rtl:rotate-180" />
              Back
            </Button>
          ) : (
            <Button
              asChild
              variant="ghost"
              size="sm"
              className="-ms-3 mb-2 gap-1.5 rounded-lg text-muted-foreground hover:text-foreground"
            >
              <Link to={'/actions' as never}>
                <ArrowLeft className="size-4 rtl:rotate-180" />
                Back to Actions
              </Link>
            </Button>
          )}
          <h1 className="font-medium text-2xl text-foreground tracking-tight">Add runner</h1>
          <p className="mt-1 text-muted-foreground text-sm">
            Choose a server and a ready-to-use job environment.
          </p>
        </div>
        <Button
          type="button"
          variant="ghost"
          onClick={() => (onClose ? onClose() : navigate({ to: '/actions' as never }))}
          className="rounded-xl"
        >
          Cancel
        </Button>
      </header>

      <form onSubmit={handleSubmit} className="grid items-start gap-6 lg:grid-cols-[minmax(0,1fr)_340px]">
        <fieldset disabled={isPending} className="m-0 min-w-0 space-y-5 rounded-2xl border-0 bg-card p-5">
          <div className="space-y-1.5">
            <label className="font-medium text-foreground text-sm">Server</label>
            {initialServer ? (
              <div className="flex items-center gap-2.5 rounded-xl border border-border/60 bg-muted/40 px-3.5 py-2.5 text-sm">
                <Server className="size-4 text-muted-foreground" />
                <span className="font-medium text-foreground">{initialServer.name}</span>
              </div>
            ) : (
              <select
                value={serverId}
                onChange={(e) => {
                  setServerId(e.target.value);
                  const s = servers.find((srv) => srv.id === e.target.value);
                  if (s && !name) setName(s.name);
                }}
                required
                className="h-10 w-full rounded-xl border border-border/60 bg-muted/40 px-3 text-foreground text-sm focus:outline-none focus:ring-2 focus:ring-ring"
              >
                <option value="">Select a server...</option>
                {servers.map((s) => (
                  <option key={s.id} value={s.id}>
                    {s.name} ({s.ipAddress || 'Connected'})
                  </option>
                ))}
              </select>
            )}
            <div className="flex items-center justify-between text-muted-foreground text-xs">
              <span>Uses the existing capacity.</span>
              <Link to={'/servers' as never} className="text-primary hover:underline">+ Connect server</Link>
            </div>
          </div>

          <div className="space-y-3 border-border/50 border-t pt-5">
            <div className="flex flex-wrap items-center justify-between gap-2">
              <h2 className="font-medium text-foreground text-sm">Default environment</h2>
              <Button
                type="button"
                size="sm"
                variant="secondary"
                onClick={() => setCatalogOpen(true)}
                className="gap-1.5 rounded-lg"
              >
                <Grid className="size-3.5" />
                Browse catalog
              </Button>
            </div>

            <div className="flex items-center gap-4 rounded-xl bg-muted/30 p-4">
              <span className="flex size-12 shrink-0 items-center justify-center rounded-xl bg-card text-foreground shadow-xs">
                <PresetIcon className="size-6" />
              </span>
              <div className="min-w-0">
                <h3 className="font-medium text-foreground text-sm">{presetInfo.title}</h3>
                <p className="mt-1 text-muted-foreground text-xs leading-relaxed">{presetInfo.badge}</p>
              </div>
            </div>
            <p className="text-muted-foreground text-xs">Choose a server to see compatible environments.</p>
          </div>

          <div className="space-y-1.5">
            <label htmlFor="runner-name" className="font-medium text-foreground text-sm">Runner name</label>
            <Input
              id="runner-name"
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="e.g. production-runner-1"
              required
              className="h-10 rounded-xl bg-muted/40"
            />
          </div>

          <RunnerAdvancedFields
            image={image}
            onImageChange={setImage}
            labels={labels}
            onLabelsChange={setLabels}
            cpu={cpu}
            onCpuChange={setCpu}
            memoryMb={memoryMb}
            onMemoryMbChange={setMemoryMb}
            maxParallel={maxParallel}
            onMaxParallelChange={setMaxParallel}
            allowDockerSocket={allowDockerSocket}
            onAllowDockerSocketChange={setAllowDockerSocket}
            isContainerMode={isContainer}
          />
        </fieldset>

        <aside className="space-y-5 rounded-2xl bg-card p-5">
          <div className="flex items-center gap-3">
            <span className="flex size-8 shrink-0 items-center justify-center rounded-full bg-blue-500/15 text-blue-500">
              <PlayCircle className="size-4" />
            </span>
            <div className="min-w-0">
              <h2 className="font-medium text-base text-foreground">Add runner</h2>
              <p className="mt-0.5 text-muted-foreground text-xs">Runner configuration</p>
            </div>
          </div>

          <dl className="space-y-3 text-sm">
            <div className="flex items-center justify-between gap-3">
              <dt className="text-muted-foreground">Default environment</dt>
              <dd className="font-medium text-foreground">{presetInfo.title}</dd>
            </div>
            {isContainer && (
              <div className="flex items-center justify-between gap-3">
                <dt className="text-muted-foreground">Per job</dt>
                <dd className="font-medium text-foreground tabular-nums">
                  {cpu} vCPU · {memoryMb / 1024} GiB
                </dd>
              </div>
            )}
            <div className="flex items-center justify-between gap-3">
              <dt className="text-muted-foreground">Concurrent jobs</dt>
              <dd className="font-medium text-foreground tabular-nums">{maxParallel}</dd>
            </div>
          </dl>

          <label className="flex items-center gap-2 pt-2 text-sm">
            <input
              type="checkbox"
              checked={enabled}
              onChange={(e) => setEnabled(e.target.checked)}
              className="size-4 rounded border-border"
            />
            <span className="font-medium text-foreground text-xs">Accept new jobs</span>
          </label>

          <Button type="submit" disabled={isPending || !serverId || !name.trim()} className="w-full rounded-full">
            {isPending ? 'Verifying runner...' : 'Verify and save runner'}
          </Button>

          <p className="text-muted-foreground text-xs leading-relaxed">
            Jobs share this server&apos;s capacity. Only run trusted workflows.
          </p>
        </aside>

        <RunnerCatalogModal
          open={catalogOpen}
          selected={preset}
          onSelect={handleSelectPreset}
          onClose={() => setCatalogOpen(false)}
        />
      </form>
    </div>
  );
}
