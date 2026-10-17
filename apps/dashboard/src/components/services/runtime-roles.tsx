import { Terminal } from 'lucide-react';
import { Badge } from '#/components/ui/badge';
import { Button } from '#/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '#/components/ui/card';
import { type RuntimeRole, type RuntimeRoleContainer, roleLogsHref } from '#/lib/runtime-roles';

export function RuntimeRoleReplicas({
  projectId,
  deploymentId,
  role,
  onNavigate,
}: {
  projectId: string;
  deploymentId: string;
  role: RuntimeRole;
  onNavigate?: (href: string) => void;
}) {
  const getBadgeVariant = (status: RuntimeRoleContainer['status']) => {
    switch (status) {
      case 'running':
        return 'default';
      case 'failed':
      case 'missing':
        return 'destructive';
      case 'restarting':
      case 'starting':
        return 'secondary';
      default:
        return 'outline';
    }
  };

  return (
    <ul className="divide-y divide-border/40">
      {role.containers.map((replica) => {
        const href = roleLogsHref(projectId, deploymentId, role.name, replica.replica);
        return (
          <li
            key={replica.replica}
            className="flex flex-wrap items-center justify-between gap-x-3 gap-y-2 py-2 text-xs"
          >
            <div className="flex items-center gap-2">
              <span className="font-medium text-foreground">Replica #{replica.replica}</span>
              <Badge variant={getBadgeVariant(replica.status)}>{replica.status}</Badge>
              {replica.restartCount !== null && (
                <span className="text-[11px] text-muted-foreground tabular-nums">
                  ({replica.restartCount} restarts)
                </span>
              )}
            </div>
            {replica.containerId && (
              <Button
                size="sm"
                variant="ghost"
                className="h-7 px-2 text-xs"
                onClick={() => onNavigate?.(href)}
              >
                <Terminal className="mr-1 h-3.5 w-3.5" />
                Logs
              </Button>
            )}
          </li>
        );
      })}
    </ul>
  );
}

export function RuntimeRolesCard({
  projectId,
  deploymentId,
  roles,
}: {
  projectId: string;
  deploymentId: string;
  roles: RuntimeRole[];
}) {
  if (!roles || roles.length === 0) return null;

  return (
    <Card className="rounded-xl border border-border/60">
      <CardHeader className="pb-3">
        <CardTitle className="font-medium text-sm">Runtime Roles</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        {roles.map((role) => (
          <div key={role.name} className="rounded-lg border border-border/40 p-3">
            <div className="mb-2 flex items-center justify-between">
              <span className="font-semibold text-foreground text-xs">
                {role.name} ({role.kind})
              </span>
              <span className="text-[11px] text-muted-foreground">
                {role.containers.length} / {role.replicas} active
              </span>
            </div>
            <RuntimeRoleReplicas projectId={projectId} deploymentId={deploymentId} role={role} />
          </div>
        ))}
      </CardContent>
    </Card>
  );
}
