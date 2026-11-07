import type { AuditLog } from '#/interfaces/audit';

const UUID_REGEX = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i;

export function isUUID(value?: string | null): boolean {
  return typeof value === 'string' && UUID_REGEX.test(value.trim());
}

export function extractNameFromDetails(details?: unknown): string | null {
  if (!details) return null;
  let parsed: Record<string, unknown> | null = null;
  if (typeof details === 'string') {
    const trimmed = details.trim();
    if (!trimmed.startsWith('{')) return null;
    try {
      parsed = JSON.parse(trimmed);
    } catch {
      return null;
    }
  } else if (typeof details === 'object') {
    parsed = details as Record<string, unknown>;
  }
  if (!parsed) return null;

  const candidateKeys = [
    'name',
    'appName',
    'projectName',
    'title',
    'domain',
    'hostname',
    'serverName',
    'slug',
  ];
  for (const key of candidateKeys) {
    const val = parsed[key];
    if (typeof val === 'string' && val.trim() && !isUUID(val)) {
      return val.trim();
    }
  }

  for (const sub of ['after', 'before']) {
    const subObj = parsed[sub];
    if (typeof subObj === 'object' && subObj !== null) {
      for (const key of candidateKeys) {
        const val = (subObj as Record<string, unknown>)[key];
        if (typeof val === 'string' && val.trim() && !isUUID(val)) {
          return val.trim();
        }
      }
    }
  }
  return null;
}

export function fallbackResourceNoun(actionName?: string, resourceType?: string): string {
  if (resourceType && resourceType.trim() && !isUUID(resourceType)) {
    return `the ${resourceType.trim().replace(/_/g, ' ')}`;
  }
  const prefix = (actionName || '').split('.')[0].toLowerCase();
  switch (prefix) {
    case 'project':
      return 'the project';
    case 'app':
    case 'service':
      return 'the application';
    case 'database':
      return 'the database';
    case 'server':
      return 'the server';
    case 'deployment':
    case 'deploy':
    case 'release':
      return 'the deployment';
    case 'domain':
    case 'ssl':
    case 'certificate':
      return 'the domain';
    case 'member':
    case 'user':
    case 'invite':
    case 'invitation':
      return 'the member';
    case 'organization':
    case 'org':
      return 'the organization';
    case 'billing':
    case 'subscription':
      return 'the subscription';
    case 'backup':
      return 'the backup';
    default:
      return 'the resource';
  }
}

export function resolveResourceDisplay(
  log?: (Partial<AuditLog> & { after?: unknown; before?: unknown }) | null,
  actionName?: string
): string {
  if (!log) return '';

  if (log.resourceName && typeof log.resourceName === 'string' && log.resourceName.trim()) {
    const trimmed = log.resourceName.trim();
    if (!isUUID(trimmed)) return trimmed;
  }

  const fromDetails =
    extractNameFromDetails(log.details) ||
    extractNameFromDetails(log.after) ||
    extractNameFromDetails(log.before);
  if (fromDetails) return fromDetails;

  if (log.resource && typeof log.resource === 'string' && log.resource.trim()) {
    const trimmed = log.resource.trim();
    if (!isUUID(trimmed)) return trimmed;
  }

  const act = actionName || log.action || '';
  return fallbackResourceNoun(act, log.resourceType);
}
