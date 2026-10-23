import { Button } from '@codedock/ui/button';
import { Check, Copy, Key, ShieldCheck, Terminal } from 'lucide-react';
import { useEffect, useState } from 'react';
import { toast } from 'sonner';
import type { McpTransportMode } from '#/interfaces/mcp';
import { SettingsSection } from '../settings-section';
import { McpClientSetup } from './mcp-client-setup';
import { McpTokenSetup } from './mcp-token-setup';
import { McpToolsList } from './mcp-tools-list';

type AuthMode = 'oauth' | 'token';

export function McpSettings() {
  const [copied, setCopied] = useState(false);
  const [copiedCurl, setCopiedCurl] = useState(false);
  const [authMode, setAuthMode] = useState<AuthMode>('oauth');
  const [transport, setTransport] = useState<McpTransportMode>('http');
  const [origin, setOrigin] = useState('http://localhost:8080');

  useEffect(() => {
    if (typeof window !== 'undefined') {
      setOrigin(window.location.origin);
    }
  }, []);

  const endpointUrl = transport === 'http' ? `${origin}/api/mcp` : `${origin}/api/mcp/sse`;
  const curlCommand = `curl -s ${origin}/api/mcp \\\n  -H "Authorization: Bearer <YOUR_PAT>" \\\n  -H "Content-Type: application/json" \\\n  -d '{"jsonrpc":"2.0","id":1,"method":"tools/list"}'`;

  const copyEndpoint = async () => {
    await navigator.clipboard.writeText(endpointUrl);
    setCopied(true);
    toast.success('MCP endpoint copied');
    setTimeout(() => setCopied(false), 2000);
  };

  const copyCurl = async () => {
    await navigator.clipboard.writeText(curlCommand);
    setCopiedCurl(true);
    toast.success('Verification command copied');
    setTimeout(() => setCopiedCurl(false), 2000);
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
        title="Model Context Protocol (MCP)"
        description="Connect Claude Code, Cursor, Codex, VS Code, Windsurf, and Zed directly to your Codedock instance."
      >
        <div className="space-y-6">
          <div className="flex items-start gap-3 rounded-xl border border-emerald-500/20 bg-emerald-500/5 p-4 text-xs">
            <div className="flex size-8 shrink-0 items-center justify-center rounded-lg bg-emerald-500/10 text-emerald-600 dark:text-emerald-400">
              <ShieldCheck className="size-4" />
            </div>
            <div className="min-w-0 flex-1 leading-relaxed">
              <p className="font-semibold text-foreground">
                Zero-trust, route-backed MCP server (Protocol 2025-06-18)
              </p>
              <p className="mt-0.5 text-muted-foreground">
                Every MCP tool call dispatches through Codedock&apos;s real permission-checked HTTP
                router. Connected AI agents can only see and invoke operations permitted by their
                credential.
              </p>
            </div>
          </div>

          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <div>
              <p className="mb-2 font-medium text-foreground text-xs">Transport Protocol</p>
              <div className="inline-flex rounded-xl border border-border/60 bg-muted/30 p-1">
                <button
                  type="button"
                  onClick={() => setTransport('http')}
                  aria-pressed={transport === 'http'}
                  className={`rounded-lg px-3 py-1.5 font-medium text-xs transition-colors ${
                    transport === 'http'
                      ? 'bg-card text-foreground shadow-xs'
                      : 'text-muted-foreground hover:text-foreground'
                  }`}
                >
                  Streamable HTTP (Recommended)
                </button>
                <button
                  type="button"
                  onClick={() => setTransport('sse')}
                  aria-pressed={transport === 'sse'}
                  className={`rounded-lg px-3 py-1.5 font-medium text-xs transition-colors ${
                    transport === 'sse'
                      ? 'bg-card text-foreground shadow-xs'
                      : 'text-muted-foreground hover:text-foreground'
                  }`}
                >
                  SSE (/api/mcp/sse)
                </button>
              </div>
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
            </div>
          </div>

          <div>
            <p className="mb-1.5 font-medium text-foreground text-xs">Endpoint URL</p>
            <div className="flex items-center gap-2">
              <code className="min-w-0 flex-1 truncate rounded-xl border border-border/60 bg-muted/30 px-3.5 py-2 font-mono text-foreground text-xs">
                {endpointUrl}
              </code>
              <Button
                variant="outline"
                size="sm"
                onClick={copyEndpoint}
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
          </div>

          <div className="border-border/40 border-t pt-4">
            {authMode === 'oauth' ? (
              <McpClientSetup endpoint={endpointUrl} />
            ) : (
              <McpTokenSetup endpoint={endpointUrl} />
            )}
          </div>

          <div className="rounded-xl border border-border/50 bg-muted/15 p-3.5">
            <div className="flex items-center justify-between gap-2">
              <span className="inline-flex items-center gap-1.5 font-medium text-foreground text-xs">
                <Terminal className="size-3.5 text-primary" />
                Verify from terminal
              </span>
              <Button variant="ghost" size="sm" onClick={copyCurl} className="h-7 gap-1.5 text-xs">
                {copiedCurl ? (
                  <Check className="size-3 text-emerald-500" />
                ) : (
                  <Copy className="size-3" />
                )}
                {copiedCurl ? 'Copied' : 'Copy curl'}
              </Button>
            </div>
            <pre className="mt-2 overflow-x-auto font-mono text-[11px] text-muted-foreground leading-relaxed">
              {curlCommand}
            </pre>
          </div>
        </div>
      </SettingsSection>

      <SettingsSection
        icon={<ShieldCheck className="size-4 text-primary" />}
        title="MCP Tools & Guided Prompts Catalog"
        description="Live route-backed JSON-RPC tools and multi-step playbooks exposed to connected AI clients."
        collapsible
        defaultOpen
      >
        <McpToolsList />
      </SettingsSection>
    </div>
  );
}
