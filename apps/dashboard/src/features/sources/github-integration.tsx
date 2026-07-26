import { useNavigate } from '@tanstack/react-router';
import { Edit, Plus, Trash } from 'lucide-react';
import { useEffect, useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import {
  useDeleteGitApp,
  useExchangeGithubManifest,
  useGetGitApps,
  useSaveGitApp,
} from '#/features/settings/hooks';
import type { GithubApp } from '#/features/settings/types';
import { Route } from '#/routes/_dashboard.sources';
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

  // Form state
  const [accessToken, setAccessToken] = useState('');
  const [webhookSecret, setWebhookSecret] = useState('');
  const [appId, setAppId] = useState('');
  const [clientId, setClientId] = useState('');
  const [appSlug, setAppSlug] = useState('');
  const [privateKey, setPrivateKey] = useState('');

  const navigate = useNavigate();
  const search = Route.useSearch();

  useEffect(() => {
    const code = search.code;
    if (code && !exchangeMutation.isPending && !exchangeMutation.isSuccess) {
      exchangeMutation.mutate(
        { code },
        {
          onSuccess: () => {
            navigate({ to: '/sources', replace: true });
            toast.success('GitHub App connected successfully!');
            setIsEditing(false);
            setEditingApp(null);
          },
          onError: (err) => {
            navigate({ to: '/sources', replace: true });
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
            redirect_url: `${window.location.origin}/dashboard/sources`,
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
      <div className="mb-5 flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="flex h-12 w-12 shrink-0 items-center justify-center rounded-lg border border-primary/20 bg-primary/10 text-primary">
            <GithubIcon className="h-6 w-6" />
          </div>
          <div>
            <h1 className="font-bold text-xl">Connected GitHub Apps</h1>
            <p className="text-muted-foreground text-sm">
              Connect GitHub Apps to automatically deploy pushed commits.
            </p>
          </div>
        </div>
        <Button
          className="gap-2"
          onClick={() => {
            setEditingApp(null);
            setIsEditing(true);
          }}
        >
          <Plus className="h-4 w-4" />
          ADD GITHUB APP
        </Button>
      </div>

      {apps.length > 0 ? (
        <div className="grid grid-cols-1 gap-4">
          {apps.map((app) => (
            <div key={app.id} className="rounded-xl border border bg-card p-6">
              <div className="flex items-start justify-between">
                <div className="flex items-start gap-4">
                  <div className="flex h-12 w-12 shrink-0 items-center justify-center rounded-lg border border-primary/20 bg-primary/10 text-primary">
                    <GithubIcon className="h-6 w-6" />
                  </div>
                  <div className="space-y-1">
                    <p className="font-bold text-[10px] text-muted-foreground uppercase tracking-widest">
                      GITHUB INTEGRATION
                    </p>
                    <div className="flex items-center gap-3">
                      <h2 className="font-bold text-xl tracking-tight">
                        {app.name || 'GitHub App'}
                      </h2>
                      <div className="rounded border border-primary/30 bg-primary/10 px-2 py-0.5 font-semibold text-[10px] text-primary uppercase tracking-widest">
                        CONNECTED
                      </div>
                    </div>
                  </div>
                </div>
                <div className="flex items-center gap-2">
                  <Button
                    variant="outline"
                    size="icon"
                    className="border/50 h-10 w-10 bg-transparent hover:bg-card"
                    onClick={() => {
                      setEditingApp(app);
                      setIsEditing(true);
                    }}
                  >
                    <Edit className="h-4 w-4" />
                  </Button>
                  <Button
                    variant="outline"
                    size="icon"
                    className="border/50 h-10 w-10 bg-transparent hover:bg-destructive/10 hover:text-destructive"
                    onClick={() => setDeletingApp(app.id)}
                    disabled={deleteMutation.isPending}
                  >
                    <Trash className="h-4 w-4" />
                  </Button>
                </div>
              </div>

              <div className="mt-6 grid grid-cols-1 gap-4 md:grid-cols-2">
                <div className="border/50 rounded-lg border bg-background/50 p-4">
                  <p className="font-medium text-[10px] text-muted-foreground uppercase tracking-widest">
                    APP SLUG
                  </p>
                  <p className="mt-2 font-mono text-sm">{app.name || 'Not set'}</p>
                </div>
                <div className="border/50 rounded-lg border bg-background/50 p-4">
                  <p className="font-medium text-[10px] text-muted-foreground uppercase tracking-widest">
                    APP ID
                  </p>
                  <p className="mt-2 truncate font-mono text-sm">{app.appId || 'Not set'}</p>
                </div>
              </div>
            </div>
          ))}
        </div>
      ) : (
        <div className="flex h-64 flex-col items-center justify-center rounded-xl border border border-dashed bg-card/40">
          <div className="flex h-12 w-12 items-center justify-center rounded-xl border border-primary/20 bg-primary/10">
            <GithubIcon className="h-5 w-5 text-primary" />
          </div>
          <h3 className="mt-4 font-bold text-foreground text-lg tracking-tight">
            No GitHub Apps connected
          </h3>
          <p className="mt-1 max-w-sm text-center text-muted-foreground text-sm">
            Connect a GitHub App to deploy repositories and receive webhooks.
          </p>
          <Button
            className="mt-6 gap-2"
            onClick={() => {
              setEditingApp(null);
              setIsEditing(true);
            }}
          >
            <Plus className="h-4 w-4" />
            CONNECT GITHUB APP
          </Button>
        </div>
      )}

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
