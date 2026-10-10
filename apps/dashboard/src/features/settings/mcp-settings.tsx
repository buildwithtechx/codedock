import { Check, Copy, Key, ShieldCheck } from 'lucide-react';
import { useEffect, useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import { McpClientSetup } from './mcp-client-setup';
import { McpTokenSetup } from './mcp-token-setup';
import { McpToolsList } from './mcp-tools-list';
import { SettingsSection } from './settings-section';

type AuthMode = 'oauth' | 'token';

export function McpSettings() {
  const [copied, setCopied] = useState(false);
  const [authMode, setAuthMode] = useState<AuthMode>('oauth');
  const [endpointUrl, setEndpointUrl] = useState('');

  useEffect(() => {
    setEndpointUrl(
      typeof window !== 'undefined'
        ? `${window.location.origin}/api/v1/mcp/sse`
        : 'http://localhost:8080/api/v1/mcp/sse'
    );
  }, []);

  const copyEndpoint = async () => {
    if (!endpointUrl) return;
    await navigator.clipboard.writeText(endpointUrl);
    setCopied(true);
    toast.success('MCP endpoint copied');
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <div className="space-y-6">
      <SettingsSection
        icon={
          <img
            src="/icons/mcp.png"
            alt="MCP"
            className="size-5 object-contain"
            onError={(e) => {
              (e.currentTarget as HTMLElement).style.display = 'none';
            }}
          />
        }
        iconBg="bg-emerald-500/10 dark:bg-emerald-500/15"
        iconColor="text-emerald-600 dark:text-emerald-400"
        title="Model Context Protocol"
        description="Connect AI coding assistants and IDEs directly via standard MCP JSON-RPC endpoints."
      >
        <div className="space-y-6">
          <div className="flex items-start gap-3 rounded-xl border border-emerald-500/20 bg-emerald-500/5 p-4 text-xs">
            <div className="flex size-8 shrink-0 items-center justify-center rounded-lg bg-emerald-500/10 text-emerald-600 dark:text-emerald-400">
              <ShieldCheck className="size-4" />
            </div>
            <div className="min-w-0 flex-1 leading-relaxed">
              <p className="font-semibold text-foreground">Zero-configuration AI integration</p>
              <p className="mt-0.5 text-muted-foreground">
                Connected AI assistants can query projects, inspect deployments, tail runtime logs,
                and manage services safely based on their assigned workspace roles.
              </p>
            </div>
          </div>

          <div>
            <p className="mb-1.5 font-medium text-foreground text-xs">Endpoint URL</p>
            <div className="flex items-center gap-2">
              <code className="min-w-0 flex-1 truncate rounded-xl border border-border/60 bg-muted/30 px-3.5 py-2 font-mono text-foreground text-xs">
                {endpointUrl || 'Loading endpoint...'}
              </code>
              <Button
                variant="outline"
                size="sm"
                onClick={copyEndpoint}
                disabled={!endpointUrl}
                className="shrink-0 gap-1.5"
              >
                {copied ? (
                  <Check className="size-3.5 text-emerald-500" />
                ) : (
                  <Copy className="size-3.5" />
                )}
                <span>{copied ? 'Copied' : 'Copy'}</span>
              </Button>
            </div>
            <p className="mt-1.5 text-muted-foreground text-xs">
              SSE transport endpoint for Cursor, Claude Code, VS Code, and other MCP clients.
            </p>
          </div>

          <div>
            <p className="mb-2 font-medium text-foreground text-xs">Authentication Mode</p>
            <div className="inline-flex rounded-xl border border-border/60 bg-muted/30 p-1">
              <button
                type="button"
                onClick={() => setAuthMode('oauth')}
                aria-pressed={authMode === 'oauth'}
                className={`flex items-center gap-1.5 rounded-lg px-3 py-1.5 font-medium text-xs transition-colors ${
                  authMode === 'oauth'
                    ? 'bg-card text-foreground shadow-xs'
                    : 'text-muted-foreground hover:text-foreground'
                }`}
              >
                <ShieldCheck className="size-3.5 text-emerald-500" />
                <span>OAuth (Recommended)</span>
              </button>
              <button
                type="button"
                onClick={() => setAuthMode('token')}
                aria-pressed={authMode === 'token'}
                className={`flex items-center gap-1.5 rounded-lg px-3 py-1.5 font-medium text-xs transition-colors ${
                  authMode === 'token'
                    ? 'bg-card text-foreground shadow-xs'
                    : 'text-muted-foreground hover:text-foreground'
                }`}
              >
                <Key className="size-3.5 text-amber-500" />
                <span>Personal Access Token</span>
              </button>
            </div>
            <p className="mt-2 text-muted-foreground text-xs">
              {authMode === 'oauth'
                ? 'Assistants prompt to authorize via the browser on first tool call with automatic token refresh.'
                : 'Authenticate using a long-lived API token with customized project scopes.'}
            </p>
          </div>

          <div className="border-border/40 border-t pt-2">
            {authMode === 'oauth' ? (
              <McpClientSetup endpoint={endpointUrl} />
            ) : (
              <McpTokenSetup endpoint={endpointUrl} />
            )}
          </div>
        </div>
      </SettingsSection>

      <SettingsSection
        icon={<ShieldCheck className="size-4 text-primary" />}
        title="MCP Tools Catalog"
        description="Available JSON-RPC tools and actions exposed to connected assistants."
        collapsible
        defaultOpen
      >
        <McpToolsList />
      </SettingsSection>
    </div>
  );
}
