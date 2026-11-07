import {
  AUDIT_CATEGORIES,
  AUDIT_CATEGORY_BY_PREFIX,
  AUDIT_VERBS,
  DANGER_VERBS,
  SUCCESS_VERBS,
  WARNING_VERBS,
} from '#/constants/audit';
import type { AuditCategoryId, AuditLog, AuditTone } from '#/interfaces/audit';
import { resolveResourceDisplay } from './audit-resource-helper';

const CATEGORY_IDS = new Set<string>(AUDIT_CATEGORIES.map((category) => category.id));

export function isAuditCategoryId(value: unknown): value is AuditCategoryId {
  return typeof value === 'string' && CATEGORY_IDS.has(value);
}

export function categoryLabel(id: string): string {
  return AUDIT_CATEGORIES.find((category) => category.id === id)?.label ?? 'System';
}

export function actionVerb(action?: string): string {
  if (!action || typeof action !== 'string') return 'event';
  return action.split('.').pop() ?? action;
}

export function describeAuditAction(action?: string): {
  label: string;
  action: string;
  tone: AuditTone;
} {
  const safe = action && typeof action === 'string' ? action : 'system.event';
  const verb = actionVerb(safe).toLowerCase();
  const label =
    AUDIT_VERBS[verb] ??
    (verb.endsWith('e')
      ? `${verb}d`
      : verb.endsWith('y') && !/[aeiou]y$/.test(verb)
        ? `${verb.slice(0, -1)}ied`
        : `${verb}ed`);
  let tone: AuditTone = 'neutral';
  if (DANGER_VERBS.has(verb)) tone = 'danger';
  else if (SUCCESS_VERBS.has(verb)) tone = 'success';
  else if (WARNING_VERBS.has(verb) || verb.includes('reveal') || verb.includes('secret'))
    tone = 'warning';
  else if (verb === 'update' || verb === 'scale' || verb === 'upgrade') tone = 'info';
  return { label: safe, action: label, tone };
}

export const describeAuditEvent = describeAuditAction;

export function auditCategoryOf(log?: {
  action?: string;
  resource?: string;
  category?: string;
}): AuditCategoryId {
  if (!log) return 'system';
  if (log.category && isAuditCategoryId(log.category)) return log.category;
  const action = log.action && typeof log.action === 'string' ? log.action : '';
  const first = (actionVerb(action) === action ? action : action.split('.')[0])
    .toLowerCase()
    .replace(/_.*$/, '');
  const hit = AUDIT_CATEGORY_BY_PREFIX[first];
  if (hit) return hit;
  const resource =
    log.resource && typeof log.resource === 'string' ? log.resource.toLowerCase() : '';
  for (const [prefix, category] of Object.entries(AUDIT_CATEGORY_BY_PREFIX)) {
    if (resource.includes(prefix)) return category;
  }
  return 'system';
}

export function auditResourceLabel(resourceType?: string): string {
  if (!resourceType || typeof resourceType !== 'string') return 'resource';
  const cleaned = resourceType.replace(/_/g, ' ').trim();
  return cleaned || 'resource';
}

export function toneDot(tone: AuditTone): string {
  if (tone === 'success') return 'bg-success';
  if (tone === 'danger') return 'bg-destructive';
  if (tone === 'warning') return 'bg-warning';
  if (tone === 'info') return 'bg-primary';
  return 'bg-muted-foreground/40';
}

export function actorDisplay(log?: Partial<AuditLog>): string {
  if (log?.actor?.name?.trim()) return log.actor.name.trim();
  if (log?.actorName?.trim()) return log.actorName.trim();
  if (log?.actor?.email?.trim()) return log.actor.email.trim();
  if (log?.actorEmail?.trim()) return log.actorEmail.trim();
  if (log?.userId && typeof log.userId === 'string' && log.userId.trim()) {
    if (log.userId.toLowerCase() === 'system') return 'System';
    return `User ${log.userId.trim().slice(0, 8)}`;
  }
  return 'System';
}

export function resourceDisplay(log?: Partial<AuditLog>): string {
  return resolveResourceDisplay(log);
}
