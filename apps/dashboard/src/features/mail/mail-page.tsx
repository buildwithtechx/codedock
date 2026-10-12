import { MailPlus, Server, Trash2 } from 'lucide-react';
import { useState } from 'react';
import { PageHeader } from '#/components/layout/page-header';
import { Button } from '#/components/ui/button';
import { Card, CardContent } from '#/components/ui/card';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import { QueryErrorState } from '#/components/ui/query-error-state';
import {
  useCreateMailDomain,
  useDeleteMailDomain,
  useDeleteMailServer,
  useListMailDomains,
  useListMailServers,
} from '#/hooks/use-mail';
import { MailDomainPanel } from './mail-domain-panel';
import { MailServerForm } from './mail-server-form';

export function MailPage() {
  const [showForm, setShowForm] = useState(false);
  const [selectedServer, setSelectedServer] = useState<string | null>(null);
  const [selectedDomain, setSelectedDomain] = useState<string | null>(null);
  const [newDomain, setNewDomain] = useState('');
  const { data: servers = [], isLoading, isError, refetch } = useListMailServers();
  const { data: domains = [] } = useListMailDomains(selectedServer ?? undefined);
  const createDomain = useCreateMailDomain(selectedServer ?? undefined);
  const deleteDomain = useDeleteMailDomain(selectedServer ?? undefined);
  const deleteServer = useDeleteMailServer();

  const activeServer = servers.find((server) => server.id === selectedServer) ?? servers[0] ?? null;
  const activeServerId = activeServer?.id ?? null;
  const activeDomain = domains.find((domain) => domain.id === selectedDomain) ?? domains[0] ?? null;

  if (isError) {
    return (
      <div className="space-y-6 pb-12">
        <PageHeader
          title="Mail"
          description="Self-hosted mail servers, mailboxes, DNS and deliverability."
        />
        <QueryErrorState
          title="Mail unavailable"
          description="Mail servers could not be loaded."
          onRetry={() => refetch()}
        />
      </div>
    );
  }

  return (
    <div className="space-y-6 pb-12">
      <PageHeader
        title="Mail"
        description="Self-hosted mail servers, mailboxes, DNS and deliverability."
      />
      <div className="flex justify-end">
        <Button className="gap-2" onClick={() => setShowForm((value) => !value)}>
          <MailPlus className="h-4 w-4" />
          Register server
        </Button>
      </div>
      {showForm && <MailServerForm onDone={() => setShowForm(false)} />}
      {isLoading ? (
        <p className="text-muted-foreground text-sm">Loading…</p>
      ) : servers.length === 0 ? (
        <Card className="border-border/60 bg-card">
          <CardContent className="p-6 text-muted-foreground text-sm">
            No mail servers registered. Register one to manage domains and mailboxes.
          </CardContent>
        </Card>
      ) : (
        <div className="grid grid-cols-1 gap-4 lg:grid-cols-3">
          <Card className="border-border/60 bg-card lg:col-span-1">
            <CardContent className="space-y-2 p-4">
              {servers.map((server) => (
                <div
                  key={server.id}
                  className={`flex w-full items-center justify-between rounded-lg border p-3 text-left ${
                    activeServerId === server.id ? 'border-primary/60' : 'border-border/60'
                  }`}
                >
                  <button
                    type="button"
                    className="flex flex-1 items-center gap-2 text-left"
                    onClick={() => {
                      setSelectedServer(server.id);
                      setSelectedDomain(null);
                    }}
                  >
                    <Server className="h-4 w-4 text-muted-foreground" />
                    <span>
                      <span className="block font-medium text-sm">{server.domain}</span>
                      <span className="block truncate font-mono text-muted-foreground text-xs">
                        {server.maildUrl}
                      </span>
                    </span>
                  </button>
                  <Button
                    variant="ghost"
                    size="icon"
                    aria-label={`Remove ${server.domain}`}
                    disabled={deleteServer.isPending}
                    onClick={() => deleteServer.mutate(server.id)}
                  >
                    <Trash2 className="h-4 w-4 text-muted-foreground" />
                  </Button>
                </div>
              ))}
              <div className="flex gap-2 pt-2">
                <div className="flex-1 space-y-2">
                  <Label htmlFor="mail-new-domain">New domain</Label>
                  <Input
                    id="mail-new-domain"
                    placeholder="example.com"
                    value={newDomain}
                    onChange={(e) => setNewDomain(e.target.value)}
                  />
                </div>
                <Button
                  className="mt-7"
                  disabled={newDomain.trim() === '' || !activeServerId || createDomain.isPending}
                  onClick={() => {
                    createDomain.mutate(newDomain.trim(), { onSuccess: () => setNewDomain('') });
                  }}
                >
                  Add
                </Button>
              </div>
              <ul className="space-y-1 pt-2">
                {domains.map((domain) => (
                  <li key={domain.id} className="flex items-center justify-between text-sm">
                    <button
                      type="button"
                      className={`font-mono ${activeDomain?.id === domain.id ? 'font-semibold' : ''}`}
                      onClick={() => setSelectedDomain(domain.id)}
                    >
                      {domain.name}
                    </button>
                    <Button
                      variant="ghost"
                      size="icon"
                      disabled={deleteDomain.isPending}
                      onClick={() => deleteDomain.mutate(domain.id)}
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  </li>
                ))}
              </ul>
            </CardContent>
          </Card>
          <div className="lg:col-span-2">
            {activeServerId && activeDomain ? (
              <MailDomainPanel
                serverId={activeServerId}
                domainId={activeDomain.id}
                domainName={activeDomain.name}
              />
            ) : (
              <Card className="border-border/60 bg-card">
                <CardContent className="p-6 text-muted-foreground text-sm">
                  Select a domain to manage its mailboxes, DNS and checks.
                </CardContent>
              </Card>
            )}
          </div>
        </div>
      )}
    </div>
  );
}
