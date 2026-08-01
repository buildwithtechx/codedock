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
        <p className="font-extrabold text-2xl text-foreground tracking-[-0.05em]">
          Check your email.
        </p>
        <p className="text-muted-foreground text-sm leading-6">
          If an account with that email exists, we've sent you instructions to reset your password.
        </p>
        <div className="mt-6">
          <Link
            to="/signin"
            className="font-bold text-primary text-sm underline decoration-primary/35 underline-offset-4"
          >
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
          <Label htmlFor="email" className="font-bold text-foreground text-sm">
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
              className="h-12 rounded-lg border-border bg-background pl-10 text-sm shadow-none transition-colors placeholder:text-muted-foreground focus-visible:border-primary/50 focus-visible:ring-2 focus-visible:ring-primary/20"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
              disabled={isPending}
            />
          </div>
        </div>

        <Button
          type="submit"
          disabled={isPending || !email}
          className="h-12 w-full rounded-lg bg-primary font-bold text-primary-foreground text-sm shadow-none hover:bg-primary-hover"
        >
          {isPending ? 'Sending...' : 'Send Reset Link'}
        </Button>
      </form>

      <div className="mt-6 border-border border-t pt-5 text-muted-foreground text-sm">
        <span>Remember your password? </span>
        <Link
          to="/signin"
          className="font-bold text-primary underline decoration-primary/35 underline-offset-4"
        >
          Sign in
        </Link>
      </div>
    </>
  );
};
