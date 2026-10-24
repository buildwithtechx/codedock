import { Badge } from '@codedock/ui/badge';
import { Button } from '@codedock/ui/button';
import { browserSupportsWebAuthn } from '@simplewebauthn/browser';
import { Fingerprint, KeyRound, Loader2, Pencil, Plus, Trash2 } from 'lucide-react';
import { useState } from 'react';
import {
  PasskeyAddDialog,
  PasskeyDeleteDialog,
  PasskeyRenameDialog,
} from '#/features/profile/passkeys-dialogs';
import { SettingsSection } from '#/features/settings/settings-section';
import {
  useDeletePasskey,
  usePasskeys,
  useRegisterPasskey,
  useRenamePasskey,
} from '#/hooks/use-passkeys';
import type { Passkey } from '#/interfaces/passkey';

export function PasskeysSection() {
  const { data: passkeysResponse, isLoading } = usePasskeys();
  const registerPasskey = useRegisterPasskey();
  const renamePasskey = useRenamePasskey();
  const deletePasskey = useDeletePasskey();

  const [addOpen, setAddOpen] = useState(false);
  const [passkeyName, setPasskeyName] = useState('');

  const [renameTarget, setRenameTarget] = useState<Passkey | null>(null);
  const [newName, setNewName] = useState('');

  const [deleteTarget, setDeleteTarget] = useState<Passkey | null>(null);

  const supported = typeof window !== 'undefined' && browserSupportsWebAuthn();
  const passkeys = passkeysResponse?.data || [];

  const handleCreate = () => {
    registerPasskey.mutate(
      { name: passkeyName },
      {
        onSuccess: () => {
          setAddOpen(false);
          setPasskeyName('');
        },
      }
    );
  };

  const handleRename = () => {
    if (!renameTarget || !newName.trim()) return;
    renamePasskey.mutate(
      { id: renameTarget.id, name: newName.trim() },
      {
        onSuccess: () => {
          setRenameTarget(null);
          setNewName('');
        },
      }
    );
  };

  const handleDelete = () => {
    if (!deleteTarget) return;
    deletePasskey.mutate(deleteTarget.id, {
      onSuccess: () => setDeleteTarget(null),
    });
  };

  return (
    <>
      <SettingsSection
        icon={<KeyRound className="size-4 text-primary" />}
        title="Passkeys (WebAuthn)"
        description="Sign in securely using biometric authentication or security keys."
        action={
          <Button
            type="button"
            size="sm"
            disabled={!supported || registerPasskey.isPending}
            onClick={() => {
              setPasskeyName('');
              setAddOpen(true);
            }}
          >
            <Plus className="mr-1.5 h-3.5 w-3.5" />
            Add Passkey
          </Button>
        }
      >
        {!supported && (
          <div className="rounded-xl border border-amber-500/30 bg-amber-500/10 p-3 text-amber-500 text-xs">
            WebAuthn passkeys are not supported in your current browser or context.
          </div>
        )}

        {isLoading ? (
          <div className="flex h-20 items-center justify-center">
            <Loader2 className="h-5 w-5 animate-spin text-muted-foreground" />
          </div>
        ) : passkeys.length === 0 ? (
          <div className="flex flex-col items-center justify-center rounded-xl border border-border/60 border-dashed py-8 text-center">
            <div className="flex h-10 w-10 items-center justify-center rounded-full bg-primary/10 text-primary">
              <Fingerprint className="h-5 w-5" />
            </div>
            <p className="mt-2 font-medium text-foreground text-sm">No passkeys registered</p>
            <p className="max-w-xs text-muted-foreground text-xs">
              Add a passkey to sign in without typing a password using Touch ID, Face ID, or a
              security key.
            </p>
          </div>
        ) : (
          <div className="divide-y divide-border/40">
            {passkeys.map((pk) => (
              <div key={pk.id} className="flex items-center justify-between py-3">
                <div className="flex items-center gap-3">
                  <div className="flex h-9 w-9 items-center justify-center rounded-lg bg-primary/10 text-primary">
                    <Fingerprint className="h-4 w-4" />
                  </div>
                  <div>
                    <div className="flex items-center gap-2">
                      <span className="font-medium text-foreground text-sm">{pk.name}</span>
                      <Badge variant="outline" className="text-[10px]">
                        {pk.backedUp ? 'Synced' : 'Device-bound'}
                      </Badge>
                    </div>
                    <p className="text-muted-foreground text-xs">
                      Added {new Date(pk.createdAt).toLocaleDateString()}
                      {pk.lastUsedAt &&
                        ` · Last used ${new Date(pk.lastUsedAt).toLocaleDateString()}`}
                    </p>
                  </div>
                </div>

                <div className="flex items-center gap-1">
                  <Button
                    type="button"
                    variant="ghost"
                    size="sm"
                    className="h-8 w-8 p-0"
                    onClick={() => {
                      setRenameTarget(pk);
                      setNewName(pk.name);
                    }}
                  >
                    <Pencil className="h-3.5 w-3.5 text-muted-foreground" />
                  </Button>
                  <Button
                    type="button"
                    variant="ghost"
                    size="sm"
                    className="h-8 w-8 p-0 text-destructive hover:text-destructive"
                    onClick={() => setDeleteTarget(pk)}
                  >
                    <Trash2 className="h-3.5 w-3.5" />
                  </Button>
                </div>
              </div>
            ))}
          </div>
        )}
      </SettingsSection>

      <PasskeyAddDialog
        open={addOpen}
        onOpenChange={setAddOpen}
        name={passkeyName}
        onNameChange={setPasskeyName}
        onSubmit={handleCreate}
        isPending={registerPasskey.isPending}
      />

      <PasskeyRenameDialog
        target={renameTarget}
        onClose={() => setRenameTarget(null)}
        name={newName}
        onNameChange={setNewName}
        onSubmit={handleRename}
        isPending={renamePasskey.isPending}
      />

      <PasskeyDeleteDialog
        target={deleteTarget}
        onClose={() => setDeleteTarget(null)}
        onSubmit={handleDelete}
        isPending={deletePasskey.isPending}
      />
    </>
  );
}
