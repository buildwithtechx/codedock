import { ArrowLeft, Paperclip, Send, ShieldCheck, User } from 'lucide-react';
import { useState } from 'react';
import { Button } from '#/components/ui/button';
import type { SupportTicket } from '#/interfaces/support';
import { cn } from '#/lib/utils';

export interface SupportConversationProps {
  ticket: SupportTicket;
  onBack: () => void;
  onSendReply: (message: string) => Promise<void> | void;
  className?: string;
}

export function SupportConversation({
  ticket,
  onBack,
  onSendReply,
  className,
}: SupportConversationProps) {
  const [replyText, setReplyText] = useState('');
  const [sending, setSending] = useState(false);

  async function handleSend() {
    const trimmed = replyText.trim();
    if (!trimmed || sending) return;
    try {
      setSending(true);
      await onSendReply(trimmed);
      setReplyText('');
    } finally {
      setSending(false);
    }
  }

  return (
    <div className={cn('flex flex-col space-y-5 text-left', className)}>
      <div className="flex flex-wrap items-center justify-between gap-3 border-border/60 border-b pb-4">
        <div className="flex items-center gap-3">
          <Button variant="ghost" size="sm" onClick={onBack} className="gap-1 p-2">
            <ArrowLeft className="size-4" />
            <span className="text-xs">Back</span>
          </Button>
          <div>
            <div className="flex items-center gap-2">
              <span className="font-mono text-muted-foreground text-xs">{ticket.id}</span>
              <span
                className={cn(
                  'rounded-full px-2 py-0.5 font-semibold text-[10px] uppercase',
                  ticket.status === 'open'
                    ? 'bg-amber-500/15 text-amber-600 dark:text-amber-400'
                    : ticket.status === 'resolved'
                      ? 'bg-emerald-500/15 text-emerald-600 dark:text-emerald-400'
                      : 'bg-muted text-muted-foreground'
                )}
              >
                {ticket.status}
              </span>
              {ticket.hasDiagnosticBundle && (
                <span className="inline-flex items-center gap-1 rounded-md bg-muted/60 px-1.5 py-0.5 text-[10px] text-muted-foreground">
                  <Paperclip className="size-2.5" />
                  <span>Diagnostics attached</span>
                </span>
              )}
            </div>
            <h3 className="mt-1 font-semibold text-base text-foreground">{ticket.subject}</h3>
          </div>
        </div>
      </div>

      <div className="space-y-4 py-2">
        {ticket.messages.map((msg) => {
          const isStaff = msg.sender === 'staff';
          return (
            <div
              key={msg.id}
              className={cn(
                'flex flex-col rounded-2xl border p-4 transition-colors',
                isStaff
                  ? 'border-primary/20 bg-primary/5 text-foreground'
                  : 'border-border/60 bg-card text-foreground'
              )}
            >
              <div className="flex items-center justify-between gap-2 border-border/40 border-b pb-2">
                <div className="flex items-center gap-1.5 font-medium text-xs">
                  {isStaff ? (
                    <>
                      <ShieldCheck className="size-3.5 text-primary" />
                      <span className="font-semibold text-primary">CodeDock Staff</span>
                    </>
                  ) : (
                    <>
                      <User className="size-3.5 text-muted-foreground" />
                      <span>{msg.senderName || 'You'}</span>
                    </>
                  )}
                </div>
                <span className="text-[11px] text-muted-foreground">{msg.createdAt}</span>
              </div>
              <div className="mt-3 whitespace-pre-wrap font-sans text-sm leading-relaxed">
                {msg.content}
              </div>
            </div>
          );
        })}
      </div>

      {ticket.status === 'closed' ? (
        <div className="rounded-xl border border-border/60 bg-muted/30 p-4 text-center text-muted-foreground text-xs">
          This ticket is marked as closed. You can submit a new ticket if you need further help.
        </div>
      ) : (
        <div className="space-y-3 pt-2">
          <textarea
            value={replyText}
            onChange={(e) => setReplyText(e.target.value)}
            placeholder="Type your response or update..."
            rows={3}
            className="w-full resize-none rounded-xl border border-input bg-background p-3 text-sm outline-none focus-visible:ring-2 focus-visible:ring-ring"
          />
          <div className="flex justify-end">
            <Button
              size="sm"
              disabled={!replyText.trim() || sending}
              onClick={handleSend}
              className="gap-1.5"
            >
              <Send className="size-3.5" />
              <span>{sending ? 'Sending...' : 'Send response'}</span>
            </Button>
          </div>
        </div>
      )}
    </div>
  );
}
