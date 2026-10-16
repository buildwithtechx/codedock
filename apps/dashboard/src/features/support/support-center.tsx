import { useState, useEffect } from 'react';
import { HelpCircle, Loader2 } from 'lucide-react';
import { SupportTicketList } from './support-ticket-list';
import { SupportConversation } from './support-conversation';
import { NewTicketModal } from './new-ticket-modal';
import {
  getTickets,
  createTicket,
  replyTicket,
  type SupportTicket,
  type CreateTicketInput,
} from '#/services/cloud-support';
import { cn } from '#/lib/utils';

export interface SupportCenterProps {
  className?: string;
}

export function SupportCenter({ className }: SupportCenterProps) {
  const [tickets, setTickets] = useState<SupportTicket[]>([]);
  const [selectedTicketId, setSelectedTicketId] = useState<string | null>(null);
  const [isNewModalOpen, setIsNewModalOpen] = useState(false);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    let mounted = true;
    void (async () => {
      try {
        const list = await getTickets();
        if (mounted) setTickets(list);
      } finally {
        if (mounted) setLoading(false);
      }
    })();
    return () => {
      mounted = false;
    };
  }, []);

  const selectedTicket = tickets.find((t) => t.id === selectedTicketId) || null;

  async function handleCreateTicket(input: CreateTicketInput) {
    const newTicket = await createTicket(input);
    setTickets((prev) => [newTicket, ...prev]);
    setSelectedTicketId(newTicket.id);
  }

  async function handleReply(message: string) {
    if (!selectedTicketId) return;
    const updated = await replyTicket(selectedTicketId, message);
    setTickets((prev) => prev.map((t) => (t.id === updated.id ? updated : t)));
  }

  return (
    <div className={cn('mx-auto max-w-4xl space-y-6 p-6 text-left', className)}>
      <div className="flex items-center gap-3 border-b border-border/60 pb-5">
        <div className="flex size-10 items-center justify-center rounded-xl bg-primary/10 text-primary">
          <HelpCircle className="size-5" />
        </div>
        <div>
          <h2 className="font-bold text-foreground text-xl tracking-tight">Support Center</h2>
          <p className="text-muted-foreground text-xs">
            Connect directly with CodeDock platform engineers for deployment, billing, and technical guidance.
          </p>
        </div>
      </div>

      {loading ? (
        <div className="flex items-center justify-center gap-2 py-16">
          <Loader2 className="size-5 animate-spin text-muted-foreground" />
          <p className="text-muted-foreground text-sm">Loading tickets...</p>
        </div>
      ) : selectedTicket ? (
        <SupportConversation
          ticket={selectedTicket}
          onBack={() => setSelectedTicketId(null)}
          onSendReply={handleReply}
        />
      ) : (
        <SupportTicketList
          tickets={tickets}
          selectedTicketId={selectedTicketId}
          onSelectTicket={setSelectedTicketId}
          onNewTicket={() => setIsNewModalOpen(true)}
        />
      )}

      {isNewModalOpen && (
        <NewTicketModal
          onClose={() => setIsNewModalOpen(false)}
          onSubmit={handleCreateTicket}
        />
      )}
    </div>
  );
}
