import { createFileRoute } from '@tanstack/react-router';
import { DetailPage } from '../features/content/detail-page';
import { featureContent } from '../features/content/feature-content';
import { pageHead } from '../lib/seo';

const content = featureContent['application-deployment'];
export const Route = createFileRoute('/features/application-deployment')({
  head: () => pageHead('/features/application-deployment', content.title, content.description),
  component: Page,
});
function Page() {
  return <DetailPage content={content} />;
}
