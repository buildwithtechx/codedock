import { CheckCircle2 } from 'lucide-react';
import { Button } from '#/components/ui/button';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';

interface DnsProviderFormProps {
  activeProvider: string;
  formData: Record<string, string>;
  setFormData: (data: Record<string, string>) => void;
  isPending: boolean;
  handleSaveProvider: (provider: string) => void;
  onCancel?: () => void;
}

const PROVIDER_NAMES: Record<string, string> = {
  cloudflare: 'Cloudflare',
  namecheap: 'Namecheap',
  spaceship: 'Spaceship',
};

export function DnsProviderForm({
  activeProvider,
  formData,
  setFormData,
  isPending,
  handleSaveProvider,
  onCancel,
}: DnsProviderFormProps) {
  const providerName = PROVIDER_NAMES[activeProvider] || 'Provider';

  return (
    <div className="space-y-5">
      {activeProvider === 'cloudflare' && (
        <div className="space-y-4">
          <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
            <div className="space-y-2">
              <Label
                htmlFor="cf-api-token"
                className="font-bold text-[10px] text-muted-foreground uppercase tracking-wider"
              >
                API KEY / TOKEN
              </Label>
              <Input
                id="cf-api-token"
                type="password"
                placeholder="Cloudflare API key or DNS token"
                value={formData.cloudflareApiToken}
                onChange={(e) =>
                  setFormData({
                    ...formData,
                    cloudflareApiToken: e.target.value,
                  })
                }
                className="bg-muted/30 font-medium"
              />
              <p className="text-[11px] text-muted-foreground">
                API Tokens only. Global API keys are not supported.
              </p>
            </div>
            <div className="space-y-2">
              <Label
                htmlFor="cf-email"
                className="font-bold text-[10px] text-muted-foreground uppercase tracking-wider"
              >
                ACCOUNT EMAIL
              </Label>
              <Input
                id="cf-email"
                placeholder="Only needed for global API keys"
                value={formData.cloudflareEmail}
                onChange={(e) =>
                  setFormData({
                    ...formData,
                    cloudflareEmail: e.target.value,
                  })
                }
                className="bg-muted/30 font-medium"
              />
            </div>
          </div>
          <div className="space-y-2">
            <Label
              htmlFor="cf-zone-id"
              className="font-bold text-[10px] text-muted-foreground uppercase tracking-wider"
            >
              ZONE ID
            </Label>
            <Input
              id="cf-zone-id"
              placeholder="Optional zone ID"
              value={formData.cloudflareZoneId}
              onChange={(e) =>
                setFormData({
                  ...formData,
                  cloudflareZoneId: e.target.value,
                })
              }
              className="bg-muted/30 font-medium"
            />
          </div>
        </div>
      )}

      {activeProvider === 'namecheap' && (
        <div className="space-y-4">
          <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
            <div className="space-y-2">
              <Label
                htmlFor="nc-api-user"
                className="font-bold text-[10px] text-muted-foreground uppercase tracking-wider"
              >
                API USER
              </Label>
              <Input
                id="nc-api-user"
                placeholder="username"
                value={formData.namecheapApiUser}
                onChange={(e) =>
                  setFormData({
                    ...formData,
                    namecheapApiUser: e.target.value,
                  })
                }
                className="bg-muted/30 font-medium"
              />
            </div>
            <div className="space-y-2">
              <Label
                htmlFor="nc-api-key"
                className="font-bold text-[10px] text-muted-foreground uppercase tracking-wider"
              >
                API KEY
              </Label>
              <Input
                id="nc-api-key"
                type="password"
                placeholder="••••••••••••••••"
                value={formData.namecheapApiKey}
                onChange={(e) =>
                  setFormData({
                    ...formData,
                    namecheapApiKey: e.target.value,
                  })
                }
                className="bg-muted/30 font-medium"
              />
            </div>
          </div>
          <div className="space-y-2">
            <Label
              htmlFor="nc-client-ip"
              className="font-bold text-[10px] text-muted-foreground uppercase tracking-wider"
            >
              CLIENT IP
            </Label>
            <Input
              id="nc-client-ip"
              placeholder="Whitelisted server IP"
              value={formData.namecheapClientIp}
              onChange={(e) =>
                setFormData({
                  ...formData,
                  namecheapClientIp: e.target.value,
                })
              }
              className="bg-muted/30 font-medium"
            />
          </div>
        </div>
      )}

      {activeProvider === 'spaceship' && (
        <div className="space-y-4">
          <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
            <div className="space-y-2">
              <Label
                htmlFor="ss-api-key"
                className="font-bold text-[10px] text-muted-foreground uppercase tracking-wider"
              >
                API KEY
              </Label>
              <Input
                id="ss-api-key"
                type="password"
                placeholder="Spaceship API key"
                value={formData.spaceshipApiKey}
                onChange={(e) =>
                  setFormData({
                    ...formData,
                    spaceshipApiKey: e.target.value,
                  })
                }
                className="bg-muted/30 font-medium"
              />
            </div>
            <div className="space-y-2">
              <Label
                htmlFor="ss-api-secret"
                className="font-bold text-[10px] text-muted-foreground uppercase tracking-wider"
              >
                API SECRET
              </Label>
              <Input
                id="ss-api-secret"
                type="password"
                placeholder="Spaceship API secret"
                value={formData.spaceshipApiSecret}
                onChange={(e) =>
                  setFormData({
                    ...formData,
                    spaceshipApiSecret: e.target.value,
                  })
                }
                className="bg-muted/30 font-medium"
              />
            </div>
          </div>
        </div>
      )}

      <div className="flex items-center justify-end gap-2 pt-2">
        {onCancel && (
          <Button type="button" variant="ghost" size="sm" onClick={onCancel} disabled={isPending}>
            Cancel
          </Button>
        )}
        <Button
          onClick={() => handleSaveProvider(activeProvider)}
          disabled={isPending}
          size="sm"
          className="gap-1.5"
        >
          <CheckCircle2 className="size-4" />
          Save {providerName}
        </Button>
      </div>
    </div>
  );
}
