import { ProjectSecurityCard } from '#/components/security/project-security-card';

export function ProjectSecuritySettings({ projectId }: { projectId: string }) {
  return (
    <section className="rounded-2xl border border-border/50 bg-card p-6 shadow-xs">
      <div className="mb-4">
        <h2 className="font-semibold text-lg">Runtime Security & Advisories</h2>
        <p className="text-muted-foreground text-sm">
          Cryptographically signed upstream advisories audited against running container manifests
          and images.
        </p>
      </div>
      <ProjectSecurityCard projectId={projectId} />
    </section>
  );
}
