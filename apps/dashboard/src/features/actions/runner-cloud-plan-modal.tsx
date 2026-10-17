import { AddServerModal, type AddServerModalProps } from '#/features/servers/add-server-modal';
export { type CloudPlan } from '#/features/servers/add-server-cloud-plans';

export type RunnerCloudPlanModalProps = AddServerModalProps;

export function RunnerCloudPlanModal(props: RunnerCloudPlanModalProps) {
  return <AddServerModal {...props} />;
}
