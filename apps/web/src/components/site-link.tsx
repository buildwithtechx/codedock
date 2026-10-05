import { Link } from '@tanstack/react-router';
import type { ComponentProps } from 'react';
import type { FileRoutesByTo } from '../routeTree.gen';

type SiteLinkProps = ComponentProps<'a'> & { href: string };
export function SiteLink({ href, ...props }: SiteLinkProps) {
  if (href.startsWith('/') && !href.startsWith('//')) {
    return <Link {...props} to={href as keyof FileRoutesByTo} />;
  }
  return <a {...props} href={href} />;
}
