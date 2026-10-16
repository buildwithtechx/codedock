import { createFileRoute } from '@tanstack/react-router';
import { useGetPublicSettings } from '#/features/settings';
import { CloudMailClient } from '#/features/mail/cloud-mail-client';
import { MailWizardView } from '#/features/mail/mail-wizard-view';

export const Route = createFileRoute('/_dashboard/apps/new/mail')({
  component: MailAppRouteComponent,
});

function MailAppRouteComponent() {
  const { data: publicRes } = useGetPublicSettings();
  const isCloud = publicRes?.data?.cloudMode ?? false;

  if (isCloud) {
    return <CloudMailClient />;
  }

  return <MailWizardView />;
}
