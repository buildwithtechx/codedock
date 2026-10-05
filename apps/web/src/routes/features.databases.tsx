import { createFileRoute } from '@tanstack/react-router';
import { DetailPage } from '../features/content/detail-page';
import { featureContent } from '../features/content/feature-content';
import { pageHead } from '../lib/seo';

const content = featureContent.databases;
export const Route = createFileRoute('/features/databases')({
  head: () => pageHead('/features/databases', content.title, content.description),
  component: Page,
});
function Page() {
  return <DetailPage content={content} />;
}
