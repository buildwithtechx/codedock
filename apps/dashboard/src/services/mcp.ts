import type { McpCatalogEnvelope, McpCatalogPayload } from '#/interfaces/mcp';
import { apiClient } from '#/lib/api-client';

export const mcpService = {
  async getCatalog(): Promise<McpCatalogPayload> {
    const res = await apiClient.get<McpCatalogEnvelope>('/mcp/catalog');
    return res.data ?? {};
  },
};
