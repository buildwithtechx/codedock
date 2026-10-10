export interface CapacityResources {
  cpuCores: number;
  memoryMb: number;
  diskMb: number;
}

export function CapacitySummary({
  resources,
  inline = false,
}: {
  resources: CapacityResources;
  inline?: boolean;
}) {
  const memory = (mb: number) =>
    mb >= 1024 ? `${(mb / 1024).toFixed(mb % 1024 === 0 ? 0 : 1)} GB` : `${mb} MB`;

  const items = [
    { label: 'vCPUs', value: `${resources.cpuCores}` },
    { label: 'RAM', value: memory(resources.memoryMb) },
    { label: 'Storage', value: memory(resources.diskMb) },
  ];

  return (
    <dl
      className={
        inline
          ? 'flex flex-wrap gap-x-4 gap-y-1'
          : 'grid grid-cols-3 gap-2 rounded-xl bg-muted/40 p-3'
      }
    >
      {items.map(({ label, value }) => (
        <div key={label} className={`flex min-w-0 ${inline ? 'items-baseline gap-1' : 'flex-col'}`}>
          <dt className={`order-last text-muted-foreground ${inline ? 'text-sm' : 'mt-1 text-xs'}`}>
            {label}
          </dt>
          <dd className="font-medium text-foreground text-sm tabular-nums">
            <bdi>{value}</bdi>
          </dd>
        </div>
      ))}
    </dl>
  );
}
