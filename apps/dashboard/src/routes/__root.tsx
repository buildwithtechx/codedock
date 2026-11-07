import type { QueryClient } from '@tanstack/react-query';
import { createRootRouteWithContext, HeadContent, Outlet } from '@tanstack/react-router';
import { BackendGate } from '#/components/backend-gate';
import { ThemeProvider, Toaster, TooltipProvider } from '@codedock/ui';
import { PostHogProvider } from '#/integrations/posthog-provider';

interface MyRouterContext {
  queryClient: QueryClient;
}

export const Route = createRootRouteWithContext<MyRouterContext>()({
  head: () => ({
    meta: [
      {
        name: 'description',
        content: 'Codedock - Deploy apps instantly',
      },
      { title: 'Codedock' },
    ],
  }),
  component: RootDocument,
});

function RootDocument() {
  return (
    <ThemeProvider>
      <PostHogProvider>
        <HeadContent />
        <TooltipProvider>
          <BackendGate>
            <Outlet />
          </BackendGate>
        </TooltipProvider>
        <Toaster />
      </PostHogProvider>
    </ThemeProvider>
  );
}
