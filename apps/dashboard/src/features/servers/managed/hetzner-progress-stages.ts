import type { ManagedOperation } from '#/interfaces/managed';

export const hetznerStages = [
  { id: 'creating_server', label: 'Creating server', detail: 'Reserving the Hetzner Cloud server' },
  { id: 'waiting_ip', label: 'Waiting for boot', detail: 'Waiting for first boot and public IP' },
  {
    id: 'bootstrapping_cloud_init',
    label: 'Bootstrapping',
    detail: 'cloud-init installs Docker and hardens SSH',
  },
  { id: 'ssh_verified', label: 'SSH verified', detail: 'SSH reachable and Docker socket ready' },
];

export function hetznerStageIndex(snapshot: ManagedOperation | null): number {
  if (!snapshot) return -1;
  if (snapshot.status === 'COMPLETED' || snapshot.phase === 'READY') return hetznerStages.length;
  const logs = snapshot.logs ?? '';
  if (logs.includes('Waiting for docker') || logs.includes('Waiting for ssh')) return 2;
  if (logs.includes('to boot') || snapshot.phase === 'BOOTSTRAP') return 1;
  return 0;
}

export function hetznerFailed(snapshot: ManagedOperation | null): boolean {
  return snapshot !== null && (snapshot.status === 'FAILED' || snapshot.status === 'INTERRUPTED');
}
