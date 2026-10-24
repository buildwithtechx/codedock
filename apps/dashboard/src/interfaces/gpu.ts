export interface ServerGpuDevice {
  index: number;
  name: string;
  uuid: string;
  memoryTotalMb: number;
  memoryFreeMb?: number;
  driverVersion?: string;
}

export interface ServiceGpuConfig {
  enabled: boolean;
  count: number;
  deviceUuids?: string[];
  capabilities?: string[];
}
