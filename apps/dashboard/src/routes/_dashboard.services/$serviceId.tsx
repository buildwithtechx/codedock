import { createFileRoute, Link, Outlet, useLocation } from '@tanstack/react-router';
import {
  Activity,
  AlertTriangle,
  ArrowLeft,
  BarChart2,
  Calendar,
  Code,
  GitPullRequest,
  Globe,
  HardDrive,
  Loader2,
  Network,
  Settings,
  Shield,
  Terminal,
  Variable,
  Webhook,
  Wrench,
} from 'lucide-react';
import { MobileContextNav } from '#/components/layout/mobile-context-nav';
import { PageFrame } from '#/components/layout/page-frame';
import { Button } from '#/components/ui/button';
import { ServiceIcon } from '#/components/ui/service-icon';
import { useGetDatabase } from '#/features/databases';
import {
  ServiceContextSidebar,
  type ServiceContextTab,
} from '#/features/services/service-context-sidebar';
import { useGetApp } from '#/hooks/use-apps';

export const Route = createFileRoute('/_dashboard/services/$serviceId')({
  component: ServiceLayoutRoute,
});

function ServiceLayoutRoute() {
  const { serviceId } = Route.useParams();
  const location = useLocation();

  const { data: appData, isLoading: appLoading } = useGetApp(serviceId);
  const { data: dbData, isLoading: dbLoading } = useGetDatabase(serviceId);

  if (appLoading || dbLoading) {
    return (
      <div className="flex justify-center p-12">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  const app = appData?.data;
  const db = dbData?.data;

  const appTabs = [
    {
      name: 'Configuration',
      href: `/services/${serviceId}`,
      icon: Settings,
      group: 'Setup',
      exact: true,
    },
    { name: 'Build Settings', href: `/services/${serviceId}/build`, icon: Wrench, group: 'Setup' },
    { name: 'Variables', href: `/services/${serviceId}/variables`, icon: Variable, group: 'Setup' },
    { name: 'Domains', href: `/services/${serviceId}/domains`, icon: Globe, group: 'Setup' },
    {
      name: 'Route Rules',
      href: `/services/${serviceId}/route-rules`,
      icon: Shield,
      group: 'Setup',
    },
    {
      name: 'Deployments',
      href: `/services/${serviceId}/deployments`,
      icon: Activity,
      group: 'Observe',
    },
    { name: 'Metrics', href: `/services/${serviceId}/metrics`, icon: BarChart2, group: 'Observe' },
    { name: 'Webhooks', href: `/services/${serviceId}/webhooks`, icon: Webhook, group: 'Observe' },
    {
      name: 'Log Drains',
      href: `/services/${serviceId}/log-drains`,
      icon: Network,
      group: 'Observe',
    },
    {
      name: 'Scheduled Tasks',
      href: `/services/${serviceId}/scheduled-tasks`,
      icon: Calendar,
      group: 'Operate',
    },
    { name: 'Storage', href: `/services/${serviceId}/volumes`, icon: HardDrive, group: 'Operate' },
    { name: 'Terminal', href: `/services/${serviceId}/terminal`, icon: Terminal, group: 'Operate' },
    {
      name: 'Serverless Editor',
      href: `/services/${serviceId}/serverless`,
      icon: Code,
      group: 'Operate',
    },
    {
      name: 'PR Previews',
      href: `/services/${serviceId}/previews`,
      icon: GitPullRequest,
      group: 'Operate',
    },
    {
      name: 'Danger Zone',
      href: `/services/${serviceId}/danger`,
      icon: AlertTriangle,
      group: 'Danger',
    },
  ];

  const dbTabs = [
    {
      name: 'Overview',
      href: `/services/${serviceId}`,
      icon: Settings,
      group: 'Database',
      exact: true,
    },
    {
      name: 'Danger Zone',
      href: `/services/${serviceId}/danger`,
      icon: AlertTriangle,
      group: 'Danger',
    },
  ];

  const tabs: ServiceContextTab[] = app ? appTabs : dbTabs;
  const resourceName = app?.name || db?.name || 'Service';
  const projectId = app?.projectId || db?.projectId;

  return (
    <div className="space-y-6">
      <header className="flex items-center gap-3">
        <Button variant="ghost" size="icon" asChild className="h-8 w-8 shrink-0">
          <Link to="/projects/$projectId" params={{ projectId: projectId as string }}>
            <ArrowLeft className="h-4 w-4" />
          </Link>
        </Button>
        {app?.icon && app.icon !== 'git' && (
          <ServiceIcon icon={app.icon} className="h-10 w-10 rounded-xl border border-border" />
        )}
        <div className="min-w-0">
          <p className="text-muted-foreground text-sm">
            {app ? 'Application service' : 'Database'}
          </p>
          <p className="truncate font-semibold text-base">{resourceName}</p>
        </div>
      </header>

      <MobileContextNav
        label="Service sections"
        items={tabs.map((tab) => ({
          title: tab.name,
          to: tab.href,
          icon: tab.icon,
          active: tab.exact
            ? location.pathname === tab.href || location.pathname === `${tab.href}/`
            : location.pathname.startsWith(tab.href),
        }))}
      />

      <PageFrame
        rail={
          <ServiceContextSidebar
            name={resourceName}
            type={app ? 'Application service' : 'Database'}
            status={app?.status || db?.status}
            tabs={tabs}
          />
        }
      >
        <Outlet />
      </PageFrame>
    </div>
  );
}
