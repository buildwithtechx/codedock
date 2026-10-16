import { Box, Cpu, HardDrive, Plus, Radio, Server } from 'lucide-react';
import { useState } from 'react';
import { Badge } from '#/components/ui/badge';
import { Button } from '#/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '#/components/ui/card';
import { type ActionRunner } from '#/services/actions';
import { RunnerEditor } from './runner-editor';

export { useCreateRunner } from './use-actions';

interface RunnerCatalogProps {
  runners: ActionRunner[];
  onRunnerRegistered?: (runner: ActionRunner) => void;
}

export function RunnerCatalog({ runners, onRunnerRegistered }: RunnerCatalogProps) {
  const [editorOpen, setEditorOpen] = useState(false);

  return (
    <div className="space-y-6">
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h3 className="text-lg font-semibold tracking-tight">Execution Runners</h3>
          <p className="text-xs text-muted-foreground">
            Manage containerized and native task runners for executing CI workflow jobs.
          </p>
        </div>
        <Button onClick={() => setEditorOpen(true)} size="sm" className="gap-1.5">
          <Plus className="size-4" />
          <span>Register Runner</span>
        </Button>
      </div>

      {runners.length === 0 ? (
        <Card className="border-dashed border-border/80">
          <CardContent className="flex flex-col items-center justify-center p-8 text-center">
            <Server className="size-8 text-muted-foreground/60 mb-2" />
            <h4 className="font-semibold text-sm">No runners registered</h4>
            <p className="text-xs text-muted-foreground max-w-sm mt-1 mb-4">
              Register a runner to start processing automated build and test pipelines.
            </p>
            <Button onClick={() => setEditorOpen(true)} size="sm" variant="outline" className="gap-1.5">
              <Plus className="size-4" />
              <span>Register Runner</span>
            </Button>
          </CardContent>
        </Card>
      ) : (
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {runners.map((runner) => (
            <Card key={runner.id} className="border-border/70 hover:border-border transition-colors">
              <CardHeader className="p-4 pb-2">
                <div className="flex items-start justify-between gap-2">
                  <div className="flex items-center gap-2">
                    {runner.mode === 'container' ? (
                      <Box className="size-4 text-blue-500" />
                    ) : (
                      <Cpu className="size-4 text-emerald-500" />
                    )}
                    <CardTitle className="text-sm font-semibold truncate">{runner.name}</CardTitle>
                  </div>
                  <Badge
                    variant={runner.status === 'online' ? 'success' : runner.status === 'busy' ? 'default' : 'secondary'}
                    className="capitalize text-[10px]"
                  >
                    {runner.status}
                  </Badge>
                </div>
                <CardDescription className="text-xs font-mono truncate">
                  {runner.mode === 'container' ? runner.image || 'container' : 'native host'}
                </CardDescription>
              </CardHeader>
              <CardContent className="p-4 pt-2 space-y-3">
                {runner.labels && runner.labels.length > 0 && (
                  <div className="flex flex-wrap gap-1">
                    {runner.labels.map((lbl) => (
                      <Badge key={lbl} variant="outline" className="text-[10px] px-1.5 py-0">
                        {lbl}
                      </Badge>
                    ))}
                  </div>
                )}
                <div className="flex items-center justify-between text-[11px] text-muted-foreground border-t border-border/50 pt-2">
                  <span className="flex items-center gap-1">
                    <Radio className="size-3 text-emerald-500" />
                    Heartbeat
                  </span>
                  <span>{runner.lastHeartbeatAt ? new Date(runner.lastHeartbeatAt).toLocaleTimeString() : 'N/A'}</span>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      )}

      {editorOpen && (
        <RunnerEditor
          open={editorOpen}
          onClose={() => setEditorOpen(false)}
          onSuccess={(runner) => {
            setEditorOpen(false);
            onRunnerRegistered?.(runner);
          }}
        />
      )}
    </div>
  );
}
