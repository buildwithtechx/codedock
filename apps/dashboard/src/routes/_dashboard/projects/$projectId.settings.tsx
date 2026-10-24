import { Button } from '@codedock/ui/button';
import { createFileRoute } from '@tanstack/react-router';
import { AlertTriangle, Trash2 } from 'lucide-react';
import { useState } from 'react';
import { ProjectDeleteDialog, ProjectRollbackSettings } from '#/features/projects';
import { ProjectMembers } from '#/features/projects/project-members';
import { ProjectTokens } from '#/features/projects/project-tokens';
import { ProjectRegistries } from '#/features/registries/project-registries';
import { SettingsSection } from '#/features/settings/settings-section';
import { useGetProject } from '#/hooks/use-projects';

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

      <div className="grid grid-cols-1 gap-6">
        <section className="rounded-2xl border border-border/50 bg-card p-6 shadow-xs">
          <ProjectMembers projectId={projectId} />
        </section>

        <section className="rounded-2xl border border-border/50 bg-card p-6 shadow-xs">
          <ProjectRegistries projectId={projectId} />
        </section>

        <section className="rounded-2xl border border-border/50 bg-card p-6 shadow-xs">
          <ProjectTokens projectId={projectId} />
        </section>

        <section className="rounded-2xl border border-border/50 bg-card p-6 shadow-xs">
          <ProjectRollbackSettings projectId={projectId} />
        </section>

        <SettingsSection
          icon={<AlertTriangle className="size-4 text-destructive" />}
          title="Danger Zone"
          description="Irreversible destructive actions for this project."
          iconBg="bg-destructive/10"
          iconColor="text-destructive"
        >
          <div className="flex flex-col items-start justify-between gap-4 sm:flex-row sm:items-center">
            <div className="space-y-1">
              <h3 className="font-medium text-foreground text-sm">Delete this project</h3>
              <p className="max-w-xl text-muted-foreground text-xs leading-relaxed">
                Permanently removes{' '}
                <span className="font-semibold text-foreground">{project.name}</span>, along with
                all of its apps, services, containers, and environment configurations. This action
                cannot be undone.
              </p>
            </div>
            <Button
              variant="destructive"
              size="sm"
              onClick={() => setDeleteOpen(true)}
              className="shrink-0 gap-2"
            >
              <Trash2 className="size-4" />
              Delete project
            </Button>
          </div>
        </SettingsSection>
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
