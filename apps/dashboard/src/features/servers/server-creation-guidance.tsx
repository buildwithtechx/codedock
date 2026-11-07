import { Cloud, Info, KeyRound, Network, Server as ServerIcon } from 'lucide-react';
import type { ServerCreationMode } from './server-creation-mode-toggle';

const guidance = [
  {
    icon: Network,
    text: 'Any host with SSH access. Codedock connects over port 22 by default.',
  },
  {
    icon: KeyRound,
    text: 'An SSH private key or password for the connection user.',
  },
  {
    icon: ServerIcon,
    text: 'Saved servers are verified, then appear in your fleet.',
  },
];

const hetznerGuidance = [
  {
    icon: KeyRound,
    text: 'A Hetzner Cloud API token with read/write access. It is verified live before saving.',
  },
  {
    icon: Cloud,
    text: 'Pick a location and server type at live Hetzner list prices. Hetzner bills your token directly.',
  },
  {
    icon: ServerIcon,
    text: 'cloud-init installs Docker and hardens SSH automatically, then the node joins your fleet.',
  },
];

export function ServerCreationGuidance({ mode }: { mode: ServerCreationMode }) {
  const items = mode === 'hetzner' ? hetznerGuidance : guidance;
  return (
    <aside className="xl:sticky xl:top-6 xl:col-start-2 xl:row-start-1">
      <div className="rounded-2xl border border-border/50 bg-card">
        <div className="flex items-center gap-3 border-border/50 border-b px-5 py-4">
          <div className="flex h-9 w-9 items-center justify-center rounded-xl bg-warning/10">
            <Info className="size-[18px] text-warning" />
          </div>
          <div>
            <h2 className="font-semibold text-[15px] text-foreground">Getting started</h2>
            <p className="text-muted-foreground text-xs">What you need</p>
          </div>
        </div>
        <div className="p-5">
          <ul className="space-y-3 text-muted-foreground text-sm">
            {items.map((item) => (
              <li key={item.text} className="flex items-start gap-2">
                <item.icon className="mt-0.5 size-4 shrink-0 text-muted-foreground" />
                <span>{item.text}</span>
              </li>
            ))}
          </ul>
        </div>
      </div>
    </aside>
  );
}
