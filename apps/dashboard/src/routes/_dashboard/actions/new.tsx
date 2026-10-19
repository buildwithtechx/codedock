import { createFileRoute } from '@tanstack/react-router';
import { WorkflowBuilder } from '#/features/actions/workflow-builder';

export const Route = createFileRoute('/_dashboard/actions/new')({
  component: NewWorkflowRoute,
});

function NewWorkflowRoute() {
  return <WorkflowBuilder />;
}
