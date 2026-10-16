import type { LucideIcon } from 'lucide-react';
import {
  Activity,
  Bot,
  CircleHelp,
  ClipboardList,
  Clock,
  CloudCog,
  CreditCard,
  FolderKanban,
  Globe2,
  HardDrive,
  Key,
  KeyRound,
  LayoutDashboard,
  LayoutGrid,
  Mail,
  PlayCircle,
  Rocket,
  Server,
  Settings,
  ShieldCheck,
  UserRound,
} from 'lucide-react';

export type DashboardNavigationItem = {
  title: string;
  description: string;
  to: string;
  icon: LucideIcon;
  search?: Record<string, string>;
  exact?: boolean;
};

export const primaryNavigation: DashboardNavigationItem[] = [
  {
    title: 'Home',
    description: 'Organization overview',
    to: '/',
    icon: LayoutDashboard,
    exact: true,
  },
  {
    title: 'Projects',
    description: 'Workloads and environments',
    to: '/projects',
    icon: FolderKanban,
  },
  {
    title: 'Apps',
    description: 'Deployed services across projects',
    to: '/apps',
    icon: LayoutGrid,
  },
  {
    title: 'Actions',
    description: 'Continuous integration and automation',
    to: '/actions',
    icon: PlayCircle,
  },
  {
    title: 'Issues',
    description: 'Health, system issues and metrics',
    to: '/monitoring',
    icon: Activity,
  },
];

export const infrastructureNavigation: DashboardNavigationItem[] = [
  {
    title: 'Servers',
    description: 'Deployment targets and runtime capacity',
    to: '/servers',
    icon: Server,
  },
  {
    title: 'Mail',
    description: 'Mail servers, mailboxes and deliverability',
    to: '/mail',
    icon: Mail,
  },
  {
    title: 'Jobs',
    description: 'Cron schedules and job runs',
    to: '/jobs',
    icon: Clock,
  },
  {
    title: 'Backups',
    description: 'Backup destinations and restores',
    to: '/backups',
    icon: HardDrive,
  },
];

export const settingsNavigation: DashboardNavigationItem[] = [
  {
    title: 'Settings',
    description: 'Instance and workspace configuration',
    to: '/settings',
    icon: Settings,
  },
  {
    title: 'Audit',
    description: 'Security and operational events',
    to: '/audit',
    icon: ClipboardList,
  },
];

export const getSettingsNavigation = (isCloud: boolean): DashboardNavigationItem[] => {
  if (!isCloud) {
    return settingsNavigation;
  }
  return [
    settingsNavigation[0],
    {
      title: 'Billing & Usage',
      description: 'Manage subscription and compute plans',
      to: '/billing',
      icon: CreditCard,
    },
    settingsNavigation[1],
  ];
};

export const systemNavigation = settingsNavigation;
export const getSystemNavigation = getSettingsNavigation;

export const helpNavigation: DashboardNavigationItem[] = [
  {
    title: 'Support',
    description: 'Get help and view support tickets',
    to: '/support',
    icon: CircleHelp,
  },
];

export const hiddenNavigation: DashboardNavigationItem[] = [
  {
    title: 'Deployments',
    description: 'Release activity and status',
    to: '/deployments',
    icon: Rocket,
  },
  {
    title: 'Domains',
    description: 'Domain routing and SSL verification',
    to: '/domains',
    icon: Globe2,
  },
];

export const contextualNavigation: DashboardNavigationItem[] = [
  {
    title: 'Credentials',
    description: 'DNS providers, tokens, and registry credentials',
    to: '/settings',
    icon: KeyRound,
    search: { tab: 'credentials' },
  },
  {
    title: 'MCP',
    description: 'Model Context Protocol server configuration',
    to: '/settings',
    icon: Bot,
    search: { tab: 'mcp' },
  },
  {
    title: 'Git',
    description: 'Git providers and registries',
    to: '/settings',
    icon: CloudCog,
    search: { tab: 'git' },
  },
  {
    title: 'Tokens',
    description: 'Personal access tokens',
    to: '/settings',
    icon: Key,
    search: { tab: 'tokens' },
  },
  {
    title: 'Team',
    description: 'Members and permissions',
    to: '/settings',
    icon: UserRound,
    search: { tab: 'team' },
  },
  {
    title: 'Profile & Security',
    description: 'Account settings, password, and 2FA',
    to: '/profile',
    icon: ShieldCheck,
  },
];

export const commandNavigation: DashboardNavigationItem[] = [
  ...primaryNavigation,
  ...infrastructureNavigation,
  ...settingsNavigation,
  ...helpNavigation,
  ...hiddenNavigation,
  ...contextualNavigation,
];
