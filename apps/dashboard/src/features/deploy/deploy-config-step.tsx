import { DeployComposeCard } from './deploy-compose-card';
import { DeployEnvSection } from './deploy-env-section';
import { DeployHealthChecksCard } from './deploy-health-checks-card';
import { DeployProjectNameCard } from './deploy-project-name-card';
import { DeployRuntimeCard } from './deploy-runtime-card';
import { DeployTopBar } from './deploy-top-bar';
import { MonorepoAppPicker } from './monorepo-app-picker';
import type { MonorepoCandidateApp } from '#/interfaces/git';

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
  rootDirectory?: string;
  onRootDirectoryChange?: (s: string) => void;
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
  repositoryUrl?: string;
  branch?: string;
  monorepoApps?: MonorepoCandidateApp[];
  selectedIndices?: number[];
  onToggleIndex?: (index: number) => void;
  onUpdateApp?: (index: number, patch: Partial<MonorepoCandidateApp>) => void;
  onAddCustomApp?: (app: MonorepoCandidateApp) => void;
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
  rootDirectory,
  onRootDirectoryChange,
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
  repositoryUrl,
  branch,
  monorepoApps,
  selectedIndices,
  onToggleIndex,
  onUpdateApp,
  onAddCustomApp,
}: DeployConfigStepProps) {
  const handleMigrateDraft = (cmds: Record<string, string>, vars: Record<string, string>) => {
    if (cmds.installCommand !== undefined) onInstallCommandChange(cmds.installCommand);
    if (cmds.buildCommand !== undefined) onBuildCommandChange(cmds.buildCommand);
    if (Object.keys(vars).length > 0) {
      const lines = envVars.trim() ? envVars.trim().split('\n') : [];
      for (const [k, v] of Object.entries(vars)) {
        lines.push(`${k}=${v}`);
      }
      onEnvVarsChange(lines.join('\n'));
    }
  };

  return (
    <div className="space-y-5">
      <DeployTopBar environmentName="Production" onEditTarget={onEditTarget} />

      {monorepoApps &&
        monorepoApps.length > 0 &&
        selectedIndices &&
        onToggleIndex &&
        onUpdateApp &&
        onAddCustomApp &&
        repositoryUrl && (
          <MonorepoAppPicker
            apps={monorepoApps}
            selectedIndices={selectedIndices}
            onToggleIndex={onToggleIndex}
            onUpdateApp={onUpdateApp}
            onAddCustomApp={onAddCustomApp}
            repositoryUrl={repositoryUrl}
            branch={branch}
          />
        )}

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
        rootDirectory={rootDirectory}
        onRootDirectoryChange={onRootDirectoryChange}
        onMigrateDraft={handleMigrateDraft}
        repositoryUrl={repositoryUrl}
        branch={branch}
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
