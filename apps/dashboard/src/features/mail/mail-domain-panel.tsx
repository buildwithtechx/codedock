import {
  CheckCircle2,
  KeyRound,
  Plus,
  RefreshCw,
  ShieldAlert,
  Trash2,
  XCircle,
} from 'lucide-react';
import { useState } from 'react';
import { Alert, AlertDescription } from '@codedock/ui/alert';
import { Badge } from '@codedock/ui/badge';
import { Button } from '@codedock/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@codedock/ui/card';
import { Input } from '@codedock/ui/input';
import { Label } from '@codedock/ui/label';
import { QueryErrorState } from '#/components/query-error-state';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@codedock/ui/tabs';
import {
  useCreateDkim,
  useCreateMailbox,
  useDeleteMailbox,
  useListMailboxes,
  useMailChecks,
  useMailDns,
} from '#/hooks/use-mail';

export function MailDomainPanel({
  serverId,
  domainId,
  domainName,
}: {
  serverId: string;
  domainId: string;
  domainName: string;
}) {
  const [localPart, setLocalPart] = useState('');
  const [oneTimePassword, setOneTimePassword] = useState<string | null>(null);
  const {
    data: mailboxes = [],
    isLoading,
    isError,
    refetch,
  } = useListMailboxes(serverId, domainId);
  const { data: dns } = useMailDns(serverId, domainId);
  const {
    data: checks,
    isFetching: isFetchingChecks,
    refetch: refetchChecks,
  } = useMailChecks(serverId, domainId);
  const createMailbox = useCreateMailbox(serverId, domainId);
  const deleteMailbox = useDeleteMailbox(serverId, domainId);
  const createDkim = useCreateDkim(serverId, domainId);

  const addMailbox = () => {
    if (localPart.trim() === '') return;
    createMailbox.mutate(localPart.trim(), {
      onSuccess: (credentials) => {
        setLocalPart('');
        setOneTimePassword(credentials.password ?? null);
      },
    });
  };

  if (isError) {
    return (
      <QueryErrorState
        title="Mailboxes unavailable"
        description="The mail server could not be reached."
        onRetry={() => refetch()}
      />
    );
  }

  return (
    <div className="space-y-4">
      {oneTimePassword && (
        <Alert>
          <KeyRound className="h-4 w-4" />
          <AlertDescription>
            Mailbox created. The password is shown once and never stored:{' '}
            <code className="font-mono font-semibold">{oneTimePassword}</code>
          </AlertDescription>
        </Alert>
      )}
      <Tabs defaultValue="mailboxes">
        <TabsList>
          <TabsTrigger value="mailboxes">Mailboxes</TabsTrigger>
          <TabsTrigger value="dns">DNS</TabsTrigger>
          <TabsTrigger value="checks">Checks</TabsTrigger>
        </TabsList>
        <TabsContent value="mailboxes" className="space-y-4">
          <Card className="border-border/60 bg-card">
            <CardHeader>
              <CardTitle>Mailboxes on {domainName}</CardTitle>
            </CardHeader>
            <CardContent className="space-y-3">
              <div className="flex gap-2">
                <div className="flex-1 space-y-2">
                  <Label htmlFor="mailbox-local">New mailbox</Label>
                  <Input
                    id="mailbox-local"
                    placeholder="alice"
                    value={localPart}
                    onChange={(e) => setLocalPart(e.target.value)}
                  />
                </div>
                <Button
                  className="mt-7 gap-2"
                  disabled={localPart.trim() === '' || createMailbox.isPending}
                  onClick={addMailbox}
                >
                  <Plus className="h-4 w-4" />
                  Add
                </Button>
                <Button
                  className="mt-7 gap-2"
                  variant="outline"
                  disabled={createDkim.isPending}
                  onClick={() => createDkim.mutate()}
                >
                  <ShieldAlert className="h-4 w-4" />
                  DKIM
                </Button>
              </div>
              {isLoading ? (
                <p className="text-muted-foreground text-sm">Loading…</p>
              ) : (
                <ul className="divide-y divide-border/60">
                  {mailboxes.map((mailbox) => (
                    <li key={mailbox.id} className="flex items-center justify-between py-2">
                      <span className="font-mono text-sm">{mailbox.address}</span>
                      <Button
                        variant="ghost"
                        size="icon"
                        disabled={deleteMailbox.isPending}
                        onClick={() => deleteMailbox.mutate(mailbox.id)}
                      >
                        <Trash2 className="h-4 w-4" />
                      </Button>
                    </li>
                  ))}
                  {mailboxes.length === 0 && (
                    <li className="py-2 text-muted-foreground text-sm">No mailboxes yet.</li>
                  )}
                </ul>
              )}
            </CardContent>
          </Card>
        </TabsContent>
        <TabsContent value="dns">
          <Card className="border-border/60 bg-card">
            <CardHeader>
              <CardTitle>DNS records to publish</CardTitle>
            </CardHeader>
            <CardContent className="space-y-2">
              {(dns?.records ?? []).map((record, index) => (
                <div key={`${record.type}-${record.name}-${index}`} className="text-sm">
                  <div className="flex items-center gap-2">
                    <Badge variant="outline">{record.type}</Badge>
                    <code className="font-mono">{record.name}</code>
                    <span className="text-muted-foreground text-xs">{record.source}</span>
                  </div>
                  <code className="mt-1 block break-all rounded bg-muted p-2 font-mono text-xs">
                    {record.value}
                  </code>
                </div>
              ))}
              {(dns?.records ?? []).length === 0 && (
                <p className="text-muted-foreground text-sm">No records returned.</p>
              )}
            </CardContent>
          </Card>
        </TabsContent>
        <TabsContent value="checks">
          <Card className="border-border/60 bg-card">
            <CardHeader className="flex flex-row items-center justify-between">
              <CardTitle>Deliverability checks</CardTitle>
              <Button
                variant="outline"
                size="sm"
                disabled={isFetchingChecks}
                onClick={() => void refetchChecks()}
                className="gap-1.5"
              >
                <RefreshCw className={`h-3.5 w-3.5 ${isFetchingChecks ? 'animate-spin' : ''}`} />
                Re-run
              </Button>
            </CardHeader>
            <CardContent className="space-y-2">
              {(checks?.checks ?? []).map((check) => (
                <div key={check.name} className="flex items-start gap-2 text-sm">
                  {check.passed ? (
                    <CheckCircle2 className="mt-0.5 h-4 w-4 shrink-0 text-emerald-500" />
                  ) : (
                    <XCircle className="mt-0.5 h-4 w-4 shrink-0 text-destructive" />
                  )}
                  <div>
                    <p className="font-medium">{check.name}</p>
                    <p className="break-all text-muted-foreground text-xs">{check.detail}</p>
                  </div>
                </div>
              ))}
              {(checks?.checks ?? []).length === 0 && (
                <p className="text-muted-foreground text-sm">No checks returned.</p>
              )}
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  );
}
