import { createFileRoute } from '@tanstack/react-router';
import { Loader2, Shield } from 'lucide-react';
import { PageHeader } from '#/components/layout/page-header';
import { Card, CardContent, CardHeader, CardTitle } from '#/components/ui/card';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '#/components/ui/table';
import { useAuditLogs } from '#/hooks/use-audit-logs';

export const Route = createFileRoute('/_dashboard/audit-logs')({
  component: AuditLogsPage,
});

function AuditLogsPage() {
  const { data: response, isLoading } = useAuditLogs();
  const logs = response?.data || [];

  return (
    <div className="space-y-6">
      <PageHeader
        title="Audit logs"
        description="Review security events and actions taken across the platform."
      />

      <Card className="border-border/80 bg-card shadow-sm">
        <CardHeader className="border-border/70 border-b">
          <CardTitle>Recent activity</CardTitle>
        </CardHeader>
        <CardContent>
          {isLoading ? (
            <div className="flex justify-center p-12">
              <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
            </div>
          ) : logs.length === 0 ? (
            <div className="flex flex-col items-center justify-center p-12 text-center text-muted-foreground">
              <Shield className="mb-4 h-8 w-8 opacity-20" />
              <p>No audit logs found.</p>
            </div>
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Timestamp</TableHead>
                  <TableHead>User</TableHead>
                  <TableHead>Action</TableHead>
                  <TableHead>Resource</TableHead>
                  <TableHead>IP Address</TableHead>
                  <TableHead>Details</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {logs.map((log) => (
                  <TableRow key={log.id}>
                    <TableCell className="whitespace-nowrap">
                      {new Date(log.createdAt).toLocaleString()}
                    </TableCell>
                    <TableCell>{log.userId}</TableCell>
                    <TableCell className="font-medium">{log.action}</TableCell>
                    <TableCell>{log.resource}</TableCell>
                    <TableCell>{log.ipAddress}</TableCell>
                    <TableCell className="max-w-xs truncate text-muted-foreground text-sm">
                      {log.details}
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
