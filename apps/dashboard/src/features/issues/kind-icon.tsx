import type { LucideIcon } from 'lucide-react';
import { ArrowLeftRight, Box, Database, Gauge, Rocket, TriangleAlert } from 'lucide-react';

const KIND_ICONS: Record<string, LucideIcon> = {
  'deployment-failed': Rocket,
  'service-unhealthy': Box,
  'backup-failed': Database,
  'migration-failed': ArrowLeftRight,
  'quota-pressure': Gauge,
};

export function KindIcon({ kind, className }: { kind: string; className?: string }) {
  const Icon = KIND_ICONS[kind] ?? TriangleAlert;
  return <Icon className={className} aria-hidden="true" />;
}
