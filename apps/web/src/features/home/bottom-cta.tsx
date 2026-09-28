import { Cloud } from 'lucide-react';

export function BottomCta() {
  return (
    <section className="relative overflow-hidden py-24 md:py-32">
      <div className="pointer-events-none absolute inset-0 z-0 bg-[radial-gradient(ellipse_at_center,var(--color-primary)_0%,transparent_60%)] opacity-15" />

      <div className="relative z-10 mx-auto max-w-4xl px-6 text-center">
        <h2 className="mb-6 text-balance font-extrabold text-4xl text-foreground tracking-tight md:text-5xl">
          Ready to own your infrastructure?
        </h2>
        <p className="mx-auto mb-10 max-w-2xl text-lg text-muted-foreground">
          Stop renting Heroku dynos. Take back control with an open-source daemon and a powerful
          fleet control plane. Deploy your first app in 60 seconds.
        </p>

        <div className="flex flex-col items-center justify-center gap-4 sm:flex-row">
          <div className="flex w-full items-center gap-2 rounded-lg border border-border bg-background/80 px-5 py-3.5 font-mono text-sm shadow-xl backdrop-blur-md sm:w-auto">
            <span className="shrink-0 font-bold text-primary">$</span>
            <code className="text-foreground">curl -fsSL https://get.codedock.run | sh</code>
          </div>
          <a
            href="https://app.codedock.run"
            className="flex w-full items-center justify-center gap-2 rounded-lg bg-primary px-8 py-3.5 font-semibold text-white shadow-lg shadow-primary/25 transition-colors hover:bg-primary-hover sm:w-auto"
          >
            <Cloud className="size-5" />
            Get Cloud
          </a>
        </div>
      </div>
    </section>
  );
}
