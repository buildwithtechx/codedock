import {
  Activity,
  Archive,
  Clock,
  FolderKanban,
  Globe2,
  KeyRound,
  LayoutDashboard,
  LayoutGrid,
  Mail,
  PlayCircle,
  Rocket,
  Server,
  Settings,
  UsersRound,
} from 'lucide-react';
import type { NavItemProps } from './nav-item';
import type { SidebarCounts } from './use-sidebar-counts';

export function getSidebarSections(counts: SidebarCounts): {
  workspaceItems: NavItemProps[];
  infrastructureItems: NavItemProps[];
  manageItems: NavItemProps[];
} {
  return {
    workspaceItems: [
      { title: 'Dashboard', url: '/', icon: LayoutDashboard, exact: true },
      { title: 'Projects', url: '/projects', icon: FolderKanban, count: counts.projects },
      { title: 'Apps', url: '/apps', icon: LayoutGrid, count: counts.apps },
      {
        title: 'Deployments',
        url: '/deployments',
        icon: Rocket,
        statusDot: counts.apps && counts.apps > 0 ? 'running' : null,
      },
      { title: 'Actions', url: '/actions', icon: PlayCircle },
      {
        title: 'Analytics',
        url: '/monitoring',
        icon: Activity,
        statusDot: counts.hasIssues ? 'warning' : null,
      },
    ],
    infrastructureItems: [
      {
        title: 'Servers',
        url: '/servers',
        icon: Server,
        count: counts.servers,
        statusDot: counts.serverStatus,
      },
      { title: 'Domains', url: '/domains', icon: Globe2 },
      { title: 'Emails', url: '/emails', icon: Mail },
      { title: 'Backups', url: '/backups', icon: Archive },
      { title: 'Cron Jobs', url: '/jobs', icon: Clock },
    ],
    manageItems: [
      { title: 'Team', url: '/team', icon: UsersRound },
      { title: 'Keys & Tokens', url: '/tokens', icon: KeyRound },
      { title: 'Settings', url: '/settings', icon: Settings },
    ],
  };
}

export const HEADER_SECTION_TITLES: Record<string, string> = {
  projects: 'Projects',
  apps: 'Apps',
  servers: 'Servers',
  actions: 'Actions',
  monitoring: 'Analytics',
  jobs: 'Cron Jobs',
  backups: 'Backups',
  emails: 'Emails',
  deployments: 'Deployments',
  audit: 'Audit Logs',
  team: 'Team',
  tokens: 'Keys & Tokens',
  domains: 'Domains',
  settings: 'Settings',
  profile: 'Profile & Security',
};
