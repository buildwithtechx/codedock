import {
  ChevronUp,
  Maximize2,
  Minimize2,
  ScrollText,
  Terminal as TerminalIcon,
  X,
} from 'lucide-react';
import { Suspense, lazy, useEffect, useState } from 'react';
import { Button } from '@codedock/ui/button';
import { useDockStore } from '#/stores/dock-store';

const LiveLogsViewer = lazy(() =>
  import('#/features/services/live-logs-viewer').then((m) => ({ default: m.LiveLogsViewer }))
);

const WebTerminal = lazy(() =>
  import('#/features/services/web-terminal').then((m) => ({ default: m.WebTerminal }))
);

export function BottomDock() {
  const {
    isOpen,
    activeTab,
    serviceId,
    deploymentId,
    title,
    height,
    closeDock,
    toggleDock,
    setActiveTab,
    setHeight,
  } = useDockStore();
  const [isMaximized, setIsMaximized] = useState(false);
  const [isDragging, setIsDragging] = useState(false);

  useEffect(() => {
    const handleKeyDown = (event: KeyboardEvent) => {
      if (event.key.toLowerCase() === 'j' && (event.metaKey || event.ctrlKey)) {
        event.preventDefault();
        toggleDock();
      }
    };
    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [toggleDock]);

  const handleMouseDown = (e: React.MouseEvent) => {
    e.preventDefault();
    setIsDragging(true);
    const startY = e.clientY;
    const startHeight = height;

    const handleMouseMove = (moveEvent: MouseEvent) => {
      const delta = startY - moveEvent.clientY;
      const nextHeight = Math.max(160, Math.min(window.innerHeight - 100, startHeight + delta));
      setHeight(nextHeight);
    };

    const handleMouseUp = () => {
      setIsDragging(false);
      window.removeEventListener('mousemove', handleMouseMove);
      window.removeEventListener('mouseup', handleMouseUp);
    };

    window.addEventListener('mousemove', handleMouseMove);
    window.addEventListener('mouseup', handleMouseUp);
  };

  if (!isOpen) {
    return (
      <div className="fixed right-6 bottom-4 z-40">
        <Button
          size="sm"
          variant="outline"
          onClick={toggleDock}
          className="flex h-8 items-center gap-2 rounded-full border-border/80 bg-background/90 px-3 text-xs shadow-lg backdrop-blur-md hover:bg-muted"
        >
          <TerminalIcon className="size-3.5 text-primary" />
          <span className="font-medium">Logs & Terminal</span>
          <kbd className="rounded border bg-muted px-1 py-0.2 font-mono text-[9px] text-muted-foreground">
            ⌘J
          </kbd>
          <ChevronUp className="size-3.5 text-muted-foreground" />
        </Button>
      </div>
    );
  }

  const currentHeight = isMaximized ? 'calc(100vh - 4rem)' : `${height}px`;

  return (
    <div
      style={{ height: currentHeight }}
      className={`fixed inset-x-0 bottom-0 z-40 flex flex-col border-border/80 border-t bg-[#090a0c] text-neutral-100 shadow-2xl transition-[height] ${
        isDragging ? 'duration-0' : 'duration-150'
      }`}
    >
      <div
        onMouseDown={handleMouseDown}
        className="group relative flex h-2 w-full cursor-ns-resize items-center justify-center bg-transparent transition-colors hover:bg-primary/20"
      >
        <div className="h-1 w-12 rounded-full bg-border/60 transition-colors group-hover:bg-primary" />
      </div>

      <div className="flex h-10 shrink-0 items-center justify-between border-border/60 border-b bg-[#0d0f12] px-3">
        <div className="flex items-center gap-1">
          <button
            type="button"
            onClick={() => setActiveTab('logs')}
            className={`flex items-center gap-2 rounded-md px-2.5 py-1 text-xs font-medium transition-colors ${
              activeTab === 'logs'
                ? 'bg-neutral-800 text-white'
                : 'text-neutral-400 hover:bg-neutral-800/50 hover:text-neutral-200'
            }`}
          >
            <ScrollText className="size-3.5" />
            <span>Logs</span>
            {serviceId && <span className="size-1.5 rounded-full bg-emerald-500 animate-pulse" />}
          </button>
          <button
            type="button"
            onClick={() => setActiveTab('terminal')}
            className={`flex items-center gap-2 rounded-md px-2.5 py-1 text-xs font-medium transition-colors ${
              activeTab === 'terminal'
                ? 'bg-neutral-800 text-white'
                : 'text-neutral-400 hover:bg-neutral-800/50 hover:text-neutral-200'
            }`}
          >
            <TerminalIcon className="size-3.5" />
            <span>Terminal</span>
          </button>

          {title && (
            <span className="ml-2 border-border/60 border-l pl-3 text-neutral-400 text-xs">
              {title}
            </span>
          )}
        </div>

        <div className="flex items-center gap-1 text-neutral-400">
          <button
            type="button"
            onClick={() => setIsMaximized((prev) => !prev)}
            className="flex size-7 items-center justify-center rounded-md hover:bg-neutral-800 hover:text-white"
            aria-label={isMaximized ? 'Restore size' : 'Maximize'}
          >
            {isMaximized ? <Minimize2 className="size-3.5" /> : <Maximize2 className="size-3.5" />}
          </button>
          <button
            type="button"
            onClick={closeDock}
            className="flex size-7 items-center justify-center rounded-md hover:bg-neutral-800 hover:text-white"
            aria-label="Close dock (⌘J)"
          >
            <X className="size-3.5" />
          </button>
        </div>
      </div>

      <div className="relative min-h-0 flex-1 overflow-hidden bg-[#07080a]">
        <Suspense
          fallback={
            <div className="flex h-full items-center justify-center text-neutral-500 text-xs">
              Loading runtime session...
            </div>
          }
        >
          {activeTab === 'logs' ? (
            serviceId ? (
              <LiveLogsViewer serviceId={serviceId} deploymentId={deploymentId} />
            ) : (
              <div className="flex h-full flex-col items-center justify-center gap-2 text-center text-neutral-400 text-xs">
                <ScrollText className="size-8 text-neutral-600" />
                <p className="font-medium text-neutral-300">
                  No active service selected for live logs
                </p>
                <p className="max-w-md text-neutral-500">
                  Open any service or active deployment to view real-time log output and stream
                  container events.
                </p>
              </div>
            )
          ) : serviceId ? (
            <WebTerminal serviceId={serviceId} />
          ) : (
            <div className="flex h-full flex-col items-center justify-center gap-2 text-center text-neutral-400 text-xs">
              <TerminalIcon className="size-8 text-neutral-600" />
              <p className="font-medium text-neutral-300">Terminal console ready</p>
              <p className="max-w-md text-neutral-500">
                Navigate to a service container or server to start an interactive shell session.
              </p>
            </div>
          )}
        </Suspense>
      </div>
    </div>
  );
}
