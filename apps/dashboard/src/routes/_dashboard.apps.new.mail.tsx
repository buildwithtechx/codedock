import { createFileRoute, useNavigate } from '@tanstack/react-router';
import { useEffect } from 'react';

export const Route = createFileRoute('/_dashboard/apps/new/mail')({
  component: MailAppRouteComponent,
});

function MailAppRouteComponent() {
  const navigate = useNavigate();

  useEffect(() => {
    void navigate({
      to: '/apps/new/$appId',
      params: { appId: 'mail' },
      replace: true,
    });
  }, [navigate]);

  return null;
}
