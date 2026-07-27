import { createFileRoute } from "@tanstack/react-router";
import { BottomCta } from "../components/bottom-cta";
import { Features } from "../components/features";
import { Hero } from "../components/hero";
import { SocialProof } from "../components/social-proof";
import { StackCarousel } from "../components/stack-carousel";

export const Route = createFileRoute("/")({
  component: IndexComponent,
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
