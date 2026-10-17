import { TabsContent } from '@codedock/ui/tabs';
import { ManagedServerActivity } from '#/features/servers/managed/managed-server-activity';
import { ManagedServerNetwork } from '#/features/servers/managed/managed-server-network';
import { ManagedServerRuntime } from '#/features/servers/managed/managed-server-runtime';
import { ManagedServerSettings } from '#/features/servers/managed/managed-server-settings';
import { ManagedServerSsh } from '#/features/servers/managed/managed-server-ssh';
import { ManagedServerWorkloads } from '#/features/servers/managed/managed-server-workloads';
import { ServerComponentsTab } from '#/features/servers/server-components-tab';
import { ServerOverviewTab } from '#/features/servers/server-overview-tab';
import { ServerSecurityTab } from '#/features/servers/server-security-tab';
import { ServerSettingsTab } from '#/features/servers/server-settings-tab';
import { ServerTerminalTab } from '#/features/servers/server-terminal-tab';
import type { Server } from '#/interfaces/server';

interface ServerDetailTabsProps {
  server: Server;
  isManaged: boolean;
  onRetry: () => void;
}

export function ServerDetailTabsContent({ server, isManaged, onRetry }: ServerDetailTabsProps) {
  if (isManaged) {
    return (
      <>
        <TabsContent value="overview">
          <ServerOverviewTab server={server} />
        </TabsContent>
        <TabsContent value="activity">
          <ManagedServerActivity server={server.managed} onRetry={onRetry} />
        </TabsContent>
        <TabsContent value="networking">
          <ManagedServerNetwork serverId={server.id} />
        </TabsContent>
        <TabsContent value="workloads">
          <ManagedServerWorkloads serverId={server.id} />
        </TabsContent>
        <TabsContent value="ssh">
          <ManagedServerSsh serverId={server.id} />
        </TabsContent>
        <TabsContent value="runtimeApi">
          <ManagedServerRuntime serverId={server.id} />
        </TabsContent>
        <TabsContent value="settings">
          <ManagedServerSettings serverId={server.id} />
        </TabsContent>
      </>
    );
  }

  return (
    <>
      <TabsContent value="overview">
        <ServerOverviewTab server={server} />
      </TabsContent>
      <TabsContent value="components">
        <ServerComponentsTab serverId={server.id} />
      </TabsContent>
      <TabsContent value="security">
        <ServerSecurityTab serverId={server.id} />
      </TabsContent>
      <TabsContent value="terminal">
        <ServerTerminalTab server={server} />
      </TabsContent>
      <TabsContent value="settings">
        <ServerSettingsTab server={server} />
      </TabsContent>
    </>
  );
}
