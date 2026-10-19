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
