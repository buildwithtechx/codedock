import { Link } from '@tanstack/react-router';
import type { ComponentProps } from 'react';
import type { FileRoutesByTo } from '../routeTree.gen';

export type SitePath = keyof FileRoutesByTo;
export type ExternalLink = `https://${string}` | `http://${string}` | `mailto:${string}`;
export type SiteHref = SitePath | ExternalLink;
type SiteLinkProps = Omit<ComponentProps<'a'>, 'href'> & { href: SiteHref };

function isInternalLink(href: SiteHref): href is SitePath {
  return href.startsWith('/');
}

export function SiteLink({ href, ...props }: SiteLinkProps) {
  if (isInternalLink(href)) {
    return <Link {...props} to={href} />;
  }
  return <a {...props} href={href} />;
}
