import { Button } from '@codedock/ui/button';
import { Input } from '@codedock/ui/input';
import { Skeleton } from '@codedock/ui/skeleton';
import { Check, Clock, Cpu, Globe } from 'lucide-react';
import { useEffect, useState } from 'react';
import { toast } from 'sonner';
import { useGetSettings, useUpdateSettings } from '#/hooks/use-settings';
import type { ServerSettings } from '#/interfaces/settings';
import { PreferenceRow } from './preference-row';
import {
  type GeneralFields,
  ServerNetworkSection,
  ServerSecuritySection,
} from './server-network-security-sections';
import { SettingsRow } from './settings-row';
import { SettingsSection } from './settings-section';

const EMPTY: GeneralFields = {
  siteName: '',
  dashboardDomain: '',
  defaultWildcardDomain: '',
  publicIpv4: '',
  publicIpv6: '',
  traefikWildcardIp: '',
  registrationEnabled: true,
  registrationDomainAllowlist: '',
  ipAllowlist: '',
  mcpServerEnabled: false,
  disableTwoStepConfirmation: false,
  telemetryEnabled: true,
  concurrentBuilds: 2,
  deploymentTimeout: 300,
  serverTimezone: 'UTC',
};

export function ServerGeneralSettings() {
  const { data, isLoading } = useGetSettings();
  const { mutateAsync: updateSettings, isPending } = useUpdateSettings();

  const [form, setForm] = useState<GeneralFields>(EMPTY);
  const [savingSection, setSavingSection] = useState<string | null>(null);

  useEffect(() => {
    if (data?.data) {
      const s = data.data;
      setForm({
        siteName: s.siteName ?? '',
        dashboardDomain: s.dashboardDomain ?? '',
        defaultWildcardDomain: s.defaultWildcardDomain ?? '',
        publicIpv4: s.publicIpv4 ?? '',
        publicIpv6: s.publicIpv6 ?? '',
        traefikWildcardIp: s.traefikWildcardIp ?? '',
        registrationEnabled: s.registrationEnabled,
        registrationDomainAllowlist: s.registrationDomainAllowlist ?? '',
        ipAllowlist: s.ipAllowlist ?? '',
        mcpServerEnabled: s.mcpServerEnabled,
        disableTwoStepConfirmation: s.disableTwoStepConfirmation,
        telemetryEnabled: s.telemetryEnabled,
        concurrentBuilds: s.concurrentBuilds,
        deploymentTimeout: s.deploymentTimeout,
        serverTimezone: s.serverTimezone ?? 'UTC',
      });
    }
  }, [data]);

  const set = <K extends keyof GeneralFields>(key: K, val: GeneralFields[K]) =>
    setForm((prev) => ({ ...prev, [key]: val }));

  const save = async (sectionName: string, keys: (keyof GeneralFields)[]) => {
    setSavingSection(sectionName);
    try {
      const payload: Partial<ServerSettings> = {};
      for (const k of keys) {
        (payload as Record<string, unknown>)[k] = form[k];
      }
      await updateSettings({ payload });
      toast.success(`${sectionName} saved`);
    } catch {
      toast.error(`Failed to save ${sectionName}`);
    } finally {
      setSavingSection(null);
    }
  };

  if (isLoading) {
    return (
      <div className="space-y-4">
        {[0, 1, 2, 3].map((i) => (
          <Skeleton key={i} className="h-32 w-full rounded-2xl" />
        ))}
      </div>
    );
  }

  const renderSaveFooter = (
    sectionName: string,
    buttonLabel: string,
    keys: (keyof GeneralFields)[]
  ) => (
    <div className="flex justify-end border-border/40 border-t pt-4">
      <Button
        size="sm"
        onClick={() => save(sectionName, keys)}
        disabled={isPending && savingSection === sectionName}
        className="gap-1.5"
      >
        <Check className="h-3.5 w-3.5" />
        {savingSection === sectionName ? 'Saving...' : buttonLabel}
      </Button>
    </div>
  );

  return (
    <div className="space-y-6">
      <SettingsSection
        collapsible
        icon={<Globe className="h-4 w-4" />}
        title="Domains & Routing"
        description="Configure the primary hostname and root wildcard domains for routing."
      >
        <SettingsRow label="Site Name" description="Displayed in the browser tab and emails.">
          <Input
            value={form.siteName ?? ''}
            onChange={(e) => set('siteName', e.target.value)}
            placeholder="Codedock"
            className="text-sm"
          />
        </SettingsRow>
        <SettingsRow
          label="Dashboard Domain"
          description="The domain Codedock control panel is served from."
        >
          <Input
            value={form.dashboardDomain ?? ''}
            onChange={(e) => set('dashboardDomain', e.target.value)}
            placeholder="pilot.example.com"
            className="font-mono text-xs"
          />
        </SettingsRow>
        <SettingsRow
          label="Wildcard Domain"
          description="Root domain for generated app URLs (e.g. *.apps.example.com)."
          bordered={false}
        >
          <Input
            value={form.defaultWildcardDomain ?? ''}
            onChange={(e) => set('defaultWildcardDomain', e.target.value)}
            placeholder="apps.example.com"
            className="font-mono text-xs"
          />
        </SettingsRow>
        {renderSaveFooter('Domains & Routing', 'Save Domains & Routing', [
          'siteName',
          'dashboardDomain',
          'defaultWildcardDomain',
        ])}
      </SettingsSection>

      <ServerNetworkSection form={form} set={set} renderSaveFooter={renderSaveFooter} />

      <ServerSecuritySection form={form} set={set} renderSaveFooter={renderSaveFooter} />

      <SettingsSection
        collapsible
        icon={<Cpu className="h-4 w-4" />}
        title="Build & Deployment"
        description="Parallel build worker limits and deployment execution timeouts."
      >
        <SettingsRow label="Concurrent Builds" description="Max number of parallel build jobs.">
          <Input
            type="number"
            value={form.concurrentBuilds}
            onChange={(e) => set('concurrentBuilds', Number(e.target.value))}
            min={1}
            max={20}
            className="font-mono text-xs"
          />
        </SettingsRow>
        <SettingsRow
          label="Deployment Timeout (s)"
          description="Seconds before a deployment is considered failed."
          bordered={false}
        >
          <Input
            type="number"
            value={form.deploymentTimeout}
            onChange={(e) => set('deploymentTimeout', Number(e.target.value))}
            min={60}
            className="font-mono text-xs"
          />
        </SettingsRow>
        {renderSaveFooter('Build & Deployment', 'Save Build & Deployment', [
          'concurrentBuilds',
          'deploymentTimeout',
        ])}
      </SettingsSection>

      <SettingsSection
        collapsible
        icon={<Clock className="h-4 w-4" />}
        title="System"
        description="Instance timezone and anonymous diagnostics reporting."
      >
        <div className="space-y-4">
          <SettingsRow
            label="Server Timezone"
            description="Timezone used for cron jobs and logs."
            bordered={false}
          >
            <Input
              value={form.serverTimezone ?? ''}
              onChange={(e) => set('serverTimezone', e.target.value)}
              placeholder="UTC"
              className="font-mono text-xs"
            />
          </SettingsRow>
          <PreferenceRow
            title="Telemetry"
            hint="Send anonymous usage statistics to help improve Codedock."
            checked={form.telemetryEnabled}
            onChange={(v) => set('telemetryEnabled', v)}
          />
        </div>
        {renderSaveFooter('System', 'Save System', ['serverTimezone', 'telemetryEnabled'])}
      </SettingsSection>
    </div>
  );
}
