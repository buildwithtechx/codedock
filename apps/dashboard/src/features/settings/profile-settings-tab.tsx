import { BillingSection } from '#/features/profile/billing-section';
import { PasskeysSection } from '#/features/profile/passkeys-section';
import { Security2FASetup } from '#/features/profile/security-2fa-setup';
import {
  ProfileEmailForm,
  ProfileNameForm,
  ProfilePasswordForm,
} from '#/features/profile/user-profile-form';
import { useGetPublicSettings } from '#/hooks/use-settings';

export function ProfileSettingsTab() {
  const { data: publicSettings } = useGetPublicSettings();

  return (
    <div className="space-y-6">
      <ProfileNameForm />
      <ProfileEmailForm />
      <ProfilePasswordForm />
      <PasskeysSection />
      <Security2FASetup />
      {publicSettings?.data?.cloudMode && <BillingSection />}
    </div>
  );
}
