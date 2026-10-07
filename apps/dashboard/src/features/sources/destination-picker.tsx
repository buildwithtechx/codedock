import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import type { Cluster, ClusterNode } from '#/features/servers/cluster-types';
import { BareRuntimeFields } from '#/features/services/bare-runtime-fields';
import type { RuntimeTarget } from '#/features/services/runtime-types';

export type SetupDestination = {
  kind: 'docker' | 'kubernetes' | 'bare';
  clusterId: string;
  imageRepository: string;
  registryId: string;
  bareNode: ClusterNode;
  bareReleaseUrl: string;
  bareSha256: string;
  bareCommand: string[];
};

const emptyBareNode: ClusterNode = { serverId: '', privateIp: '', interface: '', fingerprint: '' };

export const emptyDestination: SetupDestination = {
  kind: 'docker',
  clusterId: '',
  imageRepository: '',
  registryId: '',
  bareNode: emptyBareNode,
  bareReleaseUrl: '',
  bareSha256: '',
  bareCommand: ['./app'],
};

export function toRuntimeTarget(destination: SetupDestination): RuntimeTarget {
  if (destination.kind === 'kubernetes') {
    return {
      kind: 'kubernetes',
      clusterId: destination.clusterId,
      nodeIds: [],
      imageRepository: destination.imageRepository || undefined,
      registryId: destination.registryId || undefined,
      volumes: [],
    };
  }
  if (destination.kind === 'bare') {
    return {
      kind: 'bare',
      nodeIds: [],
      volumes: [],
      bareNode: destination.bareNode,
      bareReleaseUrl: destination.bareReleaseUrl,
      bareSha256: destination.bareSha256,
      bareCommand: destination.bareCommand,
    };
  }
  return { kind: 'docker', nodeIds: [], volumes: [] };
}

export function describeDestination(destination: SetupDestination, serverId: string): string {
  if (destination.kind === 'kubernetes') return 'Kubernetes cluster destination';
  if (destination.kind === 'bare') return 'Native systemd destination';
  return serverId ? 'Project SSH Docker server' : 'Local Docker';
}

export function DestinationPicker({
  destination,
  serverId,
  clusters,
  registries,
  onChange,
}: {
  destination: SetupDestination;
  serverId: string;
  clusters: Cluster[];
  registries: { id: string; registryUrl: string }[];
  onChange: (next: SetupDestination) => void;
}) {
  const update = (changes: Partial<SetupDestination>) => onChange({ ...destination, ...changes });
  return (
    <div className="space-y-3">
      <div className="space-y-2">
        <Label htmlFor="setup-destination">Destination</Label>
        <select
          id="setup-destination"
          className="w-full rounded-md border bg-background p-2"
          value={destination.kind}
          onChange={(event) => update({ kind: event.target.value as SetupDestination['kind'] })}
        >
          <option value="docker">
            {serverId ? 'SSH Docker server for this project' : 'Local Docker'}
          </option>
          <option value="kubernetes">Kubernetes cluster</option>
          <option value="bare">Native systemd service</option>
        </select>
      </div>
      {destination.kind === 'docker' && (
        <p className="text-muted-foreground text-sm">
          {serverId
            ? 'Deploys to the SSH Docker server bound to this project.'
            : 'Deploys to local Docker. Move the project to an SSH server to use remote Docker.'}
        </p>
      )}
      {destination.kind === 'kubernetes' && (
        <>
          <div className="space-y-2">
            <Label htmlFor="setup-cluster">Ready cluster</Label>
            <select
              id="setup-cluster"
              className="w-full rounded-md border bg-background p-2"
              value={destination.clusterId}
              onChange={(event) => update({ clusterId: event.target.value })}
            >
              <option value="">Select a cluster</option>
              {clusters
                .filter((cluster) => cluster.status === 'READY')
                .map((cluster) => (
                  <option key={cluster.id} value={cluster.id}>
                    {cluster.name}
                  </option>
                ))}
            </select>
          </div>
          <div className="space-y-2">
            <Label htmlFor="setup-registry">Image registry</Label>
            <select
              id="setup-registry"
              className="w-full rounded-md border bg-background p-2"
              value={destination.registryId}
              onChange={(event) => update({ registryId: event.target.value })}
            >
              <option value="">Use application registry</option>
              {registries.map((registry) => (
                <option key={registry.id} value={registry.id}>
                  {registry.registryUrl}
                </option>
              ))}
            </select>
          </div>
          <div className="space-y-2">
            <Label htmlFor="setup-image-repo">Image repository for Git builds</Label>
            <Input
              id="setup-image-repo"
              placeholder="registry.example.com/team/application"
              value={destination.imageRepository}
              onChange={(event) => update({ imageRepository: event.target.value })}
            />
          </div>
        </>
      )}
      {destination.kind === 'bare' && (
        <BareRuntimeFields
          target={{
            kind: 'bare',
            nodeIds: [],
            volumes: [],
            bareNode: destination.bareNode,
            bareReleaseUrl: destination.bareReleaseUrl,
            bareSha256: destination.bareSha256,
            bareCommand: destination.bareCommand,
          }}
          update={(changes) =>
            update({
              bareNode: changes.bareNode ?? destination.bareNode,
              bareReleaseUrl: changes.bareReleaseUrl ?? destination.bareReleaseUrl,
              bareSha256: changes.bareSha256 ?? destination.bareSha256,
              bareCommand: changes.bareCommand ?? destination.bareCommand,
            })
          }
        />
      )}
    </div>
  );
}
