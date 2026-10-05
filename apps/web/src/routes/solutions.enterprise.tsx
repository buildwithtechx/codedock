import { createFileRoute } from '@tanstack/react-router';
import { DetailPage } from '../features/content/detail-page';
import { solutionContent } from '../features/content/solution-content';
import { pageHead } from '../lib/seo';

const content = solutionContent.enterprise;
export const Route = createFileRoute('/solutions/enterprise')({
  head: () => pageHead('/solutions/enterprise', content.title, content.description),
  component: Page,
});
function Page() {
  return <DetailPage content={content} />;
}
