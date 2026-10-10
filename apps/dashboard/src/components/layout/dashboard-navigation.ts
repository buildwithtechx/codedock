import type { LucideIcon } from 'lucide-react';
import {
  Activity,
  Bot,
  ClipboardList,
  Clock,
  CloudCog,
  FolderKanban,
  Globe2,
  HardDrive,
  Key,
  KeyRound,
  LayoutDashboard,
  LayoutGrid,
  Rocket,
  Server,
  Settings,
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
    title: 'Deployments',
    description: 'Release activity and status',
    to: '/deployments',
    icon: Rocket,
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
    title: 'Domains',
    description: 'Domain routing and SSL verification',
    to: '/domains',
    icon: Globe2,
  },
  {
    title: 'Jobs',
    description: 'Cron schedules and job runs',
    to: '/jobs',
    icon: Clock,
  },
];

export const systemNavigation: DashboardNavigationItem[] = [
  {
    title: 'Backups',
    description: 'Backup destinations and restores',
    to: '/backups',
    icon: HardDrive,
  },
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

export const hiddenNavigation: DashboardNavigationItem[] = [];

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
];

export const commandNavigation: DashboardNavigationItem[] = [
  ...primaryNavigation,
  ...infrastructureNavigation,
  ...systemNavigation,
  ...hiddenNavigation,
  ...contextualNavigation,
];
