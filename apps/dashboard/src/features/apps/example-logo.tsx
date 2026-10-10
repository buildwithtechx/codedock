import { Code2 } from 'lucide-react';
import { useState } from 'react';
import type { ExampleApp } from '#/interfaces/templates';

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

  return (
    <span className="flex size-10 shrink-0 items-center justify-center overflow-hidden rounded-xl bg-muted/60">
      <img
        src={stage === 0 ? example.logo : GENERIC_LOGO}
        alt=""
        aria-hidden="true"
        className="size-3/5 object-contain"
        onError={() => setStage((current) => current + 1)}
      />
    </span>
  );
}
