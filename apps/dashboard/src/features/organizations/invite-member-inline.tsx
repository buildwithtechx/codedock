import { Check, Copy, Link, Lock, Send, Shield, Users } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import { Checkbox } from '#/components/ui/checkbox';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import { useListProjects } from '#/features/projects/hooks';
import { useInviteOrganizationMember } from './hooks';
import type { OrganizationRole } from './interfaces';

type RoleOption = {
  key: OrganizationRole;
  title: string;
  description: string;
  icon: typeof Users;
};

const ROLES: RoleOption[] = [
  {
    key: 'member',
    title: 'Member',
    description: 'Can deploy, view logs, and manage resources they create.',
    icon: Users,
  },
  {
    key: 'admin',
    title: 'Admin',
    description: 'Full access to manage workspace, resources, and team members.',
    icon: Shield,
  },
  {
    key: 'restricted',
    title: 'Restricted',
    description: 'Access only to specifically selected projects and resources.',
    icon: Lock,
  },
];

export function InviteMemberInline({
  organizationId,
  onInvited,
  onClose,
}: {
  organizationId: string;
  onInvited: () => void;
  onClose: () => void;
}) {
  const [email, setEmail] = useState('');
  const [role, setRole] = useState<OrganizationRole>('member');
  const [delivery, setDelivery] = useState<'email' | 'link'>('email');
  const [selectedProjects, setSelectedProjects] = useState<string[]>([]);
  const [copied, setCopied] = useState(false);
  const [created, setCreated] = useState<{ url: string; email: string; linkOnly: boolean } | null>(
    null
  );

  const { data: projectsData } = useListProjects();
  const projects = (projectsData?.data as any[]) || [];

  const inviteMutation = useInviteOrganizationMember(organizationId);

  const toggleProject = (projectId: string) => {
    setSelectedProjects((prev) =>
      prev.includes(projectId) ? prev.filter((id) => id !== projectId) : [...prev, projectId]
    );
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!email.trim()) return;

    try {
      const res = await inviteMutation.mutateAsync({
        email: email.trim(),
        permission: role,
        delivery,
        projects: role === 'restricted' ? selectedProjects : undefined,
      });

      const inviteUrl = `${window.location.origin}${res?.inviteLink || `/invite/${res?.id || ''}`}`;
      setCreated({
        url: inviteUrl,
        email: email.trim(),
        linkOnly: delivery === 'link',
      });
      onInvited();
      toast.success(delivery === 'link' ? 'Invitation link ready' : 'Invitation sent successfully');
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Failed to send invitation';
      toast.error(msg);
    }
  };

  const copyLink = async () => {
    if (!created?.url) return;
    try {
      await navigator.clipboard.writeText(created.url);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
      toast.success('Invitation link copied');
    } catch {
      toast.error('Failed to copy link');
    }
  };

  if (created) {
    return (
      <div className="space-y-4 rounded-xl border border-border/50 bg-card p-5">
        <div className="flex items-start gap-3.5">
          <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-emerald-500/10 text-emerald-500">
            {created.linkOnly ? <Link className="size-5" /> : <Send className="size-5" />}
          </div>
          <div className="min-w-0">
            <h3 className="font-semibold text-foreground text-sm">
              {created.linkOnly ? 'Invitation link ready' : 'Invitation sent'}
            </h3>
            <p className="mt-0.5 break-words text-muted-foreground text-xs">
              Recipient: <span className="font-medium text-foreground">{created.email}</span>
            </p>
          </div>
        </div>

        <div className="flex flex-col gap-2 sm:flex-row">
          <Input
            value={created.url}
            readOnly
            onFocus={(e) => e.currentTarget.select()}
            className="min-w-0 flex-1 font-mono text-xs"
          />
          <Button variant="outline" size="sm" onClick={copyLink} className="shrink-0 gap-1.5">
            {copied ? (
              <Check className="size-3.5 text-emerald-500" />
            ) : (
              <Copy className="size-3.5" />
            )}
            {copied ? 'Copied' : 'Copy link'}
          </Button>
        </div>

        <div className="flex justify-end border-border/40 border-t pt-3">
          <Button size="sm" onClick={onClose}>
            Done
          </Button>
        </div>
      </div>
    );
  }

  return (
    <form
      onSubmit={handleSubmit}
      className="space-y-5 rounded-xl border border-border/50 bg-card p-5"
    >
      <div className="grid gap-5 md:grid-cols-2">
        <div className="space-y-4">
          <div className="space-y-1.5">
            <Label className="text-xs">Email Address</Label>
            <Input
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="colleague@example.com"
              disabled={inviteMutation.isPending}
              required
              className="text-xs"
            />
          </div>

          <div className="space-y-2">
            <Label className="text-xs">Delivery Method</Label>
            <div className="grid grid-cols-2 gap-2">
              <button
                type="button"
                onClick={() => setDelivery('email')}
                className={`flex items-center justify-center gap-2 rounded-xl border p-2.5 font-medium text-xs transition-all ${
                  delivery === 'email'
                    ? 'border-primary bg-primary/10 text-primary'
                    : 'border-border/50 bg-muted/20 text-muted-foreground hover:bg-muted/40'
                }`}
              >
                <Send className="size-3.5" />
                Email
              </button>
              <button
                type="button"
                onClick={() => setDelivery('link')}
                className={`flex items-center justify-center gap-2 rounded-xl border p-2.5 font-medium text-xs transition-all ${
                  delivery === 'link'
                    ? 'border-primary bg-primary/10 text-primary'
                    : 'border-border/50 bg-muted/20 text-muted-foreground hover:bg-muted/40'
                }`}
              >
                <Link className="size-3.5" />
                Invite Link
              </button>
            </div>
            <p className="text-[11px] text-muted-foreground">
              {delivery === 'email'
                ? 'An invitation email will be delivered to the recipient with an onboarding link.'
                : 'Generate a shareable link that you can send directly to your teammate.'}
            </p>
          </div>
        </div>

        <div className="space-y-2">
          <Label className="text-xs">Workspace Role</Label>
          <div className="space-y-2">
            {ROLES.map((item) => {
              const Icon = item.icon;
              const isSelected = role === item.key;
              return (
                <button
                  key={item.key}
                  type="button"
                  onClick={() => setRole(item.key)}
                  className={`flex w-full items-start gap-3 rounded-xl border p-3 text-left transition-all ${
                    isSelected
                      ? 'border-primary bg-primary/5 text-foreground ring-1 ring-primary'
                      : 'border-border/50 bg-muted/10 text-muted-foreground hover:bg-muted/20'
                  }`}
                >
                  <div
                    className={`mt-0.5 flex size-7 shrink-0 items-center justify-center rounded-lg border ${
                      isSelected
                        ? 'border-primary/40 bg-primary/10 text-primary'
                        : 'border-border/60 bg-muted/40 text-muted-foreground'
                    }`}
                  >
                    <Icon className="size-3.5" />
                  </div>
                  <div className="min-w-0 flex-1">
                    <p className={`font-semibold text-xs ${isSelected ? 'text-foreground' : ''}`}>
                      {item.title}
                    </p>
                    <p className="mt-0.5 text-[11px] text-muted-foreground leading-relaxed">
                      {item.description}
                    </p>
                  </div>
                </button>
              );
            })}
          </div>
        </div>
      </div>

      {role === 'restricted' && (
        <div className="space-y-3 border-border/40 border-t pt-4">
          <div>
            <h4 className="font-medium text-foreground text-xs">Granted Projects</h4>
            <p className="mt-0.5 text-[11px] text-muted-foreground">
              Select which projects this restricted teammate can view and deploy.
            </p>
          </div>
          {projects.length === 0 ? (
            <p className="text-muted-foreground text-xs italic">
              No projects found in this workspace.
            </p>
          ) : (
            <div className="grid max-h-48 grid-cols-1 gap-2 overflow-y-auto rounded-xl border border-border/50 bg-muted/10 p-3 sm:grid-cols-2">
              {projects.map((proj: any) => {
                const checked = selectedProjects.includes(proj.id);
                const inputId = `proj-${proj.id}`;
                return (
                  <div
                    key={proj.id}
                    className="flex items-center gap-2.5 rounded-lg border border-border/40 bg-background/50 p-2 text-xs transition-colors hover:bg-muted/30"
                  >
                    <Checkbox
                      id={inputId}
                      checked={checked}
                      onCheckedChange={() => toggleProject(proj.id)}
                    />
                    <Label
                      htmlFor={inputId}
                      className="cursor-pointer truncate font-medium text-foreground text-xs"
                    >
                      {proj.name}
                    </Label>
                  </div>
                );
              })}
            </div>
          )}
        </div>
      )}

      <div className="flex items-center justify-end gap-2 border-border/40 border-t pt-3">
        <Button
          type="button"
          variant="ghost"
          size="sm"
          onClick={onClose}
          disabled={inviteMutation.isPending}
        >
          Cancel
        </Button>
        <Button
          type="submit"
          size="sm"
          disabled={inviteMutation.isPending || !email.trim()}
          className="gap-1.5"
        >
          {delivery === 'link' ? <Link className="size-3.5" /> : <Send className="size-3.5" />}
          {inviteMutation.isPending
            ? 'Inviting...'
            : delivery === 'link'
              ? 'Create Invite Link'
              : 'Send Invite'}
        </Button>
      </div>
    </form>
  );
}
