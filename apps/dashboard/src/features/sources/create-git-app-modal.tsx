import { ApplicationSetupModal } from './application-setup-modal';
import type { ApplicationSetupProps } from './application-setup-types';

export function CreateGitAppModal(props: ApplicationSetupProps) {
  return props.isOpen ? <ApplicationSetupModal {...props} initialSource="git" /> : null;
}
