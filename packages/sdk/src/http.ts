import { CodedockApiError, CodedockBillingGateError, CodedockRateLimitError } from './errors';

export interface CodedockClientOptions {
  baseUrl?: string;
  token?: string;
  tenantId?: string;
  fetch?: typeof fetch;
}

export class HttpClient {
  private baseUrl: string;
  private token?: string;
  private tenantId?: string;
  private fetchImpl: typeof fetch;

  constructor(options: CodedockClientOptions = {}) {
    this.baseUrl = (options.baseUrl || 'http://localhost:8080/api').replace(/\/+$/, '');
    this.token = options.token;
    this.tenantId = options.tenantId;
    this.fetchImpl = options.fetch || globalThis.fetch;
  }

  async request<T>(path: string, options: RequestInit = {}): Promise<T> {
    const url = `${this.baseUrl}${path.startsWith('/') ? path : `/${path}`}`;
    const headers = new Headers(options.headers);

    if (this.token && !headers.has('Authorization')) {
      headers.set('Authorization', `Bearer ${this.token}`);
    }
    if (this.tenantId && !headers.has('X-Tenant-ID')) {
      headers.set('X-Tenant-ID', this.tenantId);
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

      if (res.status === 429) {
        let retryAfter = 60;
        const retryHeader = res.headers.get('Retry-After');
        if (retryHeader) {
          const parsed = Number.parseInt(retryHeader, 10);
          if (!Number.isNaN(parsed) && parsed > 0) {
            retryAfter = parsed;
          }
        } else if (body && typeof body === 'object' && 'retryAfterSeconds' in body) {
          const parsed = Number((body as { retryAfterSeconds?: unknown }).retryAfterSeconds);
          if (!Number.isNaN(parsed) && parsed > 0) {
            retryAfter = parsed;
          }
        }
        const planTier =
          body && typeof body === 'object' && 'planTier' in body
            ? String((body as { planTier?: unknown }).planTier)
            : undefined;
        throw new CodedockRateLimitError(res.status, res.statusText, body, retryAfter, planTier);
      }

      const isBillingGate =
        res.status === 402 ||
        (res.status === 403 &&
          body !== null &&
          typeof body === 'object' &&
          ('code' in body || 'error' in body) &&
          (String((body as Record<string, unknown>).code)
            .toUpperCase()
            .includes('BILLING') ||
            String((body as Record<string, unknown>).code)
              .toUpperCase()
              .includes('PAYMENT') ||
            String((body as Record<string, unknown>).error)
              .toUpperCase()
              .includes('BILLING')));

      if (isBillingGate) {
        const upgradeUrl =
          body && typeof body === 'object' && 'upgradeUrl' in body
            ? String((body as { upgradeUrl?: unknown }).upgradeUrl)
            : undefined;
        const requiredPlan =
          body && typeof body === 'object' && 'requiredPlan' in body
            ? String((body as { requiredPlan?: unknown }).requiredPlan)
            : undefined;
        throw new CodedockBillingGateError(
          res.status,
          res.statusText,
          body,
          upgradeUrl,
          requiredPlan
        );
      }

      throw new CodedockApiError(res.status, res.statusText, body);
    }

    if (res.status === 204) {
      return {} as T;
    }

    const json = (await res.json()) as { data?: T } | T;
    if (
      json &&
      typeof json === 'object' &&
      'data' in json &&
      (json as { data: T }).data !== undefined
    ) {
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
