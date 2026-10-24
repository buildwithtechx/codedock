import { AlertTriangle, Cpu, HardDrive, Terminal } from 'lucide-react';
import { Link } from '@tanstack/react-router';
import { Button } from '@codedock/ui/button';

interface DeploymentResourceNoticeProps {
  issue?: string | null;
  serviceId?: string;
  onViewLogs?: () => void;
}

const ISSUE_CONFIG: Record<
  string,
  { title: string; description: string; icon: typeof AlertTriangle }
> = {
  'container-oom': {
    title: 'Container ran out of memory',
    description:
      'Docker confirmed an out-of-memory kill. The container exceeded its configured memory limit. Adjust the service memory limit or upgrade server capacity before retrying.',
    icon: Cpu,
  },
  'memory-allocation': {
    title: 'Memory allocation failed',
    description:
      'The process could not allocate system memory (ENOMEM). Review memory limits and available server capacity before choosing a larger instance.',
    icon: Cpu,
  },
  'process-heap': {
    title: 'Application heap limit reached',
    description:
      'The process reached its runtime heap limit. Adjust NODE_OPTIONS="--max-old-space-size=..." in environment variables before upgrading server capacity.',
    icon: Terminal,
  },
  'storage-full': {
    title: 'Storage allocation failed',
    description:
      'The operation ran out of disk space or inodes (ENOSPC). Review disk usage and prune unused Docker artifacts, preserving data volumes.',
    icon: HardDrive,
  },
  'storage-low': {
    title: 'Low build storage',
    description:
      'The Docker filesystem has too little free space or inodes to safely complete the build.',
    icon: HardDrive,
  },
  'process-killed': {
    title: 'Build process was terminated',
    description:
      'The process received SIGKILL (exit 137). Check command timeouts, external termination signals, or host logs.',
    icon: AlertTriangle,
  },
};

export function DeploymentResourceNotice({
  issue,
  serviceId,
  onViewLogs,
}: DeploymentResourceNoticeProps) {
  if (!issue) return null;

  const info = ISSUE_CONFIG[issue] || {
    title: 'Resource failure detected',
    description: `The deployment stopped due to a resource constraint: ${issue}.`,
    icon: AlertTriangle,
  };
  const Icon = info.icon;

  return (
    <section
      aria-label="Deployment resource guidance"
      className="rounded-2xl border border-amber-500/20 bg-amber-500/5 p-5"
    >
      <div className="flex items-start gap-3">
        <Icon className="mt-0.5 size-5 shrink-0 text-amber-500" aria-hidden="true" />
        <div className="min-w-0 flex-1">
          <h2 className="font-semibold text-foreground text-sm">{info.title}</h2>
          <p className="mt-1 text-muted-foreground text-sm leading-relaxed">{info.description}</p>
          <div className="mt-3 flex flex-wrap items-center gap-2">
            {serviceId && (
              <Button asChild variant="outline" size="sm" className="h-8 text-xs">
                <Link to="/services/$serviceId" params={{ serviceId }}>
                  Configure Service
                </Link>
              </Button>
            )}
            {onViewLogs && (
              <Button variant="ghost" size="sm" onClick={onViewLogs} className="h-8 text-xs">
                View Logs
              </Button>
            )}
          </div>
        </div>
      </div>
    </section>
  );
}
