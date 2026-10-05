import { useQueries } from '@tanstack/react-query';
import { createFileRoute } from '@tanstack/react-router';
import { Calendar, Clock, Loader2 } from 'lucide-react';
import { useState } from 'react';
import { PageHeader } from '#/components/layout/page-header';
import { Card, CardContent, CardHeader, CardTitle } from '#/components/ui/card';
import { QueryErrorState } from '#/components/ui/query-error-state';
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
import { scheduledTasksService } from '#/services/scheduled-tasks';

export const Route = createFileRoute('/_dashboard/jobs')({
  component: JobsPage,
});

export function JobsPage() {
  const [selectedProjectId, setSelectedProjectId] = useState('all');
  const {
    data: projectsResponse,
    isLoading: isLoadingProjects,
    isError: projectsError,
    refetch: refetchProjects,
  } = useListProjects();
  const projects = projectsResponse?.data?.records || [];
  const visibleProjects =
    selectedProjectId === 'all'
      ? projects
      : projects.filter((project) => project.id === selectedProjectId);
  const taskQueries = useQueries({
    queries: visibleProjects.map((project) => ({
      queryKey: ['scheduled-tasks', 'project', project.id],
      queryFn: () => scheduledTasksService.listByProject(project.id),
    })),
  });
  const isLoadingTasks = taskQueries.some((query) => query.isLoading);
  const tasks = taskQueries.flatMap((query) => query.data?.data || []);
  if (projectsError || taskQueries.some((query) => query.isError)) {
    return (
      <QueryErrorState
        title="Jobs are unavailable"
        description="Could not load scheduled tasks."
        onRetry={() => {
          void refetchProjects();
          for (const query of taskQueries) void query.refetch();
        }}
      />
    );
  }

  return (
    <div className="space-y-6">
      <PageHeader
        title="Jobs"
        description="Manage and monitor scheduled tasks, cron schedules, and maintenance jobs."
      />

      <Card className="border-border/80 bg-card shadow-sm">
        <CardHeader className="flex flex-row items-center justify-between space-y-0 border-border/70 border-b">
          <div className="flex items-center gap-2">
            <Clock className="h-4 w-4 text-primary" />
            <CardTitle>Scheduled jobs & cron tasks</CardTitle>
          </div>
          <div className="w-52">
            <Select value={selectedProjectId} onValueChange={setSelectedProjectId}>
              <SelectTrigger>
                <SelectValue placeholder="All Projects" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All Projects</SelectItem>
                {projects.map((project: { id: string; name: string }) => (
                  <SelectItem key={project.id} value={project.id}>
                    {project.name}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
        </CardHeader>
        <CardContent>
          {isLoadingProjects || isLoadingTasks ? (
            <div className="flex justify-center p-12">
              <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
            </div>
          ) : tasks.length === 0 ? (
            <div className="flex flex-col items-center justify-center p-12 text-center text-muted-foreground">
              <Calendar className="mb-4 h-8 w-8 opacity-20" />
              <p>No scheduled tasks found.</p>
            </div>
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Name</TableHead>
                  <TableHead>Schedule</TableHead>
                  <TableHead>Command</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead>Last Run</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {tasks.map((task) => (
                  <TableRow key={task.id}>
                    <TableCell className="font-medium">{task.name}</TableCell>
                    <TableCell>{task.schedule}</TableCell>
                    <TableCell className="font-mono text-xs">{task.command}</TableCell>
                    <TableCell>
                      <span className="inline-flex items-center rounded-full bg-muted px-2 py-0.5 font-medium text-xs capitalize">
                        {task.status}
                      </span>
                    </TableCell>
                    <TableCell>
                      {task.lastRunAt ? new Date(task.lastRunAt).toLocaleString() : 'Never'}
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
