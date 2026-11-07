import { ArchiveRestore, Loader2 } from 'lucide-react';
import { useEffect, useState } from 'react';
import { Button } from '@codedock/ui/button';
import { Dialog, DialogContent, DialogDescription, DialogTitle } from '@codedock/ui/dialog';
import {
  useApplyRestore,
  useCancelRestore,
  usePrepareRestore,
  useRestoreOperation,
} from '#/hooks/use-backups';
import type { BackupRecord, RestoreOperation } from '#/interfaces/backups';
import {
  isTerminal,
  RestoreApplyingStep,
  RestoreDoneStep,
  RestorePreparedStep,
  RestoreReviewStep,
  RestoreStepIndicator,
  terminalTone,
  type WizardStep,
} from './restore-wizard-steps';

export function RestoreWizard({
  record,
  sourceName,
  onClose,
  onDone,
}: {
  record: BackupRecord;
  sourceName: string;
  onClose: () => void;
  onDone?: () => void;
}) {
  const [step, setStep] = useState<WizardStep>('review');
  const [operationId, setOperationId] = useState<string | null>(null);
  const [confirmation, setConfirmation] = useState<string | null>(null);
  const [typed, setTyped] = useState('');
  const [busy, setBusy] = useState(false);
  const [cancelRequested, setCancelRequested] = useState(false);
  const [actionError, setActionError] = useState<string | null>(null);

  const prepare = usePrepareRestore();
  const apply = useApplyRestore();
  const cancel = useCancelRestore();
  const operationQuery = useRestoreOperation(operationId, step === 'applying' || step === 'done');
  const operation: RestoreOperation | null = operationQuery.data?.data ?? null;

  useEffect(() => {
    if (!operation) return;
    if (step === 'applying' && isTerminal(operation.status)) setStep('done');
  }, [operation, step]);

  useEffect(() => {
    if (step === 'done' && operation && terminalTone(operation.status) === 'success') onDone?.();
  }, [step, operation, onDone]);

  const startPrepare = async () => {
    if (busy) return;
    setBusy(true);
    setActionError(null);
    setStep('preparing');
    try {
      const res = await prepare.mutateAsync({ recordId: record.id });
      setOperationId(res.data.operation.id);
      setConfirmation(res.data.confirmation);
      setStep(res.data.operation.status === 'REVIEWED' ? 'prepared' : 'preparing');
      if (res.data.operation.status !== 'REVIEWED') {
        await operationQuery.refetch();
        setStep('prepared');
      }
    } catch (error) {
      setActionError(error instanceof Error ? error.message : 'Failed to prepare restore');
      setStep('review');
    } finally {
      setBusy(false);
    }
  };

  const applyRestore = async () => {
    if (busy || !operationId || !confirmation) return;
    if (!sourceName || typed !== sourceName) {
      setActionError(`Type "${sourceName}" to confirm`);
      return;
    }
    setBusy(true);
    setActionError(null);
    try {
      await apply.mutateAsync({ operationId, confirmation });
      setStep('applying');
      void operationQuery.refetch();
    } catch (error) {
      setActionError(error instanceof Error ? error.message : 'Failed to apply restore');
    } finally {
      setBusy(false);
    }
  };

  const cancelRestore = async () => {
    if (busy) return;
    if (!operationId) {
      onClose();
      return;
    }
    setBusy(true);
    setActionError(null);
    try {
      await cancel.mutateAsync(operationId);
      setCancelRequested(true);
      await operationQuery.refetch();
      const status = operationQuery.data?.data?.status ?? '';
      if (step !== 'applying' || isTerminal(status)) {
        if (step !== 'applying') onClose();
        else setStep('done');
      }
    } catch (error) {
      setActionError(error instanceof Error ? error.message : 'Failed to cancel restore');
    } finally {
      setBusy(false);
    }
  };

  const takenAt = record.startedAt ? new Date(record.startedAt).toLocaleString() : 'Unknown date';

  return (
    <Dialog open onOpenChange={(openState) => !openState && !busy && onClose()}>
      <DialogContent className="sm:max-w-xl">
        <DialogTitle className="flex items-center gap-2">
          <ArchiveRestore className="h-5 w-5 text-primary" />
          Restore snapshot
        </DialogTitle>
        <DialogDescription>
          Restoring <span className="font-medium text-foreground">{sourceName}</span> from {takenAt}
          .
        </DialogDescription>

        <RestoreStepIndicator currentStep={step} />

        {actionError && (
          <p role="alert" className="rounded-lg bg-destructive/10 p-3 text-destructive text-sm">
            {actionError}
          </p>
        )}

        {step === 'review' && (
          <RestoreReviewStep
            record={record}
            takenAt={takenAt}
            busy={busy}
            onClose={onClose}
            onStartPrepare={startPrepare}
          />
        )}

        {step === 'preparing' && (
          <div className="space-y-3">
            <div className="flex items-center gap-3 rounded-xl bg-muted/40 p-4 text-sm">
              <Loader2 className="h-4 w-4 animate-spin text-primary" />
              <div>
                <p className="font-medium">Verifying snapshot integrity</p>
                <p className="text-muted-foreground text-xs">
                  Checksums and restore targets are reviewed before anything changes.
                </p>
              </div>
            </div>
            <div className="flex items-center justify-end">
              <Button variant="ghost" onClick={cancelRestore} disabled={busy}>
                Cancel restore
              </Button>
            </div>
          </div>
        )}

        {step === 'prepared' && (
          <RestorePreparedStep
            sourceName={sourceName}
            typed={typed}
            busy={busy}
            onTypedChange={setTyped}
            onCancel={cancelRestore}
            onApply={applyRestore}
          />
        )}

        {step === 'applying' && (
          <RestoreApplyingStep
            operation={operation}
            cancelRequested={cancelRequested}
            busy={busy}
            onCancel={cancelRestore}
          />
        )}

        {step === 'done' && (
          <RestoreDoneStep
            operation={operation}
            sourceName={sourceName}
            takenAt={takenAt}
            onClose={onClose}
          />
        )}
      </DialogContent>
    </Dialog>
  );
}
