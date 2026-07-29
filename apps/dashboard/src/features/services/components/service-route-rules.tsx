import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { Plus, Shield, Trash2 } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Badge } from '#/components/ui/badge';
import { Button } from '#/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '#/components/ui/card';
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
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '#/components/ui/select';
import { apiClient } from '#/lib/api-client';

export interface RouteRule {
  id: string;
  serviceId: string;
  name: string;
  ruleType: 'rate_limit' | 'ip_allowlist' | 'headers';
  spec: Record<string, unknown>;
  enabled: boolean;
  createdAt: string;
}

export function ServiceRouteRules({ serviceId }: { serviceId: string }) {
  const queryClient = useQueryClient();
  const [open, setOpen] = useState(false);
  const [name, setName] = useState('');
  const [ruleType, setRuleType] = useState<'rate_limit' | 'ip_allowlist' | 'headers'>('rate_limit');
  const [rateLimitAvg, setRateLimitAvg] = useState('100');
  const [rateLimitBurst, setRateLimitBurst] = useState('50');
  const [ipList, setIpList] = useState('192.168.1.0/24');
  const [headerKey, setHeaderKey] = useState('X-Custom-Header');
  const [headerVal, setHeaderVal] = useState('value');

  const { data: rules = [], isLoading } = useQuery({
    queryKey: ['route-rules', serviceId],
    queryFn: async () => {
      const res = await apiClient.get<{ data: RouteRule[] }>(`/services/${serviceId}/route-rules`);
      return res.data || [];
    },
  });

  const createMutation = useMutation({
    mutationFn: async () => {
      let spec: Record<string, unknown> = {};
      if (ruleType === 'rate_limit') {
        spec = {
          average: Number.parseInt(rateLimitAvg, 10),
          burst: Number.parseInt(rateLimitBurst, 10),
        };
      } else if (ruleType === 'ip_allowlist') {
        spec = { ips: ipList.split(',').map((s) => s.trim()) };
      } else if (ruleType === 'headers') {
        spec = { customRequestHeaders: { [headerKey]: headerVal } };
      }

      return apiClient.post(`/services/${serviceId}/route-rules`, {
        name,
        ruleType,
        spec,
        enabled: true,
      });
    },
    onSuccess: () => {
      toast.success('Route rule created');
      setOpen(false);
      setName('');
      queryClient.invalidateQueries({ queryKey: ['route-rules', serviceId] });
    },
    onError: (err: Error) => {
      toast.error(err.message || 'Failed to create route rule');
    },
  });

  const deleteMutation = useMutation({
    mutationFn: async (ruleId: string) => {
      return apiClient.delete(`/services/${serviceId}/route-rules/${ruleId}`);
    },
    onSuccess: () => {
      toast.success('Route rule deleted');
      queryClient.invalidateQueries({ queryKey: ['route-rules', serviceId] });
    },
  });

  return (
    <Card className="border-border/50">
      <CardHeader className="flex flex-row items-center justify-between">
        <div>
          <CardTitle className="flex items-center gap-2">
            <Shield className="h-5 w-5 text-primary" />
            Route Rules & Middleware
          </CardTitle>
          <CardDescription>
            Configure rate limiting, IP allowlists, and custom headers for this service.
          </CardDescription>
        </div>
        <Dialog open={open} onOpenChange={setOpen}>
          <DialogTrigger asChild>
            <Button size="sm" className="gap-2">
              <Plus className="h-4 w-4" />
              Add Rule
            </Button>
          </DialogTrigger>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Create Route Rule</DialogTitle>
              <DialogDescription>
                Add Traefik reverse proxy middleware rule for this service.
              </DialogDescription>
            </DialogHeader>
            <div className="grid gap-4 py-4">
              <div className="grid gap-2">
                <Label>Rule Name</Label>
                <Input
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  placeholder="e.g. Rate Limit API"
                />
              </div>
              <div className="grid gap-2">
                <Label>Rule Type</Label>
                <Select value={ruleType} onValueChange={(v: any) => setRuleType(v)}>
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="rate_limit">Rate Limit</SelectItem>
                    <SelectItem value="ip_allowlist">IP Allowlist</SelectItem>
                    <SelectItem value="headers">Custom Headers</SelectItem>
                  </SelectContent>
                </Select>
              </div>

              {ruleType === 'rate_limit' && (
                <div className="grid grid-cols-2 gap-4">
                  <div className="grid gap-2">
                    <Label>Average Req/sec</Label>
                    <Input value={rateLimitAvg} onChange={(e) => setRateLimitAvg(e.target.value)} />
                  </div>
                  <div className="grid gap-2">
                    <Label>Burst</Label>
                    <Input
                      value={rateLimitBurst}
                      onChange={(e) => setRateLimitBurst(e.target.value)}
                    />
                  </div>
                </div>
              )}

              {ruleType === 'ip_allowlist' && (
                <div className="grid gap-2">
                  <Label>IP CIDR (comma separated)</Label>
                  <Input
                    value={ipList}
                    onChange={(e) => setIpList(e.target.value)}
                    placeholder="192.168.1.0/24, 10.0.0.1/32"
                  />
                </div>
              )}

              {ruleType === 'headers' && (
                <div className="grid grid-cols-2 gap-4">
                  <div className="grid gap-2">
                    <Label>Header Key</Label>
                    <Input value={headerKey} onChange={(e) => setHeaderKey(e.target.value)} />
                  </div>
                  <div className="grid gap-2">
                    <Label>Header Value</Label>
                    <Input value={headerVal} onChange={(e) => setHeaderVal(e.target.value)} />
                  </div>
                </div>
              )}

              <Button
                onClick={() => createMutation.mutate()}
                disabled={!name || createMutation.isPending}
              >
                {createMutation.isPending ? 'Saving...' : 'Save Rule'}
              </Button>
            </div>
          </DialogContent>
        </Dialog>
      </CardHeader>
      <CardContent>
        {isLoading ? (
          <p className="text-muted-foreground text-sm">Loading rules...</p>
        ) : rules.length === 0 ? (
          <p className="text-muted-foreground text-sm">
            No custom route rules defined. Traffic will pass standard proxy.
          </p>
        ) : (
          <div className="grid gap-3">
            {rules.map((r) => (
              <div
                key={r.id}
                className="flex items-center justify-between rounded-lg border border-border/50 p-3"
              >
                <div className="grid gap-1">
                  <div className="flex items-center gap-2">
                    <span className="font-medium text-sm">{r.name}</span>
                    <Badge variant="outline" className="text-xs capitalize">
                      {r.ruleType.replace('_', ' ')}
                    </Badge>
                  </div>
                  <pre className="font-mono text-muted-foreground text-xs">
                    {JSON.stringify(r.spec)}
                  </pre>
                </div>
                <Button
                  variant="ghost"
                  size="icon"
                  className="text-destructive hover:text-destructive"
                  onClick={() => deleteMutation.mutate(r.id)}
                >
                  <Trash2 className="h-4 w-4" />
                </Button>
              </div>
            ))}
          </div>
        )}
      </CardContent>
    </Card>
  );
}
