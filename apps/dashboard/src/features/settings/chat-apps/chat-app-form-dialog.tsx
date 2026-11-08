import type { ChatAppConnection, ChatAppProvider, ChatAppTransport } from '@codedock/contracts';
import {
  Button,
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  Input,
  Label,
} from '@codedock/ui';
import { useState } from 'react';
import { useCreateChatConnection, useUpdateChatConnection } from '#/hooks/use-chat-apps';

interface ChatAppFormDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  provider: ChatAppProvider;
  connection?: ChatAppConnection | null;
}

export function ChatAppFormDialog({
  open,
  onOpenChange,
  provider,
  connection,
}: ChatAppFormDialogProps) {
  const isEditing = Boolean(connection);
  const createMutation = useCreateChatConnection();
  const updateMutation = useUpdateChatConnection();
  const isSaving = createMutation.isPending || updateMutation.isPending;

  const [name, setName] = useState(connection?.name || `${provider.toUpperCase()} Bot`);
  const [transport, setTransport] = useState<ChatAppTransport>(connection?.transport || 'webhook');
  const [botUsername, setBotUsername] = useState(connection?.botUsername || '');
  const [botToken, setBotToken] = useState('');
  const [appToken, setAppToken] = useState('');
  const [signingSecret, setSigningSecret] = useState('');
  const [publicKey, setPublicKey] = useState('');
  const [appId, setAppId] = useState('');
  const [appPassword, setAppPassword] = useState('');
  const [secretToken, setSecretToken] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const credentials: Record<string, string> = {};

    if (provider === 'slack') {
      if (botToken) credentials.bot_token = botToken;
      if (appToken) credentials.app_token = appToken;
      if (signingSecret) credentials.signing_secret = signingSecret;
    } else if (provider === 'discord') {
      if (botToken) credentials.bot_token = botToken;
      if (publicKey) credentials.public_key = publicKey;
      if (appId) credentials.app_id = appId;
    } else if (provider === 'telegram') {
      if (botToken) credentials.bot_token = botToken;
      if (secretToken) credentials.secret_token = secretToken;
    } else if (provider === 'teams') {
      if (appId) credentials.app_id = appId;
      if (appPassword) credentials.app_password = appPassword;
    }

    if (isEditing && connection) {
      await updateMutation.mutateAsync({
        id: connection.id,
        input: {
          name,
          transport,
          botUsername,
          credentials: Object.keys(credentials).length > 0 ? credentials : undefined,
        },
      });
    } else {
      await createMutation.mutateAsync({
        provider,
        transport,
        purpose: 'both',
        name,
        botUsername,
        credentials,
      });
    }
    onOpenChange(false);
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>
            {isEditing ? `Edit ${provider} Integration` : `Connect ${provider}`}
          </DialogTitle>
          <DialogDescription>
            Configure your bot credentials and communication transport.
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit} className="space-y-3.5 py-2">
          <div className="space-y-1">
            <Label htmlFor="conn-name" className="text-xs">
              Integration Name
            </Label>
            <Input
              id="conn-name"
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="e.g. Production Slack Bot"
              required
            />
          </div>

          <div className="space-y-1">
            <Label htmlFor="bot-username" className="text-xs">
              Bot Display Username
            </Label>
            <Input
              id="bot-username"
              value={botUsername}
              onChange={(e) => setBotUsername(e.target.value)}
              placeholder="@codedock-bot"
            />
          </div>

          <div className="space-y-1">
            <Label htmlFor="transport" className="text-xs">
              Transport Type
            </Label>
            <select
              id="transport"
              value={transport}
              onChange={(e) => setTransport(e.target.value as ChatAppTransport)}
              className="w-full rounded-md border border-input bg-background px-3 py-1.5 text-xs shadow-xs focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"
            >
              <option value="webhook">Webhook (Incoming HTTPS)</option>
              <option value="outbound">Outbound / Socket Mode (NAT & Local Friendly)</option>
            </select>
          </div>

          {/* Provider Specific Inputs */}
          {provider === 'slack' && (
            <>
              <div className="space-y-1">
                <Label htmlFor="slack-bot-token" className="text-xs">
                  Bot User Token ({isEditing ? 'leave blank to keep existing' : 'xoxb-...'})
                </Label>
                <Input
                  id="slack-bot-token"
                  type="password"
                  value={botToken}
                  onChange={(e) => setBotToken(e.target.value)}
                  placeholder="xoxb-..."
                />
              </div>
              {transport === 'outbound' ? (
                <div className="space-y-1">
                  <Label htmlFor="slack-app-token" className="text-xs">
                    App-Level Token (Socket Mode: xapp-...)
                  </Label>
                  <Input
                    id="slack-app-token"
                    type="password"
                    value={appToken}
                    onChange={(e) => setAppToken(e.target.value)}
                    placeholder="xapp-..."
                  />
                </div>
              ) : (
                <div className="space-y-1">
                  <Label htmlFor="slack-signing-secret" className="text-xs">
                    Signing Secret (x-slack-signature verification)
                  </Label>
                  <Input
                    id="slack-signing-secret"
                    type="password"
                    value={signingSecret}
                    onChange={(e) => setSigningSecret(e.target.value)}
                    placeholder="Slack Signing Secret"
                  />
                </div>
              )}
            </>
          )}

          {provider === 'discord' && (
            <>
              <div className="space-y-1">
                <Label htmlFor="discord-bot-token" className="text-xs">
                  Bot Token
                </Label>
                <Input
                  id="discord-bot-token"
                  type="password"
                  value={botToken}
                  onChange={(e) => setBotToken(e.target.value)}
                  placeholder="Bot token from Discord Developer Portal"
                />
              </div>
              <div className="space-y-1">
                <Label htmlFor="discord-public-key" className="text-xs">
                  Public Key (Ed25519)
                </Label>
                <Input
                  id="discord-public-key"
                  value={publicKey}
                  onChange={(e) => setPublicKey(e.target.value)}
                  placeholder="64-character hex public key"
                />
              </div>
              <div className="space-y-1">
                <Label htmlFor="discord-app-id" className="text-xs">
                  Application ID
                </Label>
                <Input
                  id="discord-app-id"
                  value={appId}
                  onChange={(e) => setAppId(e.target.value)}
                  placeholder="Discord Application ID"
                />
              </div>
            </>
          )}

          {provider === 'telegram' && (
            <>
              <div className="space-y-1">
                <Label htmlFor="telegram-token" className="text-xs">
                  Telegram Bot Token
                </Label>
                <Input
                  id="telegram-token"
                  type="password"
                  value={botToken}
                  onChange={(e) => setBotToken(e.target.value)}
                  placeholder="123456789:ABCdefGhIJKlmNoPQRsTUVwxyZ"
                />
              </div>
              <div className="space-y-1">
                <Label htmlFor="telegram-secret" className="text-xs">
                  Webhook Secret Token
                </Label>
                <Input
                  id="telegram-secret"
                  type="password"
                  value={secretToken}
                  onChange={(e) => setSecretToken(e.target.value)}
                  placeholder="Optional secret header verification token"
                />
              </div>
            </>
          )}

          {provider === 'teams' && (
            <>
              <div className="space-y-1">
                <Label htmlFor="teams-app-id" className="text-xs">
                  Microsoft App ID
                </Label>
                <Input
                  id="teams-app-id"
                  value={appId}
                  onChange={(e) => setAppId(e.target.value)}
                  placeholder="Azure Bot Registration App ID"
                />
              </div>
              <div className="space-y-1">
                <Label htmlFor="teams-app-password" className="text-xs">
                  App Client Secret
                </Label>
                <Input
                  id="teams-app-password"
                  type="password"
                  value={appPassword}
                  onChange={(e) => setAppPassword(e.target.value)}
                  placeholder="Azure Bot Registration Client Secret"
                />
              </div>
            </>
          )}

          <DialogFooter className="pt-2">
            <Button type="button" variant="outline" size="sm" onClick={() => onOpenChange(false)}>
              Cancel
            </Button>
            <Button type="submit" size="sm" disabled={isSaving}>
              {isSaving ? 'Saving...' : isEditing ? 'Save Changes' : 'Connect Bot'}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
