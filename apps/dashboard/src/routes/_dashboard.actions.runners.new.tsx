import { createFileRoute, useNavigate } from '@tanstack/react-router';
import { RunnerEditor } from '#/features/actions/runner-editor';

export const Route = createFileRoute('/_dashboard/actions/runners/new')({
  component: AddRunnerPage,
});

function AddRunnerPage() {
  const navigate = useNavigate();

  return (
    <RunnerEditor
      onBack={() => navigate({ to: '/actions' as never })}
      onClose={() => navigate({ to: '/actions' as never })}
    />
  );
}
