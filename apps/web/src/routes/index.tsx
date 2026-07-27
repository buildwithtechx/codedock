import { createFileRoute } from "@tanstack/react-router";
import { BottomCta } from "../components/bottom-cta";
import { Features } from "../components/features";
import { Hero } from "../components/hero";
import { SocialProof } from "../components/social-proof";
import { StackCarousel } from "../components/stack-carousel";
import { createMeta } from "../lib/seo";

export const Route = createFileRoute("/")({
  component: IndexComponent,
  head: () => ({
    meta: createMeta({
      title: "Codedock — Ship fast, own your infrastructure",
      description:
        "Codedock is a sub-30MB daemon that turns any Linux server into a full deployment platform — apps, databases, backups, domains, and a fleet control plane.",
    }),
  }),
});

function IndexComponent() {
  return (
    <>
      <Hero />
      <StackCarousel />
      <Features />
      <SocialProof />
      <BottomCta />
    </>
  );
}
