import {
  applyNodeChanges,
  Background,
  Controls,
  type EdgeChange,
  type Node,
  type NodeChange,
  ReactFlow,
} from '@xyflow/react';
import '@xyflow/react/dist/style.css';
import { useState } from 'react';
import { Button } from '#/components/ui/button';
import { AppServiceNode } from './app-service-node';
import { CanvasBindingForm } from './canvas-binding-form';
import { CanvasResourcePanel } from './canvas-resource-panel';
import { DatabaseNode } from './database-node';
import type { OperationalCanvas } from './topology-types';
import { useCanvasSync } from './use-canvas-sync';
import { useTopologyEdits } from './use-topology-edits';

const nodeTypes = { appService: AppServiceNode, database: DatabaseNode };

export function EnvironmentCanvas({ envData }: { envData: OperationalCanvas }) {
  const { nodes, setNodes } = useCanvasSync(envData);
  const edit = useTopologyEdits(envData);
  const [selectedId, setSelected] = useState('');
  const selected = nodes.find((node) => node.id === selectedId);
  const edges = [
    ...envData.edges.filter(
      (edge) => edge.kind !== 'dependency' || edge.id.startsWith('stack-dependency:')
    ),
    ...edit.dependencies,
  ].map((edge) => ({
    ...edge,
    deletable: edge.kind === 'dependency' && !edge.id.startsWith('stack-dependency:'),
    style: {
      stroke: edge.kind === 'binding' ? '#a855f7' : edge.kind === 'routing' ? '#22c55e' : '#3b82f6',
    },
    animated: edge.kind === 'routing',
  }));
  const onEdgesChange = (changes: EdgeChange[]) => {
    const removed = new Set(
      changes.filter((change) => change.type === 'remove').map((change) => change.id)
    );
    if (removed.size) edit.remove(removed);
  };
  return (
    <div className="flex h-150 flex-col">
      <div className="flex flex-wrap items-center gap-3 border-b p-3 text-sm">
        <span>Blue: dependencies. Purple: variable bindings. Green: routing.</span>
        <span className="text-muted-foreground">
          Connect a prerequisite to its application. Dragging only changes your local layout.
        </span>
        {edit.draft && (
          <>
            <Button
              size="sm"
              variant="outline"
              disabled={edit.apply.isPending}
              onClick={edit.discard}
            >
              Discard
            </Button>
            <Button size="sm" disabled={edit.apply.isPending} onClick={() => edit.setReview(true)}>
              Review changes
            </Button>
          </>
        )}
      </div>
      <CanvasBindingForm canvas={envData} />
      {edit.error && (
        <p role="alert" className="p-3 text-destructive text-sm">
          {edit.error}
        </p>
      )}
      {edit.review && (
        <div className="space-y-2 border-b p-3">
          <p className="text-sm">
            Save these prerequisites. Deployments wait until each prerequisite is running. Current
            workloads keep running until explicitly redeployed.
          </p>
          <pre className="max-h-32 overflow-auto text-xs">
            {JSON.stringify(
              edit.dependencies.map(({ source, target }) => ({
                prerequisite: source,
                application: target,
              })),
              null,
              2
            )}
          </pre>
          <Button size="sm" disabled={edit.apply.isPending} onClick={() => edit.apply.mutate()}>
            Apply dependencies
          </Button>
        </div>
      )}
      <div className="flex min-h-0 flex-1">
        <div className="min-w-0 flex-1">
          <ReactFlow
            nodes={nodes}
            edges={edges}
            onNodesChange={(changes: NodeChange<Node>[]) =>
              setNodes((previous) => applyNodeChanges(changes, previous))
            }
            onEdgesChange={onEdgesChange}
            onConnect={edit.connect}
            onNodeClick={(_, node) => setSelected(node.id)}
            nodeTypes={nodeTypes}
            fitView
          >
            <Controls />
            <Background gap={12} size={1} />
          </ReactFlow>
        </div>
        {selected && (
          <CanvasResourcePanel
            environmentId={envData.environment.id}
            node={selected}
            onClose={() => setSelected('')}
          />
        )}
      </div>
    </div>
  );
}
