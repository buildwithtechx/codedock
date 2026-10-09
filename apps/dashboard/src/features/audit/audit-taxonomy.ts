import type { AuditLog } from '#/interfaces/audit';

export type AuditCategoryId =
  | 'deployments'
  | 'apps'
  | 'domains'
  | 'servers'
  | 'members'
  | 'agent'
  | 'security'
  | 'billing'
  | 'system';

export type AuditTone = 'neutral' | 'info' | 'success' | 'warning' | 'danger';

export interface AuditCategory {
  id: AuditCategoryId;
  label: string;
  description: string;
}

export const AUDIT_CATEGORIES: AuditCategory[] = [
  { id: 'deployments', label: 'Deployments', description: 'Builds and releases.' },
  { id: 'apps', label: 'Apps & services', description: 'Projects, services, and runtime health.' },
  { id: 'domains', label: 'Domains & SSL', description: 'Custom domains, DNS, and certificates.' },
  { id: 'servers', label: 'Servers', description: 'Deployment targets and shell access.' },
  {
    id: 'members',
    label: 'Members & access',
    description: 'Who is in this workspace and what they can do.',
  },
  { id: 'agent', label: 'AI agents', description: 'What connected assistants did over MCP.' },
  {
    id: 'security',
    label: 'Security',
    description: 'Credentials, auth, exports, and audit recording.',
  },
  { id: 'billing', label: 'Billing', description: 'Subscriptions, credits, and quotas.' },
  {
    id: 'system',
    label: 'System',
    description: 'Settings, integrations, backups, and maintenance.',
  },
];

const CATEGORY_IDS = new Set<string>(AUDIT_CATEGORIES.map((category) => category.id));

export function isAuditCategoryId(value: unknown): value is AuditCategoryId {
  return typeof value === 'string' && CATEGORY_IDS.has(value);
}

export function categoryLabel(id: string): string {
  return AUDIT_CATEGORIES.find((category) => category.id === id)?.label ?? 'System';
}

const CATEGORY_BY_PREFIX: Record<string, AuditCategoryId> = {
  deployment: 'deployments',
  deploy: 'deployments',
  release: 'deployments',
  build: 'deployments',
  app: 'apps',
  project: 'apps',
  service: 'apps',
  database: 'apps',
  container: 'apps',
  job: 'apps',
  scheduled: 'apps',
  cron: 'apps',
  domain: 'domains',
  dns: 'domains',
  ssl: 'domains',
  certificate: 'domains',
  server: 'servers',
  ssh: 'servers',
  host: 'servers',
  member: 'members',
  user: 'members',
  team: 'members',
  invitation: 'members',
  invite: 'members',
  role: 'members',
  agent: 'agent',
  mcp: 'agent',
  tool: 'agent',
  auth: 'security',
  login: 'security',
  token: 'security',
  apikey: 'security',
  credential: 'security',
  secret: 'security',
  audit: 'security',
  billing: 'billing',
  subscription: 'billing',
  invoice: 'billing',
  credit: 'billing',
  quota: 'billing',
  backup: 'system',
  restore: 'system',
  organization: 'system',
  settings: 'system',
  notification: 'system',
  integration: 'system',
  webhook: 'system',
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
  join: 'joined',
  leave: 'left',
  revoke: 'revoked',
  reveal: 'revealed credentials for',
  query: 'queried',
  insert: 'added a row to',
  backup: 'backed up',
  restore: 'restored',
  setup: 'set up',
  upgrade: 'upgraded',
  scale: 'scaled',
  verify: 'verified',
  issue: 'issued',
  renew: 'renewed',
  suspend: 'suspended',
  resume: 'resumed',
  rotate: 'rotated',
  sync: 'synced',
  cancel: 'cancelled',
  cancelled: 'cancelled',
  fail: 'failed',
  failed: 'failed',
  succeed: 'succeeded',
  succeeded: 'succeeded',
  complete: 'completed',
  completed: 'completed',
  approve: 'approved',
  transfer: 'transferred',
  attach: 'attached',
  detach: 'detached',
  link: 'linked',
  unlink: 'unlinked',
  install: 'installed',
  execute: 'executed',
  run: 'ran',
  grant: 'granted',
  deny: 'denied',
};

const DANGER_VERBS = new Set([
  'delete',
  'remove',
  'revoke',
  'rollback',
  'disable',
  'fail',
  'failed',
  'cancel',
  'cancelled',
  'suspend',
  'deny',
]);
const SUCCESS_VERBS = new Set([
  'create',
  'deploy',
  'redeploy',
  'enable',
  'trigger',
  'backup',
  'invite',
  'succeed',
  'succeeded',
  'complete',
  'completed',
  'verify',
  'approve',
  'renew',
]);
const WARNING_VERBS = new Set([
  'reveal',
  'login',
  'logout',
  'rollback',
  'credentials_reveal',
  'rotate',
  'suspend',
]);

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

export function auditCategoryOf(log: {
  action: string;
  resource: string;
  category?: string;
}): AuditCategoryId {
  if (log.category && isAuditCategoryId(log.category)) return log.category;
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
