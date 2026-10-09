import { useEffect, useState } from 'react';

const DEMO_MODE_KEY = 'codedock_demo_mode';

export function isDemoModeEnabled(): boolean {
  if (typeof window === 'undefined') return false;
  return localStorage.getItem(DEMO_MODE_KEY) === 'true';
}

export function setDemoModeEnabled(enabled: boolean): void {
  if (typeof window === 'undefined') return;
  localStorage.setItem(DEMO_MODE_KEY, String(enabled));
  if (enabled) {
    document.documentElement.classList.add('demo-mode');
  } else {
    document.documentElement.classList.remove('demo-mode');
  }
  window.dispatchEvent(new Event('demo-mode-change'));
}

export function useDemoMode(): [boolean, (val: boolean) => void] {
  const [enabled, setEnabled] = useState(isDemoModeEnabled);

  useEffect(() => {
    const handler = () => setEnabled(isDemoModeEnabled());
    window.addEventListener('demo-mode-change', handler);
    return () => window.removeEventListener('demo-mode-change', handler);
  }, []);

  const update = (next: boolean) => {
    setDemoModeEnabled(next);
    setEnabled(next);
  };

  return [enabled, update];
}
