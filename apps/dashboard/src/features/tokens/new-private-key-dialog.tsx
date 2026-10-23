import { Key } from 'lucide-react';
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

export interface SSHPrivateKey {
  id: string;
  name: string;
  description?: string;
  publicKey: string;
  fingerprint: string;
  createdAt: string;
}

interface NewPrivateKeyDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onAddKey: (key: SSHPrivateKey) => void;
}

export function NewPrivateKeyDialog({ open, onOpenChange, onAddKey }: NewPrivateKeyDialogProps) {
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [privateKey, setPrivateKey] = useState('');
  const [publicKey, setPublicKey] = useState('');

  const handleGenerate = () => {
    // Generate a mock realistic ED25519 fingerprint and public key for self-hosting setups
    const mockHash = Array.from({ length: 8 }, () =>
      Math.floor(Math.random() * 256)
        .toString(16)
        .padStart(2, '0')
    ).join(':');
    const mockPub = `ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAI${btoa(mockHash).replace(/=/g, '')} codedock@instance`;
    setPublicKey(mockPub);
    setPrivateKey(
      '-----BEGIN OPENSSH PRIVATE KEY-----\nb3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAMwAAAAtzc2gt\nZWQyNTUxOQAAAC...[AUTOMATICALLY GENERATED DEPLOY KEY]...\n-----END OPENSSH PRIVATE KEY-----'
    );
    toast.success('Generated new key pair');
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim()) {
      toast.error('Key name is required');
      return;
    }
    if (!privateKey.trim()) {
      toast.error('Private key content is required');
      return;
    }

    const fingerprint = publicKey.trim()
      ? `SHA256:${btoa(name).slice(0, 16)}...`
      : 'SHA256:imported_key';

    const derivedPublic =
      publicKey.trim() || `ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAA... ${name}@codedock`;

    const newKey: SSHPrivateKey = {
      id: `key-${Date.now()}`,
      name: name.trim(),
      description: description.trim() || undefined,
      publicKey: derivedPublic,
      fingerprint,
      createdAt: new Date().toISOString(),
    };

    onAddKey(newKey);
    toast.success(`Private key "${name}" added`);
    setName('');
    setDescription('');
    setPrivateKey('');
    setPublicKey('');
    onOpenChange(false);
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-xl">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <Key className="size-4 text-primary" />
            Add Private Key
          </DialogTitle>
          <DialogDescription>
            Add an SSH private key to connect to destination servers or deploy private git
            repositories.
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit} className="space-y-4 py-2">
          <div className="grid gap-3 sm:grid-cols-2">
            <div className="space-y-1.5">
              <Label htmlFor="key-name">Name</Label>
              <Input
                id="key-name"
                value={name}
                onChange={(e) => setName(e.target.value)}
                placeholder="e.g. Production Hetzner Key"
                required
              />
            </div>
            <div className="space-y-1.5">
              <Label htmlFor="key-desc">Description (optional)</Label>
              <Input
                id="key-desc"
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                placeholder="e.g. Used for root ssh access"
              />
            </div>
          </div>

          <div className="space-y-1.5">
            <div className="flex items-center justify-between">
              <Label htmlFor="priv-key">Private Key (OpenSSH format)</Label>
              <button
                type="button"
                onClick={handleGenerate}
                className="text-xs text-primary hover:underline font-medium"
              >
                Auto-generate key pair
              </button>
            </div>
            <Textarea
              id="priv-key"
              rows={5}
              value={privateKey}
              onChange={(e) => setPrivateKey(e.target.value)}
              placeholder="-----BEGIN OPENSSH PRIVATE KEY-----&#10;...&#10;-----END OPENSSH PRIVATE KEY-----"
              className="font-mono text-xs"
              required
            />
          </div>

          <div className="space-y-1.5">
            <Label htmlFor="pub-key">Public Key (optional / derived)</Label>
            <Input
              id="pub-key"
              value={publicKey}
              onChange={(e) => setPublicKey(e.target.value)}
              placeholder="ssh-ed25519 AAAAC3NzaC1... or ssh-rsa AAAAB3..."
              className="font-mono text-xs"
            />
            <p className="text-[11px] text-muted-foreground">
              Add this public key to ~/.ssh/authorized_keys on your destination servers.
            </p>
          </div>

          <DialogFooter className="pt-2">
            <Button type="button" variant="outline" size="sm" onClick={() => onOpenChange(false)}>
              Cancel
            </Button>
            <Button type="submit" size="sm">
              Save Private Key
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
