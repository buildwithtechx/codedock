import { FileCode, Plus, ScrollText, Trash2 } from 'lucide-react';
import { useState, useEffect } from 'react';
import { toast } from 'sonner';
import { Button } from '@codedock/ui/button';
import { Input } from '@codedock/ui/input';
import { Label } from '@codedock/ui/label';
import { Textarea } from '@codedock/ui/textarea';
import { SettingsSection } from '#/features/settings/settings-section';

export interface CloudInitScript {
  id: string;
  name: string;
  description?: string;
  script: string;
  updatedAt: string;
}

const STORAGE_KEY = 'codedock_cloud_init_scripts_store';

const DEFAULT_SCRIPT_TEMPLATE = `#cloud-config
package_update: true
package_upgrade: true
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
  - ufw --force enable
`;

const DEFAULT_SCRIPTS: CloudInitScript[] = [
  {
    id: 'ci-default-1',
    name: 'Default Docker & Firewall Bootstrap',
    description:
      'Installs Docker CE, sets up non-root daemon, and locks down UFW ports (22, 80, 443)',
    script: DEFAULT_SCRIPT_TEMPLATE,
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

  const [showForm, setShowForm] = useState(false);
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [script, setScript] = useState(DEFAULT_SCRIPT_TEMPLATE);

  useEffect(() => {
    try {
      localStorage.setItem(STORAGE_KEY, JSON.stringify(scripts));
    } catch {
      // Storage unavailable
    }
  }, [scripts]);

  const resetForm = () => {
    setName('');
    setDescription('');
    setScript(DEFAULT_SCRIPT_TEMPLATE);
    setShowForm(false);
  };

  const handleCreate = (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim()) {
      toast.error('Script name is required');
      return;
    }
    if (!script.trim()) {
      toast.error('Script contents cannot be empty');
      return;
    }

    const newScript: CloudInitScript = {
      id: `ci-${Date.now()}`,
      name: name.trim(),
      description: description.trim() || undefined,
      script: script.trim(),
      updatedAt: new Date().toISOString(),
    };

    setScripts((prev) => [newScript, ...prev]);
    toast.success(`Cloud-init script "${name}" created`);
    resetForm();
  };

  const handleDelete = (id: string, scriptName: string) => {
    setScripts((prev) => prev.filter((s) => s.id !== id));
    toast.success(`Cloud-init script "${scriptName}" deleted`);
  };

  return (
    <SettingsSection
      icon={<ScrollText className="size-4 text-primary" />}
      title="Cloud-Init Scripts"
      description="Reusable cloud-init and bootstrap configuration scripts applied during automatic server provisioning."
      action={
        !showForm ? (
          <Button size="sm" onClick={() => setShowForm(true)} className="gap-1.5">
            <Plus className="size-3.5" />
            <span>New script</span>
          </Button>
        ) : undefined
      }
    >
      {showForm && (
        <form
          onSubmit={handleCreate}
          className="mb-5 space-y-4 rounded-xl border border-border/50 bg-muted/20 p-4"
        >
          <p className="font-semibold text-xs text-foreground border-b border-border/40 pb-2">
            Create Cloud-Init Template
          </p>

          <div className="grid gap-3 sm:grid-cols-2">
            <div className="space-y-1">
              <Label htmlFor="ci-name" className="text-xs">
                Script Name
              </Label>
              <Input
                id="ci-name"
                value={name}
                onChange={(e) => setName(e.target.value)}
                placeholder="e.g. Production Hardening"
                className="bg-muted/40 text-xs h-9"
                required
              />
            </div>
            <div className="space-y-1">
              <Label htmlFor="ci-desc" className="text-xs">
                Description (Optional)
              </Label>
              <Input
                id="ci-desc"
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                placeholder="Brief summary of installed services"
                className="bg-muted/40 text-xs h-9"
              />
            </div>
          </div>

          <div className="space-y-1">
            <Label htmlFor="ci-script" className="text-xs">
              YAML / Shell Script Payload
            </Label>
            <Textarea
              id="ci-script"
              value={script}
              onChange={(e) => setScript(e.target.value)}
              rows={8}
              className="font-mono text-xs bg-muted/40"
              placeholder="#cloud-config..."
              required
            />
          </div>

          <div className="flex justify-end gap-2 pt-2 border-t border-border/40">
            <Button type="button" variant="outline" size="sm" onClick={resetForm}>
              Cancel
            </Button>
            <Button type="submit" size="sm">
              Save Script
            </Button>
          </div>
        </form>
      )}

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
    </SettingsSection>
  );
}
