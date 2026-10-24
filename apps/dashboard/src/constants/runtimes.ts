import type { ProjectRuntimeRoles, RuntimeRole, RuntimeRoleContainer } from '#/interfaces/roles';

export const DEFAULT_ROLE_KINDS = ['web', 'worker', 'cron', 'queue', 'api'] as const;

export function runtimeRoleState(role: RuntimeRole): RuntimeRoleContainer['status'] {
  const states = role.containers.map((container) => container.status);
  if (states.includes('failed')) return 'failed';
  if (states.includes('missing')) return 'missing';
  if (states.includes('restarting')) return 'restarting';
  if (!states.length || states.includes('unknown')) return 'unknown';
  if (states.every((state) => state === 'running')) return 'running';
  if (states.every((state) => state === 'stopped')) return 'stopped';
  if (states.includes('starting')) return 'starting';
  return 'failed';
}

export function roleLogsHref(
  projectId: string,
  deploymentId: string,
  name: string,
  replica: number
): string {
  const query = new URLSearchParams({ role: name, replica: String(replica), deploymentId });
  return `/projects/${encodeURIComponent(projectId)}/logs?${query}`;
}

export function roleLogTargetKey(deploymentId: string, name: string, replica: number): string {
  return `role:${JSON.stringify([deploymentId, name, replica])}`;
}

export function runtimeRoleLogTargets(data: ProjectRuntimeRoles | null) {
  if (!data?.deploymentId) return [];
  return data.roles.flatMap((role) =>
    role.containers.map((c) => ({
      key: roleLogTargetKey(data.deploymentId, role.name, c.replica),
      name: `${role.name} #${c.replica}`,
      role: role.name,
      replica: c.replica,
      status: c.status,
    }))
  );
}
