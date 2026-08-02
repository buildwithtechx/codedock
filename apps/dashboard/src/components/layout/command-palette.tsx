import { useNavigate } from '@tanstack/react-router';
import { Command } from 'cmdk';
import { CornerDownLeft, FolderKanban, Plus, Search, Server, Sparkles } from 'lucide-react';
import { useEffect } from 'react';
import { Dialog, DialogContent, DialogDescription, DialogTitle } from '#/components/ui/dialog';
import { commandNavigation } from './dashboard-navigation';

interface CommandPaletteProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

const creationItems = [
  {
    title: 'New project',
    description: 'Create a project workspace',
    to: '/projects/new',
    icon: FolderKanban,
  },
  {
    title: 'New app',
    description: 'Choose a project for a new app',
    to: '/apps/new',
    icon: Sparkles,
  },
  {
    title: 'Add server',
    description: 'Register a deployment target',
    to: '/servers/new',
    icon: Server,
  },
  {
    title: 'Add backup destination',
    description: 'Connect compatible object storage',
    to: '/backups',
    icon: Plus,
  },
];

export function CommandPalette({ open, onOpenChange }: CommandPaletteProps) {
  const navigate = useNavigate();

  useEffect(() => {
    const handleKeyDown = (event: KeyboardEvent) => {
      if (event.key.toLowerCase() !== 'k') {
        return;
      }

      if (!event.metaKey && !event.ctrlKey) {
        return;
      }

      event.preventDefault();
      onOpenChange(!open);
    };

    window.addEventListener('keydown', handleKeyDown);

    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [onOpenChange, open]);

  const navigateTo = (to: string, search?: Record<string, string>) => {
    navigate({ to: to as never, search: search as never });
    onOpenChange(false);
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent
        className="top-[12%] w-[min(920px,calc(100vw-2rem))] max-w-none translate-y-0 gap-0 overflow-hidden rounded-xl p-0 shadow-2xl sm:max-w-none"
        showCloseButton={false}
      >
        <DialogTitle className="sr-only">Command palette</DialogTitle>
        <DialogDescription className="sr-only">
          Search navigation and run dashboard commands.
        </DialogDescription>
        <Command className="bg-popover text-popover-foreground">
          <div className="flex items-center gap-3 border-b px-5 py-4">
            <Search className="size-4 text-muted-foreground" />
            <Command.Input
              autoFocus
              placeholder="Search pages, resources, and actions"
              className="h-9 flex-1 bg-transparent text-base outline-none placeholder:text-muted-foreground"
            />
            <kbd className="rounded border bg-muted px-1.5 py-0.5 font-mono text-[10px] text-muted-foreground">
              Esc
            </kbd>
          </div>

          <Command.List className="scrollbar-none max-h-140 overflow-y-auto p-3 [&::-webkit-scrollbar]:hidden">
            <Command.Empty className="px-3 py-8 text-center text-muted-foreground text-sm">
              No matching command.
            </Command.Empty>

            <Command.Group
              heading="Navigate"
              className="px-1 py-2 text-muted-foreground text-xs **:[[cmdk-group-heading]]:px-2 **:[[cmdk-group-heading]]:pb-2"
            >
              <div className="grid gap-1 md:grid-cols-2">
                {commandNavigation.map((item) => (
                  <Command.Item
                    key={item.to}
                    value={`${item.title} ${item.description}`}
                    onSelect={() => navigateTo(item.to, item.search)}
                    className="flex cursor-pointer items-center gap-3 rounded-lg border border-transparent px-3 py-3 text-foreground text-sm data-[selected=true]:border-border data-[selected=true]:bg-muted/70"
                  >
                    <div className="flex size-9 items-center justify-center rounded-md border bg-background">
                      <item.icon className="size-4 text-muted-foreground" />
                    </div>
                    <div className="min-w-0 flex-1">
                      <p className="font-medium">{item.title}</p>
                      <p className="truncate text-muted-foreground text-xs">{item.description}</p>
                    </div>
                  </Command.Item>
                ))}
              </div>
            </Command.Group>

            <Command.Separator className="my-2 h-px bg-border" />

            <Command.Group
              heading="Create"
              className="px-1 py-2 text-muted-foreground text-xs **:[[cmdk-group-heading]]:px-2 **:[[cmdk-group-heading]]:pb-2"
            >
              <div className="grid gap-1 md:grid-cols-2">
                {creationItems.map((item) => (
                  <Command.Item
                    key={item.to}
                    value={`${item.title} ${item.description}`}
                    onSelect={() => navigateTo(item.to)}
                    className="flex cursor-pointer items-center gap-3 rounded-lg border border-transparent px-3 py-3 text-foreground text-sm data-[selected=true]:border-border data-[selected=true]:bg-muted/70"
                  >
                    <div className="flex size-9 items-center justify-center rounded-md border bg-background">
                      <item.icon className="size-4 text-muted-foreground" />
                    </div>
                    <div className="min-w-0 flex-1">
                      <p className="font-medium">{item.title}</p>
                      <p className="truncate text-muted-foreground text-xs">{item.description}</p>
                    </div>
                  </Command.Item>
                ))}
              </div>
            </Command.Group>
          </Command.List>

          <div className="flex items-center justify-between border-t bg-muted/30 px-5 py-3 text-muted-foreground text-xs">
            <div className="flex items-center gap-2">
              <kbd className="rounded border bg-background px-1.5 py-0.5 font-mono text-[10px]">
                ↑↓
              </kbd>
              <span>Navigate</span>
            </div>
            <div className="flex items-center gap-2">
              <CornerDownLeft className="size-3.5" />
              <span>Open selected command</span>
            </div>
          </div>
        </Command>
      </DialogContent>
    </Dialog>
  );
}
