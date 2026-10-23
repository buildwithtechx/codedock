import { Check } from 'lucide-react';
import { Button } from '@codedock/ui/button';

interface AppearanceFloatingBarProps {
  isDirty: boolean;
  isSaving: boolean;
  onReset: () => void;
  onDiscard: () => void;
  onSave: () => void;
}

export function AppearanceFloatingBar({
  isDirty,
  isSaving,
  onReset,
  onDiscard,
  onSave,
}: AppearanceFloatingBarProps) {
  return (
    <div className="fixed bottom-6 left-1/2 z-40 flex w-[calc(100%-2rem)] max-w-3xl -translate-x-1/2 items-center justify-between rounded-2xl border border-border/80 bg-card/95 px-5 py-3 shadow-2xl backdrop-blur-md">
      <div className="flex items-center gap-2">
        {isDirty ? (
          <>
            <span className="size-2 rounded-full bg-amber-500 animate-pulse" />
            <span className="font-medium text-amber-500 text-xs">Unsaved changes</span>
          </>
        ) : (
          <span className="font-medium text-muted-foreground text-xs">Saved</span>
        )}
      </div>

      <div className="flex items-center gap-2">
        <Button
          type="button"
          variant="ghost"
          size="sm"
          onClick={onReset}
          className="text-muted-foreground hover:text-foreground text-xs"
        >
          Reset
        </Button>

        <Button
          type="button"
          variant="ghost"
          size="sm"
          onClick={onDiscard}
          disabled={!isDirty || isSaving}
          className="text-muted-foreground hover:text-foreground text-xs"
        >
          Discard
        </Button>

        <Button
          type="button"
          size="sm"
          onClick={onSave}
          disabled={isSaving}
          className="gap-1.5 rounded-xl px-4 font-medium text-xs shadow-md"
        >
          <Check className="size-3.5" />
          <span>{isSaving ? 'Saving...' : 'Save'}</span>
        </Button>
      </div>
    </div>
  );
}
