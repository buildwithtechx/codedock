import { useNavigate } from '@tanstack/react-router';
import { useEffect, useRef } from 'react';
import { PageHeader } from '#/components/layout/page-header';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '#/components/ui/tabs';
import { AppCatalog } from './app-catalog';
import { useAppCatalog } from './hooks';
import { TemplatesGallery } from './templates-gallery';

export function NewAppTabs({ deepLinkAppId }: { deepLinkAppId?: string }) {
  const navigate = useNavigate();
  const { data, isLoading } = useAppCatalog();
  const autoStarted = useRef(false);
  const catalog = Array.isArray(data) ? data : [];

  useEffect(() => {
    if (autoStarted.current || isLoading || !deepLinkAppId) return;
    const entry = catalog.find((app) => app.id === deepLinkAppId);
    if (entry) {
      autoStarted.current = true;
      void navigate({ to: '/apps/new/$appId', params: { appId: entry.id } });
    }
  }, [isLoading, catalog, deepLinkAppId, navigate]);

  return (
    <div>
      <PageHeader
        title="Create app"
        description="Install a one-click app or deploy from a template."
      />
      <Tabs defaultValue="one-click" className="mt-6 w-full">
        <TabsList variant="line">
          <TabsTrigger value="one-click">One-click apps</TabsTrigger>
          <TabsTrigger value="templates">Templates</TabsTrigger>
        </TabsList>
        <TabsContent value="one-click">
          <AppCatalog embedded deepLinkAppId={deepLinkAppId} />
        </TabsContent>
        <TabsContent value="templates">
          <TemplatesGallery />
        </TabsContent>
      </Tabs>
    </div>
  );
}
