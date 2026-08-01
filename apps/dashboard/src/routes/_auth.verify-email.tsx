import { useMutation } from '@tanstack/react-query';
import { createFileRoute, useNavigate, useSearch } from '@tanstack/react-router';
import { CheckCircle, Loader2, XCircle } from 'lucide-react';
import { useEffect, useState } from 'react';
import { toast } from 'sonner';
import { AuthPageFrame } from '#/features/auth';
import { apiClient } from '#/lib/api-client';
import { Button } from '@/components/ui/button';

export const Route = createFileRoute('/_auth/verify-email')({
  component: VerifyEmailPage,
});

function VerifyEmailPage() {
  const { token } = useSearch({ from: '/_auth/verify-email' }) as { token?: string };
  const navigate = useNavigate();
  const [status, setStatus] = useState<'loading' | 'success' | 'error' | 'idle'>('idle');
  const [errorMessage, setErrorMessage] = useState('');

  const verifyMutation = useMutation({
    mutationFn: async (t: string) => {
      return apiClient.post('/auth/email/verify', { token: t });
    },
    onSuccess: () => {
      setStatus('success');
      toast.success('Email verified successfully!');
      setTimeout(() => {
        navigate({ to: '/' });
      }, 2000);
    },
    onError: (error: any) => {
      setStatus('error');
      setErrorMessage(
        error.response?.data?.error || 'Verification failed. The token may be invalid or expired.'
      );
    },
  });

  useEffect(() => {
    if (token && status === 'idle') {
      setStatus('loading');
      verifyMutation.mutate(token);
    }
  }, [token, status, verifyMutation]);

  if (!token) {
    return (
      <AuthPageFrame
        eyebrow="Email verification"
        title="This link is incomplete."
        description="No verification token was provided. Return home and request a new verification email."
      >
        <Button
          onClick={() => navigate({ to: '/' })}
          className="rounded-lg bg-primary text-primary-foreground hover:bg-primary-hover"
        >
          Return home
        </Button>
      </AuthPageFrame>
    );
  }

  return (
    <AuthPageFrame
      eyebrow="Email verification"
      title={
        status === 'success'
          ? 'You are verified.'
          : status === 'error'
            ? 'Verification failed.'
            : 'Verifying your email.'
      }
      description={
        status === 'success'
          ? 'Your account is ready. Taking you to the dashboard now.'
          : status === 'error'
            ? 'The link may be invalid or expired.'
            : 'Confirming this email address for your Codedock account.'
      }
    >
      <div className="flex flex-col items-start gap-4">
        {status === 'loading' && (
          <div className="flex items-center gap-3 text-muted-foreground text-sm">
            <Loader2 className="h-5 w-5 animate-spin text-primary" />
            Please wait while we verify your email.
          </div>
        )}
        {status === 'success' && (
          <div className="flex items-center gap-3 text-muted-foreground text-sm">
            <CheckCircle className="h-5 w-5 text-[#3b6e32]" />
            Redirecting you to the dashboard.
          </div>
        )}
        {status === 'error' && (
          <div className="border-[#b42318] border-l-2 pl-4">
            <div className="flex items-center gap-3">
              <XCircle className="h-5 w-5 text-[#b42318]" />
              <p className="font-medium text-[#b42318] text-sm">{errorMessage}</p>
            </div>
            <Button
              onClick={() => navigate({ to: '/' })}
              variant="outline"
              className="mt-5 rounded-lg border-border bg-background text-foreground hover:bg-muted"
            >
              Return Home
            </Button>
          </div>
        )}
      </div>
    </AuthPageFrame>
  );
}
