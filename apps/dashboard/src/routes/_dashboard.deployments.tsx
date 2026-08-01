import { createFileRoute, Link } from '@tanstack/react-router';
import { Activity, Box, FolderKanban, Loader2, Rocket } from 'lucide-react';
import { useState } from 'react';
import { Button } from '#/components/ui/button';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '#/components/ui/select';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '#/components/ui/table';
import { useListProjects } from '#/features/projects';
import { useListByProject } from '#/hooks/use-apps';
import { useListByService } from '#/hooks/use-deployments';

export const Route = createFileRoute('/_dashboard/deployments')({
  component: DeploymentsPage,
});

function DeploymentsPage() {
  const [selectedProjectId, setSelectedProjectId] = useState<string>('');
  const [selectedServiceId, setSelectedServiceId] = useState<string>('');

  const { data: projectsResponse, isLoading: isLoadingProjects } = useListProjects();
  const projects = projectsResponse?.data?.records || [];

  const { data: appsResponse, isLoading: isLoadingApps } = useListByProject(selectedProjectId);
  const apps = appsResponse?.data || [];

  const { data: deploymentsResponse, isLoading: isLoadingDeployments } =
    useListByService(selectedServiceId);
  const deploymentsPaginated = deploymentsResponse?.data;
  const deployments = deploymentsPaginated?.records || [];

  return (
    <div className="space-y-6">
      <header>
        <h1 className="font-semibold text-2xl tracking-tight">Deployments</h1>
        <p className="mt-1 text-muted-foreground text-sm">
          Inspect release history across applications in the active organization.
        </p>
      </header>

      <div className="grid min-w-0 gap-6 xl:grid-cols-[minmax(0,1fr)_18rem]">
        <section className="min-w-0 space-y-4">
          <div className="flex flex-col gap-3 rounded-2xl border border-border/80 bg-card p-4 sm:flex-row sm:items-center">
            <Select
              value={selectedProjectId}
              onValueChange={(val) => {
                setSelectedProjectId(val);
                setSelectedServiceId('');
              }}
            >
              <SelectTrigger className="w-full sm:w-52">
                <SelectValue placeholder="Select Project" />
              </SelectTrigger>
              <SelectContent>
                {projects.map((project: any) => (
                  <SelectItem key={project.id} value={project.id}>
                    {project.name}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>

            <Select
              value={selectedServiceId}
              onValueChange={setSelectedServiceId}
              disabled={!selectedProjectId || apps.length === 0}
            >
              <SelectTrigger className="w-full sm:w-52">
                <SelectValue placeholder="Select App" />
              </SelectTrigger>
              <SelectContent>
                {apps.map((app: any) => (
                  <SelectItem key={app.id} value={app.id}>
                    {app.name}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
          <section className="overflow-hidden rounded-2xl border border-border/80 bg-card shadow-sm">
            <div className="flex items-center justify-between border-border/70 border-b px-5 py-4">
              <div className="flex items-center gap-3">
                <div className="flex h-9 w-9 items-center justify-center rounded-xl bg-primary/12 text-primary">
                  <Rocket className="h-4 w-4" />
                </div>
                <div>
                  <h2 className="font-semibold text-sm">Release history</h2>
                  <p className="text-muted-foreground text-xs">
                    {selectedServiceId ? `${deployments.length} deployments` : 'Select an app'}
                  </p>
                </div>
              </div>
            </div>
            {isLoadingProjects || isLoadingApps || isLoadingDeployments ? (
              <div className="flex justify-center p-12">
                <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
              </div>
            ) : !selectedServiceId ? (
              <div className="flex min-h-80 flex-col items-center justify-center px-6 text-center">
                <Rocket className="h-8 w-8 text-primary/35" />
                <h3 className="mt-4 font-semibold">Choose an app to inspect releases</h3>
                <p className="mt-1 max-w-sm text-muted-foreground text-sm">
                  Select a project, then an app, to view build status, commits, and deployment
                  activity.
                </p>
              </div>
            ) : deployments.length === 0 ? (
              <div className="flex min-h-80 flex-col items-center justify-center px-6 text-center">
                <Activity className="h-8 w-8 text-muted-foreground/45" />
                <h3 className="mt-4 font-semibold">No deployments for this app</h3>
                <p className="mt-1 max-w-sm text-muted-foreground text-sm">
                  A deployment will appear here after the app is built or redeployed.
                </p>
              </div>
            ) : (
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Status</TableHead>
                    <TableHead>Branch</TableHead>
                    <TableHead>Commit</TableHead>
                    <TableHead>Trigger</TableHead>
                    <TableHead>Created</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {deployments.map((deployment: any) => (
                    <TableRow key={deployment.id}>
                      <TableCell className="font-medium">{deployment.status}</TableCell>
                      <TableCell>{deployment.branch || '-'}</TableCell>
                      <TableCell className="max-w-25 truncate font-mono text-xs">
                        {deployment.commitHash ? deployment.commitHash.substring(0, 7) : '-'}
                      </TableCell>
                      <TableCell>{deployment.trigger || '-'}</TableCell>
                      <TableCell>{new Date(deployment.createdAt).toLocaleString()}</TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            )}
          </section>
        </section>

        <aside className="hidden xl:sticky xl:top-6 xl:block xl:self-start">
          <section className="rounded-2xl border border-border/80 bg-card p-5 shadow-sm">
            <div className="flex items-center gap-2">
              <Activity className="h-4 w-4 text-primary" />
              <h2 className="font-semibold text-sm">Release overview</h2>
            </div>
            <div className="mt-5 space-y-4 text-sm">
              <div className="flex items-center justify-between">
                <span className="flex items-center gap-2 text-muted-foreground">
                  <FolderKanban className="h-4 w-4" />
                  Projects
                </span>
                <span className="font-semibold">{isLoadingProjects ? '-' : projects.length}</span>
              </div>
              <div className="flex items-center justify-between">
                <span className="flex items-center gap-2 text-muted-foreground">
                  <Box className="h-4 w-4" />
                  Apps
                </span>
                <span className="font-semibold">{selectedProjectId ? apps.length : '-'}</span>
              </div>
              <div className="flex items-center justify-between">
                <span className="flex items-center gap-2 text-muted-foreground">
                  <Rocket className="h-4 w-4" />
                  Selected
                </span>
                <span className="font-semibold">
                  {selectedServiceId ? deployments.length : '-'}
                </span>
              </div>
            </div>
            <div className="mt-5 border-border/70 border-t pt-4">
              <Link to="/apps/new">
                <Button variant="outline" className="w-full">
                  Add an app
                </Button>
              </Link>
            </div>
          </section>
        </aside>
      </div>
    </div>
  );
}
