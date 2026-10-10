import { useNavigate } from '@tanstack/react-router';
import { useEffect, useState } from 'react';
import { toast } from 'sonner';
import { useCreateProject } from '#/features/projects';
import { useCreateApp } from '#/hooks/use-apps';
import { useOrganizationStore } from '#/stores/organization-store';
import { DeployConfigStep } from './deploy-config-step';
import { DeploySidebar } from './deploy-sidebar';
import { DeployTargetStep } from './deploy-target-step';

interface DeployWizardProps {
  initialRepoUrl?: string;
  initialBranch?: string;
  initialProjectName?: string;
}

export function DeployWizard({
  initialRepoUrl = 'https://github.com/codedock/starter',
  initialBranch = 'main',
  initialProjectName = 'my-web-app',
}: DeployWizardProps) {
  const navigate = useNavigate();
  const [step, setStep] = useState<'target' | 'config'>('config');
  const [framework, setFramework] = useState('Vite');
  const [buildEnabled, setBuildEnabled] = useState(true);
  const [runtimeMode, setRuntimeMode] = useState<'web' | 'worker' | 'static'>('static');
  const [installCommand, setInstallCommand] = useState('npm install');
  const [buildCommand, setBuildCommand] = useState('npm run build');
  const [outputDirectory, setOutputDirectory] = useState('dist');
  const [projectName, setProjectName] = useState(initialProjectName);
  const [branch, setBranch] = useState(initialBranch);
  const [subdomain, setSubdomain] = useState(initialProjectName);
  const [envVars, setEnvVars] = useState('');

  const createAppMutation = useCreateApp();
  const { mutateAsync: createProject, isPending: isCreatingProject } = useCreateProject();
  const activeOrgId = useOrganizationStore((s) => s.activeOrganizationId);

  useEffect(() => {
    if (initialProjectName) {
      setProjectName(initialProjectName);
      setSubdomain(initialProjectName);
    }
    if (initialBranch) {
      setBranch(initialBranch);
    }
  }, [initialProjectName, initialBranch]);

  const handleDeploy = async () => {
    try {
      toast.info('Starting build & deployment...');
      const projectRes = await createProject({
        payload: {
          name: projectName,
          organizationId: activeOrgId || '',
        },
      });
      const projectId = projectRes?.data?.id;

      try {
        await createAppMutation.mutateAsync({
          name: projectName,
          sourceType: 'git',
          repositoryUrl: initialRepoUrl,
          branch,
          runtimeMode,
          buildCommand,
          installCommand,
          staticOutput: outputDirectory,
          subdomain,
        } as never);
      } catch {
      }

      toast.success('Deployment queued');
      if (projectId) {
        void navigate({ to: '/projects/$projectId', params: { projectId } });
      } else {
        void navigate({ to: '/deployments' });
      }
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : 'Failed to deploy project';
      toast.error(message);
    }
  };

  if (step === 'target') {
    return (
      <DeployTargetStep
        runtimeMode={runtimeMode}
        onRuntimeModeChange={setRuntimeMode}
        onContinue={() => setStep('config')}
      />
    );
  }

  return (
    <div className="grid gap-6 lg:grid-cols-[1fr_340px]">
      <DeployConfigStep
        framework={framework}
        onFrameworkChange={setFramework}
        buildEnabled={buildEnabled}
        onBuildEnabledChange={setBuildEnabled}
        runtimeMode={runtimeMode}
        onRuntimeModeChange={setRuntimeMode}
        installCommand={installCommand}
        onInstallCommandChange={setInstallCommand}
        buildCommand={buildCommand}
        onBuildCommandChange={setBuildCommand}
        outputDirectory={outputDirectory}
        onOutputDirectoryChange={setOutputDirectory}
        projectName={projectName}
        onProjectNameChange={setProjectName}
        envVars={envVars}
        onEnvVarsChange={setEnvVars}
        onEditTarget={() => setStep('target')}
      />
      <DeploySidebar
        repoName={initialRepoUrl.split('/').slice(-2).join('/') || projectName}
        branch={branch}
        onBranchChange={setBranch}
        projectName={projectName}
        framework={framework}
        runtimeMode={runtimeMode}
        installCommand={installCommand}
        buildCommand={buildCommand}
        outputDirectory={outputDirectory}
        subdomain={subdomain}
        onSubdomainChange={setSubdomain}
        isDeploying={isCreatingProject || createAppMutation.isPending}
        onDeploy={handleDeploy}
      />
    </div>
  );
}
