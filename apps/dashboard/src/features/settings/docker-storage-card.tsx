import { Badge } from '#/components/ui/badge';
import type { SystemStats } from '#/features/settings';
import { cn } from '#/lib/utils';

type Props = {
  stats?: SystemStats;
};

export const DockerStorageCard = ({ stats }: Props) => {
  const items = [
    {
      label: 'IMAGES',
      active: stats?.docker?.images?.active || 0,
      total: stats?.docker?.images?.totalCount || 0,
      widthClass: 'w-1/3',
      size: stats?.docker?.images?.size || '0 B',
      reclaimable: stats?.docker?.images?.reclaimable || '0 B',
      candidateColor: 'text-yellow-500',
    },
    {
      label: 'CONTAINERS',
      active: stats?.docker?.containers?.active || 0,
      total: stats?.docker?.containers?.totalCount || 0,
      widthClass: 'w-1/12',
      size: stats?.docker?.containers?.size || '0 B',
      reclaimable: stats?.docker?.containers?.reclaimable || '0 B',
      candidateColor: 'text-muted-foreground/50',
    },
    {
      label: 'LOCAL VOLUMES',
      active: stats?.docker?.volumes?.active || 0,
      total: stats?.docker?.volumes?.totalCount || 0,
      widthClass: 'w-1/6',
      size: stats?.docker?.volumes?.size || '0 B',
      reclaimable: stats?.docker?.volumes?.reclaimable || '0 B',
      candidateColor: 'text-muted-foreground/50',
    },
    {
      label: 'BUILD CACHE',
      active: stats?.docker?.buildCache?.active || 0,
      total: stats?.docker?.buildCache?.totalCount || 0,
      widthClass: 'w-[80%]',
      size: stats?.docker?.buildCache?.size || '0 B',
      reclaimable: stats?.docker?.buildCache?.reclaimable || '0 B',
      candidateColor: 'text-yellow-500',
    },
  ];

  return (
    <div className="w-full overflow-hidden rounded-2xl border border-border/80 bg-card shadow-sm">
      <div className="flex items-center justify-between border-border/70 border-b p-6">
        <h3 className="font-bold text-xl">Docker storage</h3>
        <Badge
          variant="outline"
          className="border-primary/50 bg-primary/10 px-3 py-1 font-bold text-[10px] text-primary uppercase tracking-widest"
        >
          AVAILABLE
        </Badge>
      </div>

      <div className="divide-y divide-border/50">
        {items.map((item) => (
          <div
            key={item.label}
            className="grid grid-cols-1 gap-4 p-6 sm:grid-cols-3 sm:items-center"
          >
            <div>
              <p className="font-bold text-[10px] text-muted-foreground uppercase tracking-[0.15em]">
                {item.label}
              </p>
              <p className="mt-1 text-muted-foreground text-xs">
                {item.active}/{item.total} active
              </p>
            </div>
            <div className="col-span-2 flex items-center justify-between gap-6">
              <div className="flex max-w-xs flex-1 items-center gap-4">
                <div className="flex h-2 w-full overflow-hidden rounded-full bg-background">
                  <div className={cn('h-full bg-muted-foreground', item.widthClass)} />
                </div>
                <span className="shrink-0 font-mono text-foreground text-sm">{item.size}</span>
              </div>
              <div className="shrink-0 text-right font-mono text-sm">
                <span className={item.candidateColor}>{item.reclaimable} candidate</span>
              </div>
            </div>
          </div>
        ))}
        <div className="bg-background/30 p-6 text-muted-foreground text-xs">
          Docker can keep image layers listed as candidates after safe cleanup when running services
          still reference them.
        </div>
      </div>
    </div>
  );
};
