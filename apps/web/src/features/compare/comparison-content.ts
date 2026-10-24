import type { ExternalLink } from '../../components/site-link';
export interface ComparisonContent {
  name: string;
  tagline: string;
  description: string;
  alternativeFit: string;
  codedockFit: string;
  rows: { criterion: string; codedock: string; alternative: string }[];
  sources: { title: string; href: ExternalLink }[];
}
