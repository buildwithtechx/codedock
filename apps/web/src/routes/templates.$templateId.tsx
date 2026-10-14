import { createFileRoute } from '@tanstack/react-router';
import { TemplateDetail } from '../features/templates/template-detail';
import { pageHead } from '../lib/seo';

export const Route = createFileRoute('/templates/$templateId')({
  component: TemplateDetailPage,
  head: () =>
    pageHead(
      '/templates',
      'App template details',
      'Services, secrets, and setup guide for a Codedock one-click application stack.'
    ),
});

function TemplateDetailPage() {
  const { templateId } = Route.useParams();
  return <TemplateDetail templateId={templateId} />;
}
