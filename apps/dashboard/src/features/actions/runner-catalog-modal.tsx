import type { ComponentType } from 'react';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@codedock/ui/dialog';
import { CustomImageLogo, LinuxLogo, MacosLogo, UbuntuLogo } from './runner-logos';

export type RunnerPreset = 'ubuntu' | 'macos' | 'linux' | 'custom';

export interface RunnerPresetInfo {
  id: RunnerPreset;
  title: string;
  badge: string;
  description: string;
  runtimeType: string;
  icon: ComponentType<{ className?: string }>;
}

export const RUNNER_PRESETS: RunnerPresetInfo[] = [
  {
    id: 'ubuntu',
    title: 'Ubuntu 22.04',
    badge: 'Ready-to-use Linux image with Node.js and Git.',
    description: 'Choose a server to see compatible environments.',
    runtimeType: 'Container',
    icon: UbuntuLogo,
  },
  {
    id: 'macos',
    title: 'macOS',
    badge: 'Use the tools installed on your connected Mac.',
    description: 'Choose a server to see compatible environments.',
    runtimeType: 'Native',
    icon: MacosLogo,
  },
  {
    id: 'linux',
    title: 'Linux',
    badge: "Run directly with your server's installed tools.",
    description: 'Choose a server to see compatible environments.',
    runtimeType: 'Native',
    icon: LinuxLogo,
  },
  {
    id: 'custom',
    title: 'Custom image',
    badge: 'Bring a Docker runner image with Node.js and Git.',
    description: 'Choose a server to see compatible environments.',
    runtimeType: 'Container',
    icon: CustomImageLogo,
  },
];

interface RunnerCatalogModalProps {
  open: boolean;
  selected: RunnerPreset;
  onSelect: (preset: RunnerPreset) => void;
  onClose: () => void;
}

export function RunnerCatalogModal({ open, selected, onSelect, onClose }: RunnerCatalogModalProps) {
  return (
    <Dialog open={open} onOpenChange={(val) => !val && onClose()}>
      <DialogContent
        className="w-[calc(100vw-2rem)] p-5 sm:max-w-[720px] sm:p-6 outline-none"
        showCloseButton={true}
      >
        <div className="space-y-5">
          <DialogHeader className="text-left">
            <DialogTitle className="font-medium text-foreground text-xl">
              Runner catalog
            </DialogTitle>
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
                  className={`relative flex h-full w-full rounded-xl border p-4 text-start transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring ${
                    isSelected
                      ? 'border-primary bg-primary/5 ring-1 ring-primary/20'
                      : 'border-border/50 bg-card hover:border-primary/30 hover:bg-primary/[0.02]'
                  }`}
                >
                  <div className="flex w-full items-start gap-3">
                    <div
                      className={`shrink-0 rounded-lg p-2.5 ${
                        isSelected
                          ? 'bg-primary/10 text-primary'
                          : 'bg-muted/50 text-muted-foreground'
                      }`}
                    >
                      <Icon className="size-6" />
                    </div>

                    <div className="min-w-0 flex-1">
                      <div className="flex flex-wrap items-center gap-2">
                        <span
                          className={`font-semibold text-sm ${
                            isSelected ? 'text-foreground' : 'text-foreground/80'
                          }`}
                        >
                          {preset.title}
                        </span>
                      </div>

                      <div className="mt-1 space-y-3 text-muted-foreground text-xs leading-relaxed">
                        <span className="block">{preset.badge}</span>
                        <span
                          className={`block font-medium text-xs ${
                            isSelected ? 'text-foreground/80' : 'text-muted-foreground'
                          }`}
                        >
                          {preset.runtimeType}
                        </span>
                      </div>
                    </div>

                    {isSelected && (
                      <span
                        aria-hidden="true"
                        className="mt-0.5 flex size-5 shrink-0 items-center justify-center rounded-full bg-primary"
                      >
                        <span className="size-2 rounded-full bg-primary-foreground" />
                      </span>
                    )}
                  </div>
                </button>
              );
            })}
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
}
