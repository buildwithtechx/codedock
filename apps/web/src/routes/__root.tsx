import { TanStackDevtools } from "@tanstack/react-devtools";
import type { QueryClient } from "@tanstack/react-query";
import { createRootRouteWithContext, HeadContent, Scripts } from "@tanstack/react-router";
import { TanStackRouterDevtoolsPanel } from "@tanstack/react-router-devtools";
import type { ReactNode } from "react";
import { Footer } from "../components/main/footer";
import { Header } from "../components/main/header";
import { PosthogTracking } from "../components/posthog-tracking";
import { ThemeProvider } from "../components/theme-provider";
import { Toaster } from "../components/ui/sonner";
import { TooltipProvider } from "../components/ui/tooltip";
import { createMeta, globalLinks } from "../lib/seo";
import appCss from "../styles.css?url";

interface MyRouterContext {
  queryClient: QueryClient;
}

export const Route = createRootRouteWithContext<MyRouterContext>()({
  head: () => ({
    meta: [
      { charSet: "utf-8" },
      { name: "viewport", content: "width=device-width, initial-scale=1" },
      ...createMeta(),
    ],
    links: [{ rel: "stylesheet", href: appCss }, ...globalLinks],
  }),
  shellComponent: RootDocument,
});

function RootDocument({ children }: { children: ReactNode }) {
  return (
    <html lang="en" suppressHydrationWarning>
      <head>
        <HeadContent />
      </head>
      <body className="bg-background text-foreground font-sans antialiased overflow-x-hidden min-h-screen">
        <PosthogTracking />
        <ThemeProvider defaultTheme="dark">
          <TooltipProvider>
            <div className="pointer-events-none fixed inset-0 bg-[radial-gradient(#6d28d9_1.5px,transparent_1.5px)] bg-size-[40px_40px] opacity-30 -z-10" />
            <Header />
            <main>{children}</main>
            <Footer />
            <Toaster richColors closeButton />
          </TooltipProvider>

          <TanStackDevtools
            config={{
              position: "bottom-right",
            }}
            plugins={[
              {
                name: "Tanstack Router",
                render: <TanStackRouterDevtoolsPanel />,
              },
            ]}
          />
        </ThemeProvider>
        <Scripts />
      </body>
    </html>
  );
}
