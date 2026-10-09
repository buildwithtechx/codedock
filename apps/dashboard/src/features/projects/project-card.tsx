import { Link } from '@tanstack/react-router';
import { ArrowRight, Box, Cloud, Database, Folder, FolderKanban } from 'lucide-react';
import type { CanvasSummary } from '#/features/projects';

const GithubIcon = ({ className }: { className?: string }) => (
  <svg
    viewBox="0 0 24 24"
    fill="currentColor"
    stroke="currentColor"
    strokeWidth="2"
    strokeLinecap="round"
    strokeLinejoin="round"
    className={className}
  >
    <path d="M15 22v-4a4.8 4.8 0 0 0-1-3.2c3-.3 6-1.5 6-6.5 0-1.4-.5-2.5-1.5-3.4.1-.3.6-1.6-.1-3.3 0 0-1.2-.4-3.8 1.4a12.8 12.8 0 0 0-7 0C3.9 1.5 2.7 1.9 2.7 1.9c-.7 1.7-.2 3 .1 3.3-1 1-1.5 2-1.5 3.4 0 5 3 6.2 6 6.5-.4.4-.7 1-.8 2.2-.8.4-2.8.9-4-1.1 0 0-.7-1.3-2-1.4 0 0-1.3-.1-.1 1.2 0 0 1.2 1.8 3 2.5 1.5.5 3.3.4 3.3.4z" />
  </svg>
);

const IconMap: Record<string, React.ReactNode> = {
  github: <GithubIcon className="h-4 w-4" />,
  postgres: <Database className="h-4 w-4 text-blue-500" />,
  mysql: <Database className="h-4 w-4 text-blue-400" />,
  redis: <Database className="h-4 w-4 text-red-500" />,
  s3: <Cloud className="h-4 w-4 text-amber-500" />,
  local: <Folder className="h-4 w-4 text-gray-500" />,
};

function getIcon(iconName: string) {
  return IconMap[iconName.toLowerCase()] || <Box className="h-4 w-4 text-primary" />;
}

function formatUpdated(value: string) {
  const parsed = new Date(value);
  if (Number.isNaN(parsed.getTime())) return null;
  return parsed.toLocaleDateString(undefined, { month: 'short', day: 'numeric' });
}

export const ProjectCard = ({ project }: { project: CanvasSummary }) => {
  const updated = formatUpdated(project.updatedAt || project.createdAt);
  const allOnline = project.totalServices > 0 && project.onlineServices === project.totalServices;

  return (
    <div className="group relative flex h-full flex-col gap-3.5 rounded-2xl border border-border/70 bg-card p-4 shadow-sm transition-colors hover:border-primary/45 hover:bg-primary/[0.03]">
      <Link
        to="/projects/$projectId"
        params={{ projectId: project.id }}
        aria-label={project.name}
        className="absolute inset-0 z-0 rounded-2xl"
      />

      <div className="flex items-start gap-3">
        <div className="flex h-10 w-10 shrink-0 items-center justify-center overflow-hidden rounded-xl bg-muted/60 transition-colors group-hover:bg-muted">
          {project.serviceIcons?.[0] ? (
            getIcon(project.serviceIcons[0])
          ) : (
            <FolderKanban className="h-5 w-5 text-muted-foreground" />
          )}
        </div>
        <div className="min-w-0 flex-1 text-start">
          <div className="flex min-w-0 items-center gap-1.5">
            <p
              className="min-w-0 truncate font-medium text-foreground text-sm"
              title={project.name}
            >
              {project.name}
            </p>
            {project.totalServices > 0 && (
              <span className="shrink-0 rounded-md bg-muted/60 px-1.5 py-0.5 font-medium font-mono text-[10px] text-muted-foreground">
                {project.onlineServices}/{project.totalServices}
              </span>
            )}
          </div>
          {project.description && (
            <p className="mt-0.5 truncate text-muted-foreground text-xs">{project.description}</p>
          )}
        </div>
      </div>

      <div className="flex min-w-0 flex-wrap items-center gap-x-3 gap-y-1.5 whitespace-nowrap text-muted-foreground">
        {project.defaultEnvironment ? (
          <span className="inline-flex min-w-0 max-w-full items-center gap-1.5 text-xs">
            <span className="h-1.5 w-1.5 shrink-0 rounded-full bg-emerald-500/80" />
            <span className="truncate">{project.defaultEnvironment.name}</span>
          </span>
        ) : (
          <span className="inline-flex min-w-0 max-w-full items-center gap-1.5 text-xs">
            <span className="h-1.5 w-1.5 shrink-0 rounded-full bg-zinc-500/80" />
            <span className="truncate">No environment</span>
          </span>
        )}
        <span className="inline-flex items-center rounded-md bg-secondary px-2 py-0.5 text-secondary-foreground text-xs">
          {project.totalServices} service{project.totalServices === 1 ? '' : 's'}
        </span>
        {project.databasesCount > 0 && (
          <span className="inline-flex items-center gap-1.5 text-xs">
            <Database className="h-3.5 w-3.5 shrink-0" />
            <span>
              {project.databasesCount} database{project.databasesCount === 1 ? '' : 's'}
            </span>
          </span>
        )}
        {project.serviceIcons && project.serviceIcons.length > 1 && (
          <span className="flex items-center">
            {project.serviceIcons.slice(1, 5).map((icon, index) => (
              <span
                key={`${icon}-${index}`}
                className="-ml-1 flex h-6 w-6 items-center justify-center rounded-full border border-border/70 bg-background first:ml-0"
              >
                {getIcon(icon)}
              </span>
            ))}
            {project.serviceIcons.length > 5 && (
              <span className="-ml-1 flex h-6 items-center rounded-full border border-border/70 bg-background px-1.5 font-mono text-[10px]">
                +{project.serviceIcons.length - 5}
              </span>
            )}
          </span>
        )}
      </div>

      <div className="mt-auto flex items-center justify-between gap-2 pt-0.5">
        <span
          className={`shrink-0 whitespace-nowrap rounded-full px-2 py-0.5 font-medium text-[10px] ${
            project.totalServices === 0
              ? 'bg-muted text-muted-foreground'
              : allOnline
                ? 'bg-emerald-500/10 text-emerald-600 dark:text-emerald-400'
                : 'bg-amber-500/10 text-amber-600 dark:text-amber-400'
          }`}
        >
          {project.totalServices === 0 ? 'Empty' : allOnline ? 'All online' : 'Partial outage'}
        </span>
        <span className="flex min-w-0 items-center gap-2">
          {updated && (
            <span className="truncate text-muted-foreground text-xs">Updated {updated}</span>
          )}
          <ArrowRight className="h-4 w-4 shrink-0 text-muted-foreground/40 transition-colors group-hover:text-muted-foreground" />
        </span>
      </div>
    </div>
  );
};
