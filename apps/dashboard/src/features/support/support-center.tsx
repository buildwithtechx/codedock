import { useState, useEffect } from 'react';
import {
  ArrowRight,
  BookOpen,
  CreditCard,
  HelpCircle,
  Loader2,
  Mail,
  Plus,
  Rocket,
  ShieldCheck,
  User,
  Users,
} from 'lucide-react';
import { SupportTicketList } from './support-ticket-list';
import { SupportConversation } from './support-conversation';
import { NewTicketModal } from './new-ticket-modal';
import { Button } from '#/components/ui/button';
import { PageContainer } from '#/components/ui/page-container';
import {
  getTickets,
  createTicket,
  replyTicket,
  type SupportTicket,
  type CreateTicketInput,
  type SupportCategory,
} from '#/services/cloud-support';
import { cn } from '#/lib/utils';

export interface SupportCenterProps {
  className?: string;
}

const SUPPORT_TOPICS: { category: SupportCategory; label: string; description: string; icon: typeof Rocket }[] = [
  {
    category: 'deployment',
    label: 'Deployments',
    description: 'Build failures, runtime logs, cache mounts, and Docker containers.',
    icon: Rocket,
  },
  {
    category: 'billing',
    label: 'Billing',
    description: 'Invoices, plans, compute quota, and payment methods.',
    icon: CreditCard,
  },
  {
    category: 'account',
    label: 'Account & Security',
    description: 'Organization membership, access tokens, 2FA, and permissions.',
    icon: User,
  },
  {
    category: 'general',
    label: 'General Guidance',
    description: 'Platform architecture, migrations, and best practices.',
    icon: HelpCircle,
  },
];

export function SupportCenter({ className }: SupportCenterProps) {
  const [tickets, setTickets] = useState<SupportTicket[]>([]);
  const [selectedTicketId, setSelectedTicketId] = useState<string | null>(null);
  const [isNewModalOpen, setIsNewModalOpen] = useState(false);
  const [initialCategory, setInitialCategory] = useState<SupportCategory>('general');
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

  const openNewTicketWithCategory = (cat: SupportCategory) => {
    setInitialCategory(cat);
    setIsNewModalOpen(true);
  };

  return (
    <PageContainer className={cn('@container/support-center space-y-6 text-left', className)}>
      {/* Header matching OpenShip */}
      <header className="flex flex-wrap items-center justify-between gap-4">
        <div className="flex items-center gap-3.5">
          <span className="flex size-11 shrink-0 items-center justify-center rounded-2xl border border-border/60 bg-card text-muted-foreground shadow-xs">
            <HelpCircle className="size-6 text-primary" />
          </span>
          <div>
            <h1 className="text-2xl font-semibold tracking-tight text-foreground">Support Center</h1>
            <p className="mt-0.5 text-sm text-muted-foreground">
              Connect directly with CodeDock platform engineers for deployment, billing, and technical guidance.
            </p>
          </div>
        </div>
        <Button onClick={() => openNewTicketWithCategory('general')} className="gap-1.5 rounded-xl">
          <Plus className="size-4" />
          <span>New ticket</span>
        </Button>
      </header>

      {loading ? (
        <div className="flex items-center justify-center gap-2 py-16">
          <Loader2 className="size-5 animate-spin text-muted-foreground" />
          <p className="text-muted-foreground text-sm">Loading tickets...</p>
        </div>
      ) : (
        <div className="grid items-start gap-6 @[820px]/support-center:grid-cols-[340px_minmax(0,1fr)]">
          {/* Left Column: Tickets List */}
          <div className={selectedTicket ? 'hidden @[820px]/support-center:block' : ''}>
            <SupportTicketList
              tickets={tickets}
              selectedTicketId={selectedTicketId}
              onSelectTicket={setSelectedTicketId}
              onNewTicket={() => openNewTicketWithCategory('general')}
            />
          </div>

          {/* Right Column: Conversation or Topic Overview */}
          <div className="min-w-0">
            {selectedTicket ? (
              <SupportConversation
                ticket={selectedTicket}
                onBack={() => setSelectedTicketId(null)}
                onSendReply={handleReply}
              />
            ) : (
              <>
                <section className="space-y-6 rounded-2xl border border-border/50 bg-card p-6 sm:p-8 shadow-xs">
                  <div className="flex size-12 items-center justify-center rounded-2xl bg-primary/10 text-primary">
                    <Mail className="size-6" />
                  </div>
                  <div>
                    <h2 className="text-xl font-semibold tracking-tight text-foreground">
                      Need help with CodeDock?
                    </h2>
                    <p className="mt-1 text-sm text-muted-foreground">
                      Select a topic below to submit a priority support request directly to our platform engineers.
                    </p>
                  </div>

                  <div className="grid gap-3 sm:grid-cols-2">
                    {SUPPORT_TOPICS.map((topic) => (
                      <button
                        key={topic.category}
                        type="button"
                        onClick={() => openNewTicketWithCategory(topic.category)}
                        className="group flex items-start gap-3.5 rounded-xl border border-border/50 bg-muted/30 p-4 text-left transition-all hover:border-primary/40 hover:bg-muted/60 hover:shadow-xs cursor-pointer"
                      >
                        <topic.icon className="mt-0.5 size-5 shrink-0 text-muted-foreground group-hover:text-primary transition-colors" />
                        <div className="min-w-0 flex-1">
                          <div className="flex items-center justify-between">
                            <p className="text-sm font-semibold text-foreground">{topic.label}</p>
                            <ArrowRight className="size-4 text-muted-foreground/40 transition-transform group-hover:translate-x-0.5 group-hover:text-foreground" />
                          </div>
                          <p className="mt-1 line-clamp-2 text-xs text-muted-foreground leading-relaxed">
                            {topic.description}
                          </p>
                        </div>
                      </button>
                    ))}
                  </div>
                </section>

                <section
                  className="mt-6 rounded-2xl border border-border/50 bg-card p-5 sm:p-6 shadow-xs"
                  aria-labelledby="support-resources"
                >
                  <h2 id="support-resources" className="mb-4 text-sm font-semibold text-foreground">
                    Resources
                  </h2>
                  <div className="grid gap-4 sm:grid-cols-2">
                    <a
                      href="https://docs.codedock.io"
                      target="_blank"
                      rel="noreferrer"
                      className="group flex items-start gap-3 rounded-lg p-2 transition-colors hover:bg-muted/40"
                    >
                      <BookOpen className="mt-0.5 size-4 shrink-0 text-muted-foreground group-hover:text-foreground" />
                      <div>
                        <span className="text-sm font-medium text-foreground group-hover:underline">
                          Documentation
                        </span>
                        <span className="mt-0.5 block text-xs leading-5 text-muted-foreground">
                          Read architecture guides, API references, and deployment cookbooks.
                        </span>
                      </div>
                    </a>
                    <a
                      href="https://github.com/codedock/codedock/discussions"
                      target="_blank"
                      rel="noreferrer"
                      className="group flex items-start gap-3 rounded-lg p-2 transition-colors hover:bg-muted/40"
                    >
                      <Users className="mt-0.5 size-4 shrink-0 text-muted-foreground group-hover:text-foreground" />
                      <div>
                        <span className="text-sm font-medium text-foreground group-hover:underline">
                          Community Forum
                        </span>
                        <span className="mt-0.5 block text-xs leading-5 text-muted-foreground">
                          Ask questions, share workflows, and discuss features with other engineers.
                        </span>
                      </div>
                    </a>
                  </div>
                </section>

                <div className="mt-5 flex flex-wrap items-start justify-between gap-3 px-1 text-xs leading-5 text-muted-foreground">
                  <p className="flex max-w-md items-start gap-2">
                    <ShieldCheck className="mt-0.5 size-3.5 shrink-0 text-emerald-500" />
                    <span>Your tickets and diagnostic data are encrypted and visible only to your organization and authorized engineers.</span>
                  </p>
                  <a
                    href="mailto:support@codedock.io"
                    className="shrink-0 font-medium underline-offset-4 hover:underline hover:text-foreground"
                  >
                    Email support: support@codedock.io
                  </a>
                </div>
              </>
            )}
          </div>
        </div>
      )}

      {isNewModalOpen && (
        <NewTicketModal
          onClose={() => setIsNewModalOpen(false)}
          onSubmit={handleCreateTicket}
          defaultCategory={initialCategory}
        />
      )}
    </PageContainer>
  );
}

