import { CheckCircle2, CloudUpload, Database, GitBranch } from "lucide-react";

export function DeepDiveFeatures() {
  return (
    <section className="py-24 overflow-hidden relative border-t border-border">
      <div className="max-w-7xl mx-auto px-6 space-y-32">
        <div className="flex flex-col lg:flex-row items-center gap-12 lg:gap-20">
          <div className="flex-1 space-y-6">
            <div className="inline-flex items-center gap-2 px-3 py-1 rounded-md bg-primary/10 text-primary text-xs font-bold uppercase tracking-wider">
              <GitBranch className="size-4" /> Git-Native
            </div>
            <h2 className="text-3xl md:text-4xl font-extrabold text-foreground leading-tight">
              Push to deploy. No CI pipelines required.
            </h2>
            <p className="text-muted-foreground text-lg leading-relaxed">
              Connect your GitHub, GitLab, or Gitea repositories. Codedock detects your stack (Node,
              Go, Python, etc.), builds the image locally via Railpack, and deploys it. Zero YAML
              required.
            </p>
            <ul className="space-y-3 pt-4">
              {[
                "Automatic PR preview environments",
                "Instant rollbacks to previous commits",
                "Build logs streamed live to the dashboard",
              ].map((t) => (
                <li key={t} className="flex items-start gap-3 text-sm text-foreground font-medium">
                  <CheckCircle2 className="size-5 text-primary shrink-0" />
                  {t}
                </li>
              ))}
            </ul>
          </div>
          <div className="flex-1 w-full">
            <div className="bg-background/80 backdrop-blur-md border border-border rounded-xl p-4 shadow-2xl relative overflow-hidden">
              <div className="absolute inset-x-0 top-0 h-1 bg-gradient-to-r from-primary to-blue-500" />
              <div className="flex items-center gap-2 mb-4 px-2 pt-2 text-xs text-muted-foreground font-mono">
                <span>Terminal</span>
              </div>
              <div className="bg-background rounded-lg p-4 font-mono text-xs text-muted-foreground space-y-2 border border-border">
                <p>
                  <span className="text-green-400">codedock</span>{" "}
                  <span className="text-foreground">build:</span> Detected Node.js application
                </p>
                <p>
                  <span className="text-blue-400">docker</span>{" "}
                  <span className="text-foreground">layer:</span> Resolving package.json
                  dependencies...
                </p>
                <p>
                  <span className="text-blue-400">docker</span>{" "}
                  <span className="text-foreground">layer:</span> Running npm install
                </p>
                <p>
                  <span className="text-green-400">codedock</span>{" "}
                  <span className="text-foreground">deploy:</span> Starting container app_web_1
                </p>
                <p>
                  <span className="text-primary">traefik</span>{" "}
                  <span className="text-foreground">router:</span> Requesting Let's Encrypt SSL
                  certificate
                </p>
                <p className="text-green-400 font-bold pt-2">
                  ✓ Deployed to https://app.acme.com in 14s
                </p>
              </div>
            </div>
          </div>
        </div>

        <div className="flex flex-col lg:flex-row-reverse items-center gap-12 lg:gap-20">
          <div className="flex-1 space-y-6">
            <div className="inline-flex items-center gap-2 px-3 py-1 rounded-md bg-primary/10 text-primary text-xs font-bold uppercase tracking-wider">
              <Database className="size-4" /> Databases & Backups
            </div>
            <h2 className="text-3xl md:text-4xl font-extrabold text-foreground leading-tight">
              Stateful services, handled.
            </h2>
            <p className="text-muted-foreground text-lg leading-relaxed">
              Deploying Postgres or Redis is easy. Keeping them safe is hard. Codedock provisions
              your databases and automatically streams encrypted backups to your own S3 bucket.
            </p>
            <ul className="space-y-3 pt-4">
              {[
                "Point-in-time automated snapshots",
                "Zero-downtime volume mounts",
                "Bring your own AWS/R2 bucket",
              ].map((t) => (
                <li key={t} className="flex items-start gap-3 text-sm text-foreground font-medium">
                  <CheckCircle2 className="size-5 text-primary shrink-0" />
                  {t}
                </li>
              ))}
            </ul>
          </div>
          <div className="flex-1 w-full">
            <div className="bg-background/80 backdrop-blur-md border border-border rounded-xl p-6 shadow-2xl relative overflow-hidden">
              <div className="absolute inset-y-0 left-0 w-1 bg-gradient-to-b from-primary to-blue-500" />
              <div className="space-y-4">
                <div className="flex items-center justify-between border-b border-border pb-4">
                  <div className="flex items-center gap-3">
                    <div className="size-10 rounded bg-primary/10 flex items-center justify-center text-primary">
                      <Database className="size-5" />
                    </div>
                    <div>
                      <div className="text-sm font-bold text-foreground">prod-postgres</div>
                      <div className="text-xs text-green-400">Healthy • 1.2 TB</div>
                    </div>
                  </div>
                  <div className="px-3 py-1 rounded bg-surface border border-border text-xs text-muted-foreground">
                    PostgreSQL 16
                  </div>
                </div>

                <div className="space-y-2">
                  <div className="text-xs font-semibold text-muted-foreground uppercase tracking-widest mb-3">
                    Recent Backups (S3)
                  </div>
                  {[
                    { time: "Today, 02:00 AM", size: "1.2 TB", status: "Success" },
                    { time: "Yesterday, 02:00 AM", size: "1.18 TB", status: "Success" },
                    { time: "Jul 18, 02:00 AM", size: "1.15 TB", status: "Success" },
                  ].map((b) => (
                    <div
                      key={b.time}
                      className="flex items-center justify-between bg-surface rounded p-3 border border-border"
                    >
                      <div className="flex items-center gap-2">
                        <CloudUpload className="size-4 text-muted-foreground" />
                        <span className="text-xs text-foreground font-medium">{b.time}</span>
                      </div>
                      <div className="flex items-center gap-4">
                        <span className="text-xs text-muted-foreground">{b.size}</span>
                        <span className="text-[10px] text-green-400 font-bold bg-green-400/10 px-2 py-0.5 rounded">
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
