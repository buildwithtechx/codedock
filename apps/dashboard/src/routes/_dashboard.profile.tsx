import { createFileRoute } from '@tanstack/react-router';
import { PageHeader } from '#/components/layout/page-header';
import { BillingSection } from '#/features/profile/billing-section';
import { Security2FASetup } from '#/features/profile/security-2fa-setup';
import {
  ProfileEmailForm,
  ProfileNameForm,
  ProfilePasswordForm,
} from '#/features/profile/user-profile-form';

export const Route = createFileRoute('/_dashboard/profile')({
  component: ProfilePage,
});

function ProfilePage() {
  return (
    <div className="space-y-6">
      <PageHeader
        title="Profile & security"
        description="Manage your personal profile and security preferences."
      />

      <div className="grid gap-6">
        <ProfileNameForm />
        <ProfileEmailForm />
        <ProfilePasswordForm />
        <Security2FASetup />
        <BillingSection />
      </div>
    </div>
  );
}
