import { describe, expect, it } from 'vitest';
import { badgeVariants, buttonVariants, cn, statusDotVariants } from '../src/index';

describe('@codedock/ui', () => {
  describe('cn utility', () => {
    it('merges class names and handles conditionals correctly', () => {
      expect(cn('base', true && 'is-active', false && 'is-disabled')).toBe('base is-active');
    });

    it('resolves conflicting Tailwind utility classes', () => {
      expect(cn('px-4 py-2', 'px-6')).toBe('py-2 px-6');
      expect(cn('bg-red-500', 'bg-blue-500')).toBe('bg-blue-500');
    });
  });

  describe('buttonVariants', () => {
    it('produces default button styles', () => {
      const classes = buttonVariants();
      expect(classes).toContain('bg-primary');
      expect(classes).toContain('h-10');
    });

    it('produces variant-specific styles', () => {
      expect(buttonVariants({ variant: 'destructive' })).toContain('bg-destructive');
      expect(buttonVariants({ variant: 'outline' })).toContain('border-border');
      expect(buttonVariants({ size: 'sm' })).toContain('h-8');
    });
  });

  describe('badgeVariants', () => {
    it('produces default and success badge styles', () => {
      expect(badgeVariants({ variant: 'default' })).toContain('bg-primary');
      expect(badgeVariants({ variant: 'success' })).toContain('emerald');
      expect(badgeVariants({ variant: 'destructive' })).toContain('bg-destructive');
    });
  });

  describe('statusDotVariants', () => {
    it('produces dot classes for each status', () => {
      expect(statusDotVariants({ status: 'ready' })).toContain('bg-emerald-500');
      expect(statusDotVariants({ status: 'running' })).toContain('bg-blue-500');
      expect(statusDotVariants({ status: 'error' })).toContain('bg-red-500');
      expect(statusDotVariants({ size: 'lg' })).toContain('size-2.5');
    });
  });
});
