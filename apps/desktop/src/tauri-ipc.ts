import { invoke as tauriInvoke } from '@tauri-apps/api/core';

export type InvokeFn = (cmd: string, args?: Record<string, unknown>) => Promise<unknown>;

let currentInvoke: InvokeFn = tauriInvoke as InvokeFn;

export function setInvokeForTesting(fn: InvokeFn | null): void {
  currentInvoke = fn ?? (tauriInvoke as InvokeFn);
}

export async function invokeCommand<T>(cmd: string, args?: Record<string, unknown>): Promise<T> {
  return (await currentInvoke(cmd, args)) as T;
}
