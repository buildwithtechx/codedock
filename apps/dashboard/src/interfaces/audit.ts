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
