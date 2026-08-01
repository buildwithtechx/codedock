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

  // Tabs for Application Services
  const appTabs = [
    { name: 'Configuration', href: `/services/${serviceId}`, icon: Settings, exact: true },
    { name: 'Build Settings', href: `/services/${serviceId}/build`, icon: Wrench },
    { name: 'Deployments', href: `/services/${serviceId}/deployments`, icon: Activity },
    { name: 'Metrics', href: `/services/${serviceId}/metrics`, icon: BarChart2 },
    { name: 'Webhooks', href: `/services/${serviceId}/webhooks`, icon: Webhook },
    { name: 'Scheduled Tasks', href: `/services/${serviceId}/scheduled-tasks`, icon: Calendar },
    { name: 'Storage', href: `/services/${serviceId}/volumes`, icon: HardDrive },
    { name: 'Domains', href: `/services/${serviceId}/domains`, icon: Globe },
    { name: 'Route Rules', href: `/services/${serviceId}/route-rules`, icon: Shield },
    { name: 'Variables', href: `/services/${serviceId}/variables`, icon: Variable },
    { name: 'Terminal', href: `/services/${serviceId}/terminal`, icon: Terminal },
    { name: 'Serverless Editor', href: `/services/${serviceId}/serverless`, icon: Code },
    { name: 'PR Previews', href: `/services/${serviceId}/previews`, icon: GitPullRequest },
    { name: 'Log Drains', href: `/services/${serviceId}/log-drains`, icon: Network },
    { name: 'Danger Zone', href: `/services/${serviceId}/danger`, icon: AlertTriangle },
  ];

  // Tabs for Databases
  const dbTabs = [
    { name: 'Overview', href: `/services/${serviceId}`, icon: Settings, exact: true },
    { name: 'Danger Zone', href: `/services/${serviceId}/danger`, icon: AlertTriangle },
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
          <h1 className="truncate font-semibold text-2xl tracking-tight">{resourceName}</h1>
        </div>
      </header>

      <div className="lg:hidden">
        <div className="flex gap-1 overflow-x-auto rounded-xl border border-border/80 bg-card p-1.5">
          {tabs.map((tab) => {
            const isActive = tab.exact
              ? location.pathname === tab.href || location.pathname === `${tab.href}/`
              : location.pathname.startsWith(tab.href);
            return (
              <Link
                key={tab.name}
                to={tab.href}
                className={`flex shrink-0 items-center gap-2 whitespace-nowrap rounded-lg px-3 py-2 font-medium text-sm transition-colors ${
                  isActive
                    ? 'bg-primary/12 text-foreground'
                    : 'text-muted-foreground hover:bg-muted hover:text-foreground'
                }`}
              >
                <tab.icon className="h-4 w-4" />
                {tab.name}
              </Link>
            );
          })}
        </div>
      </div>

      <div className="grid min-w-0 gap-6 lg:grid-cols-[minmax(0,1fr)_18rem]">
        <div className="min-w-0">
          <Outlet />
        </div>
        <div className="hidden lg:block">
          <ServiceContextSidebar
            name={resourceName}
            type={app ? 'Application service' : 'Database'}
            status={app?.status || db?.status}
            tabs={tabs}
          />
        </div>
      </div>
    </div>
  );
}
