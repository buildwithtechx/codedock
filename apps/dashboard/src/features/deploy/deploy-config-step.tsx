import { DeployComposeCard } from './deploy-compose-card';
import { DeployEnvSection } from './deploy-env-section';
import { DeployHealthChecksCard } from './deploy-health-checks-card';
import { DeployProjectNameCard } from './deploy-project-name-card';
import { DeployRuntimeCard } from './deploy-runtime-card';
import { DeployTopBar } from './deploy-top-bar';

interface DeployConfigStepProps {
  framework: string;
  onFrameworkChange: (f: string) => void;
  exposedPort: number;
  onExposedPortChange: (p: number) => void;
  installCommand: string;
  onInstallCommandChange: (s: string) => void;
  buildCommand: string;
  onBuildCommandChange: (s: string) => void;
  outputDirectory: string;
  onOutputDirectoryChange: (s: string) => void;
  projectName: string;
  onProjectNameChange: (s: string) => void;
  envVars: string;
  onEnvVarsChange: (s: string) => void;
  onEditTarget: () => void;
}

export function DeployConfigStep({
  framework,
  onFrameworkChange,
  exposedPort,
  onExposedPortChange,
  installCommand,
  onInstallCommandChange,
  buildCommand,
  onBuildCommandChange,
  outputDirectory,
  onOutputDirectoryChange,
  projectName,
  onProjectNameChange,
  envVars,
  onEnvVarsChange,
  onEditTarget,
}: DeployConfigStepProps) {
  return (
    <div className="space-y-5">
      <DeployTopBar environmentName="Production" onEditTarget={onEditTarget} />

      <DeployRuntimeCard
        framework={framework}
        onFrameworkChange={onFrameworkChange}
        exposedPort={exposedPort}
        onExposedPortChange={onExposedPortChange}
        installCommand={installCommand}
        onInstallCommandChange={onInstallCommandChange}
        buildCommand={buildCommand}
        onBuildCommandChange={onBuildCommandChange}
        outputDirectory={outputDirectory}
        onOutputDirectoryChange={onOutputDirectoryChange}
      />

      <DeployEnvSection envVars={envVars} onEnvVarsChange={onEnvVarsChange} />

      <DeployComposeCard />

      <DeployHealthChecksCard />

      <DeployProjectNameCard projectName={projectName} onProjectNameChange={onProjectNameChange} />
    </div>
  );
}
