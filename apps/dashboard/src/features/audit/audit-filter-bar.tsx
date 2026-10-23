import { Activity, Bot, Search, Terminal, Webhook, Wrench } from 'lucide-react';
import { Input } from '@codedock/ui/input';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@codedock/ui/select';
import type { AuditFacets, AuditSource } from '#/interfaces/audit';

export type PeriodKey = 'all' | 'today' | '7d' | '30d' | '90d';
export type SourceKey = 'all' | AuditSource;

interface AuditFilterBarProps {
  search: string;
  onSearchChange: (val: string) => void;
  actor: string;
  onActorChange: (val: string) => void;
  period: PeriodKey;
  onPeriodChange: (val: PeriodKey) => void;
  source: SourceKey;
  onSourceChange: (val: SourceKey) => void;
  facets?: AuditFacets | null;
}

const SOURCES: { key: SourceKey; label: string; icon?: typeof Activity }[] = [
  { key: 'all', label: 'Everything' },
  { key: 'dashboard', label: 'Dashboard', icon: Activity },
  { key: 'mcp', label: 'MCP', icon: Bot },
  { key: 'cli', label: 'CLI', icon: Terminal },
  { key: 'api', label: 'API', icon: Webhook },
  { key: 'system', label: 'System', icon: Wrench },
];

export function AuditFilterBar({
  search,
  onSearchChange,
  actor,
  onActorChange,
  period,
  onPeriodChange,
  source,
  onSourceChange,
  facets,
}: AuditFilterBarProps) {
  return (
    <div className="space-y-3">
      <div className="grid gap-3 sm:grid-cols-[1fr_180px_160px]">
        <div className="relative">
          <Search className="pointer-events-none absolute start-3.5 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" />
          <Input
            placeholder="Search apps, servers, domains, people..."
            value={search}
            onChange={(e) => onSearchChange(e.target.value)}
            className="h-10 rounded-xl bg-muted/60 ps-10 pe-4 text-sm"
          />
        </div>

        <div>
          <Select
            value={actor || 'all'}
            onValueChange={(val) => onActorChange(val === 'all' ? '' : val)}
          >
            <SelectTrigger className="h-10 rounded-xl border border-border/60 bg-muted/40 px-3 text-xs">
              <SelectValue placeholder="Anyone" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">Anyone</SelectItem>
              {facets?.actors?.map((a) => (
                <SelectItem key={a.id} value={a.id}>
                  {a.name || a.email || a.id.slice(0, 8)}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>

        <div>
          <Select value={period} onValueChange={(val) => onPeriodChange(val as PeriodKey)}>
            <SelectTrigger className="h-10 rounded-xl border border-border/60 bg-muted/40 px-3 text-xs">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">All time</SelectItem>
              <SelectItem value="today">Today</SelectItem>
              <SelectItem value="7d">Last 7 days</SelectItem>
              <SelectItem value="30d">Last 30 days</SelectItem>
              <SelectItem value="90d">Last 90 days</SelectItem>
            </SelectContent>
          </Select>
        </div>
      </div>

      <div className="flex flex-wrap items-center gap-1.5 pt-1">
        {SOURCES.map((item) => {
          const isSelected = source === item.key;
          const Icon = item.icon;

          return (
            <button
              key={item.key}
              type="button"
              onClick={() => onSourceChange(item.key)}
              className={`inline-flex items-center gap-1.5 rounded-xl border px-3 py-1.5 font-medium text-xs transition-colors ${
                isSelected
                  ? 'border-border bg-card text-foreground shadow-xs'
                  : 'border-transparent text-muted-foreground hover:bg-muted/40 hover:text-foreground'
              }`}
            >
              {Icon && <Icon className="size-3.5" />}
              {item.label}
            </button>
          );
        })}
      </div>
    </div>
  );
}
