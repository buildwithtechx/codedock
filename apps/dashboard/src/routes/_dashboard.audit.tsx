import { createFileRoute } from '@tanstack/react-router';
import { PageHeader } from '#/components/layout/page-header';
import { AuditLogList } from '#/features/audit/audit-log-list';

export const Route = createFileRoute('/_dashboard/audit')({
  component: AuditPage,
});

export function AuditPage() {
  return (
    <div className="space-y-6">
      <PageHeader
        title="Audit log"
        description="A readable history of who did what in this workspace, newest first."
      />
      <AuditLogList />
    </div>
  );
}
