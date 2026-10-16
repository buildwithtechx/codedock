export interface CodedockClientOptions {
  baseUrl?: string;
  token?: string;
  fetch?: typeof fetch;
}

export class CodedockApiError extends Error {
  constructor(
    public status: number,
    public statusText: string,
    public body: unknown
  ) {
    super(`API error ${status} (${statusText}): ${typeof body === 'string' ? body : JSON.stringify(body)}`);
    this.name = 'CodedockApiError';
  }
}

export class HttpClient {
  private baseUrl: string;
  private token?: string;
  private fetchImpl: typeof fetch;

  constructor(options: CodedockClientOptions = {}) {
    this.baseUrl = (options.baseUrl || 'http://localhost:8080/api').replace(/\/+$/, '');
    this.token = options.token;
    this.fetchImpl = options.fetch || globalThis.fetch;
  }

  async request<T>(path: string, options: RequestInit = {}): Promise<T> {
    const url = `${this.baseUrl}${path.startsWith('/') ? path : `/${path}`}`;
    const headers = new Headers(options.headers);

    if (this.token && !headers.has('Authorization')) {
      headers.set('Authorization', `Bearer ${this.token}`);
    }
    if (!headers.has('Accept')) {
      headers.set('Accept', 'application/json');
    }
    if (options.body && typeof options.body === 'string' && !headers.has('Content-Type')) {
      headers.set('Content-Type', 'application/json');
    }

    const res = await this.fetchImpl(url, {
      ...options,
      headers,
    });

    if (!res.ok) {
      let body: unknown;
      try {
        body = await res.json();
      } catch {
        body = await res.text();
      }
      throw new CodedockApiError(res.status, res.statusText, body);
    }

    if (res.status === 204) {
      return {} as T;
    }

    const json = (await res.json()) as { data?: T } | T;
    if (json && typeof json === 'object' && 'data' in json && (json as { data: T }).data !== undefined) {
      return (json as { data: T }).data;
    }
    return json as T;
  }

  get<T>(path: string, headers?: HeadersInit): Promise<T> {
    return this.request<T>(path, { method: 'GET', headers });
  }

  post<T>(path: string, body?: unknown, headers?: HeadersInit): Promise<T> {
    return this.request<T>(path, {
      method: 'POST',
      body: body !== undefined ? JSON.stringify(body) : undefined,
      headers,
    });
  }

  put<T>(path: string, body?: unknown, headers?: HeadersInit): Promise<T> {
    return this.request<T>(path, {
      method: 'PUT',
      body: body !== undefined ? JSON.stringify(body) : undefined,
      headers,
    });
  }

  patch<T>(path: string, body?: unknown, headers?: HeadersInit): Promise<T> {
    return this.request<T>(path, {
      method: 'PATCH',
      body: body !== undefined ? JSON.stringify(body) : undefined,
      headers,
    });
  }

  delete<T>(path: string, headers?: HeadersInit): Promise<T> {
    return this.request<T>(path, { method: 'DELETE', headers });
  }
}
