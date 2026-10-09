import { Bot, Check, Copy, Terminal } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '#/components/ui/card';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '#/components/ui/tabs';

export function McpSettings() {
  const [copied, setCopied] = useState<string | null>(null);

  const endpointUrl =
    typeof window !== 'undefined'
      ? `${window.location.origin}/api/v1/mcp/sse`
      : 'http://localhost:8080/api/v1/mcp/sse';

  const copyToClipboard = (text: string, id: string) => {
    void navigator.clipboard.writeText(text);
    setCopied(id);
    toast.success('Copied to clipboard');
    setTimeout(() => setCopied(null), 2000);
  };

  const claudeConfig = JSON.stringify(
    {
      mcpServers: {
        codedock: {
          url: endpointUrl,
        },
      },
    },
    null,
    2
  );

  const cursorConfig = JSON.stringify(
    {
      mcpServers: {
        codedock: {
          url: endpointUrl,
        },
      },
    },
    null,
    2
  );

  return (
    <div className="space-y-6">
      <div>
        <h2 className="font-semibold text-foreground text-lg tracking-tight">
          Model Context Protocol (MCP)
        </h2>
        <p className="text-muted-foreground text-sm">
          Connect AI assistants directly to your Codedock platform via standard MCP JSON-RPC
          endpoints.
        </p>
      </div>

      <Card>
        <CardHeader className="pb-3">
          <div className="flex items-center gap-3">
            <div className="flex size-9 items-center justify-center rounded-xl bg-primary/10 text-primary">
              <Bot className="size-5" />
            </div>
            <div>
              <CardTitle className="text-base">MCP SSE Endpoint</CardTitle>
              <p className="text-muted-foreground text-xs">
                Server-Sent Events endpoint for real-time tool calling and resource inspection.
              </p>
            </div>
          </div>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="flex items-center gap-2">
            <code className="flex-1 truncate rounded-lg border bg-muted/50 px-3 py-2 font-mono text-xs">
              {endpointUrl}
            </code>
            <Button
              variant="outline"
              size="sm"
              className="gap-1.5"
              onClick={() => copyToClipboard(endpointUrl, 'endpoint')}
            >
              {copied === 'endpoint' ? (
                <Check className="size-3.5" />
              ) : (
                <Copy className="size-3.5" />
              )}
              {copied === 'endpoint' ? 'Copied' : 'Copy'}
            </Button>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader className="pb-3">
          <div className="flex items-center gap-3">
            <div className="flex size-9 items-center justify-center rounded-xl bg-muted text-muted-foreground">
              <Terminal className="size-5" />
            </div>
            <div>
              <CardTitle className="text-base">Client Configuration</CardTitle>
              <p className="text-muted-foreground text-xs">
                Add Codedock tools to your preferred AI coding environment.
              </p>
            </div>
          </div>
        </CardHeader>
        <CardContent>
          <Tabs defaultValue="claude" className="w-full">
            <TabsList variant="line">
              <TabsTrigger value="claude">Claude Desktop</TabsTrigger>
              <TabsTrigger value="cursor">Cursor</TabsTrigger>
              <TabsTrigger value="windsurf">Windsurf</TabsTrigger>
            </TabsList>

            <div className="mt-4">
              <TabsContent value="claude" className="space-y-3">
                <p className="text-muted-foreground text-xs">
                  Add this to your{' '}
                  <code className="rounded bg-muted px-1 py-0.5 font-mono text-[11px]">
                    claude_desktop_config.json
                  </code>
                  :
                </p>
                <div className="relative">
                  <pre className="overflow-x-auto rounded-lg border bg-muted/40 p-3 font-mono text-xs">
                    {claudeConfig}
                  </pre>
                  <Button
                    variant="ghost"
                    size="icon"
                    className="absolute top-2 right-2 h-7 w-7"
                    onClick={() => copyToClipboard(claudeConfig, 'claude')}
                  >
                    {copied === 'claude' ? (
                      <Check className="size-3.5" />
                    ) : (
                      <Copy className="size-3.5" />
                    )}
                  </Button>
                </div>
              </TabsContent>

              <TabsContent value="cursor" className="space-y-3">
                <p className="text-muted-foreground text-xs">
                  Add this to your{' '}
                  <code className="rounded bg-muted px-1 py-0.5 font-mono text-[11px]">
                    .cursor/mcp.json
                  </code>{' '}
                  or Cursor MCP Settings:
                </p>
                <div className="relative">
                  <pre className="overflow-x-auto rounded-lg border bg-muted/40 p-3 font-mono text-xs">
                    {cursorConfig}
                  </pre>
                  <Button
                    variant="ghost"
                    size="icon"
                    className="absolute top-2 right-2 h-7 w-7"
                    onClick={() => copyToClipboard(cursorConfig, 'cursor')}
                  >
                    {copied === 'cursor' ? (
                      <Check className="size-3.5" />
                    ) : (
                      <Copy className="size-3.5" />
                    )}
                  </Button>
                </div>
              </TabsContent>

              <TabsContent value="windsurf" className="space-y-3">
                <p className="text-muted-foreground text-xs">
                  Configure MCP in Windsurf using the SSE endpoint:
                </p>
                <div className="relative">
                  <pre className="overflow-x-auto rounded-lg border bg-muted/40 p-3 font-mono text-xs">
                    {cursorConfig}
                  </pre>
                  <Button
                    variant="ghost"
                    size="icon"
                    className="absolute top-2 right-2 h-7 w-7"
                    onClick={() => copyToClipboard(cursorConfig, 'windsurf')}
                  >
                    {copied === 'windsurf' ? (
                      <Check className="size-3.5" />
                    ) : (
                      <Copy className="size-3.5" />
                    )}
                  </Button>
                </div>
              </TabsContent>
            </div>
          </Tabs>
        </CardContent>
      </Card>
    </div>
  );
}
