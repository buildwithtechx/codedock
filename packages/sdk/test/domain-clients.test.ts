import { describe, expect, it, vi } from 'vitest';
import { createCodedockClient } from '../src/index';

describe('CodedockClient Domain Operations', () => {
  it('handles projects operations', async () => {
    const mockFetch = vi.fn().mockImplementation(async (url: string) => {
      if (url.includes('/projects?organizationId=org-1')) {
        return {
          ok: true,
          status: 200,
          json: async () => ({
            data: {
              records: [{ id: 'proj-1', name: 'Frontend' }],
              total: 1,
              page: 1,
              totalPages: 1,
            },
          }),
        };
      }
      if (url.endsWith('/projects/proj-1')) {
        return {
          ok: true,
          status: 200,
          json: async () => ({
            data: { id: 'proj-1', name: 'Frontend' },
          }),
        };
      }
      if (url.endsWith('/projects/proj-1/domains')) {
        return {
          ok: true,
          status: 200,
          json: async () => ({
            data: [
              {
                id: 'dom-1',
                domain: 'app.codedock.io',
                domainName: 'app.codedock.io',
                isCustom: true,
              },
            ],
          }),
        };
      }
      return { ok: true, status: 200, json: async () => ({ data: {} }) };
    });

    const client = createCodedockClient({ fetch: mockFetch });
    const listRes = await client.projects.list({ organizationId: 'org-1' });
    expect(listRes.records).toHaveLength(1);
    expect(listRes.records[0].name).toBe('Frontend');

    const project = await client.projects.get('proj-1');
    expect(project.id).toBe('proj-1');

    const domains = await client.projects.listDomains('proj-1');
    expect(domains[0].domain).toBe('app.codedock.io');
  });

  it('handles deployments operations', async () => {
    const mockFetch = vi.fn().mockImplementation(async (url: string) => {
      if (url.includes('/deployments?organizationId=org-1')) {
        return {
          ok: true,
          status: 200,
          json: async () => ({
            data: [
              {
                id: 'dep-1',
                serviceId: 'svc-1',
                status: 'SUCCESS',
                projectName: 'API',
                serviceName: 'Backend',
              },
            ],
          }),
        };
      }
      if (url.endsWith('/services/svc-1/deploy')) {
        return {
          ok: true,
          status: 200,
          json: async () => ({
            data: { id: 'dep-2', serviceId: 'svc-1', status: 'BUILDING' },
          }),
        };
      }
      if (url.endsWith('/deployments/dep-1/explain-failure')) {
        return {
          ok: true,
          status: 200,
          json: async () => ({
            data: {
              summary: 'OOM Killed',
              cause: 'Process exceeded container memory limits',
              suggestedFix: 'Increase memory limit to 512MB',
              confidence: 'high',
              commands: [],
              relatedLogLines: [],
            },
          }),
        };
      }
      return { ok: true, status: 200, json: async () => ({ data: {} }) };
    });

    const client = createCodedockClient({ fetch: mockFetch });
    const orgDeployments = await client.deployments.listByOrganization('org-1');
    expect(orgDeployments).toHaveLength(1);
    expect(orgDeployments[0].projectName).toBe('API');

    const triggered = await client.deployments.trigger('svc-1', { commitId: 'abc1234' });
    expect(triggered.status).toBe('BUILDING');

    const explanation = await client.deployments.explainFailure('dep-1');
    expect(explanation.confidence).toBe('high');
    expect(explanation.summary).toBe('OOM Killed');
  });

  it('handles servers operations', async () => {
    const mockFetch = vi.fn().mockImplementation(async (url: string, init?: RequestInit) => {
      if (url.endsWith('/servers') && init?.method === 'POST') {
        return {
          ok: true,
          status: 200,
          json: async () => ({
            data: { id: 'srv-1', name: 'Edge Node 1', ipAddress: '192.168.1.10', status: 'online' },
          }),
        };
      }
      if (url.endsWith('/servers/test-ssh')) {
        return {
          ok: true,
          status: 200,
          json: async () => ({
            data: { success: true, message: 'SSH handshake established' },
          }),
        };
      }
      return { ok: true, status: 200, json: async () => ({ data: [] }) };
    });

    const client = createCodedockClient({ fetch: mockFetch });
    const server = await client.servers.create({
      name: 'Edge Node 1',
      ipAddress: '192.168.1.10',
    });
    expect(server.name).toBe('Edge Node 1');

    const sshRes = await client.servers.testSSH({ sshHost: '192.168.1.10' });
    expect(sshRes.success).toBe(true);
    expect(sshRes.message).toBe('SSH handshake established');
  });

  it('handles services operations', async () => {
    const mockFetch = vi.fn().mockImplementation(async (url: string) => {
      if (url.endsWith('/apps/svc-1')) {
        return {
          ok: true,
          status: 200,
          json: async () => ({
            data: { id: 'svc-1', name: 'Web App', status: 'running' },
          }),
        };
      }
      if (url.endsWith('/services/svc-1/variables')) {
        return {
          ok: true,
          status: 200,
          json: async () => ({
            data: [
              { id: 'var-1', serviceId: 'svc-1', key: 'PORT', value: '3000', isSecret: false },
            ],
          }),
        };
      }
      return { ok: true, status: 200, json: async () => ({ data: {} }) };
    });

    const client = createCodedockClient({ fetch: mockFetch });
    const service = await client.services.get('svc-1');
    expect(service.id).toBe('svc-1');
    expect(service.name).toBe('Web App');

    const vars = await client.services.listVariables('svc-1');
    expect(vars).toHaveLength(1);
    expect(vars[0].key).toBe('PORT');
  });
});
