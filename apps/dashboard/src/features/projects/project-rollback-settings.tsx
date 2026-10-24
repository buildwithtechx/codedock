import { Loader2 } from 'lucide-react';
import { RollbackRetentionCards } from '#/components/rollback/rollback-retention-cards';
import { useRollbackSettings } from '#/components/rollback/use-rollback-settings';

export function ProjectRollbackSettings({ projectId }: { projectId: string }) {
  const { capacity, loading, error, savingWindow, togglingStrategy, changeWindow, toggleStrategy } =
    useRollbackSettings(projectId);

  return (
    <div className="space-y-4">
      <div>
        <h2 className="font-semibold text-lg">Rollback & Retention</h2>
        <p className="text-muted-foreground text-sm">
          Configure how previous deployment artifacts are retained for instant rollback or storage
          reclamation.
        </p>
      </div>

      {loading && !capacity ? (
        <div className="flex h-32 items-center justify-center">
          <Loader2 className="size-5 animate-spin text-muted-foreground" />
        </div>
      ) : error && !capacity ? (
        <div className="rounded-lg border border-destructive/20 bg-destructive/10 p-4 text-destructive text-sm">
          {error}
        </div>
      ) : (
        <RollbackRetentionCards
          strategy={capacity?.strategy ?? 'snapshot'}
          capacity={capacity}
          onToggleStrategy={toggleStrategy}
          onChangeWindow={changeWindow}
          togglingStrategy={togglingStrategy}
          savingWindow={savingWindow}
        />
      )}
    </div>
  );
}

export default ProjectRollbackSettings;
