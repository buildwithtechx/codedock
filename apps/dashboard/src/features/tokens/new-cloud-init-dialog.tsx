import { ScrollText } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Button } from '@codedock/ui/button';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@codedock/ui/dialog';
import { Input } from '@codedock/ui/input';
import { Label } from '@codedock/ui/label';
import { Textarea } from '@codedock/ui/textarea';

export interface CloudInitScript {
  id: string;
  name: string;
  description?: string;
  script: string;
  updatedAt: string;
}

interface NewCloudInitDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onAddScript: (script: CloudInitScript) => void;
}

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

export function NewCloudInitDialog({ open, onOpenChange, onAddScript }: NewCloudInitDialogProps) {
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [script, setScript] = useState(DEFAULT_SCRIPT_TEMPLATE);

  const handleSubmit = (e: React.FormEvent) => {
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

    onAddScript(newScript);
    toast.success(`Cloud-init script "${name}" created`);
    setName('');
    setDescription('');
    setScript(DEFAULT_SCRIPT_TEMPLATE);
    onOpenChange(false);
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-xl">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <ScrollText className="size-4 text-primary" />
            New Cloud-Init Script
          </DialogTitle>
          <DialogDescription>
            Reusable initialization scripts executed during server provisioning (Docker, firewalls,
            packages).
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit} className="space-y-4 py-2">
          <div className="grid gap-3 sm:grid-cols-2">
            <div className="space-y-1.5">
              <Label htmlFor="ci-name">Name</Label>
              <Input
                id="ci-name"
                value={name}
                onChange={(e) => setName(e.target.value)}
                placeholder="e.g. Standard Ubuntu Harden"
                required
              />
            </div>
            <div className="space-y-1.5">
              <Label htmlFor="ci-desc">Description (optional)</Label>
              <Input
                id="ci-desc"
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                placeholder="e.g. Installs UFW and sets up standard firewall"
              />
            </div>
          </div>

          <div className="space-y-1.5">
            <Label htmlFor="ci-script">Script (cloud-config YAML or bash)</Label>
            <Textarea
              id="ci-script"
              rows={8}
              value={script}
              onChange={(e) => setScript(e.target.value)}
              className="font-mono text-xs"
              required
            />
          </div>

          <DialogFooter className="pt-2">
            <Button type="button" variant="outline" size="sm" onClick={() => onOpenChange(false)}>
              Cancel
            </Button>
            <Button type="submit" size="sm">
              Save Script
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
