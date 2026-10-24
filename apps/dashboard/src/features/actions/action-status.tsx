import { AlertCircle, CheckCircle2, Circle, Clock, Loader2, XCircle } from 'lucide-react';
import { Badge } from '@codedock/ui/badge';

export type ActionStatusType = 'queued' | 'running' | 'ready' | 'failed' | 'cancelled';

interface ActionStatusProps {
  status: ActionStatusType | string;
  showIcon?: boolean;
}

export function ActionStatus({ status, showIcon = true }: ActionStatusProps) {
  const getVariant = () => {
    switch (status) {
      case 'ready':
      case 'success':
        return 'success';
      case 'failed':
      case 'failure':
        return 'destructive';
      case 'running':
        return 'default';
      case 'queued':
        return 'secondary';
      case 'cancelled':
        return 'outline';
      default:
        return 'secondary';
    }
  };

  const getIcon = () => {
    switch (status) {
      case 'ready':
      case 'success':
        return <CheckCircle2 className="size-3 text-emerald-500" />;
      case 'failed':
      case 'failure':
        return <XCircle className="size-3 text-red-500" />;
      case 'running':
        return <Loader2 className="size-3 animate-spin text-blue-500" />;
      case 'queued':
        return <Clock className="size-3 text-amber-500" />;
      case 'cancelled':
        return <AlertCircle className="size-3 text-muted-foreground" />;
      default:
        return <Circle className="size-3 text-muted-foreground" />;
    }
  };

  return (
    <Badge
      variant={getVariant()}
      className="inline-flex items-center gap-1.5 px-2 py-0.5 font-medium text-xs uppercase tracking-wider"
    >
      {showIcon && getIcon()}
      <span>{status}</span>
    </Badge>
  );
}
