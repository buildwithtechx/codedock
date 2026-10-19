import { createFileRoute, redirect } from '@tanstack/react-router';
import { z } from 'zod';

export const Route = createFileRoute('/_dashboard/projects/new')({
  beforeLoad: ({ search }) => {
    throw redirect({
      to: '/library',
      search: {
        tab:
          search.template === 'examples'
            ? 'examples'
            : search.template === 'one-click'
              ? 'apps'
              : undefined,
      },
    });
  },
  validateSearch: z.object({
    template: z.enum(['one-click', 'examples']).optional(),
  }),
});
