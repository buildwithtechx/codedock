import { Loader2 } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
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
import { Switch } from '#/components/ui/switch';
import { actionsService, type ActionRunner } from '#/services/actions';

interface RunnerEditorProps {
  open: boolean;
  onClose: () => void;
  onSuccess: (runner: ActionRunner) => void;
}

export function RunnerEditor({ open, onClose, onSuccess }: RunnerEditorProps) {
  const [name, setName] = useState('');
  const [mode, setMode] = useState<'container' | 'native'>('container');
  const [image, setImage] = useState('catthehacker/ubuntu:act-22.04');
  const [labels, setLabels] = useState('ubuntu-latest, docker');
  const [enabled, setEnabled] = useState(true);
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim()) {
      toast.error('Runner name is required');
      return;
    }

    setIsSubmitting(true);
    try {
      const parsedLabels = labels
        .split(',')
        .map((l) => l.trim())
        .filter(Boolean);

      const runner = await actionsService.createRunner({
        name: name.trim(),
        mode,
        image: mode === 'container' ? image.trim() : undefined,
        labels: parsedLabels,
        enabled,
      });

      toast.success(`Runner ${runner.name} registered successfully`);
      onSuccess(runner);
    } catch {
      toast.error('Failed to register runner');
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <Dialog open={open} onOpenChange={(val) => !val && onClose()}>
      <DialogContent className="sm:max-w-[480px]">
        <form onSubmit={handleSubmit}>
          <DialogHeader>
            <DialogTitle>Register Action Runner</DialogTitle>
            <DialogDescription>
              Connect a container or host environment to process queued workflow runs.
            </DialogDescription>
          </DialogHeader>

          <div className="grid gap-4 py-4">
            <div className="space-y-1.5">
              <Label htmlFor="runner-name">Runner Name</Label>
              <Input
                id="runner-name"
                placeholder="e.g. Ubuntu High-CPU Runner"
                value={name}
                onChange={(e) => setName(e.target.value)}
                required
              />
            </div>

            <div className="space-y-1.5">
              <Label>Execution Mode</Label>
              <div className="flex gap-4">
                <label className="flex items-center gap-2 text-sm font-medium cursor-pointer">
                  <input
                    type="radio"
                    name="mode"
                    value="container"
                    checked={mode === 'container'}
                    onChange={() => setMode('container')}
                    className="accent-primary"
                  />
                  <span>Container (Docker)</span>
                </label>
                <label className="flex items-center gap-2 text-sm font-medium cursor-pointer">
                  <input
                    type="radio"
                    name="mode"
                    value="native"
                    checked={mode === 'native'}
                    onChange={() => setMode('native')}
                    className="accent-primary"
                  />
                  <span>Native Host</span>
                </label>
              </div>
            </div>

            {mode === 'container' && (
              <div className="space-y-1.5">
                <Label htmlFor="runner-image">Container Image</Label>
                <Input
                  id="runner-image"
                  placeholder="e.g. catthehacker/ubuntu:act-22.04"
                  value={image}
                  onChange={(e) => setImage(e.target.value)}
                />
                <p className="text-[11px] text-muted-foreground">
                  Official Act or Docker image equipped with runner binaries and node runtimes.
                </p>
              </div>
            )}

            <div className="space-y-1.5">
              <Label htmlFor="runner-labels">Labels (comma-separated)</Label>
              <Input
                id="runner-labels"
                placeholder="e.g. ubuntu-latest, linux, x86_64"
                value={labels}
                onChange={(e) => setLabels(e.target.value)}
              />
            </div>

            <div className="flex items-center justify-between rounded-lg border border-border p-3">
              <div className="space-y-0.5">
                <Label htmlFor="runner-enabled" className="text-sm font-medium">
                  Enable Runner
                </Label>
                <p className="text-xs text-muted-foreground">
                  Accept incoming workflow jobs immediately upon registration.
                </p>
              </div>
              <Switch
                id="runner-enabled"
                checked={enabled}
                onCheckedChange={setEnabled}
              />
            </div>
          </div>

          <DialogFooter>
            <Button type="button" variant="outline" onClick={onClose} disabled={isSubmitting}>
              Cancel
            </Button>
            <Button type="submit" disabled={isSubmitting}>
              {isSubmitting && <Loader2 className="mr-2 size-4 animate-spin" />}
              Register Runner
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
