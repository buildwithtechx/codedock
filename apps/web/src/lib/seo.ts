import { env } from '../env';

const SITE_NAME = env.VITE_SITE_NAME || 'Codedock';
const SITE_DESCRIPTION =
  'The open-source Heroku & Vercel alternative. Deploy apps, databases, backups, and services to your own server with zero vendor lock-in.';
const SITE_URL = env.VITE_SITE_URL || 'https://codedock.run';
const TWITTER_HANDLE = '@codedockdotdev';

export interface SeoProps {
  title?: string;
  description?: string;
  image?: string;
  url?: string;
  type?: 'website' | 'article';
  publishedTime?: string;
  tags?: string[];
  noIndex?: boolean;
}

export function createMeta(props: SeoProps = {}) {
  const title = props.title
    ? `${props.title} | ${SITE_NAME}`
    : `${SITE_NAME} — Ship fast, own your infrastructure`;
  const description = props.description || SITE_DESCRIPTION;
  const image = props.image ? new URL(props.image, SITE_URL).href : undefined;
  const url = props.url || SITE_URL;
  const type = props.type || 'website';

  return [
    { title },
    { name: 'description', content: description },
    {
      name: 'keywords',
      content:
        'codedock, self-hosted, PaaS alternative, heroku alternative, vercel alternative, docker deployment, traefik, deployment platform',
    },
    { property: 'og:type', content: type },
    { property: 'og:site_name', content: SITE_NAME },
    { property: 'og:title', content: title },
    { property: 'og:description', content: description },
    { property: 'og:url', content: url },
    ...(image ? [{ property: 'og:image', content: image }] : []),
    { name: 'twitter:card', content: image ? 'summary_large_image' : 'summary' },
    { name: 'twitter:site', content: TWITTER_HANDLE },
    { name: 'twitter:title', content: title },
    { name: 'twitter:description', content: description },
    ...(image ? [{ name: 'twitter:image', content: image }] : []),
    ...(props.noIndex ? [{ name: 'robots', content: 'noindex,nofollow' }] : []),
    ...(props.publishedTime
      ? [{ property: 'article:published_time', content: props.publishedTime }]
      : []),
    ...(props.tags?.map((tag) => ({ property: 'article:tag', content: tag })) ?? []),
  ];
}

export const globalLinks = [
  { rel: 'icon', type: 'image/svg+xml', href: '/favicon.svg' },
  { rel: 'icon', href: '/favicon.ico' },
];

export function pageLinks(canonicalPath: string) {
  return [{ rel: 'canonical', href: new URL(canonicalPath, SITE_URL).href }];
}

export function pageHead(path: string, title: string, description: string) {
  return {
    meta: createMeta({ title, description, url: new URL(path, SITE_URL).href }),
    links: pageLinks(path),
  };
}
