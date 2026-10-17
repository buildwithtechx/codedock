import { createFileRoute, useNavigate } from '@tanstack/react-router';
import { RunnerDestination } from '#/features/actions/runner-destination';
import { RunnerEditor } from '#/features/actions/runner-editor';

export const Route = createFileRoute('/_dashboard/actions/runners/new')({
  component: AddRunnerPage,
});

function AddRunnerPage() {
  const navigate = useNavigate();

  return (
    <RunnerDestination
      onCancel={() => navigate({ to: '/actions' as never })}
      renderForm={(server, back) => (
        <RunnerEditor
          initialServer={server}
          onBack={back ?? (() => navigate({ to: '/actions' as never }))}
          onClose={() => navigate({ to: '/actions' as never })}
        />
      )}
    />
  );
}
