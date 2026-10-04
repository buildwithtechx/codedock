import { createFileRoute, useNavigate } from '@tanstack/react-router';
import { Activity, Box, Folder, Plus, Rocket, Server, Settings } from 'lucide-react';
import { useState } from 'react';
import { Button } from '#/components/ui/button';
import { QueryErrorState } from '#/components/ui/query-error-state';
import { ServiceIcon } from '#/components/ui/service-icon';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '#/components/ui/tabs';
import {
  EnvironmentSwitcher,
  ProjectDeploymentsTab,
  ProjectVariablesTab,
  ServiceDetailDrawer,
  useGetProject,
} from '#/features/projects';
import type { AppService } from '#/features/services';
import { useListByProject as useListAppsByProject } from '#/hooks/use-apps';
import { useTriggerProject } from '#/hooks/use-deployments';
import { useListByProject as useListEnvironments } from '#/hooks/use-environments';

export const Route = createFileRoute('/_dashboard/projects/$projectId/')({
  component: ProjectOverviewComponent,
});

function StatusBadge({ status }: { status: string }) {
  const s = (status || '').toLowerCase();
  let color = 'bg-gray-500/10 text-gray-500 border-gray-500/20';
  let dot = 'bg-gray-500';

  if (s === 'running' || s === 'online' || s === 'healthy') {
    color = 'bg-emerald-500/10 text-emerald-500 border-emerald-500/20';
    dot = 'bg-emerald-500 shadow-[0_0_8px_rgba(16,185,129,0.4)]';
  } else if (s === 'failed' || s === 'error' || s === 'stopped') {
    color = 'bg-red-500/10 text-red-500 border-red-500/20';
    dot = 'bg-red-500';
  } else if (s === 'deploying' || s === 'pending' || s === 'building') {
    color = 'bg-amber-500/10 text-amber-500 border-amber-500/20';
    dot = 'bg-amber-500 animate-pulse';
  }

  return (
    <div
      className={`flex items-center gap-1.5 rounded-full border px-2.5 py-0.5 font-medium text-xs ${color}`}
    >
      <div className={`h-1.5 w-1.5 rounded-full ${dot}`} />
      <span className="capitalize">{status || 'Unknown'}</span>
    </div>
  );
}

function ProjectOverviewComponent() {
  const { projectId } = Route.useParams();
  const navigate = useNavigate();

  const {
    data: projectRes,
    isLoading: projectLoading,
    isError: projectError,
    refetch: refetchProject,
  } = useGetProject(projectId);
  const { data: envsRes, isLoading: envsLoading } = useListEnvironments(projectId);

  const environments = envsRes?.data || [];
  const [selectedEnvId, setSelectedEnvId] = useState<string | undefined>(undefined);
  const activeEnvId = selectedEnvId || environments[0]?.id;

  const {
    data: appsRes,
    isLoading: appsLoading,
    refetch: refetchApps,
  } = useListAppsByProject(projectId, activeEnvId);

  const [selectedService, setSelectedService] = useState<AppService | null>(null);
  const [drawerOpen, setDrawerOpen] = useState(false);

  const triggerProjectMutation = useTriggerProject();

  const handleDeployAll = async () => {
    try {
      await triggerProjectMutation.mutateAsync({
        projectId,
        environmentId: activeEnvId,
      });
      refetchApps();
    } catch {
      // Handled
    }
  };

  const handleOpenDrawer = (svc: AppService) => {
    setSelectedService(svc);
    setDrawerOpen(true);
  };

  if (projectLoading || envsLoading) {
    return (
      <div className="flex h-full min-h-100 items-center justify-center">
        <div className="flex flex-col items-center gap-4">
          <Activity className="h-8 w-8 animate-pulse text-primary" />
          <p className="text-muted-foreground text-sm">Loading project workspace...</p>
        </div>
      </div>
    );
  }

  if (projectError) {
    return (
      <QueryErrorState
        title="Project data is unavailable"
        description="Could not load this project workspace."
        onRetry={() => {
          void refetchProject();
        }}
      />
    );
  }

  const project = projectRes?.data;
  const services = appsRes?.data || [];

  return (
    <div className="flex flex-col gap-6">
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <div className="flex items-center gap-3">
            <h1 className="font-semibold text-2xl text-foreground tracking-tight">
              {project?.name || 'Project Overview'}
            </h1>
            <EnvironmentSwitcher
              projectId={projectId}
              activeEnvironmentId={activeEnvId}
              onSelectEnvironment={setSelectedEnvId}
            />
          </div>
          <div className="mt-1 flex items-center gap-2 text-muted-foreground text-xs">
            <Folder className="h-3.5 w-3.5" />
            <span>{project?.description || 'No description provided'}</span>
          </div>
        </div>

        <div className="flex flex-wrap items-center gap-2">
          <Button
            variant="outline"
            size="sm"
            onClick={handleDeployAll}
            disabled={triggerProjectMutation.isPending || services.length === 0}
            className="h-9 gap-1.5 text-xs"
          >
            <Rocket className="h-3.5 w-3.5 text-primary" />
            {triggerProjectMutation.isPending ? 'Deploying...' : 'Deploy Project'}
          </Button>
          <Button
            variant="outline"
            size="sm"
            onClick={() => navigate({ to: '/projects/$projectId/settings', params: { projectId } })}
            className="h-9 gap-1.5 text-xs"
          >
            <Settings className="h-3.5 w-3.5" />
            Settings
          </Button>
          <Button
            size="sm"
            onClick={() => navigate({ to: '/projects/$projectId/new', params: { projectId } })}
            className="h-9 gap-1.5 text-xs"
          >
            <Plus className="h-4 w-4" />
            New Service
          </Button>
        </div>
      </div>

      <Tabs defaultValue="services" className="space-y-4">
        <TabsList>
          <TabsTrigger value="services">Services ({services.length})</TabsTrigger>
          <TabsTrigger value="deployments">Deployments</TabsTrigger>
          <TabsTrigger value="variables">Variables</TabsTrigger>
        </TabsList>

        <TabsContent value="services" className="space-y-4">
          {appsLoading ? (
            <div className="flex h-40 items-center justify-center">
              <Activity className="h-6 w-6 animate-pulse text-primary" />
            </div>
          ) : services.length === 0 ? (
            <section className="flex min-h-80 flex-col items-center justify-center rounded-2xl border border-border/80 border-dashed bg-card/40 p-8 text-center">
              <span className="flex h-12 w-12 items-center justify-center rounded-xl border bg-card text-muted-foreground">
                <Box className="h-6 w-6" />
              </span>
              <h3 className="mt-4 font-semibold text-foreground text-lg">
                No services in this environment
              </h3>
              <p className="mt-1.5 max-w-sm text-muted-foreground text-xs leading-5">
                Deploy an application, container image, or background service to this environment.
              </p>
              <Button
                size="sm"
                className="mt-5"
                onClick={() => navigate({ to: '/projects/$projectId/new', params: { projectId } })}
              >
                <Plus className="mr-1.5 h-4 w-4" />
                Add service
              </Button>
            </section>
          ) : (
            <div className="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
              {services.map((svc: AppService) => (
                <button
                  type="button"
                  key={svc.id}
                  onClick={() => handleOpenDrawer(svc)}
                  className="group flex cursor-pointer flex-col rounded-xl border border-border/70 bg-card p-5 text-left transition-all hover:border-primary/50 hover:shadow-md"
                >
                  <div className="flex items-start justify-between">
                    <div className="flex h-11 w-11 items-center justify-center overflow-hidden rounded-lg border bg-background/50 shadow-xs transition-colors group-hover:border-primary/20 group-hover:bg-primary/5">
                      {svc.icon && svc.icon !== 'git' ? (
                        <ServiceIcon icon={svc.icon} className="h-6 w-6" />
                      ) : (
                        <Box className="h-6 w-6 text-primary" />
                      )}
                    </div>
                    <StatusBadge status={svc.status} />
                  </div>

                  <div className="mt-4 flex-1">
                    <h3 className="font-semibold text-foreground transition-colors group-hover:text-primary">
                      {svc.name}
                    </h3>
                    <p className="mt-1 line-clamp-1 text-muted-foreground text-xs uppercase tracking-wider">
                      {svc.runtimeMode || 'Service'}
                    </p>
                  </div>

                  <div className="mt-5 flex w-full items-center justify-between border-border/60 border-t pt-3 text-muted-foreground text-xs">
                    <span className="flex items-center gap-1.5 font-mono">
                      <Server className="h-3.5 w-3.5" />
                      {svc.internalPort ? `:${svc.internalPort}` : 'No Port'}
                    </span>
                    {svc.domain && (
                      <span className="max-w-[120px] truncate text-[11px] text-primary">
                        {svc.domain}
                      </span>
                    )}
                  </div>
                </button>
              ))}
            </div>
          )}
        </TabsContent>

        <TabsContent value="deployments">
          <ProjectDeploymentsTab projectId={projectId} />
        </TabsContent>

        <TabsContent value="variables">
          <ProjectVariablesTab projectId={projectId} />
        </TabsContent>
      </Tabs>

      <ServiceDetailDrawer
        service={selectedService}
        open={drawerOpen}
        onOpenChange={setDrawerOpen}
        onRefresh={refetchApps}
      />
    </div>
  );
}
