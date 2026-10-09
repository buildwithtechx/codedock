import { Check, ChevronDown, Mail, MessageSquare, Send, Webhook } from 'lucide-react';
import { useEffect, useRef, useState } from 'react';

const CHANNELS = [
  { id: 'email', label: 'Email', icon: Mail },
  { id: 'webhook', label: 'Webhook', icon: Webhook },
  { id: 'slack', label: 'Slack', icon: MessageSquare },
  { id: 'discord', label: 'Discord', icon: MessageSquare },
  { id: 'telegram', label: 'Telegram', icon: Send },
];

export function ChannelMultiSelect({
  value,
  onChange,
}: {
  value: string[];
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
    const next = value.includes(id)
      ? value.length > 1
        ? value.filter((c) => c !== id)
        : value
      : [...value, id];
    onChange(next);
  };

  return (
    <div ref={containerRef} className="relative inline-block text-start">
      <button
        type="button"
        onClick={() => setOpen((prev) => !prev)}
        className="inline-flex items-center gap-2 rounded-lg border border-border/60 bg-background px-2.5 py-1.5 text-foreground text-xs transition-colors hover:bg-muted/40"
      >
        <span className="flex items-center gap-1.5">
          <span className="flex -space-x-1">
            {value.slice(0, 2).map((ch) => {
              const item = CHANNELS.find((c) => c.id === ch) || CHANNELS[0];
              const Icon = item.icon;
              return (
                <span
                  key={ch}
                  className="grid size-4 place-items-center rounded-full bg-muted ring-1 ring-background"
                >
                  <Icon className="size-2.5 text-muted-foreground" />
                </span>
              );
            })}
            {value.length > 2 && (
              <span className="grid size-4 place-items-center rounded-full bg-muted font-semibold text-[8px] text-muted-foreground ring-1 ring-background">
                +{value.length - 2}
              </span>
            )}
          </span>
          <span className="text-muted-foreground">
            {value.length} channel{value.length === 1 ? '' : 's'}
          </span>
        </span>
        <ChevronDown className="size-3 text-muted-foreground" />
      </button>

      {open && (
        <div className="absolute right-0 z-50 mt-1 w-44 rounded-xl border border-border/70 bg-popover/95 p-1 shadow-lg backdrop-blur-md">
          {CHANNELS.map((ch) => {
            const active = value.includes(ch.id);
            const Icon = ch.icon;
            return (
              <button
                key={ch.id}
                type="button"
                onClick={() => toggle(ch.id)}
                className="flex w-full items-center gap-2 rounded-lg px-2 py-1.5 text-foreground text-xs transition-colors hover:bg-muted/50"
              >
                <span
                  className={`grid size-3.5 place-items-center rounded border transition-colors ${
                    active ? 'border-primary bg-primary text-primary-foreground' : 'border-border'
                  }`}
                >
                  {active && <Check className="size-2.5" />}
                </span>
                <Icon className="size-3.5 text-muted-foreground" />
                <span className="flex-1 text-left">{ch.label}</span>
              </button>
            );
          })}
        </div>
      )}
    </div>
  );
}
