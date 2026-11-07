import type { RefObject } from 'react';
import { Button } from '@codedock/ui/button';

interface FolderDropzoneProps {
  dropRef: RefObject<HTMLDivElement | null>;
  archiveRef: RefObject<HTMLInputElement | null>;
  folderRef: { current: HTMLInputElement | null };
  dragging: boolean;
  busy: boolean;
  onArchiveFiles: (files: FileList | null) => void;
  onFolderFiles: (files: FileList | null) => void;
}

export function FolderDropzone({
  dropRef,
  archiveRef,
  folderRef,
  dragging,
  busy,
  onArchiveFiles,
  onFolderFiles,
}: FolderDropzoneProps) {
  return (
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
        onChange={(event) => onArchiveFiles(event.target.files)}
      />
      <input
        ref={(element) => {
          folderRef.current = element;
          element?.setAttribute('webkitdirectory', '');
        }}
        type="file"
        className="hidden"
        multiple
        onChange={(event) => void onFolderFiles(event.target.files)}
      />
      <p className="font-medium text-foreground text-sm">Drop a folder or a .tar.gz archive</p>
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
  );
}
