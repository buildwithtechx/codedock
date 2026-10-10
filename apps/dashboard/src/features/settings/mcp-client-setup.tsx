import { Check, Copy, ExternalLink } from 'lucide-react';
import { useState } from 'react';
import { Button } from '#/components/ui/button';

type ClientDefinition = {
  id: string;
  name: string;
  badge?: string;
  icon: string;
  getSetup: (endpoint: string) => {
    label: string;
    code?: string;
    steps?: string[];
    note?: string;
    deeplink?: string;
    deeplinkText?: string;
  };
};

function encodeConfig(config: Record<string, unknown>): string {
  try {
    return encodeURIComponent(btoa(JSON.stringify(config)));
  } catch {
    return '';
  }
}

const CLIENTS: ClientDefinition[] = [
  {
    id: 'claude-code',
    name: 'Claude Code',
    badge: 'CLI',
    icon: '/icons/claude.png',
    getSetup: (e) => ({
      label: 'Run command in your project directory',
      code: `claude mcp add --transport sse codedock ${e}`,
      note: 'Claude Code automatically prompts to authorize access in your browser on first command execution.',
    }),
  },
  {
    id: 'cursor',
    name: 'Cursor',
    badge: 'IDE',
    icon: '/icons/cursor.png',
    getSetup: (e) => ({
      label: 'Add to ~/.cursor/mcp.json or Cursor Features Settings',
      code: JSON.stringify({ mcpServers: { codedock: { url: e } } }, null, 2),
      deeplink: `cursor://anysphere.cursor-deeplink/mcp/install?name=codedock&config=${encodeConfig({ url: e })}`,
      deeplinkText: 'Install in Cursor',
      note: 'Save the configuration and restart Cursor to activate Codedock tools in agent mode.',
    }),
  },
  {
    id: 'codex',
    name: 'Codex',
    badge: 'CLI',
    icon: '/icons/openai.png',
    getSetup: (e) => ({
      label: 'Register server and authenticate',
      code: `codex mcp add codedock --url ${e}\ncodex mcp login codedock`,
      note: 'Registers Codedock across local Codex CLI and editor tools.',
    }),
  },
  {
    id: 'vscode',
    name: 'VS Code',
    badge: 'IDE',
    icon: '/icons/copilot.png',
    getSetup: (e) => ({
      label: 'Register in VS Code terminal or .vscode/mcp.json',
      code: `code --add-mcp '{"name":"codedock","type":"sse","url":"${e}"}'`,
      note: 'Works seamlessly with GitHub Copilot agent mode and MCP extensions.',
    }),
  },
  {
    id: 'claude-desktop',
    name: 'Claude Desktop',
    icon: '/icons/claude.png',
    getSetup: (e) => ({
      label: 'Add to claude_desktop_config.json',
      code: JSON.stringify({ mcpServers: { codedock: { url: e } } }, null, 2),
      steps: [
        'Open Claude Desktop Settings → Developer or Connectors.',
        'Click Edit Config to open claude_desktop_config.json.',
        'Paste the JSON block below under mcpServers and restart Claude.',
      ],
      note: 'Requires Claude Desktop with MCP connector support enabled.',
    }),
  },
  {
    id: 'windsurf',
    name: 'Windsurf',
    badge: 'IDE',
    icon: '/icons/windsurf.png',
    getSetup: (e) => ({
      label: 'Add to ~/.codeium/windsurf/mcp_config.json',
      code: JSON.stringify({ mcpServers: { codedock: { serverUrl: e } } }, null, 2),
      note: 'Cascade AI in Windsurf can immediately query Codedock services and deployments.',
    }),
  },
  {
    id: 'zed',
    name: 'Zed',
    icon: '/icons/zed.png',
    getSetup: (e) => ({
      label: 'Add to Zed settings.json',
      code: JSON.stringify(
        {
          context_servers: {
            codedock: {
              source: 'custom',
              command: { path: 'npx', args: ['-y', 'mcp-remote', e] },
            },
          },
        },
        null,
        2
      ),
      note: 'Bridges remote Codedock SSE endpoints through mcp-remote in Zed.',
    }),
  },
  {
    id: 'generic',
    name: 'Generic / Other',
    icon: '/icons/mcp.png',
    getSetup: (e) => ({
      label: 'Universal MCP SSE client configuration',
      code: JSON.stringify({ mcpServers: { codedock: { url: e } } }, null, 2),
      note: 'Compatible with any client supporting standard Model Context Protocol SSE transports.',
    }),
  },
];

export function McpClientSetup({ endpoint }: { endpoint: string }) {
  const [selectedId, setSelectedId] = useState(CLIENTS[0].id);
  const [copied, setCopied] = useState(false);

  const selectedClient = CLIENTS.find((c) => c.id === selectedId) || CLIENTS[0];
  const setup = selectedClient.getSetup(endpoint);

  const copyCode = async () => {
    if (!setup.code) return;
    await navigator.clipboard.writeText(setup.code);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <div className="space-y-4">
      <div>
        <p className="mb-2 font-medium text-foreground text-xs">Choose Your AI Assistant or IDE</p>
        <div className="flex flex-wrap gap-1.5">
          {CLIENTS.map((client) => {
            const isActive = client.id === selectedId;
            return (
              <button
                key={client.id}
                type="button"
                onClick={() => setSelectedId(client.id)}
                className={`flex items-center gap-2 rounded-lg border px-3 py-1.5 font-medium text-xs transition-colors ${
                  isActive
                    ? 'border-primary/60 bg-primary/10 text-foreground'
                    : 'border-border/60 text-muted-foreground hover:bg-muted/40 hover:text-foreground'
                }`}
              >
                <img
                  src={client.icon}
                  alt={client.name}
                  className="size-4 object-contain brightness-90 dark:brightness-100"
                  onError={(e) => {
                    e.currentTarget.style.display = 'none';
                  }}
                />
                <span>{client.name}</span>
                {client.badge && (
                  <span className="rounded bg-muted px-1 py-0.2 font-mono text-[9px] text-muted-foreground">
                    {client.badge}
                  </span>
                )}
              </button>
            );
          })}
        </div>
      </div>

      <div className="space-y-2">
        <div className="flex items-center justify-between">
          <p className="font-medium text-foreground text-xs">{setup.label}</p>
          {setup.deeplink && (
            <a
              href={setup.deeplink}
              className="inline-flex items-center gap-1 font-medium text-primary text-xs hover:underline"
            >
              {setup.deeplinkText}
              <ExternalLink className="size-3" />
            </a>
          )}
        </div>

        {setup.steps && (
          <ol className="list-decimal space-y-1.5 rounded-xl border border-border/50 bg-muted/20 py-3 ps-8 pe-4 text-muted-foreground text-xs leading-relaxed">
            {setup.steps.map((step, i) => (
              <li key={i}>{step}</li>
            ))}
          </ol>
        )}

        {setup.code && (
          <div className="relative">
            <pre className="overflow-x-auto rounded-xl border border-border/60 bg-muted/30 p-3.5 pe-14 font-mono text-xs leading-relaxed">
              {setup.code}
            </pre>
            <Button
              variant="ghost"
              size="icon"
              onClick={copyCode}
              className="absolute top-2 right-2 size-7 text-muted-foreground hover:text-foreground"
              aria-label="Copy code"
            >
              {copied ? (
                <Check className="size-3.5 text-emerald-500" />
              ) : (
                <Copy className="size-3.5" />
              )}
            </Button>
          </div>
        )}

        {setup.note && (
          <p className="text-muted-foreground text-xs leading-relaxed">{setup.note}</p>
        )}
      </div>
    </div>
  );
}
