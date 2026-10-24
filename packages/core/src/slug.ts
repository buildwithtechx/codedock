/**
 * Encodes a name or identifier into a URL-friendly, DNS-safe deployment slug.
 * Transforms to lowercase, converts whitespace and non-alphanumeric characters to hyphens,
 * collapses consecutive hyphens, and trims leading/trailing hyphens.
 */
export function encodeDeploySlug(name: string): string {
  if (!name) return '';
  return name
    .toLowerCase()
    .trim()
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/^-+|-+$/g, '');
}

/**
 * Decodes a deployment slug into human-readable words.
 * Replaces hyphens with spaces and capitalizes each word.
 */
export function decodeDeploySlug(slug: string): string {
  if (!slug) return '';
  return slug
    .replace(/-/g, ' ')
    .trim()
    .replace(/\b\w/g, (char) => char.toUpperCase());
}
