import { ArrowRight, FolderPlus, FolderUp, Loader2, Package, X } from 'lucide-react';
import { useEffect, useRef, useState } from 'react';
import { toast } from 'sonner';
import { Button } from '@codedock/ui/button';
import { Input } from '@codedock/ui/input';
import { Label } from '@codedock/ui/label';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@codedock/ui/select';
import { useDeployLibraryArchive } from '#/hooks/use-library';
import { useListOrganizations } from '#/hooks/use-organizations';
import { useCreateProject, useListAllProjects } from '#/hooks/use-projects';
import { useOrganizationStore } from '#/stores/organization-store';
import { FolderDropzone } from './folder-dropzone';
import {
  collectFolderFiles,
  type FolderEntry,
  type FolderInspection,
  inspectFolderFiles,
  packTarGz,
} from './folder-pack';
import { FolderUploadDone } from './folder-upload-done';
import { rememberLocalImport } from './local-projects';

type Phase = 'idle' | 'packing' | 'uploading' | 'done';

const NEW_PROJECT_VALUE = '__new__';

interface PickedFolder extends FolderInspection {
  entries: FolderEntry[];
}

export function FolderUpload({ fixedProjectId }: { fixedProjectId?: string }) {
  const projectsQuery = useListAllProjects();
  const { data: organizations = [] } = useListOrganizations();
  const activeOrgId = useOrganizationStore((s) => s.activeOrganizationId);
  const { mutateAsync: createProject } = useCreateProject();
  const deployArchiveMutation = useDeployLibraryArchive();
  const [projectId, setProjectId] = useState(fixedProjectId || NEW_PROJECT_VALUE);
  const [deployedProjectId, setDeployedProjectId] = useState(fixedProjectId || '');
  const [name, setName] = useState('');
  const [archive, setArchive] = useState<File | null>(null);
  const [folder, setFolder] = useState<PickedFolder | null>(null);
  const [phase, setPhase] = useState<Phase>('idle');
  const [error, setError] = useState('');
  const [dragging, setDragging] = useState(false);
  const archiveRef = useRef<HTMLInputElement>(null);
  const folderRef = useRef<HTMLInputElement>(null);
  const dropRef = useRef<HTMLDivElement>(null);
  const handlersRef = useRef({
    handleArchiveFiles: (_list: FileList | null) => {},
    handleFolderFiles: (_list: FileList | null): Promise<void> => Promise.resolve(),
  });
  const projects = projectsQuery.data ?? [];

  useEffect(() => {
    const node = dropRef.current;
    if (!node) return;
    const over = (event: DragEvent) => {
      event.preventDefault();
      setDragging(true);
    };
    const leave = () => setDragging(false);
    const drop = (event: DragEvent) => {
      event.preventDefault();
      setDragging(false);
      const files = event.dataTransfer?.files;
      if (!files?.length) return;
      if (/\.t(ar\.)?gz$/i.test(files[0].name)) {
        handlersRef.current.handleArchiveFiles(files);
      } else {
        void handlersRef.current.handleFolderFiles(files);
      }
    };
    node.addEventListener('dragover', over);
    node.addEventListener('dragleave', leave);
    node.addEventListener('drop', drop);
    return () => {
      node.removeEventListener('dragover', over);
      node.removeEventListener('dragleave', leave);
      node.removeEventListener('drop', drop);
    };
  }, []);

  const clear = () => {
    setArchive(null);
    setFolder(null);
    setError('');
  };

  const handleArchiveFiles = (list: FileList | null) => {
    const file = list?.[0];
    if (!file) return;
    if (!/\.t(ar\.)?gz$/i.test(file.name) && !/\.tgz$/i.test(file.name)) {
      setError('Upload a .tar.gz archive, or pick a folder to pack one automatically.');
      return;
    }
    clear();
    setArchive(file);
    setName(file.name.replace(/\.tar\.gz$/i, '').replace(/\.tgz$/i, ''));
  };

  const handleFolderFiles = async (list: FileList | null) => {
    if (!list || list.length === 0) return;
    const entries = collectFolderFiles(list);
    if (entries.length === 0) {
      setError('The selected folder has no uploadable files.');
      return;
    }
    const inspection = await inspectFolderFiles(entries);
    clear();
    setFolder({ ...inspection, entries });
    setName(inspection.name);
  };

  handlersRef.current = { handleArchiveFiles, handleFolderFiles };

  const handleDeploy = async () => {
    if ((!archive && !folder) || !name.trim() || phase === 'packing' || phase === 'uploading') {
      return;
    }
    setError('');
    try {
      let file = archive;
      if (!file && folder) {
        setPhase('packing');
        const blob = await packTarGz(folder.entries);
        file = new File([blob], `${name.trim()}.tar.gz`, { type: 'application/gzip' });
      }
      if (!file) return;
      setPhase('uploading');

      let targetProjectId = fixedProjectId || projectId;
      if (!targetProjectId || targetProjectId === NEW_PROJECT_VALUE) {
        const orgId = activeOrgId || organizations[0]?.id || '';
        const created = await createProject({
          payload: { name: name.trim(), organizationId: orgId },
        });
        targetProjectId = created.data.id;
      }

      const response = await deployArchiveMutation.mutateAsync({
        projectId: targetProjectId,
        name: name.trim(),
        file,
      });
      setDeployedProjectId(targetProjectId);
      if (folder) {
        rememberLocalImport({
          name: name.trim(),
          root: name.trim(),
          fileCount: folder.fileCount,
          framework: folder.framework,
          packageManager: folder.packageManager,
        });
      }
      setPhase('done');
      toast.success(response.message || 'Archive deployed');
    } catch (err) {
      setPhase('idle');
      setError(err instanceof Error ? err.message : 'Archive upload failed');
    }
  };

  const busy = phase === 'packing' || phase === 'uploading';

  if (phase === 'done') {
    return (
      <FolderUploadDone
        name={name}
        projectId={fixedProjectId || deployedProjectId}
        onReset={() => {
          clear();
          setPhase('idle');
        }}
      />
    );
  }

  return (
    <div
      id="library-folder-upload"
      className="scroll-mt-6 rounded-2xl border border-border/50 bg-card"
    >
      <div className="flex items-center gap-3 border-border/50 border-b px-5 py-4">
        <div className="flex h-9 w-9 items-center justify-center rounded-xl bg-muted/60">
          <FolderUp className="size-4.5 text-muted-foreground" />
        </div>
        <div>
          <h2 className="font-semibold text-[15px] text-foreground">Upload a folder</h2>
          <p className="text-muted-foreground text-xs">
            Pack a local folder in your browser and deploy it as an archive.
          </p>
        </div>
      </div>
      <div className="space-y-4 px-5 py-4">
        {!archive && !folder ? (
          <FolderDropzone
            dropRef={dropRef}
            archiveRef={archiveRef}
            folderRef={folderRef}
            dragging={dragging}
            busy={busy}
            onArchiveFiles={handleArchiveFiles}
            onFolderFiles={handleFolderFiles}
          />
        ) : (
          <div className="space-y-4">
            <div className="flex items-center gap-3 rounded-xl border border-border/50 bg-muted/40 px-3 py-3">
              <div className="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-muted">
                <Package className="size-4.5 text-muted-foreground" />
              </div>
              <div className="min-w-0 flex-1">
                <p className="truncate font-medium text-foreground text-sm">
                  {folder ? folder.name : (archive?.name ?? '')}
                </p>
                <p className="text-muted-foreground text-xs">
                  {folder
                    ? `${folder.fileCount} files · ${folder.framework} · ${folder.packageManager}`
                    : 'Ready to upload'}
                </p>
              </div>
              {!busy && (
                <button
                  type="button"
                  onClick={clear}
                  aria-label="Clear selection"
                  className="rounded-lg p-1.5 text-muted-foreground/50 transition-colors hover:bg-muted hover:text-foreground"
                >
                  <X className="size-4" />
                </button>
              )}
            </div>

            <div className={fixedProjectId ? 'space-y-2' : 'grid gap-4 sm:grid-cols-2'}>
              <div className="space-y-2">
                <Label htmlFor="folder-name">Application name *</Label>
                <Input
                  id="folder-name"
                  value={name}
                  disabled={busy}
                  onChange={(event) => setName(event.target.value)}
                  placeholder="my-app"
                />
              </div>
              {!fixedProjectId && (
                <div className="space-y-2">
                  <Label htmlFor="folder-project">Destination project</Label>
                  <Select value={projectId} onValueChange={setProjectId} disabled={busy}>
                    <SelectTrigger id="folder-project" className="w-full">
                      <SelectValue placeholder="Create new project" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value={NEW_PROJECT_VALUE}>
                        <span className="inline-flex items-center gap-1.5">
                          <FolderPlus className="size-3.5 text-primary" />
                          Create new project ({name.trim() || 'my-app'})
                        </span>
                      </SelectItem>
                      {projects.map((project) => (
                        <SelectItem key={project.id} value={project.id}>
                          {project.name}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>
              )}
            </div>

            <Button
              onClick={() => void handleDeploy()}
              disabled={busy || !name.trim()}
              className="w-full"
            >
              {busy ? (
                <Loader2 className="size-4 animate-spin" />
              ) : (
                <ArrowRight className="size-4" />
              )}
              {phase === 'packing'
                ? 'Packing archive…'
                : phase === 'uploading'
                  ? 'Uploading…'
                  : 'Deploy archive'}
            </Button>
          </div>
        )}

        {error && (
          <p role="alert" className="text-destructive text-xs">
            {error}
          </p>
        )}
      </div>
    </div>
  );
}
