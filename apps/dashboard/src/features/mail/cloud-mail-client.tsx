import { Link } from '@tanstack/react-router';
import { LayoutGrid, Sparkles } from 'lucide-react';
import { Button } from '@codedock/ui/button';

interface CloudMailClientProps {
  onSelfHost?: () => void;
}

export function CloudMailClient({ onSelfHost }: CloudMailClientProps = {}) {
  return (
    <div className="space-y-6">
      <header>
        <h1 className="font-medium text-2xl text-foreground tracking-tight">Emails</h1>
        <p className="mt-1 text-muted-foreground text-sm">
          Use your existing mailboxes with Codedock Webmail.
        </p>
      </header>

      <div className="grid items-start gap-6 lg:grid-cols-[minmax(0,1fr)_340px]">
        <section className="flex flex-col items-center rounded-2xl bg-card px-5 py-10 text-center sm:px-8 sm:py-12">
          <svg
            aria-hidden="true"
            viewBox="0 0 240 130"
            className="mb-7 w-60 max-w-full text-muted-foreground"
            fill="none"
          >
            <path
              d="M49 65H191"
              stroke="currentColor"
              strokeOpacity=".25"
              strokeWidth="1.5"
              strokeDasharray="4 5"
            />
            <rect
              x="14"
              y="36"
              width="52"
              height="58"
              rx="12"
              fill="currentColor"
              fillOpacity=".05"
            />
            <rect
              x="28"
              y="49"
              width="24"
              height="30"
              rx="4"
              stroke="currentColor"
              strokeOpacity=".7"
              strokeWidth="1.5"
            />
            <path
              d="M34 56H46M34 62H46"
              stroke="currentColor"
              strokeOpacity=".7"
              strokeWidth="1.5"
              strokeLinecap="round"
            />
            <circle cx="40" cy="72" r="1.5" fill="currentColor" fillOpacity=".7" />
            <rect
              x="77"
              y="20"
              width="86"
              height="90"
              rx="20"
              fill="currentColor"
              fillOpacity=".05"
            />
            <rect
              x="95"
              y="46"
              width="50"
              height="37"
              rx="7"
              stroke="currentColor"
              strokeWidth="1.5"
            />
            <path
              d="M97 51L116 65C118.4 66.8 121.6 66.8 124 65L143 51"
              stroke="currentColor"
              strokeWidth="1.5"
              strokeLinecap="round"
            />
            <rect
              x="174"
              y="36"
              width="52"
              height="58"
              rx="12"
              fill="currentColor"
              fillOpacity=".05"
            />
            <rect
              x="185"
              y="50"
              width="30"
              height="25"
              rx="4"
              stroke="currentColor"
              strokeOpacity=".7"
              strokeWidth="1.5"
            />
            <path
              d="M193 82H207M200 75V82"
              stroke="currentColor"
              strokeOpacity=".7"
              strokeWidth="1.5"
              strokeLinecap="round"
            />
          </svg>

          <h2 className="font-medium text-foreground text-xl">Connect your mail server</h2>
          <p className="mt-2 max-w-md text-muted-foreground text-sm leading-relaxed">
            Connect any IMAP/SMTP server. Your provider keeps your mailboxes and mail delivery.
          </p>

          <Button asChild className="mt-6 rounded-full px-5 font-medium">
            <Link to={'/apps/new/mail' as never}>
              <Sparkles className="size-4" />
              <span>Connect existing mail server</span>
            </Link>
          </Button>

          <p className="mt-4 max-w-md text-muted-foreground text-xs leading-relaxed">
            Sign in to webmail with your mailbox email address and password.
          </p>

          <Button
            asChild
            variant="ghost"
            size="sm"
            className="mt-2 text-muted-foreground hover:text-foreground"
          >
            <Link to={'/apps' as never}>
              <LayoutGrid className="size-3.5" />
              <span>Manage webmail apps</span>
            </Link>
          </Button>
        </section>

        <aside className="rounded-2xl bg-card p-5">
          {onSelfHost ? (
            <>
              <h2 className="font-medium text-foreground text-sm">Self-host Codedock Mail</h2>
              <p className="mt-3 text-muted-foreground text-xs leading-relaxed">
                Deploy the complete mail engine with SMTP, IMAP, DKIM and mailbox management on your
                own server.
              </p>
              <Button
                type="button"
                variant="secondary"
                onClick={onSelfHost}
                className="mt-4 w-full rounded-full text-xs"
              >
                Deploy mail server
              </Button>
            </>
          ) : (
            <>
              <div className="flex flex-wrap items-center justify-between gap-3">
                <h2 className="font-medium text-foreground text-sm">Managed email</h2>
                <span className="rounded-full bg-muted/60 px-2.5 py-0.5 text-muted-foreground text-xs">
                  Coming soon
                </span>
              </div>
              <p className="mt-4 text-muted-foreground text-xs leading-relaxed">
                Managed email on Codedock Cloud is coming soon. Cloud server plans are for hosting
                applications and do not provision a mail server.
              </p>
            </>
          )}
        </aside>
      </div>
    </div>
  );
}
