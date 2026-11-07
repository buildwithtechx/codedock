import { ReviewedRestoreDialog } from './reviewed-restore-dialog';

export function VolumeRestoreDialog({
  recordId,
  onClose,
}: {
  recordId: string | null;
  onClose: () => void;
}) {
  return recordId ? (
    <ReviewedRestoreDialog key={recordId} recordId={recordId} onClose={onClose} />
  ) : null;
}
