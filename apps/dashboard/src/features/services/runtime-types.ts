export type RuntimeVolume = {
  name: string;
  mountPath: string;
  sizeGiB: number;
  storageClass: string;
  shared: boolean;
};
export type RuntimeTarget = {
  kind: 'docker' | 'kubernetes' | 'bare';
  bareNode?: import('../servers/cluster-types').ClusterNode;
  bareReleaseUrl?: string;
  bareSha256?: string;
  bareCommand?: string[];
  bareRepoUrl?: string;
  bareBranch?: string;
  bareToolchain?: 'go' | 'node' | 'python' | 'static';
  bareInstallCommand?: string;
  bareBuildCommand?: string;
  bareOutput?: string;
  bareStaticDir?: string;
  clusterId?: string;
  nodeIds: string[];
  imageRepository?: string;
  registryId?: string;
  volumes: RuntimeVolume[];
};
export type ServiceRuntime = {
  serviceId: string;
  target: RuntimeTarget;
  revision: number;
  status: string;
  error: string;
};
export type RuntimeReview = {
  operation: { id: string; status: string; effects: string; expiresAt: number };
  confirmation: string;
};
export type WorkloadObservation = {
  kind: string;
  status: string;
  desired: number;
  available: number;
  metricsAvailable: boolean;
  metricsError?: string;
  pods: {
    name: string;
    node: string;
    ready: boolean;
    phase: string;
    restarts: number;
    cpu: number;
    memoryBytes: number;
  }[];
};
