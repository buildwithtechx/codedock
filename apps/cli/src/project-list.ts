import type { ApiClient } from './api-client.js';
import type { ProjectRecord } from './types.js';

export async function listProjects(
  client: ApiClient,
  organizationId?: string
): Promise<ProjectRecord[]> {
  const organizations = organizationId
    ? [{ id: organizationId }]
    : (await client.get<{ id: string }[]>('/api/organizations')).data || [];
  const projects: ProjectRecord[] = [];
  for (const organization of organizations) {
    let page = 1;
    let totalPages = 1;
    do {
      const result = await client.get<{ records: ProjectRecord[]; totalPages: number }>(
        `/api/projects?organizationId=${encodeURIComponent(organization.id)}&limit=100&page=${page}`
      );
      projects.push(...(result.data?.records || []));
      totalPages = result.data?.totalPages || 1;
      page++;
    } while (page <= totalPages);
  }
  return projects;
}
