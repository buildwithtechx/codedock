import { Skeleton } from '#/components/ui/skeleton';

export function ServerListSkeleton({ rows = 4 }: { rows?: number }) {
  return (
    <div role="status" aria-label="Loading servers" className="divide-y divide-border/50">
      {Array.from({ length: rows }, (_, index) => (
        <div key={index} className="flex items-center gap-3.5 px-5 py-3">
          <Skeleton className="size-9 shrink-0 rounded-xl" />
          <div className="w-44 shrink-0 space-y-1.5 lg:w-56">
            <Skeleton className="h-4 w-3/4" />
            <Skeleton className="h-3 w-1/2" />
          </div>
          <div className="hidden flex-1 items-center gap-2 sm:flex">
            <Skeleton className="h-5 w-16 rounded-md" />
            <Skeleton className="h-5 w-12 rounded-md" />
          </div>
          <Skeleton className="ml-auto h-4 w-14 rounded-full" />
        </div>
      ))}
    </div>
  );
}
