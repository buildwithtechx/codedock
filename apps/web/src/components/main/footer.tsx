import { footerGroups } from '../../lib/navigation';
import { productLinks } from '../../lib/product-links';
import { DiscordIcon } from '../icons/discord-icon';
import { GithubIcon } from '../icons/github-icon';
import { SiteLink } from '../site-link';

export function Footer() {
  return (
    <footer className="mt-12 border-border border-t bg-card/40">
      <div className="mx-auto max-w-7xl px-6 py-14 lg:px-8">
        <div className="mb-12 flex flex-col justify-between gap-6 sm:flex-row sm:items-center">
          <div>
            <SiteLink href="/" className="font-extrabold text-xl">
              Codedock<span className="text-primary">.</span>
            </SiteLink>
            <p className="mt-2 text-muted-foreground text-sm">
              Ship your applications. Own your infrastructure.
            </p>
          </div>
          <SiteLink
            href={productLinks.installation}
            className="w-fit rounded-xl border border-border px-5 py-3 font-semibold text-sm hover:border-primary/50"
          >
            Start self-hosting →
          </SiteLink>
        </div>
        <nav
          aria-label="Footer navigation"
          className="grid grid-cols-2 gap-8 sm:grid-cols-3 lg:grid-cols-5"
        >
          {footerGroups.map((group) => (
            <div key={group.label}>
              <h2 className="mb-4 font-semibold text-sm">{group.label}</h2>
              <ul className="space-y-3">
                {group.items.map((item) => (
                  <li key={item.href}>
                    <SiteLink
                      href={item.href}
                      className="text-muted-foreground text-sm transition-colors hover:text-foreground"
                    >
                      {item.label}
                    </SiteLink>
                  </li>
                ))}
              </ul>
            </div>
          ))}
        </nav>
        <div className="mt-12 flex flex-col justify-between gap-4 border-border border-t pt-6 text-muted-foreground text-xs sm:flex-row sm:items-center">
          <p>&copy; {new Date().getFullYear()} Codedock. Source licensed under Apache-2.0.</p>
          <div className="flex items-center gap-5">
            <SiteLink
              href={productLinks.github}
              aria-label="View the source on GitHub"
              className="hover:text-foreground"
            >
              <GithubIcon className="size-4" />
            </SiteLink>
            <SiteLink
              href={productLinks.discord}
              aria-label="Join the Discord community"
              className="hover:text-foreground"
            >
              <DiscordIcon className="size-4" />
            </SiteLink>
            <span>Your server. Your data.</span>
          </div>
        </div>
      </div>
    </footer>
  );
}
