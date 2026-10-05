import { createFileRoute } from '@tanstack/react-router';
import { TemplateCatalog } from '../features/templates/template-catalog';
import { pageHead } from '../lib/seo';
export const Route = createFileRoute('/templates')({
  component: TemplateCatalog,
  head: () =>
    pageHead(
      '/templates',
      'Service templates and recipes',
      'Browse database, storage, automation, monitoring and AI service recipes. Find the supported workflow and setup guide for each service.'
    ),
});
