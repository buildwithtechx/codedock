import { createFileRoute } from '@tanstack/react-router';
import { NotFoundComponent } from '#/components/not-found';

export const Route = createFileRoute('/$')({
  component: NotFoundComponent,
});
