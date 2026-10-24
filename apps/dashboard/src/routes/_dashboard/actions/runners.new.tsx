import { createFileRoute } from '@tanstack/react-router';
import { RunnerEditor } from '#/features/actions/runner-editor';

export const Route = createFileRoute('/_dashboard/actions/runners/new')({
  component: AddRunnerPage,
});

function AddRunnerPage() {
  return <RunnerEditor />;
}
