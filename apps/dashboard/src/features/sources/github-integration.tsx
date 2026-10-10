import { useNavigate, useRouterState } from '@tanstack/react-router';
import { Edit, Plus, Trash } from 'lucide-react';
import { useEffect, useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import type { GithubApp } from '#/features/settings';
import {
  useDeleteGitApp,
  useExchangeGithubManifest,
  useGetGitApps,
  useSaveGitApp,
} from '#/features/settings';
import { SettingsSection } from '#/features/settings/settings-section';
import { GithubAppDialogs, GithubIcon } from './github-app-dialogs';

export function GithubIntegration() {
  const { data, isLoading } = useGetGitApps();
  const saveMutation = useSaveGitApp();
  const deleteMutation = useDeleteGitApp();
  const exchangeMutation = useExchangeGithubManifest();

  const apps = (data?.data as GithubApp[]) || [];

  const [isEditing, setIsEditing] = useState(false);
  const [editingApp, setEditingApp] = useState<GithubApp | null>(null);
  const [deletingApp, setDeletingApp] = useState<string | null>(null);

  const [accessToken, setAccessToken] = useState('');
  const [webhookSecret, setWebhookSecret] = useState('');
  const [appId, setAppId] = useState('');
  const [clientId, setClientId] = useState('');
  const [appSlug, setAppSlug] = useState('');
  const [privateKey, setPrivateKey] = useState('');

  const navigate = useNavigate();
  const search = useRouterState({ select: (state) => state.location.search as { code?: string } });

  useEffect(() => {
    const code = search.code;
    if (code && !exchangeMutation.isPending && !exchangeMutation.isSuccess) {
      exchangeMutation.mutate(
        { code },
        {
          onSuccess: () => {
            navigate({ to: '/settings', search: { tab: 'git' } as never, replace: true });
            toast.success('GitHub App connected successfully!');
            setIsEditing(false);
            setEditingApp(null);
          },
          onError: (err) => {
            navigate({ to: '/settings', search: { tab: 'git' } as never, replace: true });
            toast.error(err.message || 'Failed to connect GitHub App');
          },
        }
      );
    }
  }, [
    exchangeMutation.isPending,
    exchangeMutation.isSuccess,
    exchangeMutation.mutate,
    search.code,
    navigate,
  ]);

  useEffect(() => {
    if (editingApp) {
      setAppId(editingApp.appId || '');
      setClientId(editingApp.clientId || '');
      setWebhookSecret(editingApp.webhookSecret ? '********' : '');
      setAccessToken(editingApp.clientSecret ? '********' : '');
      setAppSlug(editingApp.name || '');
      setPrivateKey(editingApp.privateKey ? '********' : '');
    } else {
      setAppId('');
      setClientId('');
      setWebhookSecret('');
      setAccessToken('');
      setAppSlug('');
      setPrivateKey('');
    }
  }, [editingApp]);

  const handleSave = (e: React.FormEvent) => {
    e.preventDefault();

    const payload = {
      ...(editingApp?.id ? { id: editingApp.id } : {}),
      appId,
      clientId,
      name: appSlug,
      ...(accessToken !== '********' ? { clientSecret: accessToken } : {}),
      ...(webhookSecret !== '********' ? { webhookSecret } : {}),
      ...(privateKey !== '********' ? { privateKey } : {}),
    };

    saveMutation.mutate(payload, {
      onSuccess: () => {
        setIsEditing(false);
        setEditingApp(null);
        toast.success('GitHub settings saved successfully');
      },
      onError: (err: Error) => {
        toast.error(err.message || 'Failed to save GitHub settings');
      },
    });
  };

  const confirmDelete = () => {
    if (!deletingApp) return;
    deleteMutation.mutate(deletingApp, {
      onSuccess: () => {
        if (editingApp?.id === deletingApp) {
          setIsEditing(false);
          setEditingApp(null);
        }
        toast.success('GitHub connection removed');
        setDeletingApp(null);
      },
      onError: (err: Error) => {
        toast.error(err.message || 'Failed to remove GitHub connection');
      },
    });
  };

  const manifestStr =
    typeof window !== 'undefined'
      ? (() => {
          const isLocalhost =
            window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1';
          const baseManifest = {
            name: `codedock-${Math.random().toString(36).substring(7)}`,
            url: window.location.origin,
            redirect_url: `${window.location.origin}/settings?tab=sources`,
            public: false,
            default_permissions: {
              contents: 'read',
              metadata: 'read',
              pull_requests: 'read',
              emails: 'read',
            },
          };

          if (!isLocalhost) {
            return JSON.stringify({
              ...baseManifest,
              hook_attributes: {
                url: `${window.location.origin}/api/webhooks/github/services/generic`,
              },
              default_events: ['push', 'pull_request'],
            });
          }

          return JSON.stringify(baseManifest);
        })()
      : '{}';

  if (isLoading) {
    return <div className="h-64 animate-pulse rounded-xl bg-card" />;
  }

  return (
    <div className="space-y-6">
      <SettingsSection
        icon={<GithubIcon className="size-4 text-primary" />}
        title="Connected GitHub Apps"
        description="Connect GitHub Apps to automatically deploy pushed commits."
        action={
          <Button
            size="sm"
            className="gap-1.5"
            onClick={() => {
              setEditingApp(null);
              setIsEditing(true);
            }}
          >
            <Plus className="h-4 w-4" />
            Add GitHub App
          </Button>
        }
      >
        {apps.length > 0 ? (
          <div className="space-y-4">
            {apps.map((app) => (
              <div key={app.id} className="rounded-xl border border-border/50 bg-background/50 p-5">
                <div className="flex items-start justify-between">
                  <div className="flex items-start gap-4">
                    <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl border border-border/50 bg-background">
                      <GithubIcon className="h-5 w-5" />
                    </div>
                    <div className="space-y-1">
                      <div className="flex items-center gap-2">
                        <h3 className="font-semibold text-base text-foreground">
                          {app.name || 'GitHub App'}
                        </h3>
                        <div className="rounded border border-primary/30 bg-primary/10 px-2 py-0.5 font-semibold text-[10px] text-primary uppercase tracking-widest">
                          CONNECTED
                        </div>
                      </div>
                      <p className="font-mono text-muted-foreground text-xs">
                        App ID: {app.appId || 'Not set'}
                      </p>
                    </div>
                  </div>
                  <div className="flex items-center gap-2">
                    <Button
                      variant="outline"
                      size="icon"
                      className="h-8 w-8"
                      onClick={() => {
                        setEditingApp(app);
                        setIsEditing(true);
                      }}
                    >
                      <Edit className="h-3.5 w-3.5" />
                    </Button>
                    <Button
                      variant="outline"
                      size="icon"
                      className="h-8 w-8 text-destructive hover:bg-destructive/10 hover:text-destructive"
                      onClick={() => setDeletingApp(app.id)}
                      disabled={deleteMutation.isPending}
                    >
                      <Trash className="h-3.5 w-3.5" />
                    </Button>
                  </div>
                </div>
              </div>
            ))}
          </div>
        ) : (
          <div className="flex flex-col items-center justify-center rounded-xl border border-border/60 border-dashed py-10 text-center">
            <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-muted text-muted-foreground">
              <GithubIcon className="h-5 w-5" />
            </div>
            <h4 className="mt-3 font-medium text-foreground text-sm">No GitHub Apps connected</h4>
            <p className="mt-1 max-w-sm text-muted-foreground text-xs">
              Connect a GitHub App to deploy repositories and receive webhooks automatically.
            </p>
          </div>
        )}
      </SettingsSection>

      <GithubAppDialogs
        isEditing={isEditing}
        setIsEditing={setIsEditing}
        editingApp={editingApp}
        deletingApp={deletingApp}
        setDeletingApp={setDeletingApp}
        accessToken={accessToken}
        setAccessToken={setAccessToken}
        webhookSecret={webhookSecret}
        setWebhookSecret={setWebhookSecret}
        appId={appId}
        setAppId={setAppId}
        clientId={clientId}
        setClientId={setClientId}
        appSlug={appSlug}
        setAppSlug={setAppSlug}
        privateKey={privateKey}
        setPrivateKey={setPrivateKey}
        handleSave={handleSave}
        confirmDelete={confirmDelete}
        isSaving={saveMutation.isPending}
        isDeleting={deleteMutation.isPending}
        manifestStr={manifestStr}
      />
    </div>
  );
}
