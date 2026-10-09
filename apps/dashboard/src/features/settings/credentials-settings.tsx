import { Globe, Key, ShieldCheck } from 'lucide-react';
import { useState } from 'react';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '#/components/ui/tabs';
import { DnsSettings } from '#/features/dns/dns-settings';
import { GitProviders } from '#/features/sources';
import { ApiKeysList } from '#/features/users/api-keys-list';

type CredentialSection = 'dns' | 'tokens' | 'git';

export function CredentialsSettings() {
  const [section, setSection] = useState<CredentialSection>('dns');

  return (
    <div className="space-y-6">
      <div>
        <h2 className="font-semibold text-foreground text-lg tracking-tight">Credentials</h2>
        <p className="text-muted-foreground text-sm">
          Manage third-party tokens, DNS provider credentials, and access keys in one unified view.
        </p>
      </div>

      <Tabs
        value={section}
        onValueChange={(val) => setSection(val as CredentialSection)}
        className="w-full"
      >
        <TabsList variant="line">
          <TabsTrigger value="dns" className="gap-2">
            <Globe className="size-4" />
            DNS Providers
          </TabsTrigger>
          <TabsTrigger value="tokens" className="gap-2">
            <Key className="size-4" />
            Access Tokens
          </TabsTrigger>
          <TabsTrigger value="git" className="gap-2">
            <ShieldCheck className="size-4" />
            Git Credentials
          </TabsTrigger>
        </TabsList>

        <div className="mt-6">
          <TabsContent value="dns" className="mt-0">
            <DnsSettings />
          </TabsContent>
          <TabsContent value="tokens" className="mt-0">
            <ApiKeysList />
          </TabsContent>
          <TabsContent value="git" className="mt-0">
            <GitProviders />
          </TabsContent>
        </div>
      </Tabs>
    </div>
  );
}
