import { createFileRoute } from '@tanstack/react-router';
import { AlertTriangle } from 'lucide-react';
import { useState } from 'react';
import { Button } from '#/components/ui/button';
import { ProjectDeleteDialog, useGetProject } from '#/features/projects';
import { ProjectMembers } from '#/features/projects/project-members';
import { ProjectTokens } from '#/features/projects/project-tokens';
import { ProjectRegistries } from '#/features/registries/project-registries';

export const Route = createFileRoute('/_dashboard/projects/$projectId/settings')({
  component: SettingsRouteComponent,
});

function SettingsRouteComponent() {
  const { projectId } = Route.useParams();
  const { data: projectResponse } = useGetProject(projectId);
  const [deleteOpen, setDeleteOpen] = useState(false);

  const project = projectResponse?.data ?? { id: projectId, name: 'Project' };

  return (
    <div className="mx-auto flex w-full max-w-5xl flex-col gap-6">
      <h1 className="font-bold text-2xl">Project Settings</h1>

      <div className="grid grid-cols-1 gap-8">
        <section className="rounded-lg border border-border/60 bg-card p-6 shadow-sm">
          <ProjectMembers projectId={projectId} />
        </section>

        <section className="rounded-lg border border-border/60 bg-card p-6 shadow-sm">
          <ProjectRegistries projectId={projectId} />
        </section>

        <section className="rounded-lg border border-border/60 bg-card p-6 shadow-sm">
          <ProjectTokens projectId={projectId} />
        </section>

        <section className="rounded-lg border border-destructive/40 bg-card p-6 shadow-sm">
          <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
            <div className="space-y-1">
              <div className="flex items-center gap-2">
                <AlertTriangle className="size-5 text-destructive" />
                <h2 className="font-semibold text-destructive text-lg">Danger Zone</h2>
              </div>
              <p className="text-muted-foreground text-sm">
                Permanently delete this project, along with all of its apps, services, containers,
                and environment configurations.
              </p>
            </div>
            <Button variant="destructive" onClick={() => setDeleteOpen(true)} className="shrink-0">
              Delete project
            </Button>
          </div>
        </section>
      </div>

      <ProjectDeleteDialog
        isOpen={deleteOpen}
        onOpenChange={setDeleteOpen}
        project={project}
        redirectToProjectsOnDelete
      />
    </div>
  );
}
