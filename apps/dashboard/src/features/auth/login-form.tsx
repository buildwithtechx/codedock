import { zodResolver } from '@hookform/resolvers/zod';
import { Link } from '@tanstack/react-router';
import { Eye, EyeOff, KeyRound, Lock, Mail } from 'lucide-react';
import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { z } from 'zod';
import { Button } from '#/components/ui/button';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import { useLogin } from '#/features/auth';

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
  const [showPassword, setShowPassword] = useState(false);
  const [pendingCredentials, setPendingCredentials] = useState<LoginSchema | null>(null);

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

  const onTOTPSubmit = async (data: TOTPSchema) => {
    if (!pendingCredentials) return;
    try {
      await login({ ...pendingCredentials, totpCode: data.totpCode });
    } catch {
      totpForm.setError('totpCode', { message: 'Invalid verification code' });
    }
  };

  if (pendingCredentials) {
    return (
      <form onSubmit={totpForm.handleSubmit(onTOTPSubmit)} className="space-y-4">
        <div className="space-y-1 text-center">
          <p className="font-semibold text-foreground text-sm">Two-Factor Authentication</p>
          <p className="text-muted-foreground text-xs">
            Enter the 6-digit code from your authenticator app.
          </p>
        </div>
        <div className="space-y-1.5">
          <Label htmlFor="totpCode" className="font-medium text-foreground/90 text-sm">
            Verification Code
          </Label>
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
              className="h-11 rounded-xl border bg-background/80 pl-10 text-center text-sm tracking-widest transition-all duration-300 focus:border-primary/50 focus:ring-2 focus:ring-primary/20"
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
            onClick={() => setPendingCredentials(null)}
            className="h-11 flex-1 rounded-xl text-sm"
          >
            Back
          </Button>
          <Button
            type="submit"
            disabled={isPending}
            className="h-11 flex-1 rounded-xl bg-linear-to-r from-primary to-purple-600 font-semibold text-sm shadow-lg shadow-primary/30 transition-all duration-200 hover:brightness-110 active:scale-[0.985]"
          >
            {isPending ? 'Verifying...' : 'Verify'}
          </Button>
        </div>
      </form>
    );
  }

  return (
    <form onSubmit={loginForm.handleSubmit(onLoginSubmit)} className="space-y-4">
      <div className="space-y-1.5">
        <Label htmlFor="email" className="font-medium text-foreground/90 text-sm">
          Email
        </Label>
        <div className="group relative">
          <div className="absolute top-1/2 left-3.5 -translate-y-1/2 text-muted-foreground transition-colors group-focus-within:text-primary">
            <Mail className="h-4 w-4" />
          </div>
          <Input
            id="email"
            type="email"
            placeholder="name@example.com"
            className="h-11 rounded-xl border bg-background/80 pl-10 text-sm transition-all duration-300 focus:border-primary/50 focus:ring-2 focus:ring-primary/20"
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
          <Label htmlFor="password" className="font-medium text-foreground/90 text-sm">
            Password
          </Label>
          <Link
            to="/forgot-password"
            className="text-primary text-sm transition-colors hover:text-primary-hover"
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
            className="h-11 rounded-xl border bg-background/80 pr-10 pl-10 text-sm transition-all duration-300 focus:border-primary/50 focus:ring-2 focus:ring-primary/20"
            {...loginForm.register('password')}
          />
          <button
            type="button"
            onClick={() => setShowPassword(!showPassword)}
            className="absolute top-1/2 right-3.5 -translate-y-1/2 text-muted-foreground transition-colors hover:text-foreground"
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

      <Button
        type="submit"
        disabled={isPending}
        className="h-11 w-full rounded-xl bg-linear-to-r from-primary to-purple-600 font-semibold text-sm shadow-lg shadow-primary/30 transition-all duration-200 hover:brightness-110 active:scale-[0.985]"
      >
        {isPending ? 'Signing in...' : 'Sign In'}
      </Button>
    </form>
  );
};
