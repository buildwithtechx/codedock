import { AlertTriangle, CheckCircle2, Loader2, Play, Terminal, XCircle } from 'lucide-react';
import { useEffect, useRef, useState } from 'react';
import { toast } from 'sonner';
import { Avatar, AvatarFallback, AvatarImage } from '@codedock/ui/avatar';
import { Button } from '@codedock/ui/button';
import { CodedockLogo } from '#/components/layout/codedock-logo';
import { useConfirmAiAction } from '#/hooks/use-ai';
import type { AiPendingAction, AiTurnMessage } from '#/interfaces/ai';
import { useAiSidebarStore } from '#/stores/ai-sidebar-store';

export function AskAiMessages({
  messages,
  isThinking,
}: {
  messages: AiTurnMessage[];
  isThinking: boolean;
}) {
  const bottomRef = useRef<HTMLDivElement>(null);
  const confirmActionMutation = useConfirmAiAction();
  const updateActionStatus = useAiSidebarStore((s) => s.updateActionStatus);
  const [executingId, setExecutingId] = useState<string | null>(null);

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages, isThinking]);

  const handleConfirmAction = async (messageId: string, action: AiPendingAction) => {
    setExecutingId(action.actionId);
    try {
      await confirmActionMutation.mutateAsync(action.actionId);
      updateActionStatus(messageId, action.actionId, 'approved');
      toast.success(`Executed ${action.tool}`);
    } catch (err) {
      updateActionStatus(messageId, action.actionId, 'failed');
      toast.error(err instanceof Error ? err.message : 'Action execution failed');
    } finally {
      setExecutingId(null);
    }
  };

  const handleCancelAction = (messageId: string, actionId: string) => {
    updateActionStatus(messageId, actionId, 'cancelled');
    toast.info('Action cancelled');
  };

  return (
    <div className="flex-1 space-y-4 overflow-y-auto px-4 py-4">
      {messages.map((msg) => {
        const isUser = msg.role === 'user';
        return (
          <div
            key={msg.id}
            className={`flex flex-col ${isUser ? 'items-end' : 'items-start'} space-y-2`}
          >
            {isUser ? (
              <div className="max-w-[88%] rounded-2xl rounded-br-sm border border-primary/25 bg-primary/10 px-3.5 py-2.5 text-foreground text-xs leading-relaxed">
                {msg.mentions && msg.mentions.length > 0 && (
                  <div className="mb-1.5 flex flex-wrap gap-1">
                    {msg.mentions.map((m) => (
                      <span
                        key={`${m.type}-${m.id}`}
                        className="inline-flex items-center gap-1 rounded-md bg-primary/20 px-1.5 py-0.5 font-mono text-[10px] text-primary"
                      >
                        @{m.type}:{m.name}
                      </span>
                    ))}
                  </div>
                )}
                <p className="whitespace-pre-wrap">{msg.content}</p>
              </div>
            ) : (
              <div className="w-full space-y-2.5">
                <div className="flex items-center gap-1.5 text-[11px] text-muted-foreground">
                  <Avatar size="sm" className="size-4 rounded-md">
                    <AvatarImage src="/favicon.svg" alt="Codedock Copilot" className="rounded-md" />
                    <AvatarFallback className="rounded-md bg-transparent p-0">
                      <CodedockLogo className="size-4 rounded-md" />
                    </AvatarFallback>
                  </Avatar>
                  <span className="font-medium text-foreground">Codedock Copilot</span>
                </div>

                {/* Executed MCP Tool Steps */}
                {msg.steps && msg.steps.length > 0 && (
                  <div className="space-y-1">
                    {msg.steps.map((step) => (
                      <div
                        key={step.id}
                        className="flex items-center gap-2 rounded-lg border border-border/50 bg-muted/20 px-2.5 py-1.5 text-[11px] text-muted-foreground"
                      >
                        {step.status === 'done' ? (
                          <CheckCircle2 className="size-3.5 shrink-0 text-emerald-500" />
                        ) : step.status === 'error' ? (
                          <XCircle className="size-3.5 shrink-0 text-destructive" />
                        ) : (
                          <Terminal className="size-3.5 shrink-0 text-amber-500" />
                        )}
                        <code className="font-mono text-[10px] text-foreground">{step.tool}</code>
                        <span className="truncate">{step.label}</span>
                      </div>
                    ))}
                  </div>
                )}

                {/* Assistant Markdown / Text Body */}
                <div className="rounded-2xl rounded-tl-sm border border-border/60 bg-card px-3.5 py-3 text-foreground text-xs leading-relaxed">
                  <div className="whitespace-pre-wrap">{msg.content}</div>
                </div>

                {/* Pending Action Confirmation Cards */}
                {msg.pendingActions?.map((action) => {
                  const isRunning = executingId === action.actionId;
                  return (
                    <div
                      key={action.actionId}
                      className="rounded-xl border border-amber-500/30 bg-amber-500/5 p-3 text-xs"
                    >
                      <div className="flex items-start justify-between gap-2">
                        <div className="min-w-0">
                          <div className="flex items-center gap-1.5">
                            {action.destructive ? (
                              <AlertTriangle className="size-3.5 shrink-0 text-amber-500" />
                            ) : (
                              <Play className="size-3.5 shrink-0 text-primary" />
                            )}
                            <span className="font-semibold text-foreground">{action.title}</span>
                          </div>
                          <p className="mt-1 text-[11px] text-muted-foreground">
                            {action.description}
                          </p>
                          <p className="mt-1 font-mono text-[10px] text-muted-foreground/80">
                            {action.method} {action.path}
                          </p>
                        </div>
                      </div>

                      <div className="mt-3 flex items-center gap-2">
                        {action.resolvedStatus === 'approved' ? (
                          <span className="inline-flex items-center gap-1 font-medium text-[11px] text-emerald-500">
                            <CheckCircle2 className="size-3.5" />
                            Executed successfully
                          </span>
                        ) : action.resolvedStatus === 'cancelled' ? (
                          <span className="inline-flex items-center gap-1 text-[11px] text-muted-foreground">
                            <XCircle className="size-3.5" />
                            Cancelled
                          </span>
                        ) : (
                          <>
                            <Button
                              size="sm"
                              disabled={isRunning}
                              onClick={() => void handleConfirmAction(msg.id, action)}
                              className="h-7 gap-1.5 px-2.5 text-[11px]"
                            >
                              {isRunning ? (
                                <Loader2 className="size-3 animate-spin" />
                              ) : (
                                <Play className="size-3" />
                              )}
                              Approve & Execute
                            </Button>
                            <Button
                              size="sm"
                              variant="outline"
                              disabled={isRunning}
                              onClick={() => handleCancelAction(msg.id, action.actionId)}
                              className="h-7 px-2.5 text-[11px]"
                            >
                              Cancel
                            </Button>
                          </>
                        )}
                      </div>
                    </div>
                  );
                })}
              </div>
            )}
          </div>
        );
      })}

      {isThinking && (
        <div className="flex items-center gap-2 rounded-xl border border-border/50 bg-card px-3.5 py-2.5 text-muted-foreground text-xs">
          <Loader2 className="size-3.5 animate-spin text-primary" />
          <span>Inspecting instance state & MCP tools...</span>
        </div>
      )}

      <div ref={bottomRef} />
    </div>
  );
}
