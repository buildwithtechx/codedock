import type { BackupRecord } from '#/interfaces/backups';
import { ReviewedRestoreDialog } from '#/features/backups/reviewed-restore-dialog';

type DatabaseRestoreDialogProps = {
  record: BackupRecord | null;
  onClose: () => void;
};
export function DatabaseRestoreDialog({ record, onClose }: DatabaseRestoreDialogProps) {
  return record ? (
    <ReviewedRestoreDialog
      key={record.id}
      recordId={record.id}
      sourceDatabaseId={record.databaseId}
      onClose={onClose}
      allowDatabaseTarget
    />
  ) : null;
}
