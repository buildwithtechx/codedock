import { Bot, Check, Copy, Key, ShieldCheck } from 'lucide-react';
import { useEffect, useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import { McpClientSetup } from './mcp-client-setup';
import { McpTokenSetup } from './mcp-token-setup';
import { McpToolsList } from './mcp-tools-list';

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
      <div>
        <h2 className="font-semibold text-foreground text-lg tracking-tight">
          Model Context Protocol (MCP)
        </h2>
        <p className="mt-1 text-muted-foreground text-sm">
          Connect AI coding assistants and IDEs directly to your Codedock platform via standard MCP
          JSON-RPC endpoints.
        </p>
      </div>

      <div className="flex items-start gap-3 rounded-2xl border border-emerald-500/20 bg-emerald-500/5 p-4 text-xs">
        <div className="flex size-8 shrink-0 items-center justify-center rounded-xl bg-emerald-500/10 text-emerald-600 dark:text-emerald-400">
          <ShieldCheck className="size-4" />
        </div>
        <div className="min-w-0 flex-1 leading-relaxed">
          <p className="font-semibold text-foreground">Zero-configuration AI integration</p>
          <p className="mt-0.5 text-muted-foreground">
            Connected AI assistants can query projects, inspect deployments, tail runtime logs, and
            manage services safely based on their assigned workspace roles.
          </p>
        </div>
      </div>

      <section className="space-y-4 rounded-2xl border border-border/50 bg-card p-5">
        <div className="flex items-center gap-3">
          <div className="flex size-9 items-center justify-center rounded-xl bg-primary/10 text-primary">
            <Bot className="size-5" />
          </div>
          <div>
            <h3 className="font-semibold text-foreground text-sm">SSE Endpoint</h3>
            <p className="text-muted-foreground text-xs">
              Server-Sent Events endpoint used by MCP client transports.
            </p>
          </div>
        </div>

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
      </section>

      <section className="space-y-5 rounded-2xl border border-border/50 bg-card p-5">
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

        <div className="border-border/40 border-t pt-4">
          {authMode === 'oauth' ? (
            <McpClientSetup endpoint={endpointUrl} />
          ) : (
            <McpTokenSetup endpoint={endpointUrl} />
          )}
        </div>
      </section>

      <section className="rounded-2xl border border-border/50 bg-card p-5">
        <McpToolsList />
      </section>
    </div>
  );
}
