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

export interface AuditLog {
  id: string;
  organizationId?: string;
  userId: string;
  action: string;
  resource: string;
  resourceName?: string;
  resourceType?: string;
  details: string;
  ipAddress: string;
  createdAt: string;
  actorName?: string;
  actorEmail?: string;
  actor?: { id?: string; name?: string; email?: string };
}

export type AuditSource = 'dashboard' | 'mcp' | 'cli' | 'api' | 'webhook' | 'system';

export interface AuditLogRow extends AuditLog {
  eventType?: string;
  resourceType?: string;
  resourceId?: string;
  resourceName?: string;
  actorUserId?: string;
  actor?: { id?: string; name?: string; email?: string };
  category?: string;
  source?: AuditSource;
  sourceClientName?: string;
  sourceClientId?: string;
  actorName?: string;
  actorEmail?: string;
  userAgent?: string;
  before?: unknown;
  after?: unknown;
}

export type AuditEventRow = AuditLogRow;

export interface AuditCategoryFacet {
  id: string;
  label: string;
  description: string;
  count: number;
}

export interface AuditFacets {
  total: number;
  categories: AuditCategoryFacet[];
  actors?: Array<{ id: string; name?: string; email?: string }>;
  sources?: Array<{ source: string; count: number }>;
  settings?: AuditSettings;
  canManage?: boolean;
}

export interface AuditListParams {
  category?: string;
  source?: string;
  actor?: string;
  period?: string;
  q?: string;
  limit?: number;
  offset?: number;
  page?: number;
  perPage?: number;
}

export interface AuditSettings {
  organizationId: string;
  enabled: boolean;
  retentionDays: number;
  updatedAt?: string;
}
