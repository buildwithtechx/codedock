import { listen } from '@tauri-apps/api/event';
import { invokeCommand } from './tauri-ipc.js';

export interface SidecarStatus {
  running: boolean;
  healthy: boolean;
  restarts: number;
  pid: number | null;
  lastError: string | null;
  port: number;
}

export const SIDECAR_STATUS_EVENT = 'sidecar-status';

export function getSidecarStatus(): Promise<SidecarStatus> {
  return invokeCommand<SidecarStatus>('sidecar_status');
}

export function restartSidecar(): Promise<void> {
  return invokeCommand<void>('sidecar_restart');
}

export function stopSidecar(): Promise<void> {
  return invokeCommand<void>('sidecar_stop');
}

export function startSidecar(): Promise<void> {
  return invokeCommand<void>('sidecar_start');
}

export function setSidecarPort(port: number): Promise<void> {
  return invokeCommand<void>('sidecar_set_port', { port });
}

export function onSidecarStatus(handler: (status: SidecarStatus) => void): Promise<() => void> {
  return listen<SidecarStatus>(SIDECAR_STATUS_EVENT, (event) => handler(event.payload));
}
