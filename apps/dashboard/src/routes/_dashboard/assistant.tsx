import { createFileRoute } from '@tanstack/react-router';
import { AssistantPage } from '#/features/assistant/assistant-page';

export const Route = createFileRoute('/_dashboard/assistant')({
  component: AssistantPage,
});
