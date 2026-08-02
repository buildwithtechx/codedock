import type { LucideIcon } from 'lucide-react';
import {
  CalendarDays,
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

export const infrastructureNavigation: DashboardNavigationItem[] = [
  {
    title: 'Servers',
    description: 'Deployment targets and runtime capacity',
    to: '/servers',
    icon: Server,
  },
  {
    title: 'Domains & DNS',
    description: 'Domain audit and provider setup',
    to: '/dns',
    icon: Globe2,
  },
  {
    title: 'Backups',
    description: 'Backup destinations and restores',
    to: '/backups',
    icon: HardDrive,
    exact: true,
  },
];

export const systemNavigation: DashboardNavigationItem[] = [
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
    to: '/settings',
    icon: CloudCog,
    search: { tab: 'sources' },
  },
  {
    title: 'Schedules',
    description: 'Recurring tasks and service jobs',
    to: '/scheduled-tasks',
    icon: CalendarDays,
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
  ...infrastructureNavigation,
  ...systemNavigation,
  ...contextualNavigation,
];
