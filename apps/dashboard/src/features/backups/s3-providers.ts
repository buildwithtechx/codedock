import { Cloud, Database, HardDrive, Server } from 'lucide-react';

export type S3ProviderId = 'aws' | 'r2' | 'b2' | 'wasabi' | 'do' | 'minio';

export interface S3Provider {
  id: S3ProviderId;
  label: string;
  mode: 'region' | 'accountId' | 'endpoint';
  regionPlaceholder: string;
  icon: typeof Cloud;
}

export const s3Providers: S3Provider[] = [
  { id: 'aws', label: 'AWS S3', mode: 'region', regionPlaceholder: 'us-east-1', icon: Cloud },
  { id: 'r2', label: 'Cloudflare R2', mode: 'accountId', regionPlaceholder: 'auto', icon: Cloud },
  {
    id: 'b2',
    label: 'Backblaze B2',
    mode: 'region',
    regionPlaceholder: 'us-west-004',
    icon: HardDrive,
  },
  {
    id: 'wasabi',
    label: 'Wasabi',
    mode: 'region',
    regionPlaceholder: 'us-east-1',
    icon: Database,
  },
  { id: 'do', label: 'DO Spaces', mode: 'region', regionPlaceholder: 'nyc3', icon: Cloud },
  {
    id: 'minio',
    label: 'MinIO / Custom',
    mode: 'endpoint',
    regionPlaceholder: 'us-east-1',
    icon: Server,
  },
];

export function providerById(id: S3ProviderId): S3Provider {
  return s3Providers.find((provider) => provider.id === id) ?? s3Providers[0];
}

export function providerIdFromName(provider: string): S3ProviderId {
  const normalized = provider.toLowerCase();
  if (normalized === 'r2' || normalized.includes('cloudflare')) return 'r2';
  if (normalized === 'b2' || normalized.includes('backblaze')) return 'b2';
  if (normalized.includes('wasabi')) return 'wasabi';
  if (normalized === 'do' || normalized.includes('digitalocean') || normalized.includes('spaces'))
    return 'do';
  if (normalized.includes('minio') || normalized === 'other' || normalized === 'custom')
    return 'minio';
  return 'aws';
}

export function buildS3Endpoint(
  provider: S3ProviderId,
  fields: { region: string; accountId: string; endpoint: string }
): { endpoint: string; region: string } {
  const region = fields.region.trim();
  const accountId = fields.accountId.trim();
  const endpoint = fields.endpoint.trim();
  switch (provider) {
    case 'aws':
      return { endpoint: '', region: region || 'us-east-1' };
    case 'r2':
      return {
        endpoint: accountId ? `https://${accountId}.r2.cloudflarestorage.com` : '',
        region: 'auto',
      };
    case 'b2':
      return {
        endpoint: region ? `https://s3.${region}.backblazeb2.com` : '',
        region: region || '',
      };
    case 'wasabi':
      return {
        endpoint: region ? `https://s3.${region}.wasabisys.com` : '',
        region: region || '',
      };
    case 'do':
      return {
        endpoint: region ? `https://${region}.digitaloceanspaces.com` : '',
        region: region || '',
      };
    case 'minio':
      return { endpoint, region: region || 'us-east-1' };
  }
}

export function deriveS3Provider(
  endpoint: string,
  region: string
): { provider: S3ProviderId; region: string; accountId: string; endpoint: string } {
  const fallback = { provider: 'aws' as S3ProviderId, region, accountId: '', endpoint: '' };
  if (!endpoint || endpoint.includes('.amazonaws.com')) return fallback;
  const r2 = endpoint.match(/^https?:\/\/([a-z0-9]+)\.r2\.cloudflarestorage\.com/i);
  if (r2?.[1]) return { provider: 'r2', region: 'auto', accountId: r2[1], endpoint: '' };
  const b2 = endpoint.match(/^https?:\/\/s3\.([a-z0-9-]+)\.backblazeb2\.com/i);
  if (b2?.[1]) return { provider: 'b2', region: b2[1], accountId: '', endpoint: '' };
  if (endpoint.includes('.wasabisys.com')) {
    const match = endpoint.match(/s3\.([a-z0-9-]+)\.wasabisys\.com/i);
    return { provider: 'wasabi', region: match?.[1] ?? region, accountId: '', endpoint: '' };
  }
  if (endpoint.includes('.digitaloceanspaces.com')) {
    const match = endpoint.match(/^https?:\/\/([a-z0-9-]+)\.digitaloceanspaces\.com/i);
    return { provider: 'do', region: match?.[1] ?? region, accountId: '', endpoint: '' };
  }
  return { provider: 'minio', region, accountId: '', endpoint };
}
