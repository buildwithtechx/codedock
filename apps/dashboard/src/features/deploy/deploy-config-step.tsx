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
  composePath?: string;
  onComposePathChange?: (p: string) => void;
  waitForPort?: boolean;
  onWaitForPortChange?: (v: boolean) => void;
  watchCrashLoop?: boolean;
  onWatchCrashLoopChange?: (v: boolean) => void;
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
  composePath,
  onComposePathChange,
  waitForPort,
  onWaitForPortChange,
  watchCrashLoop,
  onWatchCrashLoopChange,
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

      <DeployComposeCard composePath={composePath} onComposePathChange={onComposePathChange} />

      <DeployHealthChecksCard
        waitForPort={waitForPort}
        onWaitForPortChange={onWaitForPortChange}
        watchCrashLoop={watchCrashLoop}
        onWatchCrashLoopChange={onWatchCrashLoopChange}
      />

      <DeployProjectNameCard projectName={projectName} onProjectNameChange={onProjectNameChange} />
    </div>
  );
}
