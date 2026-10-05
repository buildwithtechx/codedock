import { createFileRoute } from '@tanstack/react-router';
import { DetailPage } from '../features/content/detail-page';
import { featureContent } from '../features/content/feature-content';
import { pageHead } from '../lib/seo';

const content = featureContent['ai-deployment'];
export const Route = createFileRoute('/features/ai-deployment')({
  head: () => pageHead('/features/ai-deployment', content.title, content.description),
  component: Page,
});
function Page() {
  return <DetailPage content={content} />;
}
