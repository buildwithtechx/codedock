import { Link } from '@tanstack/react-router';
import { BookOpen, Clock, DollarSign, Lightbulb, Menu, X } from 'lucide-react';
import { useState } from 'react';
import { DiscordIcon } from '../icons/discord-icon';
import { GithubIcon } from '../icons/github-icon';

export function Header() {
  const [isOpen, setIsOpen] = useState(false);

  const navLinks = [
    { to: '/philosophy', label: 'Philosophy', icon: Lightbulb },
    { to: '/pricing', label: 'Pricing', icon: DollarSign },
    { to: '/changelog', label: 'Changelog', icon: Clock },
    { href: 'https://docs.codedock.run', label: 'Docs', icon: BookOpen },
    { href: 'https://github.com/techxteam/codedock', label: 'GitHub', icon: GithubIcon },
    { href: 'https://discord.gg/codedock', label: 'Community', icon: DiscordIcon },
  ];

  return (
    <header className="sticky top-0 z-50 border-border border-b bg-background/80 backdrop-blur-md">
      <nav className="mx-auto flex h-16 max-w-7xl items-center justify-between px-4">
        <Link
          to="/"
          className="relative z-50 flex shrink-0 items-center gap-2 font-bold text-lg tracking-tight"
        >
          <span className="flex size-7 items-center justify-center rounded-lg bg-primary font-black text-white text-xs">
            V
          </span>
          Codedock
        </Link>

        <div className="hidden items-center gap-0.5 text-sm lg:flex">
          {navLinks.map((link) => {
            const Icon = link.icon;
            if (link.to) {
              return (
                <Link
                  key={link.label}
                  to={link.to}
                  className="flex items-center gap-1.5 rounded-md px-3 py-2 font-semibold text-foreground transition-colors hover:bg-surface hover:text-foreground"
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
                className="flex items-center gap-1.5 rounded-md px-3 py-2 font-semibold text-foreground transition-colors hover:bg-surface hover:text-foreground"
              >
                <Icon className="size-4 text-primary" />
                {link.label}
              </a>
            );
          })}
        </div>

        <a
          href="https://app.codedock.run"
          className="hidden shrink-0 rounded-lg bg-primary px-4 py-2 font-semibold text-sm text-white transition-colors hover:bg-primary-hover lg:inline-flex"
        >
          Get Cloud
        </a>

        <button
          type="button"
          onClick={() => setIsOpen(!isOpen)}
          className="relative z-50 -mr-2 p-2 text-foreground lg:hidden"
          aria-label="Toggle Menu"
        >
          {isOpen ? <X className="size-6" /> : <Menu className="size-6" />}
        </button>
      </nav>

      {isOpen && (
        <div className="absolute top-full right-0 left-0 z-40 flex flex-col border-border border-b bg-[#111111] px-4 py-6 shadow-2xl lg:hidden">
          <div className="flex flex-col gap-1">
            {navLinks.map((link) => {
              const Icon = link.icon;
              if (link.to) {
                return (
                  <Link
                    key={link.label}
                    to={link.to}
                    onClick={() => setIsOpen(false)}
                    className="flex items-center gap-4 rounded-md px-3 py-3 font-medium text-foreground text-sm transition-colors hover:bg-white/5"
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
                  className="flex items-center gap-4 rounded-md px-3 py-3 font-medium text-foreground text-sm transition-colors hover:bg-white/5"
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
              className="flex w-full items-center justify-center rounded-md border border-[#3e257f] bg-[#231247] py-3.5 font-bold text-sm text-white transition-colors hover:bg-[#2d185e]"
            >
              To Cloud
            </a>
          </div>
        </div>
      )}
    </header>
  );
}
