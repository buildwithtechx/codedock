import { createFileRoute, useNavigate } from '@tanstack/react-router';
import { RunnerDestination } from '#/features/actions/runner-destination';
import { RunnerEditor } from '#/features/actions/runner-editor';
import { useGetPublicSettings } from '#/features/settings';

export const Route = createFileRoute('/_dashboard/actions/runners/new')({
  component: AddRunnerPage,
});

function AddRunnerPage() {
  const navigate = useNavigate();
  const { data: publicRes } = useGetPublicSettings();
  const isCloud = publicRes?.data?.cloudMode ?? false;

  return (
    <RunnerDestination
      selfHosted={!isCloud}
      onCancel={() => navigate({ to: '/actions' as never })}
      renderForm={(server, back) => (
        <RunnerEditor
          initialServer={server}
          defaultCloud={isCloud}
          onBack={back ?? (() => navigate({ to: '/actions' as never }))}
          onClose={() => navigate({ to: '/actions' as never })}
        />
      )}
    />
  );
}
