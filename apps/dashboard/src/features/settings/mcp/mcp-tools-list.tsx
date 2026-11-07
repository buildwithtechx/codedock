import {
  BookOpen,
  Box,
  Database,
  HardDrive,
  RefreshCw,
  Search,
  ShieldAlert,
  Terminal,
  Wrench,
} from 'lucide-react';
import { useMemo, useState } from 'react';
import { Button } from '@codedock/ui/button';
import { Input } from '@codedock/ui/input';
import { useMcpCatalog } from '#/hooks/use-mcp';
import type { McpCatalogPrompt, McpCatalogTool } from '#/interfaces/mcp';

const FALLBACK_TOOLS: McpCatalogTool[] = [
  {
    name: 'get_projects',
    description: 'List all deployment projects registered in this Codedock instance.',
    method: 'GET',
    path: '/api/projects',
    category: 'Projects',
    auth: 'user',
    annotations: { readOnlyHint: true, destructiveHint: false },
  },
  {
    name: 'get_projects_by_id',
    description: 'Get detailed information about a project including environments and services.',
    method: 'GET',
    path: '/api/projects/:id',
    category: 'Projects',
    auth: 'user',
    annotations: { readOnlyHint: true, destructiveHint: false },
  },
  {
    name: 'get_apps',
    description: 'List all application services across projects or scoped to a project.',
    method: 'GET',
    path: '/api/apps',
    category: 'Services',
    auth: 'user',
    annotations: { readOnlyHint: true, destructiveHint: false },
  },
  {
    name: 'get_apps_by_id',
    description: 'Get service runtime status, ports, domains, branch, and configuration.',
    method: 'GET',
    path: '/api/apps/:id',
    category: 'Services',
    auth: 'user',
    annotations: { readOnlyHint: true, destructiveHint: false },
  },
  {
    name: 'post_apps_by_id_redeploy',
    description: 'Trigger a fresh zero-downtime deployment build for an application service.',
    method: 'POST',
    path: '/api/apps/:id/redeploy',
    category: 'Deployments',
    auth: 'user',
    annotations: { readOnlyHint: false, destructiveHint: false },
  },
  {
    name: 'post_apps_by_id_restart',
    description: 'Gracefully restart an active service container instance.',
    method: 'POST',
    path: '/api/apps/:id/restart',
    category: 'Operations',
    auth: 'user',
    annotations: { readOnlyHint: false, destructiveHint: true },
  },
  {
    name: 'post_apps_by_id_stop',
    description: 'Safely halt an application service workload.',
    method: 'POST',
    path: '/api/apps/:id/stop',
    category: 'Operations',
    auth: 'user',
    annotations: { readOnlyHint: false, destructiveHint: true },
  },
  {
    name: 'get_databases',
    description: 'List managed database instances and their current provisioning status.',
    method: 'GET',
    path: '/api/databases',
    category: 'Databases',
    auth: 'user',
    annotations: { readOnlyHint: true, destructiveHint: false },
  },
  {
    name: 'get_databases_by_id',
    description: 'Get database connection parameters, storage volumes, and credentials.',
    method: 'GET',
    path: '/api/databases/:id',
    category: 'Databases',
    auth: 'user',
    annotations: { readOnlyHint: true, destructiveHint: false },
  },
  {
    name: 'get_deployments',
    description: 'Retrieve real-time build logs and release history for a service.',
    method: 'GET',
    path: '/api/deployments',
    category: 'Deployments',
    auth: 'user',
    annotations: { readOnlyHint: true, destructiveHint: false },
  },
  {
    name: 'get_system_stats',
    description: 'Inspect CPU, memory usage, Docker daemon health, and host metrics.',
    method: 'GET',
    path: '/api/system/stats',
    category: 'System',
    auth: 'user',
    annotations: { readOnlyHint: true, destructiveHint: false },
  },
];

const FALLBACK_PROMPTS: McpCatalogPrompt[] = [
  {
    name: 'codedock-overview',
    title: 'Codedock: How to drive the platform via MCP',
    description:
      'Orientation on Codedock tool groups, workspace scoping, and safe operational workflows.',
  },
  {
    name: 'deploy-from-git',
    title: 'Deploy an application from a Git repository',
    description:
      'End-to-end flow: verify Git connection, inspect repository, create project/service, deploy, and tail logs.',
  },
  {
    name: 'install-catalog-app',
    title: 'Install a one-click catalog application',
    description: 'Browse the one-click template catalog, review configuration, and deploy.',
  },
  {
    name: 'diagnose-and-fix-issues',
    title: 'Diagnose unhealthy services and failed deployments',
    description:
      'Inspect active platform issues, read deployment/runtime logs, and remediate or rollback.',
  },
  {
    name: 'backup-and-restore',
    title: 'Configure and run S3/SFTP backups and restores',
    description:
      'Verify backup storage destinations, trigger backup runs, and review/apply restores.',
  },
];

function iconForCategory(cat: string) {
  const lower = cat.toLowerCase();
  if (lower.includes('project')) return Box;
  if (lower.includes('database')) return Database;
  if (lower.includes('deploy') || lower.includes('backup')) return HardDrive;
  if (lower.includes('service') || lower.includes('app')) return Terminal;
  return Wrench;
}

export function McpToolsList() {
  const [search, setSearch] = useState('');
  const [category, setCategory] = useState('All');
  const [view, setView] = useState<'tools' | 'prompts'>('tools');

  const { data: catalog, isFetching, refetch } = useMcpCatalog();

  const tools = catalog?.tools?.length ? catalog.tools : FALLBACK_TOOLS;
  const prompts = catalog?.prompts?.length ? catalog.prompts : FALLBACK_PROMPTS;

  const categories = useMemo(() => {
    const set = new Set<string>(['All']);
    for (const t of tools) {
      if (t.category) set.add(t.category);
    }
    return Array.from(set);
  }, [tools]);

  const filteredTools = useMemo(() => {
    const q = search.trim().toLowerCase();
    return tools.filter((t) => {
      if (category !== 'All' && t.category !== category) return false;
      if (!q) return true;
      return `${t.name} ${t.description} ${t.path ?? ''} ${t.category}`.toLowerCase().includes(q);
    });
  }, [tools, category, search]);

  return (
    <div className="space-y-4">
      <div className="flex flex-wrap items-center justify-between gap-3">
        <div className="inline-flex rounded-xl border border-border/60 bg-muted/30 p-1">
          <button
            type="button"
            onClick={() => setView('tools')}
            className={`rounded-lg px-3 py-1.5 font-medium text-xs transition-colors ${
              view === 'tools'
                ? 'bg-card text-foreground shadow-xs'
                : 'text-muted-foreground hover:text-foreground'
            }`}
          >
            Available Tools ({tools.length})
          </button>
          <button
            type="button"
            onClick={() => setView('prompts')}
            className={`rounded-lg px-3 py-1.5 font-medium text-xs transition-colors ${
              view === 'prompts'
                ? 'bg-card text-foreground shadow-xs'
                : 'text-muted-foreground hover:text-foreground'
            }`}
          >
            Guided Flows ({prompts.length})
          </button>
        </div>

        <Button
          type="button"
          variant="outline"
          size="sm"
          onClick={() => void refetch()}
          disabled={isFetching}
          className="h-8 gap-1.5 text-xs"
        >
          <RefreshCw className={`size-3.5 ${isFetching ? 'animate-spin' : ''}`} />
          Refresh
        </Button>
      </div>

      {view === 'tools' ? (
        <>
          <div className="flex flex-col gap-3 sm:flex-row sm:items-center">
            <div className="relative flex-1">
              <Search className="pointer-events-none absolute top-1/2 left-3 size-3.5 -translate-y-1/2 text-muted-foreground" />
              <Input
                type="search"
                value={search}
                onChange={(e) => setSearch(e.target.value)}
                placeholder="Search tools by name, route, or description..."
                className="h-9 pl-9 text-xs"
              />
            </div>
            <div className="scrollbar-none flex items-center gap-1 overflow-x-auto pb-1">
              {categories.slice(0, 8).map((cat) => (
                <button
                  key={cat}
                  type="button"
                  onClick={() => setCategory(cat)}
                  className={`shrink-0 rounded-lg px-2.5 py-1 font-medium text-[11px] transition-colors ${
                    category === cat
                      ? 'bg-primary/15 text-primary'
                      : 'text-muted-foreground hover:bg-muted/40 hover:text-foreground'
                  }`}
                >
                  {cat}
                </button>
              ))}
            </div>
          </div>

          <div className="grid max-h-105 grid-cols-1 gap-2.5 overflow-y-auto pr-1 sm:grid-cols-2">
            {filteredTools.map((tool) => {
              const Icon = iconForCategory(tool.category);
              return (
                <div
                  key={tool.name}
                  className="flex items-start gap-3 rounded-xl border border-border/50 bg-muted/15 p-3 transition-colors hover:bg-muted/30"
                >
                  <div className="flex size-7 shrink-0 items-center justify-center rounded-lg bg-primary/10 text-primary">
                    <Icon className="size-3.5" />
                  </div>
                  <div className="min-w-0 flex-1">
                    <div className="flex items-center justify-between gap-2">
                      <code className="truncate font-mono font-semibold text-[11px] text-foreground">
                        {tool.name}
                      </code>
                      <div className="flex shrink-0 items-center gap-1">
                        {tool.annotations?.destructiveHint && (
                          <span className="inline-flex items-center gap-0.5 rounded bg-destructive/10 px-1.5 py-0.5 font-medium text-[9px] text-destructive">
                            <ShieldAlert className="size-2.5" />
                            destructive
                          </span>
                        )}
                        <span className="text-[10px] text-muted-foreground/70 uppercase tracking-wider">
                          {tool.category}
                        </span>
                      </div>
                    </div>
                    {tool.path && (
                      <p className="mt-0.5 truncate font-mono text-[10px] text-muted-foreground/80">
                        {tool.method} {tool.path}
                      </p>
                    )}
                    <p className="mt-1 line-clamp-2 text-muted-foreground text-xs leading-relaxed">
                      {tool.description}
                    </p>
                  </div>
                </div>
              );
            })}
          </div>
        </>
      ) : (
        <div className="grid grid-cols-1 gap-2.5 sm:grid-cols-2">
          {prompts.map((prompt) => (
            <div
              key={prompt.name}
              className="flex items-start gap-3 rounded-xl border border-border/50 bg-muted/15 p-3.5"
            >
              <div className="flex size-7 shrink-0 items-center justify-center rounded-lg bg-emerald-500/10 text-emerald-500">
                <BookOpen className="size-3.5" />
              </div>
              <div className="min-w-0 flex-1">
                <code className="font-mono font-semibold text-foreground text-xs">
                  /{prompt.name}
                </code>
                <p className="mt-0.5 font-medium text-foreground/90 text-xs">{prompt.title}</p>
                <p className="mt-1 text-muted-foreground text-xs leading-relaxed">
                  {prompt.description}
                </p>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
