import { useState } from 'react';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import { ClusterNodeFields } from '../servers/cluster-node-fields';
import type { RuntimeTarget } from './runtime-types';
export function BareRuntimeFields({
  target,
  update,
}: {
  target: RuntimeTarget;
  update: (change: Partial<RuntimeTarget>) => void;
}) {
  const [command, setCommand] = useState(JSON.stringify(target.bareCommand ?? ['./app']));
  const [error, setError] = useState('');
  const source = !!target.bareRepoUrl;
  return (
    <fieldset className="space-y-3">
      <legend className="font-medium">Native systemd service</legend>
      <p className="text-muted-foreground text-sm">
        Deploy one executable release artifact as an isolated service user. Set one instance and
        clear Docker volumes and cluster placement. Web services use their native port and must bind
        0.0.0.0 for managed HTTP routing; cluster scheduling and database replication are separate
        capabilities. Persistent files belong in CODEDOCK_DATA_DIR.
      </p>
      <Label htmlFor="bare-source-mode">Release source</Label>
      <select
        id="bare-source-mode"
        className="w-full rounded border bg-background p-2"
        value={source ? 'git' : 'artifact'}
        onChange={(event) =>
          update(
            event.target.value === 'git'
              ? { bareRepoUrl: 'https://', bareReleaseUrl: undefined, bareSha256: undefined }
              : { bareRepoUrl: undefined, bareReleaseUrl: '', bareSha256: '' }
          )
        }
      >
        <option value="artifact">HTTPS artifact with SHA256</option>
        <option value="git">Git repository with toolchain build</option>
      </select>
      <ClusterNodeFields
        nodes={[target.bareNode ?? { serverId: '', privateIp: '', interface: '', fingerprint: '' }]}
        locked={0}
        onChange={(nodes) => update({ bareNode: nodes[0] })}
      />
      {source ? (
        <>
          <Label htmlFor="bare-repo">Git repository URL</Label>
          <Input
            id="bare-repo"
            type="url"
            placeholder="https://github.com/team/application.git"
            value={target.bareRepoUrl ?? ''}
            onChange={(event) => update({ bareRepoUrl: event.target.value })}
          />
          <Label htmlFor="bare-branch">Branch</Label>
          <Input
            id="bare-branch"
            placeholder="main"
            value={target.bareBranch ?? ''}
            onChange={(event) => update({ bareBranch: event.target.value })}
          />
          <Label htmlFor="bare-toolchain">Toolchain</Label>
          <select
            id="bare-toolchain"
            className="w-full rounded border bg-background p-2"
            value={target.bareToolchain ?? 'go'}
            onChange={(event) =>
              update({ bareToolchain: event.target.value as RuntimeTarget['bareToolchain'] })
            }
          >
            <option value="go">Go (go build)</option>
            <option value="node">Node (npm install and build)</option>
            <option value="python">Python (pip install)</option>
            <option value="static">Static site (served over HTTP)</option>
          </select>
          <Label htmlFor="bare-install">Install command (optional)</Label>
          <Input
            id="bare-install"
            placeholder="npm ci"
            value={target.bareInstallCommand ?? ''}
            onChange={(event) => update({ bareInstallCommand: event.target.value })}
          />
          <Label htmlFor="bare-build">Build command (optional)</Label>
          <Input
            id="bare-build"
            placeholder="npm run build"
            value={target.bareBuildCommand ?? ''}
            onChange={(event) => update({ bareBuildCommand: event.target.value })}
          />
          <Label htmlFor="bare-output">Build output path (default app)</Label>
          <Input
            id="bare-output"
            placeholder="app"
            value={target.bareOutput ?? ''}
            onChange={(event) => update({ bareOutput: event.target.value })}
          />
        </>
      ) : (
        <>
          <Label htmlFor="bare-artifact">Executable artifact URL</Label>
          <Input
            id="bare-artifact"
            type="url"
            placeholder="https://example.com/releases/application"
            value={target.bareReleaseUrl ?? ''}
            onChange={(event) => update({ bareReleaseUrl: event.target.value })}
          />
          <Label htmlFor="bare-checksum">Artifact SHA256</Label>
          <Input
            id="bare-checksum"
            value={target.bareSha256 ?? ''}
            onChange={(event) => update({ bareSha256: event.target.value })}
          />
        </>
      )}
      <Label htmlFor="bare-command">Command arguments as a JSON array</Label>
      <Input
        id="bare-command"
        value={command}
        onChange={(event) => {
          setCommand(event.target.value);
          try {
            const value: unknown = JSON.parse(event.target.value);
            if (!Array.isArray(value) || value.some((arg) => typeof arg !== 'string'))
              throw new Error('Provide an array of strings starting with ./app');
            update({ bareCommand: value as string[] });
            setError('');
          } catch {
            update({ bareCommand: [] });
            setError('Provide an array of strings starting with ./app');
          }
        }}
      />
      {error && <p role="alert">{error}</p>}
    </fieldset>
  );
}
