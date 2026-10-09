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
  const environmentsQuery = useListByProject(projectId);
  const projects = projectsQuery.data ?? [];
  const environments = environmentsQuery.data?.data ?? [];

  return (
    <div className="space-y-4">
      <div className="space-y-2">
        <Label htmlFor="wizard-project">Project *</Label>
        <Select value={projectId} onValueChange={onProjectChange} disabled={disabled}>
          <SelectTrigger id="wizard-project" className="w-full">
            <SelectValue placeholder="Select a project" />
          </SelectTrigger>
          <SelectContent>
            {projects.map((project) => (
              <SelectItem key={project.id} value={project.id}>
                {project.name}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
        {projectsQuery.isError && (
          <p role="alert" className="text-destructive text-xs">
            Projects could not be loaded.
          </p>
        )}
      </div>
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
    </div>
  );
}
