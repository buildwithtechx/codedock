import { useListAppsByOrganization } from '#/hooks/use-apps';
import { useListAttentionIssues } from '#/hooks/use-issues';
import { useListProjects } from '#/hooks/use-projects';
import { useListServers } from '#/hooks/use-servers';

export type SidebarCounts = {
  projects: number | null;
  apps: number | null;
  servers: number | null;
  serverStatus: 'healthy' | 'warning' | 'error' | null;
  hasIssues: boolean;
};

export function useSidebarCounts(): SidebarCounts {
  const { data: projects } = useListProjects({ page: 1, limit: 1 });
  const { data: apps } = useListAppsByOrganization();
  const { data: servers } = useListServers();
  const { data: issuesRes } = useListAttentionIssues();

  const serverList = Array.isArray(servers) ? servers : [];
  const serverStatus =
    serverList.length === 0
      ? null
      : serverList.every((s) => s.ipAddress || s.sshHost)
        ? 'healthy'
        : 'warning';

  const issueList = issuesRes?.data?.issues ?? [];
  const hasIssues = issueList.some((i) => i.status === 'open');

  return {
    projects: projects?.data.total ?? null,
    apps: apps ? apps.data.length : null,
    servers: serverList.length > 0 ? serverList.length : null,
    serverStatus,
    hasIssues,
  };
}
