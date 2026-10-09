import type { AuditLog } from '#/interfaces/audit';

export type AuditCategoryId =
  | 'deployments'
  | 'projects'
  | 'services'
  | 'databases'
  | 'domains'
  | 'members'
  | 'security'
  | 'backups'
  | 'jobs'
  | 'servers'
  | 'system';

export type AuditTone = 'neutral' | 'info' | 'success' | 'warning' | 'danger';

export interface AuditCategory {
  id: AuditCategoryId;
  label: string;
}

export const AUDIT_CATEGORIES: AuditCategory[] = [
  { id: 'deployments', label: 'Deployments' },
  { id: 'projects', label: 'Projects' },
  { id: 'services', label: 'Services' },
  { id: 'databases', label: 'Databases' },
  { id: 'domains', label: 'Domains' },
  { id: 'members', label: 'Members' },
  { id: 'security', label: 'Security' },
  { id: 'backups', label: 'Backups' },
  { id: 'jobs', label: 'Jobs' },
  { id: 'servers', label: 'Servers' },
  { id: 'system', label: 'System' },
];

const CATEGORY_BY_PREFIX: Record<string, AuditCategoryId> = {
  deployment: 'deployments',
  project: 'projects',
  service: 'services',
  database: 'databases',
  domain: 'domains',
  user: 'members',
  member: 'members',
  team: 'members',
  invitation: 'members',
  auth: 'security',
  login: 'security',
  token: 'security',
  apikey: 'security',
  backup: 'backups',
  restore: 'backups',
  job: 'jobs',
  scheduled: 'jobs',
  server: 'servers',
  billing: 'system',
  organization: 'system',
  settings: 'system',
};

const VERBS: Record<string, string> = {
  create: 'created',
  update: 'updated',
  delete: 'deleted',
  remove: 'removed',
  trigger: 'triggered',
  rollback: 'rolled back',
  deploy: 'deployed',
  redeploy: 'redeployed',
  restart: 'restarted',
  stop: 'stopped',
  start: 'started',
  enable: 'enabled',
  disable: 'disabled',
  login: 'signed in',
  logout: 'signed out',
  invite: 'invited',
  revoke: 'revoked',
  reveal: 'revealed credentials for',
  query: 'queried',
  insert: 'added a row to',
  backup: 'backed up',
  restore: 'restored',
  setup: 'set up',
  upgrade: 'upgraded',
  scale: 'scaled',
};

const DANGER_VERBS = new Set(['delete', 'remove', 'revoke', 'rollback', 'disable']);
const SUCCESS_VERBS = new Set([
  'create',
  'deploy',
  'redeploy',
  'enable',
  'trigger',
  'backup',
  'invite',
]);
const WARNING_VERBS = new Set(['reveal', 'login', 'logout', 'rollback', 'credentials_reveal']);

export function actionVerb(action: string): string {
  return action.split('.').pop() ?? action;
}

export function describeAuditAction(action: string): {
  label: string;
  action: string;
  tone: AuditTone;
} {
  const verb = actionVerb(action).toLowerCase();
  const label = VERBS[verb] ?? `${verb}d`;
  let tone: AuditTone = 'neutral';
  if (DANGER_VERBS.has(verb)) tone = 'danger';
  else if (SUCCESS_VERBS.has(verb)) tone = 'success';
  else if (WARNING_VERBS.has(verb) || verb.includes('reveal') || verb.includes('secret'))
    tone = 'warning';
  else if (verb === 'update' || verb === 'scale' || verb === 'upgrade') tone = 'info';
  return { label: action, action: label, tone };
}

export function auditCategoryOf(log: AuditLog): AuditCategoryId {
  const first = (actionVerb(log.action) === log.action ? log.action : log.action.split('.')[0])
    .toLowerCase()
    .replace(/_.*$/, '');
  const hit = CATEGORY_BY_PREFIX[first];
  if (hit) return hit;
  const resource = log.resource.toLowerCase();
  for (const [prefix, category] of Object.entries(CATEGORY_BY_PREFIX)) {
    if (resource.includes(prefix)) return category;
  }
  return 'system';
}

export function auditResourceLabel(resourceType: string): string {
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

export function actorDisplay(log: AuditLog): string {
  if (log.userId?.trim()) return log.userId;
  return 'System';
}

export function resourceDisplay(log: AuditLog): string {
  if (log.resource?.trim()) return log.resource;
  if (log.details?.trim()) return log.details.slice(0, 80);
  return 'resource';
}

export function relativeTime(iso: string): string {
  const date = new Date(iso);
  if (Number.isNaN(date.getTime())) return iso;
  const diff = Date.now() - date.getTime();
  const min = Math.floor(diff / 60_000);
  if (min < 1) return 'just now';
  if (min < 60) return `${min}m ago`;
  const hr = Math.floor(min / 60);
  if (hr < 24) return `${hr}h ago`;
  const days = Math.floor(hr / 24);
  if (days < 7) return `${days}d ago`;
  return date.toLocaleDateString();
}

export function absoluteTime(iso: string): string {
  const date = new Date(iso);
  if (Number.isNaN(date.getTime())) return iso;
  return date.toLocaleString();
}

export function clockTime(iso: string): string {
  const date = new Date(iso);
  if (Number.isNaN(date.getTime())) return '';
  return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
}

export function dayKey(iso: string): string {
  const date = new Date(iso);
  return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`;
}
