import { createFileRoute } from '@tanstack/react-router';
import { SupportCenter } from '#/features/support';

export const Route = createFileRoute('/_dashboard/support')({
  component: SupportPage,
});

function SupportPage() {
  return <SupportCenter />;
}
