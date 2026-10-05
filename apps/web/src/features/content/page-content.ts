export interface PageSection {
  title: string;
  body: string;
  points: string[];
}
export interface PageContent {
  eyebrow: string;
  title: string;
  description: string;
  highlights: string[];
  sections: PageSection[];
  steps: string[];
  docsPath: string;
  docsLabel: string;
}
