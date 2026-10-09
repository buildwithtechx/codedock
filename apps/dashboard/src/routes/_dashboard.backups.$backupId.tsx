import { createFileRoute } from '@tanstack/react-router';
import { DestinationDetail } from '#/features/backups/destination-detail';

export const Route = createFileRoute('/_dashboard/backups/$backupId')({
  component: DestinationDetailPage,
});

function DestinationDetailPage() {
  const { backupId } = Route.useParams();
  return <DestinationDetail destinationId={backupId} />;
}
