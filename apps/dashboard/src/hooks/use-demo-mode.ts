import { useEffect, useState } from 'react';
import { DEMO_MODE_CLASS, DEMO_MODE_EVENT } from '#/constants/demo-mode';
import { isDemoModeEnabled, setDemoModeEnabled } from '#/services/demo-mode';

export function useDemoMode(): [boolean, (val: boolean) => void] {
  const [enabled, setEnabled] = useState(isDemoModeEnabled);

  useEffect(() => {
    if (enabled) {
      document.documentElement.classList.add(DEMO_MODE_CLASS);
    }
    const handler = () => {
      const isEnabled = isDemoModeEnabled();
      setEnabled(isEnabled);
      if (isEnabled) {
        document.documentElement.classList.add(DEMO_MODE_CLASS);
      } else {
        document.documentElement.classList.remove(DEMO_MODE_CLASS);
      }
    };
    window.addEventListener(DEMO_MODE_EVENT, handler);
    return () => window.removeEventListener(DEMO_MODE_EVENT, handler);
  }, [enabled]);

  const update = (next: boolean) => {
    setDemoModeEnabled(next);
    setEnabled(next);
  };

  return [enabled, update];
}
