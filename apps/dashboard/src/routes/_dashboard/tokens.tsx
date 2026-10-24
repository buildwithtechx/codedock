import { createFileRoute } from '@tanstack/react-router';
import { TokensLayout } from '#/features/tokens/tokens-page';

export const Route = createFileRoute('/_dashboard/tokens')({
  component: TokensPage,
});

export function TokensPage() {
  return <TokensLayout />;
}
