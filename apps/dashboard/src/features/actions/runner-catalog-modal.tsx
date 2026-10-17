import { Laptop, Server, Terminal } from 'lucide-react';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '#/components/ui/dialog';

export type RunnerPreset = 'ubuntu' | 'macos' | 'linux' | 'custom';

export interface RunnerPresetInfo {
  id: RunnerPreset;
  title: string;
  badge: string;
  description: string;
  icon: typeof Server;
}

export const RUNNER_PRESETS: RunnerPresetInfo[] = [
  {
    id: 'ubuntu',
    title: 'Ubuntu 22.04',
    badge: 'Ready-to-use Linux image with Node.js and Git.',
    description: 'Choose a server to see compatible environments.',
    icon: Server,
  },
  {
    id: 'macos',
    title: 'macOS',
    badge: 'Use the tools installed on your connected Mac.',
    description: 'Choose a server to see compatible environments.',
    icon: Laptop,
  },
  {
    id: 'linux',
    title: 'Linux',
    badge: "Run directly with your server's installed tools.",
    description: 'Choose a server to see compatible environments.',
    icon: Terminal,
  },
  {
    id: 'custom',
    title: 'Custom image',
    badge: 'Bring a Docker runner image with Node.js and Git.',
    description: 'Choose a server to see compatible environments.',
    icon: Server,
  },
];

interface RunnerCatalogModalProps {
  open: boolean;
  selected: RunnerPreset;
  onSelect: (preset: RunnerPreset) => void;
  onClose: () => void;
}

export function RunnerCatalogModal({
  open,
  selected,
  onSelect,
  onClose,
}: RunnerCatalogModalProps) {
  return (
    <Dialog open={open} onOpenChange={(val) => !val && onClose()}>
      <DialogContent className="max-w-2xl p-5 sm:p-6" showCloseButton={true}>
        <div className="space-y-5">
          <DialogHeader className="text-left">
            <DialogTitle className="font-medium text-foreground text-xl">Runner catalog</DialogTitle>
            <DialogDescription className="mt-1 text-muted-foreground text-sm leading-relaxed">
              Choose an environment that works on your selected server.
            </DialogDescription>
          </DialogHeader>

          <div className="grid gap-3 sm:grid-cols-2">
            {RUNNER_PRESETS.map((preset) => {
              const isSelected = selected === preset.id;
              const Icon = preset.icon;

              return (
                <button
                  key={preset.id}
                  type="button"
                  onClick={() => {
                    onSelect(preset.id);
                    onClose();
                  }}
                  className={`relative flex flex-col items-start rounded-2xl border p-4 text-start transition-all hover:bg-muted/40 ${
                    isSelected
                      ? 'border-foreground/80 bg-muted/20 shadow-xs'
                      : 'border-border/60 bg-card'
                  }`}
                >
                  <div className="flex w-full items-start justify-between gap-2">
                    <div className="flex items-center gap-3">
                      <span className="flex size-9 shrink-0 items-center justify-center rounded-xl bg-muted/60 text-foreground">
                        <Icon className="size-4" />
                      </span>
                      <span className="font-medium text-foreground text-sm">{preset.title}</span>
                    </div>
                    <span
                      className={`flex size-4 shrink-0 items-center justify-center rounded-full border-2 ${
                        isSelected ? 'border-primary' : 'border-muted-foreground/30'
                      }`}
                    >
                      {isSelected && <span className="size-2 rounded-full bg-primary" />}
                    </span>
                  </div>
                  <p className="mt-2 text-muted-foreground text-xs leading-relaxed">
                    {preset.badge}
                  </p>
                  <p className="mt-2 text-muted-foreground/70 text-[11px]">
                    {preset.description}
                  </p>
                </button>
              );
            })}
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
}
