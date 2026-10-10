import { DnsSettings } from '#/features/dns/dns-settings';

export function CredentialsSettings() {
  return (
    <div className="space-y-6">
      <div>
        <h2 className="font-semibold text-foreground text-lg tracking-tight">Credentials</h2>
        <p className="text-muted-foreground text-sm">
          Manage DNS provider credentials and automated domain records.
        </p>
      </div>

      <DnsSettings />
    </div>
  );
}
