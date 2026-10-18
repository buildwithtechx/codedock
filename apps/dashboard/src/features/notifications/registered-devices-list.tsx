import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { Apple, Smartphone, Trash2 } from 'lucide-react';
import { toast } from 'sonner';
import { Badge } from '@codedock/ui/badge';
import { Button } from '@codedock/ui/button';
import { Empty, EmptyDescription, EmptyHeader, EmptyMedia, EmptyTitle } from '@codedock/ui/empty';
import { Skeleton } from '@codedock/ui/skeleton';
import { settingsService } from '#/features/settings/api';

export function RegisteredDevicesList() {
  const queryClient = useQueryClient();

  const { data, isLoading } = useQuery({
    queryKey: ['notification-devices'],
    queryFn: () => settingsService.getDeviceTokens(),
  });

  const { mutateAsync: revoke, isPending: revoking } = useMutation({
    mutationFn: (token: string) => settingsService.deleteDeviceToken(token),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['notification-devices'] });
      toast.success('Device token revoked successfully');
    },
    onError: (err: unknown) => {
      const msg = err instanceof Error ? err.message : 'Failed to revoke device token';
      toast.error(msg);
    },
  });

  const devices = data?.data ?? [];

  return (
    <section className="rounded-2xl border border-border/60 bg-card p-6">
      <div className="mb-5 flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="flex h-9 w-9 items-center justify-center rounded-xl bg-primary/10 text-primary">
            <Smartphone className="h-4 w-4" />
          </div>
          <div>
            <h2 className="font-semibold text-foreground text-sm">Connected Mobile Devices</h2>
            <p className="text-muted-foreground text-xs">
              iOS and Android devices paired via the Codedock Expo mobile app receiving native push
              notifications.
            </p>
          </div>
        </div>
        <Badge variant="outline" className="text-xs">
          {devices.length} {devices.length === 1 ? 'device' : 'devices'}
        </Badge>
      </div>

      {isLoading ? (
        <div className="space-y-3">
          {Array.from({ length: 2 }).map((_, index) => (
            <Skeleton key={index} className="h-16 w-full rounded-xl" />
          ))}
        </div>
      ) : devices.length === 0 ? (
        <Empty>
          <EmptyHeader>
            <EmptyMedia variant="icon">
              <Smartphone />
            </EmptyMedia>
            <EmptyTitle>No mobile devices paired</EmptyTitle>
            <EmptyDescription>
              To receive instant push notifications, open the Codedock Expo mobile app, go to
              Settings, and connect to this server with your Personal Access Token.
            </EmptyDescription>
          </EmptyHeader>
        </Empty>
      ) : (
        <ul className="divide-y divide-border/40">
          {devices.map((device) => {
            const isApple = device.platform?.toLowerCase() === 'ios';
            return (
              <li key={device.id} className="flex items-center justify-between gap-3 px-1 py-3.5">
                <div className="flex items-center gap-3 min-w-0">
                  <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-muted text-muted-foreground">
                    {isApple ? <Apple className="h-4 w-4" /> : <Smartphone className="h-4 w-4" />}
                  </div>
                  <div className="min-w-0">
                    <div className="flex items-center gap-2">
                      <span className="font-medium text-foreground text-sm capitalize">
                        {device.platform || 'Expo'} Device
                      </span>
                      <Badge variant="secondary" className="text-[10px] uppercase">
                        {device.platform || 'expo'}
                      </Badge>
                    </div>
                    <p className="mt-0.5 truncate font-mono text-muted-foreground text-xs">
                      {device.token.replace(/(ExponentPushToken\[.{6}).+(.{4}\])/, '$1••••$2')}
                    </p>
                  </div>
                </div>

                <div className="flex items-center gap-3 shrink-0">
                  <span className="text-muted-foreground text-xs hidden sm:inline">
                    {device.createdAt ? new Date(device.createdAt).toLocaleDateString() : 'Active'}
                  </span>
                  <Button
                    variant="ghost"
                    size="sm"
                    className="text-muted-foreground hover:text-destructive"
                    onClick={() => revoke(device.token)}
                    disabled={revoking}
                    aria-label="Revoke device"
                  >
                    <Trash2 className="h-4 w-4" />
                  </Button>
                </div>
              </li>
            );
          })}
        </ul>
      )}
    </section>
  );
}
