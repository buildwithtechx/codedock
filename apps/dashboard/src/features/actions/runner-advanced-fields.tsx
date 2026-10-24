import { ChevronDown, Settings } from 'lucide-react';
import { useState } from 'react';
import { Input } from '@codedock/ui/input';

interface RunnerAdvancedFieldsProps {
  image: string;
  onImageChange: (image: string) => void;
  labels: string;
  onLabelsChange: (labels: string) => void;
  cpu: number;
  onCpuChange: (cpu: number) => void;
  memoryMb: number;
  onMemoryMbChange: (memoryMb: number) => void;
  maxParallel: number;
  onMaxParallelChange: (maxParallel: number) => void;
  allowDockerSocket: boolean;
  onAllowDockerSocketChange: (allow: boolean) => void;
  isContainerMode?: boolean;
}

export function RunnerAdvancedFields({
  image,
  onImageChange,
  labels,
  onLabelsChange,
  cpu,
  onCpuChange,
  memoryMb,
  onMemoryMbChange,
  maxParallel,
  onMaxParallelChange,
  allowDockerSocket,
  onAllowDockerSocketChange,
  isContainerMode = true,
}: RunnerAdvancedFieldsProps) {
  const [open, setOpen] = useState(false);

  return (
    <details
      open={open}
      onToggle={(event) => setOpen(event.currentTarget.open)}
      className="group border-border/50 border-t pt-4"
    >
      <summary className="flex cursor-pointer list-none items-center justify-between gap-3 font-medium text-sm">
        <span className="flex items-center gap-2 text-foreground">
          <Settings className="size-4 text-muted-foreground" />
          Advanced settings
        </span>
        <ChevronDown className="size-4 text-muted-foreground transition-transform group-open:rotate-180" />
      </summary>

      <div className="mt-5 space-y-5">
        {isContainerMode && (
          <div className="space-y-1.5">
            <label className="font-medium text-foreground text-xs">Runner image</label>
            <Input
              value={image}
              onChange={(e) => onImageChange(e.target.value)}
              placeholder="catthehacker/ubuntu:act-22.04"
              className="h-10 rounded-xl bg-muted/40"
            />
            <p className="text-muted-foreground text-xs leading-relaxed">
              Use an image with Node.js and Git installed.
            </p>
          </div>
        )}

        <div className="space-y-1.5">
          <label className="font-medium text-foreground text-xs">Additional labels</label>
          <Input
            value={labels}
            onChange={(e) => onLabelsChange(e.target.value)}
            placeholder="self-hosted, linux, x64"
            className="h-10 rounded-xl bg-muted/40"
          />
          <p className="text-muted-foreground text-xs leading-relaxed">
            Separate labels with commas. OS and architecture come from the server.
          </p>
        </div>

        <div className={`grid gap-4 ${isContainerMode ? 'sm:grid-cols-3' : 'sm:grid-cols-1'}`}>
          {isContainerMode && (
            <>
              <div className="space-y-1.5">
                <label className="font-medium text-foreground text-xs">vCPU per container</label>
                <Input
                  type="number"
                  min={0.25}
                  step={0.25}
                  max={64}
                  value={cpu}
                  onChange={(e) => onCpuChange(Number(e.target.value))}
                  className="h-10 rounded-xl bg-muted/40"
                />
              </div>
              <div className="space-y-1.5">
                <label className="font-medium text-foreground text-xs">
                  Memory per container (MiB)
                </label>
                <Input
                  type="number"
                  min={256}
                  step={256}
                  value={memoryMb}
                  onChange={(e) => onMemoryMbChange(Number(e.target.value))}
                  className="h-10 rounded-xl bg-muted/40"
                />
              </div>
            </>
          )}

          <div className="space-y-1.5">
            <label className="font-medium text-foreground text-xs">Concurrent jobs</label>
            <Input
              type="number"
              min={1}
              max={16}
              value={maxParallel}
              onChange={(e) => onMaxParallelChange(Number(e.target.value))}
              className="h-10 rounded-xl bg-muted/40"
            />
          </div>
        </div>

        {isContainerMode && (
          <label className="flex items-start gap-2.5 pt-1 text-sm">
            <input
              type="checkbox"
              checked={allowDockerSocket}
              onChange={(e) => onAllowDockerSocketChange(e.target.checked)}
              className="mt-1 size-4 rounded border-border"
            />
            <span>
              <span className="block font-medium text-foreground text-xs">Allow Docker access</span>
              <span className="mt-0.5 block text-muted-foreground text-xs leading-relaxed">
                Enable image builds inside job containers. Only run trusted workflows on this
                server.
              </span>
            </span>
          </label>
        )}
      </div>
    </details>
  );
}
