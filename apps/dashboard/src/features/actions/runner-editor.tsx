import { useNavigate } from '@tanstack/react-router';
import { ArrowLeft, Grid } from 'lucide-react';
import { useId, useState } from 'react';
import { toast } from 'sonner';
import { Button } from '@codedock/ui/button';
import { Input } from '@codedock/ui/input';
import { useCreateRunner, useInspectDestination } from '#/hooks/use-actions';
import { useListServers } from '#/hooks/use-servers';
import { useGetPublicSettings } from '#/hooks/use-settings';
import type { Server as ServerModel } from '#/interfaces/server';
import { RunnerAdvancedFields } from './runner-advanced-fields';
import { RUNNER_PRESETS, RunnerCatalogModal, type RunnerPreset } from './runner-catalog-modal';
import { RunnerServerField } from './runner-server-field';
import { RunnerSummaryCard } from './runner-summary-card';

export type DestinationServer = { id: string; name: string };

interface RunnerEditorProps {
  initialServer?: DestinationServer;
  onBack?: () => void;
  open?: boolean;
  onClose?: () => void;
  onSuccess?: (runner: any) => void;
  defaultCloud?: boolean;
}

const PRESET_DEFAULTS: Record<RunnerPreset, { image: string; labels: string }> = {
  ubuntu: { image: 'catthehacker/ubuntu:act-22.04', labels: 'self-hosted, ubuntu, x64' },
  macos: { image: '', labels: 'self-hosted, macos, arm64' },
  linux: { image: '', labels: 'self-hosted, linux, x64' },
  custom: { image: '', labels: 'self-hosted, custom' },
};

export function RunnerEditor({
  initialServer,
  onBack,
  onClose,
  onSuccess,
}: RunnerEditorProps) {
  const navigate = useNavigate();
  const runnerNameId = useId();
  const { data: publicRes } = useGetPublicSettings();
  const isCloud = publicRes?.data?.cloudMode ?? false;
  const { data: serversData } = useListServers();
  const servers: ServerModel[] = serversData ?? [];
  const { mutateAsync: createRunner, isPending } = useCreateRunner();

  const [serverId, setServerId] = useState<string>(initialServer?.id ?? '');
  const [preset, setPreset] = useState<RunnerPreset>('ubuntu');
  const [catalogOpen, setCatalogOpen] = useState(false);
  const [name, setName] = useState(initialServer?.name ?? '');
  const [enabled, setEnabled] = useState(true);

  const [image, setImage] = useState('catthehacker/ubuntu:act-22.04');
  const [labels, setLabels] = useState('self-hosted, ubuntu, x64');
  const [cpu, setCpu] = useState(1);
  const [memoryMb, setMemoryMb] = useState(2048);
  const [maxParallel, setMaxParallel] = useState(1);
  const [allowDockerSocket, setAllowDockerSocket] = useState(false);

  const {
    data: capabilities,
    refetch: recheckCapabilities,
    isFetching: probing,
  } = useInspectDestination(serverId || undefined);

  const handleSelectPreset = (nextPreset: RunnerPreset) => {
    setPreset(nextPreset);
    if (PRESET_DEFAULTS[nextPreset]) {
      setImage(PRESET_DEFAULTS[nextPreset].image);
      setLabels(PRESET_DEFAULTS[nextPreset].labels);
    }
  };

  const isContainer = preset === 'ubuntu' || preset === 'custom';
  const presetInfo = RUNNER_PRESETS.find((p) => p.id === preset) ?? RUNNER_PRESETS[0];
  const PresetIcon = presetInfo.icon;

  const matchingLabels = capabilities
    ? [
        'self-hosted',
        capabilities?.os ?? 'linux',
        capabilities?.architecture ?? 'x64',
        ...labels
          .split(',')
          .map((l) => l.trim())
          .filter((l) => l && l !== 'self-hosted' && l !== 'linux' && l !== 'x64'),
      ]
    : [];

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim()) {
      toast.error('Please enter a runner name');
      return;
    }

    try {
      const created = await createRunner({
        serverId: serverId || undefined,
        kind: 'cloud',
        name: name.trim(),
        enabled,
        config: {
          mode: isContainer ? 'container' : 'native',
          image: isContainer ? image.trim() : null,
          labels: matchingLabels,
          maxParallel,
          cpu,
          memoryMb,
          allowDockerSocket,
        },
      });
      toast.success('Runner verified and saved');
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
          <Button
            type="button"
            variant="ghost"
            size="sm"
            onClick={onBack ?? (() => navigate({ to: '/actions' as never }))}
            className="-ms-3 mb-2 gap-1.5 rounded-lg text-muted-foreground hover:text-foreground"
          >
            <ArrowLeft className="size-4 rtl:rotate-180" />
            {onBack ? 'Back' : 'Back to Actions'}
          </Button>
          <h1 className="font-medium text-2xl text-foreground tracking-tight">Add runner</h1>
          <p className="mt-1 text-muted-foreground text-sm">
            {isCloud
              ? 'Choose a Cloud space and a ready-to-use job environment.'
              : 'Choose a server and a ready-to-use job environment.'}
          </p>
        </div>
        <Button
          type="button"
          variant="ghost"
          onClick={() => (onClose ? onClose() : navigate({ to: '/actions' as never }))}
          className="rounded-xl text-muted-foreground hover:text-foreground"
        >
          Cancel
        </Button>
      </header>

      <form
        onSubmit={handleSubmit}
        className="grid items-start gap-6 lg:grid-cols-[minmax(0,1fr)_340px]"
      >
        <fieldset
          disabled={isPending}
          className="m-0 min-w-0 space-y-5 rounded-2xl border-0 bg-card p-5"
        >
          <RunnerServerField
            isCloud={isCloud}
            initialServer={initialServer}
            serverId={serverId}
            onServerIdChange={(id, sName) => {
              setServerId(id);
              if (sName && !name) setName(sName);
            }}
            servers={servers}
            probing={probing}
            onRecheckCapabilities={() => recheckCapabilities()}
          />

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

            <div className="flex items-center gap-4 rounded-xl border border-border/40 bg-card p-4">
              <span className="flex size-12 shrink-0 items-center justify-center rounded-xl bg-muted/50 p-2.5 text-foreground">
                <PresetIcon className="size-7" />
              </span>
              <div className="min-w-0">
                <h3 className="font-medium text-foreground text-sm">{presetInfo.title}</h3>
                <p className="mt-1 text-muted-foreground text-xs leading-relaxed">
                  {presetInfo.badge}
                </p>
              </div>
            </div>
            <p className="text-muted-foreground text-xs leading-relaxed">
              {isCloud
                ? 'Choose a Cloud space to see compatible environments.'
                : 'Choose a server to see compatible environments.'}
            </p>
          </div>

          <div className="space-y-1.5">
            <label htmlFor={runnerNameId} className="font-medium text-foreground text-sm">
              Runner name
            </label>
            <Input
              id={runnerNameId}
              value={name}
              onChange={(e) => setName(e.target.value)}
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

        <RunnerSummaryCard
          title={name}
          presetTitle={presetInfo.title}
          isContainer={isContainer}
          cpu={cpu}
          memoryMb={memoryMb}
          maxParallel={maxParallel}
          matchingLabels={matchingLabels}
          enabled={enabled}
          onEnabledChange={setEnabled}
          isPending={isPending}
          canSubmit={Boolean(name.trim())}
          isCloud={isCloud}
        />

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
