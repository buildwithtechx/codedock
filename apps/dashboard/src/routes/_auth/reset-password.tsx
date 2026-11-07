import { createFileRoute, useNavigate } from '@tanstack/react-router';
import { AuthPageFrame } from '#/features/auth';
import { ResetPasswordForm } from '#/features/auth/reset-password-form';
import { useSystemStore } from '#/stores/system-store';

export const Route = createFileRoute('/_auth/reset-password')({
  validateSearch: (search: Record<string, unknown>) => {
    return {
      token: (search.token as string) || '',
    };
  },
  component: ResetPasswordPage,
  head: () => {
    const siteName = useSystemStore.getState().siteName;
    return { meta: [{ title: `New Password - ${siteName}` }] };
  },
});

function ResetPasswordPage() {
  const { token } = Route.useSearch();
  const navigate = useNavigate();

  if (!token) {
    return (
      <AuthPageFrame
        eyebrow="Account recovery"
        title="This link is incomplete."
        description="The password reset token is missing. Request another reset link and try again."
      >
        <button
          type="button"
          onClick={() => navigate({ to: '/signin' })}
          className="text-primary text-sm underline-offset-4 hover:underline"
        >
          Return to sign in
        </button>
      </AuthPageFrame>
    );
  }

  return (
    <AuthPageFrame
      eyebrow="Account recovery"
      title="Set a new password."
      description="Use a strong password you have not used for this account before."
    >
      <ResetPasswordForm token={token} />
    </AuthPageFrame>
  );
}
