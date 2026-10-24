import { Laptop, Server, Terminal } from 'lucide-react';

export type RunnerPresetKey = 'ubuntu' | 'macos' | 'linux' | 'custom';

export interface RunnerPresetMetadata {
  id: RunnerPresetKey;
  title: string;
  badge: string;
  description: string;
  icon: typeof Server;
}

export const RUNNER_PRESETS: RunnerPresetMetadata[] = [
  {
    id: 'ubuntu',
    title: 'Ubuntu 22.04',
    badge: 'Ready-to-use Linux image with Node.js and Git.',
    description: 'Choose a server to see compatible environments.',
    icon: Server,
  },
  {
    id: 'macos',
    title: 'macOS',
    badge: 'Use the tools installed on your connected Mac.',
    description: 'Choose a server to see compatible environments.',
    icon: Laptop,
  },
  {
    id: 'linux',
    title: 'Linux',
    badge: "Run directly with your server's installed tools.",
    description: 'Choose a server to see compatible environments.',
    icon: Terminal,
  },
  {
    id: 'custom',
    title: 'Custom image',
    badge: 'Bring a Docker runner image with Node.js and Git.',
    description: 'Choose a server to see compatible environments.',
    icon: Server,
  },
];
