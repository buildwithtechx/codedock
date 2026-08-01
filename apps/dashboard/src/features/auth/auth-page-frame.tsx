import type { ReactNode } from 'react';

interface AuthPageFrameProps {
  eyebrow: string;
  title: string;
  description: string;
  children: ReactNode;
}

export function AuthPageFrame({ eyebrow, title, description, children }: AuthPageFrameProps) {
  return (
    <div className="auth-page-enter w-full max-w-md">
      <div className="mb-8">
        <p className="mb-4 font-semibold text-[10px] text-primary uppercase tracking-[0.18em]">
          {eyebrow}
        </p>
        <h1 className="max-w-md font-bold text-3xl text-foreground tracking-[-0.045em] sm:text-4xl">
          {title}
        </h1>
        <p className="mt-3 max-w-md text-muted-foreground text-sm leading-6">{description}</p>
      </div>
      {children}
    </div>
  );
}
