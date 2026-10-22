export type RuntimeImageStatus = 'current' | 'update_available' | 'pinned' | 'unknown';

export interface RuntimeImageInfo {
  serviceId: string | null;
  serviceName: string | null;
  role: 'image' | 'base';
  ref: string;
  currentDigest: string | null;
  latestDigest: string | null;
  status: RuntimeImageStatus;
}

export interface RuntimeImageCheckResult {
  checkedAt: string;
  images: RuntimeImageInfo[];
}

export interface UpdateRuntimeInput {
  serviceIds?: string[];
}
