const DOMAIN_REGEX = /^(?:[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}$/;

/**
 * Validates whether a given string is a syntactically valid domain name.
 */
export function isValidDomain(domain: string): boolean {
  if (!domain || typeof domain !== 'string') return false;
  const clean = domain.trim().toLowerCase();
  if (clean.length > 253) return false;
  return DOMAIN_REGEX.test(clean);
}

/**
 * Validates whether a given string is a syntactically valid HTTP or HTTPS URL.
 */
export function isValidUrl(input: string): boolean {
  if (!input || typeof input !== 'string') return false;
  try {
    const parsed = new URL(input);
    return parsed.protocol === 'http:' || parsed.protocol === 'https:';
  } catch {
    return false;
  }
}

/** Resolves the parent wildcard pattern (*.example.com) for a subdomain, or null. */
export function parentWildcardDomain(domain: string): string | null {
  if (!domain || domain.startsWith('*.')) return null;
  const clean = domain.trim().toLowerCase();
  const firstDot = clean.indexOf('.');
  if (firstDot <= 0) return null;
  const parent = clean.slice(firstDot + 1);
  if (!parent.includes('.')) return null;
  return `*.${parent}`;
}

/** Validates that disabling swap requires a positive memory limit. */
export function validateSwapDisabled(
  resources?: { memoryLimit?: number; swapDisabled?: boolean } | null
): boolean {
  if (resources?.swapDisabled === true && (!resources.memoryLimit || resources.memoryLimit <= 0)) {
    return false;
  }
  return true;
}
