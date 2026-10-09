import { Link } from '@tanstack/react-router';
import { ArrowRight, FolderUp, Loader2, Package, X } from 'lucide-react';
import { useEffect, useRef, useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '#/components/ui/select';
import { useListAllProjects } from '#/features/projects';
import { libraryApi } from './api';
import {
  collectFolderFiles,
  type FolderEntry,
  type FolderInspection,
  inspectFolderFiles,
  packTarGz,
} from './folder-pack';
import { rememberLocalImport } from './local-projects';

type Phase = 'idle' | 'packing' | 'uploading' | 'done';

interface PickedFolder extends FolderInspection {
  entries: FolderEntry[];
}

export function FolderUpload() {
  const projectsQuery = useListAllProjects();
  const [projectId, setProjectId] = useState('');
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
      if (!files || files.length === 0) return;
      const [first] = Array.from(files);
      if (first && /\.t(ar\.)?gz$/i.test(first.name)) {
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
    if (
      (!archive && !folder) ||
      !projectId ||
      !name.trim() ||
      phase === 'packing' ||
      phase === 'uploading'
    ) {
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
      const response = await libraryApi.deployArchive(projectId, name.trim(), file);
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
      <div className="rounded-2xl border border-border/50 bg-card p-8 text-center">
        <div className="mx-auto mb-4 flex h-14 w-14 items-center justify-center rounded-2xl bg-emerald-500/10">
          <Package className="size-7 text-emerald-500" />
        </div>
        <h3 className="font-semibold text-base text-foreground">Archive deployed</h3>
        <p className="mx-auto mt-1 max-w-sm text-muted-foreground text-sm">
          {name} is deploying in your project.
        </p>
        <div className="mt-6 flex flex-wrap items-center justify-center gap-3">
          <Button asChild>
            <Link to="/projects/$projectId" params={{ projectId }}>
              Open project
            </Link>
          </Button>
          <Button
            variant="outline"
            onClick={() => {
              clear();
              setPhase('idle');
            }}
          >
            Upload another
          </Button>
        </div>
      </div>
    );
  }

  return (
    <div
      id="library-folder-upload"
      className="scroll-mt-6 rounded-2xl border border-border/50 bg-card"
    >
      <div className="flex items-center gap-3 border-border/50 border-b px-5 py-4">
        <div className="flex h-9 w-9 items-center justify-center rounded-xl bg-muted/60">
          <FolderUp className="size-[18px] text-muted-foreground" />
        </div>
        <div>
          <h2 className="font-semibold text-[15px] text-foreground">Upload a folder</h2>
          <p className="text-muted-foreground text-xs">
            Pack a local folder in your browser and deploy it as an archive.
          </p>
        </div>
      </div>
      <div className="space-y-4 px-5 py-4">
        <div className="grid gap-4 sm:grid-cols-2">
          <div className="space-y-2">
            <Label htmlFor="folder-project">Project *</Label>
            <Select value={projectId} onValueChange={setProjectId} disabled={busy}>
              <SelectTrigger id="folder-project" className="w-full">
                <SelectValue placeholder="Select a project" />
              </SelectTrigger>
              <SelectContent>
                {projects.map((project) => (
                  <SelectItem key={project.id} value={project.id}>
                    {project.name}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
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
        </div>

        {!archive && !folder ? (
          <div
            ref={dropRef}
            className={`rounded-xl border-2 border-dashed p-8 text-center transition-all ${
              dragging
                ? 'border-primary bg-primary/5'
                : 'border-border/60 hover:border-border hover:bg-muted/30'
            }`}
          >
            <input
              ref={archiveRef}
              type="file"
              className="hidden"
              accept=".tar.gz,.tgz"
              onChange={(event) => handleArchiveFiles(event.target.files)}
            />
            <input
              ref={(element) => {
                folderRef.current = element;
                element?.setAttribute('webkitdirectory', '');
              }}
              type="file"
              className="hidden"
              multiple
              onChange={(event) => void handleFolderFiles(event.target.files)}
            />
            <p className="font-medium text-foreground text-sm">
              Drop a folder or a .tar.gz archive
            </p>
            <p className="mt-1 text-muted-foreground text-xs">
              Folders are packed into a .tar.gz archive locally before upload.
            </p>
            <div className="mt-4 flex flex-wrap items-center justify-center gap-2">
              <Button
                type="button"
                variant="outline"
                size="sm"
                disabled={busy}
                onClick={() => folderRef.current?.click()}
              >
                Pick a folder
              </Button>
              <Button
                type="button"
                variant="outline"
                size="sm"
                disabled={busy}
                onClick={() => archiveRef.current?.click()}
              >
                Pick an archive
              </Button>
            </div>
          </div>
        ) : (
          <div className="space-y-4">
            <div className="flex items-center gap-3 rounded-xl border border-border/50 bg-muted/40 px-3 py-3">
              <div className="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-muted">
                <Package className="size-[18px] text-muted-foreground" />
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
            <Button
              onClick={() => void handleDeploy()}
              disabled={busy || !projectId || !name.trim()}
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
