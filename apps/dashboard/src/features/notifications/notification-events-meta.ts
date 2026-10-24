export type NotificationGroup =
  | 'deployment'
  | 'app_health'
  | 'backups'
  | 'jobs'
  | 'actions'
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
  { id: 'actions', label: 'Actions' },
  { id: 'domains', label: 'Domains & SSL' },
  { id: 'members', label: 'Members' },
  { id: 'billing', label: 'Billing' },
];

export const NOTIFICATION_EVENTS: NotificationEventDef[] = [
  // Deployment (3)
  {
    id: 'deploy.failed',
    group: 'deployment',
    title: 'Deploy failed',
    description:
      'A build or deploy errored out. We send the error message + a snippet of the failing logs.',
    defaultChannels: ['email', 'slack', 'webhook'],
    defaultOrgEnabled: true,
    defaultNotifyMe: false,
  },
  {
    id: 'deploy.succeeded',
    group: 'deployment',
    title: 'Deploy succeeded',
    description: 'Every successful production deploy. Off by default — most teams find it noisy.',
    defaultChannels: ['email', 'slack', 'webhook', 'telegram'],
    defaultOrgEnabled: false,
    defaultNotifyMe: false,
  },
  {
    id: 'deploy.cancelled',
    group: 'deployment',
    title: 'Deploy cancelled',
    description: 'An in-flight deploy was cancelled by you or another member.',
    defaultChannels: ['email'],
    defaultOrgEnabled: false,
    defaultNotifyMe: false,
  },

  // App health (5)
  {
    id: 'service.down',
    group: 'app_health',
    title: 'App down',
    description:
      'A container that should be running isn\u2019t — exited, dead, or OOM-killed. Includes the exit code and its last log lines. Never fires for a service you stopped yourself.',
    defaultChannels: ['email'],
    defaultOrgEnabled: true,
    defaultNotifyMe: false,
  },
  {
    id: 'service.crash_loop',
    group: 'app_health',
    title: 'App crash looping',
    description:
      'A container keeps restarting. Point-in-time status always reads healthy for these, so this is the only way you\u2019d know.',
    defaultChannels: ['email'],
    defaultOrgEnabled: true,
    defaultNotifyMe: false,
  },
  {
    id: 'service.unhealthy',
    group: 'app_health',
    title: 'App unhealthy',
    description: 'A container is up but its own healthcheck is failing.',
    defaultChannels: ['email'],
    defaultOrgEnabled: true,
    defaultNotifyMe: false,
  },
  {
    id: 'service.recovered',
    group: 'app_health',
    title: 'App recovered',
    description: 'An outage ended and the container is healthy again.',
    defaultChannels: ['email'],
    defaultOrgEnabled: true,
    defaultNotifyMe: false,
  },
  {
    id: 'server.unreachable',
    group: 'app_health',
    title: 'Server unreachable',
    description:
      'We can\u2019t reach a server\u2019s Docker daemon. One alert per server — never one per service on it.',
    defaultChannels: ['email'],
    defaultOrgEnabled: true,
    defaultNotifyMe: false,
  },

  // Backups (3)
  {
    id: 'backup.failed',
    group: 'backups',
    title: 'Backup failed',
    description:
      'A scheduled or manual backup run errored out. Includes the destination + policy that failed.',
    defaultChannels: ['email'],
    defaultOrgEnabled: true,
    defaultNotifyMe: false,
  },
  {
    id: 'backup.succeeded',
    group: 'backups',
    title: 'Backup succeeded',
    description: 'Each successful backup. Off by default — high volume.',
    defaultChannels: ['email'],
    defaultOrgEnabled: false,
    defaultNotifyMe: false,
  },
  {
    id: 'backup.restore_completed',
    group: 'backups',
    title: 'Restore completed',
    description: 'A restore finished (or failed). Always worth notifying on.',
    defaultChannels: ['email'],
    defaultOrgEnabled: true,
    defaultNotifyMe: false,
  },

  // Jobs (3)
  {
    id: 'job.run.failed',
    group: 'jobs',
    title: 'Job failed',
    description: 'A scheduled or manual job run errored out. Includes the job + exit code.',
    defaultChannels: ['email'],
    defaultOrgEnabled: true,
    defaultNotifyMe: false,
  },
  {
    id: 'job.run.succeeded',
    group: 'jobs',
    title: 'Job succeeded',
    description: 'Each successful job run. Off by default — can be high volume.',
    defaultChannels: ['email'],
    defaultOrgEnabled: false,
    defaultNotifyMe: false,
  },
  {
    id: 'job.run.started',
    group: 'jobs',
    title: 'Job started',
    description: 'A job began running. Off by default — mostly useful per-job.',
    defaultChannels: ['email'],
    defaultOrgEnabled: false,
    defaultNotifyMe: false,
  },

  // Actions (3)
  {
    id: 'action.run.failed',
    group: 'actions',
    title: 'Workflow failed',
    description:
      'An Actions workflow failed or timed out. Open the run to review its checks and logs.',
    defaultChannels: ['email'],
    defaultOrgEnabled: false,
    defaultNotifyMe: false,
  },
  {
    id: 'action.run.succeeded',
    group: 'actions',
    title: 'Workflow succeeded',
    description: 'An Actions workflow completed successfully.',
    defaultChannels: ['email'],
    defaultOrgEnabled: false,
    defaultNotifyMe: false,
  },
  {
    id: 'action.run.cancelled',
    group: 'actions',
    title: 'Workflow cancelled',
    description: 'An Actions workflow was cancelled.',
    defaultChannels: ['email'],
    defaultOrgEnabled: false,
    defaultNotifyMe: false,
  },

  // Domains & SSL (2)
  {
    id: 'domain.expiring',
    group: 'domains',
    title: 'SSL cert expiring',
    description:
      'Daily check — fires when any of your custom domain certs has under 7 days left and the renewer can\u2019t reach the host.',
    defaultChannels: ['email'],
    defaultOrgEnabled: true,
    defaultNotifyMe: false,
  },
  {
    id: 'domain.verification_failed',
    group: 'domains',
    title: 'Domain verification failed',
    description: 'DNS check didn\u2019t pass during initial domain setup.',
    defaultChannels: ['email'],
    defaultOrgEnabled: true,
    defaultNotifyMe: false,
  },

  // Members (3)
  {
    id: 'member.added',
    group: 'members',
    title: 'New member joined',
    description: 'Someone accepted an invite and joined this organization.',
    defaultChannels: ['email'],
    defaultOrgEnabled: false,
    defaultNotifyMe: false,
  },
  {
    id: 'member.removed',
    group: 'members',
    title: 'Member removed',
    description: 'A member was removed from the organization.',
    defaultChannels: ['email'],
    defaultOrgEnabled: false,
    defaultNotifyMe: false,
  },
  {
    id: 'invitation.sent',
    group: 'members',
    title: 'Invitations sent',
    description: 'An admin sent invite(s) to new members. Useful as a security signal.',
    defaultChannels: ['email'],
    defaultOrgEnabled: false,
    defaultNotifyMe: false,
  },

  // Billing (2)
  {
    id: 'billing.alert',
    group: 'billing',
    title: 'Billing alert',
    description:
      'Payment failed, plan limit reached, or invoice overdue. Always notifies billing-owners.',
    defaultChannels: ['email'],
    defaultOrgEnabled: true,
    defaultNotifyMe: false,
  },
  {
    id: 'quota.warning',
    group: 'billing',
    title: 'Quota warning',
    description:
      'You\u2019re approaching a plan limit (storage, deployments, members). Heads-up before things start failing.',
    defaultChannels: ['email'],
    defaultOrgEnabled: true,
    defaultNotifyMe: false,
  },
];
