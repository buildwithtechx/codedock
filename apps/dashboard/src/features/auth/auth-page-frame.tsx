import type { ReactNode } from 'react';

interface AuthPageFrameProps {
  eyebrow: string;
  title: string;
  description: string;
  children: ReactNode;
}

export function AuthPageFrame({ eyebrow, title, description, children }: AuthPageFrameProps) {
  return (
    <div className="auth-page-enter w-full max-w-lg">
      <div className="mb-8 text-center">
        <p className="mb-3 font-semibold text-[10px] text-primary uppercase tracking-[0.18em]">
          {eyebrow}
        </p>
        <h1 className="font-bold text-2xl text-foreground tracking-[-0.04em]">{title}</h1>
        <p className="mt-3 text-muted-foreground text-sm leading-6">{description}</p>
      </div>
      {children}
    </div>
  );
}
