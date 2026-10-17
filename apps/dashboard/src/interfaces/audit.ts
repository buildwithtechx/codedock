export interface AuditLog {
  id: string;
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
  actor?: { id: string; name?: string; email?: string };
}
