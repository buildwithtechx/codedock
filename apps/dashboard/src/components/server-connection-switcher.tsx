import { Check, RefreshCw, Server } from 'lucide-react';
import { useState } from 'react';
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
import { getApiBaseUrl, setApiBaseUrl } from '#/lib/api-client';

export function ServerConnectionSwitcher() {
  const [open, setOpen] = useState(false);
  const [customUrl, setCustomUrl] = useState(getApiBaseUrl() || 'http://localhost:8080');

  const handleSave = () => {
    setApiBaseUrl(customUrl.trim());
    window.location.reload();
  };

  const handleResetLocal = () => {
    setApiBaseUrl('http://localhost:8080');
    window.location.reload();
  };

  const currentUrl = getApiBaseUrl() || 'http://localhost:8080';

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button variant="outline" size="sm" className="h-8 gap-2 border-border/50 text-xs">
          <Server className="h-3.5 w-3.5 text-blue-400" />
          <span className="max-w-30 truncate font-mono">{currentUrl}</span>
        </Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-106.25">
        <DialogHeader>
          <DialogTitle>Server Connection</DialogTitle>
          <DialogDescription>
            Switch between a remote Codedock VPS control plane and your local daemon.
          </DialogDescription>
        </DialogHeader>
        <div className="grid gap-4 py-4">
          <div className="grid gap-2">
            <Label htmlFor="server-url">Control Plane URL</Label>
            <Input
              id="server-url"
              value={customUrl}
              onChange={(e) => setCustomUrl(e.target.value)}
              placeholder="https://codedock.example.com"
            />
          </div>
          <div className="flex gap-2">
            <Button variant="outline" size="sm" className="flex-1 gap-2" onClick={handleResetLocal}>
              <RefreshCw className="h-3.5 w-3.5" />
              Reset Local (localhost:8080)
            </Button>
            <Button size="sm" className="gap-2" onClick={handleSave}>
              <Check className="h-3.5 w-3.5" />
              Connect
            </Button>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
}
