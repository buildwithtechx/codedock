import { createFileRoute, Link } from '@tanstack/react-router';
import { AlertCircle } from 'lucide-react';
import { Button } from '#/components/ui/button';
import { AuthPageFrame } from '#/features/auth';
import { ForgotPasswordForm } from '#/features/auth/forgot-password-form';
import { useGetPublicSettings } from '#/features/settings';
import { useSystemStore } from '#/stores/system-store';

export const Route = createFileRoute('/_auth/forgot-password')({
  component: ForgotPasswordPage,
  head: () => {
    const siteName = useSystemStore.getState().siteName;
    return { meta: [{ title: `Reset Password - ${siteName}` }] };
  },
});

function ForgotPasswordPage() {
  const { data, isLoading } = useGetPublicSettings();
  const emailEnabled = data?.data?.emailEnabled;

  return (
    <AuthPageFrame
      eyebrow="Account recovery"
      title="Recover your account."
      description="Enter your email and we will send a reset link if the account exists."
    >
      {!isLoading && emailEnabled === false ? (
        <div className="border-primary border-l-2 py-1 pl-5">
          <div className="flex items-center gap-3">
            <AlertCircle className="h-5 w-5 text-[#b42318]" />
            <p className="font-medium text-foreground text-sm">Email not configured</p>
          </div>
          <p className="mt-3 text-muted-foreground text-sm leading-6">
            Your team has not enabled email. Contact an administrator to restore access.
          </p>
          <Button asChild variant="outline" className="mt-5">
            <Link to="/signin">Back to sign in</Link>
          </Button>
        </div>
      ) : (
        <ForgotPasswordForm />
      )}
    </AuthPageFrame>
  );
}
