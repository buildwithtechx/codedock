import { Code2 } from 'lucide-react';
import { useState } from 'react';
import type { ExampleApp } from '#/interfaces/templates';
import { cn } from '#/lib/utils';

const GENERIC_LOGO = '/app-logos/_generic.svg';

export function ExampleLogo({ example }: { example: ExampleApp }) {
  const [stage, setStage] = useState(0);

  if (!example.logo || stage > 1) {
    return (
      <span
        aria-hidden="true"
        className="flex size-10 shrink-0 items-center justify-center rounded-xl bg-primary/10 text-primary"
      >
        <Code2 className="size-5" />
      </span>
    );
  }

  const shouldInvert = example.logo
    ? /calcom|directus|ghost|umami|vaultwarden|kafka|buzz|posthog/i.test(example.logo)
    : false;

  return (
    <span className="flex size-10 shrink-0 items-center justify-center overflow-hidden rounded-xl bg-muted/60">
      <img
        src={stage === 0 ? example.logo : GENERIC_LOGO}
        alt=""
        aria-hidden="true"
        className={cn(
          'size-3/5 object-contain',
          shouldInvert && 'dim:brightness-0 dim:invert dark:brightness-0 dark:invert'
        )}
        onError={() => setStage((current) => current + 1)}
      />
    </span>
  );
}
