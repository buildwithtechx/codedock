import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import type { CreateAppServiceRequest } from '#/features/services';

export function ApplicationSetupFields({
  draft,
  source,
  update,
  variables,
  setVariables,
}: {
  draft: CreateAppServiceRequest;
  source: 'git' | 'image';
  update: <K extends keyof CreateAppServiceRequest>(
    key: K,
    value: CreateAppServiceRequest[K]
  ) => void;
  variables: string;
  setVariables: (value: string) => void;
}) {
  const textFields =
    source === 'git'
      ? ([
          'name',
          'repositoryUrl',
          'branch',
          'rootDirectory',
          'installCommand',
          'buildCommand',
          'startCommand',
          'dockerfilePath',
          'staticOutput',
          'domain',
          'healthCheckPath',
        ] as const)
      : (['name', 'imageRef', 'startCommand', 'domain', 'healthCheckPath'] as const);
  const labels: Record<string, string> = {
    name: 'Application name',
    repositoryUrl: 'Repository URL',
    imageRef: 'Image reference',
    branch: 'Branch',
    rootDirectory: 'Repository root',
    installCommand: 'Install command',
    buildCommand: 'Build command',
    startCommand: 'Start command',
    dockerfilePath: 'Dockerfile path',
    staticOutput: 'Static output directory',
    domain: 'Domain (blank generates a domain)',
    healthCheckPath: 'Readiness path',
  };
  return (
    <div className="grid gap-4 sm:grid-cols-2">
      {textFields.map((key) => (
        <div key={key} className="space-y-2">
          <Label htmlFor={`setup-${key}`}>{labels[key]}</Label>
          <Input
            id={`setup-${key}`}
            value={draft[key] ?? ''}
            onChange={(event) => update(key, event.target.value)}
          />
        </div>
      ))}
      <div className="space-y-2">
        <Label htmlFor="setup-runtime">Runtime mode</Label>
        <select
          id="setup-runtime"
          className="w-full rounded-md border bg-background p-2"
          value={draft.runtimeMode}
          onChange={(event) =>
            update('runtimeMode', event.target.value as CreateAppServiceRequest['runtimeMode'])
          }
        >
          <option value="web">Web (Server)</option>
          <option value="worker">Worker</option>
          <option value="static">Static site</option>
        </select>
      </div>
      {source === 'git' && (
        <div className="space-y-2">
          <Label htmlFor="setup-engine">Build engine</Label>
          <select
            id="setup-engine"
            className="w-full rounded-md border bg-background p-2"
            value={draft.buildEngine}
            onChange={(event) =>
              update('buildEngine', event.target.value as CreateAppServiceRequest['buildEngine'])
            }
          >
            {['nixpacks', 'dockerfile', 'buildpacks', 'railpack'].map((engine) => (
              <option key={engine} value={engine}>
                {engine}
              </option>
            ))}
          </select>
        </div>
      )}
      <div className="space-y-2">
        <Label htmlFor="setup-port">Internal port</Label>
        <Input
          id="setup-port"
          type="number"
          min={1}
          max={65535}
          value={draft.internalPort}
          onChange={(event) => update('internalPort', Number(event.target.value))}
        />
      </div>
      <div className="space-y-2">
        <Label htmlFor="setup-memory">Memory limit (MB, 0 uses server default)</Label>
        <Input
          id="setup-memory"
          type="number"
          min={0}
          value={draft.memoryLimit ?? 0}
          onChange={(event) => update('memoryLimit', Number(event.target.value))}
        />
      </div>
      <div className="space-y-2">
        <Label htmlFor="setup-cpu">CPU limit (0 uses server default)</Label>
        <Input
          id="setup-cpu"
          type="number"
          min={0}
          step="0.1"
          value={draft.cpuLimit ?? 0}
          onChange={(event) => update('cpuLimit', Number(event.target.value))}
        />
      </div>
      <div className="space-y-2 sm:col-span-2">
        <Label htmlFor="setup-variables">Variables (KEY=value, one per line)</Label>
        <textarea
          id="setup-variables"
          className="min-h-24 w-full rounded-md border bg-background p-2 font-mono text-sm"
          value={variables}
          onChange={(event) => setVariables(event.target.value)}
          spellCheck={false}
        />
      </div>
    </div>
  );
}
