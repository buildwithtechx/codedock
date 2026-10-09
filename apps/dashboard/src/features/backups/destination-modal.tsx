import { ArrowLeft, Cloud, Lock, Server } from 'lucide-react';
import { useEffect, useState } from 'react';
import { Button } from '#/components/ui/button';
import { Dialog, DialogContent, DialogDescription, DialogTitle } from '#/components/ui/dialog';
import { DestinationConfigureForm } from './destination-configure-form';
import type { DestinationKind, S3Destination } from './interfaces';
import { S3ProviderMark } from './s3-destination-fields';
import { s3Providers } from './s3-providers';

export function DestinationModal({
  open,
  onOpenChange,
  editing = null,
  initialKind = 's3',
  onSaved,
}: {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  editing?: S3Destination | null;
  initialKind?: DestinationKind;
  onSaved?: () => void;
}) {
  const [step, setStep] = useState<'pick' | 'configure'>(editing ? 'configure' : 'pick');
  const [kind, setKind] = useState<DestinationKind>(editing ? 's3' : initialKind);

  useEffect(() => {
    if (!open) return;
    if (editing) {
      setKind('s3');
      setStep('configure');
    } else {
      setKind(initialKind);
      setStep('pick');
    }
  }, [open, editing, initialKind]);

  const title = editing
    ? `Edit ${editing.name}`
    : step === 'pick'
      ? 'Add destination'
      : kind === 's3'
        ? 'New S3 destination'
        : 'New SFTP destination';

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="gap-0 overflow-hidden p-0 sm:max-w-2xl">
        <div className="flex items-center justify-between gap-3 border-border/40 border-b px-6 py-4">
          <div className="flex min-w-0 items-center gap-2">
            {step === 'configure' && !editing && (
              <Button
                variant="ghost"
                size="icon"
                className="h-8 w-8"
                onClick={() => setStep('pick')}
              >
                <ArrowLeft className="h-4 w-4" />
              </Button>
            )}
            <div className="min-w-0">
              <DialogTitle className="truncate">{title}</DialogTitle>
              <DialogDescription>
                {step === 'pick'
                  ? 'Choose where snapshots are stored off-server.'
                  : 'Credentials are encrypted before they are saved.'}
              </DialogDescription>
            </div>
          </div>
          <span className="hidden shrink-0 items-center gap-1.5 text-muted-foreground text-xs sm:flex">
            <Lock className="h-3.5 w-3.5" />
            Encrypted at rest
          </span>
        </div>
        <div className="max-h-[75vh] overflow-y-auto px-6 py-6">
          {step === 'pick' ? (
            <div className="grid gap-4 sm:grid-cols-2">
              <button
                type="button"
                onClick={() => {
                  setKind('s3');
                  setStep('configure');
                }}
                className="group flex items-start gap-4 rounded-2xl border border-border/60 bg-card p-5 text-left transition-all hover:-translate-y-0.5 hover:border-primary/40 hover:shadow-md"
              >
                <span className="flex h-12 w-12 shrink-0 items-center justify-center rounded-xl bg-muted/60 text-muted-foreground transition-transform group-hover:scale-105">
                  <Cloud className="h-5 w-5" />
                </span>
                <span className="min-w-0">
                  <span className="block font-semibold text-base">S3 compatible</span>
                  <span className="mt-1 block text-muted-foreground text-sm leading-relaxed">
                    Object storage with per-bucket prefixes and regions.
                  </span>
                  <span className="mt-3 flex items-center gap-2.5">
                    {s3Providers
                      .filter((provider) => provider.id !== 'minio')
                      .map((provider) => (
                        <S3ProviderMark
                          key={provider.id}
                          provider={provider}
                          className="size-4 opacity-90"
                        />
                      ))}
                    <span className="font-medium text-muted-foreground/70 text-xs">
                      R2 · S3 · B2 · MinIO
                    </span>
                  </span>
                </span>
              </button>
              <button
                type="button"
                onClick={() => {
                  setKind('sftp');
                  setStep('configure');
                }}
                className="group flex items-start gap-4 rounded-2xl border border-border/60 bg-card p-5 text-left transition-all hover:-translate-y-0.5 hover:border-primary/40 hover:shadow-md"
              >
                <span className="flex h-12 w-12 shrink-0 items-center justify-center rounded-xl bg-muted/60 text-muted-foreground transition-transform group-hover:scale-105">
                  <Server className="h-5 w-5" />
                </span>
                <span className="min-w-0">
                  <span className="block font-semibold text-base">SFTP server</span>
                  <span className="mt-1 block text-muted-foreground text-sm leading-relaxed">
                    Push snapshots to any SSH server over SFTP.
                  </span>
                  <span className="mt-2 block font-medium text-muted-foreground/70 text-xs uppercase tracking-wider">
                    Password · Key auth
                  </span>
                </span>
              </button>
            </div>
          ) : (
            <DestinationConfigureForm
              key={editing?.id ?? `new-${kind}`}
              kind={kind}
              editing={editing}
              onCancel={() => onOpenChange(false)}
              onSaved={() => {
                onOpenChange(false);
                onSaved?.();
              }}
            />
          )}
        </div>
      </DialogContent>
    </Dialog>
  );
}
