import { TanStackDevtools } from '@tanstack/react-devtools';
import type { QueryClient } from '@tanstack/react-query';
import { createRootRouteWithContext, HeadContent, Outlet } from '@tanstack/react-router';
import { TanStackRouterDevtoolsPanel } from '@tanstack/react-router-devtools';
import { Footer } from '../components/main/footer';
import { Header } from '../components/main/header';
import { PosthogTracking } from '../components/posthog-tracking';
import { ThemeProvider } from '../components/theme-provider';
import { Toaster } from '../components/ui/sonner';
import { TooltipProvider } from '../components/ui/tooltip';
import { createMeta } from '../lib/seo';

interface MyRouterContext {
  queryClient: QueryClient;
}

export const Route = createRootRouteWithContext<MyRouterContext>()({
  head: () => ({
    meta: [
      { charSet: 'utf-8' },
      { name: 'viewport', content: 'width=device-width, initial-scale=1' },
      ...createMeta(),
    ],
  }),
  component: RootComponent,
});

function RootComponent() {
  return (
    <>
      <HeadContent />
      <PosthogTracking />
      <ThemeProvider defaultTheme="dark">
        <TooltipProvider>
          <div className="pointer-events-none fixed inset-0 -z-10 bg-[radial-gradient(#6d28d9_1.5px,transparent_1.5px)] bg-size-[40px_40px] opacity-30" />
          <a
            href="#main-content"
            className="sr-only focus:not-sr-only focus:absolute focus:z-50 focus:rounded-lg focus:bg-primary focus:p-4 focus:text-white"
          >
            Skip to content
          </a>
          <Header />
          <main id="main-content" tabIndex={-1}>
            <Outlet />
          </main>
          <Footer />
          <Toaster richColors closeButton />
        </TooltipProvider>

        {import.meta.env.DEV && (
          <TanStackDevtools
            config={{
              position: 'bottom-right',
            }}
            plugins={[
              {
                name: 'Tanstack Router',
                render: <TanStackRouterDevtoolsPanel />,
              },
            ]}
          />
        )}
      </ThemeProvider>
    </>
  );
}
