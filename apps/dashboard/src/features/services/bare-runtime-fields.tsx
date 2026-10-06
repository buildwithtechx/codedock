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
  return (
    <fieldset className="space-y-3">
      <legend className="font-medium">Native systemd service</legend>
      <p className="text-muted-foreground text-sm">
        Deploy one executable release artifact as an isolated service user. Set one instance and
        clear the managed domain and Docker volumes. Web services use their native port; cluster
        scheduling and database replication are separate capabilities. Persistent files belong in
        CODEDOCK_DATA_DIR.
      </p>
      <ClusterNodeFields
        nodes={[target.bareNode ?? { serverId: '', privateIp: '', interface: '', fingerprint: '' }]}
        locked={0}
        onChange={(nodes) => update({ bareNode: nodes[0] })}
      />
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
