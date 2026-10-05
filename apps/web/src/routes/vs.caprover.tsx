import { createFileRoute } from '@tanstack/react-router';
import { ComparisonPage } from '../features/compare/comparison-page';
import { comparisons } from '../features/compare/comparisons';
import { pageHead } from '../lib/seo';

const content = comparisons.caprover;
export const Route = createFileRoute('/vs/caprover')({
  head: () => pageHead('/vs/caprover', `vs ${content.name}`, content.description),
  component: Page,
});
function Page() {
  return <ComparisonPage content={content} />;
}
