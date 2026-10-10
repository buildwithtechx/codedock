import { Link } from '@tanstack/react-router';
import {
  ArrowRight,
  Database,
  ExternalLink,
  FolderKanban,
  GitBranch,
  HardDrive,
  MoreVertical,
  Server,
  Settings,
  Trash2,
} from 'lucide-react';
import { useState } from 'react';
import { Button } from '#/components/ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '#/components/ui/dropdown-menu';
import type { CanvasSummary } from '#/features/projects';
import { ProjectDeleteDialog } from './project-delete-dialog';

function statusTone(project: CanvasSummary) {
  if (project.totalServices === 0) return 'bg-muted text-muted-foreground';
  if (project.onlineServices === project.totalServices)
    return 'bg-emerald-500/10 text-emerald-600 dark:text-emerald-400';
  return 'bg-amber-500/10 text-amber-600 dark:text-amber-400';
}

function statusLabel(project: CanvasSummary) {
  if (project.totalServices === 0) return 'Empty';
  if (project.onlineServices === project.totalServices) return 'All online';
  return 'Partial outage';
}

function timeAgo(value: string) {
  const diff = Date.now() - new Date(value).getTime();
  const minutes = Math.floor(diff / 60_000);
  if (minutes < 1) return 'just now';
  if (minutes < 60) return `${minutes}m ago`;
  const hours = Math.floor(minutes / 60);
  if (hours < 24) return `${hours}h ago`;
  const days = Math.floor(hours / 24);
  if (days < 30) return `${days}d ago`;
  return `${Math.floor(days / 30)}mo ago`;
}

export function ProjectListRow({ project }: { project: CanvasSummary }) {
  const [deleteOpen, setDeleteOpen] = useState(false);
  const repoIcon = project.serviceIcons.includes('github');
  return (
    <div className="group relative transition-colors hover:bg-muted/40">
      <Link
        to="/projects/$projectId"
        params={{ projectId: project.id }}
        aria-label={project.name}
        className="absolute inset-0 z-0"
      />
      <div className="flex items-center gap-4 px-5 py-3.5">
        <div className="flex h-10 w-10 shrink-0 items-center justify-center overflow-hidden rounded-xl bg-muted/60 transition-colors group-hover:bg-muted">
          <FolderKanban className="h-5 w-5 text-muted-foreground" />
        </div>
        <div className="w-44 min-w-0 flex-none text-start xl:w-56">
          <p className="truncate font-medium text-foreground text-sm" title={project.name}>
            {project.name}
          </p>
          {project.description && (
            <p className="mt-0.5 truncate text-muted-foreground text-xs">{project.description}</p>
          )}
        </div>
        <div className="hidden min-w-0 flex-1 items-center gap-3 whitespace-nowrap xl:flex">
          <span
            className="inline-flex min-w-0 max-w-[40%] shrink-0 items-center gap-1.5 rounded-md bg-secondary px-2 py-0.5 text-secondary-foreground text-xs"
            title={`${project.totalServices} services`}
          >
            <span className="truncate">
              {project.totalServices} service{project.totalServices === 1 ? '' : 's'}
            </span>
          </span>
          <span
            className="inline-flex min-w-0 items-center gap-1.5 text-muted-foreground text-xs"
            title={project.deployTarget === 'server' ? project.serverName || 'Server' : 'Local'}
          >
            {project.deployTarget === 'server' ? (
              <Server className="h-3.5 w-3.5 shrink-0" />
            ) : (
              <HardDrive className="h-3.5 w-3.5 shrink-0" />
            )}
            <span className="truncate">
              {project.deployTarget === 'server' ? project.serverName || 'Server' : 'Local'}
            </span>
          </span>
          {repoIcon && (
            <span className="hidden min-w-0 items-center gap-1.5 text-muted-foreground text-xs 2xl:inline-flex">
              <GitBranch className="h-3.5 w-3.5 shrink-0" />
              <span className="truncate">Git</span>
            </span>
          )}
          {project.databasesCount > 0 && (
            <span className="hidden min-w-0 items-center gap-1.5 text-muted-foreground text-xs 2xl:inline-flex">
              <Database className="h-3.5 w-3.5 shrink-0" />
              <span className="truncate">
                {project.databasesCount} database{project.databasesCount === 1 ? '' : 's'}
              </span>
            </span>
          )}
        </div>
        <div className="flex min-w-0 max-w-full shrink-0 items-center gap-3 whitespace-nowrap">
          <span className="hidden text-muted-foreground text-xs 2xl:block">
            {timeAgo(project.updatedAt || project.createdAt)}
          </span>
          <span
            className={`shrink-0 rounded-full px-2 py-0.5 font-medium text-[10px] ${statusTone(project)}`}
          >
            {statusLabel(project)}
          </span>
          <div className="relative z-10 flex items-center gap-1">
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button
                  variant="ghost"
                  size="icon"
                  className="h-8 w-8 text-muted-foreground hover:text-foreground"
                  onClick={(e) => e.stopPropagation()}
                >
                  <MoreVertical className="h-4 w-4" />
                  <span className="sr-only">Open menu</span>
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end" className="w-44">
                <DropdownMenuItem asChild>
                  <Link to="/projects/$projectId" params={{ projectId: project.id }}>
                    <ExternalLink className="mr-2 h-4 w-4" />
                    Open project
                  </Link>
                </DropdownMenuItem>
                <DropdownMenuItem asChild>
                  <Link to="/projects/$projectId/settings" params={{ projectId: project.id }}>
                    <Settings className="mr-2 h-4 w-4" />
                    Settings
                  </Link>
                </DropdownMenuItem>
                <DropdownMenuSeparator />
                <DropdownMenuItem
                  variant="destructive"
                  onClick={(e) => {
                    e.stopPropagation();
                    setDeleteOpen(true);
                  }}
                >
                  <Trash2 className="mr-2 h-4 w-4" />
                  Delete project
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
          <ArrowRight className="h-4 w-4 shrink-0 text-muted-foreground/40 transition-colors group-hover:text-muted-foreground" />
        </div>
      </div>

      <ProjectDeleteDialog isOpen={deleteOpen} onOpenChange={setDeleteOpen} project={project} />
    </div>
  );
}
