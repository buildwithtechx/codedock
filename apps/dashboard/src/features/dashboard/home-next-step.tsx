import { Link } from '@tanstack/react-router';
import { ArrowRight, Globe, HardDriveDownload, Sparkles } from 'lucide-react';
import { Button } from '@codedock/ui/button';

export function HomeNextStep({ hasProjects }: { hasProjects: boolean }) {
  if (!hasProjects) {
    return (
      <section className="relative overflow-hidden rounded-2xl border border-primary/20 bg-linear-to-br from-primary/5 via-card to-card p-5 shadow-xs">
        <div className="flex items-center gap-2.5">
          <span className="flex size-8 items-center justify-center rounded-xl bg-primary/10 text-primary">
            <Sparkles className="size-4" />
          </span>
          <div>
            <p className="font-semibold text-[10px] text-primary uppercase tracking-wider">
              Quick Start
            </p>
            <h2 className="font-semibold text-foreground text-sm">Deploy Your First Project</h2>
          </div>
        </div>
        <p className="mt-3 text-muted-foreground text-xs leading-relaxed">
          Launch from GitHub repositories, Dockerfiles, or ready-to-run templates with automated SSL.
        </p>
        <div className="mt-4 border-border/40 border-t pt-3">
          <Button asChild size="sm" className="h-8 w-full justify-between font-medium text-xs">
            <Link to="/library">
              <span>Explore Templates</span>
              <ArrowRight className="size-3.5 opacity-80" />
            </Link>
          </Button>
        </div>
      </section>
    );
  }

  return (
    <section className="relative overflow-hidden rounded-2xl border border-border/50 bg-card p-5 shadow-xs">
      <div className="flex items-center gap-2.5">
        <span className="flex size-8 items-center justify-center rounded-xl bg-primary/10 text-primary">
          <Sparkles className="size-4" />
        </span>
        <div>
          <p className="font-semibold text-[10px] text-primary uppercase tracking-wider">
            Recommendations
          </p>
          <h2 className="font-semibold text-foreground text-sm">Harden Your Stack</h2>
        </div>
      </div>
      <p className="mt-3 text-muted-foreground text-xs leading-relaxed">
        Attach custom ingress domains with wildcard TLS and schedule automated volume backups.
      </p>

      <div className="mt-4 grid grid-cols-2 gap-2">
        <Button
          asChild
          variant="outline"
          size="sm"
          className="h-8 w-full justify-start text-[11px] font-medium"
        >
          <Link to="/domains">
            <Globe className="mr-1.5 size-3 text-muted-foreground" />
            <span>Add Domain</span>
          </Link>
        </Button>
        <Button
          asChild
          variant="outline"
          size="sm"
          className="h-8 w-full justify-start text-[11px] font-medium"
        >
          <Link to="/backups">
            <HardDriveDownload className="mr-1.5 size-3 text-muted-foreground" />
            <span>Backups</span>
          </Link>
        </Button>
      </div>
    </section>
  );
}
