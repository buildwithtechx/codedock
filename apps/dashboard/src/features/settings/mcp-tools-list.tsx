import { Box, Database, HardDrive, Play, RefreshCw, Terminal, Wrench } from 'lucide-react';

const MCP_TOOLS = [
  {
    name: 'list_projects',
    description: 'List all deployment projects registered in this Codedock instance.',
    icon: Box,
    category: 'Projects',
  },
  {
    name: 'get_project',
    description: 'Get detailed information about a project including environments and services.',
    icon: Box,
    category: 'Projects',
  },
  {
    name: 'list_apps',
    description: 'List all application services across projects or scoped to a project.',
    icon: Terminal,
    category: 'Services',
  },
  {
    name: 'get_app',
    description: 'Get service runtime status, ports, domains, branch, and configuration.',
    icon: Terminal,
    category: 'Services',
  },
  {
    name: 'redeploy_app',
    description: 'Trigger a fresh zero-downtime deployment build for an application service.',
    icon: RefreshCw,
    category: 'Deployments',
  },
  {
    name: 'restart_app',
    description: 'Gracefully restart an active service container instance.',
    icon: Play,
    category: 'Operations',
  },
  {
    name: 'stop_app',
    description: 'Safely halt an application service workload.',
    icon: Wrench,
    category: 'Operations',
  },
  {
    name: 'list_databases',
    description: 'List managed database instances and their current provisioning status.',
    icon: Database,
    category: 'Databases',
  },
  {
    name: 'get_database',
    description: 'Get database connection parameters, storage volumes, and credentials.',
    icon: Database,
    category: 'Databases',
  },
  {
    name: 'list_deployments',
    description: 'Retrieve real-time build logs and release history for a service.',
    icon: HardDrive,
    category: 'Deployments',
  },
  {
    name: 'get_system_status',
    description: 'Inspect CPU, memory usage, Docker daemon health, and host metrics.',
    icon: Wrench,
    category: 'System',
  },
];

export function McpToolsList() {
  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <div>
          <h3 className="font-semibold text-foreground text-sm">Available MCP Tools</h3>
          <p className="text-muted-foreground text-xs">
            Exposed tools that connected AI assistants can safely invoke.
          </p>
        </div>
        <span className="rounded-full bg-primary/10 px-2.5 py-0.5 font-medium text-primary text-xs">
          {MCP_TOOLS.length} Tools
        </span>
      </div>

      <div className="grid grid-cols-1 gap-2.5 sm:grid-cols-2">
        {MCP_TOOLS.map((tool) => {
          const Icon = tool.icon;
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
                  <code className="font-mono font-semibold text-[11px] text-foreground">
                    {tool.name}
                  </code>
                  <span className="text-[10px] text-muted-foreground/70 uppercase tracking-wider">
                    {tool.category}
                  </span>
                </div>
                <p className="mt-1 line-clamp-2 text-muted-foreground text-xs leading-relaxed">
                  {tool.description}
                </p>
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
}
