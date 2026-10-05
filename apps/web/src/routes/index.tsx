import { createFileRoute } from '@tanstack/react-router';
import { BottomCta } from '../features/home/bottom-cta';
import { Features } from '../features/home/features';
import { Hero } from '../features/home/hero';
import { HowItWorks } from '../features/home/how-it-works';
import { SocialProof } from '../features/home/social-proof';
import { StackCarousel } from '../features/home/stack-carousel';
import { pageHead } from '../lib/seo';

export const Route = createFileRoute('/')({
  component: IndexComponent,
  head: () =>
    pageHead(
      '/',
      'Ship fast. Own your infrastructure.',
      'Deploy applications, databases and services on your own Linux servers with Codedock, an open-source Go control plane with Docker and Traefik.'
    ),
});

function IndexComponent() {
  return (
    <>
      <Hero />
      <StackCarousel />
      <HowItWorks />
      <Features />
      <SocialProof />
      <BottomCta />
    </>
  );
}
