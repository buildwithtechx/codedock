export const HETZNER_WIZARD_STEPS = [
  { id: 'credential', label: 'Credential' },
  { id: 'plan', label: 'Location & type' },
  { id: 'review', label: 'Review' },
  { id: 'progress', label: 'Provision' },
] as const;

export type HetznerWizardStep = (typeof HETZNER_WIZARD_STEPS)[number]['id'];

export const CONNECTED_SERVER_STEPS = [
  { id: 'connection', label: 'Connection' },
  { id: 'credentials', label: 'Authentication' },
  { id: 'preflight', label: 'Preflight' },
] as const;

export type ConnectedServerStep = (typeof CONNECTED_SERVER_STEPS)[number]['id'];

export function nextStep<T extends string>(steps: readonly T[], current: T): T | null {
  const index = steps.indexOf(current);
  if (index === -1 || index >= steps.length - 1) {
    return null;
  }
  return steps[index + 1];
}

export function prevStep<T extends string>(steps: readonly T[], current: T): T | null {
  const index = steps.indexOf(current);
  if (index <= 0) {
    return null;
  }
  return steps[index - 1];
}

export function isFirstStep<T extends string>(steps: readonly T[], current: T): boolean {
  return steps.indexOf(current) === 0;
}

export function isLastStep<T extends string>(steps: readonly T[], current: T): boolean {
  return steps.indexOf(current) === steps.length - 1;
}
