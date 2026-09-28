import { Star } from 'lucide-react';

export function SocialProof() {
  const testimonials = [
    {
      text: `"Migrated 15 Heroku apps to a $20 Hetzner instance using Codedock. The zero-config Traefik routing is actual magic."`,
      author: 'JD',
      role: 'Full-stack Developer',
    },
    {
      text: `"Finally, a PaaS alternative that doesn't feel like maintaining Kubernetes. The Yamux tunnel setup is incredibly secure."`,
      author: 'AS',
      role: 'DevOps Engineer',
    },
    {
      text: `"The Canvas view makes onboarding junior devs so easy. They can see exactly how the Redis cache connects to the API."`,
      author: 'MR',
      role: 'CTO @ Startup',
    },
  ];

  return (
    <section className="border-border border-b py-20">
      <div className="mx-auto max-w-7xl px-6 text-center">
        <p className="mb-8 font-semibold text-primary text-sm uppercase tracking-widest">
          Trusted by developers moving off PaaS
        </p>

        <div className="mx-auto grid max-w-5xl grid-cols-1 gap-6 text-left md:grid-cols-3">
          {testimonials.map((t) => (
            <div
              key={t.author}
              className="rounded-xl border border-border bg-background/80 p-6 shadow-sm backdrop-blur-md"
            >
              <div className="mb-4 flex items-center gap-2">
                <div className="flex text-yellow-400">
                  <Star className="size-4 fill-current" />
                  <Star className="size-4 fill-current" />
                  <Star className="size-4 fill-current" />
                  <Star className="size-4 fill-current" />
                  <Star className="size-4 fill-current" />
                </div>
              </div>
              <p className="mb-6 text-muted-foreground text-sm">{t.text}</p>
              <div className="flex items-center gap-3">
                <div className="flex size-8 items-center justify-center rounded-full border border-border bg-surface font-bold text-foreground text-xs">
                  {t.author}
                </div>
                <div className="text-muted-foreground text-xs">{t.role}</div>
              </div>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
