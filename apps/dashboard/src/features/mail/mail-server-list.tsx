import {
  Activity,
  AlertCircle,
  ArrowRight,
  Loader2,
  Mail,
  MoreVertical,
  Plus,
  Server,
  Trash2,
} from 'lucide-react';
import { Button } from '#/components/ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '#/components/ui/dropdown-menu';

export interface MailServerItem {
  id: string;
  name?: string;
  domain: string;
  maildUrl: string;
  status?: 'running' | 'installing' | 'incomplete';
}

interface MailServerListProps {
  servers: MailServerItem[];
  onOpen: (serverId: string) => void;
  onAddNew: () => void;
  onRemove: (serverId: string) => void;
}

export function MailServerList({ servers, onOpen, onAddNew, onRemove }: MailServerListProps) {
  const running = servers.filter((s) => s.status === 'running' || !s.status).length;
  const installing = servers.filter((s) => s.status === 'installing').length;
  const incomplete = servers.filter((s) => s.status === 'incomplete').length;
  const runningPct = servers.length > 0 ? Math.round((running / servers.length) * 100) : 0;

  return (
    <div className="grid grid-cols-1 gap-6 lg:grid-cols-[1fr_320px]">
      {/* Registry Rows */}
      <div className="min-w-0">
        <div className="mb-3 flex items-center justify-between gap-4 px-1">
          <h2 className="inline-flex items-baseline gap-2 text-sm font-medium text-muted-foreground">
            Mail servers
            <span className="tabular-nums text-muted-foreground/60">({servers.length})</span>
          </h2>
          <Button onClick={onAddNew} size="sm" className="gap-2 rounded-xl">
            <Plus className="size-4" />
            Add mail server
          </Button>
        </div>

        <div className="rounded-2xl border border-border/50 bg-card overflow-hidden">
          <ul className="divide-y divide-border/50">
            {servers.map((s) => {
              const label = s.domain || s.name || s.id;
              const status = s.status ?? 'running';

              return (
                <li
                  key={s.id}
                  className="group relative flex items-center gap-4 px-5 py-4 transition-colors hover:bg-muted/30 cursor-pointer"
                  onClick={() => onOpen(s.id)}
                >
                  <div className="flex size-10 shrink-0 items-center justify-center rounded-xl bg-muted/60 text-muted-foreground">
                    <Mail className="size-4" />
                  </div>

                  <div className="min-w-0 flex-1">
                    <p className="truncate text-sm font-medium text-foreground">{label}</p>
                    <p className="mt-0.5 truncate font-mono text-xs text-muted-foreground">
                      {s.maildUrl}
                    </p>
                  </div>

                  {/* Status Indicator */}
                  <span
                    className={`inline-flex shrink-0 items-center gap-2 text-xs font-medium ${
                      status === 'running'
                        ? 'text-emerald-500'
                        : status === 'installing'
                          ? 'text-sky-500'
                          : 'text-amber-500'
                    }`}
                  >
                    <span
                      className={`size-2.5 rounded-full border-2 ${
                        status === 'running'
                          ? 'border-emerald-500'
                          : status === 'installing'
                            ? 'border-sky-500 animate-pulse'
                            : 'border-amber-500'
                      }`}
                    />
                    <span className="capitalize">{status}</span>
                  </span>

                  {/* Actions Dropdown */}
                  <div
                    className="relative z-10 flex shrink-0 items-center gap-1"
                    onClick={(e) => e.stopPropagation()}
                  >
                    <DropdownMenu>
                      <DropdownMenuTrigger asChild>
                        <Button variant="ghost" size="icon" className="size-8">
                          <MoreVertical className="size-4 text-muted-foreground" />
                        </Button>
                      </DropdownMenuTrigger>
                      <DropdownMenuContent align="end">
                        <DropdownMenuItem onClick={() => onOpen(s.id)}>
                          <ArrowRight className="mr-2 size-4" />
                          Open Console
                        </DropdownMenuItem>
                        <DropdownMenuSeparator />
                        <DropdownMenuItem
                          className="text-destructive focus:text-destructive"
                          onClick={() => onRemove(s.id)}
                        >
                          <Trash2 className="mr-2 size-4" />
                          Remove from list
                        </DropdownMenuItem>
                      </DropdownMenuContent>
                    </DropdownMenu>
                  </div>
                </li>
              );
            })}
          </ul>
        </div>
      </div>

      {/* Right Rail Overview Card */}
      <aside className="space-y-4 lg:sticky lg:top-6 lg:self-start">
        <div className="rounded-2xl border border-border/50 bg-card p-5 space-y-4">
          <div className="flex items-center gap-2">
            <Activity className="size-4 text-muted-foreground" />
            <h3 className="text-sm font-semibold text-foreground">Registry Overview</h3>
          </div>

          <div className="flex items-baseline justify-between">
            <div className="flex items-baseline gap-1.5">
              <span className="text-2xl font-semibold tabular-nums text-foreground">{running}</span>
              <span className="text-sm text-muted-foreground">
                / {servers.length} running
              </span>
            </div>
            <span className="font-mono text-xs text-muted-foreground">{runningPct}%</span>
          </div>

          <div className="h-1.5 w-full overflow-hidden rounded-full bg-muted">
            <div
              className="h-full rounded-full bg-emerald-500 transition-all duration-500"
              style={{ width: `${runningPct}%` }}
            />
          </div>

          <div className="space-y-2 border-t border-border/50 pt-4 text-xs">
            <div className="flex items-center justify-between text-muted-foreground">
              <span className="flex items-center gap-2">
                <Server className="size-3.5" />
                Total servers
              </span>
              <span className="font-medium text-foreground">{servers.length}</span>
            </div>

            {installing > 0 && (
              <div className="flex items-center justify-between text-sky-500">
                <span className="flex items-center gap-2">
                  <Loader2 className="size-3.5 animate-spin" />
                  Installing
                </span>
                <span className="font-medium">{installing}</span>
              </div>
            )}

            {incomplete > 0 && (
              <div className="flex items-center justify-between text-amber-500">
                <span className="flex items-center gap-2">
                  <AlertCircle className="size-3.5" />
                  Incomplete
                </span>
                <span className="font-medium">{incomplete}</span>
              </div>
            )}
          </div>
        </div>
      </aside>
    </div>
  );
}
