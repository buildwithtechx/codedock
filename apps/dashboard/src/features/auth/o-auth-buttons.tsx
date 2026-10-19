import { Button } from '@codedock/ui/button';
import { useEnabledOAuthProviders, useTriggerOAuthLogin } from '#/hooks/use-oauth';

export const OAuthButtons = () => {
  const { data: enabledProvidersData } = useEnabledOAuthProviders();
  const triggerOAuthLogin = useTriggerOAuthLogin();
  const enabledProviders = (enabledProvidersData?.data || []).map((p) =>
    p.providerName.toLowerCase()
  );

  const handleOAuthLogin = (provider: string) => {
    triggerOAuthLogin(provider);
  };

  if (!enabledProviders.length) return null;

  return (
    <>
      <div className="space-y-2">
        {enabledProviders.map((provider) => (
          <Button
            key={provider}
            variant="outline"
            type="button"
            onClick={() => handleOAuthLogin(provider)}
            className="w-full capitalize"
          >
            Continue with {provider}
          </Button>
        ))}
      </div>

      <div className="relative">
        <div className="absolute inset-0 flex items-center">
          <div className="h-px w-full bg-border" />
        </div>
        <div className="relative flex justify-center">
          <span className="bg-background px-3 text-[10px] text-muted-foreground uppercase tracking-[0.16em]">
            or
          </span>
        </div>
      </div>
    </>
  );
};
