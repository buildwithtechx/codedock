import { createFileRoute, useNavigate } from '@tanstack/react-router';
import {
  ArrowRight,
  Box,
  Boxes,
  Code2,
  Container,
  Database,
  FolderUp,
  GitBranch,
  LayoutGrid,
  Plus,
} from 'lucide-react';
import { useEffect, useState } from 'react';
import { z } from 'zod';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@codedock/ui/tabs';
import { AppCatalog } from '#/features/apps';
import { TemplatesGallery } from '#/features/apps/templates-gallery';
import { CreateDatabaseModal } from '#/features/databases/create-database-modal';
import { FolderUpload } from '#/features/library/folder-upload';
import { AddServiceModal } from '#/features/services';
import { CreateDockerImageModal } from '#/features/sources/create-docker-image-modal';
import { CreateGitAppModal } from '#/features/sources/create-git-app-modal';
import { useListByProject as useListEnvironments } from '#/hooks/use-environments';

export const Route = createFileRoute('/_dashboard/projects/$projectId/new')({
  component: ProjectNewResourcePage,
  validateSearch: z.object({
    tab: z.enum(['resources', 'one-click', 'examples']).optional(),
    resource: z.enum(['git', 'database', 'docker', 'folder', 'companion']).optional(),
  }),
});

function ProjectNewResourcePage() {
  const { projectId } = Route.useParams();
  const search = Route.useSearch();
  const navigate = useNavigate();
  const { data: envsRes } = useListEnvironments(projectId);
  const activeEnvId = envsRes?.data?.[0]?.id || '';

  const [tab, setTab] = useState<'resources' | 'one-click' | 'examples'>(search.tab || 'resources');
  const [dbModalOpen, setDbModalOpen] = useState(search.resource === 'database');
  const [gitModalOpen, setGitModalOpen] = useState(search.resource === 'git');
  const [dockerModalOpen, setDockerModalOpen] = useState(search.resource === 'docker');
  const [companionModalOpen, setCompanionModalOpen] = useState(search.resource === 'companion');
  const [showFolderUpload, setShowFolderUpload] = useState(search.resource === 'folder');

  useEffect(() => {
    if (search.tab) setTab(search.tab);
    if (search.resource === 'database') setDbModalOpen(true);
    if (search.resource === 'git') setGitModalOpen(true);
    if (search.resource === 'docker') setDockerModalOpen(true);
    if (search.resource === 'companion') setCompanionModalOpen(true);
    if (search.resource === 'folder') setShowFolderUpload(true);
  }, [search.tab, search.resource]);

  const resourceCards = [
    {
      id: 'git',
      title: 'Git Repository',
      description:
        'Deploy source code from GitHub, GitLab, Bitbucket, or Gitea with automatic builds.',
      hint: 'Best for applications you develop and push.',
      icon: GitBranch,
      onClick: () => setGitModalOpen(true),
    },
    {
      id: 'database',
      title: 'Database',
      description: 'Provision PostgreSQL, MySQL, Redis, MongoDB, or another managed database.',
      hint: 'Best for stateful backing services.',
      icon: Database,
      onClick: () => setDbModalOpen(true),
    },
    {
      id: 'docker',
      title: 'Docker Image',
      description: 'Run a prebuilt image from any public or private registry without a build step.',
      hint: 'Best for published images and third-party tools.',
      icon: Container,
      onClick: () => setDockerModalOpen(true),
    },
    {
      id: 'companion',
      title: 'Companion Service',
      description: 'Add a preconfigured cache, search engine, queue, or companion container.',
      hint: 'Quick setup with pre-filled ports and volumes.',
      icon: Box,
      onClick: () => setCompanionModalOpen(true),
    },
    {
      id: 'compose',
      title: 'Docker Compose',
      description: 'Import a Compose file and review every service before Codedock creates them.',
      hint: 'Best for multi-service stacks you already defined.',
      icon: Boxes,
      onClick: () => void navigate({ to: '/projects/$projectId/compose', params: { projectId } }),
    },
    {
      id: 'folder',
      title: 'Upload Folder',
      description: 'Pack a local folder or .tar.gz archive in your browser and deploy it directly.',
      hint: 'Best for local prototypes and static or Node/Go builds.',
      icon: FolderUp,
      onClick: () => setShowFolderUpload((prev) => !prev),
    },
  ];

  return (
    <div className="space-y-6">
      <div>
        <h1 className="font-semibold text-2xl text-foreground tracking-tight">Add to project</h1>
        <p className="mt-1 text-muted-foreground text-sm">
          Deploy a service, database, one-click application, or starter template into this project.
        </p>
      </div>

      <Tabs
        value={tab}
        onValueChange={(val) => setTab(val as 'resources' | 'one-click' | 'examples')}
        className="w-full"
      >
        <TabsList variant="line">
          <TabsTrigger value="resources">
            <Plus className="size-4" />
            Resources
          </TabsTrigger>
          <TabsTrigger value="one-click">
            <LayoutGrid className="size-4" />
            Apps
          </TabsTrigger>
          <TabsTrigger value="examples">
            <Code2 className="size-4" />
            Examples
          </TabsTrigger>
        </TabsList>

        <TabsContent value="resources" className="mt-6 space-y-6">
          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
            {resourceCards.map((card) => (
              <button
                key={card.id}
                type="button"
                onClick={card.onClick}
                className="group flex flex-col justify-between gap-3 rounded-2xl border border-border/50 bg-card p-5 text-left transition-all hover:border-primary/40 hover:shadow-md"
              >
                <div className="space-y-3">
                  <div className="flex items-center justify-between gap-2">
                    <span className="flex size-10 shrink-0 items-center justify-center rounded-xl bg-primary/10 text-primary">
                      <card.icon className="size-5" />
                    </span>
                    <ArrowRight className="size-4 shrink-0 text-muted-foreground/40 transition-colors group-hover:text-foreground" />
                  </div>
                  <div>
                    <p className="font-medium text-foreground">{card.title}</p>
                    <p className="mt-1 line-clamp-2 text-muted-foreground text-sm">
                      {card.description}
                    </p>
                  </div>
                </div>
                <p className="text-[11px] text-muted-foreground/80">{card.hint}</p>
              </button>
            ))}
          </div>

          {showFolderUpload && <FolderUpload fixedProjectId={projectId} />}
        </TabsContent>

        <TabsContent value="one-click" className="mt-6">
          <AppCatalog embedded projectId={projectId} />
        </TabsContent>

        <TabsContent value="examples" className="mt-6">
          <TemplatesGallery projectId={projectId} />
        </TabsContent>
      </Tabs>

      <CreateDatabaseModal
        isOpen={dbModalOpen}
        onOpenChange={setDbModalOpen}
        projectId={projectId}
        environmentId={activeEnvId}
      />
      <CreateGitAppModal
        isOpen={gitModalOpen}
        onOpenChange={setGitModalOpen}
        projectId={projectId}
      />
      <CreateDockerImageModal
        isOpen={dockerModalOpen}
        onOpenChange={setDockerModalOpen}
        projectId={projectId}
      />
      {activeEnvId && (
        <AddServiceModal
          open={companionModalOpen}
          onOpenChange={setCompanionModalOpen}
          projectId={projectId}
          environmentId={activeEnvId}
          onCreated={() => {
            void navigate({ to: '/projects/$projectId', params: { projectId } });
          }}
        />
      )}
    </div>
  );
}
