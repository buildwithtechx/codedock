import { listen } from '@tauri-apps/api/event';
import { invokeCommand } from './tauri-ipc.js';

export interface UpdateInfo {
  available: boolean;
  version: string;
  notes: string;
  date: string;
}

export interface UpdateProgress {
  state: string;
  chunk: number | null;
  total: number | null;
}

export const UPDATE_STATUS_EVENT = 'update-status';

export function checkUpdate(): Promise<UpdateInfo> {
  return invokeCommand<UpdateInfo>('check_update');
}

export function downloadAndInstall(): Promise<void> {
  return invokeCommand<void>('download_and_install');
}

export function restartApp(): Promise<void> {
  return invokeCommand<void>('restart_app');
}

export function onUpdateStatus(handler: (progress: UpdateProgress) => void): Promise<() => void> {
  return listen<UpdateProgress>(UPDATE_STATUS_EVENT, (event) => handler(event.payload));
}
