import { Button } from '#/components/ui/button';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '#/components/ui/select';
import { useListServers } from '#/hooks/use-servers';
import type { ClusterNode } from './cluster-types';

export function ClusterNodeFields({
  nodes,
  onChange,
  locked = 0,
}: {
  nodes: ClusterNode[];
  onChange: (nodes: ClusterNode[]) => void;
  locked?: number;
}) {
  const servers = useListServers();
  const update = (index: number, key: keyof ClusterNode, value: string) =>
    onChange(nodes.map((node, i) => (i === index ? { ...node, [key]: value } : node)));
  return (
    <div className="space-y-4">
      {nodes.map((node, index) => (
        <fieldset key={index} disabled={index < locked} className="space-y-2 rounded-md border p-3">
          <legend className="px-1 text-sm">
            {index === 0 ? 'Control plane' : 'Worker'} {index + 1}
          </legend>
          <Label htmlFor={`cluster-server-${index}`}>Registered SSH server</Label>
          <Select value={node.serverId} onValueChange={(value) => update(index, 'serverId', value)}>
            <SelectTrigger id={`cluster-server-${index}`}>
              <SelectValue placeholder="Select server" />
            </SelectTrigger>
            <SelectContent>
              {servers.data
                ?.filter((server) => !server.isLocal)
                .map((server) => (
                  <SelectItem key={server.id} value={server.id}>
                    {server.name}
                  </SelectItem>
                ))}
            </SelectContent>
          </Select>
          <Label htmlFor={`cluster-ip-${index}`}>Private IPv4 address</Label>
          <Input
            id={`cluster-ip-${index}`}
            value={node.privateIp}
            placeholder="192.168.1.10"
            onChange={(event) => update(index, 'privateIp', event.target.value)}
          />
          <Label htmlFor={`cluster-interface-${index}`}>Private network interface</Label>
          <Input
            id={`cluster-interface-${index}`}
            value={node.interface}
            placeholder="ens3"
            onChange={(event) => update(index, 'interface', event.target.value)}
          />
          <Label htmlFor={`cluster-key-${index}`}>Verified SSH host fingerprint</Label>
          <Input
            id={`cluster-key-${index}`}
            value={node.fingerprint}
            placeholder="SHA256:…"
            onChange={(event) => update(index, 'fingerprint', event.target.value)}
          />
          {index >= locked && nodes.length > 1 && (
            <Button
              type="button"
              variant="outline"
              onClick={() => onChange(nodes.filter((_, i) => i !== index))}
            >
              Remove node
            </Button>
          )}
        </fieldset>
      ))}
      {servers.error && <p role="alert">{servers.error.message}</p>}
      <Button
        type="button"
        variant="outline"
        disabled={nodes.length >= 32}
        onClick={() =>
          onChange([...nodes, { serverId: '', privateIp: '', interface: '', fingerprint: '' }])
        }
      >
        Add worker
      </Button>
    </div>
  );
}
