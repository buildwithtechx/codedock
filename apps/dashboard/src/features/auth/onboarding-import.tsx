import { ArrowUpRight, FileUp, Loader2, LockKeyhole, Upload } from 'lucide-react';
import { useRef, useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '#/components/ui/dialog';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import { useGetSetupStatus } from '#/features/settings';
import { apiClient } from '#/lib/api-client';

export function OnboardingImport() {
  const { data: setupStatus } = useGetSetupStatus();
  const [open, setOpen] = useState(false);
  const [passphrase, setPassphrase] = useState('');
  const [isImporting, setIsImporting] = useState(false);
  const [fileName, setFileName] = useState('');
  const fileRef = useRef<HTMLInputElement>(null);

  if (!setupStatus?.data?.setupRequired) return null;

  const handleImport = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const bundle = fileRef.current?.files?.[0];
    if (!bundle || !passphrase) return;

    const payload = new FormData();
    payload.set('bundle', bundle);
    payload.set('passphrase', passphrase);
    setIsImporting(true);

    try {
      await apiClient.post('/system/setup/import', payload);
      toast.success('Codedock instance imported. Sign in with an imported account.');
      window.location.assign('/signin');
    } catch (error) {
      toast.error(error instanceof Error ? error.message : 'Failed to import Codedock instance.');
    } finally {
      setIsImporting(false);
    }
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <button
          type="button"
          className="group flex w-fit items-center gap-2 rounded-lg px-1 py-2 font-semibold text-white text-xs uppercase tracking-[0.12em] transition-colors hover:bg-white/10 hover:text-[#d8c7ff]"
        >
          <FileUp className="size-4" />
          Import existing Codedock
          <ArrowUpRight className="size-3.5 transition-transform group-hover:translate-x-0.5 group-hover:-translate-y-0.5" />
        </button>
      </DialogTrigger>
      <DialogContent className="overflow-hidden p-0 sm:max-w-2xl">
        <div className="grid sm:grid-cols-[0.82fr_1.18fr]">
          <div className="bg-primary-soft p-7 sm:p-9">
            <div className="flex size-11 items-center justify-center rounded-xl bg-primary text-primary-foreground">
              <FileUp className="size-5" />
            </div>
            <DialogHeader className="mt-7 text-left">
              <DialogTitle className="text-2xl tracking-[-0.04em]">
                Bring your workspace home.
              </DialogTitle>
              <DialogDescription className="text-sm leading-6">
                Restore projects, services, and settings from an encrypted Codedock bundle.
              </DialogDescription>
            </DialogHeader>
            <div className="mt-8 space-y-4 text-muted-foreground text-sm">
              <p className="flex gap-3">
                <span className="font-semibold text-primary">01</span>Choose your exported bundle.
              </p>
              <p className="flex gap-3">
                <span className="font-semibold text-primary">02</span>Enter the bundle passphrase.
              </p>
              <p className="flex gap-3">
                <span className="font-semibold text-primary">03</span>Sign in with a restored
                account.
              </p>
            </div>
          </div>
          <form onSubmit={handleImport} className="p-7 sm:p-9">
            <div className="space-y-6">
              <div className="space-y-2">
                <Label htmlFor="migration-bundle">Migration bundle</Label>
                <Input
                  ref={fileRef}
                  id="migration-bundle"
                  type="file"
                  accept=".codedock,application/octet-stream"
                  className="sr-only"
                  onChange={(event) => setFileName(event.target.files?.[0]?.name ?? '')}
                  required
                />
                <label
                  htmlFor="migration-bundle"
                  className="flex min-h-34 cursor-pointer flex-col items-center justify-center rounded-xl border border-primary/30 border-dashed bg-primary/5 px-5 text-center transition-colors hover:border-primary/60 hover:bg-primary/10"
                >
                  <Upload className="size-5 text-primary" />
                  <span className="mt-3 font-medium text-foreground text-sm">
                    {fileName || 'Choose bundle file'}
                  </span>
                  <span className="mt-1 text-muted-foreground text-xs">
                    .codedock bundle, up to 500 MB
                  </span>
                </label>
              </div>
              <div className="space-y-2">
                <Label htmlFor="migration-passphrase">Bundle passphrase</Label>
                <div className="relative">
                  <LockKeyhole className="pointer-events-none absolute top-1/2 left-3.5 size-4 -translate-y-1/2 text-muted-foreground" />
                  <Input
                    id="migration-passphrase"
                    type="password"
                    value={passphrase}
                    onChange={(event) => setPassphrase(event.target.value)}
                    placeholder="Enter the export passphrase"
                    className="h-12 pl-10"
                    required
                  />
                </div>
              </div>
            </div>
            <div className="mt-8">
              <Button
                type="submit"
                className="h-12 w-full"
                disabled={isImporting || !fileName || !passphrase}
              >
                {isImporting && <Loader2 className="animate-spin" />}
                {isImporting ? 'Importing...' : 'Import bundle'}
              </Button>
            </div>
          </form>
        </div>
      </DialogContent>
    </Dialog>
  );
}
