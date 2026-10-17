import { Bot, Globe, Layers, LayoutGrid, Rocket, Server, Shield, Users } from 'lucide-react';
import type { AuditEventRow } from './api';
import type { PeriodKey } from './audit-filter-bar';
import type { AuditCategoryId } from './audit-taxonomy';

export type CategoryKey = 'all' | AuditCategoryId;

export const CATEGORY_ICONS: Record<AuditCategoryId, typeof Rocket> = {
  deployments: Rocket,
  apps: LayoutGrid,
  domains: Globe,
  servers: Server,
  members: Users,
  agent: Bot,
  security: Shield,
  billing: Layers,
  system: Layers,
};

export function periodStart(period: PeriodKey): string | undefined {
  if (period === 'all') return undefined;
  const start = new Date();
  start.setHours(0, 0, 0, 0);
  const back = period === 'today' ? 0 : period === '7d' ? 6 : period === '30d' ? 29 : 89;
  start.setDate(start.getDate() - back);
  return start.toISOString();
}

export function groupByDay(events: AuditEventRow[]): Record<string, AuditEventRow[]> {
  const groups: Record<string, AuditEventRow[]> = {};
  for (const event of events) {
    const d = new Date(event.createdAt);
    const key = `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`;
    if (!groups[key]) groups[key] = [];
    groups[key].push(event);
  }
  return groups;
}
