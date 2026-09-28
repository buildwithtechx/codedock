import { Link } from '@tanstack/react-router';
import { Mail } from 'lucide-react';
import { useState } from 'react';
import { Button } from '#/components/ui/button';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import { useForgotPassword } from '#/features/auth';

export const ForgotPasswordForm = () => {
  const [email, setEmail] = useState('');
  const { mutate, isPending, isSuccess } = useForgotPassword();

  const handleSubmit = (e: React.SyntheticEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (!email) return;
    mutate(email);
  };

  if (isSuccess) {
    return (
      <div className="space-y-4">
        <div className="flex size-11 items-center justify-center rounded-lg bg-primary-soft text-primary">
          <Mail className="h-5 w-5" />
        </div>
        <p className="font-bold text-2xl text-foreground tracking-[-0.04em]">Check your email.</p>
        <p className="text-muted-foreground text-sm leading-6">
          If an account with that email exists, we've sent you instructions to reset your password.
        </p>
        <div className="mt-6">
          <Link to="/signin" className="text-primary text-sm underline-offset-4 hover:underline">
            Back to sign in
          </Link>
        </div>
      </div>
    );
  }

  return (
    <>
      <form onSubmit={handleSubmit} className="space-y-5">
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
              className="h-12 pl-10"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
              disabled={isPending}
            />
          </div>
        </div>

        <Button type="submit" disabled={isPending || !email} className="h-12 w-full">
          {isPending ? 'Sending...' : 'Send Reset Link'}
        </Button>
      </form>

      <div className="mt-6 border-border border-t pt-5 text-muted-foreground text-sm">
        <span>Remember your password? </span>
        <Link to="/signin" className="text-primary underline-offset-4 hover:underline">
          Sign in
        </Link>
      </div>
    </>
  );
};
