export type AttentionSeverity = 'info' | 'warning' | 'critical';
export type AttentionStatus = 'open' | 'acked' | 'resolved';
export type DisplaySeverity = 'outage' | 'action_required' | 'advisory';
export type MonitoringTab = 'open' | 'resolved' | 'health';
export type SeverityFilter = 'all' | DisplaySeverity;

export interface AttentionIssue {
  id: string;
  organizationId: string;
  kind: string;
  subject: string;
  projectId?: string;
  serviceId?: string;
  severity: AttentionSeverity;
  status: AttentionStatus;
  title: string;
  detail: string;
  remediation: string;
  action?: string;
  occurrences: number;
  firstSeen: string;
  lastSeen: string;
  updatedAt: string;
}

export interface ActIssueInput {
  issueId: string;
  action: string;
}

export interface ActIssuePayload {
  action: string;
  params?: Record<string, string>;
}

export interface ActIssueOutcome {
  outcome: string;
}
