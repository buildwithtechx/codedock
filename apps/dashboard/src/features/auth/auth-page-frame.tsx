import type { ReactNode } from 'react';
import { useSystemStore } from '#/stores/system-store';

interface AuthPageFrameProps {
  eyebrow: string;
  title: string;
  description: string;
  children: ReactNode;
}

export function AuthPageFrame({ eyebrow, title, description, children }: AuthPageFrameProps) {
  const siteName = useSystemStore((state) => state.siteName);

  return (
    <div className="auth-page-enter w-full max-w-lg">
      <div className="mb-10 text-center">
        <div className="mb-5 flex items-center justify-center gap-2.5">
          <img src="/apple-touch-icon.png" alt="" className="size-8 rounded-lg" />
          <span className="font-bold text-foreground text-lg tracking-[-0.04em]">{siteName}</span>
        </div>
        <p className="mb-4 font-semibold text-[10px] text-primary uppercase tracking-[0.18em]">
          {eyebrow}
        </p>
        <h1 className="font-bold text-3xl text-foreground tracking-[-0.045em] sm:text-4xl">
          {title}
        </h1>
        <p className="mt-3 text-muted-foreground text-sm leading-6">{description}</p>
      </div>
      {children}
    </div>
  );
}
