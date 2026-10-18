import { useQueryClient } from '@tanstack/react-query';
import { useNavigate } from '@tanstack/react-router';
import { useEffect, useState } from 'react';
import { toast } from 'sonner';
import { useCreateApp } from '#/hooks/use-apps';
import { useTrigger } from '#/hooks/use-deployments';
import { projectEnvironmentsQueryOptions } from '#/hooks/use-environments';
import { useCreateProject, useSetVars } from '#/hooks/use-projects';
import { useOrganizationStore } from '#/stores/organization-store';
import { DeployConfigStep } from './deploy-config-step';
import { DeploySidebar } from './deploy-sidebar';
import { DeployTargetStep } from './deploy-target-step';

interface DeployWizardProps {
  initialRepoUrl?: string;
  initialBranch?: string;
  initialProjectName?: string;
  initialDir?: string;
}

export function DeployWizard({
  initialRepoUrl = 'https://github.com/codedock/starter',
  initialBranch = 'main',
  initialProjectName = 'my-web-app',
  initialDir = '',
}: DeployWizardProps) {
  const navigate = useNavigate();
  const [step, setStep] = useState<'target' | 'config'>('config');
  const [framework, setFramework] = useState('Dockerfile');
  const runtimeMode: 'web' | 'static' = framework.toLowerCase() === 'dockerfile' ? 'web' : 'static';
  const [installCommand, setInstallCommand] = useState('npm install');
  const [buildCommand, setBuildCommand] = useState('npm run build');
  const [outputDirectory, setOutputDirectory] = useState('dist');
  const [rootDirectory, setRootDirectory] = useState(initialDir);
  const [projectName, setProjectName] = useState(initialProjectName);
  const [branch, setBranch] = useState(initialBranch);
  const [subdomain, setSubdomain] = useState(initialProjectName);
  const [exposedPort, setExposedPort] = useState(8080);
  const [envVars, setEnvVars] = useState('');
  const [composePath, setComposePath] = useState('');
  const [waitForPort, setWaitForPort] = useState(false);
  const [watchCrashLoop, setWatchCrashLoop] = useState(false);
  const [isDeploying, setIsDeploying] = useState(false);

  const queryClient = useQueryClient();
  const { mutateAsync: createProject } = useCreateProject();
  const setVarsMutation = useSetVars();
  const createAppMutation = useCreateApp();
  const triggerDeploymentMutation = useTrigger();
  const activeOrgId = useOrganizationStore((s) => s.activeOrganizationId);

  useEffect(() => {
    if (initialProjectName) {
      setProjectName(initialProjectName);
      setSubdomain(initialProjectName);
    }
    if (initialBranch) {
      setBranch(initialBranch);
    }
    if (initialDir !== undefined) {
      setRootDirectory(initialDir);
    }
  }, [initialProjectName, initialBranch, initialDir]);

  const handleDeploy = async () => {
    setIsDeploying(true);
    try {
      toast.info('Starting build & deployment...');
      const projectRes = await createProject({
        payload: {
          name: projectName,
          organizationId: activeOrgId || '',
          gitUrl: initialRepoUrl,
          gitBranch: branch,
          gitProvider: 'github',
        },
      });
      const projectId = projectRes?.data?.id;

      if (!projectId) {
        throw new Error('Project creation failed: missing project ID');
      }

      if (envVars.trim()) {
        const parsedVars: Record<string, string> = {};
        for (const line of envVars.split('\n')) {
          const trimmed = line.trim();
          if (!trimmed || trimmed.startsWith('#')) continue;
          const idx = trimmed.indexOf('=');
          if (idx > 0) {
            parsedVars[trimmed.slice(0, idx).trim()] = trimmed.slice(idx + 1).trim();
          }
        }
        if (Object.keys(parsedVars).length > 0) {
          try {
            await setVarsMutation.mutateAsync({
              id: projectId,
              payload: { variables: parsedVars },
            });
          } catch {}
        }
      }

      const envsRes = await queryClient.fetchQuery(projectEnvironmentsQueryOptions(projectId));
      const defaultEnv = envsRes?.data?.[0];
      if (defaultEnv?.id) {
        const appRes = await createAppMutation.mutateAsync({
          environmentId: defaultEnv.id,
          payload: {
            projectId,
            name: projectName,
            repositoryUrl: initialRepoUrl,
            branch,
            runtimeMode,
            installCommand,
            buildCommand,
            startCommand: '',
            staticOutput: outputDirectory,
            domain: subdomain,
            buildEngine: runtimeMode === 'static' ? 'static' : 'nixpacks',
            internalPort: exposedPort,
            rootDirectory: rootDirectory || '',
            dockerfilePath: '',
            healthCheckPath: '',
          },
        });
        const serviceId = appRes?.data?.id;
        if (serviceId) {
          try {
            await triggerDeploymentMutation.mutateAsync({ serviceId });
          } catch {}
        }
      }

      toast.success('Deployment queued');
      void navigate({ to: '/projects/$projectId', params: { projectId } });
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : 'Failed to deploy project';
      toast.error(message);
    } finally {
      setIsDeploying(false);
    }
  };

  if (step === 'target') {
    return <DeployTargetStep onContinue={() => setStep('config')} />;
  }

  return (
    <div className="grid gap-5 lg:grid-cols-[1fr_340px]">
      <DeployConfigStep
        framework={framework}
        onFrameworkChange={setFramework}
        exposedPort={exposedPort}
        onExposedPortChange={setExposedPort}
        installCommand={installCommand}
        onInstallCommandChange={setInstallCommand}
        buildCommand={buildCommand}
        onBuildCommandChange={setBuildCommand}
        outputDirectory={outputDirectory}
        onOutputDirectoryChange={setOutputDirectory}
        rootDirectory={rootDirectory}
        onRootDirectoryChange={setRootDirectory}
        projectName={projectName}
        onProjectNameChange={setProjectName}
        envVars={envVars}
        onEnvVarsChange={setEnvVars}
        composePath={composePath}
        onComposePathChange={setComposePath}
        waitForPort={waitForPort}
        onWaitForPortChange={setWaitForPort}
        watchCrashLoop={watchCrashLoop}
        onWatchCrashLoopChange={setWatchCrashLoop}
        onEditTarget={() => setStep('target')}
      />
      <DeploySidebar
        repoName={initialRepoUrl.split('/').slice(-2).join('/') || projectName}
        branch={branch}
        onBranchChange={setBranch}
        projectName={projectName}
        framework={framework}
        subdomain={subdomain}
        onSubdomainChange={setSubdomain}
        exposedPort={exposedPort}
        isDeploying={isDeploying}
        onDeploy={handleDeploy}
      />
    </div>
  );
}
