import { Input } from '@codedock/ui/input';
import { Info, Lock } from 'lucide-react';
import type { ReactNode } from 'react';
import type { ServerSettings } from '#/interfaces/settings';
import { PreferenceRow } from './preference-row';
import { SettingsRow } from './settings-row';
import { SettingsSection } from './settings-section';

export type GeneralFields = Pick<
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

interface NetworkSecurityProps {
  form: GeneralFields;
  set: <K extends keyof GeneralFields>(key: K, val: GeneralFields[K]) => void;
  renderSaveFooter: (
    sectionName: string,
    buttonLabel: string,
    keys: (keyof GeneralFields)[]
  ) => ReactNode;
}

export function ServerNetworkSection({ form, set, renderSaveFooter }: NetworkSecurityProps) {
  return (
    <SettingsSection
      collapsible
      icon={<Info className="h-4 w-4" />}
      title="Network"
      description="Public IP addresses, Traefik edge routing, and access control allowlists."
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
      <SettingsRow label="Traefik Wildcard IP" description="IP Traefik routes wildcard domains to.">
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
        bordered={false}
      >
        <Input
          value={form.ipAllowlist ?? ''}
          onChange={(e) => set('ipAllowlist', e.target.value)}
          placeholder="0.0.0.0/0"
          className="font-mono text-xs"
        />
      </SettingsRow>
      {renderSaveFooter('Network', 'Save Network', [
        'publicIpv4',
        'publicIpv6',
        'traefikWildcardIp',
        'ipAllowlist',
      ])}
    </SettingsSection>
  );
}

export function ServerSecuritySection({ form, set, renderSaveFooter }: NetworkSecurityProps) {
  return (
    <SettingsSection
      collapsible
      icon={<Lock className="h-4 w-4" />}
      title="Security & Access"
      description="User registration policies, domain restrictions, and safety confirmations."
    >
      <div className="space-y-4">
        <PreferenceRow
          title="User Registration"
          hint="Allow new users to sign up without an explicit invite."
          checked={form.registrationEnabled}
          onChange={(v) => set('registrationEnabled', v)}
        />
        <SettingsRow
          label="Registration Domain Allowlist"
          description="Comma-separated domains allowed to register (e.g. acme.com)."
          bordered={false}
        >
          <Input
            value={form.registrationDomainAllowlist ?? ''}
            onChange={(e) => set('registrationDomainAllowlist', e.target.value)}
            placeholder="company.com, partner.com"
            className="font-mono text-xs"
          />
        </SettingsRow>
        <PreferenceRow
          title="Two-Step Confirmation"
          hint="Require confirmation dialog before destructive actions."
          checked={!form.disableTwoStepConfirmation}
          onChange={(v) => set('disableTwoStepConfirmation', !v)}
        />
      </div>
      {renderSaveFooter('Security & Access', 'Save Security & Access', [
        'registrationEnabled',
        'registrationDomainAllowlist',
        'disableTwoStepConfirmation',
      ])}
    </SettingsSection>
  );
}
