import { Link } from '@tanstack/react-router';
import { ArrowUp, Box, Database, Edit3, SlidersHorizontal, Terminal, X, Zap } from 'lucide-react';
import { type KeyboardEvent, useMemo, useState } from 'react';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@codedock/ui';
import { useAiResources } from '#/hooks/use-ai';
import type { AiProviderStatusSummary, AiResourceTag } from '#/interfaces/ai';
import { useAiSidebarStore } from '#/stores/ai-sidebar-store';

const SHORTCUTS = [
  {
    key: '?status',
    label: 'Instance health & metrics',
    prompt: 'Check the overall health, projects, and resource status of this Codedock instance.',
  },
  {
    key: '?diagnose',
    label: 'Diagnose failing services',
    prompt: 'Inspect all services and recent deployments for errors or unhealthy containers.',
  },
  {
    key: '?redeploy',
    label: 'Redeploy a service',
    prompt: 'Help me trigger a zero-downtime redeploy on one of my application services.',
  },
  {
    key: '?domains',
    label: 'Audit domains & SSL',
    prompt:
      'List configured domains and check if any TLS certificates or DNS records need attention.',
  },
  {
    key: '?backups',
    label: 'Check backup runs',
    prompt: 'Show my S3/SFTP backup destinations and recent backup policies.',
  },
];

export function AskAiComposer({
  onSend,
  disabled,
  providerStatus,
}: {
  onSend: (text: string, mentions: AiResourceTag[]) => void;
  disabled: boolean;
  providerStatus: AiProviderStatusSummary;
}) {
  const [text, setText] = useState('');
  const [mentions, setMentions] = useState<AiResourceTag[]>([]);
  const mode = useAiSidebarStore((s) => s.mode);
  const setMode = useAiSidebarStore((s) => s.setMode);

  const { data: resources = [] } = useAiResources();

  const activeTrigger = useMemo(() => {
    const trimmedEnd = text;
    const atIdx = trimmedEnd.lastIndexOf('@');
    if (atIdx !== -1 && (atIdx === 0 || trimmedEnd[atIdx - 1] === ' ')) {
      const query = trimmedEnd.slice(atIdx + 1);
      if (!query.includes(' '))
        return { kind: 'at' as const, query: query.toLowerCase(), index: atIdx };
    }
    if (trimmedEnd.startsWith('?')) {
      return { kind: 'shortcut' as const, query: trimmedEnd.slice(1).toLowerCase(), index: 0 };
    }
    return null;
  }, [text]);

  const matchingResources = useMemo(() => {
    if (activeTrigger?.kind !== 'at') return [];
    const q = activeTrigger.query;
    return resources
      .filter((r) => !q || `${r.type} ${r.name} ${r.subtitle ?? ''}`.toLowerCase().includes(q))
      .slice(0, 6);
  }, [activeTrigger, resources]);

  const matchingShortcuts = useMemo(() => {
    if (activeTrigger?.kind !== 'shortcut') return [];
    const q = activeTrigger.query;
    return SHORTCUTS.filter(
      (s) => !q || s.key.toLowerCase().includes(q) || s.label.toLowerCase().includes(q)
    );
  }, [activeTrigger]);

  const selectResource = (res: AiResourceTag) => {
    if (!activeTrigger) return;
    const before = text.slice(0, activeTrigger.index);
    setText(before);
    if (!mentions.some((m) => m.id === res.id && m.type === res.type)) {
      setMentions((prev) => [...prev, res]);
    }
  };

  const selectShortcut = (prompt: string) => {
    setText(prompt);
  };

  const handleSubmit = () => {
    const trimmed = text.trim();
    if (!trimmed || disabled) return;
    onSend(trimmed, mentions);
    setText('');
    setMentions([]);
  };

  const handleKeyDown = (e: KeyboardEvent<HTMLTextAreaElement>) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSubmit();
    }
  };

  return (
    <div className="relative border-border/60 border-t bg-sidebar p-3">
      {/* @Resource Tagging Autocomplete Popover */}
      {activeTrigger?.kind === 'at' && matchingResources.length > 0 && (
        <div className="absolute right-3 bottom-full left-3 z-50 mb-2 overflow-hidden rounded-xl border border-border bg-popover p-1.5 shadow-xl">
          <p className="px-2 py-1 font-medium text-[10px] text-muted-foreground uppercase tracking-wider">
            Tag Instance Resource
          </p>
          {matchingResources.map((r) => (
            <button
              key={`${r.type}-${r.id}`}
              type="button"
              onClick={() => selectResource(r)}
              className="flex w-full items-center justify-between gap-2 rounded-lg px-2.5 py-1.5 text-left text-xs transition-colors hover:bg-muted"
            >
              <div className="flex min-w-0 items-center gap-2">
                {r.type === 'project' ? (
                  <Box className="size-3.5 shrink-0 text-primary" />
                ) : r.type === 'database' ? (
                  <Database className="size-3.5 shrink-0 text-amber-500" />
                ) : (
                  <Terminal className="size-3.5 shrink-0 text-emerald-500" />
                )}
                <span className="truncate font-medium text-foreground">{r.name}</span>
              </div>
              <span className="shrink-0 rounded bg-muted px-1.5 py-0.5 font-mono text-[10px] text-muted-foreground">
                {r.type}
              </span>
            </button>
          ))}
        </div>
      )}

      {/* ? Shortcuts Autocomplete Popover */}
      {activeTrigger?.kind === 'shortcut' && matchingShortcuts.length > 0 && (
        <div className="absolute right-3 bottom-full left-3 z-50 mb-2 overflow-hidden rounded-xl border border-border bg-popover p-1.5 shadow-xl">
          <p className="px-2 py-1 font-medium text-[10px] text-muted-foreground uppercase tracking-wider">
            Quick Shortcuts
          </p>
          {matchingShortcuts.map((s) => (
            <button
              key={s.key}
              type="button"
              onClick={() => selectShortcut(s.prompt)}
              className="flex w-full items-center justify-between gap-2 rounded-lg px-2.5 py-1.5 text-left text-xs transition-colors hover:bg-muted"
            >
              <code className="font-mono font-semibold text-primary text-xs">{s.key}</code>
              <span className="truncate text-muted-foreground">{s.label}</span>
            </button>
          ))}
        </div>
      )}

      <div className="rounded-2xl border border-border/70 bg-card/90 p-3 shadow-xs transition-colors focus-within:border-border">
        {mentions.length > 0 && (
          <div className="mb-2 flex flex-wrap gap-1.5">
            {mentions.map((m) => (
              <span
                key={`${m.type}-${m.id}`}
                className="inline-flex items-center gap-1 rounded-md border border-primary/30 bg-primary/10 px-2 py-0.5 font-mono text-[11px] text-primary"
              >
                @{m.type}:{m.name}
                <button
                  type="button"
                  onClick={() =>
                    setMentions((prev) => prev.filter((x) => !(x.id === m.id && x.type === m.type)))
                  }
                  aria-label={`Remove ${m.name}`}
                  className="hover:text-foreground"
                >
                  <X className="size-3" />
                </button>
              </span>
            ))}
          </div>
        )}

        <textarea
          rows={2}
          value={text}
          onChange={(e) => setText(e.target.value)}
          onKeyDown={handleKeyDown}
          placeholder="Type @ to tag a resource or ? for shortcuts"
          className="w-full resize-none bg-transparent text-foreground text-xs leading-relaxed placeholder:text-muted-foreground focus:outline-none"
        />

        <div className="mt-2 flex items-center justify-between gap-2">
          <div className="flex min-w-0 items-center gap-1.5">
            <button
              type="button"
              onClick={() => setMode(mode === 'ask' ? 'agent' : 'ask')}
              title={
                mode === 'ask'
                  ? 'Ask Mode: Proposes actions for confirmation'
                  : 'Agent Mode: Automatically executes safe MCP actions'
              }
              className="inline-flex h-7 shrink-0 items-center gap-1.5 rounded-full border border-border/80 bg-muted/40 px-2.5 font-medium text-[11px] text-foreground transition-colors hover:bg-muted"
            >
              {mode === 'ask' ? (
                <>
                  <Edit3 className="size-3 text-muted-foreground" />
                  <span>Ask</span>
                </>
              ) : (
                <>
                  <Zap className="size-3 text-amber-400" />
                  <span>Agent</span>
                </>
              )}
            </button>

            <span
              title={`${providerStatus.activeProviderName} (${providerStatus.activeModel})`}
              className="truncate rounded-md bg-muted/40 px-2 py-0.5 font-mono text-[10px] text-muted-foreground"
            >
              {providerStatus.activeProviderName}
            </span>
          </div>

          <div className="flex shrink-0 items-center gap-2">
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <button
                  type="button"
                  aria-label="AI and MCP settings"
                  className="inline-flex size-7 items-center justify-center rounded-lg text-muted-foreground transition-colors hover:bg-muted hover:text-foreground"
                >
                  <SlidersHorizontal className="size-3.5" />
                </button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end" className="w-52 p-1.5">
                <DropdownMenuLabel className="font-normal text-[11px] text-muted-foreground">
                  AI & MCP Configuration
                </DropdownMenuLabel>
                <DropdownMenuSeparator />
                <DropdownMenuItem asChild>
                  <Link
                    to="/settings"
                    search={{ tab: 'ai' } as never}
                    className="cursor-pointer text-xs"
                  >
                    Configure AI Models & Keys
                  </Link>
                </DropdownMenuItem>
                <DropdownMenuItem asChild>
                  <Link
                    to="/settings"
                    search={{ tab: 'mcp' } as never}
                    className="cursor-pointer text-xs"
                  >
                    MCP Server & IDE Plugins
                  </Link>
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>

            <button
              type="button"
              onClick={handleSubmit}
              disabled={!text.trim() || disabled}
              aria-label="Send message"
              className="inline-flex size-7 items-center justify-center rounded-full bg-primary text-primary-foreground transition-all hover:opacity-90 disabled:opacity-40"
            >
              <ArrowUp className="size-3.5" />
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
