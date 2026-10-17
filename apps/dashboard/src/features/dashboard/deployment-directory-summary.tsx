import { Link } from '@tanstack/react-router';
import {
  Activity,
  ArrowRight,
  CheckCircle2,
  CircleX,
  Rocket,
  Zap,
} from 'lucide-react';
import type { OrganizationDeployment } from '#/features/services';

const activeStatuses = new Set(['PENDING', 'CLONING', 'PULLING', 'BUILDING']);
const successfulStatuses = new Set(['READY', 'ACTIVE', 'SUCCESS']);

export function DeploymentSummary({
  deployments,
  isLoading,
}: {
  deployments: OrganizationDeployment[];
  isLoading: boolean;
}) {
  const successful = deployments.filter((deployment) =>
    successfulStatuses.has(deployment.status.toUpperCase())
  ).length;
  const active = deployments.filter((deployment) =>
    activeStatuses.has(deployment.status.toUpperCase())
  ).length;
  const failed = deployments.filter(
    (deployment) => deployment.status.toUpperCase() === 'FAILED'
  ).length;
  const items = [
    { label: 'Total', value: deployments.length, icon: Rocket, tone: 'text-primary bg-primary/12' },
    {
      label: 'Successful',
      value: successful,
      icon: CheckCircle2,
      tone: 'text-emerald-500 bg-emerald-500/12',
    },
    { label: 'Failed', value: failed, icon: CircleX, tone: 'text-rose-500 bg-rose-500/12' },
  ];

  return (
    <div className="space-y-4">
      <aside className="rounded-2xl bg-card p-5">
        <div className="flex items-center gap-2">
          <Activity className="h-4 w-4 text-muted-foreground" />
          <h2 className="font-semibold text-sm">Overview</h2>
        </div>
        <div className="mt-5 space-y-3.5">
          {items.map((item) => (
            <div key={item.label} className="flex items-center justify-between">
              <span className="flex items-center gap-2.5 text-muted-foreground text-sm">
                <span
                  className={`flex h-8 w-8 items-center justify-center rounded-lg ${item.tone}`}
                >
                  <item.icon className="h-4 w-4" />
                </span>
                {item.label}
              </span>
              <span className="font-semibold">{isLoading ? '–' : item.value}</span>
            </div>
          ))}
          {active > 0 && (
            <div className="flex items-center justify-between">
              <span className="flex items-center gap-2.5 text-muted-foreground text-sm">
                <span className="flex h-8 w-8 items-center justify-center rounded-lg bg-amber-500/12 text-amber-500">
                  <Activity className="h-4 w-4" />
                </span>
                In progress
              </span>
              <span className="font-semibold">{isLoading ? '–' : active}</span>
            </div>
          )}
        </div>
        {!isLoading && deployments.length > 0 && (
          <div className="mt-4 flex h-1.5 overflow-hidden rounded-full bg-muted/40">
            {successful > 0 && (
              <div
                className="bg-emerald-500"
                style={{ width: `${(successful / deployments.length) * 100}%` }}
              />
            )}
            {failed > 0 && (
              <div
                className="bg-rose-500"
                style={{ width: `${(failed / deployments.length) * 100}%` }}
              />
            )}
            {active > 0 && (
              <div
                className="bg-amber-400"
                style={{ width: `${(active / deployments.length) * 100}%` }}
              />
            )}
          </div>
        )}
      </aside>
      <DeploymentTip hasDeployments={deployments.length > 0} />
    </div>
  );
}

function DeploymentTip({ hasDeployments }: { hasDeployments: boolean }) {
  return (
    <section className="rounded-2xl border border-primary/10 bg-linear-to-br from-primary/10 via-primary/5 to-transparent p-5">
      <div className="flex items-center gap-2">
        <Zap className="h-4 w-4 text-primary" />
        <h2 className="font-semibold text-sm">{hasDeployments ? 'Auto deploy' : 'Get started'}</h2>
      </div>
      <p className="mt-3 text-muted-foreground text-sm leading-6">
        {hasDeployments
          ? 'Push to a connected branch and Codedock builds and ships the new release automatically.'
          : 'Connect a repository or choose a template to create the first deployment.'}
      </p>
      <Link
        to={hasDeployments ? '/projects' : '/library'}
        className="mt-4 inline-flex items-center gap-1.5 font-medium text-sm transition-colors hover:text-primary"
      >
        {hasDeployments ? 'View projects' : 'Create project'}
        <ArrowRight className="h-3.5 w-3.5" />
      </Link>
    </section>
  );
}
