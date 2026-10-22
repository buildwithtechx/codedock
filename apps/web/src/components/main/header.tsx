import { useEffect, useState } from 'react';
import { productLinks } from '../../lib/product-links';
import { GithubIcon } from '../icons/github-icon';
import { SiteLink } from '../site-link';
import { DesktopNavigation } from './desktop-navigation';
import { MobileNavigation } from './mobile-navigation';

export function Header() {
  const [stars, setStars] = useState<number | null>(null);

  useEffect(() => {
    const controller = new AbortController();
    fetch('/api/github-stars', { signal: controller.signal })
      .then((res) => (res.ok ? res.json() : null))
      .then((data) => {
        if (!controller.signal.aborted && Number.isSafeInteger(data?.stars) && data.stars > 0) {
          setStars(data.stars);
        }
      })
      .catch(() => {});
    return () => controller.abort();
  }, []);

  const formattedStars =
    stars != null
      ? new Intl.NumberFormat('en-US', { notation: 'compact', maximumFractionDigits: 1 })
          .format(stars)
          .toLowerCase()
      : null;

  return (
    <header className="sticky top-0 z-40 border-border border-b bg-background/90 backdrop-blur-xl">
      <div className="mx-auto flex h-18 max-w-7xl items-center justify-between gap-4 px-5 lg:px-8">
        <SiteLink
          href="/"
          aria-label="Codedock home"
          className="flex shrink-0 items-center gap-2.5 font-extrabold text-lg tracking-tight"
        >
          <img
            src="/favicon.svg"
            alt="Codedock"
            className="size-8 rounded-xl object-contain shadow-xs"
          />
          Codedock
          <span className="ml-1 hidden rounded-full border border-border px-2 py-0.5 font-medium text-[10px] text-muted-foreground sm:inline">
            OPEN SOURCE
          </span>
        </SiteLink>
        <DesktopNavigation />
        <div className="flex items-center gap-3">
          <SiteLink
            href={productLinks.github}
            aria-label="Codedock on GitHub"
            className="flex items-center gap-1.5 rounded-lg p-2 text-muted-foreground hover:bg-muted hover:text-foreground"
          >
            <GithubIcon className="size-5" />
            {formattedStars && (
              <span className="font-semibold text-xs tracking-tight">★ {formattedStars}</span>
            )}
          </SiteLink>
          <SiteLink
            href={productLinks.cloud}
            className="hidden rounded-xl bg-primary px-4 py-2.5 font-semibold text-sm text-white transition-colors hover:bg-primary/90 sm:block"
          >
            Open Cloud
          </SiteLink>
          <MobileNavigation />
        </div>
      </div>
    </header>
  );
}
