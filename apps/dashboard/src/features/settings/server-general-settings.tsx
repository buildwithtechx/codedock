import { Check, Clock, Cpu, Globe, Info, Lock } from 'lucide-react';
import { useEffect, useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import { Input } from '#/components/ui/input';
import { Skeleton } from '#/components/ui/skeleton';
import { Switch } from '#/components/ui/switch';
import type { ServerSettings } from '#/features/settings';
import { useGetSettings, useUpdateSettings } from '#/features/settings';
import { SettingsRow } from './settings-row';
import { SettingsSection } from './settings-section';

type GeneralFields = Pick<
  ServerSettings,
  | 'siteName'
  | 'dashboardDomain'
  | 'defaultWildcardDomain'
  | 'publicIpv4'
  | 'publicIpv6'
  | 'traefikWildcardIp'
  | 'registrationEnabled'
  | 'registrationDomainAllowlist'
  | 'ipAllowlist'
  | 'mcpServerEnabled'
  | 'disableTwoStepConfirmation'
  | 'telemetryEnabled'
  | 'concurrentBuilds'
  | 'deploymentTimeout'
  | 'serverTimezone'
>;

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

  const saveAction = (sectionName: string, keys: (keyof GeneralFields)[]) => (
    <Button
      size="sm"
      variant="outline"
      onClick={() => save(sectionName, keys)}
      disabled={isPending && savingSection === sectionName}
      className="gap-1.5"
    >
      <Check className="h-3.5 w-3.5" />
      {savingSection === sectionName ? 'Saving...' : 'Save'}
    </Button>
  );

  return (
    <div className="space-y-6">
      <SettingsSection
        icon={<Globe className="h-4 w-4" />}
        title="Domains & Routing"
        action={saveAction('Domains', ['siteName', 'dashboardDomain', 'defaultWildcardDomain'])}
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
        >
          <Input
            value={form.defaultWildcardDomain ?? ''}
            onChange={(e) => set('defaultWildcardDomain', e.target.value)}
            placeholder="apps.example.com"
            className="font-mono text-xs"
          />
        </SettingsRow>
      </SettingsSection>

      <SettingsSection
        icon={<Info className="h-4 w-4" />}
        title="Network"
        action={saveAction('Network', [
          'publicIpv4',
          'publicIpv6',
          'traefikWildcardIp',
          'ipAllowlist',
        ])}
      >
        <SettingsRow
          label="Public IPv4"
          description="Server's public IPv4 address for DNS A records."
        >
          <Input
            value={form.publicIpv4 ?? ''}
            onChange={(e) => set('publicIpv4', e.target.value)}
            placeholder="1.2.3.4"
            className="font-mono text-xs"
          />
        </SettingsRow>
        <SettingsRow label="Public IPv6" description="Server's public IPv6 address (optional).">
          <Input
            value={form.publicIpv6 ?? ''}
            onChange={(e) => set('publicIpv6', e.target.value)}
            placeholder="2001:db8::1"
            className="font-mono text-xs"
          />
        </SettingsRow>
        <SettingsRow
          label="Traefik Wildcard IP"
          description="IP Traefik routes wildcard domains to."
        >
          <Input
            value={form.traefikWildcardIp ?? ''}
            onChange={(e) => set('traefikWildcardIp', e.target.value)}
            placeholder="1.2.3.4"
            className="font-mono text-xs"
          />
        </SettingsRow>
        <SettingsRow
          label="IP Allowlist"
          description="Comma-separated CIDRs that can access the control plane."
        >
          <Input
            value={form.ipAllowlist ?? ''}
            onChange={(e) => set('ipAllowlist', e.target.value)}
            placeholder="0.0.0.0/0"
            className="font-mono text-xs"
          />
        </SettingsRow>
      </SettingsSection>

      <SettingsSection
        icon={<Lock className="h-4 w-4" />}
        title="Security & Access"
        action={saveAction('Security', [
          'registrationEnabled',
          'registrationDomainAllowlist',
          'disableTwoStepConfirmation',
        ])}
      >
        <SettingsRow
          label="User Registration"
          description="Allow new users to sign up without an explicit invite."
        >
          <Switch
            checked={form.registrationEnabled}
            onCheckedChange={(v) => set('registrationEnabled', v)}
          />
        </SettingsRow>
        <SettingsRow
          label="Registration Domain Allowlist"
          description="Comma-separated domains allowed to register (e.g. acme.com)."
        >
          <Input
            value={form.registrationDomainAllowlist ?? ''}
            onChange={(e) => set('registrationDomainAllowlist', e.target.value)}
            placeholder="company.com, partner.com"
            className="font-mono text-xs"
          />
        </SettingsRow>
        <SettingsRow
          label="Two-Step Confirmation"
          description="Require confirmation dialog before destructive actions."
        >
          <Switch
            checked={!form.disableTwoStepConfirmation}
            onCheckedChange={(v) => set('disableTwoStepConfirmation', !v)}
          />
        </SettingsRow>
      </SettingsSection>

      <SettingsSection
        icon={<Cpu className="h-4 w-4" />}
        title="Build & Deployment"
        action={saveAction('Build & deployment', ['concurrentBuilds', 'deploymentTimeout'])}
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
        >
          <Input
            type="number"
            value={form.deploymentTimeout}
            onChange={(e) => set('deploymentTimeout', Number(e.target.value))}
            min={60}
            className="font-mono text-xs"
          />
        </SettingsRow>
      </SettingsSection>

      <SettingsSection
        icon={<Clock className="h-4 w-4" />}
        title="System"
        action={saveAction('System', ['serverTimezone', 'telemetryEnabled'])}
      >
        <SettingsRow label="Server Timezone" description="Timezone used for cron jobs and logs.">
          <Input
            value={form.serverTimezone ?? ''}
            onChange={(e) => set('serverTimezone', e.target.value)}
            placeholder="UTC"
            className="font-mono text-xs"
          />
        </SettingsRow>
        <SettingsRow
          label="Telemetry"
          description="Send anonymous usage statistics to help improve Codedock."
        >
          <Switch
            checked={form.telemetryEnabled}
            onCheckedChange={(v) => set('telemetryEnabled', v)}
          />
        </SettingsRow>
      </SettingsSection>
    </div>
  );
}
