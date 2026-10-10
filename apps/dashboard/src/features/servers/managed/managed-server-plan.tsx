import { Link } from '@tanstack/react-router';
import { Button } from '#/components/ui/button';
import type { CloudWorkspaceSummary } from '#/interfaces/server';
import { CapacitySummary } from './capacity-summary';

export function ManagedServerPlan({ server }: { server?: CloudWorkspaceSummary | null }) {
  const resources = server?.resources ?? {
    cpuCores: 4,
    memoryMb: 8192,
    diskMb: 80000,
  };

  const isFree = server?.planTierId === 'free';

  return (
    <section className="space-y-4 rounded-2xl bg-card p-5">
      <h2 className="font-medium text-base">Provisioned capacity</h2>
      <CapacitySummary resources={resources} />
      <p className="text-muted-foreground text-sm">
        Resources dynamically allocated across workloads.
      </p>
      <Button asChild className="w-full">
        <Link to="/billing">{isFree ? 'Choose a plan' : 'Change plan'}</Link>
      </Button>
      {!isFree && (
        <Button variant="secondary" className="w-full" asChild>
          <Link to="/billing">Resize capacity</Link>
        </Button>
      )}
      <Link
        to="/billing"
        className="block text-center font-medium text-muted-foreground text-sm hover:text-foreground"
      >
        Billing & usage
      </Link>
    </section>
  );
}
