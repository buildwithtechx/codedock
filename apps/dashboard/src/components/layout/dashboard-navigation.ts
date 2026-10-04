import type { LucideIcon } from 'lucide-react';
import {
  Activity,
  ClipboardList,
  Clock,
  CloudCog,
  FolderKanban,
  Globe2,
  HardDrive,
  Key,
  LayoutDashboard,
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
    title: 'Deployments',
    description: 'Release activity and status',
    to: '/deployments',
    icon: Rocket,
  },
  {
    title: 'Monitoring',
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
    title: 'Domains & DNS',
    description: 'Domain routing and SSL verification',
    to: '/dns',
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
    exact: true,
  },
  {
    title: 'Audit Logs',
    description: 'Security and operational events',
    to: '/audit',
    icon: ClipboardList,
  },
  {
    title: 'API Access',
    description: 'Personal access tokens',
    to: '/api-access',
    icon: Key,
  },
  {
    title: 'Settings',
    description: 'Instance and workspace configuration',
    to: '/settings',
    icon: Settings,
    exact: true,
  },
];

export const contextualNavigation: DashboardNavigationItem[] = [
  {
    title: 'Sources',
    description: 'Git providers and registries',
    to: '/settings',
    icon: CloudCog,
    search: { tab: 'sources' },
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
];
