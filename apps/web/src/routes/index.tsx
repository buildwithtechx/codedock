import { createFileRoute } from '@tanstack/react-router';
import { BottomCta } from '../features/home/bottom-cta';
import { Features } from '../features/home/features';
import { Hero } from '../features/home/hero';
import { SocialProof } from '../features/home/social-proof';
import { StackCarousel } from '../features/home/stack-carousel';
import { createMeta } from '../lib/seo';

export const Route = createFileRoute('/')({
  component: IndexComponent,
  head: () => ({
    meta: createMeta({
      title: 'Codedock — Ship fast, own your infrastructure',
      description:
        'Codedock is a sub-30MB daemon that turns any Linux server into a full deployment platform — apps, databases, backups, domains, and a fleet control plane.',
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
