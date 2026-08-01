import { useState } from 'react';
import { CreateProjectModal } from '#/features/projects/create-project-modal';
import { useListCanvasSummaries } from '#/hooks/use-canvas';
import { useAuthStore } from '#/stores/auth-store';
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
        <h1 className="font-semibold text-2xl tracking-tight sm:text-3xl">
          {greeting}, {firstName}
        </h1>
        <p className="max-w-xl text-muted-foreground text-sm">
          Here is the current shape of your Codedock workspace.
        </p>
      </header>

      <div className="grid gap-5 xl:grid-cols-[minmax(0,1fr)_19rem]">
        <HomeProjectList
          projects={projects}
          isLoading={isLoading}
          onCreateProject={() => setCreateProjectOpen(true)}
        />
        <HomeRuntimeSummary projects={projects} isLoading={isLoading} />
      </div>

      <HomeShortcuts />
      <CreateProjectModal open={createProjectOpen} onOpenChange={setCreateProjectOpen} />
    </div>
  );
}
