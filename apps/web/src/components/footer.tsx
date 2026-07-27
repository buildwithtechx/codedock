import { Link } from "@tanstack/react-router";
import { DiscordIcon } from "./discord-icon";
import { GithubIcon } from "./github-icon";
import { XIcon } from "./x-icon";

export function Footer() {
  return (
    <footer className="mt-24 border-t border-border bg-background/80 backdrop-blur-md">
      <div className="max-w-7xl mx-auto px-4 py-12 flex flex-col md:flex-row items-center justify-between gap-8">
        <div className="flex flex-col sm:flex-row items-center gap-3 sm:gap-2 text-center sm:text-left">
          <span className="size-6 rounded-md bg-primary flex items-center justify-center text-white text-xs font-black">
            V
          </span>
          <span className="text-sm text-muted-foreground">
            © 2026 Codedock. Open-source &amp; free forever.
          </span>
        </div>

        <div className="flex flex-col sm:flex-row flex-wrap justify-center items-center gap-4 sm:gap-2 text-sm text-muted-foreground">
          <div className="flex flex-wrap justify-center items-center gap-1 sm:gap-2">
            <Link
              to="/philosophy"
              className="px-3 py-2 rounded-md hover:text-foreground hover:bg-surface transition-colors"
            >
              Philosophy
            </Link>
            <Link
              to="/pricing"
              className="px-3 py-2 rounded-md hover:text-foreground hover:bg-surface transition-colors"
            >
              Pricing
            </Link>
            <Link
              to="/changelog"
              className="px-3 py-2 rounded-md hover:text-foreground hover:bg-surface transition-colors"
            >
              Changelog
            </Link>
            <a
              href="https://docs.codedock.run"
              target="_blank"
              rel="noreferrer"
              className="px-3 py-2 rounded-md hover:text-foreground hover:bg-surface transition-colors"
            >
              Docs
            </a>
          </div>

          <div className="flex items-center gap-2 sm:gap-1 sm:ml-2 sm:pl-4 sm:border-l border-border pt-2 sm:pt-0 border-t sm:border-t-0 w-full sm:w-auto justify-center">
            <a
              href="https://github.com/techxteam/codedock"
              target="_blank"
              rel="noreferrer"
              aria-label="GitHub"
              className="p-2 rounded-md text-muted-foreground hover:text-foreground hover:bg-surface transition-colors"
            >
              <GithubIcon className="size-4" />
            </a>
            <a
              href="https://x.com/codedockdotdev"
              target="_blank"
              rel="noreferrer"
              aria-label="X / Twitter"
              className="p-2 rounded-md text-muted-foreground hover:text-foreground hover:bg-surface transition-colors"
            >
              <XIcon className="size-4" />
            </a>
            <a
              href="https://discord.gg/codedock"
              target="_blank"
              rel="noreferrer"
              aria-label="Discord"
              className="p-2 rounded-md text-muted-foreground hover:text-foreground hover:bg-surface transition-colors"
            >
              <DiscordIcon className="size-4" />
            </a>
          </div>
        </div>
      </div>
    </footer>
  );
}
