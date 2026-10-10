import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';

interface DeployProjectNameCardProps {
  projectName: string;
  onProjectNameChange: (name: string) => void;
}

export function DeployProjectNameCard({
  projectName,
  onProjectNameChange,
}: DeployProjectNameCardProps) {
  return (
    <section className="rounded-2xl border border-border/60 bg-card p-6">
      <Label htmlFor="deploy-project-name" className="font-semibold text-foreground text-sm">
        Project Name
      </Label>
      <Input
        id="deploy-project-name"
        value={projectName}
        onChange={(e) => onProjectNameChange(e.target.value)}
        placeholder="my-project"
        className="mt-3 h-10 rounded-xl border border-border/60 bg-background/50 font-mono text-sm"
      />
    </section>
  );
}
