import { Badge } from '#/components/ui/badge';
import type { SystemStats } from '#/features/settings/types';

type Props = {
  stats?: SystemStats;
};

export const DockerStorageCard = ({ stats }: Props) => {
  return (
    <div className="border/50 overflow-hidden rounded-2xl border bg-card/40 lg:col-span-2">
      <div className="border/50 flex items-center justify-between border-b p-6">
        <h3 className="font-bold text-xl">Docker storage</h3>
        <Badge
          variant="outline"
          className="border-primary/50 bg-primary/10 px-3 py-1 font-bold text-[10px] text-primary uppercase tracking-widest"
        >
          AVAILABLE
        </Badge>
      </div>

      <div className="divide-y divide-border/50">
        <div className="grid grid-cols-3 items-center p-6">
          <div>
            <p className="font-bold text-[10px] text-muted-foreground uppercase tracking-[0.15em]">
              IMAGES
            </p>
            <p className="mt-1 text-muted-foreground text-xs">
              {stats?.docker?.images?.active || 0}/{stats?.docker?.images?.totalCount || 0} active
            </p>
          </div>
          <div className="col-span-2 flex items-center gap-6">
            <div className="flex h-2 w-48 overflow-hidden rounded-full bg-background">
              <div className="h-full w-1/3 bg-muted-foreground" />
            </div>
            <div className="space-x-2 font-mono text-sm">
              <span className="text-foreground">{stats?.docker?.images?.size || '0 B'}</span>
              <span className="text-yellow-500">
                {stats?.docker?.images?.reclaimable || '0 B'} candidate
              </span>
            </div>
          </div>
        </div>
        <div className="grid grid-cols-3 items-center p-6">
          <div>
            <p className="font-bold text-[10px] text-muted-foreground uppercase tracking-[0.15em]">
              CONTAINERS
            </p>
            <p className="mt-1 text-muted-foreground text-xs">
              {stats?.docker?.containers?.active || 0}/{stats?.docker?.containers?.totalCount || 0}{' '}
              active
            </p>
          </div>
          <div className="col-span-2 flex items-center gap-6">
            <div className="flex h-2 w-48 overflow-hidden rounded-full bg-background">
              <div className="h-full w-1/12 bg-muted-foreground" />
            </div>
            <div className="space-x-2 font-mono text-sm">
              <span className="text-foreground">{stats?.docker?.containers?.size || '0 B'}</span>
              <span className="text-muted-foreground/50">
                {stats?.docker?.containers?.reclaimable || '0 B'} candidate
              </span>
            </div>
          </div>
        </div>
        <div className="grid grid-cols-3 items-center p-6">
          <div>
            <p className="font-bold text-[10px] text-muted-foreground uppercase tracking-[0.15em]">
              LOCAL VOLUMES
            </p>
            <p className="mt-1 text-muted-foreground text-xs">
              {stats?.docker?.volumes?.active || 0}/{stats?.docker?.volumes?.totalCount || 0} active
            </p>
          </div>
          <div className="col-span-2 flex items-center gap-6">
            <div className="flex h-2 w-48 overflow-hidden rounded-full bg-background">
              <div className="h-full w-1/6 bg-muted-foreground" />
            </div>
            <div className="space-x-2 font-mono text-sm">
              <span className="text-foreground">{stats?.docker?.volumes?.size || '0 B'}</span>
              <span className="text-muted-foreground/50">
                {stats?.docker?.volumes?.reclaimable || '0 B'} candidate
              </span>
            </div>
          </div>
        </div>
        <div className="grid grid-cols-3 items-center p-6">
          <div>
            <p className="font-bold text-[10px] text-muted-foreground uppercase tracking-[0.15em]">
              BUILD CACHE
            </p>
            <p className="mt-1 text-muted-foreground text-xs">
              {stats?.docker?.buildCache?.active || 0}/{stats?.docker?.buildCache?.totalCount || 0}{' '}
              active
            </p>
          </div>
          <div className="col-span-2 flex items-center gap-6">
            <div className="flex h-2 w-48 overflow-hidden rounded-full bg-background">
              <div className="h-full w-[80%] bg-muted-foreground" />
            </div>
            <div className="space-x-2 font-mono text-sm">
              <span className="text-foreground">{stats?.docker?.buildCache?.size || '0 B'}</span>
              <span className="text-yellow-500">
                {stats?.docker?.buildCache?.reclaimable || '0 B'} candidate
              </span>
            </div>
          </div>
        </div>
        <div className="bg-background/30 p-6 text-muted-foreground text-xs">
          Docker can keep image layers listed as candidates after safe cleanup when running services
          still reference them.
        </div>
      </div>
    </div>
  );
};
