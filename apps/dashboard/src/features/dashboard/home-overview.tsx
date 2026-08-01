import { useState } from 'react';
import { CreateProjectModal } from '#/features/projects/create-project-modal';
import { useListCanvasSummaries } from '#/hooks/use-canvas';
import { useAuthStore } from '#/stores/auth-store';
import { HomeNextStep } from './home-next-step';
import { HomeProjectList } from './home-project-list';
import { HomeRuntimeSummary } from './home-runtime-summary';
import { HomeShortcuts } from './home-shortcuts';

export function HomeOverview() {
  const [createProjectOpen, setCreateProjectOpen] = useState(false);
  const { user } = useAuthStore();
  const { data, isLoading } = useListCanvasSummaries();
  const projects = data?.data || [];
  const firstName = user?.name?.trim().split(/\s+/)[0] || 'there';
  const hour = new Date().getHours();
  const greeting = hour < 12 ? 'Good morning' : hour < 18 ? 'Good afternoon' : 'Good evening';

  return (
    <div className="space-y-6">
      <header className="flex flex-col gap-1.5">
        <h1 className="font-medium text-2xl tracking-tight">
          {greeting}, {firstName}
        </h1>
        <p className="max-w-xl text-muted-foreground text-sm">
          Here is what is happening across your projects.
        </p>
      </header>

      <div className="grid gap-6 xl:grid-cols-[minmax(0,1fr)_21.25rem]">
        <HomeProjectList
          projects={projects}
          isLoading={isLoading}
          onCreateProject={() => setCreateProjectOpen(true)}
        />
        <div className="space-y-4">
          <HomeRuntimeSummary projects={projects} isLoading={isLoading} />
          <HomeNextStep
            onCreateProject={() => setCreateProjectOpen(true)}
            projectCount={projects.length}
          />
        </div>
      </div>

      <HomeShortcuts />
      <CreateProjectModal open={createProjectOpen} onOpenChange={setCreateProjectOpen} />
    </div>
  );
}
