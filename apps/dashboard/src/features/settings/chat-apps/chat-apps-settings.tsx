import type { ChatAppConnection, ChatAppProvider } from '@codedock/contracts';
import {
  Button,
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
  Switch,
} from '@codedock/ui';
import { useSearch } from '@tanstack/react-router';
import {
  AlertCircle,
  Bot,
  CheckCircle2,
  HelpCircle,
  MessageSquare,
  Plus,
  Radio,
  Send,
  Trash2,
  Users,
} from 'lucide-react';
import { useState } from 'react';
import {
  useChatConnections,
  useChatIdentities,
  useDeleteChatConnection,
  useDeleteChatIdentity,
  useTestChatConnection,
  useUpdateChatConnection,
} from '#/hooks/use-chat-apps';
import { ChatAppFormDialog } from './chat-app-form-dialog';
import { ChatAppPairingDialog } from './chat-app-pairing-dialog';

interface ProviderMeta {
  provider: ChatAppProvider;
  name: string;
  description: string;
  defaultTransport: string;
  badge: string;
  icon: typeof MessageSquare;
}

const PROVIDERS: ProviderMeta[] = [
  {
    provider: 'slack',
    name: 'Slack',
    description: 'Socket Mode & Webhook bot for managing deployments and chatting in channels.',
    defaultTransport: 'Socket Mode (NAT Friendly)',
    badge: 'Popular',
    icon: MessageSquare,
  },
  {
    provider: 'discord',
    name: 'Discord',
    description: 'Discord bot gateway with Ed25519 webhook verification and slash commands.',
    defaultTransport: 'Gateway & Webhook',
    badge: 'Real-time',
    icon: Radio,
  },
  {
    provider: 'telegram',
    name: 'Telegram',
    description: 'Direct Telegram bot integration using long-polling or secure secret webhooks.',
    defaultTransport: 'Long Polling / Webhook',
    badge: 'Instant',
    icon: Send,
  },
  {
    provider: 'teams',
    name: 'Microsoft Teams',
    description: 'Enterprise Microsoft Teams bot integration via Azure Bot Framework.',
    defaultTransport: 'Bot Framework Webhook',
    badge: 'Enterprise',
    icon: Bot,
  },
];

export function ChatAppsSettings() {
  const searchParams = useSearch({ strict: false }) as Record<string, string | undefined>;
  const linkToken = searchParams?.link || null;

  const { data: connections = [] } = useChatConnections();
  const testMutation = useTestChatConnection();
  const updateMutation = useUpdateChatConnection();
  const deleteMutation = useDeleteChatConnection();

  const [activeForm, setActiveForm] = useState<{
    open: boolean;
    provider: ChatAppProvider;
    connection?: ChatAppConnection | null;
  }>({ open: false, provider: 'slack', connection: null });

  const [activeIdentitiesConnId, setActiveIdentitiesConnId] = useState<string | null>(null);
  const [pairingToken, setPairingToken] = useState<string | null>(linkToken);

  const getConnForProvider = (provider: ChatAppProvider) =>
    connections.find((c) => c.provider === provider);

  return (
    <div className="space-y-6">
      {/* Top Banner */}
      <div>
        <h2 className="font-semibold text-foreground text-lg">Inbound Chat Apps</h2>
        <p className="mt-1 text-muted-foreground text-xs">
          Connect your team's chat tools to Codedock. Inspect container health, roll out
          redeployments, and chat with the AI assistant directly from Slack, Discord, Telegram, or
          Teams.
        </p>
      </div>

      {/* Provider Connectors Grid */}
      <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
        {PROVIDERS.map((meta) => {
          const conn = getConnForProvider(meta.provider);
          const isConnected = Boolean(conn?.enabled);
          const Icon = meta.icon;

          return (
            <Card
              key={meta.provider}
              className={`relative overflow-hidden transition-all ${
                isConnected ? 'border-primary/40 bg-card' : 'border-border/60 bg-muted/20'
              }`}
            >
              <CardHeader className="p-4 pb-2">
                <div className="flex items-start justify-between gap-3">
                  <div className="flex items-center gap-2.5">
                    <div
                      className={`flex size-9 items-center justify-center rounded-lg ${
                        isConnected
                          ? 'bg-primary/10 text-primary'
                          : 'bg-muted text-muted-foreground'
                      }`}
                    >
                      <Icon className="size-5" />
                    </div>
                    <div>
                      <CardTitle className="flex items-center gap-2 font-medium text-sm">
                        {meta.name}
                        <span className="rounded bg-muted px-1.5 py-0.5 font-normal text-[10px] text-muted-foreground">
                          {meta.badge}
                        </span>
                      </CardTitle>
                      <CardDescription className="text-xs">
                        {conn?.name || meta.defaultTransport}
                      </CardDescription>
                    </div>
                  </div>

                  {/* Status Indicator */}
                  <div>
                    {conn ? (
                      <span
                        className={`inline-flex items-center gap-1 rounded-full px-2 py-0.5 font-medium text-[10px] ${
                          conn.status === 'connected'
                            ? 'bg-emerald-500/10 text-emerald-600 dark:text-emerald-400'
                            : 'bg-amber-500/10 text-amber-600 dark:text-amber-400'
                        }`}
                      >
                        {conn.status === 'connected' ? (
                          <CheckCircle2 className="size-3" />
                        ) : (
                          <AlertCircle className="size-3" />
                        )}
                        {conn.status === 'connected' ? 'Connected' : 'Action Required'}
                      </span>
                    ) : (
                      <span className="rounded-full bg-muted px-2 py-0.5 text-[10px] text-muted-foreground">
                        Not Configured
                      </span>
                    )}
                  </div>
                </div>
              </CardHeader>

              <CardContent className="space-y-3 p-4 pt-1">
                <p className="text-muted-foreground text-xs">{meta.description}</p>

                {conn && (
                  <div className="flex items-center justify-between border-border/40 border-t pt-2 text-xs">
                    <div className="flex items-center gap-2 text-[11px] text-muted-foreground">
                      <span>
                        Transport: <strong className="text-foreground">{conn.transport}</strong>
                      </span>
                      {conn.botUsername && (
                        <span>
                          • Bot: <strong className="text-foreground">{conn.botUsername}</strong>
                        </span>
                      )}
                    </div>
                    <div className="flex items-center gap-2">
                      <span className="text-[11px] text-muted-foreground">Enabled</span>
                      <Switch
                        checked={conn.enabled}
                        onCheckedChange={(checked) =>
                          updateMutation.mutate({ id: conn.id, input: { enabled: checked } })
                        }
                        aria-label={`Toggle ${meta.name} integration`}
                      />
                    </div>
                  </div>
                )}

                <div className="flex items-center justify-between gap-2 pt-1">
                  {conn ? (
                    <div className="flex items-center gap-1.5">
                      <Button
                        type="button"
                        variant="outline"
                        size="sm"
                        className="h-7 text-xs"
                        onClick={() =>
                          setActiveForm({ open: true, provider: meta.provider, connection: conn })
                        }
                      >
                        Configure
                      </Button>
                      <Button
                        type="button"
                        variant="ghost"
                        size="sm"
                        className="h-7 text-xs"
                        onClick={() => testMutation.mutate(conn.id)}
                        disabled={testMutation.isPending}
                      >
                        Test
                      </Button>
                      <Button
                        type="button"
                        variant="ghost"
                        size="sm"
                        className="h-7 text-xs"
                        onClick={() =>
                          setActiveIdentitiesConnId(
                            activeIdentitiesConnId === conn.id ? null : conn.id
                          )
                        }
                      >
                        <Users className="mr-1 size-3" />
                        Identities
                      </Button>
                      <Button
                        type="button"
                        variant="ghost"
                        size="icon"
                        className="size-7 text-destructive hover:bg-destructive/10"
                        onClick={() => deleteMutation.mutate(conn.id)}
                        aria-label={`Delete ${meta.name} connection`}
                      >
                        <Trash2 className="size-3.5" />
                      </Button>
                    </div>
                  ) : (
                    <Button
                      type="button"
                      variant="outline"
                      size="sm"
                      className="h-7 text-xs"
                      onClick={() =>
                        setActiveForm({ open: true, provider: meta.provider, connection: null })
                      }
                    >
                      <Plus className="mr-1 size-3.5" />
                      Connect {meta.name}
                    </Button>
                  )}
                </div>
              </CardContent>
            </Card>
          );
        })}
      </div>

      {/* Paired Users Drawer / Section */}
      {activeIdentitiesConnId && <IdentitiesSection connectionId={activeIdentitiesConnId} />}

      {/* Command Cheat Sheet */}
      <Card className="border-border/60 bg-muted/10">
        <CardHeader className="p-4 pb-2">
          <CardTitle className="flex items-center gap-2 font-medium text-sm">
            <HelpCircle className="size-4 text-primary" />
            Chat Command Reference
          </CardTitle>
          <CardDescription className="text-xs">
            Send commands to your connected bot in direct messages or mentioned channels.
          </CardDescription>
        </CardHeader>
        <CardContent className="p-4 pt-1">
          <div className="grid grid-cols-1 gap-2 text-xs sm:grid-cols-2">
            <div className="space-y-1 rounded border border-border/40 bg-background/60 p-2.5">
              <code className="font-semibold text-primary">/codedock link</code>
              <p className="text-[11px] text-muted-foreground">
                Initiates pairing flow. Returns a secure pairing token to connect your chat account.
              </p>
            </div>
            <div className="space-y-1 rounded border border-border/40 bg-background/60 p-2.5">
              <code className="font-semibold text-primary">/codedock verify &lt;code&gt;</code>
              <p className="text-[11px] text-muted-foreground">
                Completes pairing from within the chat channel with your 6-digit dashboard code.
              </p>
            </div>
            <div className="space-y-1 rounded border border-border/40 bg-background/60 p-2.5">
              <code className="font-semibold text-primary">/codedock projects</code>
              <p className="text-[11px] text-muted-foreground">
                Lists all accessible projects, active services, health status, and domains.
              </p>
            </div>
            <div className="space-y-1 rounded border border-border/40 bg-background/60 p-2.5">
              <code className="font-semibold text-primary">/codedock status &lt;id&gt;</code>
              <p className="text-[11px] text-muted-foreground">
                Displays live container status, git branch, and deployed image digests.
              </p>
            </div>
            <div className="space-y-1 rounded border border-border/40 bg-background/60 p-2.5">
              <code className="font-semibold text-primary">/codedock logs &lt;id&gt;</code>
              <p className="text-[11px] text-muted-foreground">
                Safely inspects container exit codes and gives an authenticated web logs link.
              </p>
            </div>
            <div className="space-y-1 rounded border border-border/40 bg-background/60 p-2.5">
              <code className="font-semibold text-primary">/codedock redeploy &lt;id&gt;</code>
              <p className="text-[11px] text-muted-foreground">
                Triggers redeployment. Generates a confirmation token to prevent accidental
                rollouts.
              </p>
            </div>
            <div className="space-y-1 rounded border border-border/40 bg-background/60 p-2.5">
              <code className="font-semibold text-primary">/codedock confirm &lt;token&gt;</code>
              <p className="text-[11px] text-muted-foreground">
                Confirms and executes a pending mutating action within the 5-minute timeout window.
              </p>
            </div>
            <div className="space-y-1 rounded border border-border/40 bg-background/60 p-2.5">
              <code className="font-semibold text-primary">/codedock ask &lt;query&gt;</code>
              <p className="text-[11px] text-muted-foreground">
                Asks the Codedock AI assistant for infrastructure analysis or advice.
              </p>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Form Dialog */}
      {activeForm.open && (
        <ChatAppFormDialog
          open={activeForm.open}
          onOpenChange={(open) => setActiveForm((prev) => ({ ...prev, open }))}
          provider={activeForm.provider}
          connection={activeForm.connection}
        />
      )}

      {/* Pairing Dialog */}
      {pairingToken && (
        <ChatAppPairingDialog token={pairingToken} onClose={() => setPairingToken(null)} />
      )}
    </div>
  );
}

function IdentitiesSection({ connectionId }: { connectionId: string }) {
  const { data: identities = [], isLoading } = useChatIdentities(connectionId);
  const deleteMutation = useDeleteChatIdentity();

  return (
    <Card className="border-border/80 bg-background">
      <CardHeader className="p-4 pb-2">
        <CardTitle className="font-medium text-sm">Paired Chat Identities</CardTitle>
        <CardDescription className="text-xs">
          Authorized team members paired via <code>/codedock link</code>.
        </CardDescription>
      </CardHeader>
      <CardContent className="p-4 pt-1">
        {isLoading ? (
          <p className="text-muted-foreground text-xs">Loading identities...</p>
        ) : identities.length === 0 ? (
          <p className="text-muted-foreground text-xs">
            No external chat users have paired with this connection yet. Run{' '}
            <code className="font-mono text-foreground">/codedock link</code> in chat to pair.
          </p>
        ) : (
          <div className="divide-y divide-border/40 rounded border border-border/60">
            {identities.map((idt) => (
              <div key={idt.id} className="flex items-center justify-between p-2.5 text-xs">
                <div>
                  <p className="font-medium text-foreground">
                    {idt.externalUsername || idt.externalUserId}{' '}
                    <span className="text-[10px] text-muted-foreground">({idt.provider})</span>
                  </p>
                  <p className="text-[11px] text-muted-foreground">
                    Role: <strong className="text-foreground">{idt.role}</strong> • Paired:{' '}
                    {new Date(idt.pairedAt).toLocaleDateString()}
                  </p>
                </div>
                <Button
                  type="button"
                  variant="ghost"
                  size="sm"
                  className="h-7 text-destructive text-xs hover:bg-destructive/10"
                  onClick={() => deleteMutation.mutate({ connectionId, identityId: idt.id })}
                >
                  Unlink
                </Button>
              </div>
            ))}
          </div>
        )}
      </CardContent>
    </Card>
  );
}
