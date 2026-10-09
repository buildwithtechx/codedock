import { useListProjects } from '#/features/projects';
import { useListAppsByOrganization } from '#/hooks/use-apps';
import { useListServers } from '#/hooks/use-servers';

export type SidebarCounts = {
  projects: number | null;
  apps: number | null;
  servers: number | null;
};

export function useSidebarCounts(): SidebarCounts {
  const { data: projects } = useListProjects({ page: 1, limit: 1 });
  const { data: apps } = useListAppsByOrganization();
  const { data: servers } = useListServers();

  return {
    projects: projects?.data.total ?? null,
    apps: apps ? apps.data.length : null,
    servers: servers ? servers.length : null,
  };
}
