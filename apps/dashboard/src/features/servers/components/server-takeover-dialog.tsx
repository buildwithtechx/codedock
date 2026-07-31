import { useMutation, useQueryClient } from '@tanstack/react-query';
import { Check, Download, RefreshCw, Server } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Badge } from '#/components/ui/badge';
import { Button } from '#/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '#/components/ui/dialog';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import { Textarea } from '#/components/ui/textarea';
import { apiClient } from '#/lib/api-client';

interface DiscoveredContainer {
  id: string;
  name: string;
  image: string;
  detectedPlatform: string;
  composeProject: string;
  status: string;
}

export function ServerTakeoverDialog() {
  const queryClient = useQueryClient();
  const [open, setOpen] = useState(false);
  const [step, setStep] = useState<'credentials' | 'discovered'>('credentials');
  const [host, setHost] = useState('');
  const [port, setPort] = useState('22');
  const [user, setUser] = useState('root');
  const [key, setKey] = useState('');
  const [password, setPassword] = useState('');
  const [runId, setRunId] = useState('');
  const [containers, setContainers] = useState<DiscoveredContainer[]>([]);
  const [selectedIds, setSelectedIds] = useState<string[]>([]);
  const [targetProject, setTargetProject] = useState('takeover-project');

  const scanMutation = useMutation({
    mutationFn: async () => {
      const parsedPort = Number.parseInt(port, 10);
      if (Number.isNaN(parsedPort) || parsedPort < 1 || parsedPort > 65535) {
        throw new Error('Invalid SSH port. Must be between 1 and 65535');
      }
      const res = await apiClient.post<{
        data: { runId: string; containers: DiscoveredContainer[] };
      }>('/system/takeover/scan', {
        host,
        port: parsedPort,
        user,
        key,
        password,
      });
      return res.data;
    },
    onSuccess: (data) => {
      setRunId(data.runId);
      setContainers(data.containers || []);
      setSelectedIds((data.containers || []).map((c) => c.id));
      setStep('discovered');
      toast.success(`Discovered ${data.containers?.length || 0} containers`);
    },
    onError: (err: Error) => {
      toast.error(err.message || 'SSH scan failed');
    },
  });

  const adoptMutation = useMutation({
    mutationFn: async () => {
      return apiClient.post('/system/takeover/adopt', {
        runId,
        projectID: targetProject,
        containerIDs: selectedIds,
      });
    },
    onSuccess: () => {
      toast.success('Successfully adopted server stacks!');
      setOpen(false);
      setKey('');
      setPassword('');
      setStep('credentials');
      queryClient.invalidateQueries({ queryKey: ['projects'] });
    },
    onError: (err: Error) => {
      toast.error(err.message || 'Failed to adopt stacks');
    },
  });

  const toggleSelect = (id: string) => {
    setSelectedIds((prev) =>
      prev.includes(id) ? prev.filter((item) => item !== id) : [...prev, id]
    );
  };

  return (
    <Dialog
      open={open}
      onOpenChange={(isOpen) => {
        setOpen(isOpen);
        if (!isOpen) {
          setKey('');
          setPassword('');
        }
      }}
    >
      <DialogTrigger asChild>
        <Button variant="outline" className="gap-2 border-border/50">
          <Download className="h-4 w-4 text-blue-400" />
          Takeover Server Wizard
        </Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-150">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <Server className="h-5 w-5 text-primary" />
            Server Takeover & Migration Wizard
          </DialogTitle>
          <DialogDescription>
            Import running containers from Dokploy, Coolify, Dokku, or plain Docker via SSH.
          </DialogDescription>
        </DialogHeader>

        {step === 'credentials' ? (
          <div className="grid gap-4 py-4">
            <div className="grid grid-cols-3 gap-4">
              <div className="col-span-2 grid gap-2">
                <Label>Host / IP Address</Label>
                <Input
                  value={host}
                  onChange={(e) => setHost(e.target.value)}
                  placeholder="192.168.1.100"
                />
              </div>
              <div className="grid gap-2">
                <Label>SSH Port</Label>
                <Input value={port} onChange={(e) => setPort(e.target.value)} />
              </div>
            </div>
            <div className="grid gap-2">
              <Label>SSH User</Label>
              <Input value={user} onChange={(e) => setUser(e.target.value)} />
            </div>
            <div className="grid gap-2">
              <Label>SSH Private Key (Optional)</Label>
              <Textarea
                value={key}
                onChange={(e) => setKey(e.target.value)}
                placeholder="-----BEGIN OPENSSH PRIVATE KEY-----"
                rows={5}
              />
            </div>
            <div className="grid gap-2">
              <Label>SSH Password (Optional)</Label>
              <Input
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
              />
            </div>
            <Button
              onClick={() => scanMutation.mutate()}
              disabled={!host || scanMutation.isPending}
              className="gap-2"
            >
              {scanMutation.isPending ? (
                <RefreshCw className="h-4 w-4 animate-spin" />
              ) : (
                <Download className="h-4 w-4" />
              )}
              {scanMutation.isPending ? 'Scanning Server...' : 'Scan Server Stacks'}
            </Button>
          </div>
        ) : (
          <div className="grid gap-4 py-4">
            <div className="grid gap-2">
              <Label>Target Codedock Project Name</Label>
              <Input value={targetProject} onChange={(e) => setTargetProject(e.target.value)} />
            </div>
            <div className="max-h-75 space-y-2 overflow-y-auto">
              <Label>Discovered Containers ({containers.length})</Label>
              {containers.map((c) => (
                <button
                  type="button"
                  key={c.id}
                  onClick={() => toggleSelect(c.id)}
                  className={`flex w-full cursor-pointer items-center justify-between rounded-lg border p-3 text-left transition-colors ${
                    selectedIds.includes(c.id) ? 'border-primary bg-primary/5' : 'border-border/50'
                  }`}
                >
                  <div>
                    <div className="flex items-center gap-2">
                      <span className="font-medium text-sm">{c.name}</span>
                      <Badge variant="outline" className="text-[10px] capitalize">
                        {c.detectedPlatform || 'Docker'}
                      </Badge>
                    </div>
                    <p className="font-mono text-muted-foreground text-xs">{c.image}</p>
                  </div>
                  {selectedIds.includes(c.id) && <Check className="h-4 w-4 text-primary" />}
                </button>
              ))}
            </div>
            <div className="flex gap-2">
              <Button variant="outline" onClick={() => setStep('credentials')}>
                Back
              </Button>
              <Button
                onClick={() => adoptMutation.mutate()}
                disabled={selectedIds.length === 0 || adoptMutation.isPending}
                className="flex-1"
              >
                {adoptMutation.isPending
                  ? 'Adopting Stacks...'
                  : `Adopt ${selectedIds.length} Stacks into Codedock`}
              </Button>
            </div>
          </div>
        )}
      </DialogContent>
    </Dialog>
  );
}
