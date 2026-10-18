import type { OrganizationDeployment } from '#/features/services';
import type { DeploymentHistorySummary } from './deployment-overview';

export function computeDeploymentSummary(
  deployments: OrganizationDeployment[],
  days = 30,
  now = new Date()
): DeploymentHistorySummary {
  const emptyCounts = () => ({
    total: 0,
    success: 0,
    failed: 0,
    active: 0,
    canceled: 0,
  });

  const daily = Array.from({ length: days }, (_, i) => {
    const d = new Date(now);
    d.setUTCDate(d.getUTCDate() - days + 1 + i);
    return {
      date: d.toISOString().slice(0, 10),
      ...emptyCounts(),
    };
  });

  const projectMap = new Map<
    string,
    {
      projectId: string;
      name: string;
      total: number;
      success: number;
      failed: number;
      active: number;
      canceled: number;
    }
  >();

  const summary: DeploymentHistorySummary = {
    ...emptyCounts(),
    avgBuildDurationMs: null,
    daily,
    projects: [],
  };

  const classify = (statusStr: string): 'success' | 'failed' | 'active' | 'canceled' => {
    const s = (statusStr || '').toLowerCase();
    if (['ready', 'active', 'success', 'deployed'].includes(s)) return 'success';
    if (['failed', 'error', 'action_required', 'rejected'].includes(s)) return 'failed';
    if (['building', 'cloning', 'pulling', 'pending', 'deploying'].includes(s)) return 'active';
    if (['canceled', 'cancelled', 'removed', 'slept'].includes(s)) return 'canceled';
    return 'active';
  };

  let totalDurationMs = 0;
  let durationCount = 0;

  for (const dep of deployments) {
    const createdAt = dep.createdAt ? new Date(dep.createdAt) : null;
    const cat = classify(dep.status);

    summary.total++;
    summary[cat]++;

    if (createdAt) {
      const dateStr = createdAt.toISOString().slice(0, 10);
      const day = daily.find((d) => d.date === dateStr);
      if (day) {
        day.total++;
        day[cat]++;
      }
    }

    const pId = dep.projectId || 'unknown';
    const pName = dep.projectName || pId;
    let proj = projectMap.get(pId);
    if (!proj) {
      proj = { projectId: pId, name: pName, ...emptyCounts() };
      projectMap.set(pId, proj);
    }
    proj.total++;
    proj[cat]++;

    if (dep.finishedAt && dep.createdAt) {
      const ms = new Date(dep.finishedAt).getTime() - new Date(dep.createdAt).getTime();
      if (ms > 0) {
        totalDurationMs += ms;
        durationCount++;
      }
    }
  }

  summary.avgBuildDurationMs =
    durationCount > 0 ? Math.round(totalDurationMs / durationCount) : null;
  summary.projects = Array.from(projectMap.values()).sort(
    (a, b) => b.failed - a.failed || b.active - a.active || b.total - a.total
  );

  return summary;
}
