import type { ClassValue } from 'clsx';
import { clsx } from 'clsx';
import { twMerge } from 'tailwind-merge';

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function formatBytes(bytes: number): string {
  if (!Number.isFinite(bytes) || bytes <= 0) return '0 B';
  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
  const i = Math.min(Math.floor(Math.log(bytes) / Math.log(k)), sizes.length - 1);
  return `${(bytes / k ** i).toFixed(1)} ${sizes[i]}`;
}

export function encodeDeploySlug(input: string): string {
  try {
    const base64 = btoa(
      encodeURIComponent(input).replace(/%([0-9A-F]{2})/g, (_, p1) =>
        String.fromCharCode(parseInt(p1, 16))
      )
    );
    return base64.replace(/\+/g, '-').replace(/\//g, '_').replace(/=+$/, '');
  } catch {
    return encodeURIComponent(input);
  }
}

export function decodeDeploySlug(slug: string): string {
  const raw = decodeURIComponent(slug);
  try {
    const padded = raw.replace(/-/g, '+').replace(/_/g, '/');
    const padLength = (4 - (padded.length % 4)) % 4;
    const base64 = padded + '='.repeat(padLength);
    const decoded = decodeURIComponent(
      Array.prototype.map
        .call(atob(base64), (c: string) => `%${(`00${c.charCodeAt(0).toString(16)}`).slice(-2)}`)
        .join('')
    );
    if (decoded && /^[\w.\-:/@]+$/.test(decoded)) {
      return decoded;
    }
  } catch {}
  return raw;
}
