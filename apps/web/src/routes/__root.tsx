import { createRootRoute, HeadContent, Scripts } from "@tanstack/react-router";
import type { ReactNode } from "react";
import { Footer } from "../components/footer";
import { Header } from "../components/header";
import { PosthogTracking } from "../components/posthog-tracking";
import appCss from "../styles/global.css?url";

export const Route = createRootRoute({
  head: () => ({
    meta: [
      { charSet: "utf-8" },
      { name: "viewport", content: "width=device-width, initial-scale=1" },
      { title: "Codedock — Self-hosting with superpowers" },
      {
        name: "description",
        content:
          "Codedock is a sub-30MB daemon that turns any Linux server into a full deployment platform — apps, databases, backups, domains, and a fleet control plane.",
      },
    ],
    links: [
      { rel: "stylesheet", href: appCss },
      { rel: "icon", type: "image/svg+xml", href: "/favicon.svg" },
      { rel: "icon", href: "/favicon.ico" },
    ],
  }),
  shellComponent: RootDocument,
});

function RootDocument({ children }: { children: ReactNode }) {
  return (
    <html lang="en" className="dark">
      <head>
        <HeadContent />
      </head>
      <body className="bg-background text-foreground font-sans antialiased overflow-x-hidden max-w-[100vw]">
        <PosthogTracking />
        <div className="pointer-events-none fixed inset-0 bg-[radial-gradient(#6d28d9_1.5px,transparent_1.5px)] bg-size-[40px_40px] opacity-30 -z-10" />
        <Header />
        <main>{children}</main>
        <Footer />
        <Scripts />
      </body>
    </html>
  );
}
