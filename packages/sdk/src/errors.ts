export class CodedockApiError extends Error {
  constructor(
    public status: number,
    public statusText: string,
    public body: unknown
  ) {
    super(
      `API error ${status} (${statusText}): ${typeof body === 'string' ? body : JSON.stringify(body)}`
    );
    this.name = 'CodedockApiError';
  }
}

export class CodedockRateLimitError extends CodedockApiError {
  constructor(
    status: number,
    statusText: string,
    body: unknown,
    public retryAfterSeconds: number,
    public planTier?: string
  ) {
    super(status, statusText, body);
    this.name = 'CodedockRateLimitError';
  }
}

export class CodedockBillingGateError extends CodedockApiError {
  constructor(
    status: number,
    statusText: string,
    body: unknown,
    public upgradeUrl?: string,
    public requiredPlan?: string
  ) {
    super(status, statusText, body);
    this.name = 'CodedockBillingGateError';
  }
}
