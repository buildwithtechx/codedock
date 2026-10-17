import { useEffect, useState } from 'react';
import type { ManagedOperation } from '#/interfaces/managed';
import { getApiBaseUrl } from '#/lib/api-client';
import { getAuthHeaders } from '#/lib/auth-refresh';
import { managedService } from '#/services/managed';

const terminalStatuses = new Set(['COMPLETED', 'FAILED', 'INTERRUPTED']);
const pollIntervalMs = 3000;

function parseFrames(buffer: string): { frames: { event: string; data: string }[]; rest: string } {
  const frames: { event: string; data: string }[] = [];
  const chunks = buffer.split('\n\n');
  const rest = chunks.pop() ?? '';
  for (const chunk of chunks) {
    let event = '';
    const data: string[] = [];
    for (const line of chunk.split('\n')) {
      if (line.startsWith('event:')) event = line.slice(6).trim();
      else if (line.startsWith('data:')) data.push(line.slice(5).trimStart());
    }
    if (event) frames.push({ event, data: data.join('\n') });
  }
  return { frames, rest };
}

export function useOperationStream(operationId: string | null) {
  const [snapshot, setSnapshot] = useState<ManagedOperation | null>(null);
  const [streaming, setStreaming] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!operationId) return;
    let cancelled = false;
    const controller = new AbortController();
    let pollTimer: ReturnType<typeof setInterval> | null = null;

    const stopPolling = () => {
      if (pollTimer) {
        clearInterval(pollTimer);
        pollTimer = null;
      }
    };

    const pollUntilDone = async () => {
      stopPolling();
      const tick = async () => {
        if (cancelled) return;
        try {
          const op = await managedService.getOperation(operationId);
          if (cancelled) return;
          setSnapshot(op);
          if (terminalStatuses.has(op.status)) stopPolling();
        } catch (err) {
          if (cancelled) return;
          setError(err instanceof Error ? err.message : 'Failed to load operation');
          stopPolling();
        }
      };
      await tick();
      if (!cancelled) pollTimer = setInterval(() => void tick(), pollIntervalMs);
    };

    const run = async () => {
      setStreaming(true);
      setError(null);
      try {
        const response = await fetch(
          `${getApiBaseUrl()}/managed-operations/${operationId}/stream`,
          { headers: getAuthHeaders(), credentials: 'include', signal: controller.signal }
        );
        if (!response.ok || !response.body) {
          await pollUntilDone();
          return;
        }
        const reader = response.body.getReader();
        const decoder = new TextDecoder();
        let buffer = '';
        let finished = false;
        while (!finished) {
          const { done, value } = await reader.read();
          if (cancelled) {
            reader.cancel().catch(() => undefined);
            return;
          }
          buffer += decoder.decode(value ?? new Uint8Array(), { stream: !done });
          const { frames, rest } = parseFrames(buffer);
          buffer = rest;
          for (const frame of frames) {
            if (frame.event === 'progress' || frame.event === 'done') {
              try {
                const op = JSON.parse(frame.data) as ManagedOperation;
                setSnapshot(op);
                if (frame.event === 'done' || terminalStatuses.has(op.status)) finished = true;
              } catch {
                setError('Received an unreadable operation update');
                finished = true;
              }
            } else if (frame.event === 'error') {
              setError(frame.data || 'Operation stream failed');
              finished = true;
            }
          }
          if (done) break;
        }
        if (!cancelled && !finished) {
          await pollUntilDone();
          return;
        }
      } catch (_err) {
        if (!cancelled && !controller.signal.aborted) await pollUntilDone();
      } finally {
        if (!cancelled) setStreaming(false);
      }
    };

    void run();
    return () => {
      cancelled = true;
      stopPolling();
      controller.abort();
    };
  }, [operationId]);

  return { snapshot, streaming, error };
}
