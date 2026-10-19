import { zodResolver } from '@hookform/resolvers/zod';
import { browserSupportsWebAuthn } from '@simplewebauthn/browser';
import { Link } from '@tanstack/react-router';
import { Eye, EyeOff, Fingerprint, KeyRound, Loader2, Lock, Mail } from 'lucide-react';
import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { z } from 'zod';
import { Button } from '@codedock/ui/button';
import { Input } from '@codedock/ui/input';
import { Label } from '@codedock/ui/label';
import { useLogin } from '#/hooks/use-auth';
import { usePasskeyLogin } from '#/hooks/use-passkeys';
import type { AuthenticationResponseJSON } from '#/interfaces/passkey';

const loginSchema = z.object({
  email: z.string().email('Please enter a valid email address'),
  password: z.string().min(1, 'Password is required'),
});

const totpSchema = z.object({
  totpCode: z.string().min(6, 'Code must be at least 6 characters'),
});

type LoginSchema = z.infer<typeof loginSchema>;
type TOTPSchema = z.infer<typeof totpSchema>;

export const LoginForm = () => {
  const { mutateAsync: login, isPending } = useLogin();
  const { mutateAsync: passkeyLogin, isPending: isPasskeyPending } = usePasskeyLogin();
  const [showPassword, setShowPassword] = useState(false);
  const [pendingCredentials, setPendingCredentials] = useState<LoginSchema | null>(null);
  const [pendingPasskeyAssertion, setPendingPasskeyAssertion] =
    useState<AuthenticationResponseJSON | null>(null);

  const webAuthnSupported = typeof window !== 'undefined' && browserSupportsWebAuthn();

  const loginForm = useForm<LoginSchema>({
    resolver: zodResolver(loginSchema),
    defaultValues: { email: '', password: '' },
  });

  const totpForm = useForm<TOTPSchema>({
    resolver: zodResolver(totpSchema),
    defaultValues: { totpCode: '' },
  });

  const onLoginSubmit = async (data: LoginSchema) => {
    try {
      await login(data);
    } catch (err) {
      const message = err instanceof Error ? err.message : '';
      if (message.toLowerCase().includes('2fa code required')) {
        setPendingCredentials(data);
      }
    }
  };

  const handlePasskeySignIn = async () => {
    try {
      const result = await passkeyLogin();
      if (result?.requires2FA && result.assertion) {
        setPendingPasskeyAssertion(result.assertion);
      }
    } catch {
      // Error handling is handled in usePasskeyLogin
    }
  };

  const onTOTPSubmit = async (data: TOTPSchema) => {
    if (pendingPasskeyAssertion) {
      try {
        await passkeyLogin({ totpCode: data.totpCode, assertion: pendingPasskeyAssertion });
      } catch {
        totpForm.setError('totpCode', { message: 'Invalid verification code' });
      }
      return;
    }

    if (!pendingCredentials) return;
    try {
      await login({ ...pendingCredentials, totpCode: data.totpCode });
    } catch {
      totpForm.setError('totpCode', { message: 'Invalid verification code' });
    }
  };

  if (pendingCredentials || pendingPasskeyAssertion) {
    return (
      <form onSubmit={totpForm.handleSubmit(onTOTPSubmit)} className="space-y-5">
        <div className="space-y-1">
          <p className="font-medium text-foreground text-sm">Two-factor authentication</p>
          <p className="text-muted-foreground text-xs">
            Enter the code from your authenticator app.
          </p>
        </div>
        <div className="space-y-1.5">
          <Label htmlFor="totpCode">Verification Code</Label>
          <div className="group relative">
            <div className="absolute top-1/2 left-3.5 -translate-y-1/2 text-muted-foreground transition-colors group-focus-within:text-primary">
              <KeyRound className="h-4 w-4" />
            </div>
            <Input
              id="totpCode"
              type="text"
              inputMode="numeric"
              autoComplete="one-time-code"
              placeholder="000000"
              maxLength={8}
              className="pl-10 text-center tracking-[0.28em]"
              {...totpForm.register('totpCode')}
            />
          </div>
          {totpForm.formState.errors.totpCode && (
            <p className="pl-1 text-destructive text-xs">
              {totpForm.formState.errors.totpCode.message}
            </p>
          )}
        </div>
        <div className="flex gap-2">
          <Button
            type="button"
            variant="outline"
            onClick={() => {
              setPendingCredentials(null);
              setPendingPasskeyAssertion(null);
            }}
            className="flex-1"
          >
            Back
          </Button>
          <Button type="submit" disabled={isPending || isPasskeyPending} className="flex-1">
            {isPending || isPasskeyPending ? 'Verifying...' : 'Verify'}
          </Button>
        </div>
      </form>
    );
  }

  return (
    <form onSubmit={loginForm.handleSubmit(onLoginSubmit)} className="space-y-5">
      <div className="space-y-1.5">
        <Label htmlFor="email">Email</Label>
        <div className="group relative">
          <div className="absolute top-1/2 left-3.5 -translate-y-1/2 text-muted-foreground transition-colors group-focus-within:text-primary">
            <Mail className="h-4 w-4" />
          </div>
          <Input
            id="email"
            type="email"
            placeholder="name@example.com"
            className="pl-10"
            {...loginForm.register('email')}
          />
        </div>
        {loginForm.formState.errors.email && (
          <p className="pl-1 text-destructive text-xs">
            {loginForm.formState.errors.email.message}
          </p>
        )}
      </div>

      <div className="space-y-1.5">
        <div className="flex items-center justify-between">
          <Label htmlFor="password">Password</Label>
          <Link
            to="/forgot-password"
            className="text-primary text-sm underline-offset-4 hover:underline"
          >
            Forgot password?
          </Link>
        </div>
        <div className="group relative">
          <div className="absolute top-1/2 left-3.5 -translate-y-1/2 text-muted-foreground transition-colors group-focus-within:text-primary">
            <Lock className="h-4 w-4" />
          </div>
          <Input
            id="password"
            type={showPassword ? 'text' : 'password'}
            placeholder="Enter your password"
            className="pr-10 pl-10"
            {...loginForm.register('password')}
          />
          <button
            type="button"
            onClick={() => setShowPassword(!showPassword)}
            className="absolute top-1/2 right-3.5 -translate-y-1/2 text-muted-foreground transition-colors hover:text-primary"
          >
            {showPassword ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
          </button>
        </div>
        {loginForm.formState.errors.password && (
          <p className="pl-1 text-destructive text-xs">
            {loginForm.formState.errors.password.message}
          </p>
        )}
      </div>

      <Button type="submit" disabled={isPending || isPasskeyPending} className="w-full">
        {isPending ? 'Signing in...' : 'Sign In'}
      </Button>

      {webAuthnSupported && (
        <>
          <div className="relative my-4">
            <div className="absolute inset-0 flex items-center">
              <span className="w-full border-t border-border/60" />
            </div>
            <div className="relative flex justify-center text-xs uppercase">
              <span className="bg-card px-2 font-mono text-muted-foreground text-xs">Or</span>
            </div>
          </div>

          <Button
            type="button"
            variant="outline"
            disabled={isPending || isPasskeyPending}
            onClick={handlePasskeySignIn}
            className="w-full gap-2"
          >
            {isPasskeyPending ? (
              <Loader2 className="h-4 w-4 animate-spin" />
            ) : (
              <Fingerprint className="h-4 w-4 text-primary" />
            )}
            Sign in with a passkey
          </Button>
        </>
      )}
    </form>
  );
};
