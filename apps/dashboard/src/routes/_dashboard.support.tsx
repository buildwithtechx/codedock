import { createFileRoute, useNavigate } from '@tanstack/react-router';
import { useEffect } from 'react';
import { SupportCenter } from '#/features/support';
import { useGetPublicSettings } from '#/features/settings';

export const Route = createFileRoute('/_dashboard/support')({
  component: SupportPage,
});

function SupportPage() {
  const navigate = useNavigate();
  const { data: publicRes, isLoading } = useGetPublicSettings();
  const isCloud = publicRes?.data?.cloudMode ?? false;

  useEffect(() => {
    if (!isLoading && !isCloud) {
      void navigate({ to: '/' });
    }
  }, [isLoading, isCloud, navigate]);

  if (!isCloud) {
    return null;
  }

  return <SupportCenter />;
}
