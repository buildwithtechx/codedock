import { ApplicationSetupModal } from './application-setup-modal';
import type { ApplicationSetupProps } from './application-setup-types';

export function CreateDockerImageModal(props: ApplicationSetupProps) {
  return props.isOpen ? <ApplicationSetupModal {...props} initialSource="image" /> : null;
}
