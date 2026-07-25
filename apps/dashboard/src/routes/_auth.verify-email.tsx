import { useMutation } from '@tanstack/react-query';
import { createFileRoute, useNavigate, useSearch } from '@tanstack/react-router';
import { CheckCircle, Loader2, XCircle } from 'lucide-react';
import { useEffect, useState } from 'react';
import { toast } from 'sonner';
import { apiClient } from '#/lib/apiClient';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';

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
      <Card className="mx-auto w-full max-w-md">
        <CardHeader className="text-center">
          <CardTitle className="text-2xl">Missing Token</CardTitle>
          <CardDescription>No verification token provided.</CardDescription>
        </CardHeader>
        <CardContent className="flex justify-center">
          <Button onClick={() => navigate({ to: '/' })}>Return Home</Button>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card className="mx-auto w-full max-w-md">
      <CardHeader className="text-center">
        <CardTitle className="text-2xl">Email Verification</CardTitle>
        <CardDescription>
          {status === 'loading' && 'Verifying your email address...'}
          {status === 'success' && 'Your email has been verified!'}
          {status === 'error' && 'Verification failed'}
        </CardDescription>
      </CardHeader>
      <CardContent className="flex flex-col items-center justify-center space-y-4">
        {status === 'loading' && (
          <div className="flex flex-col items-center justify-center py-6">
            <Loader2 className="mb-4 h-12 w-12 animate-spin text-primary" />
            <p className="text-muted-foreground text-sm">Please wait while we verify your email.</p>
          </div>
        )}
        {status === 'success' && (
          <div className="flex flex-col items-center justify-center py-6">
            <CheckCircle className="mb-4 h-12 w-12 text-green-500" />
            <p className="text-muted-foreground text-sm">Redirecting you to the dashboard...</p>
          </div>
        )}
        {status === 'error' && (
          <div className="flex flex-col items-center justify-center py-6">
            <XCircle className="mb-4 h-12 w-12 text-destructive" />
            <p className="mb-4 text-center font-medium text-destructive text-sm">{errorMessage}</p>
            <Button onClick={() => navigate({ to: '/' })} variant="outline">
              Return Home
            </Button>
          </div>
        )}
      </CardContent>
    </Card>
  );
}
