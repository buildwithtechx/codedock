import { Boxes, Loader2, Plus, Server, Trash2 } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Button } from '@codedock/ui/button';
import { Input } from '@codedock/ui/input';
import { Label } from '@codedock/ui/label';
import type { ClusterNode, ComputeCluster } from '#/interfaces/clusters';
import { ClusterScalingPanel } from './cluster-scaling-panel';

interface ClusterEditorProps {
  initialCluster?: ComputeCluster | null;
  onSaved?: (cluster: ComputeCluster) => void;
}

export function ClusterEditor({ initialCluster, onSaved }: ClusterEditorProps) {
  const [cluster, setCluster] = useState<ComputeCluster | null>(initialCluster ?? null);
  const [name, setName] = useState(initialCluster?.name ?? 'production-pool');
  const [controls, setControls] = useState<number>(initialCluster?.controls ?? 1);
  const [nodes, setNodes] = useState<ClusterNode[]>(
    initialCluster?.nodes ?? [
      {
        serverId: 'srv-1',
        serverName: 'Primary Host',
        role: 'server',
        privateIp: '10.42.0.1',
      },
      {
        serverId: 'srv-2',
        serverName: 'Worker Node A',
        role: 'agent',
        privateIp: '10.42.0.2',
      },
      {
        serverId: 'srv-3',
        serverName: 'Worker Node B',
        role: 'agent',
        privateIp: '10.42.0.3',
      },
    ]
  );
  const [isSubmitting, setIsSubmitting] = useState(false);

  const addNode = () => {
    const idx = nodes.length + 1;
    setNodes([
      ...nodes,
      {
        serverId: `srv-${idx}`,
        serverName: `Worker Node ${idx}`,
        role: 'agent',
        privateIp: `10.42.0.${idx}`,
      },
    ]);
  };

  const removeNode = (index: number) => {
    if (nodes.length <= 1) {
      toast.error('A cluster requires at least 1 node');
      return;
    }
    setNodes(nodes.filter((_, i) => i !== index));
  };

  const toggleRole = (index: number) => {
    const updated = [...nodes];
    updated[index].role = updated[index].role === 'server' ? 'agent' : 'server';
    setNodes(updated);
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim()) {
      toast.error('Cluster name is required');
      return;
    }

    setIsSubmitting(true);
    try {
      const serverNodes = nodes.filter((n) => n.role === 'server').length;
      const savedCluster: ComputeCluster = {
        id: cluster?.id ?? `cls-${Date.now()}`,
        name: name.trim(),
        version: 'v1.31.1+k3s1',
        controls: serverNodes || controls,
        nodes,
        revision: (cluster?.revision ?? 0) + 1,
        status: 'ready',
        updatedAt: new Date().toISOString(),
      };

      setCluster(savedCluster);
      toast.success('Cluster configuration saved');
      onSaved?.(savedCluster);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="space-y-6" data-testid="cluster-editor">
      {cluster && <ClusterScalingPanel clusterId={cluster.id} />}

      <form
        onSubmit={handleSubmit}
        className="space-y-6 rounded-2xl border border-border/50 bg-card p-5 shadow-xs"
      >
        <div className="flex items-center justify-between border-border/50 border-b pb-4">
          <div className="flex items-center gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-primary/10">
              <Boxes className="size-5 text-primary" />
            </div>
            <div>
              <h3 className="font-semibold text-base text-foreground">K3s Server Cluster Pool</h3>
              <p className="text-muted-foreground text-xs">
                Manage control-plane consensus and scalable worker agents
              </p>
            </div>
          </div>

          <div className="flex items-center gap-2">
            <span className="font-mono text-muted-foreground text-xs">
              Control nodes: {nodes.filter((n) => n.role === 'server').length}
            </span>
          </div>
        </div>

        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
          <div className="space-y-2">
            <Label htmlFor="cluster-name">Cluster Name</Label>
            <Input
              id="cluster-name"
              value={name}
              onChange={(e) => setName(e.target.value)}
              disabled={isSubmitting}
              required
            />
          </div>

          <div className="space-y-2">
            <Label>Control Plane Quorum</Label>
            <div className="grid grid-cols-2 gap-2">
              <Button
                type="button"
                variant={controls === 1 ? 'default' : 'outline'}
                size="sm"
                onClick={() => setControls(1)}
              >
                1 Node (Standalone)
              </Button>
              <Button
                type="button"
                variant={controls === 3 ? 'default' : 'outline'}
                size="sm"
                onClick={() => setControls(3)}
              >
                3 Nodes (HA Quorum)
              </Button>
            </div>
          </div>
        </div>

        <div className="space-y-3 pt-2">
          <div className="flex items-center justify-between">
            <Label className="font-semibold text-foreground text-sm">
              Assigned Server Nodes ({nodes.length})
            </Label>
            <Button type="button" variant="outline" size="sm" onClick={addNode} className="gap-1.5">
              <Plus className="size-3.5" />
              Add Node
            </Button>
          </div>

          <div className="divide-y divide-border/50 rounded-xl border border-border/60 bg-muted/20">
            {nodes.map((node, idx) => (
              <div
                key={node.serverId}
                className="flex flex-wrap items-center justify-between gap-3 p-3"
              >
                <div className="flex items-center gap-3">
                  <div className="flex h-7 w-7 items-center justify-center rounded-lg bg-muted text-muted-foreground">
                    <Server className="size-3.5" />
                  </div>
                  <div>
                    <p className="font-medium text-foreground text-xs">
                      {node.serverName || node.serverId}
                    </p>
                    <p className="font-mono text-[11px] text-muted-foreground">{node.privateIp}</p>
                  </div>
                </div>

                <div className="flex items-center gap-2">
                  <Button
                    type="button"
                    variant={node.role === 'server' ? 'default' : 'secondary'}
                    size="sm"
                    className="h-7 text-xs"
                    onClick={() => toggleRole(idx)}
                  >
                    {node.role === 'server' ? 'Control Plane' : 'Worker Agent'}
                  </Button>

                  <Button
                    type="button"
                    variant="ghost"
                    size="icon"
                    className="h-8 w-8 text-destructive hover:bg-destructive/10"
                    onClick={() => removeNode(idx)}
                    disabled={nodes.length <= 1}
                  >
                    <Trash2 className="size-3.5" />
                  </Button>
                </div>
              </div>
            ))}
          </div>
        </div>

        <div className="flex justify-end border-border/50 border-t pt-4">
          <Button type="submit" disabled={isSubmitting}>
            {isSubmitting && <Loader2 className="mr-1.5 size-4 animate-spin" />}
            Save Cluster Pool
          </Button>
        </div>
      </form>
    </div>
  );
}
