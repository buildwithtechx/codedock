import { Cpu } from 'lucide-react';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@codedock/ui/dialog';
import { ClusterScalingPanel } from '#/features/servers/clusters/cluster-scaling-panel';

interface TopologyScalingSetupProps {
  clusterId: string;
  isOpen: boolean;
  onClose: () => void;
  onReady?: (ready: boolean) => void;
}

export function TopologyScalingSetup({
  clusterId,
  isOpen,
  onClose,
  onReady,
}: TopologyScalingSetupProps) {
  return (
    <Dialog open={isOpen} onOpenChange={(open) => !open && onClose()}>
      <DialogContent className="max-w-2xl rounded-2xl p-6">
        <DialogHeader icon={Cpu}>
          <DialogTitle className="text-base font-semibold">
            Cluster Scaling Runtime (K3s)
          </DialogTitle>
          <DialogDescription className="text-xs text-muted-foreground">
            Configure and verify K3s operators, cert-manager, and CloudNativePG to enable database
            and service auto-scaling.
          </DialogDescription>
        </DialogHeader>

        <div className="pt-2">
          <ClusterScalingPanel
            clusterId={clusterId}
            onReadyChange={(ready) => {
              onReady?.(ready);
              if (ready) {
                onClose();
              }
            }}
          />
        </div>
      </DialogContent>
    </Dialog>
  );
}
