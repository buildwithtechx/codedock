import { createFileRoute, Link, Navigate } from '@tanstack/react-router';
import { useEffect } from 'react';
import { AuthPageFrame, OAuthButtons, RegisterForm } from '#/features/auth';
import { useGetPublicSettings, useGetSetupStatus } from '#/features/settings';
import { useSystemStore } from '#/stores/system-store';

export const Route = createFileRoute('/_auth/signup')({
  component: RegisterPage,
  head: () => {
    const siteName = useSystemStore.getState().siteName;
    return { meta: [{ title: `Sign Up - ${siteName}` }] };
  },
});

function RegisterPage() {
  const siteName = useSystemStore((state) => state.siteName);
  const { data: publicSettings } = useGetPublicSettings();
  const { data: setupStatus, isLoading } = useGetSetupStatus();
  const registrationEnabled = publicSettings?.data?.registrationEnabled ?? true;
  const isOnboarding = setupStatus?.data?.setupRequired ?? false;

  useEffect(() => {
    document.title = `${isOnboarding ? 'Onboarding' : 'Sign Up'} - ${siteName}`;
  }, [isOnboarding, siteName]);

  if (!isLoading && !registrationEnabled && !setupStatus?.data?.setupRequired) {
    return <Navigate to="/signin" replace />;
  }

  return (
    <AuthPageFrame
      eyebrow={isOnboarding ? 'First-time setup' : 'Create account'}
      title={isOnboarding ? 'Set up your Codedock workspace.' : 'Create your Codedock account.'}
      description="Create your account, then organize the services your team runs."
    >
      <div className="space-y-6">
        {!isOnboarding && <OAuthButtons />}
        <RegisterForm />

        <p className="border-border border-t pt-5 text-muted-foreground text-sm">
          Already have an account?{' '}
          <Link
            to="/signin"
            className="font-medium text-primary underline-offset-4 hover:underline"
          >
            Sign in
          </Link>
        </p>
      </div>
    </AuthPageFrame>
  );
}
