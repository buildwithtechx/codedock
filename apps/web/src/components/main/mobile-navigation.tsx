import { Menu, X } from 'lucide-react';
import { Dialog } from 'radix-ui';
import { navigationGroups } from '../../lib/navigation';
import { productLinks } from '../../lib/product-links';

export function MobileNavigation() {
  return (
    <Dialog.Root>
      <Dialog.Trigger asChild>
        <button
          type="button"
          aria-label="Open menu"
          className="rounded-lg p-2 hover:bg-muted lg:hidden"
        >
          <Menu className="size-5" />
        </button>
      </Dialog.Trigger>
      <Dialog.Portal>
        <Dialog.Overlay className="fixed inset-0 z-50 bg-black/60 backdrop-blur-sm" />
        <Dialog.Content className="fixed inset-y-0 right-0 z-50 w-[min(25rem,100vw)] overflow-y-auto border-border border-l bg-background p-6">
          <div className="mb-7 flex items-center justify-between">
            <Dialog.Title className="font-bold text-lg">Explore Codedock</Dialog.Title>
            <Dialog.Close asChild>
              <button
                type="button"
                aria-label="Close menu"
                className="rounded-lg p-2 hover:bg-muted"
              >
                <X className="size-5" />
              </button>
            </Dialog.Close>
          </div>
          <Dialog.Description className="sr-only">
            Features, solutions, and resources for deploying on your own infrastructure.
          </Dialog.Description>
          <nav aria-label="Mobile navigation" className="space-y-6">
            {navigationGroups.map((group) => (
              <div key={group.label}>
                <h2 className="mb-2 font-semibold text-primary text-xs uppercase tracking-widest">
                  {group.label}
                </h2>
                <div className="space-y-1">
                  {group.items.map((item) => (
                    <Dialog.Close key={item.href} asChild>
                      <a
                        href={item.href}
                        className="block rounded-lg px-3 py-2 text-sm hover:bg-muted"
                      >
                        {item.label}
                      </a>
                    </Dialog.Close>
                  ))}
                </div>
              </div>
            ))}
            <Dialog.Close asChild>
              <a href="/pricing" className="block rounded-lg px-3 py-2 text-sm hover:bg-muted">
                Pricing
              </a>
            </Dialog.Close>
            <Dialog.Close asChild>
              <a
                href={productLinks.cloud}
                className="block rounded-xl bg-primary px-4 py-3 text-center font-semibold text-sm text-white"
              >
                Open Cloud
              </a>
            </Dialog.Close>
          </nav>
        </Dialog.Content>
      </Dialog.Portal>
    </Dialog.Root>
  );
}
