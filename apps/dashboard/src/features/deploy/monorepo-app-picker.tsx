import { Check, FolderGit2, Layers, Loader2, Plus } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Badge } from '@codedock/ui/badge';
import { Button } from '@codedock/ui/button';
import { Checkbox } from '@codedock/ui/checkbox';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@codedock/ui/dialog';
import { Input } from '@codedock/ui/input';
import { Label } from '@codedock/ui/label';
import { useInspectRepository } from '#/hooks/use-git';
import type { MonorepoCandidateApp } from '#/interfaces/git';
import { RepositoryPathBrowser } from '#/features/sources/repository-path-browser';

export interface MonorepoAppPickerProps {
  apps: MonorepoCandidateApp[];
  selectedIndices: number[];
  onToggleIndex: (index: number) => void;
  onUpdateApp: (index: number, patch: Partial<MonorepoCandidateApp>) => void;
  onAddCustomApp: (app: MonorepoCandidateApp) => void;
  repositoryUrl: string;
  branch?: string;
}

export function MonorepoAppPicker({
  apps,
  selectedIndices,
  onToggleIndex,
  onUpdateApp,
  onAddCustomApp,
  repositoryUrl,
  branch,
}: MonorepoAppPickerProps) {
  const [modalOpen, setModalOpen] = useState(false);
  const [pickedPath, setPickedPath] = useState('');
  const inspectRepo = useInspectRepository();
  const inspecting = inspectRepo.isPending;

  const handleAddFromPath = async () => {
    const cleanPath = pickedPath.replace(/^\/+|\/+$/g, '');
    if (!cleanPath) {
      toast.error('Please select a subdirectory');
      return;
    }

    if (apps.some((a) => a.rootDirectory === cleanPath)) {
      toast.error('This folder is already in the list');
      return;
    }

    try {
      const res = await inspectRepo.mutateAsync({
        repositoryUrl,
        branch,
        rootDirectory: cleanPath,
      });
      const inspection = res.data;
      const folderName = cleanPath.split('/').pop() || 'app';

      onAddCustomApp({
        name: folderName,
        rootDirectory: cleanPath,
        framework: inspection.framework || 'Custom',
        internalPort: inspection.internalPort || 3000,
        buildEngine: inspection.buildEngine || 'nixpacks',
        packageManager: inspection.packageManager,
        installCommand: inspection.installCommand,
        buildCommand: inspection.buildCommand,
        startCommand: inspection.startCommand,
        dockerfilePath: inspection.dockerfilePath,
      });

      toast.success(`Added ${folderName} to deployment`);
      setModalOpen(false);
      setPickedPath('');
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Failed to inspect subfolder');
    }
  };

  return (
    <div className="space-y-4 rounded-2xl border border-border/60 bg-card p-5">
      <div className="flex items-start justify-between gap-3">
        <div className="flex items-center gap-3">
          <div className="flex size-10 shrink-0 items-center justify-center rounded-xl bg-primary/10 text-primary">
            <Layers className="size-5" />
          </div>
          <div>
            <h3 className="font-semibold text-foreground text-sm">Monorepo Applications</h3>
            <p className="text-muted-foreground text-xs">
              Select which services to provision and deploy from this monorepo
            </p>
          </div>
        </div>

        <Button
          type="button"
          variant="outline"
          size="sm"
          onClick={() => setModalOpen(true)}
          className="h-8 gap-1.5 rounded-xl text-xs"
        >
          <Plus className="size-3.5" />
          Add Folder
        </Button>
      </div>

      <div className="divide-y divide-border/50 rounded-xl border border-border/60 bg-background/50">
        {apps.map((app, idx) => {
          const isSelected = selectedIndices.includes(idx);
          return (
            <div
              key={`${app.rootDirectory}-${idx}`}
              className={`p-4 transition-colors ${isSelected ? 'bg-primary/2' : 'opacity-70'}`}
            >
              <div className="flex items-start gap-3">
                <div className="pt-1">
                  <Checkbox
                    id={`app-select-${idx}`}
                    checked={isSelected}
                    onCheckedChange={() => onToggleIndex(idx)}
                    className="rounded-md"
                  />
                </div>

                <div className="flex-1 space-y-3">
                  <div className="flex flex-wrap items-center justify-between gap-2">
                    <div className="flex items-center gap-2">
                      <span className="font-mono text-xs text-muted-foreground">
                        {app.rootDirectory}
                      </span>
                      <Badge variant="outline" className="text-[10px] uppercase">
                        {app.framework}
                      </Badge>
                      <Badge variant="secondary" className="text-[10px]">
                        {app.buildEngine}
                      </Badge>
                    </div>

                    <div className="flex items-center gap-2">
                      <Label htmlFor={`port-${idx}`} className="text-xs text-muted-foreground">
                        Port:
                      </Label>
                      <Input
                        id={`port-${idx}`}
                        type="number"
                        value={app.internalPort}
                        disabled={!isSelected}
                        onChange={(e) =>
                          onUpdateApp(idx, { internalPort: Number(e.target.value) || 3000 })
                        }
                        className="h-7 w-20 rounded-lg font-mono text-xs"
                      />
                    </div>
                  </div>

                  {isSelected && (
                    <div className="grid gap-3 sm:grid-cols-2">
                      <div className="space-y-1">
                        <Label
                          htmlFor={`name-${idx}`}
                          className="text-[11px] text-muted-foreground"
                        >
                          Service Name
                        </Label>
                        <Input
                          id={`name-${idx}`}
                          value={app.name}
                          onChange={(e) => onUpdateApp(idx, { name: e.target.value })}
                          className="h-8 rounded-lg font-mono text-xs"
                        />
                      </div>
                      <div className="space-y-1">
                        <Label
                          htmlFor={`build-${idx}`}
                          className="text-[11px] text-muted-foreground"
                        >
                          Build Command
                        </Label>
                        <Input
                          id={`build-${idx}`}
                          value={app.buildCommand || ''}
                          onChange={(e) => onUpdateApp(idx, { buildCommand: e.target.value })}
                          placeholder="npm run build"
                          className="h-8 rounded-lg font-mono text-xs"
                        />
                      </div>
                    </div>
                  )}
                </div>
              </div>
            </div>
          );
        })}
      </div>

      <Dialog open={modalOpen} onOpenChange={setModalOpen}>
        <DialogContent className="max-w-md rounded-2xl p-6">
          <DialogHeader icon={FolderGit2}>
            <DialogTitle className="text-base font-semibold">Browse Repository Folders</DialogTitle>
            <DialogDescription className="text-xs text-muted-foreground">
              Select any subfolder in this repository to configure as an additional service.
            </DialogDescription>
          </DialogHeader>

          <div className="space-y-4 pt-2">
            <RepositoryPathBrowser
              repositoryUrl={repositoryUrl}
              branch={branch}
              kind="directory"
              value={pickedPath}
              onChange={setPickedPath}
            />

            <div className="flex items-center justify-between border-t border-border/40 pt-3">
              <span className="truncate font-mono text-xs text-muted-foreground">
                Selected: {pickedPath || 'none'}
              </span>
              <Button
                type="button"
                size="sm"
                onClick={handleAddFromPath}
                disabled={inspecting || !pickedPath}
                className="gap-1.5 rounded-xl text-xs"
              >
                {inspecting ? (
                  <Loader2 className="size-3.5 animate-spin" />
                ) : (
                  <Check className="size-3.5" />
                )}
                Add Service
              </Button>
            </div>
          </div>
        </DialogContent>
      </Dialog>
    </div>
  );
}
