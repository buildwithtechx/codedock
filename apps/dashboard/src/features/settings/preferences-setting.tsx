import { Sliders } from 'lucide-react';
import { useEffect, useState } from 'react';
import { useDemoMode } from '#/lib/demo-mode';
import { PreferenceRow } from './preference-row';
import { SettingsSection } from './settings-section';

export function PreferencesSetting() {
  const [demoMode, setDemoMode] = useDemoMode();
  const [compactTables, setCompactTables] = useState(false);
  const [animations, setAnimations] = useState(true);

  useEffect(() => {
    const compactSaved = localStorage.getItem('codedock_compact_tables');
    if (compactSaved !== null) {
      setCompactTables(compactSaved === 'true');
    }
    const animSaved = localStorage.getItem('codedock_animations');
    if (animSaved !== null) {
      setAnimations(animSaved === 'true');
    }
  }, []);

  const handleToggleCompact = (val: boolean) => {
    setCompactTables(val);
    localStorage.setItem('codedock_compact_tables', String(val));
  };

  const handleToggleAnimations = (val: boolean) => {
    setAnimations(val);
    localStorage.setItem('codedock_animations', String(val));
  };

  return (
    <SettingsSection
      icon={<Sliders className="size-4 text-primary" />}
      title="Preferences"
      description="Personal display preferences for this browser — nothing is sent to the server."
      collapsible
    >
      <div className="space-y-3">
        <PreferenceRow
          title="Demo mode"
          hint="Blur IP addresses, hosts, and environment values during screen-shares and recordings."
          checked={demoMode}
          onChange={setDemoMode}
        />
        <PreferenceRow
          title="Compact table view"
          hint="Reduce row spacing and display more items per page across tables."
          checked={compactTables}
          onChange={handleToggleCompact}
        />
        <PreferenceRow
          title="Animated transitions"
          hint="Enable smooth UI transitions and chart animations."
          checked={animations}
          onChange={handleToggleAnimations}
        />
      </div>
    </SettingsSection>
  );
}
