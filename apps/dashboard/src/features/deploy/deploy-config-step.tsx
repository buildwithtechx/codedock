import { Boxes, Check, ChevronDown, ChevronUp, Flame, HeartPulse, Sliders } from 'lucide-react';
import { useState } from 'react';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import { Switch } from '#/components/ui/switch';
import { DeployDestinationCard } from './deploy-destination-card';
import { DeployEnvSection } from './deploy-env-section';

const FRAMEWORKS = [
  'Vite',
  'Next.js',
  'Remix',
  'Astro',
  'Node.js',
  'Go',
  'Python',
  'Static HTML',
  'Dockerfile',
];

interface DeployConfigStepProps {
  framework: string;
  onFrameworkChange: (f: string) => void;
  buildEnabled: boolean;
  onBuildEnabledChange: (b: boolean) => void;
  runtimeMode: 'web' | 'worker' | 'static';
  onRuntimeModeChange: (m: 'web' | 'worker' | 'static') => void;
  installCommand: string;
  onInstallCommandChange: (s: string) => void;
  buildCommand: string;
  onBuildCommandChange: (s: string) => void;
  outputDirectory: string;
  onOutputDirectoryChange: (s: string) => void;
  projectName: string;
  onProjectNameChange: (s: string) => void;
  envVars: string;
  onEnvVarsChange: (s: string) => void;
  onEditTarget: () => void;
}

export function DeployConfigStep({
  framework,
  onFrameworkChange,
  buildEnabled,
  onBuildEnabledChange,
  runtimeMode,
  onRuntimeModeChange,
  installCommand,
  onInstallCommandChange,
  buildCommand,
  onBuildCommandChange,
  outputDirectory,
  onOutputDirectoryChange,
  projectName,
  onProjectNameChange,
  envVars,
  onEnvVarsChange,
  onEditTarget,
}: DeployConfigStepProps) {
  const [advancedOpen, setAdvancedOpen] = useState(false);
  const [frameworkSelectOpen, setFrameworkSelectOpen] = useState(false);

  return (
    <div className="space-y-5">
      <DeployDestinationCard runtimeMode={runtimeMode} onEdit={onEditTarget} />

      <section className="rounded-2xl border border-border/60 bg-card p-5">
        <div className="flex items-center justify-between">
          <div>
            <h3 className="font-semibold text-foreground text-sm">Framework</h3>
          </div>
          <button
            type="button"
            onClick={() => setFrameworkSelectOpen((prev) => !prev)}
            className="flex items-center gap-1 font-medium text-muted-foreground text-xs hover:text-foreground"
          >
            Change
            <ChevronDown className="size-3.5" />
          </button>
        </div>

        <div className="mt-3 flex items-center gap-3">
          <div className="grid size-9 place-items-center rounded-xl bg-violet-500/10 text-violet-500">
            <Flame className="size-5" />
          </div>
          <div className="flex items-center gap-2">
            <span className="font-semibold text-base text-foreground">{framework}</span>
            <span className="inline-flex items-center gap-1 rounded-full bg-emerald-500/10 px-2 py-0.5 font-medium text-[11px] text-emerald-500">
              <Check className="size-3" />
              Detected
            </span>
          </div>
        </div>

        {frameworkSelectOpen && (
          <div className="mt-4 grid grid-cols-2 gap-2 border-border/50 border-t pt-3 sm:grid-cols-3">
            {FRAMEWORKS.map((fw) => (
              <button
                key={fw}
                type="button"
                onClick={() => {
                  onFrameworkChange(fw);
                  setFrameworkSelectOpen(false);
                }}
                className={`rounded-lg px-3 py-1.5 text-left font-medium text-xs transition-colors ${
                  framework === fw
                    ? 'bg-primary text-primary-foreground'
                    : 'bg-muted/40 text-muted-foreground hover:bg-muted hover:text-foreground'
                }`}
              >
                {fw}
              </button>
            ))}
          </div>
        )}
      </section>

      <section className="rounded-2xl border border-border/60 bg-card p-6">
        <div className="flex items-center gap-3">
          <div className="grid size-8 place-items-center rounded-lg bg-primary/10 text-primary">
            <Sliders className="size-4" />
          </div>
          <div>
            <h3 className="font-semibold text-foreground text-sm">Deploy Configuration</h3>
            <p className="text-[11px] text-muted-foreground">
              {framework.toLowerCase()} defaults applied
            </p>
          </div>
        </div>

        <div className="mt-5 flex flex-wrap items-center justify-between gap-4 border-border/60 border-b pb-4">
          <div className="flex items-center gap-3">
            <span className="font-medium text-foreground text-xs">Build</span>
            <Switch checked={buildEnabled} onCheckedChange={onBuildEnabledChange} />
          </div>

          <div className="flex items-center gap-2">
            <span className="font-medium text-foreground text-xs">Start</span>
            <div className="flex rounded-lg bg-muted/50 p-1">
              {(['web', 'worker', 'static'] as const).map((mode) => (
                <button
                  key={mode}
                  type="button"
                  onClick={() => onRuntimeModeChange(mode)}
                  className={`rounded-md px-3 py-1 font-medium text-xs capitalize transition-colors ${
                    runtimeMode === mode
                      ? 'bg-card text-foreground shadow-sm'
                      : 'text-muted-foreground hover:text-foreground'
                  }`}
                >
                  {mode === 'web' ? 'Server' : mode}
                </button>
              ))}
            </div>
          </div>
        </div>

        <div className="mt-5 space-y-4">
          <div>
            <Label className="text-xs">Install Command</Label>
            <Input
              value={installCommand}
              onChange={(e) => onInstallCommandChange(e.target.value)}
              className="mt-1 font-mono text-xs"
              placeholder="npm install"
            />
          </div>

          <div>
            <Label className="text-xs">Build Command</Label>
            <Input
              value={buildCommand}
              onChange={(e) => onBuildCommandChange(e.target.value)}
              className="mt-1 font-mono text-xs"
              placeholder="npm run build"
            />
          </div>

          <div>
            <Label className="text-xs">Output Directory</Label>
            <Input
              value={outputDirectory}
              onChange={(e) => onOutputDirectoryChange(e.target.value)}
              className="mt-1 font-mono text-xs"
              placeholder="dist"
            />
          </div>

          {runtimeMode === 'static' && (
            <p className="rounded-lg bg-muted/30 px-3 py-2 text-[11px] text-muted-foreground">
              Served as static files from the edge — nothing runs as a process.
            </p>
          )}

          <div className="border-border/60 border-t pt-3">
            <button
              type="button"
              onClick={() => setAdvancedOpen((prev) => !prev)}
              className="flex items-center gap-1.5 font-medium text-muted-foreground text-xs hover:text-foreground"
            >
              Advanced
              {advancedOpen ? (
                <ChevronUp className="size-3.5" />
              ) : (
                <ChevronDown className="size-3.5" />
              )}
            </button>
            {advancedOpen && (
              <div className="mt-3 space-y-3">
                <div>
                  <Label className="text-xs">Root Directory</Label>
                  <Input placeholder="./" className="mt-1 font-mono text-xs" />
                </div>
                <div>
                  <Label className="text-xs">Dockerfile Path</Label>
                  <Input placeholder="Dockerfile" className="mt-1 font-mono text-xs" />
                </div>
              </div>
            )}
          </div>
        </div>
      </section>

      <DeployEnvSection envVars={envVars} onEnvVarsChange={onEnvVarsChange} />

      <section className="rounded-2xl border border-border/60 bg-card p-5">
        <div className="flex items-center gap-3">
          <div className="grid size-8 place-items-center rounded-lg bg-muted text-muted-foreground">
            <Boxes className="size-4" />
          </div>
          <div>
            <h3 className="font-semibold text-foreground text-sm">Compose file</h3>
            <p className="text-[11px] text-muted-foreground">
              Deploying with Docker Compose? Point to the file.
            </p>
          </div>
        </div>
      </section>

      <section className="rounded-2xl border border-border/60 bg-card p-5">
        <div className="flex items-center gap-3">
          <div className="grid size-8 place-items-center rounded-lg bg-muted text-muted-foreground">
            <HeartPulse className="size-4" />
          </div>
          <div>
            <h3 className="font-semibold text-foreground text-sm">Health checks</h3>
            <p className="text-[11px] text-muted-foreground">
              Off — the deploy finishes as soon as your app is running
            </p>
          </div>
        </div>
      </section>

      <section className="rounded-2xl border border-border/60 bg-card p-5">
        <Label className="font-semibold text-foreground text-sm">Project Name</Label>
        <Input
          value={projectName}
          onChange={(e) => onProjectNameChange(e.target.value)}
          placeholder="my-awesome-project"
          className="mt-2 text-sm"
        />
        <p className="mt-1.5 text-[11px] text-muted-foreground">
          A unique identifier for your deployment
        </p>
      </section>
    </div>
  );
}
