import { createFileRoute } from '@tanstack/react-router';
import { ComparisonPage } from '../features/compare/comparison-page';
import { comparisons } from '../features/compare/comparisons';
import { pageHead } from '../lib/seo';

const content = comparisons.dokku;
export const Route = createFileRoute('/vs/dokku')({
  head: () => pageHead('/vs/dokku', `vs ${content.name}`, content.description),
  component: Page,
});
function Page() {
  return <ComparisonPage content={content} />;
}
