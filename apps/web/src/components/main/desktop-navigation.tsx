import { ChevronDown } from 'lucide-react';
import { NavigationMenu } from 'radix-ui';
import { navigationGroups } from '../../lib/navigation';

export function DesktopNavigation() {
  return (
    <NavigationMenu.Root className="relative hidden lg:block" aria-label="Main navigation">
      <NavigationMenu.List className="flex items-center gap-1">
        {navigationGroups.map((group) => (
          <NavigationMenu.Item key={group.label}>
            <NavigationMenu.Trigger className="group flex items-center gap-1 rounded-lg px-3 py-2 text-muted-foreground text-sm transition-colors hover:bg-muted hover:text-foreground focus-visible:outline-2 focus-visible:outline-primary data-[state=open]:text-foreground">
              {group.label}
              <ChevronDown className="size-3.5 transition-transform group-data-[state=open]:rotate-180" />
            </NavigationMenu.Trigger>
            <NavigationMenu.Content className="absolute top-full left-0 mt-3 w-[min(36rem,80vw)] rounded-2xl border border-border bg-card p-3 shadow-2xl">
              <div className="grid gap-1 sm:grid-cols-2">
                {group.items.map((item) => (
                  <NavigationMenu.Link key={item.href} asChild>
                    <a
                      href={item.href}
                      className="rounded-xl p-4 transition-colors hover:bg-muted focus-visible:outline-2 focus-visible:outline-primary"
                    >
                      <span className="block font-semibold text-sm">{item.label}</span>
                      <span className="mt-1 block text-muted-foreground text-xs leading-relaxed">
                        {item.description}
                      </span>
                    </a>
                  </NavigationMenu.Link>
                ))}
              </div>
            </NavigationMenu.Content>
          </NavigationMenu.Item>
        ))}
        <NavigationMenu.Item>
          <NavigationMenu.Link asChild>
            <a
              href="/pricing"
              className="rounded-lg px-3 py-2 text-muted-foreground text-sm hover:text-foreground"
            >
              Pricing
            </a>
          </NavigationMenu.Link>
        </NavigationMenu.Item>
      </NavigationMenu.List>
    </NavigationMenu.Root>
  );
}
