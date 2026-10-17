import { Check, ChevronDown, Mail, MessageSquare, Send, Users, Webhook } from 'lucide-react';
import type React from 'react';
import { useEffect, useRef, useState } from 'react';

export type ChannelKind = 'email' | 'webhook' | 'slack' | 'discord' | 'msteams' | 'telegram';

interface ChannelChoice {
  id: ChannelKind;
  label: string;
  icon: React.ComponentType<{ className?: string }>;
}

const CHANNELS: ChannelChoice[] = [
  { id: 'email', label: 'Email', icon: Mail },
  { id: 'webhook', label: 'Webhook', icon: Webhook },
  { id: 'slack', label: 'Slack', icon: MessageSquare },
  { id: 'discord', label: 'Discord', icon: MessageSquare },
  { id: 'msteams', label: 'Microsoft Teams', icon: Users },
  { id: 'telegram', label: 'Telegram', icon: Send },
];

export function ChannelMultiSelect({
  value,
  disabled,
  onChange,
}: {
  value: string[];
  disabled?: boolean;
  onChange: (channels: string[]) => void;
}) {
  const [open, setOpen] = useState(false);
  const containerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!open) return;
    const handleClickOutside = (e: MouseEvent) => {
      if (containerRef.current && !containerRef.current.contains(e.target as Node)) {
        setOpen(false);
      }
    };
    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, [open]);

  const toggle = (id: string) => {
    const has = value.includes(id);
    if (has && value.length === 1) return;
    const next = has ? value.filter((c) => c !== id) : [...value, id];
    onChange(next);
  };

  return (
    <div ref={containerRef} className="relative inline-flex justify-end text-start">
      <button
        type="button"
        disabled={disabled}
        onClick={() => setOpen((prev) => !prev)}
        className="inline-flex items-center gap-2 rounded-lg border border-border/50 bg-secondary/80 px-2.5 py-1.5 font-normal text-secondary-foreground text-xs transition-colors hover:bg-secondary disabled:cursor-not-allowed disabled:opacity-50"
        aria-expanded={open}
      >
        <span className="flex items-center gap-1.5">
          <span className="flex -space-x-1">
            {value.slice(0, 2).map((ch) => {
              const item = CHANNELS.find((c) => c.id === ch) || CHANNELS[0];
              const Icon = item.icon;
              return (
                <span
                  key={ch}
                  className="grid size-5 place-items-center rounded-full bg-muted ring-1 ring-background"
                >
                  <Icon className="size-3 text-muted-foreground" />
                </span>
              );
            })}
            {value.length > 2 && (
              <span className="grid size-5 place-items-center rounded-full bg-muted font-semibold text-[10px] text-muted-foreground ring-1 ring-background">
                +{value.length - 2}
              </span>
            )}
          </span>
          <span className="whitespace-nowrap text-muted-foreground">{value.length} channels</span>
        </span>
        <ChevronDown className="size-3.5 shrink-0 text-muted-foreground" />
      </button>

      {open && (
        <div
          role="menu"
          className="fade-in slide-in-from-top-1 absolute right-0 z-50 mt-1.5 w-52 animate-in overflow-hidden rounded-xl border border-border/60 bg-popover/95 p-1 shadow-2xl backdrop-blur-xl duration-150"
        >
          {CHANNELS.map((ch) => {
            const active = value.includes(ch.id);
            const Icon = ch.icon;
            return (
              <button
                key={ch.id}
                type="button"
                role="menuitemcheckbox"
                aria-checked={active}
                onClick={() => toggle(ch.id)}
                className="flex w-full items-center gap-2.5 rounded-lg px-2.5 py-1.5 text-start text-foreground text-xs transition-colors hover:bg-muted/50"
              >
                <span
                  className={`grid size-4 shrink-0 place-items-center rounded border transition-colors ${
                    active
                      ? 'border-primary bg-primary text-primary-foreground'
                      : 'border-border/70 bg-transparent'
                  }`}
                >
                  {active && <Check className="size-3" />}
                </span>
                <Icon className="size-4 shrink-0 text-muted-foreground" />
                <span className="flex-1 font-normal text-foreground text-sm">{ch.label}</span>
              </button>
            );
          })}
        </div>
      )}
    </div>
  );
}
