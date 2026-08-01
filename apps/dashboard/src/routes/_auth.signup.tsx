import { createFileRoute, Link, Navigate } from '@tanstack/react-router';
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
  const { data: publicSettings } = useGetPublicSettings();
  const { data: setupStatus, isLoading } = useGetSetupStatus();
  const registrationEnabled = publicSettings?.data?.registrationEnabled ?? true;

  if (!isLoading && !registrationEnabled && !setupStatus?.data?.setupRequired) {
    return <Navigate to="/signin" replace />;
  }

  return (
    <AuthPageFrame
      eyebrow={setupStatus?.data?.setupRequired ? 'First-time setup' : 'Create account'}
      title={
        setupStatus?.data?.setupRequired
          ? 'Set up your Codedock workspace.'
          : 'Create your Codedock account.'
      }
      description="Create your account, then organize the services your team runs."
    >
      <div className="space-y-6">
        {!setupStatus?.data?.setupRequired && <OAuthButtons />}
        <RegisterForm />

        <p className="border-border border-t pt-5 text-muted-foreground text-sm">
          Already have an account?{' '}
          <Link
            to="/signin"
            className="font-bold text-primary underline decoration-primary/35 underline-offset-4"
          >
            Sign in
          </Link>
        </p>
      </div>
    </AuthPageFrame>
  );
}
