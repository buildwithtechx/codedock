import { ArrowLeft, Check, ChevronRight, FileCode, Folder, Loader2, RefreshCw } from 'lucide-react';
import { useEffect, useMemo, useState } from 'react';
import { Button } from '@codedock/ui/button';
import { useRepoTree } from '#/hooks/use-git';
import type { RepoTreeEntry } from '#/interfaces/git';

export interface RepositoryPathBrowserProps {
  repositoryUrl: string;
  branch?: string;
  kind: 'directory' | 'file';
  value: string;
  onChange: (path: string) => void;
  disabled?: boolean;
}

interface TreeNode {
  name: string;
  path: string;
  isDirectory: boolean;
  children: TreeNode[];
}

function buildTree(entries: RepoTreeEntry[]): TreeNode[] {
  const root: TreeNode[] = [];
  const map = new Map<string, TreeNode>();

  for (const entry of entries) {
    const parts = entry.path.split('/');
    let currentPath = '';

    for (let i = 0; i < parts.length; i++) {
      const part = parts[i];
      const prevPath = currentPath;
      currentPath = prevPath ? `${prevPath}/${part}` : part;
      const isLast = i === parts.length - 1;
      const isDir = !isLast || entry.type === 'tree';

      if (!map.has(currentPath)) {
        const node: TreeNode = {
          name: part,
          path: currentPath,
          isDirectory: isDir,
          children: [],
        };
        map.set(currentPath, node);

        if (i === 0) {
          root.push(node);
        } else {
          const parent = map.get(prevPath);
          if (parent) {
            parent.isDirectory = true;
            if (!parent.children.some((c) => c.name === part)) {
              parent.children.push(node);
            }
          }
        }
      }
    }
  }

  const sortNodes = (nodes: TreeNode[]) => {
    nodes.sort((a, b) => {
      if (a.isDirectory !== b.isDirectory) {
        return a.isDirectory ? -1 : 1;
      }
      return a.name.localeCompare(b.name);
    });
    for (const n of nodes) {
      if (n.children.length > 0) sortNodes(n.children);
    }
  };

  sortNodes(root);
  return root;
}

export function RepositoryPathBrowser({
  repositoryUrl,
  branch,
  kind,
  value,
  onChange,
  disabled,
}: RepositoryPathBrowserProps) {
  const {
    data: treeRes,
    isLoading: loading,
    error: treeError,
    refetch,
  } = useRepoTree({
    repositoryUrl,
    branch,
  });

  const entries = treeRes?.data || null;
  const error =
    treeError instanceof Error
      ? treeError.message
      : treeError
        ? 'Failed to load repository tree'
        : null;

  const initialFolder = useMemo(() => {
    if (!value || value === '.' || value === '/') return '.';
    if (kind === 'directory') return value.replace(/^\/+|\/+$/g, '');
    const parts = value.split('/');
    return parts.slice(0, -1).join('/') || '.';
  }, [value, kind]);

  const [folder, setFolder] = useState<string>(initialFolder);

  useEffect(() => {
    setFolder(initialFolder);
  }, [initialFolder]);

  const tree = useMemo(() => buildTree(entries || []), [entries]);

  const visibleNodes = useMemo(() => {
    if (folder === '.' || folder === '') {
      return tree.filter((n) => (kind === 'file' ? true : n.isDirectory));
    }

    const segments = folder.split('/').filter(Boolean);
    let current = tree;

    for (const seg of segments) {
      const match = current.find((n) => n.name === seg && n.isDirectory);
      if (!match) return [];
      current = match.children;
    }

    return current.filter((n) => (kind === 'file' ? true : n.isDirectory));
  }, [tree, folder, kind]);

  const navigateUp = () => {
    if (folder === '.' || folder === '') return;
    const parts = folder.split('/').filter(Boolean);
    parts.pop();
    const nextFolder = parts.length > 0 ? parts.join('/') : '.';
    setFolder(nextFolder);
    if (kind === 'directory') {
      onChange(nextFolder === '.' ? '' : nextFolder);
    }
  };

  const handleSelectNode = (node: TreeNode) => {
    if (node.isDirectory) {
      setFolder(node.path);
      if (kind === 'directory') {
        onChange(node.path);
      }
    } else if (kind === 'file') {
      onChange(node.path);
    }
  };

  return (
    <div className="space-y-2 rounded-xl border border-border/60 bg-muted/20 p-3 text-xs">
      <div className="flex items-center justify-between border-b border-border/40 pb-2">
        <div className="flex items-center gap-1.5 font-mono text-[11px] text-muted-foreground">
          <Folder className="size-3.5 text-primary" />
          <span className="font-semibold text-foreground">
            {folder === '.' || folder === '' ? 'root' : folder}
          </span>
        </div>
        {loading && <Loader2 className="size-3.5 animate-spin text-muted-foreground" />}
      </div>

      {folder !== '.' && folder !== '' && (
        <Button
          type="button"
          variant="ghost"
          size="sm"
          onClick={navigateUp}
          disabled={disabled || loading}
          className="h-7 w-full justify-start gap-1.5 rounded-lg px-2 text-xs text-muted-foreground hover:text-foreground"
        >
          <ArrowLeft className="size-3" />
          <span>.. parent directory</span>
        </Button>
      )}

      {loading && !entries ? (
        <div className="flex items-center justify-center py-6 text-muted-foreground">
          <Loader2 className="size-5 animate-spin" />
        </div>
      ) : error ? (
        <div className="space-y-2 p-3 text-center">
          <p className="text-xs text-destructive">{error}</p>
          <Button
            type="button"
            variant="outline"
            size="sm"
            onClick={() => refetch()}
            className="h-7 gap-1.5 rounded-lg text-xs"
          >
            <RefreshCw className="size-3" />
            Retry
          </Button>
        </div>
      ) : visibleNodes.length === 0 ? (
        <p className="py-4 text-center text-xs text-muted-foreground">
          {kind === 'directory' ? 'No subdirectories found' : 'No matching files found'}
        </p>
      ) : (
        <div className="max-h-56 space-y-0.5 overflow-y-auto pr-1">
          {visibleNodes.map((node) => {
            const isSelected = value === node.path;
            return (
              <button
                key={node.path}
                type="button"
                disabled={disabled}
                onClick={() => handleSelectNode(node)}
                className={`flex w-full items-center justify-between rounded-lg px-2 py-1.5 text-left text-xs transition-colors ${
                  isSelected
                    ? 'bg-primary/10 text-primary font-medium'
                    : 'text-foreground hover:bg-muted/60'
                }`}
              >
                <div className="flex min-w-0 items-center gap-2">
                  {node.isDirectory ? (
                    <Folder className="size-3.5 shrink-0 text-primary/80" />
                  ) : (
                    <FileCode className="size-3.5 shrink-0 text-muted-foreground" />
                  )}
                  <span className="truncate font-mono">{node.name}</span>
                </div>
                <div className="flex shrink-0 items-center gap-1">
                  {isSelected && <Check className="size-3 text-primary" />}
                  {node.isDirectory && <ChevronRight className="size-3 text-muted-foreground" />}
                </div>
              </button>
            );
          })}
        </div>
      )}
    </div>
  );
}
