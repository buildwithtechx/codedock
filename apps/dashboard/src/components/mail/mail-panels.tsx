import { type KeyboardEvent, type ReactNode, useEffect, useRef, useState } from 'react';
import { cn } from '#/lib/utils';

export interface MailPanelsProps {
  mailList: ReactNode;
  readingPane: ReactNode;
  hasThread: boolean;
  className?: string;
  storageKey?: string;
}

export function MailPanels({
  mailList,
  readingPane,
  hasThread,
  className,
  storageKey = 'mail-split-panel-layout',
}: MailPanelsProps) {
  const [listWidth, setListWidth] = useState<number>(() => {
    if (typeof window !== 'undefined') {
      try {
        const saved = localStorage.getItem(storageKey);
        if (saved) {
          const num = Number(saved);
          if (!Number.isNaN(num) && num >= 20 && num <= 60) {
            return num;
          }
        }
      } catch {
        // ignore
      }
    }
    return 35;
  });

  const [isDragging, setIsDragging] = useState(false);
  const containerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!isDragging) return;

    const onMouseMove = (e: MouseEvent) => {
      if (!containerRef.current) return;
      const rect = containerRef.current.getBoundingClientRect();
      const currentX = e.clientX - rect.left;
      const newWidthPercent = Math.min(Math.max((currentX / rect.width) * 100, 20), 60);
      setListWidth(newWidthPercent);
      try {
        localStorage.setItem(storageKey, String(Math.round(newWidthPercent)));
      } catch {
        // ignore
      }
    };

    const onMouseUp = () => {
      setIsDragging(false);
    };

    window.addEventListener('mousemove', onMouseMove);
    window.addEventListener('mouseup', onMouseUp);
    return () => {
      window.removeEventListener('mousemove', onMouseMove);
      window.removeEventListener('mouseup', onMouseUp);
    };
  }, [isDragging, storageKey]);

  const handleKeyDown = (e: KeyboardEvent) => {
    let nextWidth = listWidth;
    if (e.key === 'ArrowLeft') {
      nextWidth = Math.max(listWidth - 2, 20);
    } else if (e.key === 'ArrowRight') {
      nextWidth = Math.min(listWidth + 2, 60);
    } else if (e.key === 'Home') {
      nextWidth = 20;
    } else if (e.key === 'End') {
      nextWidth = 60;
    } else {
      return;
    }
    e.preventDefault();
    setListWidth(nextWidth);
    try {
      localStorage.setItem(storageKey, String(Math.round(nextWidth)));
    } catch {
      // ignore
    }
  };

  return (
    <div
      ref={containerRef}
      className={cn('relative flex h-full w-full select-none overflow-hidden', className)}
      data-testid="mail-panels-container"
    >
      {/* Desktop split view */}
      <div
        className="hidden h-full flex-col overflow-hidden border-border/60 border-r lg:flex"
        style={{ width: `${listWidth}%` }}
        data-panel="mail-list"
        data-panel-size={listWidth.toFixed(1)}
      >
        {mailList}
      </div>

      <div
        role="separator"
        aria-orientation="vertical"
        aria-label="Resize email list and reading pane"
        aria-valuenow={Math.round(listWidth)}
        aria-valuemin={20}
        aria-valuemax={60}
        tabIndex={0}
        title="Drag to resize, or use the arrow keys"
        onMouseDown={() => setIsDragging(true)}
        onKeyDown={handleKeyDown}
        className={cn(
          'hidden w-1.5 cursor-col-resize items-center justify-center transition-colors hover:bg-primary/20 focus:outline-hidden focus-visible:bg-primary/30 lg:flex',
          isDragging && 'bg-primary/40'
        )}
      />

      <div
        className="hidden h-full flex-1 flex-col overflow-hidden lg:flex"
        style={{ width: `${100 - listWidth}%` }}
        data-panel="mail-reader"
        data-panel-size={(100 - listWidth).toFixed(1)}
      >
        {readingPane}
      </div>

      {/* Mobile single pane view */}
      <div className="flex h-full w-full flex-1 overflow-hidden lg:hidden">
        {hasThread ? (
          <div className="h-full w-full overflow-y-auto" data-panel="mail-reader-mobile">
            {readingPane}
          </div>
        ) : (
          <div className="h-full w-full overflow-y-auto" data-panel="mail-list-mobile">
            {mailList}
          </div>
        )}
      </div>
    </div>
  );
}
