export function runtimeLogStreamUrl(serviceId: string, pod?: string): string {
  const base = `/apps/${serviceId}/runtime/logs/stream`;
  return pod ? `${base}?pod=${encodeURIComponent(pod)}` : base;
}

export function runtimeExecTerminalUrl(serviceId: string): string {
  return `/apps/${serviceId}/runtime/exec-terminal`;
}

export function openRuntimeSocket(path: string): WebSocket {
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
  return new WebSocket(`${protocol}//${window.location.host}${path}`);
}
