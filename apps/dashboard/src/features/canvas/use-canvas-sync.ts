import type { Node } from '@xyflow/react';
import { useEffect, useState } from 'react';
import type { OperationalCanvas } from './topology-types';

export function useCanvasSync(canvas: OperationalCanvas) {
  const [nodes, setNodes] = useState<Node[]>([]);
  useEffect(() => {
    setNodes((previous) =>
      canvas.nodes.map((node) => {
        const existing = previous.find((item) => item.id === node.id);
        return {
          ...node,
          deletable: false,
          position: existing?.position ?? node.position,
          selected: existing?.selected ?? false,
        };
      })
    );
  }, [canvas]);
  return { nodes, setNodes };
}
