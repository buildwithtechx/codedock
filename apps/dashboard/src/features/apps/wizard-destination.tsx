import { useEffect } from 'react';
import { Label } from '#/components/ui/label';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '#/components/ui/select';
import { useListAllProjects } from '#/features/projects';
import { useListByProject } from '#/hooks/use-environments';

export function WizardDestination({
  projectId,
  environmentId,
  onProjectChange,
  onEnvironmentChange,
  disabled = false,
}: {
  projectId: string;
  environmentId: string;
  onProjectChange: (projectId: string) => void;
  onEnvironmentChange: (environmentId: string) => void;
  disabled?: boolean;
}) {
  const projectsQuery = useListAllProjects();
  const activeProjectId = projectId === 'standalone' ? '' : projectId;
  const environmentsQuery = useListByProject(activeProjectId);
  const projects = projectsQuery.data ?? [];
  const environments = environmentsQuery.data?.data ?? [];

  useEffect(() => {
    if (!projectId && !projectsQuery.isLoading) {
      if (projects.length > 0) {
        onProjectChange(projects[0].id);
      } else {
        onProjectChange('standalone');
      }
    }
  }, [projectId, projects, projectsQuery.isLoading, onProjectChange]);

  return (
    <div className="space-y-4">
      <div className="space-y-2">
        <Label htmlFor="wizard-project">Project</Label>
        <Select
          value={projectId || 'standalone'}
          onValueChange={onProjectChange}
          disabled={disabled}
        >
          <SelectTrigger id="wizard-project" className="w-full">
            <SelectValue placeholder="Select destination" />
          </SelectTrigger>
          <SelectContent>
            {projects.map((project) => (
              <SelectItem key={project.id} value={project.id}>
                {project.name}
              </SelectItem>
            ))}
            <SelectItem value="standalone">
              {projects.length === 0
                ? 'Dedicated Project (Auto-created)'
                : '+ Dedicated Project (Auto-created)'}
            </SelectItem>
          </SelectContent>
        </Select>
        {projectId === 'standalone' && (
          <p className="text-[11px] text-muted-foreground">
            A dedicated project will be auto-created for this app upon deployment.
          </p>
        )}
        {projectsQuery.isError && (
          <p role="alert" className="text-destructive text-xs">
            Projects could not be loaded.
          </p>
        )}
      </div>
      {projectId !== 'standalone' && (
        <div className="space-y-2">
          <Label htmlFor="wizard-environment">Environment</Label>
          <Select
            value={environmentId}
            onValueChange={onEnvironmentChange}
            disabled={disabled || !projectId}
          >
            <SelectTrigger id="wizard-environment" className="w-full">
              <SelectValue placeholder="Default environment" />
            </SelectTrigger>
            <SelectContent>
              {environments.map((environment) => (
                <SelectItem key={environment.id} value={environment.id}>
                  {environment.name}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
          {projectId && environmentsQuery.isError && (
            <p role="alert" className="text-destructive text-xs">
              Environments could not be loaded.
            </p>
          )}
        </div>
      )}
    </div>
  );
}
