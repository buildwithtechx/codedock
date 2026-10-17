import { Paperclip, Send, X } from 'lucide-react';
import { useState } from 'react';
import { Button } from '#/components/ui/button';
import { cn } from '#/lib/utils';
import type { CreateTicketInput, SupportCategory } from '#/interfaces/support';

export interface NewTicketModalProps {
  onClose: () => void;
  onSubmit: (input: CreateTicketInput) => Promise<void> | void;
  defaultCategory?: SupportCategory;
  className?: string;
}

const CATEGORIES: { id: SupportCategory; label: string }[] = [
  { id: 'deployment', label: 'Deployments & Builds' },
  { id: 'billing', label: 'Billing & Subscriptions' },
  { id: 'account', label: 'Account & Teams' },
  { id: 'general', label: 'General Questions' },
];

export function NewTicketModal({
  onClose,
  onSubmit,
  defaultCategory,
  className,
}: NewTicketModalProps) {
  const [subject, setSubject] = useState('');
  const [category, setCategory] = useState<SupportCategory>(defaultCategory ?? 'deployment');
  const [message, setMessage] = useState('');
  const [attachDiagnostics, setAttachDiagnostics] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (!subject.trim() || !message.trim() || submitting) return;
    try {
      setSubmitting(true);
      setError(null);
      await onSubmit({
        subject: subject.trim(),
        category,
        message: message.trim(),
        attachDiagnostics,
      });
      onClose();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to submit ticket');
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <div
      role="dialog"
      aria-modal="true"
      aria-labelledby="new-ticket-title"
      className="fixed inset-0 z-50 flex items-center justify-center bg-background/80 p-4 backdrop-blur-sm"
    >
      <div
        className={cn(
          'relative w-full max-w-lg rounded-2xl border border-border/60 bg-card p-6 text-left shadow-2xl',
          className
        )}
      >
        <button
          onClick={onClose}
          aria-label="Close modal"
          className="absolute top-4 right-4 rounded-lg p-1.5 text-muted-foreground hover:bg-muted hover:text-foreground"
        >
          <X className="size-4" />
        </button>

        <h3 id="new-ticket-title" className="font-bold text-foreground text-lg">
          Create Support Ticket
        </h3>
        <p className="mt-1 text-muted-foreground text-xs leading-relaxed">
          Reach out to the CodeDock team. We typically respond within a few hours.
        </p>

        {error && (
          <p
            role="alert"
            className="mt-3 rounded-lg bg-destructive/10 p-2 text-destructive text-xs"
          >
            {error}
          </p>
        )}

        <form onSubmit={handleSubmit} className="mt-4 space-y-4">
          <div>
            <label className="block font-medium text-foreground text-xs">Subject</label>
            <input
              type="text"
              required
              value={subject}
              onChange={(e) => setSubject(e.target.value)}
              placeholder="e.g. Build timeout on Node 22 service"
              className="mt-1.5 w-full rounded-xl border border-input bg-background px-3 py-2 text-sm outline-none focus-visible:ring-2 focus-visible:ring-ring"
            />
          </div>

          <div>
            <label className="block font-medium text-foreground text-xs">Category</label>
            <select
              value={category}
              onChange={(e) => setCategory(e.target.value as SupportCategory)}
              className="mt-1.5 w-full rounded-xl border border-input bg-background px-3 py-2 text-sm outline-none focus-visible:ring-2 focus-visible:ring-ring"
            >
              {CATEGORIES.map((cat) => (
                <option key={cat.id} value={cat.id}>
                  {cat.label}
                </option>
              ))}
            </select>
          </div>

          <div>
            <label className="block font-medium text-foreground text-xs">Message</label>
            <textarea
              required
              rows={4}
              value={message}
              onChange={(e) => setMessage(e.target.value)}
              placeholder="Describe the issue, what you expected, and any error messages..."
              className="mt-1.5 w-full resize-none rounded-xl border border-input bg-background p-3 text-sm outline-none focus-visible:ring-2 focus-visible:ring-ring"
            />
          </div>

          <label className="flex cursor-pointer items-center gap-2 rounded-xl bg-muted/40 p-3 text-xs">
            <input
              type="checkbox"
              checked={attachDiagnostics}
              onChange={(e) => setAttachDiagnostics(e.target.checked)}
              className="size-4 rounded text-primary focus:ring-primary"
            />
            <Paperclip className="size-3.5 text-muted-foreground" />
            <span className="text-muted-foreground">
              Attach server runtime diagnostics and deployment error logs
            </span>
          </label>

          <div className="flex items-center justify-end gap-2 pt-2">
            <Button variant="outline" size="sm" type="button" onClick={onClose}>
              Cancel
            </Button>
            <Button
              size="sm"
              type="submit"
              disabled={!subject.trim() || !message.trim() || submitting}
              className="gap-1.5"
            >
              <Send className="size-3.5" />
              <span>{submitting ? 'Submitting...' : 'Submit ticket'}</span>
            </Button>
          </div>
        </form>
      </div>
    </div>
  );
}
