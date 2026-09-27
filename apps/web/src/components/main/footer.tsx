import { Link } from '@tanstack/react-router';
import { DiscordIcon } from '../icons/discord-icon';
import { GithubIcon } from '../icons/github-icon';
import { XIcon } from '../icons/x-icon';

export function Footer() {
  return (
    <footer className="mt-24 border-border border-t bg-background/80 backdrop-blur-md">
      <div className="mx-auto flex max-w-7xl flex-col items-center justify-between gap-8 px-4 py-12 md:flex-row">
        <div className="flex flex-col items-center gap-3 text-center sm:flex-row sm:gap-2 sm:text-left">
          <span className="flex size-6 items-center justify-center rounded-md bg-primary font-black text-white text-xs">
            V
          </span>
          <span className="text-muted-foreground text-sm">
            © 2026 Codedock. Open-source &amp; free forever.
          </span>
        </div>

        <div className="flex flex-col flex-wrap items-center justify-center gap-4 text-muted-foreground text-sm sm:flex-row sm:gap-2">
          <div className="flex flex-wrap items-center justify-center gap-1 sm:gap-2">
            <Link
              to="/philosophy"
              className="rounded-md px-3 py-2 transition-colors hover:bg-surface hover:text-foreground"
            >
              Philosophy
            </Link>
            <Link
              to="/pricing"
              className="rounded-md px-3 py-2 transition-colors hover:bg-surface hover:text-foreground"
            >
              Pricing
            </Link>
            <Link
              to="/changelog"
              className="rounded-md px-3 py-2 transition-colors hover:bg-surface hover:text-foreground"
            >
              Changelog
            </Link>
            <a
              href="https://docs.codedock.run"
              target="_blank"
              rel="noreferrer"
              className="rounded-md px-3 py-2 transition-colors hover:bg-surface hover:text-foreground"
            >
              Docs
            </a>
          </div>

          <div className="flex w-full items-center justify-center gap-2 border-border border-t pt-2 sm:ml-2 sm:w-auto sm:gap-1 sm:border-t-0 sm:border-l sm:pt-0 sm:pl-4">
            <a
              href="https://github.com/techxteam/codedock"
              target="_blank"
              rel="noreferrer"
              aria-label="GitHub"
              className="rounded-md p-2 text-muted-foreground transition-colors hover:bg-surface hover:text-foreground"
            >
              <GithubIcon className="size-4" />
            </a>
            <a
              href="https://x.com/codedockdotdev"
              target="_blank"
              rel="noreferrer"
              aria-label="X / Twitter"
              className="rounded-md p-2 text-muted-foreground transition-colors hover:bg-surface hover:text-foreground"
            >
              <XIcon className="size-4" />
            </a>
            <a
              href="https://discord.gg/codedock"
              target="_blank"
              rel="noreferrer"
              aria-label="Discord"
              className="rounded-md p-2 text-muted-foreground transition-colors hover:bg-surface hover:text-foreground"
            >
              <DiscordIcon className="size-4" />
            </a>
          </div>
        </div>
      </div>
    </footer>
  );
}
