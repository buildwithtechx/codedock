import { invokeCommand } from './tauri-ipc.js';

export interface DaemonConnection {
  daemonUrl: string;
  daemonPort: number;
  hasToken: boolean;
}

export const DEFAULT_DAEMON_URL = 'http://localhost:8080';
export const DEFAULT_DAEMON_PORT = 8080;

export function defaultConnection(): DaemonConnection {
  return { daemonUrl: DEFAULT_DAEMON_URL, daemonPort: DEFAULT_DAEMON_PORT, hasToken: false };
}

export function isDefaultConnection(connection: DaemonConnection): boolean {
  return (
    connection.daemonUrl === DEFAULT_DAEMON_URL &&
    connection.daemonPort === DEFAULT_DAEMON_PORT &&
    connection.hasToken === false
  );
}

export function getConnection(): Promise<DaemonConnection> {
  return invokeCommand<DaemonConnection>('get_connection');
}

export function setDaemonUrl(url: string): Promise<void> {
  return invokeCommand<void>('set_daemon_url', { url });
}

export function setDaemonPort(port: number): Promise<void> {
  return invokeCommand<void>('set_daemon_port', { port });
}

export function saveToken(token: string): Promise<void> {
  return invokeCommand<void>('save_token', { token });
}

export function loadToken(): Promise<string> {
  return invokeCommand<string>('load_token');
}

export function clearToken(): Promise<void> {
  return invokeCommand<void>('clear_token');
}
