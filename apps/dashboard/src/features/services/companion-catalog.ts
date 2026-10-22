export interface CompanionCatalogItem {
  id: string;
  name: string;
  image: string;
  category: 'database' | 'cache' | 'queue' | 'search' | 'tools' | 'storage';
  description: string;
  port: number;
  defaultEnv: { key: string; value?: string; generate?: 'secret' }[];
  defaultVolume?: string;
}

/** Generates a cryptographically random secret string. */
export function generateCatalogSecret(bytes = 24): string {
  const arr = new Uint8Array(bytes);
  if (typeof crypto !== 'undefined' && crypto.getRandomValues) {
    crypto.getRandomValues(arr);
  } else {
    for (let i = 0; i < bytes; i++) arr[i] = Math.floor(Math.random() * 256);
  }
  return Array.from(arr, (b) => b.toString(16).padStart(2, '0')).join('');
}

/** Resolves initial environment variables, automatically populating generated secrets. */
export function resolveCatalogEnvironment(
  item: CompanionCatalogItem
): { key: string; value: string }[] {
  return item.defaultEnv.map((env) => ({
    key: env.key,
    value: env.generate === 'secret' ? generateCatalogSecret() : (env.value ?? ''),
  }));
}

export const COMPANION_CATALOG: CompanionCatalogItem[] = [
  {
    id: 'postgres',
    name: 'PostgreSQL',
    image: 'postgres:16-alpine',
    category: 'database',
    description: 'Open-source relational database.',
    port: 5432,
    defaultEnv: [
      { key: 'POSTGRES_USER', value: 'postgres' },
      { key: 'POSTGRES_PASSWORD', generate: 'secret' },
      { key: 'POSTGRES_DB', value: 'app' },
    ],
    defaultVolume: 'pgdata:/var/lib/postgresql/data',
  },
  {
    id: 'mysql',
    name: 'MySQL',
    image: 'mysql:8',
    category: 'database',
    description: 'Relational database common with PHP and Node apps.',
    port: 3306,
    defaultEnv: [
      { key: 'MYSQL_ROOT_PASSWORD', generate: 'secret' },
      { key: 'MYSQL_DATABASE', value: 'app' },
    ],
    defaultVolume: 'mysql_data:/var/lib/mysql',
  },
  {
    id: 'mariadb',
    name: 'MariaDB',
    image: 'mariadb:11',
    category: 'database',
    description: 'Fast, community-driven MySQL-compatible drop-in.',
    port: 3306,
    defaultEnv: [
      { key: 'MARIADB_ROOT_PASSWORD', generate: 'secret' },
      { key: 'MARIADB_DATABASE', value: 'app' },
    ],
    defaultVolume: 'mariadb_data:/var/lib/mysql',
  },
  {
    id: 'mongodb',
    name: 'MongoDB',
    image: 'mongo:7',
    category: 'database',
    description: 'Document-oriented NoSQL database.',
    port: 27017,
    defaultEnv: [
      { key: 'MONGO_INITDB_ROOT_USERNAME', value: 'root' },
      { key: 'MONGO_INITDB_ROOT_PASSWORD', generate: 'secret' },
    ],
    defaultVolume: 'mongo_data:/data/db',
  },
  {
    id: 'redis',
    name: 'Redis',
    image: 'redis:7-alpine',
    category: 'cache',
    description: 'In-memory key-value data structure store.',
    port: 6379,
    defaultEnv: [],
    defaultVolume: 'redis_data:/data',
  },
  {
    id: 'valkey',
    name: 'Valkey',
    image: 'valkey/valkey:8-alpine',
    category: 'cache',
    description: 'High-performance Redis-compatible in-memory store.',
    port: 6379,
    defaultEnv: [],
    defaultVolume: 'valkey_data:/data',
  },
  {
    id: 'memcached',
    name: 'Memcached',
    image: 'memcached:alpine',
    category: 'cache',
    description: 'Distributed memory object caching system.',
    port: 11211,
    defaultEnv: [],
  },
  {
    id: 'rabbitmq',
    name: 'RabbitMQ',
    image: 'rabbitmq:3-management-alpine',
    category: 'queue',
    description: 'Message broker with management plugin enabled.',
    port: 5672,
    defaultEnv: [
      { key: 'RABBITMQ_DEFAULT_USER', value: 'guest' },
      { key: 'RABBITMQ_DEFAULT_PASS', value: 'guest' },
    ],
    defaultVolume: 'rabbitmq_data:/var/lib/rabbitmq',
  },
  {
    id: 'meilisearch',
    name: 'Meilisearch',
    image: 'getmeili/meilisearch:v1.7',
    category: 'search',
    description: 'Fast, typo-tolerant search engine.',
    port: 7700,
    defaultEnv: [{ key: 'MEILI_NO_ANALYTICS', value: 'true' }],
    defaultVolume: 'meili_data:/meili_data',
  },
  {
    id: 'minio',
    name: 'MinIO',
    image: 'minio/minio:latest',
    category: 'storage',
    description: 'S3-compatible object storage for files and assets.',
    port: 9000,
    defaultEnv: [
      { key: 'MINIO_ROOT_USER', value: 'minioadmin' },
      { key: 'MINIO_ROOT_PASSWORD', generate: 'secret' },
    ],
    defaultVolume: 'minio_data:/data',
  },
  {
    id: 'mailpit',
    name: 'Mailpit',
    image: 'axllent/mailpit:latest',
    category: 'tools',
    description: 'Email and SMTP testing tool for development.',
    port: 1025,
    defaultEnv: [],
  },
];
