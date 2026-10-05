export interface TemplateRecipe {
  id: string;
  name: string;
  category: string;
  workflow: string;
  description: string;
  docsPath: string;
}
export const templateRecipes: TemplateRecipe[] = [
  {
    id: 'postgres',
    name: 'PostgreSQL',
    category: 'Databases',
    workflow: 'Database workflow',
    description:
      'Relational storage with database provisioning, SQL Studio, and backup configuration.',
    docsPath: '/deployments/templates/#postgres',
  },
  {
    id: 'mysql',
    name: 'MySQL',
    category: 'Databases',
    workflow: 'Database workflow',
    description:
      'Provision a relational database and connect application services using the generated credentials.',
    docsPath: '/deployments/templates/#mysql',
  },
  {
    id: 'redis',
    name: 'Redis',
    category: 'Databases',
    workflow: 'Database workflow',
    description:
      'Run a persistent cache or queue and configure authentication and storage for your use case.',
    docsPath: '/deployments/templates/#redis',
  },
  {
    id: 'mongodb',
    name: 'MongoDB',
    category: 'Databases',
    workflow: 'Database workflow',
    description:
      'Document storage with project connection details and engine-specific backup tooling.',
    docsPath: '/deployments/templates/#mongodb',
  },
  {
    id: 'supabase',
    name: 'Supabase',
    category: 'Application stacks',
    workflow: 'Compose recipe',
    description:
      'A multi-service stack: review Compose compatibility, credentials, routing, and persistence before deploying.',
    docsPath: '/deployments/templates/#supabase',
  },
  {
    id: 'minio',
    name: 'MinIO',
    category: 'Storage',
    workflow: 'Docker recipe',
    description:
      'S3-compatible object storage. Configure a persistent volume and add the endpoint as a backup destination.',
    docsPath: '/deployments/templates/#minio',
  },
  {
    id: 'n8n',
    name: 'n8n',
    category: 'Automation',
    workflow: 'Docker recipe',
    description:
      'Run workflow automation as an image service with persistent state and an application domain.',
    docsPath: '/deployments/templates/#n8n',
  },
  {
    id: 'grafana',
    name: 'Grafana',
    category: 'Monitoring',
    workflow: 'Docker recipe',
    description: 'Deploy a dashboard service and connect your own metrics or log data sources.',
    docsPath: '/deployments/templates/#grafana',
  },
  {
    id: 'ollama',
    name: 'Ollama',
    category: 'AI',
    workflow: 'Docker recipe',
    description:
      'Run a model-serving container on appropriately sized hardware with persistent model storage.',
    docsPath: '/deployments/templates/#ollama',
  },
  {
    id: 'qdrant',
    name: 'Qdrant',
    category: 'AI',
    workflow: 'Docker recipe',
    description:
      'Run a vector search service with a persistent data volume and private application connections.',
    docsPath: '/deployments/templates/#qdrant',
  },
];
