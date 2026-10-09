import { useEffect, useMemo, useState } from 'react';
import { env } from '#/env';
import type { ServerMetrics } from '#/interfaces/server';
import { useAuthStore } from '#/stores/auth-store';

export interface LiveMetricPoint {
  time: number;
  cpu: number;
  memory: number;
  disk: number;
}

export interface LiveMetricsSnapshot {
  cpu: number;
  memory: number;
  disk: number;
  memoryUsedBytes: number;
  memoryLimitBytes: number;
  diskUsedBytes: number;
  diskTotalBytes: number;
}

function toFiniteNumber(value: unknown): number {
  return typeof value === 'number' && Number.isFinite(value) ? value : 0;
}

function snapshotOf(payload: Partial<ServerMetrics>): LiveMetricsSnapshot {
  const memoryLimit = toFiniteNumber(payload.memory_limit_bytes);
  const memoryUsed = toFiniteNumber(payload.memory_usage_bytes);
  const diskTotal = toFiniteNumber(payload.disk_total_bytes);
  const diskUsed = toFiniteNumber(payload.disk_usage_bytes);
  return {
    cpu: toFiniteNumber(payload.cpu_usage_percentage),
    memory: memoryLimit > 0 ? (memoryUsed / memoryLimit) * 100 : 0,
    disk: diskTotal > 0 ? (diskUsed / diskTotal) * 100 : 0,
    memoryUsedBytes: memoryUsed,
    memoryLimitBytes: memoryLimit,
    diskUsedBytes: diskUsed,
    diskTotalBytes: diskTotal,
  };
}

export function formatBytes(bytes: number): string {
  if (bytes === 0) return '0 B';
  const units = ['B', 'KB', 'MB', 'GB', 'TB'];
  const i = Math.floor(Math.log(bytes) / Math.log(1024));
  const value = bytes / 1024 ** i;
  return `${value.toFixed(i > 1 ? 1 : 0)} ${units[i]}`;
}

export function useServerLiveMetrics(serverId: string | undefined) {
  const [points, setPoints] = useState<LiveMetricPoint[]>([]);
  const [latest, setLatest] = useState<LiveMetricsSnapshot | null>(null);
  const [connected, setConnected] = useState(false);
  const [generation, setGeneration] = useState(0);

  useEffect(() => {
    if (!serverId) return;
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsHost = env.VITE_API_URL.replace(/^http(s?):\/\//, '');
    const socket = new WebSocket(
      `${protocol}//${wsHost}/api/ws/servers/${serverId}/metrics`,
      (() => {
        const token = useAuthStore.getState().token;
        return token ? ['auth', token] : undefined;
      })()
    );

    socket.onopen = () => setConnected(true);
    socket.onclose = () => setConnected(false);
    socket.onerror = () => setConnected(false);
    socket.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data) as Partial<ServerMetrics>;
        const snapshot = snapshotOf(data);
        setLatest(snapshot);
        setPoints((prev) => {
          const next = [
            ...prev,
            { time: Date.now(), cpu: snapshot.cpu, memory: snapshot.memory, disk: snapshot.disk },
          ];
          if (next.length > 30) next.shift();
          return next;
        });
      } catch {
        setConnected(false);
      }
    };

    return () => socket.close();
  }, [serverId, generation]);

  const cpuSeries = useMemo(
    () => points.map((point) => ({ time: point.time, value: point.cpu })),
    [points]
  );
  const memorySeries = useMemo(
    () => points.map((point) => ({ time: point.time, value: point.memory })),
    [points]
  );
  const diskSeries = useMemo(
    () => points.map((point) => ({ time: point.time, value: point.disk })),
    [points]
  );

  return {
    points,
    snapshot: latest,
    cpuSeries,
    memorySeries,
    diskSeries,
    connected,
    reconnect: () => setGeneration((value) => value + 1),
  };
}

export function snapshotFromStored(metrics: ServerMetrics | null): LiveMetricsSnapshot | null {
  if (!metrics) return null;
  return snapshotOf(metrics);
}
