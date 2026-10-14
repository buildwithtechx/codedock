import { render, screen } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import { BackupStorageSummary } from './backup-storage-summary';
import type { BackupConfig, BackupRecord, S3Destination } from './interfaces';

const records: BackupRecord[] = [
  {
    id: 'r-1',
    backupConfigId: 'c-1',
    projectId: 'p-1',
    status: 'completed',
    filePath: '/snap/a.tgz',
    fileSizeBytes: 1024,
    logs: '',
    startedAt: new Date().toISOString(),
    completedAt: new Date().toISOString(),
  },
  {
    id: 'r-2',
    backupConfigId: 'c-1',
    projectId: 'p-1',
    status: 'failed',
    filePath: '',
    fileSizeBytes: 0,
    logs: 'boom',
    startedAt: new Date().toISOString(),
    completedAt: new Date().toISOString(),
  },
];

const destinations: S3Destination[] = [
  {
    id: 'd-1',
    projectId: 'p-1',
    name: 'offsite',
    description: '',
    provider: 's3',
    endpoint: 'https://s3.example.com',
    bucket: 'backups',
    region: 'us-east-1',
    accessKeyId: 'key',
    secretAccessKey: 'secret',
    createdAt: new Date().toISOString(),
  },
];

const configs: BackupConfig[] = [
  {
    id: 'c-1',
    projectId: 'p-1',
    name: 'nightly',
    description: '',
    dbUser: 'app',
    backupEnabled: true,
    s3Enabled: true,
    disableLocal: false,
    schedule: '0 2 * * *',
    timezone: 'UTC',
    timeout: 3600,
    retentionDays: 7,
    maxBackups: 0,
    maxStorageGb: 0,
    status: 'active',
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString(),
  },
];

describe('backup storage summary', () => {
  it('renders full breakdown with destination count', () => {
    render(
      <BackupStorageSummary
        records={records}
        destinations={destinations}
        configs={configs}
        sftpCount={1}
      />
    );
    expect(screen.getByText('Storage summary')).toBeDefined();
    expect(screen.getByText('Snapshots')).toBeDefined();
    expect(screen.getByText('Destinations')).toBeDefined();
    expect(screen.getByText('2')).toBeDefined();
  });

  it('renders compact strip with totals', () => {
    const { container } = render(
      <BackupStorageSummary
        records={records}
        destinations={destinations}
        configs={configs}
        compact
      />
    );
    expect(screen.queryByText('Storage summary')).toBeNull();
    expect(screen.getByText('snapshots')).toBeDefined();
    expect(screen.getByText('destinations')).toBeDefined();
    expect(container.querySelector('dl')).toBeNull();
  });
});
