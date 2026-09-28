import { CircleAlert, RefreshCw } from 'lucide-react';
import { cn } from '#/lib/utils';
import { Button } from './button';

type QueryErrorStateProps = {
  title: string;
  description: string;
  onRetry: () => void;
  className?: string;
};

export function QueryErrorState({ title, description, onRetry, className }: QueryErrorStateProps) {
  return (
    <section
      className={cn(
        'flex min-h-[25rem] flex-col items-center justify-center px-6 py-12 text-center',
        className
      )}
    >
      <div className="relative flex h-14 w-14 items-center justify-center rounded-2xl border border-destructive/25 bg-destructive/10 text-destructive shadow-[0_0_0_9px_color-mix(in_oklch,var(--destructive)_8%,transparent)]">
        <CircleAlert className="h-6 w-6" />
      </div>
      <h2 className="mt-5 font-semibold text-xl tracking-tight">{title}</h2>
      <p className="mt-2 max-w-sm text-muted-foreground text-sm leading-6">{description}</p>
      <Button variant="outline" className="mt-6 gap-2" onClick={onRetry}>
        <RefreshCw className="h-4 w-4" />
        Try again
      </Button>
    </section>
  );
}
