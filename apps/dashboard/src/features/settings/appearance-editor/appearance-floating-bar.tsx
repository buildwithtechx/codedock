import { Button } from '@codedock/ui/button';
import { Check } from 'lucide-react';

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
    <div className="sticky bottom-3 z-10 flex w-full flex-wrap items-center justify-between gap-3 rounded-2xl border border-border/80 bg-popover/90 px-4 py-3 shadow-lg backdrop-blur-xl">
      <div className="flex items-center gap-2">
        {isDirty ? (
          <>
            <span className="size-2 animate-pulse rounded-full bg-amber-500" />
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
          className="text-muted-foreground text-xs hover:text-foreground"
        >
          Reset
        </Button>

        <Button
          type="button"
          variant="ghost"
          size="sm"
          onClick={onDiscard}
          disabled={!isDirty || isSaving}
          className="text-muted-foreground text-xs hover:text-foreground"
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
