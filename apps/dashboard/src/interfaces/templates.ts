export interface OneClickEnvVar {
  key: string;
  label: string;
  defaultValue?: string;
  secret: boolean;
  required: boolean;
  input: boolean;
}

export interface OneClickAppDetails {
  id: string;
  name: string;
  description: string;
  icon: string;
  category: string;
  dockerImage: string;
  defaultPort: number;
  services: string[];
  volumes: string[];
  envVariables: OneClickEnvVar[];
  verified: boolean;
}

export interface InstallAppInput {
  appId: string;
  projectId: string;
  environmentId?: string;
  name: string;
  secrets?: Record<string, string>;
  environment?: Record<string, string>;
  hostPort?: number;
  domain?: string;
  digest?: string;
}

export interface InstallPreviewService {
  service: string;
  image: string;
  env: string[];
  ports: string[];
}

export interface InstallPreviewVolume {
  service: string;
  name: string;
  target: string;
}

export interface InstallPreview {
  appId: string;
  name: string;
  services: InstallPreviewService[];
  volumes: InstallPreviewVolume[];
  generatedSecrets: string[];
  composeYaml: string;
  digest: string;
  kind: string;
}

export interface AppInstallResult {
  kind: string;
  stack?: { id: string; name: string; status: string };
}

export interface ExampleApp {
  id: string;
  name: string;
  description: string;
  repo: string;
  icon?: string;
  logo?: string;
}

export interface OneClickDeployRequest {
  appId: string;
  projectId: string;
  name: string;
  environmentId?: string;
  secrets?: Record<string, string>;
  environment?: Record<string, string>;
  hostPort?: number;
  domain?: string;
  digest?: string;
}

export interface OneClickDeployResponse {
  success: boolean;
  message: string;
  serviceId?: string;
  data?: AppInstallResult;
}

export interface ComposeDeployResponse {
  success: boolean;
  message: string;
  serviceIds?: string[];
}

export interface ArchiveDeployResponse {
  success: boolean;
  message: string;
  serviceId?: string;
}
