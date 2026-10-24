/** Freshness status of a runtime container image. */
export type RuntimeImageStatus = 'current' | 'update_available' | 'pinned' | 'unknown';

/** Version metadata and status for a service or base container image. */
export interface RuntimeImageInfo {
  serviceId: string | null;
  serviceName: string | null;
  role: 'image' | 'base';
  ref: string;
  currentDigest: string | null;
  latestDigest: string | null;
  status: RuntimeImageStatus;
}

/** Result returned by a runtime image freshness inspection. */
export interface RuntimeImageCheckResult {
  checkedAt: string;
  images: RuntimeImageInfo[];
}

/** Input to trigger runtime base image updates and redeployments. */
export interface UpdateRuntimeInput {
  serviceIds?: string[];
}
