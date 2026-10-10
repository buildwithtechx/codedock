import { createFileRoute, useNavigate } from '@tanstack/react-router';
import { Activity, Box, Folder, Plus, Rocket, Server, Settings } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
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
import { ProjectDatabaseInventory } from '#/features/projects/project-database-inventory';
import { ProjectGitTab } from '#/features/projects/project-git-tab';
import { StatusBadge } from '#/features/projects/service-status-badge';
import { ProjectClusters } from '#/features/servers/project-clusters';
import type { AppService } from '#/features/services';
import { AddServiceModal } from '#/features/services';
import { useListByProject as useListAppsByProject } from '#/hooks/use-apps';
import { useTriggerProject } from '#/hooks/use-deployments';
import { useListByProject as useListEnvironments } from '#/hooks/use-environments';

export const Route = createFileRoute('/_dashboard/projects/$projectId/')({
  component: ProjectOverviewComponent,
});

function ProjectOverviewComponent() {
  const { projectId } = Route.useParams();
  const navigate = useNavigate();

  const {
    data: projectRes,
    isLoading: projectLoading,
    isError: projectError,
    refetch: refetchProject,
  } = useGetProject(projectId);
  const {
    data: envsRes,
    isLoading: envsLoading,
    isError: envsError,
    refetch: refetchEnvironments,
  } = useListEnvironments(projectId);

  const environments = envsRes?.data || [];
  const [selectedEnvId, setSelectedEnvId] = useState<string | undefined>(undefined);
  const activeEnvId =
    environments.find((environment) => environment.id === selectedEnvId)?.id || environments[0]?.id;

  const {
    data: appsRes,
    isLoading: appsLoading,
    refetch: refetchApps,
  } = useListAppsByProject(projectId, activeEnvId, Boolean(activeEnvId) && !envsError);

  const [selectedService, setSelectedService] = useState<AppService | null>(null);
  const [drawerOpen, setDrawerOpen] = useState(false);
  const [addServiceOpen, setAddServiceOpen] = useState(false);

  const triggerProjectMutation = useTriggerProject();

  const handleDeployAll = async () => {
    if (!activeEnvId || envsError) return;
    try {
      await triggerProjectMutation.mutateAsync({
        projectId,
        environmentId: activeEnvId,
      });
      refetchApps();
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Failed to deploy project');
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

  if (envsError) {
    return (
      <QueryErrorState
        title="Environments are unavailable"
        description="Select a valid project environment before deploying."
        onRetry={() => {
          void refetchEnvironments();
        }}
      />
    );
  }

  if (environments.length === 0) {
    return (
      <div className="space-y-4 p-6">
        <h2 className="font-semibold">Create an environment</h2>
        <p className="text-muted-foreground text-sm">
          This project needs an environment before services can be deployed.
        </p>
        <EnvironmentSwitcher projectId={projectId} onSelectEnvironment={setSelectedEnvId} />
      </div>
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
            disabled={
              !activeEnvId || envsError || triggerProjectMutation.isPending || services.length === 0
            }
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
          <Button size="sm" onClick={() => setAddServiceOpen(true)} className="h-9 gap-1.5 text-xs">
            <Plus className="h-4 w-4" />
            Add Service
          </Button>
        </div>
      </div>

      <Tabs defaultValue="services" className="space-y-4">
        <TabsList>
          <TabsTrigger value="services">Services ({services.length})</TabsTrigger>
          <TabsTrigger value="git">Git</TabsTrigger>
          <TabsTrigger value="deployments">Deployments</TabsTrigger>
          <TabsTrigger value="variables">Variables</TabsTrigger>
          <TabsTrigger value="clusters">Clusters</TabsTrigger>
        </TabsList>

        <TabsContent value="clusters">
          <ProjectClusters projectId={projectId} />
        </TabsContent>
        <TabsContent value="services" className="space-y-4">
          <ProjectDatabaseInventory projectId={projectId} environmentId={activeEnvId as string} />
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
                      <span className="max-w-30 truncate text-[11px] text-primary">
                        {svc.domain}
                      </span>
                    )}
                  </div>
                </button>
              ))}
            </div>
          )}
        </TabsContent>

        <TabsContent value="git">
          <ProjectGitTab projectId={projectId} environmentId={activeEnvId} />
        </TabsContent>

        <TabsContent value="deployments">
          <ProjectDeploymentsTab projectId={projectId} />
        </TabsContent>

        <TabsContent value="variables">
          <ProjectVariablesTab key={projectId} projectId={projectId} />
        </TabsContent>
      </Tabs>

      <ServiceDetailDrawer
        service={selectedService}
        open={drawerOpen}
        onOpenChange={setDrawerOpen}
        onRefresh={refetchApps}
      />

      {activeEnvId && (
        <AddServiceModal
          open={addServiceOpen}
          onOpenChange={setAddServiceOpen}
          projectId={projectId}
          environmentId={activeEnvId}
          onCreated={() => {
            void refetchApps();
          }}
        />
      )}
    </div>
  );
}
