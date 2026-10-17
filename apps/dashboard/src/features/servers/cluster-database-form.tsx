import { Button } from '#/components/ui/button';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import type { ClusterDatabaseSpec, DatabaseAction } from '#/interfaces/databases';

export interface ClusterDatabaseFormProps {
  action: DatabaseAction;
  spec: ClusterDatabaseSpec;
  setSpec: (spec: ClusterDatabaseSpec) => void;
  versions: {
    operatorVersion: string;
    barmanVersion: string;
    certManagerVersion: string;
  };
  setVersions: (versions: {
    operatorVersion: string;
    barmanVersion: string;
    certManagerVersion: string;
  }) => void;
  environments: { id: string; name: string }[];
  destinations: { id: string; name: string }[];
  restoreTime: string;
  setRestoreTime: (time: string) => void;
  onPrepare: () => void;
  isPreparing: boolean;
}

export function ClusterDatabaseForm({
  action,
  spec,
  setSpec,
  versions,
  setVersions,
  environments,
  destinations,
  restoreTime,
  setRestoreTime,
  onPrepare,
  isPreparing,
}: ClusterDatabaseFormProps) {
  return (
    <div className="space-y-3">
      {action === 'operators'
        ? Object.entries(versions).map(([key, value]) => (
            <div key={key}>
              <Label htmlFor={key}>{key.replace(/Version$/, ' version')}</Label>
              <Input
                id={key}
                value={value}
                onChange={(event) => setVersions({ ...versions, [key]: event.target.value })}
              />
            </div>
          ))
        : (action === 'create' || action === 'restore') && (
            <>
              <Label htmlFor="cluster-db-name">Database name</Label>
              <Input
                id="cluster-db-name"
                value={spec.name}
                onChange={(event) => setSpec({ ...spec, name: event.target.value })}
              />
              <Label htmlFor="cluster-db-env">Environment</Label>
              <select
                id="cluster-db-env"
                className="w-full rounded border bg-background p-2"
                value={spec.environmentId}
                onChange={(event) => setSpec({ ...spec, environmentId: event.target.value })}
              >
                <option value="">Select an environment</option>
                {environments.map((environment) => (
                  <option key={environment.id} value={environment.id}>
                    {environment.name}
                  </option>
                ))}
              </select>
              <Label htmlFor="cluster-db-engine">Engine</Label>
              <select
                id="cluster-db-engine"
                className="w-full rounded border bg-background p-2"
                disabled={action === 'restore'}
                value={spec.engine}
                onChange={(event) => {
                  const engine = event.target.value as 'postgres' | 'redis';
                  setSpec({
                    ...spec,
                    engine,
                    image:
                      engine === 'redis' ? 'redis:7.4.2' : 'ghcr.io/cloudnative-pg/postgresql:17.6',
                    s3DestinationId: '',
                  });
                }}
              >
                <option value="postgres">PostgreSQL</option>
                <option value="redis">Redis</option>
              </select>
              <Label htmlFor="cluster-db-image">Exact engine image tag</Label>
              <Input
                id="cluster-db-image"
                value={spec.image}
                onChange={(event) => setSpec({ ...spec, image: event.target.value })}
              />
              <Label htmlFor="cluster-db-instances">Database mode</Label>
              <select
                id="cluster-db-instances"
                className="w-full rounded border bg-background p-2"
                value={spec.instances}
                onChange={(event) => setSpec({ ...spec, instances: Number(event.target.value) })}
              >
                <option value={1}>Standalone</option>
                <option value={3}>3 replicated instances</option>
              </select>
              <Label htmlFor="cluster-db-storage">Storage per instance (GiB)</Label>
              <Input
                id="cluster-db-storage"
                type="number"
                min={1}
                value={spec.storageGiB}
                onChange={(event) => setSpec({ ...spec, storageGiB: Number(event.target.value) })}
              />
              <Label htmlFor="cluster-db-class">Storage class (empty uses cluster default)</Label>
              <Input
                id="cluster-db-class"
                value={spec.storageClass}
                onChange={(event) => setSpec({ ...spec, storageClass: event.target.value })}
              />
              {spec.engine === 'postgres' && (
                <>
                  <Label htmlFor="cluster-db-s3">Object-store backups</Label>
                  <select
                    id="cluster-db-s3"
                    className="w-full rounded border bg-background p-2"
                    value={spec.s3DestinationId}
                    onChange={(event) => setSpec({ ...spec, s3DestinationId: event.target.value })}
                  >
                    <option value="">No object-store backups</option>
                    {destinations.map((destination) => (
                      <option key={destination.id} value={destination.id}>
                        {destination.name}
                      </option>
                    ))}
                  </select>
                  <p className="text-muted-foreground text-xs">
                    Daily base backups and WAL archiving use a separate database prefix with 30-day
                    retention. Generated credentials stay encrypted in Codedock and Kubernetes
                    Secrets.
                  </p>
                </>
              )}
              {action === 'restore' && (
                <>
                  <Label htmlFor="cluster-db-time">Recovery time (UTC ISO timestamp)</Label>
                  <Input
                    id="cluster-db-time"
                    value={restoreTime}
                    onChange={(event) => setRestoreTime(event.target.value)}
                  />
                  <p className="text-sm">
                    The source remains unchanged. The source must have a completed base backup and
                    continuous WAL coverage for this time.
                  </p>
                </>
              )}
            </>
          )}
      <Button disabled={isPreparing} onClick={onPrepare}>
        Prepare review
      </Button>
    </div>
  );
}
