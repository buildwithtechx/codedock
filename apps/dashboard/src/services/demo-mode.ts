import { DEMO_MODE_CLASS, DEMO_MODE_EVENT, DEMO_MODE_STORAGE_KEY } from '#/constants/demo-mode';

export function isDemoModeEnabled(): boolean {
  if (typeof window === 'undefined') return false;
  return localStorage.getItem(DEMO_MODE_STORAGE_KEY) === 'true';
}

export function setDemoModeEnabled(enabled: boolean): void {
  if (typeof window === 'undefined') return;
  localStorage.setItem(DEMO_MODE_STORAGE_KEY, String(enabled));
  if (enabled) {
    document.documentElement.classList.add(DEMO_MODE_CLASS);
  } else {
    document.documentElement.classList.remove(DEMO_MODE_CLASS);
  }
  window.dispatchEvent(new Event(DEMO_MODE_EVENT));
}

if (typeof window !== 'undefined' && isDemoModeEnabled()) {
  document.documentElement.classList.add(DEMO_MODE_CLASS);
}
