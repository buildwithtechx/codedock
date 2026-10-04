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
        title="Audit logs"
        description="Review security events, administrative changes, and operator actions across your workspace."
      />
      <AuditLogList />
    </div>
  );
}
