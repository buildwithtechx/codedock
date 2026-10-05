import { Menu, X } from 'lucide-react';
import { Dialog } from 'radix-ui';
import { useState } from 'react';
import { navigationGroups } from '../../lib/navigation';
import { productLinks } from '../../lib/product-links';
import { SiteLink } from '../site-link';

export function MobileNavigation() {
  const [open, setOpen] = useState(false);
  return (
    <Dialog.Root open={open} onOpenChange={setOpen}>
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
                    <SiteLink
                      key={item.href}
                      href={item.href}
                      onClick={() => setOpen(false)}
                      className="block rounded-lg px-3 py-2 text-sm hover:bg-muted"
                    >
                      {item.label}
                    </SiteLink>
                  ))}
                </div>
              </div>
            ))}
            <SiteLink
              href="/pricing"
              onClick={() => setOpen(false)}
              className="block rounded-lg px-3 py-2 text-sm hover:bg-muted"
            >
              Pricing
            </SiteLink>
            <SiteLink
              href={productLinks.cloud}
              onClick={() => setOpen(false)}
              className="block rounded-xl bg-primary px-4 py-3 text-center font-semibold text-sm text-white"
            >
              Open Cloud
            </SiteLink>
          </nav>
        </Dialog.Content>
      </Dialog.Portal>
    </Dialog.Root>
  );
}
