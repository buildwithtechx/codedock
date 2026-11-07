import { useRouterState } from '@tanstack/react-router';
import { ChevronDown, Maximize2, MessageSquare, Minimize2, Plus, Trash2, X } from 'lucide-react';
import { useCallback, useEffect } from 'react';
import { toast } from 'sonner';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@codedock/ui';
import { useAiProviderStatus, useSendAiChat } from '#/hooks/use-ai';
import type { AiResourceTag } from '#/interfaces/ai';
import { useAiSidebarStore } from '#/stores/ai-sidebar-store';
import { AskAiComposer } from './ask-ai-composer';
import { AskAiMessages } from './ask-ai-messages';
import { AskAiWelcome } from './ask-ai-welcome';

export function AskAiSidebar() {
  const pathname = useRouterState({
    select: (state) => state.location.pathname,
  });
  const {
    isOpen,
    isExpanded,
    mode,
    threads,
    activeThreadId,
    toggleSidebar,
    closeSidebar,
    toggleExpanded,
    startNewThread,
    selectThread,
    deleteThread,
    appendMessage,
  } = useAiSidebarStore();

  const { status: providerStatus } = useAiProviderStatus();
  const sendChatMutation = useSendAiChat();
  const isSending = sendChatMutation.isPending;

  const activeThread = threads.find((t) => t.id === activeThreadId) || threads[0];
  const hasMessages = (activeThread?.messages.length || 0) > 0;

  useEffect(() => {
    const onKeyDown = (e: KeyboardEvent) => {
      if ((e.metaKey || e.ctrlKey) && e.key.toLowerCase() === 'i') {
        e.preventDefault();
        toggleSidebar();
      }
    };
    window.addEventListener('keydown', onKeyDown);
    return () => window.removeEventListener('keydown', onKeyDown);
  }, [toggleSidebar]);

  const handleSendMessage = useCallback(
    async (text: string, mentions: AiResourceTag[]) => {
      if (!text.trim() || isSending) return;
      const priorMessages = (activeThread?.messages || []).slice(-12).map((m) => ({
        role: m.role,
        content: m.content,
      }));

      appendMessage({
        id: `msg_${Date.now().toString(36)}_u`,
        role: 'user',
        content: text,
        createdAt: new Date().toISOString(),
        mentions: mentions.length > 0 ? mentions : undefined,
      });

      try {
        const payload = await sendChatMutation.mutateAsync({
          messages: [...priorMessages, { role: 'user', content: text }],
          mode,
          provider: providerStatus.activeProviderId ?? undefined,
          model: providerStatus.activeProviderId != null ? providerStatus.activeModel : undefined,
          currentPath: pathname,
          mentions,
        });
        appendMessage({
          id: `msg_${Date.now().toString(36)}_a`,
          role: 'assistant',
          content:
            payload?.reply || 'I inspected your Codedock instance and completed the request.',
          createdAt: new Date().toISOString(),
          steps: payload?.steps,
          pendingActions: payload?.pendingActions,
        });
      } catch (err) {
        const message =
          err instanceof Error ? err.message : 'Unable to reach Codedock AI assistant';
        toast.error(message);
        appendMessage({
          id: `msg_${Date.now().toString(36)}_err`,
          role: 'assistant',
          content: `I encountered an error processing that request: **${message}**. Check **[Settings → AI Models](/settings?tab=ai)** or try again.`,
          createdAt: new Date().toISOString(),
        });
      }
    },
    [
      activeThread?.messages,
      appendMessage,
      isSending,
      mode,
      pathname,
      providerStatus.activeModel,
      providerStatus.activeProviderId,
      sendChatMutation,
    ]
  );

  if (!isOpen) return null;

  return (
    <aside
      aria-label="Ask AI Assistant Sidebar"
      className={`flex h-dvh shrink-0 flex-col border-sidebar-border/80 border-l bg-sidebar text-sidebar-foreground transition-[width] duration-200 ease-out ${
        isExpanded ? 'w-full sm:w-135 lg:w-155' : 'w-full sm:w-97.5 lg:w-102.5'
      }`}
    >
      {/* Top Header Bar */}
      <div className="flex h-14 shrink-0 items-center justify-between border-sidebar-border/80 border-b px-3.5">
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <button
              type="button"
              className="inline-flex max-w-57.5 items-center gap-1.5 rounded-md px-2 py-1 font-semibold text-sidebar-foreground text-xs transition-colors hover:bg-sidebar-accent"
            >
              <span className="truncate">{activeThread?.title || 'New conversation'}</span>
              <ChevronDown className="size-3.5 shrink-0 text-muted-foreground" />
            </button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="start" className="w-64 p-1">
            <DropdownMenuLabel className="font-normal text-[11px] text-muted-foreground">
              Recent conversations
            </DropdownMenuLabel>
            <DropdownMenuSeparator />
            {threads.map((thread) => (
              <DropdownMenuItem
                key={thread.id}
                onSelect={() => selectThread(thread.id)}
                className="flex cursor-pointer items-center justify-between gap-2 py-1.5 text-xs"
              >
                <div className="flex min-w-0 items-center gap-2">
                  <MessageSquare className="size-3.5 shrink-0 text-muted-foreground" />
                  <span className="truncate">{thread.title}</span>
                </div>
                {threads.length > 1 && (
                  <button
                    type="button"
                    onClick={(e) => {
                      e.stopPropagation();
                      deleteThread(thread.id);
                    }}
                    aria-label="Delete conversation"
                    className="rounded p-0.5 text-muted-foreground hover:bg-destructive/10 hover:text-destructive"
                  >
                    <Trash2 className="size-3" />
                  </button>
                )}
              </DropdownMenuItem>
            ))}
          </DropdownMenuContent>
        </DropdownMenu>

        <div className="flex items-center gap-1">
          <button
            type="button"
            onClick={() => startNewThread()}
            title="New conversation"
            aria-label="New conversation"
            className="flex size-7 items-center justify-center rounded-md text-muted-foreground transition-colors hover:bg-sidebar-accent hover:text-sidebar-foreground"
          >
            <Plus className="size-4" />
          </button>
          <button
            type="button"
            onClick={toggleExpanded}
            title={isExpanded ? 'Restore sidebar width' : 'Expand sidebar width'}
            aria-label={isExpanded ? 'Restore sidebar width' : 'Expand sidebar width'}
            className="hidden size-7 items-center justify-center rounded-md text-muted-foreground transition-colors hover:bg-sidebar-accent hover:text-sidebar-foreground sm:flex"
          >
            {isExpanded ? <Minimize2 className="size-3.5" /> : <Maximize2 className="size-3.5" />}
          </button>
          <button
            type="button"
            onClick={closeSidebar}
            title="Close Ask AI (⌘I)"
            aria-label="Close Ask AI"
            className="flex size-7 items-center justify-center rounded-md text-muted-foreground transition-colors hover:bg-sidebar-accent hover:text-sidebar-foreground"
          >
            <X className="size-4" />
          </button>
        </div>
      </div>

      {/* Main Content Area: Welcome State or Active Conversation */}
      <div className="flex min-h-0 flex-1 flex-col overflow-hidden">
        {hasMessages && activeThread ? (
          <AskAiMessages messages={activeThread.messages} isThinking={isSending} />
        ) : (
          <AskAiWelcome
            providerStatus={providerStatus}
            onSelectPrompt={(prompt) => void handleSendMessage(prompt, [])}
          />
        )}
      </div>

      {/* Bottom Composer */}
      <AskAiComposer
        disabled={isSending}
        providerStatus={providerStatus}
        onSend={handleSendMessage}
      />
    </aside>
  );
}
