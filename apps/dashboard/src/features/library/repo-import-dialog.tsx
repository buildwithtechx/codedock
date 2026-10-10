import { GitFork } from 'lucide-react';
import { useEffect, useState } from 'react';
import { Button } from '#/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '#/components/ui/dialog';
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
import type { ImportTarget, PendingImport } from './types';

export function RepoImportDialog({
  pending,
  onClose,
  onConfirm,
}: {
  pending: PendingImport | null;
  onClose: () => void;
  onConfirm: (target: ImportTarget) => void;
}) {
  const projectsQuery = useListAllProjects();
  const [projectId, setProjectId] = useState('');
  const [name, setName] = useState('');
  const [branch, setBranch] = useState('');

  useEffect(() => {
    if (pending) {
      setProjectId('');
      setName(pending.name);
      setBranch(pending.branch);
    }
  }, [pending]);

  const projects = projectsQuery.data ?? [];

  return (
    <Dialog open={pending !== null} onOpenChange={(open) => !open && onClose()}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader icon={GitFork}>
          <DialogTitle>Import repository</DialogTitle>
          <DialogDescription>
            Choose where this repository deploys, then configure the application.
          </DialogDescription>
        </DialogHeader>
        {pending && (
          <div className="space-y-4">
            <p className="truncate rounded-lg bg-muted/60 px-3 py-2 font-mono text-xs">
              {pending.repositoryUrl}
            </p>
            <div className="space-y-2">
              <Label htmlFor="import-project">Project *</Label>
              <Select value={projectId} onValueChange={setProjectId}>
                <SelectTrigger id="import-project" className="w-full">
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
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="import-name">Application name</Label>
                <Input
                  id="import-name"
                  value={name}
                  onChange={(event) => setName(event.target.value)}
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="import-branch">Branch</Label>
                <Input
                  id="import-branch"
                  value={branch}
                  onChange={(event) => setBranch(event.target.value)}
                />
              </div>
            </div>
          </div>
        )}
        <DialogFooter>
          <Button variant="outline" onClick={onClose}>
            Cancel
          </Button>
          <Button
            disabled={!pending || !projectId || !branch.trim()}
            onClick={() => {
              if (pending && projectId) {
                onConfirm({
                  projectId,
                  repositoryUrl: pending.repositoryUrl,
                  name: name.trim() || pending.name,
                  branch: branch.trim() || pending.branch,
                });
              }
            }}
          >
            Continue
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
