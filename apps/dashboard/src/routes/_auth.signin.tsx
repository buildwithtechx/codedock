import { createFileRoute, Link } from '@tanstack/react-router';
import { AuthPageFrame, LoginForm, OAuthButtons } from '#/features/auth';
import { useGetPublicSettings } from '#/features/settings';
import { useSystemStore } from '#/stores/system-store';

export const Route = createFileRoute('/_auth/signin')({
  component: LoginPage,
  head: () => {
    const siteName = useSystemStore.getState().siteName;
    return { meta: [{ title: `Sign In - ${siteName}` }] };
  },
});

function LoginPage() {
  const { data: publicSettings } = useGetPublicSettings();
  const registrationEnabled = publicSettings?.data?.registrationEnabled ?? true;

  return (
    <AuthPageFrame
      eyebrow="Account access"
      title="Pick up where you left off."
      description="Sign in to manage your services, releases, and servers."
    >
      <div className="space-y-6">
        <OAuthButtons />
        <LoginForm />

        {registrationEnabled && (
          <p className="border-border border-t pt-5 text-center text-muted-foreground text-sm">
            Don't have an account?{' '}
            <Link
              to="/signup"
              className="font-medium text-primary underline-offset-4 hover:underline"
            >
              Create one
            </Link>
          </p>
        )}
      </div>
    </AuthPageFrame>
  );
}
