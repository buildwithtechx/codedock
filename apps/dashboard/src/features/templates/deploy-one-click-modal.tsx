import { Loader2 } from 'lucide-react';
import { useEffect, useMemo, useState } from 'react';
import { toast } from 'sonner';
import { Badge } from '#/components/ui/badge';
import { Button } from '#/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '#/components/ui/dialog';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import { QueryErrorState } from '#/components/ui/query-error-state';
import {
  useDeployOneClickApp,
  useOneClickApp,
  useReviewOneClickInstall,
} from '#/hooks/use-templates';
import type { InstallPreview, OneClickEnvVar } from '#/interfaces/templates';

interface Props {
  isOpen: boolean;
  onOpenChange: (open: boolean) => void;
  projectId: string;
  appId: string;
  appName: string;
  onDeployed: (appName: string) => void;
}

export function DeployOneClickModal({
  isOpen,
  onOpenChange,
  projectId,
  appId,
  appName,
  onDeployed,
}: Props) {
  return (
    <Dialog open={isOpen} onOpenChange={onOpenChange}>
      <DialogContent className="max-h-[85vh] overflow-y-auto sm:max-w-lg">
        <DialogHeader>
          <DialogTitle>Deploy {appName}</DialogTitle>
          <DialogDescription>
            Configure install settings, review the plan, then deploy.
          </DialogDescription>
        </DialogHeader>
        {isOpen && appId ? (
          <DeployOneClickBody
            projectId={projectId}
            appId={appId}
            appName={appName}
            onDone={() => {
              onOpenChange(false);
              onDeployed(appName);
            }}
          />
        ) : null}
      </DialogContent>
    </Dialog>
  );
}

function DeployOneClickBody({
  projectId,
  appId,
  appName,
  onDone,
}: {
  projectId: string;
  appId: string;
  appName: string;
  onDone: () => void;
}) {
  const { data: app, isLoading, isError, refetch } = useOneClickApp(appId);
  const reviewInstall = useReviewOneClickInstall();
  const deployApp = useDeployOneClickApp();
  const [name, setName] = useState(appName);
  const [values, setValues] = useState<Record<string, string>>({});
  const [hostPort, setHostPort] = useState('');
  const [domain, setDomain] = useState('');
  const [preview, setPreview] = useState<InstallPreview | null>(null);

  useEffect(() => {
    setName(appName);
    setValues({});
    setHostPort('');
    setDomain('');
    setPreview(null);
  }, [appId, appName]);

  const variables = useMemo(() => app?.envVariables ?? [], [app]);
  const secrets = variables.filter((variable) => variable.secret);
  const inputs = variables.filter((variable) => !variable.secret && variable.input);
  const fixed = variables.filter((variable) => !variable.secret && !variable.input);

  const setValue = (key: string, value: string) =>
    setValues((current) => ({ ...current, [key]: value }));

  const buildPayload = (digest?: string) => {
    const secretsPayload: Record<string, string> = {};
    const environmentPayload: Record<string, string> = {};
    for (const variable of variables) {
      const value = (values[variable.key] ?? '').trim();
      if (!value) {
        continue;
      }
      if (variable.secret) {
        secretsPayload[variable.key] = value;
      } else {
        environmentPayload[variable.key] = value;
      }
    }
    const port = Number.parseInt(hostPort, 10);
    return {
      appId,
      projectId,
      name: name.trim() || appName,
      secrets: secretsPayload,
      environment: environmentPayload,
      hostPort: Number.isNaN(port) ? undefined : port,
      domain: domain.trim() || undefined,
      digest,
    };
  };

  const configurationValid = () => {
    if (!name.trim()) {
      toast.error('Install name is required');
      return false;
    }
    for (const variable of variables) {
      if (variable.required && !(values[variable.key] ?? '').trim() && !variable.defaultValue) {
        toast.error(`${variable.label || variable.key} is required`);
        return false;
      }
    }
    if (hostPort.trim()) {
      const port = Number.parseInt(hostPort, 10);
      if (Number.isNaN(port) || port < 1 || port > 65535) {
        toast.error('Host port must be between 1 and 65535');
        return false;
      }
    }
    return true;
  };

  const handleReview = async () => {
    if (!configurationValid()) {
      return;
    }
    try {
      setPreview(await reviewInstall.mutateAsync(buildPayload()));
    } catch {
      toast.error('Install review failed');
    }
  };

  const handleDeploy = async () => {
    if (!preview) {
      return;
    }
    try {
      const response = await deployApp.mutateAsync(buildPayload(preview.digest));
      toast.success(response.message || `${appName} deployed`);
      onDone();
    } catch {
      toast.error(`Failed to deploy ${appName}`);
    }
  };

  if (isLoading) {
    return (
      <div className="flex justify-center p-12">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  if (isError || !app) {
    return (
      <QueryErrorState
        title="App details are unavailable"
        description="Codedock could not load this catalogue entry."
        onRetry={() => void refetch()}
      />
    );
  }

  if (preview) {
    return (
      <div className="space-y-4">
        <div className="flex flex-wrap gap-2">
          {preview.services.map((service) => (
            <Badge key={service.service} variant="secondary">
              {service.service} · {service.image}
            </Badge>
          ))}
        </div>
        {preview.generatedSecrets.length > 0 && (
          <p className="text-muted-foreground text-sm">
            Generated secrets: {preview.generatedSecrets.join(', ')}
          </p>
        )}
        <pre className="max-h-72 overflow-auto rounded-md border bg-muted p-3 font-mono text-xs">
          {preview.composeYaml}
        </pre>
        <div className="flex justify-end gap-2">
          <Button variant="outline" onClick={() => setPreview(null)} disabled={deployApp.isPending}>
            Back
          </Button>
          <Button onClick={() => void handleDeploy()} disabled={deployApp.isPending}>
            {deployApp.isPending ? 'Deploying...' : 'Deploy'}
          </Button>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      <div className="space-y-2">
        <Label htmlFor="one-click-name">Install name</Label>
        <Input id="one-click-name" value={name} onChange={(event) => setName(event.target.value)} />
      </div>

      {inputs.map((variable) => (
        <VariableField
          key={variable.key}
          variable={variable}
          value={values[variable.key] ?? ''}
          onChange={setValue}
        />
      ))}

      {secrets.length > 0 && (
        <div className="space-y-4">
          <p className="font-medium text-sm">Secrets</p>
          {secrets.map((variable) => (
            <VariableField
              key={variable.key}
              variable={variable}
              value={values[variable.key] ?? ''}
              onChange={setValue}
            />
          ))}
        </div>
      )}

      {fixed.length > 0 && (
        <div className="space-y-4">
          <p className="font-medium text-sm">Environment overrides</p>
          {fixed.map((variable) => (
            <VariableField
              key={variable.key}
              variable={variable}
              value={values[variable.key] ?? ''}
              onChange={setValue}
            />
          ))}
        </div>
      )}

      <div className="grid grid-cols-2 gap-4">
        <div className="space-y-2">
          <Label htmlFor="one-click-port">Host port</Label>
          <Input
            id="one-click-port"
            inputMode="numeric"
            placeholder={String(app.defaultPort)}
            value={hostPort}
            onChange={(event) => setHostPort(event.target.value)}
          />
        </div>
        <div className="space-y-2">
          <Label htmlFor="one-click-domain">Domain</Label>
          <Input
            id="one-click-domain"
            placeholder="app.example.com"
            value={domain}
            onChange={(event) => setDomain(event.target.value)}
          />
        </div>
      </div>

      <div className="flex justify-end gap-2">
        <Button onClick={handleReview} disabled={reviewInstall.isPending}>
          {reviewInstall.isPending ? 'Reviewing...' : 'Review plan'}
        </Button>
      </div>
    </div>
  );
}

function VariableField({
  variable,
  value,
  onChange,
}: {
  variable: OneClickEnvVar;
  value: string;
  onChange: (key: string, value: string) => void;
}) {
  const placeholder = variable.secret ? 'Leave blank to generate' : (variable.defaultValue ?? '');
  return (
    <div className="space-y-2">
      <Label htmlFor={`one-click-${variable.key}`}>
        {variable.label || variable.key}
        {variable.required && !variable.defaultValue ? ' *' : ''}
      </Label>
      <Input
        id={`one-click-${variable.key}`}
        type={variable.secret ? 'password' : 'text'}
        autoComplete="off"
        placeholder={placeholder}
        value={value}
        onChange={(event) => onChange(variable.key, event.target.value)}
      />
    </div>
  );
}
