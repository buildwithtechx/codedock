import { ArrowRightLeft, Download, Loader2, Upload } from 'lucide-react';
import { useRef, useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import { ServerTakeoverDialog } from '#/features/servers/components/server-takeover-dialog';
import { useExportSystem, useImportSystem } from '#/features/settings';
import { SettingsSection } from './settings-section';

export function MigrationSettings() {
  const [passphrase, setPassphrase] = useState('');
  const [confirmPassphrase, setConfirmPassphrase] = useState('');
  const [importPassphrase, setImportPassphrase] = useState('');
  const [importFile, setImportFile] = useState<File | null>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const exportSystem = useExportSystem();
  const importSystem = useImportSystem();

  const handleExport = async (e: React.FormEvent) => {
    e.preventDefault();
    if (passphrase !== confirmPassphrase) {
      toast.error('Passphrases do not match');
      return;
    }
    if (passphrase.length < 8) {
      toast.error('Passphrase must be at least 8 characters long');
      return;
    }
    try {
      const blob = await exportSystem.mutateAsync(passphrase);
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `codedock-export-${new Date().toISOString().split('T')[0]}.codedock`;
      a.click();
      window.URL.revokeObjectURL(url);

      toast.success('Instance bundle downloaded successfully');
      setPassphrase('');
      setConfirmPassphrase('');
    } catch {
      toast.error('Failed to export instance bundle');
    }
  };

  const handleImport = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!importFile) {
      toast.error('Please select an export bundle file');
      return;
    }
    if (!importPassphrase) {
      toast.error('Passphrase is required to decrypt the archive');
      return;
    }

    try {
      const formData = new FormData();
      formData.append('bundle', importFile);
      formData.append('passphrase', importPassphrase);
      await importSystem.mutateAsync(formData);
      toast.success('Instance restored successfully. Reloading...');
      setTimeout(() => window.location.reload(), 1500);
    } catch {
      toast.error('Failed to import instance bundle. Verify your passphrase.');
    }
  };

  const passphrasesMismatch = confirmPassphrase.length > 0 && passphrase !== confirmPassphrase;

  return (
    <SettingsSection
      collapsible
      icon={<ArrowRightLeft className="size-4 text-primary" />}
      title="Data Transfer & Migration"
      description="Export encrypted backup bundles or restore and take over an existing deployment."
    >
      <div className="space-y-5">
        <div className="rounded-xl border border-amber-500/20 bg-amber-500/10 p-3.5 text-amber-600 text-xs leading-relaxed dark:text-amber-400">
          Archive files include the Codedock database and credentials. Docker volume contents remain
          on their servers; use project backups for service data.
        </div>

        <div className="rounded-xl border border-border/50 bg-background/50 p-5">
          <form onSubmit={handleExport} className="space-y-4">
            <div>
              <h4 className="font-medium text-foreground text-sm">Export instance bundle</h4>
              <p className="mt-0.5 text-muted-foreground text-xs">
                Package database schemas, environment keys, and project configurations into an
                encrypted archive.
              </p>
            </div>

            <div className="grid grid-cols-1 gap-3 sm:grid-cols-2">
              <div className="space-y-1.5">
                <Label className="text-xs">Migration Passphrase</Label>
                <Input
                  type="password"
                  value={passphrase}
                  onChange={(e) => setPassphrase(e.target.value)}
                  placeholder="Enter a secure passphrase (min 8 chars)"
                  className="bg-muted/30 font-mono text-xs"
                  required
                />
              </div>
              <div className="space-y-1.5">
                <Label className="text-xs">Confirm Passphrase</Label>
                <Input
                  type="password"
                  value={confirmPassphrase}
                  onChange={(e) => setConfirmPassphrase(e.target.value)}
                  placeholder="Re-enter your passphrase"
                  className={`bg-muted/30 font-mono text-xs ${
                    passphrasesMismatch
                      ? 'border-destructive focus-visible:ring-destructive/20'
                      : ''
                  }`}
                  required
                />
              </div>
            </div>

            <div className="flex justify-end border-border/40 border-t pt-3">
              <Button
                type="submit"
                size="sm"
                disabled={exportSystem.isPending || !passphrase || passphrasesMismatch}
                className="gap-1.5 text-xs"
              >
                {exportSystem.isPending ? (
                  <Loader2 className="size-3.5 animate-spin" />
                ) : (
                  <Download className="size-3.5" />
                )}
                {exportSystem.isPending ? 'Exporting...' : 'Download Bundle'}
              </Button>
            </div>
          </form>
        </div>

        <div className="rounded-xl border border-border/50 bg-background/50 p-5">
          <form onSubmit={handleImport} className="space-y-4">
            <div>
              <h4 className="font-medium text-foreground text-sm">Restore from bundle</h4>
              <p className="mt-0.5 text-muted-foreground text-xs">
                Restore instance database, projects, and secrets from an existing .codedock export
                archive.
              </p>
            </div>

            <div className="grid grid-cols-1 gap-3 sm:grid-cols-2">
              <div className="space-y-1.5">
                <Label className="text-xs">Bundle File (.codedock)</Label>
                <Input
                  ref={fileInputRef}
                  type="file"
                  accept=".codedock,.tar.gz,.bundle"
                  onChange={(e) => setImportFile(e.target.files?.[0] || null)}
                  className="bg-muted/30 text-xs"
                  required
                />
              </div>
              <div className="space-y-1.5">
                <Label className="text-xs">Decryption Passphrase</Label>
                <Input
                  type="password"
                  value={importPassphrase}
                  onChange={(e) => setImportPassphrase(e.target.value)}
                  placeholder="Enter archive passphrase"
                  className="bg-muted/30 font-mono text-xs"
                  required
                />
              </div>
            </div>

            <div className="flex justify-end border-border/40 border-t pt-3">
              <Button
                type="submit"
                size="sm"
                variant="outline"
                disabled={importSystem.isPending || !importFile || !importPassphrase}
                className="gap-1.5 text-xs"
              >
                {importSystem.isPending ? (
                  <Loader2 className="size-3.5 animate-spin" />
                ) : (
                  <Upload className="size-3.5" />
                )}
                {importSystem.isPending ? 'Restoring...' : 'Restore Bundle'}
              </Button>
            </div>
          </form>
        </div>

        <div className="flex flex-col gap-3 rounded-xl border border-border/50 bg-background/50 p-5 sm:flex-row sm:items-center sm:justify-between">
          <div>
            <h4 className="font-medium text-foreground text-sm">Server Takeover & Discovery</h4>
            <p className="mt-0.5 text-muted-foreground text-xs">
              Take over existing servers and migrations from another Codedock instance.
            </p>
          </div>
          <ServerTakeoverDialog />
        </div>
      </div>
    </SettingsSection>
  );
}
