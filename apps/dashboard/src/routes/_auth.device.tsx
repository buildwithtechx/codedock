import { createFileRoute, Link } from '@tanstack/react-router';
import { ArrowRight, CheckCircle2, Laptop, ShieldCheck, XCircle } from 'lucide-react';
import { useState } from 'react';
import { z } from 'zod';
import { Button } from '@codedock/ui';
import { AuthPageFrame, LoginForm } from '#/features/auth';
import { authService } from '#/services/auth';
import { useAuthStore } from '#/stores/auth-store';
import { useSystemStore } from '#/stores/system-store';

const deviceSearchSchema = z.object({
  code: z.string().optional(),
});

export const Route = createFileRoute('/_auth/device')({
  validateSearch: (search) => deviceSearchSchema.parse(search),
  component: DeviceAuthorizePage,
  head: () => {
    const siteName = useSystemStore.getState().siteName;
    return { meta: [{ title: `Authorize Device - ${siteName}` }] };
  },
});

function DeviceAuthorizePage() {
  const { code: initialCode } = Route.useSearch();
  const { isAuthenticated, user } = useAuthStore();

  const [inputCode, setInputCode] = useState(initialCode || '');
  const [activeCode, setActiveCode] = useState(initialCode || '');
  const [deviceInfo, setDeviceInfo] = useState<{
    valid: boolean;
    clientName: string;
    userCode: string;
  } | null>(null);
  const [loadingVerify, setLoadingVerify] = useState(false);
  const [errorMsg, setErrorMsg] = useState<string | null>(null);
  const [authResult, setAuthResult] = useState<{
    status: 'approved' | 'denied';
    token?: string;
  } | null>(null);
  const [submitting, setSubmitting] = useState(false);

  // Auto-verify if code is present in URL
  useState(() => {
    if (initialCode && isAuthenticated) {
      void verifyCode(initialCode);
    }
  });

  async function verifyCode(codeToVerify: string) {
    const clean = codeToVerify.trim().toUpperCase();
    if (!clean) return;
    setLoadingVerify(true);
    setErrorMsg(null);
    try {
      const data = await authService.verifyDeviceCode(clean);
      setDeviceInfo(data);
      setActiveCode(clean);
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Invalid or expired code';
      setErrorMsg(msg);
      setDeviceInfo(null);
    } finally {
      setLoadingVerify(false);
    }
  }

  async function handleAction(action: 'approve' | 'deny') {
    if (!activeCode) return;
    setSubmitting(true);
    setErrorMsg(null);
    try {
      const res = await authService.authorizeDevice(activeCode, action);
      setAuthResult({
        status: action === 'approve' ? 'approved' : 'denied',
        token: res.token,
      });
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Action failed';
      setErrorMsg(msg);
    } finally {
      setSubmitting(false);
    }
  }

  if (!isAuthenticated) {
    return (
      <AuthPageFrame
        eyebrow="Device Authorization"
        title="Sign in to authorize your device"
        description="A client or application is requesting access to your Codedock account."
      >
        <div className="space-y-6">
          <LoginForm />
          <p className="border-border border-t pt-5 text-center text-muted-foreground text-sm">
            Need an account?{' '}
            <Link
              to="/signup"
              className="font-medium text-primary underline-offset-4 hover:underline"
            >
              Sign up
            </Link>
          </p>
        </div>
      </AuthPageFrame>
    );
  }

  return (
    <AuthPageFrame
      eyebrow="Device Authorization"
      title="Authorize Client Access"
      description="Connect Codedock Mobile, VS Code Extension, or remote tools to your account."
    >
      <div className="space-y-6">
        {authResult ? (
          <div className="rounded-xl border border-border bg-card p-6 text-center space-y-4">
            {authResult.status === 'approved' ? (
              <>
                <div className="mx-auto flex h-14 w-14 items-center justify-center rounded-full bg-emerald-500/10 text-emerald-500">
                  <CheckCircle2 className="h-8 w-8" />
                </div>
                <div className="space-y-1">
                  <h3 className="text-lg font-bold text-foreground">Authorization Granted!</h3>
                  <p className="text-sm text-muted-foreground">
                    You have successfully authorized{' '}
                    <span className="font-semibold text-foreground">
                      {deviceInfo?.clientName || 'your device'}
                    </span>
                    .
                  </p>
                </div>
                {authResult.token && (
                  <div className="pt-2">
                    <a
                      href={`codedock://auth/callback?token=${encodeURIComponent(authResult.token)}&email=${encodeURIComponent(user?.email || '')}`}
                      className="inline-flex items-center gap-2 rounded-lg bg-primary px-4 py-2.5 text-sm font-semibold text-white shadow hover:opacity-90 transition"
                    >
                      Open in Codedock App <ArrowRight className="h-4 w-4" />
                    </a>
                  </div>
                )}
                <p className="text-xs text-muted-foreground pt-2">
                  You can now close this tab and return to your app or editor.
                </p>
              </>
            ) : (
              <>
                <div className="mx-auto flex h-14 w-14 items-center justify-center rounded-full bg-destructive/10 text-destructive">
                  <XCircle className="h-8 w-8" />
                </div>
                <h3 className="text-lg font-bold text-foreground">Request Declined</h3>
                <p className="text-sm text-muted-foreground">
                  The authorization request was declined. No tokens were minted.
                </p>
              </>
            )}
          </div>
        ) : !deviceInfo ? (
          <div className="space-y-4">
            <div className="rounded-lg border border-border bg-card p-4 space-y-3">
              <label htmlFor="user-code" className="block text-sm font-medium text-foreground">
                Enter 8-character verification code
              </label>
              <input
                id="user-code"
                type="text"
                placeholder="ABCD-1234"
                value={inputCode}
                onChange={(e) => setInputCode(e.target.value.toUpperCase())}
                className="w-full rounded-md border border-input bg-background px-3 py-2 text-center font-mono text-lg tracking-widest uppercase focus:outline-none focus:ring-2 focus:ring-primary"
              />
              {errorMsg && <p className="text-xs text-destructive">{errorMsg}</p>}
              <Button
                type="button"
                className="w-full"
                disabled={!inputCode.trim() || loadingVerify}
                onClick={() => verifyCode(inputCode)}
              >
                {loadingVerify ? 'Verifying...' : 'Continue'}
              </Button>
            </div>
          </div>
        ) : (
          <div className="rounded-xl border border-border bg-card p-6 space-y-5">
            <div className="flex items-center gap-3">
              <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-primary/10 text-primary">
                <Laptop className="h-6 w-6" />
              </div>
              <div>
                <h3 className="font-bold text-foreground text-base">{deviceInfo.clientName}</h3>
                <p className="text-xs text-muted-foreground">Code: {deviceInfo.userCode}</p>
              </div>
            </div>

            <div className="rounded-lg bg-muted/50 p-3 text-xs text-muted-foreground space-y-1">
              <div className="flex items-center gap-1.5 text-foreground font-medium">
                <ShieldCheck className="h-4 w-4 text-primary" />
                <span>Signed in as {user?.email}</span>
              </div>
              <p>This client will receive an API token to manage workloads and resources.</p>
            </div>

            {errorMsg && <p className="text-xs text-destructive">{errorMsg}</p>}

            <div className="flex gap-3 pt-2">
              <Button
                type="button"
                variant="outline"
                className="flex-1"
                disabled={submitting}
                onClick={() => handleAction('deny')}
              >
                Decline
              </Button>
              <Button
                type="button"
                className="flex-1"
                disabled={submitting}
                onClick={() => handleAction('approve')}
              >
                {submitting ? 'Authorizing...' : 'Authorize'}
              </Button>
            </div>
          </div>
        )}
      </div>
    </AuthPageFrame>
  );
}
