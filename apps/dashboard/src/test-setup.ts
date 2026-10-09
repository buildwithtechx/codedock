import { vi } from 'vitest';

const matchMediaStub = vi.fn().mockImplementation((query: string) => ({
  matches: false,
  media: query,
  onchange: null,
  addListener: vi.fn(),
  removeListener: vi.fn(),
  addEventListener: vi.fn(),
  removeEventListener: vi.fn(),
  dispatchEvent: vi.fn(),
}));

vi.stubGlobal('matchMedia', matchMediaStub);

vi.stubGlobal(
  'ResizeObserver',
  class ResizeObserver {
    observe() {}
    unobserve() {}
    disconnect() {}
  }
);

vi.stubGlobal(
  'IntersectionObserver',
  class IntersectionObserver {
    observe() {}
    unobserve() {}
    disconnect() {}
  }
);

vi.stubGlobal(
  'WebSocket',
  class WebSocket {
    url: string;
    onopen: ((ev: Event) => void) | null = null;
    onmessage: ((ev: MessageEvent) => void) | null = null;
    onerror: ((ev: Event) => void) | null = null;
    onclose: ((ev: CloseEvent) => void) | null = null;

    constructor(url: string) {
      this.url = url;
    }

    send() {}
    close() {}
    addEventListener() {}
    removeEventListener() {}
  }
);
