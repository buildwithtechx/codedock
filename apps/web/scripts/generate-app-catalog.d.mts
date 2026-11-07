export function parseScalar(line: string): [string, string | boolean] | null;
export function parseCatalogEntry(id: string, text: string): {
  id: string;
  name: string;
  description: string;
  category: string;
  icon: string;
  verified: boolean;
  services: string[];
  secrets: Array<{ var: string; label: string; length: number }>;
} | null;
