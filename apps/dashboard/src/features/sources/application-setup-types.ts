import type { CreateAppServiceRequest } from '#/features/services';

export interface ApplicationSetupProps {
  isOpen: boolean;
  onOpenChange: (open: boolean) => void;
  projectId: string;
  initialSource?: 'git' | 'image';
}

export interface RepositoryInspection {
  framework: string;
  packageManager: string;
  installCommand: string;
  buildCommand: string;
  startCommand: string;
  internalPort: number;
  staticOutput: string;
  buildEngine: CreateAppServiceRequest['buildEngine'];
  dockerfilePath: string;
  warnings: string[];
}

export const setupDefaults: Omit<CreateAppServiceRequest, 'projectId'> = {
  name: '',
  repositoryUrl: '',
  imageRef: '',
  branch: 'main',
  rootDirectory: '/',
  runtimeMode: 'web',
  installCommand: '',
  buildCommand: '',
  startCommand: '',
  dockerfilePath: 'Dockerfile',
  buildEngine: 'nixpacks',
  internalPort: 3000,
  domain: '',
  staticOutput: '',
  healthCheckPath: '',
};

export function parseSetupVariables(text: string) {
  const variables = new Map<string, string>();
  for (const line of text.split('\n')) {
    if (!line.trim()) continue;
    const separator = line.indexOf('=');
    const key = line.slice(0, separator).trim();
    if (separator < 1 || !/^[A-Za-z_][A-Za-z0-9_]*$/.test(key))
      throw new Error('Variables must use KEY=value, one per line.');
    if (variables.has(key)) throw new Error(`Duplicate variable: ${key}`);
    variables.set(key, line.slice(separator + 1));
  }
  return Array.from(variables, ([key, value]) => ({ key, value, isSecret: true }));
}
