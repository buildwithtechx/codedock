import { ChevronDown, Loader2, RefreshCw, Search } from 'lucide-react';
import { Button } from '@codedock/ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@codedock/ui/dropdown-menu';
import type { useRuntimeUpdate } from '#/hooks/use-runtime-update';

export function RuntimeRedeployMenu({
  runtime,
  disabled = false,
  onStandardRedeploy,
}: {
  runtime: ReturnType<typeof useRuntimeUpdate>;
  disabled?: boolean;
  onStandardRedeploy?: () => void;
}) {
  const isBusy = runtime.updating || runtime.checking;

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button
          variant="outline"
          size="icon"
          disabled={disabled || isBusy}
          aria-label="Redeploy options"
          className="size-9 rounded-xl text-muted-foreground hover:text-foreground"
        >
          {isBusy ? (
            <Loader2 className="size-4 animate-spin" />
          ) : (
            <ChevronDown className="size-4" />
          )}
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end" className="w-64 rounded-xl p-1.5 shadow-lg">
        {onStandardRedeploy && (
          <>
            <DropdownMenuItem
              onClick={onStandardRedeploy}
              className="gap-2 rounded-lg text-xs"
            >
              <RefreshCw className="size-3.5" />
              <span>Redeploy (current image)</span>
            </DropdownMenuItem>
            <DropdownMenuSeparator />
          </>
        )}
        <DropdownMenuItem
          onClick={() => void runtime.update()}
          className="gap-2 rounded-lg text-xs"
        >
          <RefreshCw className="size-3.5 text-amber-500" />
          <span>Redeploy with latest base image</span>
        </DropdownMenuItem>
        <DropdownMenuItem
          onClick={() => void runtime.checkUpdates()}
          className="gap-2 rounded-lg text-xs"
        >
          <Search className="size-3.5 text-primary" />
          <span>Check for image updates</span>
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
