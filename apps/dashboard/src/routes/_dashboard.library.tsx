import { createFileRoute } from '@tanstack/react-router';
import { LibraryPage } from '#/features/library';

export const Route = createFileRoute('/_dashboard/library')({
  component: LibraryRoutePage,
});

function LibraryRoutePage() {
  return <LibraryPage />;
}
