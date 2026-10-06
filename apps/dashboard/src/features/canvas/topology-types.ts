import type { Edge, Node } from '@xyflow/react';
import type { EnvironmentCanvas } from '#/features/projects';

export interface TopologyEdge extends Edge {
  kind: 'dependency' | 'binding' | 'routing';
}

export interface OperationalCanvas extends EnvironmentCanvas {
  revision: string;
  nodes: Node[];
  edges: TopologyEdge[];
}

export function createsCycle(edges: Edge[], source: string, target: string) {
  const visited = new Set<string>();
  const visit = (node: string): boolean => {
    if (node === source) return true;
    if (visited.has(node)) return false;
    visited.add(node);
    return edges.filter((edge) => edge.source === node).some((edge) => visit(edge.target));
  };
  return source === target || visit(target);
}
