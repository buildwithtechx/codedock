export type AnalyticsSummary = {
  projectId: string;
  from: string;
  to: string;
  requests: number;
  bytes: number;
  errorRate: number;
  avgDurationMs: number;
};
export type AnalyticsOverview = {
  projectId: string;
  from: string;
  to: string;
  series: { time: string; requests: number; bytes: number; errors: number }[];
  statuses: { status: number; requests: number }[];
  topPaths: { path: string; requests: number; bytes: number }[];
};
export type AnalyticsGeo = {
  projectId: string;
  from: string;
  to: string;
  countries: { country: string; requests: number; visitors: number; bytes: number }[];
};
export type AnalyticsDeploymentStats = {
  projectId: string;
  total: number;
  succeeded: number;
  failed: number;
  successRate: number;
  avgDurationSeconds: number;
};
export type AnalyticsDashboard = {
  organizationId: string;
  requests24h: number;
  errorRate24h: number;
  deployments7d: number;
  deploySuccess: number;
  openIssues: number;
};
export type ServiceUsage = {
  serviceId: string;
  name: string;
  status: string;
  cpuPercent: number;
  memoryBytes: number;
  memoryLimit: number;
  uptimeSeconds: number;
};
export type AttentionIssue = {
  id: string;
  organizationId: string;
  kind: string;
  subject: string;
  projectId: string;
  serviceId: string;
  severity: string;
  status: string;
  title: string;
  detail: string;
  remediation: string;
  action: string;
  occurrences: number;
  firstSeen: string;
  lastSeen: string;
  updatedAt: string;
};
