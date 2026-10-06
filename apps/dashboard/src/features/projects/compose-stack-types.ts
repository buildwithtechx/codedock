export interface ComposeStack {
  id: string;
  projectId: string;
  environmentId: string;
  name: string;
  revision: number;
  status: string;
  error: string;
  results: string;
}

export interface ComposeServiceResult {
  name: string;
  containerId: string;
  state: string;
  health: string;
  exitCode: number;
}

export interface ComposeStackReview {
  config: string;
  digest: string;
  services: string[];
  effects: string[];
}

export const activeStackStates = new Set(['PREPARING', 'BUILDING', 'STARTING', 'READINESS']);

export function stackResults(stack: ComposeStack): ComposeServiceResult[] {
  try {
    const parsed: unknown = JSON.parse(stack.results);
    return Array.isArray(parsed) ? (parsed as ComposeServiceResult[]) : [];
  } catch {
    return [];
  }
}
