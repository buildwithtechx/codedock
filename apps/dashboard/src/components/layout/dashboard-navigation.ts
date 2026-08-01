import type { LucideIcon } from 'lucide-react';
import {
  CalendarDays,
  ClipboardList,
  CloudCog,
  FolderKanban,
  Globe2,
  HardDrive,
  Key,
  LayoutDashboard,
  Rocket,
  Server,
  Settings,
  Sparkles,
  UserRound,
} from 'lucide-react';

export type DashboardNavigationItem = {
  title: string;
  description: string;
  to: string;
  icon: LucideIcon;
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
    description: 'Services and environments',
    to: '/projects',
    icon: FolderKanban,
  },
  {
    title: 'Apps',
    description: 'Deployed application workloads',
    to: '/apps',
    icon: Sparkles,
  },
  {
    title: 'Deployments',
    description: 'Release activity and status',
    to: '/deployments',
    icon: Rocket,
  },
];

export const platformNavigation: DashboardNavigationItem[] = [
  {
    title: 'Backups',
    description: 'Backup destinations and restores',
    to: '/backups',
    icon: HardDrive,
    exact: true,
  },
  {
    title: 'API Access',
    description: 'Personal access tokens',
    to: '/api-access',
    icon: Key,
  },
  {
    title: 'Settings',
    description: 'Instance configuration',
    to: '/settings',
    icon: Settings,
    exact: true,
  },
];

export const contextualNavigation: DashboardNavigationItem[] = [
  {
    title: 'Sources',
    description: 'Git providers and registries',
    to: '/sources',
    icon: CloudCog,
  },
  {
    title: 'Servers',
    description: 'Deployment targets',
    to: '/servers',
    icon: Server,
  },
  {
    title: 'Domains and DNS',
    description: 'Domain audit and provider setup',
    to: '/dns',
    icon: Globe2,
  },
  {
    title: 'Schedules',
    description: 'Recurring tasks and service jobs',
    to: '/scheduled-tasks',
    icon: CalendarDays,
  },
  {
    title: 'Audit log',
    description: 'Security and account activity',
    to: '/audit-logs',
    icon: ClipboardList,
  },
  {
    title: 'Profile',
    description: 'Personal profile and security',
    to: '/profile',
    icon: UserRound,
  },
];

export const commandNavigation = [
  ...primaryNavigation,
  ...platformNavigation,
  ...contextualNavigation,
];
