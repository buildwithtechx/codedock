import { Link } from '@tanstack/react-router';
import { ArrowRight, Zap } from 'lucide-react';
import { useEffect, useMemo, useState } from 'react';
import { useGetPublicSettings } from '#/hooks/use-settings';

interface ProductTip {
  id: string;
  text: string;
  label: string;
  href: string;
  selfHostedOnly?: boolean;
}

const PRODUCT_TIPS: ProductTip[] = [
  {
    id: 'envVars',
    text: 'Keep environment variables encrypted and shared across services in your project settings.',
    label: 'View projects',
    href: '/projects',
  },
  {
    id: 'customDomain',
    text: 'Attach a custom domain with automatic Let’s Encrypt SSL to any web service.',
    label: 'Manage domains',
    href: '/domains',
  },
  {
    id: 'rollback',
    text: 'Instantly roll back to any previous deployment snapshot with zero downtime.',
    label: 'View deployments',
    href: '/deployments',
  },
  {
    id: 'apps',
    text: 'Deploy pre-configured stacks like Supabase, PocketBase, and Grafana in one click.',
    label: 'Browse apps',
    href: '/apps/new',
  },
  {
    id: 'servers',
    text: 'Connect additional Linux servers over SSH to distribute workloads across your fleet.',
    label: 'Manage servers',
    href: '/servers',
    selfHostedOnly: true,
  },
  {
    id: 'jobs',
    text: 'Schedule recurring cron jobs to run maintenance tasks or scripts inside your containers.',
    label: 'View cron jobs',
    href: '/jobs',
  },
  {
    id: 'backups',
    text: 'Configure automated S3 or SFTP backups to protect your databases and persistent volumes.',
    label: 'Configure backups',
    href: '/backups',
  },
  {
    id: 'mail',
    text: 'Run your own self-hosted mail server with automated DKIM, SPF, and DMARC checks.',
    label: 'Open emails',
    href: '/emails',
    selfHostedOnly: true,
  },
  {
    id: 'team',
    text: 'Invite collaborators and review workspace activity in the Team Audit Log.',
    label: 'Manage team',
    href: '/team',
  },
];

export function HomeNextStep({
  hasProjects,
  isLoading = false,
}: {
  hasProjects: boolean;
  isLoading?: boolean;
}) {
  const { data: publicSettings } = useGetPublicSettings();
  const isCloud = Boolean(publicSettings?.data?.cloudMode);
  const selfHosted = !isCloud;

  const pool = useMemo(
    () => PRODUCT_TIPS.filter((tip) => selfHosted || !tip.selfHostedOnly),
    [selfHosted]
  );

  const [idx, setIdx] = useState(0);
  useEffect(() => {
    if (pool.length > 1) {
      setIdx(Math.floor(Math.random() * pool.length));
    }
  }, [pool.length]);

  const tip = !isLoading && !hasProjects
    ? {
        text: 'Create your first project from a Git repository, Dockerfile, or ready-to-use template.',
        label: 'Create a project',
        href: '/library',
      }
    : (pool[idx] ?? pool[0]);

  if (!tip) return null;

  return (
    <div className="rounded-2xl border border-primary/10 bg-linear-to-br from-primary/5 via-primary/3 to-transparent p-5">
      <div className="mb-3 flex items-center gap-2">
        <Zap className="size-4 text-primary" />
        <h3 className="font-semibold text-foreground text-sm">Quick tip</h3>
      </div>
      <p className="text-muted-foreground text-sm leading-relaxed">{tip.text}</p>
      <Link
        to={tip.href}
        className="mt-3 inline-flex items-center gap-1 font-medium text-primary text-sm transition-colors hover:text-primary/80"
      >
        {tip.label}
        <ArrowRight className="size-3.5 rtl:rotate-180" />
      </Link>
    </div>
  );
}
