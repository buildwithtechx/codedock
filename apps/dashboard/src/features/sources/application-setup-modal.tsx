import { Button } from '#/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '#/components/ui/dialog';
import { Label } from '#/components/ui/label';
import { ApplicationSetupFields } from './application-setup-fields';
import { type ApplicationSetupProps, parseSetupVariables } from './application-setup-types';
import { DeploymentProgress } from './deployment-progress';
import { DestinationPicker } from './destination-picker';
import { RepositorySelector } from './repository-selector';
import { useApplicationSetup } from './use-application-setup';

export function ApplicationSetupModal(props: ApplicationSetupProps) {
  const setup = useApplicationSetup(props);
  return (
    <Dialog open={props.isOpen} onOpenChange={props.onOpenChange}>
      <DialogContent className="max-h-[90vh] overflow-y-auto sm:max-w-3xl">
        <DialogHeader>
          <DialogTitle>Set up an application</DialogTitle>
          <DialogDescription>
            Configure its source, environment and runtime, then review and deploy.
          </DialogDescription>
        </DialogHeader>
        {setup.deploymentId ? (
          <DeploymentProgress deploymentId={setup.deploymentId} />
        ) : (
          <>
            <fieldset disabled={setup.locked} className="space-y-5">
              <div className="grid gap-4 sm:grid-cols-3">
                <div className="space-y-2">
                  <Label htmlFor="setup-project">Project</Label>
                  <select
                    id="setup-project"
                    className="w-full rounded-md border bg-background p-2"
                    value={setup.projectId}
                    onChange={(event) => setup.setProjectId(event.target.value)}
                  >
                    {!setup.projects.data?.data?.some(
                      (project) => project.id === setup.projectId
                    ) && <option value={setup.projectId}>Current project</option>}
                    {setup.projects.data?.data?.map((project) => (
                      <option key={project.id} value={project.id}>
                        {project.name}
                      </option>
                    ))}
                  </select>
                </div>
                <div className="space-y-2">
                  <Label htmlFor="setup-environment">Environment</Label>
                  <select
                    id="setup-environment"
                    className="w-full rounded-md border bg-background p-2"
                    value={setup.environmentId}
                    onChange={(event) => setup.setEnvironment(event.target.value)}
                  >
                    <option value="" disabled>
                      Select environment
                    </option>
                    {setup.environments.data?.data?.map((env) => (
                      <option key={env.id} value={env.id}>
                        {env.name}
                      </option>
                    ))}
                  </select>
                </div>
                <div className="space-y-2">
                  <Label htmlFor="setup-source">Source</Label>
                  <select
                    id="setup-source"
                    className="w-full rounded-md border bg-background p-2"
                    value={setup.source}
                    onChange={(event) => setup.setSource(event.target.value as 'git' | 'image')}
                  >
                    <option value="git">Git repository</option>
                    <option value="image">Container image</option>
                  </select>
                </div>
              </div>
              <p className="text-muted-foreground text-sm">
                Target: {setup.destinationSummary}. Docker uses the project server binding; cluster
                and native destinations are applied before the first deployment.
              </p>
              <DestinationPicker
                destination={setup.destination}
                serverId={setup.serverId}
                clusters={setup.clusters.data?.data ?? []}
                registries={setup.registries.data?.data ?? []}
                onChange={setup.setDestination}
              />
              {!setup.review && (
                <>
                  {setup.source === 'git' && (
                    <RepositorySelector
                      onSelect={(url, branch) => {
                        setup.update('repositoryUrl', url);
                        setup.update('branch', branch);
                      }}
                    />
                  )}
                  <ApplicationSetupFields
                    draft={setup.draft}
                    source={setup.source}
                    update={setup.update}
                    variables={setup.variables}
                    setVariables={setup.setVariables}
                  />
                  {setup.source === 'git' && (
                    <Button
                      variant="outline"
                      disabled={setup.inspect.isPending || !setup.draft.repositoryUrl}
                      onClick={() => setup.inspect.mutate()}
                    >
                      {setup.inspect.isPending ? 'Inspecting…' : 'Detect repository defaults'}
                    </Button>
                  )}
                  {setup.detection && (
                    <p className="text-sm">
                      Detected: {setup.detection.framework} {setup.detection.packageManager}.{' '}
                      {setup.detection.warnings.join(' ')}
                    </p>
                  )}
                </>
              )}
              <a className="text-primary text-sm" href={`/projects/${setup.projectId}/compose`}>
                Set up a Compose stack in this project
              </a>
            </fieldset>
            {setup.review && (
              <div className="space-y-3">
                <p className="text-sm">
                  Creates one application and saves its variables. Deploy also builds or pulls its
                  image, checks readiness and activates routing. Retry continues the same
                  application if setup is interrupted.
                </p>
                <pre className="max-h-72 overflow-auto rounded-md bg-muted p-3 text-xs">
                  {JSON.stringify(
                    {
                      application: setup.payload,
                      destination: setup.runtimeTarget,
                      variables: parseSetupVariables(setup.variables),
                    },
                    null,
                    2
                  )}
                </pre>
              </div>
            )}
            {setup.error && (
              <p role="alert" className="text-destructive text-sm">
                {setup.error}
              </p>
            )}
            <div className="flex justify-end gap-2">
              {setup.review ? (
                <>
                  <Button variant="outline" disabled={setup.locked} onClick={setup.edit}>
                    Edit
                  </Button>
                  <Button
                    variant="outline"
                    disabled={setup.apply.isPending}
                    onClick={() => setup.apply.mutate(false)}
                  >
                    Save configuration
                  </Button>
                  <Button disabled={setup.apply.isPending} onClick={() => setup.apply.mutate(true)}>
                    {setup.apply.isPending ? 'Applying…' : 'Save and deploy'}
                  </Button>
                </>
              ) : (
                <Button
                  disabled={setup.inspect.isPending || setup.locked}
                  onClick={setup.prepareReview}
                >
                  Review setup
                </Button>
              )}
            </div>
          </>
        )}
      </DialogContent>
    </Dialog>
  );
}
