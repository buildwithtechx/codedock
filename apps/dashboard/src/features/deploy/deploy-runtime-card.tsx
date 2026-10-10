import { Check, ChevronDown, Container, Flame } from 'lucide-react';
import { useState } from 'react';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';

const FRAMEWORKS = [
  'Dockerfile',
  'Vite',
  'Next.js',
  'Remix',
  'Astro',
  'Node.js',
  'Go',
  'Python',
  'Static HTML',
];

interface DeployRuntimeCardProps {
  framework: string;
  onFrameworkChange: (f: string) => void;
  exposedPort: number;
  onExposedPortChange: (p: number) => void;
  installCommand: string;
  onInstallCommandChange: (s: string) => void;
  buildCommand: string;
  onBuildCommandChange: (s: string) => void;
  outputDirectory: string;
  onOutputDirectoryChange: (s: string) => void;
}

export function DeployRuntimeCard({
  framework,
  onFrameworkChange,
  exposedPort,
  onExposedPortChange,
  installCommand,
  onInstallCommandChange,
  buildCommand,
  onBuildCommandChange,
  outputDirectory,
  onOutputDirectoryChange,
}: DeployRuntimeCardProps) {
  const [frameworkSelectOpen, setFrameworkSelectOpen] = useState(false);
  const isDocker = framework.toLowerCase() === 'dockerfile';

  return (
    <section className="rounded-2xl border border-border/60 bg-card p-6">
      <div className="flex items-start justify-between gap-4">
        <div className="flex items-start gap-3.5">
          <div className="flex size-10 shrink-0 items-center justify-center rounded-xl bg-primary/10 text-primary">
            {isDocker ? <Container className="size-5" /> : <Flame className="size-5" />}
          </div>
          <div>
            <div className="flex items-center gap-2">
              <h2 className="font-semibold text-foreground text-sm">
                {isDocker ? 'Dockerfile Detected' : `${framework} Detected`}
              </h2>
              <span className="inline-flex items-center gap-1 rounded-full bg-emerald-500/10 px-2 py-0.5 font-medium text-[10px] text-emerald-600 dark:text-emerald-400">
                <Check className="size-2.5" />
                Detected
              </span>
            </div>
            <p className="mt-0.5 text-muted-foreground text-xs">
              {isDocker
                ? 'Your project will be built and deployed using its Dockerfile'
                : `Detected from project configuration — ${framework.toLowerCase()} defaults applied`}
            </p>
          </div>
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

      <p className="mt-4 text-muted-foreground text-xs leading-relaxed">
        {isDocker
          ? 'No build settings needed - your Dockerfile defines the build process. Just configure the port your application listens on.'
          : 'Standard build settings applied. You can adjust build commands and listening port below.'}
      </p>

      <div className="mt-5 space-y-2">
        <Label
          htmlFor="exposed-port"
          className="flex items-center gap-1.5 font-medium text-foreground text-xs"
        >
          <span className="font-mono text-muted-foreground">#</span>
          Exposed Port
        </Label>
        <Input
          id="exposed-port"
          type="number"
          value={exposedPort}
          onChange={(e) => onExposedPortChange(Number(e.target.value) || 8080)}
          className="h-10 w-full rounded-xl border border-border/60 bg-background/50 px-3 font-mono text-foreground text-sm focus-visible:ring-1"
        />
        <p className="text-[11px] text-muted-foreground">
          The port your container exposes (from EXPOSE or CMD)
        </p>
      </div>

      {!isDocker && (
        <div className="mt-5 space-y-4 border-border/50 border-t pt-4">
          <div className="grid gap-3 sm:grid-cols-2">
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
        </div>
      )}
    </section>
  );
}
