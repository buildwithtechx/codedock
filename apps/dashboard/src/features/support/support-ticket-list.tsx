import { MessageSquare, Plus, Tag } from 'lucide-react';
import { Button } from '#/components/ui/button';
import type { SupportCategory, SupportTicket } from '#/interfaces/support';
import { cn } from '#/lib/utils';

export interface SupportTicketListProps {
  tickets: SupportTicket[];
  selectedTicketId?: string | null;
  onSelectTicket: (ticketId: string) => void;
  onNewTicket: () => void;
  className?: string;
}

const CATEGORY_LABELS: Record<SupportCategory, string> = {
  deployment: 'Deployments',
  billing: 'Billing',
  account: 'Account',
  general: 'General',
  technical: 'Technical',
  urgent: 'Urgent',
};

export function SupportTicketList({
  tickets,
  selectedTicketId,
  onSelectTicket,
  onNewTicket,
  className,
}: SupportTicketListProps) {
  return (
    <div className={cn('flex flex-col space-y-4 text-left', className)}>
      <div className="flex items-center justify-between">
        <div>
          <h3 className="font-semibold text-foreground text-lg">Support Tickets</h3>
          <p className="text-muted-foreground text-xs">
            {tickets.length} total conversation{tickets.length === 1 ? '' : 's'}
          </p>
        </div>
        <Button size="sm" onClick={onNewTicket} className="gap-1.5">
          <Plus className="size-4" />
          <span>New ticket</span>
        </Button>
      </div>

      {tickets.length === 0 ? (
        <div className="flex flex-col items-center justify-center rounded-2xl border border-border/80 border-dashed p-10 text-center">
          <div className="flex size-10 items-center justify-center rounded-xl bg-muted/60 text-muted-foreground">
            <MessageSquare className="size-5" />
          </div>
          <p className="mt-3 font-medium text-foreground text-sm">No tickets found</p>
          <p className="mt-1 text-muted-foreground text-xs">
            Have a question or running into an issue? Create a support ticket to reach the team.
          </p>
          <Button variant="outline" size="sm" onClick={onNewTicket} className="mt-4 gap-1.5">
            <Plus className="size-3.5" />
            <span>Open your first ticket</span>
          </Button>
        </div>
      ) : (
        <div className="space-y-2">
          {tickets.map((ticket) => {
            const isSelected = selectedTicketId === ticket.id;
            return (
              <div
                key={ticket.id}
                onClick={() => onSelectTicket(ticket.id)}
                className={cn(
                  'group flex cursor-pointer items-start justify-between gap-4 rounded-xl border p-4 transition-all',
                  isSelected
                    ? 'border-primary bg-primary/5 ring-1 ring-primary/20'
                    : 'border-border/60 bg-card hover:border-border hover:bg-muted/30'
                )}
              >
                <div className="min-w-0 flex-1">
                  <div className="flex items-center gap-2">
                    <span className="font-medium font-mono text-muted-foreground text-xs">
                      {ticket.id}
                    </span>
                    <span
                      className={cn(
                        'rounded-full px-2 py-0.5 font-semibold text-[10px] uppercase tracking-wider',
                        ticket.status === 'open'
                          ? 'bg-amber-500/15 text-amber-600 dark:text-amber-400'
                          : ticket.status === 'resolved'
                            ? 'bg-emerald-500/15 text-emerald-600 dark:text-emerald-400'
                            : 'bg-muted text-muted-foreground'
                      )}
                    >
                      {ticket.status}
                    </span>
                    <span className="inline-flex items-center gap-1 text-[11px] text-muted-foreground">
                      <Tag className="size-3" />
                      <span>{CATEGORY_LABELS[ticket.category] || ticket.category}</span>
                    </span>
                  </div>

                  <p className="mt-1.5 truncate font-medium text-foreground text-sm">
                    {ticket.subject}
                  </p>
                </div>

                <div className="flex shrink-0 items-center gap-1.5 text-muted-foreground text-xs">
                  <MessageSquare className="size-3.5" />
                  <span>{ticket.messages.length}</span>
                </div>
              </div>
            );
          })}
        </div>
      )}
    </div>
  );
}
