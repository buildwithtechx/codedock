import { CircleAlert, RefreshCw } from 'lucide-react';
import { Button } from './button';

type QueryErrorStateProps = {
  title: string;
  description: string;
  onRetry: () => void;
  className?: string;
};

export function QueryErrorState({ title, description, onRetry, className }: QueryErrorStateProps) {
  return (
    <div
      className={`flex min-h-80 flex-col items-center justify-center rounded-2xl border border-destructive/30 bg-card px-6 text-center ${className ?? ''}`}
    >
      <div className="flex h-11 w-11 items-center justify-center rounded-xl bg-destructive/10 text-destructive">
        <CircleAlert className="h-5 w-5" />
      </div>
      <h2 className="mt-4 font-semibold text-lg">{title}</h2>
      <p className="mt-1 max-w-sm text-muted-foreground text-sm">{description}</p>
      <Button variant="outline" className="mt-5 gap-2" onClick={onRetry}>
        <RefreshCw className="h-4 w-4" />
        Try again
      </Button>
    </div>
  );
}
