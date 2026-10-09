export type NotificationGroup =
  | 'deployment'
  | 'app_health'
  | 'backups'
  | 'jobs'
  | 'domains'
  | 'members'
  | 'billing'
  | 'mail';

export interface NotificationEventDef {
  id: string;
  group: NotificationGroup;
  title: string;
  description: string;
  defaultChannels: string[];
  defaultOrgEnabled: boolean;
  defaultNotifyMe: boolean;
}

export const NOTIFICATION_GROUPS: { id: NotificationGroup; label: string }[] = [
  { id: 'deployment', label: 'Deployment' },
  { id: 'app_health', label: 'App health' },
  { id: 'backups', label: 'Backups' },
  { id: 'jobs', label: 'Jobs' },
  { id: 'domains', label: 'Domains & SSL' },
  { id: 'members', label: 'Members' },
  { id: 'billing', label: 'Billing' },
  { id: 'mail', label: 'Mail' },
];

export const NOTIFICATION_EVENTS: NotificationEventDef[] = [
  {
    id: 'deploy.failed',
    group: 'deployment',
    title: 'Deploy failed',
    description:
      'A build or deploy errored out. We send the error message + a snippet of the failing logs.',
    defaultChannels: ['email', 'slack'],
    defaultOrgEnabled: true,
    defaultNotifyMe: true,
  },
  {
    id: 'deploy.succeeded',
    group: 'deployment',
    title: 'Deploy succeeded',
    description: 'Every successful production deploy. Off by default — most teams find it noisy.',
    defaultChannels: ['slack'],
    defaultOrgEnabled: false,
    defaultNotifyMe: false,
  },
  {
    id: 'deploy.cancelled',
    group: 'deployment',
    title: 'Deploy cancelled',
    description: 'An in-flight deploy was cancelled by you or another member.',
    defaultChannels: ['slack'],
    defaultOrgEnabled: false,
    defaultNotifyMe: false,
  },
  {
    id: 'service.down',
    group: 'app_health',
    title: 'App down',
    description:
      'A container that should be running isn\u2019t — exited, dead, or OOM-killed. Includes the exit code and its last log lines.',
    defaultChannels: ['email', 'slack', 'webhook'],
    defaultOrgEnabled: true,
    defaultNotifyMe: true,
  },
  {
    id: 'service.crash_looping',
    group: 'app_health',
    title: 'App crash looping',
    description:
      'A container keeps restarting. Point-in-time status always reads healthy for these, so this is the only way you\u2019d know.',
    defaultChannels: ['email', 'slack'],
    defaultOrgEnabled: true,
    defaultNotifyMe: true,
  },
  {
    id: 'service.unhealthy',
    group: 'app_health',
    title: 'App unhealthy',
    description: 'Container failed consecutive HTTP or TCP readiness health checks.',
    defaultChannels: ['slack'],
    defaultOrgEnabled: true,
    defaultNotifyMe: false,
  },
  {
    id: 'service.high_cpu',
    group: 'app_health',
    title: 'High CPU usage',
    description: 'Container sustained over 90% CPU utilization for more than 5 minutes.',
    defaultChannels: ['slack'],
    defaultOrgEnabled: false,
    defaultNotifyMe: false,
  },
  {
    id: 'service.high_memory',
    group: 'app_health',
    title: 'High memory usage',
    description:
      'Container reached configured memory warning limit and is at risk of OOM termination.',
    defaultChannels: ['email', 'slack'],
    defaultOrgEnabled: true,
    defaultNotifyMe: false,
  },
  {
    id: 'backup.failed',
    group: 'backups',
    title: 'Backup failed',
    description: 'Scheduled or manual database/volume snapshot failed to complete.',
    defaultChannels: ['email', 'slack'],
    defaultOrgEnabled: true,
    defaultNotifyMe: true,
  },
  {
    id: 'backup.completed',
    group: 'backups',
    title: 'Backup completed',
    description: 'Snapshot successfully created and uploaded to the backup destination.',
    defaultChannels: ['slack'],
    defaultOrgEnabled: false,
    defaultNotifyMe: false,
  },
  {
    id: 'backup.pruned',
    group: 'backups',
    title: 'Backup pruned',
    description: 'Older backup retention policy executed and retired aged archive files.',
    defaultChannels: ['slack'],
    defaultOrgEnabled: false,
    defaultNotifyMe: false,
  },
  {
    id: 'job.failed',
    group: 'jobs',
    title: 'Job failed',
    description: 'Scheduled cron execution exited with a non-zero exit code.',
    defaultChannels: ['email', 'slack'],
    defaultOrgEnabled: true,
    defaultNotifyMe: true,
  },
  {
    id: 'job.completed',
    group: 'jobs',
    title: 'Job completed',
    description: 'Scheduled job execution finished successfully.',
    defaultChannels: ['slack'],
    defaultOrgEnabled: false,
    defaultNotifyMe: false,
  },
  {
    id: 'job.timeout',
    group: 'jobs',
    title: 'Job timeout',
    description: 'Scheduled task exceeded configured maximum run timeout.',
    defaultChannels: ['email', 'slack'],
    defaultOrgEnabled: true,
    defaultNotifyMe: true,
  },
  {
    id: 'domain.cert_expiring',
    group: 'domains',
    title: 'Certificate expiring',
    description:
      'Automated TLS certificate renewal failed and certificate is within 7 days of expiration.',
    defaultChannels: ['email'],
    defaultOrgEnabled: true,
    defaultNotifyMe: true,
  },
  {
    id: 'domain.dns_failed',
    group: 'domains',
    title: 'DNS resolution failed',
    description: 'Acme challenge or domain CNAME verification could not be resolved.',
    defaultChannels: ['email', 'slack'],
    defaultOrgEnabled: true,
    defaultNotifyMe: false,
  },
  {
    id: 'member.joined',
    group: 'members',
    title: 'Member joined',
    description: 'A new user accepted an invite and joined this workspace.',
    defaultChannels: ['email'],
    defaultOrgEnabled: true,
    defaultNotifyMe: false,
  },
  {
    id: 'member.removed',
    group: 'members',
    title: 'Member removed',
    description: 'A teammate was removed or left the organization.',
    defaultChannels: ['email'],
    defaultOrgEnabled: false,
    defaultNotifyMe: false,
  },
  {
    id: 'member.role_changed',
    group: 'members',
    title: 'Role changed',
    description: 'A member was granted elevated or updated workspace permissions.',
    defaultChannels: ['email'],
    defaultOrgEnabled: true,
    defaultNotifyMe: false,
  },
  {
    id: 'billing.limit_reached',
    group: 'billing',
    title: 'Usage limit reached',
    description: 'Organization reached project, resource, or bandwidth allocation ceiling.',
    defaultChannels: ['email'],
    defaultOrgEnabled: true,
    defaultNotifyMe: true,
  },
  {
    id: 'billing.payment_failed',
    group: 'billing',
    title: 'Invoice failed',
    description: 'Subscription payment method failed processing.',
    defaultChannels: ['email'],
    defaultOrgEnabled: true,
    defaultNotifyMe: true,
  },
];
