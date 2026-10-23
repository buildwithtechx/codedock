import { useNavigate } from '@tanstack/react-router';
import {
  Activity,
  ArrowRight,
  Database,
  FolderKanban,
  LayoutGrid,
  Moon,
  Server,
  Sparkles,
  Sun,
} from 'lucide-react';
import { useEffect, useState } from 'react';
import { toast } from 'sonner';
import {
  CommandDialog,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
  CommandSeparator,
  CommandShortcut,
  useTheme,
} from '@codedock/ui';
import { useListAppsByOrganization } from '#/hooks/use-apps';
import { useListDatabases } from '#/hooks/use-databases';
import { useListProjects } from '#/hooks/use-projects';
import { useListServers } from '#/hooks/use-servers';
import { commandNavigation } from './dashboard-navigation';

interface CommandPaletteProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

function extractItems<T>(res: unknown): T[] {
  if (Array.isArray(res)) return res as T[];
  const data = (res as { data?: unknown })?.data;
  if (Array.isArray(data)) return data as T[];
  if (Array.isArray((data as { items?: unknown[] })?.items)) return (data as { items: T[] }).items;
  return [];
}

export function CommandPalette({ open, onOpenChange }: CommandPaletteProps) {
  const navigate = useNavigate();
  const { theme, setTheme } = useTheme();
  const [search, setSearch] = useState('');

  const { data: projectsRes } = useListProjects({ limit: 10 });
  const { data: appsRes } = useListAppsByOrganization();
  const { data: serversRes } = useListServers();
  const { data: databasesRes } = useListDatabases();

  useEffect(() => {
    const handleKeyDown = (event: KeyboardEvent) => {
      if (event.key.toLowerCase() !== 'k' || (!event.metaKey && !event.ctrlKey)) {
        return;
      }
      event.preventDefault();
      onOpenChange(!open);
    };

    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [onOpenChange, open]);

  const navigateTo = (to: string, searchParams?: Record<string, string>) => {
    navigate({ to: to as never, search: searchParams as never });
    onOpenChange(false);
  };

  const handleToggleTheme = () => {
    const next = theme === 'light' ? 'dim' : theme === 'dim' ? 'dark' : 'light';
    setTheme(next);
    toast.success(`Theme switched to ${next}`);
    onOpenChange(false);
  };

  const projects = extractItems<{ id: string; name: string; description?: string }>(projectsRes);
  const apps = extractItems<{ id: string; name: string; projectId: string }>(appsRes);
  const servers = extractItems<{ id: string; name: string; ip: string }>(serversRes);
  const databases = extractItems<{ id: string; name: string; projectId?: string; type?: string }>(
    databasesRes
  );

  return (
    <CommandDialog
      open={open}
      onOpenChange={onOpenChange}
      title="Global Command Palette"
      description="Search projects, services, servers, and execute commands."
      showCloseButton={false}
      className="top-[14%] max-w-2xl border-border/80 bg-popover shadow-2xl"
    >
      <CommandInput
        value={search}
        onValueChange={setSearch}
        placeholder="Type a command or search projects, apps, servers..."
      />
      <CommandList className="max-h-95 p-2">
        <CommandEmpty>No matching resources or commands.</CommandEmpty>

        <CommandGroup heading="Actions & Shortcuts">
          <CommandItem onSelect={handleToggleTheme} value="toggle theme mode dark light">
            <div className="flex size-7 items-center justify-center rounded-md border bg-background">
              {theme === 'dark' ? <Sun className="size-3.5" /> : <Moon className="size-3.5" />}
            </div>
            <span className="flex-1 font-medium">Toggle Theme</span>
            <CommandShortcut>⌘T</CommandShortcut>
          </CommandItem>
          <CommandItem onSelect={() => navigateTo('/library')} value="new project create">
            <div className="flex size-7 items-center justify-center rounded-md border bg-background">
              <FolderKanban className="size-3.5 text-primary" />
            </div>
            <span className="flex-1 font-medium">New Project</span>
            <CommandShortcut>⌘N</CommandShortcut>
          </CommandItem>
          <CommandItem
            onSelect={() => navigateTo('/library', { tab: 'apps' })}
            value="deploy app template"
          >
            <div className="flex size-7 items-center justify-center rounded-md border bg-background">
              <Sparkles className="size-3.5 text-amber-500" />
            </div>
            <span className="flex-1 font-medium">Deploy App / Template</span>
          </CommandItem>
          <CommandItem onSelect={() => navigateTo('/servers/new')} value="add server connect host">
            <div className="flex size-7 items-center justify-center rounded-md border bg-background">
              <Server className="size-3.5 text-blue-500" />
            </div>
            <span className="flex-1 font-medium">Connect Server</span>
          </CommandItem>
          <CommandItem
            onSelect={() => navigateTo('/library', { tab: 'databases' })}
            value="new database postgres redis mysql"
          >
            <div className="flex size-7 items-center justify-center rounded-md border bg-background">
              <Database className="size-3.5 text-emerald-500" />
            </div>
            <span className="flex-1 font-medium">Provision Database</span>
          </CommandItem>
          <CommandItem onSelect={() => navigateTo('/monitoring')} value="monitoring metrics health">
            <div className="flex size-7 items-center justify-center rounded-md border bg-background">
              <Activity className="size-3.5 text-rose-500" />
            </div>
            <span className="flex-1 font-medium">Monitoring & Alerts</span>
          </CommandItem>
        </CommandGroup>

        {projects.length > 0 && (
          <>
            <CommandSeparator className="my-1.5" />
            <CommandGroup heading="Projects">
              {projects.slice(0, 5).map((project) => (
                <CommandItem
                  key={`p-${project.id}`}
                  value={`project ${project.name} ${project.description ?? ''}`}
                  onSelect={() => navigateTo(`/projects/${project.id}`)}
                >
                  <div className="flex size-7 items-center justify-center rounded-md border bg-background">
                    <FolderKanban className="size-3.5 text-muted-foreground" />
                  </div>
                  <div className="min-w-0 flex-1">
                    <p className="truncate font-medium">{project.name}</p>
                    {project.description && (
                      <p className="truncate text-muted-foreground text-xs">
                        {project.description}
                      </p>
                    )}
                  </div>
                  <ArrowRight className="size-3 text-muted-foreground/60" />
                </CommandItem>
              ))}
            </CommandGroup>
          </>
        )}

        {apps.length > 0 && (
          <>
            <CommandSeparator className="my-1.5" />
            <CommandGroup heading="Apps & Services">
              {apps.slice(0, 5).map((app) => (
                <CommandItem
                  key={`app-${app.id}`}
                  value={`app ${app.name} service`}
                  onSelect={() => navigateTo(`/projects/${app.projectId}`)}
                >
                  <div className="flex size-7 items-center justify-center rounded-md border bg-background">
                    <LayoutGrid className="size-3.5 text-muted-foreground" />
                  </div>
                  <div className="min-w-0 flex-1">
                    <p className="truncate font-medium">{app.name}</p>
                    <p className="truncate text-muted-foreground text-xs">Service</p>
                  </div>
                  <ArrowRight className="size-3 text-muted-foreground/60" />
                </CommandItem>
              ))}
            </CommandGroup>
          </>
        )}

        {servers.length > 0 && (
          <>
            <CommandSeparator className="my-1.5" />
            <CommandGroup heading="Servers">
              {servers.slice(0, 5).map((srv) => (
                <CommandItem
                  key={`srv-${srv.id}`}
                  value={`server ${srv.name} ${srv.ip}`}
                  onSelect={() => navigateTo(`/servers/${srv.id}`)}
                >
                  <div className="flex size-7 items-center justify-center rounded-md border bg-background">
                    <Server className="size-3.5 text-muted-foreground" />
                  </div>
                  <div className="min-w-0 flex-1">
                    <p className="truncate font-medium">{srv.name}</p>
                    <p className="font-mono text-muted-foreground text-xs">{srv.ip}</p>
                  </div>
                  <ArrowRight className="size-3 text-muted-foreground/60" />
                </CommandItem>
              ))}
            </CommandGroup>
          </>
        )}

        {databases.length > 0 && (
          <>
            <CommandSeparator className="my-1.5" />
            <CommandGroup heading="Databases">
              {databases.slice(0, 5).map((db) => (
                <CommandItem
                  key={`db-${db.id}`}
                  value={`database ${db.name} ${db.type ?? ''}`}
                  onSelect={() =>
                    db.projectId ? navigateTo(`/projects/${db.projectId}`) : navigateTo('/projects')
                  }
                >
                  <div className="flex size-7 items-center justify-center rounded-md border bg-background">
                    <Database className="size-3.5 text-muted-foreground" />
                  </div>
                  <div className="min-w-0 flex-1">
                    <p className="truncate font-medium">{db.name}</p>
                    <p className="text-muted-foreground text-xs uppercase">
                      {db.type ?? 'Managed'}
                    </p>
                  </div>
                  <ArrowRight className="size-3 text-muted-foreground/60" />
                </CommandItem>
              ))}
            </CommandGroup>
          </>
        )}

        <CommandSeparator className="my-1.5" />
        <CommandGroup heading="Navigation">
          {commandNavigation.map((item) => (
            <CommandItem
              key={`${item.title}:${item.to}`}
              value={`${item.title} ${item.description}`}
              onSelect={() => navigateTo(item.to, item.search)}
            >
              <div className="flex size-7 items-center justify-center rounded-md border bg-background">
                <item.icon className="size-3.5 text-muted-foreground" />
              </div>
              <div className="min-w-0 flex-1">
                <p className="font-medium">{item.title}</p>
                <p className="truncate text-muted-foreground text-xs">{item.description}</p>
              </div>
            </CommandItem>
          ))}
        </CommandGroup>
      </CommandList>

      <div className="flex items-center justify-between border-t bg-muted/20 px-4 py-2.5 text-muted-foreground text-xs">
        <div className="flex items-center gap-2">
          <kbd className="rounded border bg-background px-1.5 py-0.5 font-mono text-[10px]">↑↓</kbd>
          <span>Navigate</span>
        </div>
        <div className="flex items-center gap-2">
          <kbd className="rounded border bg-background px-1.5 py-0.5 font-mono text-[10px]">↵</kbd>
          <span>Select</span>
        </div>
        <div className="flex items-center gap-2">
          <kbd className="rounded border bg-background px-1.5 py-0.5 font-mono text-[10px]">
            Esc
          </kbd>
          <span>Close</span>
        </div>
      </div>
    </CommandDialog>
  );
}
