import { CheckCircle2, CloudUpload, Database, GitBranch } from 'lucide-react';

export function DeepDiveFeatures() {
  return (
    <section className="relative overflow-hidden border-border border-t py-24">
      <div className="mx-auto max-w-7xl space-y-32 px-6">
        <div className="flex flex-col items-center gap-12 lg:flex-row lg:gap-20">
          <div className="flex-1 space-y-6">
            <div className="inline-flex items-center gap-2 rounded-md bg-primary/10 px-3 py-1 font-bold text-primary text-xs uppercase tracking-wider">
              <GitBranch className="size-4" /> Git-Native
            </div>
            <h2 className="font-extrabold text-3xl text-foreground leading-tight md:text-4xl">
              Push to deploy. No CI pipelines required.
            </h2>
            <p className="text-lg text-muted-foreground leading-relaxed">
              Connect your GitHub, GitLab, or Gitea repositories. Codedock detects your stack (Node,
              Go, Python, etc.), builds the image locally via Railpack, and deploys it. Zero YAML
              required.
            </p>
            <ul className="space-y-3 pt-4">
              {[
                'Automatic PR preview environments',
                'Instant rollbacks to previous commits',
                'Build logs streamed live to the dashboard',
              ].map((t) => (
                <li key={t} className="flex items-start gap-3 font-medium text-foreground text-sm">
                  <CheckCircle2 className="size-5 shrink-0 text-primary" />
                  {t}
                </li>
              ))}
            </ul>
          </div>
          <div className="w-full flex-1">
            <div className="relative overflow-hidden rounded-xl border border-border bg-background/80 p-4 shadow-2xl backdrop-blur-md">
              <div className="absolute inset-x-0 top-0 h-1 bg-gradient-to-r from-primary to-blue-500" />
              <div className="mb-4 flex items-center gap-2 px-2 pt-2 font-mono text-muted-foreground text-xs">
                <span>Terminal</span>
              </div>
              <div className="space-y-2 rounded-lg border border-border bg-background p-4 font-mono text-muted-foreground text-xs">
                <p>
                  <span className="text-green-400">codedock</span>{' '}
                  <span className="text-foreground">build:</span> Detected Node.js application
                </p>
                <p>
                  <span className="text-blue-400">docker</span>{' '}
                  <span className="text-foreground">layer:</span> Resolving package.json
                  dependencies...
                </p>
                <p>
                  <span className="text-blue-400">docker</span>{' '}
                  <span className="text-foreground">layer:</span> Running npm install
                </p>
                <p>
                  <span className="text-green-400">codedock</span>{' '}
                  <span className="text-foreground">deploy:</span> Starting container app_web_1
                </p>
                <p>
                  <span className="text-primary">traefik</span>{' '}
                  <span className="text-foreground">router:</span> Requesting Let's Encrypt SSL
                  certificate
                </p>
                <p className="pt-2 font-bold text-green-400">
                  ✓ Deployed to https://app.acme.com in 14s
                </p>
              </div>
            </div>
          </div>
        </div>

        <div className="flex flex-col items-center gap-12 lg:flex-row-reverse lg:gap-20">
          <div className="flex-1 space-y-6">
            <div className="inline-flex items-center gap-2 rounded-md bg-primary/10 px-3 py-1 font-bold text-primary text-xs uppercase tracking-wider">
              <Database className="size-4" /> Databases & Backups
            </div>
            <h2 className="font-extrabold text-3xl text-foreground leading-tight md:text-4xl">
              Stateful services, handled.
            </h2>
            <p className="text-lg text-muted-foreground leading-relaxed">
              Deploying Postgres or Redis is easy. Keeping them safe is hard. Codedock provisions
              your databases and automatically streams encrypted backups to your own S3 bucket.
            </p>
            <ul className="space-y-3 pt-4">
              {[
                'Point-in-time automated snapshots',
                'Zero-downtime volume mounts',
                'Bring your own AWS/R2 bucket',
              ].map((t) => (
                <li key={t} className="flex items-start gap-3 font-medium text-foreground text-sm">
                  <CheckCircle2 className="size-5 shrink-0 text-primary" />
                  {t}
                </li>
              ))}
            </ul>
          </div>
          <div className="w-full flex-1">
            <div className="relative overflow-hidden rounded-xl border border-border bg-background/80 p-6 shadow-2xl backdrop-blur-md">
              <div className="absolute inset-y-0 left-0 w-1 bg-gradient-to-b from-primary to-blue-500" />
              <div className="space-y-4">
                <div className="flex items-center justify-between border-border border-b pb-4">
                  <div className="flex items-center gap-3">
                    <div className="flex size-10 items-center justify-center rounded bg-primary/10 text-primary">
                      <Database className="size-5" />
                    </div>
                    <div>
                      <div className="font-bold text-foreground text-sm">prod-postgres</div>
                      <div className="text-green-400 text-xs">Healthy • 1.2 TB</div>
                    </div>
                  </div>
                  <div className="rounded border border-border bg-surface px-3 py-1 text-muted-foreground text-xs">
                    PostgreSQL 16
                  </div>
                </div>

                <div className="space-y-2">
                  <div className="mb-3 font-semibold text-muted-foreground text-xs uppercase tracking-widest">
                    Recent Backups (S3)
                  </div>
                  {[
                    { time: 'Today, 02:00 AM', size: '1.2 TB', status: 'Success' },
                    { time: 'Yesterday, 02:00 AM', size: '1.18 TB', status: 'Success' },
                    { time: 'Jul 18, 02:00 AM', size: '1.15 TB', status: 'Success' },
                  ].map((b) => (
                    <div
                      key={b.time}
                      className="flex items-center justify-between rounded border border-border bg-surface p-3"
                    >
                      <div className="flex items-center gap-2">
                        <CloudUpload className="size-4 text-muted-foreground" />
                        <span className="font-medium text-foreground text-xs">{b.time}</span>
                      </div>
                      <div className="flex items-center gap-4">
                        <span className="text-muted-foreground text-xs">{b.size}</span>
                        <span className="rounded bg-green-400/10 px-2 py-0.5 font-bold text-[10px] text-green-400">
                          {b.status}
                        </span>
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
