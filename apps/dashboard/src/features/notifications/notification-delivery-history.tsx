import { useQuery } from '@tanstack/react-query';
import { History } from 'lucide-react';
import { Badge } from '@codedock/ui/badge';
import { Empty, EmptyDescription, EmptyHeader, EmptyMedia, EmptyTitle } from '@codedock/ui/empty';
import { Skeleton } from '@codedock/ui/skeleton';
import { settingsService } from '#/features/settings/api';

export function NotificationDeliveryHistory() {
  const { data, isLoading } = useQuery({
    queryKey: ['notification-deliveries'],
    queryFn: () => settingsService.getDeliveries(20),
  });
  const rows = data?.data ?? [];

  return (
    <section className="rounded-2xl border border-border/60 bg-card p-6">
      <div className="mb-5 flex items-center gap-3">
        <div className="flex h-9 w-9 items-center justify-center rounded-xl bg-primary/10 text-primary">
          <History className="h-4 w-4" />
        </div>
        <div>
          <h2 className="font-semibold text-foreground text-sm">Delivery history</h2>
          <p className="text-muted-foreground text-xs">
            Recent notification attempts across every channel, newest first.
          </p>
        </div>
      </div>

      {isLoading ? (
        <div className="space-y-3">
          {Array.from({ length: 4 }).map((_, index) => (
            <Skeleton key={index} className="h-14 w-full rounded-xl" />
          ))}
        </div>
      ) : rows.length === 0 ? (
        <Empty>
          <EmptyHeader>
            <EmptyMedia variant="icon">
              <History />
            </EmptyMedia>
            <EmptyTitle>No deliveries yet</EmptyTitle>
            <EmptyDescription>
              Send a test notification from a channel above and it will appear here.
            </EmptyDescription>
          </EmptyHeader>
        </Empty>
      ) : (
        <ul className="divide-y divide-border/40">
          {rows.map((delivery) => (
            <li key={delivery.id} className="flex items-center gap-3 px-1 py-3">
              <div className="min-w-0 flex-1">
                <p className="truncate text-foreground text-sm">
                  <span className="font-medium">{delivery.channel}</span>
                  {delivery.title ? (
                    <span className="text-muted-foreground"> · {delivery.title}</span>
                  ) : null}
                </p>
                <p className="mt-0.5 truncate text-muted-foreground text-xs">
                  {delivery.eventType}
                  {delivery.error ? ` · ${delivery.error}` : null}
                </p>
              </div>
              <span className="shrink-0 text-muted-foreground text-xs">
                {new Date(delivery.createdAt).toLocaleString()}
              </span>
              <Badge variant={delivery.success ? 'secondary' : 'destructive'}>
                {delivery.success ? 'Delivered' : 'Failed'}
              </Badge>
            </li>
          ))}
        </ul>
      )}
    </section>
  );
}
