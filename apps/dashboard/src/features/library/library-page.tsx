import { useNavigate } from '@tanstack/react-router';
import { Code2, FolderUp, LayoutGrid, Link2, SquareCode } from 'lucide-react';
import { useEffect, useState } from 'react';
import { PageHeader } from '#/components/layout/page-header';
import { Skeleton } from '#/components/ui/skeleton';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '#/components/ui/tabs';
import { AppCatalog } from '#/features/apps';
import { TemplatesGallery } from '#/features/apps/templates-gallery';
import { encodeDeploySlug } from '#/lib/slug-utils';
import { ConnectPrompt } from './connect-prompt';
import { FolderUpload } from './folder-upload';
import { useGitConnections, useLibraryRepos } from './hooks';
import { LibrarySidebar } from './library-sidebar';
import { LocalProjects } from './local-projects';
import { RepositoryList } from './repository-list';
import { UrlImport } from './url-import';

export type LibraryTab = 'repositories' | 'folder' | 'url' | 'apps' | 'examples';

export function LibraryPage({ initialTab }: { initialTab?: LibraryTab }) {
  const navigate = useNavigate();
  const [tab, setTab] = useState<LibraryTab>(initialTab || 'repositories');

  useEffect(() => {
    if (initialTab) {
      setTab(initialTab);
    }
  }, [initialTab]);
  const [provider, setProvider] = useState('github');
  const connectionsQuery = useGitConnections();
  const connections = connectionsQuery.data ?? [];
  const anyConnected = connections.some((item) => item.connected);
  const countsQuery = useLibraryRepos(provider, anyConnected);

  const counts = countsQuery.counts;

  return (
    <div className="space-y-6">
      <PageHeader
        title="Library"
        description="Import code from git, a folder, or a URL — or start from a one-click app or example."
      />

      <Tabs value={tab} onValueChange={(value) => setTab(value as LibraryTab)} className="w-full">
        <TabsList>
          <TabsTrigger value="repositories">
            <SquareCode className="size-4" />
            Repositories
          </TabsTrigger>
          <TabsTrigger value="folder">
            <FolderUp className="size-4" />
            Folder
          </TabsTrigger>
          <TabsTrigger value="url">
            <Link2 className="size-4" />
            URL
          </TabsTrigger>
          <TabsTrigger value="apps">
            <LayoutGrid className="size-4" />
            Apps
          </TabsTrigger>
          <TabsTrigger value="examples">
            <Code2 className="size-4" />
            Examples
          </TabsTrigger>
        </TabsList>

        <div className="mt-6 grid grid-cols-1 items-start gap-6 lg:grid-cols-[minmax(0,1fr)_340px]">
          <div className="min-w-0">
            <TabsContent value="repositories" className="mt-0">
              {connectionsQuery.isLoading ? (
                <div className="space-y-3 rounded-2xl border border-border/50 bg-card p-5">
                  <Skeleton className="h-9 w-2/3" />
                  <Skeleton className="h-10 w-full" />
                  <Skeleton className="h-40 w-full" />
                </div>
              ) : !anyConnected ? (
                <ConnectPrompt onBrowseApps={() => setTab('apps')} />
              ) : (
                <RepositoryList
                  connections={connections}
                  provider={provider}
                  onProviderChange={setProvider}
                  onImport={(repo) => {
                    void navigate({
                      to: '/deploy/$slug',
                      params: { slug: encodeDeploySlug(repo.fullName || repo.name) },
                      search: { name: repo.name, branch: repo.defaultBranch },
                    });
                  }}
                  onImportUrl={() => setTab('url')}
                />
              )}
            </TabsContent>
            <TabsContent value="folder" className="mt-0 space-y-6">
              <FolderUpload />
              <LocalProjects
                onImportFolder={() =>
                  document
                    .getElementById('library-folder-upload')
                    ?.scrollIntoView({ behavior: 'smooth', block: 'start' })
                }
              />
            </TabsContent>
            <TabsContent value="url" className="mt-0">
              <UrlImport
                onImport={(parsed) => {
                  void navigate({
                    to: '/deploy/$slug',
                    params: { slug: encodeDeploySlug(parsed.repositoryUrl) },
                    search: { name: parsed.name, branch: parsed.branch },
                  });
                }}
              />
            </TabsContent>
            <TabsContent value="apps" className="mt-0">
              <AppCatalog embedded />
            </TabsContent>
            <TabsContent value="examples" className="mt-0">
              <TemplatesGallery />
            </TabsContent>
          </div>
          <LibrarySidebar connections={connections} counts={counts} />
        </div>
      </Tabs>
    </div>
  );
}
