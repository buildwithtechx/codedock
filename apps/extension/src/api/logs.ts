import { cleanBaseUrl } from './client';

export function streamRuntimeLogs(
  url: string,
  token: string,
  appId: string,
  onLine: (line: string) => void,
  onError: (err: Error) => void
): AbortController {
  const controller = new AbortController();
  const base = cleanBaseUrl(url);
  const streamUrl = `${base}/api/apps/${encodeURIComponent(appId)}/runtime/logs/stream`;

  (async () => {
    try {
      const response = await fetch(streamUrl, {
        method: 'GET',
        headers: {
          Authorization: `Bearer ${token.trim()}`,
          Accept: 'text/event-stream',
        },
        signal: controller.signal,
      });

      if (!response.ok || !response.body) {
        throw new Error(`Log stream failed (${response.status}: ${response.statusText})`);
      }

      const reader = response.body.getReader();
      const decoder = new TextDecoder();
      let buffer = '';

      while (!controller.signal.aborted) {
        const { done, value } = await reader.read();
        if (done) break;

        buffer += decoder.decode(value, { stream: true });
        const lines = buffer.split('\n');
        buffer = lines.pop() ?? '';

        for (const line of lines) {
          const trimmed = line.trim();
          if (trimmed.startsWith('data:')) {
            const data = trimmed.replace(/^data:\s*/, '');
            onLine(data);
          } else if (trimmed && !trimmed.startsWith(':')) {
            onLine(trimmed);
          }
        }
      }
    } catch (err: unknown) {
      if (!controller.signal.aborted) {
        onError(err instanceof Error ? err : new Error(String(err)));
      }
    }
  })();

  return controller;
}
