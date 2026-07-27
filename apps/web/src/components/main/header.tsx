import { Link } from "@tanstack/react-router";
import { BookOpen, Clock, DollarSign, Lightbulb, Menu, X } from "lucide-react";
import { useState } from "react";
import { DiscordIcon } from "./discord-icon";
import { GithubIcon } from "./github-icon";

export function Header() {
  const [isOpen, setIsOpen] = useState(false);

  const navLinks = [
    { to: "/philosophy", label: "Philosophy", icon: Lightbulb },
    { to: "/pricing", label: "Pricing", icon: DollarSign },
    { to: "/changelog", label: "Changelog", icon: Clock },
    { href: "https://docs.codedock.run", label: "Docs", icon: BookOpen },
    { href: "https://github.com/techxteam/codedock", label: "GitHub", icon: GithubIcon },
    { href: "https://discord.gg/codedock", label: "Community", icon: DiscordIcon },
  ];

  return (
    <header className="sticky top-0 z-50 border-b border-border bg-background/80 backdrop-blur-md">
      <nav className="max-w-7xl mx-auto px-4 flex items-center justify-between h-16">
        <Link
          to="/"
          className="flex items-center gap-2 font-bold text-lg tracking-tight shrink-0 relative z-50"
        >
          <span className="size-7 rounded-lg bg-primary flex items-center justify-center text-white text-xs font-black">
            V
          </span>
          Codedock
        </Link>

        <div className="hidden lg:flex items-center gap-0.5 text-sm">
          {navLinks.map((link) => {
            const Icon = link.icon;
            if (link.to) {
              return (
                <Link
                  key={link.label}
                  to={link.to}
                  className="flex items-center gap-1.5 px-3 py-2 rounded-md font-semibold text-foreground hover:text-foreground hover:bg-surface transition-colors"
                >
                  <Icon className="size-4 text-primary" />
                  {link.label}
                </Link>
              );
            }
            return (
              <a
                key={link.label}
                href={link.href}
                target="_blank"
                rel="noreferrer"
                className="flex items-center gap-1.5 px-3 py-2 rounded-md font-semibold text-foreground hover:text-foreground hover:bg-surface transition-colors"
              >
                <Icon className="size-4 text-primary" />
                {link.label}
              </a>
            );
          })}
        </div>

        <a
          href="https://app.codedock.run"
          className="hidden lg:inline-flex shrink-0 text-sm font-semibold px-4 py-2 rounded-lg bg-primary hover:bg-primary-hover text-white transition-colors"
        >
          Get Cloud
        </a>

        <button
          type="button"
          onClick={() => setIsOpen(!isOpen)}
          className="lg:hidden p-2 -mr-2 text-foreground relative z-50"
          aria-label="Toggle Menu"
        >
          {isOpen ? <X className="size-6" /> : <Menu className="size-6" />}
        </button>
      </nav>

      {isOpen && (
        <div className="absolute top-full left-0 right-0 bg-[#111111] border-b border-border z-40 flex flex-col px-4 py-6 lg:hidden shadow-2xl">
          <div className="flex flex-col gap-1">
            {navLinks.map((link) => {
              const Icon = link.icon;
              if (link.to) {
                return (
                  <Link
                    key={link.label}
                    to={link.to}
                    onClick={() => setIsOpen(false)}
                    className="flex items-center gap-4 px-3 py-3 rounded-md font-medium text-sm text-foreground hover:bg-white/5 transition-colors"
                  >
                    <Icon className="size-5 text-primary" />
                    {link.label}
                  </Link>
                );
              }
              return (
                <a
                  key={link.label}
                  href={link.href}
                  target="_blank"
                  rel="noreferrer"
                  onClick={() => setIsOpen(false)}
                  className="flex items-center gap-4 px-3 py-3 rounded-md font-medium text-sm text-foreground hover:bg-white/5 transition-colors"
                >
                  <Icon className="size-5 text-primary" />
                  {link.label}
                </a>
              );
            })}
          </div>

          <div className="mt-8 mb-2">
            <a
              href="https://app.codedock.run"
              className="flex items-center justify-center w-full py-3.5 rounded-md text-sm font-bold text-white transition-colors bg-[#231247] hover:bg-[#2d185e] border border-[#3e257f]"
            >
              To Cloud
            </a>
          </div>
        </div>
      )}
    </header>
  );
}
