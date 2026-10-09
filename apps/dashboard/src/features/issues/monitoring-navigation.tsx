import { Tabs, TabsList, TabsTrigger } from '#/components/ui/tabs';
import type { MonitoringTab } from './interfaces';

const NAV_TABS: Array<{ value: MonitoringTab; label: string }> = [
  { value: 'open', label: 'Overview' },
  { value: 'health', label: 'Health' },
  { value: 'resolved', label: 'History' },
];

export function MonitoringNavigation({
  value,
  onChange,
}: {
  value: MonitoringTab;
  onChange: (tab: MonitoringTab) => void;
}) {
  return (
    <Tabs
      value={value}
      onValueChange={(next) => onChange(next as MonitoringTab)}
      className="w-full"
    >
      <TabsList variant="line" className="w-full justify-start gap-1 border-border/60 border-b">
        {NAV_TABS.map((tab) => (
          <TabsTrigger key={tab.value} value={tab.value} className="px-3">
            {tab.label}
          </TabsTrigger>
        ))}
      </TabsList>
    </Tabs>
  );
}
