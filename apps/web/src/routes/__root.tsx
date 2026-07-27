import { TanStackDevtools } from "@tanstack/react-devtools";
import type { QueryClient } from "@tanstack/react-query";
import { createRootRouteWithContext, HeadContent, Scripts } from "@tanstack/react-router";
import { TanStackRouterDevtoolsPanel } from "@tanstack/react-router-devtools";
import type { ReactNode } from "react";
import { Footer } from "../components/footer";
import { Header } from "../components/header";
import { PosthogTracking } from "../components/posthog-tracking";
import { ThemeProvider } from "../components/theme-provider";
import { Toaster } from "../components/ui/sonner";
import { TooltipProvider } from "../components/ui/tooltip";
import { createMeta, globalLinks } from "../lib/seo";
import appCss from "../styles/global.css?url";

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
    <html lang="en" className="dark bg-[#0f0f11]">
      <head>
        <HeadContent />
      </head>
      <body className="bg-[#0f0f11] text-white antialiased">
        <PosthogTracking />
        <ThemeProvider>
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
