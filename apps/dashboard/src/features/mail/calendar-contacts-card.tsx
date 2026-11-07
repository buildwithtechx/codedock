import { Button } from '@codedock/ui/button';
import { Link } from '@tanstack/react-router';
import { ArrowRight, Calendar } from 'lucide-react';

export function CalendarContactsCard({ serverDomain }: { serverDomain?: string }) {
  return (
    <div className="space-y-3 rounded-2xl border border-border/60 bg-card p-5 shadow-xs">
      <div className="flex items-center gap-2">
        <div className="flex size-8 items-center justify-center rounded-lg bg-primary/10 text-primary">
          <Calendar className="size-4" />
        </div>
        <div>
          <h3 className="font-semibold text-foreground text-sm">
            Calendars & Contacts (CalDAV/CardDAV)
          </h3>
          <p className="text-muted-foreground text-xs">
            Sync native phone and desktop calendar and contacts apps for{' '}
            {serverDomain || 'your mailboxes'}.
          </p>
        </div>
      </div>

      <p className="text-muted-foreground text-xs leading-relaxed">
        Radicale runs alongside your mail server and delegates authentication directly to your IMAP
        mailboxes over TLS. No separate passwords or user sync required.
      </p>

      <div className="pt-1">
        <Button size="sm" variant="outline" asChild className="h-8 gap-1.5 text-xs">
          <Link to="/apps/new/$appId" params={{ appId: 'radicale' }}>
            Deploy Radicale companion
            <ArrowRight className="size-3.5" />
          </Link>
        </Button>
      </div>
    </div>
  );
}
