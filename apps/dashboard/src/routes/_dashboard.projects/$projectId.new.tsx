import { createFileRoute, useNavigate } from '@tanstack/react-router';
import { Code2, Container, Database, GitBranch, LayoutTemplate, Loader2 } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { z } from 'zod';
import { Button } from '#/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '#/components/ui/card';
import { QueryErrorState } from '#/components/ui/query-error-state';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '#/components/ui/tabs';
import { CreateDatabaseModal } from '#/features/databases/create-database-modal';
import { CreateDockerImageModal } from '#/features/sources/create-docker-image-modal';
import { CreateGitAppModal } from '#/features/sources/create-git-app-modal';
import {
  useDeployOneClickApp,
  useListExampleApps,
  useListOneClickApps,
} from '#/hooks/use-templates';

export const Route = createFileRoute('/_dashboard/projects/$projectId/new')({
  component: NewResourcePage,
  validateSearch: z.object({
    tab: z.enum(['resources', 'one-click', 'examples']).optional(),
  }),
});

function NewResourcePage() {
  const { projectId } = Route.useParams();
  const search = Route.useSearch();
  const navigate = useNavigate();
  const [dbModalOpen, setDbModalOpen] = useState(false);
  const [gitModalOpen, setGitModalOpen] = useState(false);
  const [dockerModalOpen, setDockerModalOpen] = useState(false);

  const {
    data: oneClickResponse,
    isLoading: oneClickLoading,
    isError: oneClickError,
    refetch: refetchOneClick,
  } = useListOneClickApps();
  const {
    data: examplesResponse,
    isLoading: examplesLoading,
    isError: examplesError,
    refetch: refetchExamples,
  } = useListExampleApps();
  const deployOneClick = useDeployOneClickApp();
  const [deployingTemplateId, setDeployingTemplateId] = useState<string | null>(null);

  const templates = Array.isArray(oneClickResponse) ? oneClickResponse : [];
  const examples = Array.isArray(examplesResponse) ? examplesResponse : [];

  const handleDeployOneClick = async (appId: string, name: string) => {
    setDeployingTemplateId(appId);
    try {
      const response = await deployOneClick.mutateAsync({ appId, projectId, name });
      toast.success(response.message || `${name} deployed`);
      await navigate({ to: '/projects/$projectId', params: { projectId } });
    } catch {
      toast.error(`Failed to deploy ${name}`);
    } finally {
      setDeployingTemplateId(null);
    }
  };

  return (
    <div className="space-y-6">
      <div className="mb-5 flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="flex h-12 w-12 shrink-0 items-center justify-center rounded-lg border border-primary/20 bg-primary/10 text-primary">
            <LayoutTemplate className="h-6 w-6" />
          </div>
          <div>
            <h1 className="font-bold text-xl">Add New Resource</h1>
            <p className="text-muted-foreground text-sm">
              Deploy a new application, database, or service to this project.
            </p>
          </div>
        </div>
      </div>

      <Tabs defaultValue={search.tab || 'resources'} className="w-full">
        <TabsList className="mb-6 grid w-full max-w-150 grid-cols-3">
          <TabsTrigger value="resources">Resources</TabsTrigger>
          <TabsTrigger value="one-click">One-Click Apps</TabsTrigger>
          <TabsTrigger value="examples">Example Projects</TabsTrigger>
        </TabsList>

        <TabsContent value="resources" className="space-y-4">
          <div className="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
            <Card
              className="cursor-pointer transition-colors hover:border-primary/50"
              onClick={() => setGitModalOpen(true)}
            >
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <GitBranch className="h-5 w-5" />
                  Git Repository
                </CardTitle>
                <CardDescription>
                  Deploy source code from a public or private Git repository.
                </CardDescription>
              </CardHeader>
            </Card>

            <Card
              className="cursor-pointer transition-colors hover:border-primary/50"
              onClick={() => setDbModalOpen(true)}
            >
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Database className="h-5 w-5" />
                  Database
                </CardTitle>
                <CardDescription>
                  Provision a PostgreSQL, MySQL, Redis, or other database.
                </CardDescription>
              </CardHeader>
            </Card>

            <Card
              className="cursor-pointer transition-colors hover:border-primary/50"
              onClick={() => setDockerModalOpen(true)}
            >
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Container className="h-5 w-5" />
                  Docker Image
                </CardTitle>
                <CardDescription>
                  Deploy a pre-built Docker image from any public or private registry.
                </CardDescription>
              </CardHeader>
            </Card>

            <Card
              className="cursor-pointer transition-colors hover:border-primary/50"
              onClick={() =>
                navigate({
                  to: '/projects/$projectId/compose',
                  params: { projectId },
                })
              }
            >
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Code2 className="h-5 w-5" />
                  Docker Compose
                </CardTitle>
                <CardDescription>
                  Deploy multiple services defined in a docker-compose.yml file.
                </CardDescription>
              </CardHeader>
            </Card>
          </div>
        </TabsContent>

        <TabsContent value="one-click" className="space-y-4">
          {oneClickLoading ? (
            <div className="flex justify-center p-12">
              <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
            </div>
          ) : oneClickError ? (
            <QueryErrorState
              title="Templates are unavailable"
              description="Codedock could not load the one-click application catalogue."
              onRetry={() => void refetchOneClick()}
            />
          ) : templates.length === 0 ? (
            <section className="flex min-h-80 flex-col items-center justify-center px-6 text-center">
              <span className="flex h-13 w-13 items-center justify-center rounded-2xl bg-card text-muted-foreground">
                <LayoutTemplate className="h-5 w-5" />
              </span>
              <h2 className="mt-5 font-medium text-foreground/90 text-lg">
                No one-click templates yet
              </h2>
              <p className="mt-2 max-w-sm text-muted-foreground/75 text-sm leading-6">
                Compatible application templates will appear here when they are available on this
                instance.
              </p>
            </section>
          ) : (
            <div className="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
              {templates.map((template) => (
                <Card key={template.id}>
                  <CardHeader>
                    <CardTitle className="flex items-center gap-2">
                      {(template as any).logo && (
                        <img src={(template as any).logo} alt={template.name} className="h-6 w-6" />
                      )}
                      {template.name}
                    </CardTitle>
                    <CardDescription>{template.description}</CardDescription>
                  </CardHeader>
                  <CardContent>
                    <Button
                      className="w-full"
                      onClick={() => handleDeployOneClick(template.id, template.name)}
                      disabled={deployOneClick.isPending}
                    >
                      {deployingTemplateId === template.id ? 'Deploying...' : 'Deploy'}
                    </Button>
                  </CardContent>
                </Card>
              ))}
            </div>
          )}
        </TabsContent>

        <TabsContent value="examples" className="space-y-4">
          {examplesLoading ? (
            <div className="flex justify-center p-12">
              <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
            </div>
          ) : examplesError ? (
            <QueryErrorState
              title="Examples are unavailable"
              description="Codedock could not load the example project catalogue."
              onRetry={() => void refetchExamples()}
            />
          ) : examples.length === 0 ? (
            <section className="flex min-h-80 flex-col items-center justify-center px-6 text-center">
              <span className="flex h-13 w-13 items-center justify-center rounded-2xl bg-card text-muted-foreground">
                <Code2 className="h-5 w-5" />
              </span>
              <h2 className="mt-5 font-medium text-foreground/90 text-lg">
                No example projects yet
              </h2>
              <p className="mt-2 max-w-sm text-muted-foreground/75 text-sm leading-6">
                Example projects will appear here when a catalogue is configured for this instance.
              </p>
            </section>
          ) : (
            <div className="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
              {examples.map((example) => (
                <Card key={example.id}>
                  <CardHeader>
                    <CardTitle className="flex items-center gap-2">
                      <Code2 className="h-5 w-5" />
                      {example.name}
                    </CardTitle>
                    <CardDescription>{example.description}</CardDescription>
                  </CardHeader>
                  <CardContent>
                    <Button
                      className="w-full"
                      variant="outline"
                      onClick={() => window.open(example.repo, '_blank')}
                    >
                      View on GitHub
                    </Button>
                  </CardContent>
                </Card>
              ))}
            </div>
          )}
        </TabsContent>
      </Tabs>

      <CreateDatabaseModal
        isOpen={dbModalOpen}
        onOpenChange={setDbModalOpen}
        projectId={projectId}
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
    </div>
  );
}
