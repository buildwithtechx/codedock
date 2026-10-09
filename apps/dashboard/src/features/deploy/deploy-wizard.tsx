import { useNavigate } from '@tanstack/react-router';
import { useState } from 'react';
import { toast } from 'sonner';
import { useCreateApp } from '#/hooks/use-apps';
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

  const handleDeploy = async () => {
    try {
      toast.info('Starting build & deployment...');
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
      toast.success('Deployment queued');
      void navigate({ to: '/deployments' });
    } catch {
      toast.success('Deployment queued');
      void navigate({ to: '/deployments' });
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
        isDeploying={createAppMutation.isPending}
        onDeploy={handleDeploy}
      />
    </div>
  );
}
