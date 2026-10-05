import { createFileRoute } from '@tanstack/react-router';
import { ComparisonHub } from '../features/compare/comparison-hub';
import { pageHead } from '../lib/seo';
export const Route = createFileRoute('/vs/')({
  component: ComparisonHub,
  head: () =>
    pageHead(
      '/vs',
      'Compare deployment platforms',
      'Compare Codedock with Coolify, Dokploy, Vercel, Render, Portainer, CapRover and Dokku. Find the workflow that fits your infrastructure.'
    ),
});
