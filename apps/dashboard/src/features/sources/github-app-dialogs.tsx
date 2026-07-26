import { Check, ExternalLink, Trash } from 'lucide-react';
import type React from 'react';
import { Button } from '#/components/ui/button';
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogTitle,
} from '#/components/ui/dialog';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import { Textarea } from '#/components/ui/textarea';
import type { GithubApp } from '#/interfaces/settings';

export const GithubIcon = ({ className }: { className?: string }) => (
  <svg
    xmlns="http://www.w3.org/2000/svg"
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    strokeWidth="2"
    strokeLinecap="round"
    strokeLinejoin="round"
    className={className}
  >
    <path d="M15 22v-4a4.8 4.8 0 0 0-1-3.5c3 0 6-2 6-5.5.08-1.25-.27-2.48-1-3.5.28-1.15.28-2.35 0-3.5 0 0-1 0-3 1.5-2.64-.5-5.36-.5-8 0C6 2 5 2 5 2c-.3 1.15-.3 2.35 0 3.5A5.403 5.403 0 0 0 4 9c0 3.5 3 5.5 6 5.5-.39.49-.68 1.05-.85 1.65-.17.6-.22 1.23-.15 1.85v4" />
    <path d="M9 18c-4.51 2-5-2-7-2" />
  </svg>
);

export interface GithubAppDialogsProps {
  isEditing: boolean;
  setIsEditing: (open: boolean) => void;
  editingApp: GithubApp | null;
  deletingApp: string | null;
  setDeletingApp: (id: string | null) => void;
  accessToken: string;
  setAccessToken: (val: string) => void;
  webhookSecret: string;
  setWebhookSecret: (val: string) => void;
  appId: string;
  setAppId: (val: string) => void;
  clientId: string;
  setClientId: (val: string) => void;
  appSlug: string;
  setAppSlug: (val: string) => void;
  privateKey: string;
  setPrivateKey: (val: string) => void;
  handleSave: (e: React.FormEvent) => void;
  confirmDelete: () => void;
  isSaving: boolean;
  isDeleting: boolean;
  manifestStr: string;
}

export function GithubAppDialogs({
  isEditing,
  setIsEditing,
  editingApp,
  deletingApp,
  setDeletingApp,
  accessToken,
  setAccessToken,
  webhookSecret,
  setWebhookSecret,
  appId,
  setAppId,
  clientId,
  setClientId,
  appSlug,
  setAppSlug,
  privateKey,
  setPrivateKey,
  handleSave,
  confirmDelete,
  isSaving,
  isDeleting,
  manifestStr,
}: GithubAppDialogsProps) {
  return (
    <>
      <Dialog open={isEditing} onOpenChange={setIsEditing}>
        <DialogContent className="border/50 gap-0 bg-card/95 p-0 backdrop-blur-xl sm:max-w-4xl [&>button]:hidden">
          <div className="px-5 pt-5 pb-4">
            <div className="flex items-start justify-between">
              <div className="flex flex-col">
                <DialogTitle className="font-bold text-foreground text-xl tracking-tight">
                  {editingApp ? 'Edit GitHub App' : 'Connect GitHub App'}
                </DialogTitle>
                <DialogDescription>Configure Github Integration</DialogDescription>
              </div>
              <div className="flex items-center gap-3">
                {!editingApp && (
                  <Button
                    asChild
                    variant="ghost"
                    className="h-9 gap-2 font-mono font-semibold text-[11px] text-primary uppercase tracking-wider hover:text-primary"
                  >
                    <a href="https://github.com/settings/apps/new" target="_blank" rel="noreferrer">
                      <ExternalLink className="h-3.5 w-3.5" />
                      Create App
                    </a>
                  </Button>
                )}
                <DialogClose asChild>
                  <Button
                    variant="ghost"
                    className="font-medium text-foreground/80 text-sm hover:bg-transparent hover:text-foreground"
                  >
                    CLOSE
                  </Button>
                </DialogClose>
              </div>
            </div>
          </div>

          <div className="h-px w-full bg-border/50" />

          <div className="px-5 pt-4 pb-5">
            {!editingApp && (
              <div className="mt-4 rounded-lg border border-primary/20 bg-primary/5 p-4">
                <div className="flex items-center justify-between">
                  <div className="space-y-1">
                    <p className="font-bold text-[10px] text-primary uppercase tracking-widest">
                      ONE-CLICK CONNECT
                    </p>
                    <p className="text-sm">
                      Create the GitHub App and fill in every credential automatically.
                    </p>
                  </div>
                  <form
                    action="https://github.com/settings/apps/new"
                    method="post"
                    target="_blank"
                    rel="noopener"
                  >
                    <input type="hidden" name="manifest" value={manifestStr} />
                    <Button
                      type="submit"
                      className="h-10 gap-2 bg-primary/20 font-bold text-primary text-xs uppercase tracking-wider hover:bg-primary/30"
                    >
                      <GithubIcon className="h-4 w-4" />
                      CONNECT WITH GITHUB
                    </Button>
                  </form>
                </div>
              </div>
            )}

            {!editingApp && (
              <div className="border/50 relative mt-12 mb-8 flex w-full justify-center border-t">
                <span className="absolute -top-3 bg-card px-4 font-bold text-[10px] text-muted-foreground uppercase tracking-widest">
                  OR ENTER MANUALLY
                </span>
              </div>
            )}

            <form onSubmit={handleSave} className="space-y-6">
              <div className="grid grid-cols-1 gap-6 md:grid-cols-2">
                <div className="space-y-2">
                  <Label className="font-bold text-[10px] text-muted-foreground uppercase tracking-widest">
                    GITHUB_CLIENT_SECRET
                  </Label>
                  <Input
                    type="password"
                    placeholder="GitHub personal access token"
                    value={accessToken}
                    onChange={(e) => setAccessToken(e.target.value)}
                    className="h-11 bg-background/50 font-mono"
                  />
                </div>
                <div className="space-y-2">
                  <Label className="font-bold text-[10px] text-muted-foreground uppercase tracking-widest">
                    GITHUB_WEBHOOK_SECRET
                  </Label>
                  <Input
                    type="password"
                    placeholder="Webhook secret"
                    value={webhookSecret}
                    onChange={(e) => setWebhookSecret(e.target.value)}
                    className="h-11 bg-background/50 font-mono"
                  />
                </div>
              </div>

              <div className="grid grid-cols-1 gap-6 md:grid-cols-2">
                <div className="space-y-2">
                  <Label className="font-bold text-[10px] text-muted-foreground uppercase tracking-widest">
                    GITHUB_APP_ID
                  </Label>
                  <Input
                    placeholder="e.g. 4322334"
                    value={appId}
                    onChange={(e) => setAppId(e.target.value)}
                    className="h-11 bg-background/50 font-mono"
                  />
                </div>
                <div className="space-y-2">
                  <Label className="font-bold text-[10px] text-muted-foreground uppercase tracking-widest">
                    GITHUB_CLIENT_ID
                  </Label>
                  <Input
                    placeholder="Iv1.xxxxxxxxxxxx"
                    value={clientId}
                    onChange={(e) => setClientId(e.target.value)}
                    className="h-11 bg-background/50 font-mono"
                  />
                </div>
              </div>

              <div className="space-y-2">
                <Label className="font-bold text-[10px] text-muted-foreground uppercase tracking-widest">
                  GITHUB_APP_SLUG
                </Label>
                <Input
                  placeholder="my-codedock-app"
                  value={appSlug}
                  onChange={(e) => setAppSlug(e.target.value)}
                  className="h-11 bg-background/50 font-mono"
                />
              </div>

              <div className="space-y-2">
                <Label className="font-bold text-[10px] text-muted-foreground uppercase tracking-widest">
                  GITHUB_PRIVATE_KEY
                </Label>
                <Textarea
                  placeholder="-----BEGIN RSA PRIVATE KEY-----"
                  value={privateKey}
                  onChange={(e) => setPrivateKey(e.target.value)}
                  className="min-h-40 bg-background/50 font-mono"
                />
              </div>

              <div className="mt-8 flex justify-end gap-3 pt-6">
                <Button
                  type="button"
                  variant="ghost"
                  onClick={() => setIsEditing(false)}
                  className="h-9 font-mono font-semibold text-[11px] uppercase tracking-wider"
                >
                  Cancel
                </Button>
                <Button
                  type="submit"
                  disabled={isSaving}
                  className="h-9 gap-2 font-mono font-semibold text-[11px] uppercase tracking-wider"
                >
                  <Check className="h-3.5 w-3.5" />
                  {isSaving ? 'Saving...' : 'Save Settings'}
                </Button>
              </div>
            </form>
          </div>
        </DialogContent>
      </Dialog>

      <Dialog open={!!deletingApp} onOpenChange={(open) => !open && setDeletingApp(null)}>
        <DialogContent className="border/50 gap-0 bg-card/95 p-0 backdrop-blur-xl sm:max-w-md [&>button]:hidden">
          <div className="p-5">
            <div className="flex items-start justify-between">
              <div className="flex flex-col">
                <DialogTitle className="flex items-center gap-2 font-bold text-destructive text-xl tracking-tight">
                  <Trash className="h-5 w-5" />
                  Remove GitHub App
                </DialogTitle>
                <DialogDescription>This will break existing deployments</DialogDescription>
              </div>
              <DialogClose asChild>
                <Button
                  variant="ghost"
                  className="font-medium text-foreground/80 text-sm hover:bg-transparent hover:text-foreground"
                >
                  CLOSE
                </Button>
              </DialogClose>
            </div>
          </div>

          <div className="flex items-center justify-end gap-3 p-5 pt-0">
            <Button
              variant="ghost"
              onClick={() => setDeletingApp(null)}
              className="h-9 font-mono font-semibold text-[11px] uppercase tracking-wider"
            >
              Cancel
            </Button>
            <Button
              onClick={(e) => {
                e.preventDefault();
                confirmDelete();
              }}
              disabled={isDeleting}
              variant="destructive"
              className="h-9 gap-2 font-mono font-semibold text-[11px] uppercase tracking-wider"
            >
              <Trash className="h-3.5 w-3.5" />
              {isDeleting ? 'Removing...' : 'Remove App'}
            </Button>
          </div>
        </DialogContent>
      </Dialog>
    </>
  );
}
