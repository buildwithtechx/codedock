import { FileCode, Plus, ScrollText, Trash2 } from 'lucide-react';
import { useState, useEffect } from 'react';
import { toast } from 'sonner';
import { Button } from '@codedock/ui/button';
import { SettingsSection } from '#/features/settings/settings-section';
import { NewCloudInitDialog, type CloudInitScript } from './new-cloud-init-dialog';

const STORAGE_KEY = 'codedock_cloud_init_scripts_store';

const DEFAULT_SCRIPTS: CloudInitScript[] = [
  {
    id: 'ci-default-1',
    name: 'Default Docker & Firewall Bootstrap',
    description:
      'Installs Docker CE, sets up non-root daemon, and locks down UFW ports (22, 80, 443)',
    script: `#cloud-config
package_update: true
packages:
  - curl
  - git
  - ufw

runcmd:
  - ufw default deny incoming
  - ufw default allow outgoing
  - ufw allow 22/tcp
  - ufw allow 80/tcp
  - ufw allow 443/tcp
  - ufw --force enable`,
    updatedAt: new Date().toISOString(),
  },
];

export function CloudInitTab() {
  const [scripts, setScripts] = useState<CloudInitScript[]>(() => {
    if (typeof window === 'undefined') return DEFAULT_SCRIPTS;
    try {
      const stored = localStorage.getItem(STORAGE_KEY);
      return stored ? JSON.parse(stored) : DEFAULT_SCRIPTS;
    } catch {
      return DEFAULT_SCRIPTS;
    }
  });

  const [dialogOpen, setDialogOpen] = useState(false);

  useEffect(() => {
    try {
      localStorage.setItem(STORAGE_KEY, JSON.stringify(scripts));
    } catch {
      // Storage unavailable
    }
  }, [scripts]);

  const handleAdd = (newScript: CloudInitScript) => {
    setScripts((prev) => [newScript, ...prev]);
  };

  const handleDelete = (id: string, name: string) => {
    setScripts((prev) => prev.filter((s) => s.id !== id));
    toast.success(`Cloud-init script "${name}" deleted`);
  };

  return (
    <SettingsSection
      icon={<ScrollText className="size-4 text-primary" />}
      title="Cloud-Init Scripts"
      description="Reusable cloud-init and bootstrap configuration scripts applied during automatic server provisioning."
      action={
        <Button size="sm" onClick={() => setDialogOpen(true)} className="gap-1.5">
          <Plus className="size-3.5" />
          New script
        </Button>
      }
    >
      {scripts.length === 0 ? (
        <div className="rounded-xl border border-border/60 border-dashed py-10 text-center">
          <ScrollText className="mx-auto size-8 text-muted-foreground/60" />
          <h3 className="mt-3 font-semibold text-sm">No cloud-init scripts</h3>
          <p className="mt-1 text-xs text-muted-foreground">
            Create initialization scripts to automatically configure servers on boot.
          </p>
        </div>
      ) : (
        <div className="divide-y divide-border/50 rounded-xl border border-border/50">
          {scripts.map((item) => (
            <div key={item.id} className="p-4 space-y-3">
              <div className="flex items-center justify-between">
                <div className="space-y-0.5">
                  <div className="flex items-center gap-2">
                    <FileCode className="size-4 text-primary" />
                    <span className="font-semibold text-foreground text-sm">{item.name}</span>
                  </div>
                  {item.description && (
                    <p className="text-xs text-muted-foreground">{item.description}</p>
                  )}
                </div>
                <Button
                  variant="ghost"
                  size="icon"
                  onClick={() => handleDelete(item.id, item.name)}
                  className="size-8 text-muted-foreground hover:text-destructive shrink-0"
                  aria-label={`Delete ${item.name}`}
                >
                  <Trash2 className="size-4" />
                </Button>
              </div>

              <div className="rounded-lg bg-muted/60 p-3">
                <pre className="max-h-32 overflow-y-auto font-mono text-[11px] text-muted-foreground leading-relaxed whitespace-pre-wrap">
                  {item.script}
                </pre>
              </div>
            </div>
          ))}
        </div>
      )}

      <NewCloudInitDialog open={dialogOpen} onOpenChange={setDialogOpen} onAddScript={handleAdd} />
    </SettingsSection>
  );
}
