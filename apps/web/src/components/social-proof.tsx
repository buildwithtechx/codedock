import { Star } from "lucide-react";

export function SocialProof() {
  const testimonials = [
    {
      text: `"Migrated 15 Heroku apps to a $20 Hetzner instance using Codedock. The zero-config Traefik routing is actual magic."`,
      author: "JD",
      role: "Full-stack Developer",
    },
    {
      text: `"Finally, a PaaS alternative that doesn't feel like maintaining Kubernetes. The Yamux tunnel setup is incredibly secure."`,
      author: "AS",
      role: "DevOps Engineer",
    },
    {
      text: `"The Canvas view makes onboarding junior devs so easy. They can see exactly how the Redis cache connects to the API."`,
      author: "MR",
      role: "CTO @ Startup",
    },
  ];

  return (
    <section className="py-20 border-b border-border">
      <div className="max-w-7xl mx-auto px-6 text-center">
        <p className="text-sm font-semibold text-primary uppercase tracking-widest mb-8">
          Trusted by developers moving off PaaS
        </p>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 max-w-5xl mx-auto text-left">
          {testimonials.map((t) => (
            <div
              key={t.author}
              className="bg-background/80 backdrop-blur-md border border-border p-6 rounded-xl shadow-sm"
            >
              <div className="flex items-center gap-2 mb-4">
                <div className="flex text-yellow-400">
                  <Star className="size-4 fill-current" />
                  <Star className="size-4 fill-current" />
                  <Star className="size-4 fill-current" />
                  <Star className="size-4 fill-current" />
                  <Star className="size-4 fill-current" />
                </div>
              </div>
              <p className="text-sm text-muted-foreground mb-6">{t.text}</p>
              <div className="flex items-center gap-3">
                <div className="size-8 rounded-full bg-surface border border-border flex items-center justify-center font-bold text-xs text-foreground">
                  {t.author}
                </div>
                <div className="text-xs text-muted-foreground">{t.role}</div>
              </div>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
