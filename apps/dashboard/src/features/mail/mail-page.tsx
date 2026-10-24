import { ArrowLeft, MailPlus, Plus, RefreshCw, Server, Trash2 } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { PageHeader } from '#/components/layout/page-header';
import { Button } from '@codedock/ui/button';
import { Card, CardContent } from '@codedock/ui/card';
import { Input } from '@codedock/ui/input';
import { QueryErrorState } from '#/components/query-error-state';
import { Skeleton } from '@codedock/ui/skeleton';
import { useGetPublicSettings } from '#/hooks/use-settings';
import {
  useCreateMailDomain,
  useCreateMailServer,
  useDeleteMailDomain,
  useDeleteMailServer,
  useListMailDomains,
  useListMailServers,
} from '#/hooks/use-mail';
import { CloudMailClient } from './cloud-mail-client';
import { MailDomainPanel } from './mail-domain-panel';
import { MailServerList } from './mail-server-list';
import { MailSetupForm } from './mail-setup-form';

export function MailPage() {
  const { data: publicRes } = useGetPublicSettings();
  const isCloud = publicRes?.data?.cloudMode ?? false;
  const [isAddingNew, setIsAddingNew] = useState(false);
  const [selectedServer, setSelectedServer] = useState<string | null>(null);
  const [selectedDomain, setSelectedDomain] = useState<string | null>(null);
  const [newDomain, setNewDomain] = useState('');

  const { data: servers = [], isLoading, isFetching, isError, refetch } = useListMailServers();
  const createServer = useCreateMailServer();
  const deleteServer = useDeleteMailServer();

  // Active server is either explicit selection, or the sole server if exactly 1 exists
  const activeServer =
    servers.find((server) => server.id === selectedServer) ??
    (servers.length === 1 && !isAddingNew ? servers[0] : null);
  const activeServerId = activeServer?.id ?? null;

  const { data: domains = [] } = useListMailDomains(activeServerId ?? undefined);
  const createDomain = useCreateMailDomain(activeServerId ?? undefined);
  const deleteDomain = useDeleteMailDomain(activeServerId ?? undefined);

  const activeDomain = domains.find((domain) => domain.id === selectedDomain) ?? domains[0] ?? null;

  const handleStartSetup = async (data: {
    domain: string;
    adminPassword: string;
    serverHost?: string;
    relayEnabled: boolean;
    relayProvider?: string;
    relayHost?: string;
    relayPort?: string;
    relayUsername?: string;
    relayPassword?: string;
  }) => {
    try {
      const created = await createServer.mutateAsync({
        domain: data.domain,
        maildUrl: `https://mail.${data.domain}:8443`,
        webmailUrl: `https://mail.${data.domain}`,
        apiToken: data.adminPassword,
      });

      toast.success(`Mail server for ${data.domain} installed successfully`);
      setIsAddingNew(false);
      setSelectedServer(created.id);
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Failed to install mail server');
    }
  };

  const handleDeleteServer = async (id: string) => {
    try {
      await deleteServer.mutateAsync(id);
      if (selectedServer === id) setSelectedServer(null);
      toast.success('Mail server removed');
    } catch {
      toast.error('Failed to remove mail server');
    }
  };

  if (isCloud) {
    return <CloudMailClient />;
  }

  if (isError) {
    return (
      <div className="space-y-6 pb-12">
        <PageHeader
          title="Emails"
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

  // View Resolution:
  // 1. Loading: Skeleton loader
  // 2. No servers or actively adding: Render setup hero form (MailSetupForm)
  // 3. Multiple servers and none selected: Render MailServerList
  // 4. Server selected or single server: Render Mail Server Console

  return (
    <div className="space-y-6 pb-12">
      {(servers.length > 0 || isAddingNew) && (
        <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
          <div>
            <PageHeader
              title="Emails"
              description="Self-hosted mail servers, mailboxes, DNS and deliverability."
            />
          </div>
          <div className="flex items-center gap-2">
            <Button
              variant="ghost"
              size="icon"
              disabled={isFetching}
              aria-label="Refresh mail servers"
              title="Refresh mail servers"
              onClick={() => void refetch()}
            >
              <RefreshCw className={`size-4 ${isFetching ? 'animate-spin' : ''}`} />
            </Button>
            {servers.length > 0 && !isAddingNew && !activeServerId && (
              <Button onClick={() => setIsAddingNew(true)} className="gap-2 rounded-xl">
                <MailPlus className="h-4 w-4" />
                Add mail server
              </Button>
            )}
          </div>
        </div>
      )}

      {isLoading ? (
        <div className="space-y-4">
          <Skeleton className="h-64 w-full rounded-2xl" />
        </div>
      ) : isAddingNew ? (
        <MailSetupForm
          onStart={handleStartSetup}
          isLoading={createServer.isPending}
          onCancel={() => setIsAddingNew(false)}
        />
      ) : servers.length === 0 ? (
        <CloudMailClient onSelfHost={() => setIsAddingNew(true)} />
      ) : servers.length > 1 && !activeServer ? (
        <MailServerList
          servers={servers.map((s) => ({
            id: s.id,
            domain: s.domain,
            maildUrl: s.maildUrl,
            status: 'running',
          }))}
          onOpen={(id) => {
            setSelectedServer(id);
            setSelectedDomain(null);
          }}
          onAddNew={() => setIsAddingNew(true)}
          onRemove={handleDeleteServer}
        />
      ) : activeServer ? (
        <div className="space-y-6">
          {servers.length > 1 && (
            <Button
              variant="ghost"
              size="sm"
              onClick={() => {
                setSelectedServer(null);
                setSelectedDomain(null);
              }}
              className="gap-1.5 text-muted-foreground text-xs hover:text-foreground"
            >
              <ArrowLeft className="size-4" />
              All mail servers
            </Button>
          )}

          <div className="grid grid-cols-1 gap-6 lg:grid-cols-3">
            <Card className="border-border/60 bg-card lg:col-span-1">
              <CardContent className="space-y-4 p-5">
                <div className="flex items-center gap-3 border-border/40 border-b pb-4">
                  <div className="flex size-9 shrink-0 items-center justify-center rounded-xl bg-violet-500/10 text-violet-500">
                    <Server className="size-4" />
                  </div>
                  <div className="min-w-0 flex-1">
                    <p className="truncate font-semibold text-foreground text-sm">
                      {activeServer.domain}
                    </p>
                    <p className="truncate font-mono text-muted-foreground text-xs">
                      {activeServer.maildUrl}
                    </p>
                  </div>
                  <Button
                    variant="ghost"
                    size="icon"
                    aria-label={`Remove ${activeServer.domain}`}
                    disabled={deleteServer.isPending}
                    onClick={() => handleDeleteServer(activeServer.id)}
                  >
                    <Trash2 className="size-4 text-muted-foreground hover:text-destructive" />
                  </Button>
                </div>

                <div className="space-y-2">
                  <div className="flex items-center justify-between">
                    <span className="font-semibold text-muted-foreground text-xs uppercase tracking-wider">
                      Domains ({domains.length})
                    </span>
                  </div>

                  <div className="flex gap-2">
                    <Input
                      id="mail-new-domain"
                      placeholder="Add domain, e.g. app.io"
                      value={newDomain}
                      onChange={(e) => setNewDomain(e.target.value)}
                      className="h-8 rounded-lg text-xs"
                    />
                    <Button
                      size="sm"
                      className="h-8"
                      disabled={newDomain.trim() === '' || createDomain.isPending}
                      onClick={() => {
                        createDomain.mutate(newDomain.trim(), {
                          onSuccess: () => setNewDomain(''),
                        });
                      }}
                    >
                      <Plus className="size-3.5" />
                    </Button>
                  </div>

                  <ul className="space-y-1 pt-2">
                    {domains.map((domain) => (
                      <li
                        key={domain.id}
                        className={`flex cursor-pointer items-center justify-between rounded-lg p-2 text-xs transition-colors ${
                          activeDomain?.id === domain.id
                            ? 'bg-primary/10 font-medium text-primary'
                            : 'text-foreground hover:bg-muted/40'
                        }`}
                        onClick={() => setSelectedDomain(domain.id)}
                      >
                        <span className="truncate font-mono">{domain.name}</span>
                        <Button
                          variant="ghost"
                          size="icon"
                          className="size-6 text-muted-foreground hover:text-destructive"
                          disabled={deleteDomain.isPending}
                          onClick={(e) => {
                            e.stopPropagation();
                            deleteDomain.mutate(domain.id);
                          }}
                        >
                          <Trash2 className="size-3.5" />
                        </Button>
                      </li>
                    ))}
                  </ul>
                </div>
              </CardContent>
            </Card>

            <div className="lg:col-span-2">
              {activeDomain ? (
                <MailDomainPanel
                  serverId={activeServer.id}
                  domainId={activeDomain.id}
                  domainName={activeDomain.name}
                />
              ) : (
                <Card className="border-border/60 bg-card">
                  <CardContent className="p-8 text-center text-muted-foreground text-sm">
                    Select a domain to manage its mailboxes, DNS records, and delivery checks.
                  </CardContent>
                </Card>
              )}
            </div>
          </div>
        </div>
      ) : null}
    </div>
  );
}
