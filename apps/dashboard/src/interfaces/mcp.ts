export type McpTransportMode = 'http' | 'sse';

export interface McpToolAnnotations {
  readOnlyHint: boolean;
  destructiveHint: boolean;
}

export interface McpCatalogTool {
  name: string;
  description: string;
  method: string;
  path: string;
  category: string;
  auth: string;
  annotations: McpToolAnnotations;
}

export interface McpCatalogPromptArg {
  name: string;
  description: string;
  required?: boolean;
}

export interface McpCatalogPrompt {
  name: string;
  title: string;
  description: string;
  arguments?: McpCatalogPromptArg[];
}

export interface McpCatalogResource {
  uri: string;
  name: string;
  description: string;
  mimeType: string;
}

export interface McpCatalogPayload {
  serverInfo?: {
    name: string;
    version: string;
    protocol: string;
  };
  tools?: McpCatalogTool[];
  prompts?: McpCatalogPrompt[];
  resources?: McpCatalogResource[];
}

export interface McpCatalogEnvelope {
  data?: McpCatalogPayload;
}
