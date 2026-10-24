import { createFileRoute } from '@tanstack/react-router';
import { MailPage } from '#/features/mail/mail-page';

export const Route = createFileRoute('/_dashboard/emails')({
  component: MailPage,
});
